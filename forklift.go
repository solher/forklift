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

import "github.com/solher/forklift/files"

func init() {
{{range $path, $file := .Files }}  files.Add("{{$path}}", {{$file}})
{{end -}}
}
`

var tmpl = template.Must(template.New("forklift").Parse(tmplStr))

type forklift struct {
	Package    string
	Files      map[string]string
	Extensions []string
}

func main() {
	pkg := flag.String("package", "main", "package where the forklift file is to be written")
	dir := flag.String("directory", "./", "the root directory where to look for files")
	extensions := flag.String("extensions", "", "a comma separated list of extensions to load")
	flag.Parse()

	f := &forklift{
		Package:    *pkg,
		Files:      map[string]string{},
		Extensions: []string{},
	}

	for _, extension := range strings.Split(*extensions, ",") {
		f.Extensions = append(f.Extensions, strings.TrimSpace(extension))
	}

	filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if len(f.Extensions) > 0 {
			found := false
			for _, extension := range f.Extensions {
				if strings.HasSuffix(path, extension) {
					found = true
				}
			}
			if !found {
				return nil
			}
		}
		if info != nil && info.IsDir() {
			return nil
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		file, err := ioutil.ReadFile(abs)
		if err != nil {
			panic(err)
		}
		f.Files[abs] = fmt.Sprintf("%q", file)
		return nil
	})

	if err := tmpl.Execute(os.Stdout, f); err != nil {
		panic(err)
	}
}
