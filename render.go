package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"html/template"
)

type Render interface {
	Render(w http.ResponseWriter, name string, data interface{}) error
}

func importFunc(arg ...string) string {
	if len(arg) < 1 {
		return ""
	}
	return ""
}

func marshalFunc(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

var defaultFuncs = template.FuncMap{
	"import":  importFunc,
	"marshal": marshalFunc,
}

type ProdRender struct {
	root string
	tpls map[string]*template.Template
}

func NewRender(root string, debug bool) Render {
	if debug {
		return NewDebugRender(root)
	}

	return NewProdRender(root)
}

func NewProdRender(root string) *ProdRender {
	r := &ProdRender{
		root: root,
		tpls: make(map[string]*template.Template),
	}
	return r
}

func (r *ProdRender) Render(w http.ResponseWriter, name string, data interface{}) error {
	var t *template.Template
	var has bool
	if t, has = r.tpls[name]; !has {
		t = r.load_template(name)
		r.tpls[name] = t
	}
	err := t.Execute(w, data)
	if err != nil {
		log.Printf("render %s fail: %s", name, err)
	}
	return err
}

func get_deps_from(dir, name string) []string {
	f, err := os.Open(name)
	if err != nil {
		return nil
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	onBrace := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i+1 < len(data); i++ {
			if data[i] == '}' && data[i+1] == '}' {
				return i + 2, data[:i], bufio.ErrFinalToken
			}
			if data[i] == '{' && data[i+1] == '{' {
				return i + 2, nil, nil
			}
		}
		return 0, nil, bufio.ErrFinalToken
	}

	s.Split(onBrace)
	var depLine string
	if s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "import") {
			depLine = strings.TrimSpace(line[6:])
		}
	}

	s = bufio.NewScanner(strings.NewReader(depLine))
	onQuote := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == '"' {
				return i + 1, data[:i], nil
			}
		}
		return 0, nil, bufio.ErrFinalToken
	}
	s.Split(onQuote)
	var deps []string
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) < 1 {
			continue
		}
		deps = append(deps, filepath.Join(dir, line))
	}
	return deps
}

func (r *ProdRender) load_template(name string) *template.Template {
	deps := make(map[string]bool)
	f := filepath.Join(r.root, name)
	files := []string{f}
	for i := 0; i < len(files); i++ {
		f = files[i]
		if _, has := deps[f]; has {
			continue
		}
		deps[f] = true
		for _, f := range get_deps_from(r.root, f) {
			if _, has := deps[f]; has {
				continue
			}
			files = append(files, f)
		}
	}

	return template.Must(template.New(name).Funcs(defaultFuncs).ParseFiles(files...))
}

type DebugRender struct {
	*ProdRender
}

func NewDebugRender(root string) *DebugRender {
	r := &DebugRender{
		ProdRender: NewProdRender(root),
	}
	return r
}

func (r *DebugRender) Render(w http.ResponseWriter, name string, data interface{}) error {
	t := r.load_template(name)
	err := t.Execute(w, data)
	if err != nil {
		log.Printf("render %s fail: %s", name, err)
	}
	return err
}
