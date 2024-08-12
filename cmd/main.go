package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axelrhd/htmltemplater"
)

func main() {
	// ht := htmltemplater.New(&htmltemplater.TemplaterOptions{
	// 	ImportPath: "./templates",
	// })
	ht := htmltemplater.New(nil)

	fmt.Printf("%+v\n", ht)
	ht.ReadFiles()

	for _, t := range ht.Templates {
		fmt.Println(t)
	}

	fmt.Printf("%+v\n", ht)

	tp, err := ht.PageTmpl("index", "partials/navbar")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tp.Name())
	fmt.Println(tp.DefinedTemplates())

	err = tp.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatal(err)
	}
}
