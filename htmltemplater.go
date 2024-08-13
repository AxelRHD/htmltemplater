package htmltemplater

import (
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"path"
	"path/filepath"
)

type HtmlTemplate struct {
	*template.Template
	Templater *Templater
}

type TemplaterOptions struct {
	ImportPath    string
	FileExtension string
	Layout        []string
}

type Templater struct {
	TemplaterOptions
}

var (
	globalTemplater *Templater
	ErrTmplNotFound = errors.New("template not found")
	ErrTmplNil      = errors.New("template nil / undefined")
)

func init() {
	globalTemplater = NewTemplater(nil)
}

func NewTemplater(opts *TemplaterOptions) *Templater {
	o := TemplaterOptions{
		ImportPath: "templates",
		// FileExtension: ".tmpl.html",
		FileExtension: ".gotmpl",
		Layout:        []string{"_layout"},
	}

	if opts != nil {
		if opts.ImportPath != "" {
			o.ImportPath = opts.ImportPath
		}

		if opts.FileExtension != "" {
			o.FileExtension = opts.FileExtension
		}
	}

	return &Templater{
		TemplaterOptions: o,
	}
}

func SetImportPath(pth string) {
	globalTemplater.ImportPath = pth
}

func GetImportPath() string {
	return globalTemplater.ImportPath
}

func SetFileExtension(ext string) {
	globalTemplater.FileExtension = ext
}

func GetFileExtension() string {
	return globalTemplater.FileExtension
}

func SetLayout(tmpls ...string) {
	globalTemplater.Layout = tmpls
}

func GetLayout() []string {
	return globalTemplater.Layout
}

func (tpltr *Templater) NewTemplate(withLayout, isPath bool, paths ...string) (*HtmlTemplate, error) {
	var tmpl *template.Template
	var inp []string

	if withLayout {
		for _, v := range tpltr.Layout {
			inp = append(inp, tpltr.GenerateTmplPath(v))
		}
	}

	for _, v := range paths {
		if isPath {
			inp = append(inp, v)
		} else {
			inp = append(inp, tpltr.GenerateTmplPath(v))
		}
	}

	tmpl, err := template.ParseFiles(inp...)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("NewTemplate", slog.Any("tpltr", tpltr))
	return &HtmlTemplate{
		Template:  tmpl,
		Templater: tpltr,
	}, nil
}

func NewTemplate(withLayout, isPath bool, paths ...string) (*HtmlTemplate, error) {
	return globalTemplater.NewTemplate(withLayout, isPath, paths...)
}

func (tpltr *Templater) NewTemplateWithFuncs(fncs template.FuncMap, withLayout, isPath bool, paths ...string) (*HtmlTemplate, error) {
	var tmpl *template.Template
	var inp []string
	var tname string

	if withLayout {
		for i, v := range tpltr.Layout {
			if i == 0 {
				tname = tpltr.GenerateFilename(v)
			}
			inp = append(inp, tpltr.GenerateTmplPath(v))
		}
	}

	for i, v := range paths {
		if i == 0 && tname == "" {
			if isPath {
				tname = v
			} else {
				tname = tpltr.GenerateFilename(v)
			}

		}

		if isPath {
			inp = append(inp, v)
		} else {
			inp = append(inp, tpltr.GenerateTmplPath(v))
		}
	}

	tmpl = template.New(filepath.Base(tname)).Funcs(fncs)
	_, err := tmpl.ParseFiles(inp...)
	if err != nil {
		return nil, err
	}

	return &HtmlTemplate{
		Template:  tmpl,
		Templater: tpltr,
	}, nil
}

func NewTemplateWithFuncs(fncs template.FuncMap, withLayout, isPath bool, paths ...string) (*HtmlTemplate, error) {
	return globalTemplater.NewTemplateWithFuncs(fncs, withLayout, isPath, paths...)
}

func (tpltr *Templater) GenerateTmplPaths(tmpls ...string) []string {
	var out []string

	for _, t := range tmpls {
		fmt.Println(t)
		out = append(out, tpltr.GenerateTmplPath(t))
	}

	return out
}

func GenerateTmplPaths(tmpls ...string) []string {
	return globalTemplater.GenerateTmplPaths(tmpls...)
}

func (tpltr *Templater) GenerateTmplPath(tmpl string) string {
	return path.Join(tpltr.ImportPath, tpltr.GenerateFilename(tmpl))
}

func GenerateTmplPath(tmpl string) string {
	return globalTemplater.GenerateTmplPath(tmpl)
}

func (tpltr *Templater) GenerateFilename(tmpl string) string {
	return fmt.Sprintf("%s%s", tmpl, tpltr.FileExtension)
}

func GenerateFilename(tmpl string) string {
	return globalTemplater.GenerateFilename(tmpl)
}

func (tmpl *HtmlTemplate) Parse(tmpls ...string) (*HtmlTemplate, error) {
	t, err := tmpl.ParseFiles(tmpl.Templater.GenerateTmplPaths(tmpls...)...)
	if err != nil {
		return tmpl, err
	}

	tmpl.Template = t
	return tmpl, nil
}
