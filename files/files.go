package files

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"text/template"
)

type parsedFile struct {
	file     string
	template *template.Template
}

var files = map[string]*parsedFile{}

// File reads the given file located relatively to the caller.
func File(path string) string {
	return AbsFile(absFromCaller(path))
}

// Template reads the given file located relatively to the caller and parses it.
func Template(path string, data interface{}) string {
	return AbsTemplate(absFromCaller(path), data)
}

// AbsFile reads the file located at the given absolute path.
func AbsFile(path string) string {
	if q, ok := files[path]; ok {
		return q.file
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(file)
}

// AbsTemplate reads the file located at the given absolute path and parses it.
func AbsTemplate(path string, data interface{}) string {
	output := bytes.NewBuffer(nil)
	if q, ok := files[path]; ok {
		q.template.Execute(output, data)
		return output.String()
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New(path).Parse(string(file))
	if err != nil {
		panic(err)
	}
	tmpl.Execute(output, data)
	return output.String()
}

// Add adds a new file to the cache and tries to parse it.
// The file is expected to be gzipped then base64 encoded.
func Add(path string, base64File string) {
	clearFile := base64File
	decodedFile, err := base64.StdEncoding.DecodeString(base64File)
	if err == nil {
		gz, err := gzip.NewReader(bytes.NewBuffer(decodedFile))
		if err != nil {
			panic(err)
		}
		defer gz.Close()
		unzippedFile, err := ioutil.ReadAll(gz)
		if err != nil {
			panic(err)
		}
		clearFile = string(unzippedFile)
	}
	files[path] = &parsedFile{
		file:     clearFile,
		template: template.Must(template.New(path).Parse(clearFile)),
	}
}

func absFromCaller(path string) string {
	_, f, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	abs, err := filepath.Abs(f)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s/%s", filepath.Dir(abs), path)
}
