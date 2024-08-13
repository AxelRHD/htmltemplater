package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/url"
	"os"

	tpl "github.com/axelrhd/htmltemplater"
)

var (
	Prefix = "/app"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	tpl.SetImportPath("../templates")

	stdFuncs := template.FuncMap{
		"lnk": func(p string) string {
			u, err := url.JoinPath(Prefix, p)
			if err != nil {
				return p
			}

			return u
		},
	}

	index, err := tpl.NewTemplateWithFuncs(stdFuncs, true, false, "page/index")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(index.DefinedTemplates())

	err = index.Execute(os.Stdout, nil)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
