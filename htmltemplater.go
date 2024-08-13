package htmltemplater

import (
	"errors"
	"fmt"
	"html/template"
	"path"
)

type HtmlTemplate struct {
	Template *template.Template
	Name     string
}

type Templater struct {
	Templates []*HtmlTemplate
	// Templates     []*template.Template
	ImportPath    string
	FileExtension string
	LayoutTmpls   []string
	Layout        *template.Template
	StdFuncs      *template.FuncMap
}

var (
	ErrTmplNotFound = errors.New("template not found")
)

func New() *Templater {
	return &Templater{
		ImportPath:    "templates",
		FileExtension: ".tmpl.html",
		// FileExtension: ".gotmpl",
		LayoutTmpls: []string{"_layout"},
	}
}

// func (ht *Templater) NewFromLayout() (*template.Template,error){

// }

func (ht *Templater) GenerateLayout() (*template.Template, error) {
	var tmpl *template.Template
	var err error
	for i, l := range ht.LayoutTmpls {
		p := ht.GenerateTmplPath(l)

		if i == 0 {
			tmpl = template.New(ht.GenerateFilename(l))

			if ht.StdFuncs != nil {
				tmpl = template.New(ht.GenerateFilename(l)).Funcs(*ht.StdFuncs)
				// tmpl = tmpl.Funcs(*ht.StdFuncs)
				// tmpl.Funcs(*ht.StdFuncs)
			} else {
				tmpl = template.New(ht.GenerateFilename(l))
			}

		}

		tmpl, err = template.ParseFiles(p)
		if err != nil {
			return nil, err
		}
	}

	ht.Layout = tmpl
	return tmpl, err
}

func (ht *Templater) GenerateTmplPaths(tmpls ...string) []string {
	var out []string

	for _, t := range tmpls {
		out = append(out, ht.GenerateTmplPath(t))
	}

	return out
}

func (ht *Templater) GenerateTmplPath(tmpl string) string {

	return path.Join(ht.ImportPath, ht.GenerateFilename(tmpl))
}

func (ht *Templater) GenerateFilename(tmpl string) string {
	return fmt.Sprintf("%s%s", tmpl, ht.FileExtension)
}

func (ht *Templater) GetTmplNames() []string {
	var out []string

	for _, t := range ht.Templates {
		out = append(out, t.Name)
	}

	return out
}

func (ht *Templater) GetTmpl(name string) (*template.Template, error) {
	for _, t := range ht.Templates {
		if t.Name == name {
			return t.Template, nil
		}
	}

	return nil, ErrTmplNotFound
}
