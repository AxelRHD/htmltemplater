package htmltemplater

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	DefaultTemplaterOptions = TemplaterOptions{
		ImportPath:     "templates",
		FileExtension:  "tmpl.html",
		LayoutTemplate: "_layout",
	}

	ErrTmplNotFound = errors.New("template not found")
)

type TemplaterOptions struct {
	ImportPath     string
	FileExtension  string
	LayoutTemplate string
}

type HtmlTemplate struct {
	Template *template.Template
	Name     string
}

func (ht HtmlTemplate) String() string {
	var tpls []string
	for _, v := range ht.Template.Templates() {
		tpls = append(tpls, v.Tree.Name)
	}
	return fmt.Sprintf("%s: %+v", ht.Name, tpls)
}

type HtmlTemplater struct {
	Templates      []*HtmlTemplate
	ImportPath     string
	FileExtension  string
	LayoutTemplate string
}

func New(opts *TemplaterOptions) *HtmlTemplater {
	var topts TemplaterOptions

	if opts != nil {
		topts = *opts
	}

	if topts.ImportPath == "" {
		topts.ImportPath = DefaultTemplaterOptions.ImportPath
	}

	if topts.FileExtension == "" {
		topts.FileExtension = DefaultTemplaterOptions.FileExtension
	}

	if topts.LayoutTemplate == "" {
		topts.LayoutTemplate = DefaultTemplaterOptions.LayoutTemplate
	}

	return &HtmlTemplater{
		ImportPath:     topts.ImportPath,
		FileExtension:  topts.FileExtension,
		LayoutTemplate: topts.LayoutTemplate,
	}
}

func (ht *HtmlTemplater) ReadFiles() error {
	fsys := os.DirFS(".")
	return fs.WalkDir(fsys, ht.ImportPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if strings.Contains(path, ht.FileExtension) {
				parent := filepath.Base(filepath.Dir(path))

				suffix := strings.TrimPrefix(ht.FileExtension, ".")
				suffix = fmt.Sprintf(".%s", suffix)

				tmpl := strings.TrimSuffix(filepath.Base(path), suffix)
				if parent != ht.ImportPath {
					tmpl = filepath.Join(parent, tmpl)
				}

				t, err := template.ParseFiles(path)
				if err != nil {
					return err
				}

				ht.Templates = append(ht.Templates, &HtmlTemplate{
					Name:     tmpl,
					Template: t,
				})
			}
		}

		return nil
	})
}

func (ht *HtmlTemplater) GetTmplNames() []string {
	var out []string

	for _, v := range ht.Templates {
		out = append(out, v.Name)
	}

	return out
}

func (ht *HtmlTemplater) GetTmpl(name string) (*template.Template, error) {
	for _, v := range ht.Templates {
		if v.Name == name {
			return v.Template, nil
		}
	}

	return nil, ErrTmplNotFound
}

func (ht *HtmlTemplater) PageTmpl(child string, others ...string) (*template.Template, error) {
	var tmplLayout, tmplChild *template.Template

	tmplLayout, err := ht.GetTmpl(ht.LayoutTemplate)
	if err != nil {
		return nil, err
	}

	tmplChild, err = ht.GetTmpl(child)
	if err != nil {
		return nil, err
	}

	for _, v := range others {
		tmpl, err := ht.GetTmpl(v)
		if err != nil {
			return nil, err
		}

		// tmplChild, err = tmplChild.AddParseTree(tmpl.Name(), tmpl.Tree)
		tmplLayout, err = tmplLayout.AddParseTree(tmpl.Name(), tmpl.Tree)
		if err != nil {
			return nil, err
		}
	}

	// return tmplChild.AddParseTree(ht.LayoutTemplate, tmplLayout.Tree)
	return tmplLayout.AddParseTree(tmplChild.Tree.ParseName, tmplChild.Tree)
}
