package main

import (
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/axelrhd/htmltemplater"
)

var (
	Tmpl *htmltemplater.HtmlTemplater
)

func main() {
	Tmpl = htmltemplater.New()
	Tmpl.StdFuncs = &template.FuncMap{
		"cap": strings.ToLower,
	}

	layout, err := Tmpl.GenerateLayout()
	logFatal(err)

	index, err := layout.ParseFiles(Tmpl.GenerateTmplPaths("page/index")...)
	logFatal(err)

	err = index.Execute(os.Stdout, map[string]any{
		"User": "aXeL",
	})
	logFatal(err)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
