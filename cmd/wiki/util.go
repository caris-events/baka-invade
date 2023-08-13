package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cp "github.com/otiai10/copy"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"gopkg.in/yaml.v3"
)

var (
	minifier *minify.M
)

// Copy
func Copy(src, dst string) error {
	return cp.Copy(src, dst)
}

// Minify
func Minify(contentTyp string, b []byte) ([]byte, error) {
	if minifier == nil {
		minifier = minify.New()
		minifier.AddFunc("text/css", css.Minify)
		minifier.AddFunc("text/html", html.Minify)
	}
	return minifier.Bytes(contentTyp, b)
}

// LoadFiles
func LoadFiles[T any](dir string, slice []T) ([]T, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return slice, err
	}
	for _, v := range files {
		if filepath.Ext(v.Name()) != ".yml" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, v.Name()))
		if err != nil {
			return slice, err
		}
		var object T
		if err := yaml.Unmarshal(b, &object); err != nil {
			return slice, err
		}
		slice = append(slice, object)
	}
	return slice, nil
}

// ReplaceAllStringSubmatchFunc
func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0
	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}
		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}
	return result + str[lastIndex:]
}

// BuildTmpl
func BuildTmpl(tplSrc string, dst string) error {
	funcs := template.FuncMap{
		"html": func(v string) template.HTML {
			return template.HTML(v)
		},
	}
	tpl, err := template.New(fmt.Sprintf("%s.html", tplSrc)).Funcs(funcs).ParseFiles(
		"./tmpl/_disclaimer.html",
		fmt.Sprintf("./tmpl/%s.html", tplSrc),
	)
	if err != nil {
		return err
	}
	if dst != "/" {
		if err := os.MkdirAll(fmt.Sprintf("./../../docs%s", dst), 0644); err != nil {
			return err
		}
	}
	file, err := os.Create(fmt.Sprintf("./../../docs%s", strings.TrimSuffix(dst, "/")+"/index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	var tplBuf bytes.Buffer
	if err = tpl.Execute(&tplBuf, nil); err != nil {
		return err
	}
	b, err := Minify("text/html", tplBuf.Bytes())
	if err != nil {
		return err
	}
	if _, err := file.Write(b); err != nil {
		return err
	}
	return nil
}
