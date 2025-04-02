package files

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"text/template"
)

var files = map[string]string{}
var templates = template.New("root")

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
	if file, ok := files[path]; ok {
		return file
	}
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(file)
}

// AbsTemplate reads the file located at the given absolute path and parses it.
func AbsTemplate(path string, data interface{}) string {
	output := bytes.NewBuffer(nil)

	// Try to use existing template.
	tmpl := templates.Lookup(path)
	if tmpl == nil {
		// If not found, try to read from disk.
		content, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		// Process includes and parse the template.
		processedContent := processIncludes(string(content), path)
		tmpl = template.Must(template.New(path).Parse(processedContent))
	}

	if err := tmpl.Execute(output, data); err != nil {
		panic(err)
	}
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
		unzippedFile, err := io.ReadAll(gz)
		if err != nil {
			panic(err)
		}
		clearFile = string(unzippedFile)
	}

	// Store the raw file.
	files[path] = clearFile
}

// LoadAllTemplates processes all files and adds them to the template set.
func LoadAllTemplates() {
	processedFiles := map[string]string{}

	for path, content := range files {
		// Process includes.
		processedContent := processIncludes(content, path)

		// Store the processed file.
		processedFiles[path] = processedContent

		// Add to template set.
		template.Must(templates.New(path).Parse(processedContent))
	}

	// Replace the original files with the processed ones.
	files = processedFiles
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

// includeRegex is the regex to find {{include "path"}} patterns.
var includeRegex = regexp.MustCompile(`{{\s*include\s*"([^"]+)"\s*}}`)

// processIncludes replaces all {{include "path"}} directives with the content of the referenced files.
func processIncludes(content, basePath string) string {
	return includeRegex.ReplaceAllStringFunc(content, func(match string) string {
		submatches := includeRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		includePath := submatches[1]
		if !filepath.IsAbs(includePath) {
			includePath = filepath.Join(filepath.Dir(basePath), includePath)
		}

		return AbsFile(includePath)
	})
}
