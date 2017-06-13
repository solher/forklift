package queries

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

type parsedQuery struct {
	query    string
	template *template.Template
}

var queries = map[string]*parsedQuery{}

// File reads the given sql file located relatively to the caller.
func File(path string) string {
	return AbsFile(absFromCaller(path))
}

// Template reads the given sql file located relatively to the caller and parses it.
func Template(path string) *template.Template {
	return AbsTemplate(absFromCaller(path))
}

// AbsFile reads the sql file located at the given absolute path.
func AbsFile(path string) string {
	if q, ok := queries[path]; ok {
		return q.query
	}
	query, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(query)
}

// AbsTemplate reads the sql file located at the given absolute path and parses it.
func AbsTemplate(path string) *template.Template {
	if q, ok := queries[path]; ok {
		return q.template
	}
	query, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New(path).Parse(string(query))
	if err != nil {
		panic(err)
	}
	return tmpl
}

// Add adds a new sql file to the cache and tries to parse it.
func Add(path string, query string) {
	queries[path] = &parsedQuery{
		query:    query,
		template: template.Must(template.New(path).Parse(query)),
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
