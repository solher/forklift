package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const tmplStr = `package {{.Package}}

import "github.com/solher/forklift/queries"

func init() {
{{range $path, $query := .Queries}}queries.Add("{{$path}}", {{$query}})
{{end}}
}
`

var tmpl = template.Must(template.New("forklift").Parse(tmplStr))

type forklift struct {
	Package string
	Queries map[string]string
}

func main() {
	pkg := flag.String("package", "main", "package where the query file is to be written")
	dir := flag.String("directory", "./", "the root directory where to look for sql files")
	flag.Parse()

	f := &forklift{
		Package: *pkg,
		Queries: map[string]string{},
	}

	filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if !(strings.HasSuffix(path, ".static.sql") || strings.HasSuffix(path, ".tmpl.sql") || strings.HasSuffix(path, ".lazy.sql")) {
			return nil
		}
		if info != nil && info.IsDir() {
			return nil
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		query, err := ioutil.ReadFile(abs)
		if err != nil {
			panic(err)
		}
		f.Queries[abs] = fmt.Sprintf("%q", query)
		return nil
	})

	if err := tmpl.Execute(os.Stdout, f); err != nil {
		panic(err)
	}
}
