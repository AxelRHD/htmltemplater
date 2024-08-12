package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axelrhd/htmltemplater"
)

var (
	Tmpl *htmltemplater.HtmlTemplater
)

func main() {
	// ht := htmltemplater.New(&htmltemplater.TemplaterOptions{
	// 	ImportPath: "./templates",
	// })
	Tmpl = htmltemplater.New(nil)

	t, err := Tmpl.AddPage("page/index", &[]string{"_layout", "partial/navbar"}, nil)
	logFatal(err)

	p, err := Tmpl.AddPartial("partial/navbar", nil)
	logFatal(err)

	err = t.Execute(os.Stdout, nil)
	logFatal(err)

	fmt.Printf("%+v\n", Tmpl)

	err = p.ExecuteTemplate(os.Stdout, "navbar", nil)
	logFatal(err)

}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
