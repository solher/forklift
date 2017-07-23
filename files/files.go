package files

import (
	"bytes"
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
func Add(path string, file string) {
	files[path] = &parsedFile{
		file:     file,
		template: template.Must(template.New(path).Parse(file)),
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
