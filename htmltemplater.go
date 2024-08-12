package htmltemplater

import (
	"errors"
	"fmt"
	"html/template"
	"path"
	"path/filepath"
)

// type TemplateType int

// type TemplaterOptions struct {
// 	ImportPath    string
// 	FileExtension string
// 	Layout        []string
// 	StdFuncs      *template.FuncMap
// }

type HtmlTemplate struct {
	Template *template.Template
	Name     string
	// Type     TemplateType
}

type HtmlTemplater struct {
	Templates []*HtmlTemplate
	// Templates     []*template.Template
	ImportPath    string
	FileExtension string
	LayoutTmpls   []string
	Layout        *template.Template
	StdFuncs      *template.FuncMap
}

// const (
// TMPL_PAGE TemplateType = iota
// TMPL_PART
// )

var (
	// DefaultTemplaterOptions = TemplaterOptions{
	// 	ImportPath:    "templates",
	// 	FileExtension: ".tmpl.html",
	// 	// FileExtension: ".gotmpl",
	// 	LayoutTmpls: []string{"_layout"},
	// }

	ErrTmplNotFound = errors.New("template not found")
)

func New() *HtmlTemplater {
	return &HtmlTemplater{
		ImportPath:    "templates",
		FileExtension: ".tmpl.html",
		// FileExtension: ".gotmpl",
		LayoutTmpls: []string{"_layout"},
	}
}

func (ht *HtmlTemplater) GenerateLayout() (*template.Template, error) {
	var inp []string

	for _, l := range ht.LayoutTmpls {
		inp = append(inp, generateFsPath(l, ht.ImportPath, ht.FileExtension))
	}

	base := filepath.Base(inp[0])

	var tmpl *template.Template
	var err error
	if ht.StdFuncs != nil {
		tmpl, err = template.New(base).Funcs(*ht.StdFuncs).ParseFiles(inp...)
	} else {
		tmpl, err = template.New(base).ParseFiles(inp...)
	}
	if err != nil {
		return nil, err
	}
	ht.Layout = tmpl

	return tmpl, nil
}

func (ht *HtmlTemplater) GenerateTmplPaths(tmpls ...string) []string {
	var out []string

	for _, t := range tmpls {
		out = append(out, generateFsPath(t, ht.ImportPath, ht.FileExtension))
	}

	return out
}

// func (ht HtmlTemplate) String() string {
// 	var tpls []string
// 	for _, v := range ht.Template.Templates() {
// 		tpls = append(tpls, v.Tree.Name)
// 	}
// 	return fmt.Sprintf("%s: %+v", ht.Name, tpls)
// }

// func New(opts *TemplaterOptions) *HtmlTemplater {
// 	var topts TemplaterOptions

// 	if opts != nil {
// 		topts = *opts
// 	}

// 	if topts.ImportPath == "" {
// 		topts.ImportPath = DefaultTemplaterOptions.ImportPath
// 	}

// 	if topts.FileExtension == "" {
// 		topts.FileExtension = DefaultTemplaterOptions.FileExtension
// 	}

// 	if len(topts.Layout) == 0 {
// 		topts.Layout = DefaultTemplaterOptions.Layout
// 	}

// 	return &HtmlTemplater{
// 		ImportPath:    topts.ImportPath,
// 		FileExtension: topts.FileExtension,
// 		Layout:        topts.Layout,
// 		StdFuncs:      topts.StdFuncs,
// 	}
// }

// AddTemplate adds a template to the templater with standard funcs
// func (ht *HtmlTemplater) AddTemplate(name string, withLayout bool, tmpls ...string) (*template.Template, error) {
// 	var err error
// 	var inp []string
// 	var tmpl *template.Template

// 	if withLayout {
// 		for _, l := range ht.Layout {
// 			inp = append(inp, generateFsPath(l, ht.ImportPath, ht.FileExtension))
// 		}
// 	}

// 	for _, t := range tmpls {
// 		inp = append(inp, generateFsPath(t, ht.ImportPath, ht.FileExtension))
// 	}

// 	base := filepath.Base(inp[0])
// 	if ht.StdFuncs != nil {
// 		tmpl, err = template.New(base).Funcs(*ht.StdFuncs).ParseFiles(inp...)
// 	} else {
// 		tmpl, err = template.New(base).ParseFiles(inp...)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	// ht.Templates = append(ht.Templates, tmpl)
// 	ht.Templates = append(ht.Templates, &HtmlTemplate{
// 		Template: tmpl,
// 		Name:     name,
// 	})

// 	return tmpl, nil
// }

// func (ht *HtmlTemplater) GenerateTmplPaths(pths ...string) []string {
// 	var out []string

