package tmpl

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

var tmpl *template.Template
var templatePaths = make([]string, 0)

func RenderTemplate(w io.Writer, name string, data interface{}) error {
	return tmpl.ExecuteTemplate(w, name, data)
}

// 增加 templates 加载的目录，当你需要添加新的模板目录时，此方法将新的路径加入到
// 模板集，系统会在 LoadTemplates 载入时遍历加载 *.tmpl
func AddTemplatePath(name string) {
	templatePaths = append(templatePaths, name)
}

func LoadTemplates() *template.Template {
	var tmplFiles = make([]string, 0, 1000)
	for _, dir := range templatePaths {
		ls, _ := walkPath(dir)
		tmplFiles = append(tmplFiles, ls...)
	}
	tmpl = template.Must(template.ParseFiles(tmplFiles...))
	return tmpl
}

func AddBuiltinTemplate(builtin *template.Template) {
	filename := builtin.Name()
	ext := filepath.Ext(filename)
	if _, err := tmpl.AddParseTree(filename[:len(filename)-len(ext)], builtin.Tree); err != nil {
		fmt.Printf("add parse tree %s\n", err)
	}
}

func walkPath(dir string) ([]string, error) {
	var tmplFiles = make([]string, 0, 100)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".tmpl" {
			tmplFiles = append(tmplFiles, path)
		}
		fmt.Printf("visited file: %q\n", path)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tmplFiles, nil
}

func init() {
	AddTemplatePath("./pkg/templates")
}
