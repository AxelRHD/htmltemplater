package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"strings"

	"github.com/axelrhd/htmltemplater"
)

var (
	Tmpl htmltemplater.Templater
)

func main() {
	Tmpl = *htmltemplater.New()
	Tmpl.ImportPath = "../templates"
	Tmpl.StdFuncs = &template.FuncMap{
		"cap": strings.ToTitle,
	}

	fmt.Printf("%+v\n", Tmpl)

	layout, err := Tmpl.GenerateLayout()
	logFatal(err)
	fmt.Println(layout.Tree)

	index, err := layout.ParseFiles(Tmpl.GenerateTmplPaths("page/index")...)
	logFatal(err)

	fmt.Printf("%+v\n", layout.DefinedTemplates())
	fmt.Printf("%+v\n", index.DefinedTemplates())
	fmt.Printf("%+v\n", Tmpl)
}

func logFatal(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