// 	for _, p := range pths {
// 		out = append(out, generateFsPath(p, ht.ImportPath, ht.FileExtension))
// 	}

// 	return out
// }

func (ht *HtmlTemplater) GetTmplNames() []string {
	var out []string

	for _, t := range ht.Templates {
		out = append(out, t.Name)
	}

	return out
}

func (ht *HtmlTemplater) GetTmpl(name string) (*template.Template, error) {
	for _, t := range ht.Templates {
		if t.Name == name {
			return t.Template, nil
		}
	}

	return nil, ErrTmplNotFound
}

func generateFsPath(tmpl, root, ext string) string {
	// return filepath.Join(ht.ImportPath, fmt.Sprintf("%s%s", tmpl, ht.FileExtension))
	// return path.Join(ht.ImportPath, fmt.Sprintf("%s%s", tmpl, ht.FileExtension))
	// fmt.Printf("%v%v\n", tmpl, ext)
	// fmt.Println("tmpl: ", tmpl)
	// fmt.Println("root: ", root)
	// fmt.Println("ext: ", ext)
	return path.Join(root, fmt.Sprintf("%v%v", tmpl, ext))
	// return fmt.Sprintf("%s/%s%s", root, tmpl, ext)
}

// func (ht *HtmlTemplater) addTemplate(tmpl string, tplType TemplateType, fn *template.FuncMap, parts ...string) (*template.Template, error) {
// 	var err error
// 	var tmpls []string

// 	if tplType == TMPL_PART {
// 		tmpls = append(tmpls, tmpl)
// 	}

// 	for _, p := range parts {
// 		tmpls = append(tmpls, p)
// 	}

// 	if tplType == TMPL_PAGE {
// 		tmpls = append(tmpls, tmpl)
// 	}

// 	var pths []string
// 	for _, tpl := range tmpls {
// 		pths = append(pths, generateTemplatePath(tpl, ht.ImportPath, ht.FileExtension))
// 	}

// 	fmt.Printf("tmpls: %+v\n", tmpls)
// 	fmt.Printf("pths: %+v\n", pths)

// 	fm := ht.StdFuncs
// 	if fn != nil {
// 		fm = fn
// 	}

// 	var parsed *template.Template
// 	if fm != nil {

// 		// parsed, err = template.Funcs(*fm).ParseFiles(pths...)
// 	} else {
// 		parsed, err = template.ParseFiles(pths...)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	ht.Templates = append(ht.Templates, &HtmlTemplate{
// 		Template: parsed,
// 		Name:     tmpl,
// 		Type:     tplType,
// 	})

// 	return parsed, nil
// }

// func (ht *HtmlTemplater) AddTemplate(name strimng)

// func (ht *HtmlTemplater) AddPage(tmpl string, layout, parts *[]string, funcmp *template.FuncMap) (*template.Template, error) {
// 	var tmpls []string

// 	if layout != nil {
// 		for _, l := range *layout {
// 			tmpls = append(tmpls, l)
// 		}
// 	} else {
// 		tmpls = ht.Layout
// 	}

// 	if parts != nil {
// 		for _, p := range *parts {
// 			tmpls = append(tmpls, p)
// 		}
// 	}

// 	return ht.addTemplate(tmpl, TMPL_PAGE, fm, tmpls...)
// }

// func (ht *HtmlTemplater) AddPartial(tmpl string, parts *[]string, funcmp *template.FuncMap) (*template.Template, error) {
// 	tmpls := []string{tmpl}

// 	if parts != nil {
// 		for _, p := range *parts {
// 			tmpls = append(tmpls, p)
// 		}
// 	}

// 	return ht.addTemplate(tmpl, TMPL_PART, tmpls...)
// }

// func (ht *HtmlTemplater) ReadFiles() error {
// 	fsys := os.DirFS(".")
// 	return fs.WalkDir(fsys, ht.ImportPath, func(path string, d fs.DirEntry, err error) error {
// 		if !d.IsDir() {
// 			if strings.Contains(path, ht.FileExtension) {
// 				parent := filepath.Base(filepath.Dir(path))

// 				suffix := strings.TrimPrefix(ht.FileExtension, ".")
// 				suffix = fmt.Sprintf(".%s", suffix)

// 				tmpl := strings.TrimSuffix(filepath.Base(path), suffix)
// 				if parent != ht.ImportPath {
// 					tmpl = filepath.Join(parent, tmpl)
// 				}

// 				t, err := template.ParseFiles(path)
// 				if err != nil {
// 					return err
// 				}

// 				ht.Templates = append(ht.Templates, &HtmlTemplate{
// 					Name:     tmpl,
// 					Template: t,
// 				})
// 			}
// 		}

// 		return nil
// 	})
// }
