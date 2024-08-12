package htmltemplater

import (
	"errors"
	"fmt"
	"html/template"
	"path"
)

type TemplateType int

type TemplaterOptions struct {
	ImportPath    string
	FileExtension string
	Layout        []string
}

type HtmlTemplate struct {
	Template *template.Template
	Name     string
	Type     TemplateType
}

type HtmlTemplater struct {
	Templates     []*HtmlTemplate
	ImportPath    string
	FileExtension string
	Layout        []string
}

const (
	TMPL_PAGE TemplateType = iota
	TMPL_PART
)

var (
	DefaultTemplaterOptions = TemplaterOptions{
		ImportPath: "templates",
		// FileExtension: ".tmpl.html",
		FileExtension: ".gotmpl",
		Layout:        []string{"_layout"},
	}

	ErrTmplNotFound = errors.New("template not found")
)

func (ht HtmlTemplate) String() string {
	var tpls []string
	for _, v := range ht.Template.Templates() {
		tpls = append(tpls, v.Tree.Name)
	}
	return fmt.Sprintf("%s: %+v", ht.Name, tpls)
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

	if len(topts.Layout) == 0 {
		topts.Layout = DefaultTemplaterOptions.Layout
	}

	return &HtmlTemplater{
		ImportPath:    topts.ImportPath,
		FileExtension: topts.FileExtension,
		Layout:        topts.Layout,
	}
}

func (ht *HtmlTemplater) addTemplate(tmpl string, tplType TemplateType, parts ...string) (*template.Template, error) {
	var tmpls []string

	if tplType == TMPL_PART {
		tmpls = append(tmpls, tmpl)
	}

	for _, p := range parts {
		tmpls = append(tmpls, p)
	}

	if tplType == TMPL_PAGE {
		tmpls = append(tmpls, tmpl)
	}

	var pths []string
	for _, tpl := range tmpls {
		pths = append(pths, generateTemplatePath(tpl, ht.ImportPath, ht.FileExtension))
	}

	fmt.Printf("tmpls: %+v\n", tmpls)
	fmt.Printf("pths: %+v\n", pths)

	parsed, err := template.ParseFiles(pths...)
	if err != nil {
		return nil, err
	}

	ht.Templates = append(ht.Templates, &HtmlTemplate{
		Template: parsed,
		Name:     tmpl,
		Type:     tplType,
	})

	return parsed, nil
}

func (ht *HtmlTemplater) AddPage(tmpl string, layout, parts *[]string) (*template.Template, error) {
	var tmpls []string

	if layout != nil {
		for _, l := range *layout {
			tmpls = append(tmpls, l)
		}
	} else {
		tmpls = ht.Layout
	}

	if parts != nil {
		for _, p := range *parts {
			tmpls = append(tmpls, p)
		}
	}

	return ht.addTemplate(tmpl, TMPL_PAGE, tmpls...)
}

func (ht *HtmlTemplater) AddPartial(tmpl string, parts *[]string) (*template.Template, error) {
	tmpls := []string{tmpl}

	if parts != nil {
		for _, p := range *parts {
			tmpls = append(tmpls, p)
		}
	}

	return ht.addTemplate(tmpl, TMPL_PART, tmpls...)
}

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

// func (ht *HtmlTemplater) PageTmpl(child string, others ...string) (*template.Template, error) {
// 	var tmplLayout, tmplChild *template.Template

// 	tmplLayout, err := ht.GetTmpl(ht.LayoutTemplate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tmplChild, err = ht.GetTmpl(child)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range others {
// 		tmpl, err := ht.GetTmpl(v)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// tmplChild, err = tmplChild.AddParseTree(tmpl.Name(), tmpl.Tree)
// 		_, err = tmplLayout.AddParseTree(tmpl.Name(), tmpl.Tree)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	// return tmplChild.AddParseTree(ht.LayoutTemplate, tmplLayout.Tree)
// 	_, err = tmplLayout.AddParseTree(tmplChild.Tree.ParseName, tmplChild.Tree)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return tmplLayout, nil
// }

func generateTemplatePath(tmpl, root, ext string) string {
	// return filepath.Join(ht.ImportPath, fmt.Sprintf("%s%s", tmpl, ht.FileExtension))
	// return path.Join(ht.ImportPath, fmt.Sprintf("%s%s", tmpl, ht.FileExtension))
	fmt.Printf("%v%v\n", tmpl, ext)
	fmt.Println("tmpl: ", tmpl)
	fmt.Println("root: ", root)
	fmt.Println("ext: ", ext)
	return path.Join(root, fmt.Sprintf("%v%v", tmpl, ext))
	// return fmt.Sprintf("%s/%s%s", root, tmpl, ext)
}
