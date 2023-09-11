package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/caris-events/sitemap-generator/smg"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

const (
	objectDir = "./../../database/wiki"
	dictDir   = "./../../database/dict"
)

// Update
func Update(c *cli.Context) error {
	objects, err := LoadFiles(objectDir, []*Object{})
	if err != nil {
		return err
	}
	PrepareObjects(objects)

	dicts, err := LoadFiles(dictDir, []*Dict{})
	if err != nil {
		return err
	}
	PrepareDicts(dicts)

	cache, err := LoadCache()
	if err != nil {
		return err
	}

	objectChanged, err := CompileObjects(cache, objects)
	if err != nil {
		return err
	}
	objCodes := lo.Map(objects, func(v *Object, _ int) string {
		return v.Code
	})
	objSimps := SimpleObjects(objCodes, objects)

	if objectChanged {
		codeBytes, err := json.Marshal(objCodes)
		if err != nil {
			return err
		}
		simpBytes, err := json.Marshal(objSimps)
		if err != nil {
			return err
		}
		if err = os.WriteFile("./../../docs/assets/scripts/objects.js", []byte(fmt.Sprintf(`
			var object_codes = %s;
			var objects = %s;
		`, string(codeBytes), string(simpBytes))), 0644); err != nil {
			return err
		}
	}

	dictChanged, err := CompileDicts(cache, dicts)
	if err != nil {
		return err
	}
	dictSimps := SimpleDicts(dicts)

	if dictChanged {
		simpBytes, err := json.Marshal(dictSimps)
		if err != nil {
			return err
		}
		if err = os.WriteFile("./../../docs/assets/scripts/dicts.js", []byte(fmt.Sprintf(`
			var dicts = %s;
		`, string(simpBytes))), 0644); err != nil {
			return err
		}
	}

	if objectChanged || dictChanged {
		//
		objRands := lo.Map(lo.Samples(objSimps, 60), func(v *SimpleObject, _ int) *SimpleObject {
			v.Code = v.Object.Code
			return v
		})
		objBytes, err := json.Marshal(objRands)
		if err != nil {
			return err
		}
		dictRands := lo.Samples(dictSimps, 60)
		dictBytes, err := json.Marshal(dictRands)
		if err != nil {
			return err
		}
		if err = os.WriteFile("./../../docs/assets/scripts/random_data.js", []byte(fmt.Sprintf(`
			var random_objects = %s;
			var random_dicts = %s;
		`, string(objBytes), string(dictBytes))), 0644); err != nil {
			return err
		}
		if err := CompileSitemap(objects, dicts); err != nil {
			return err
		}
	}

	if err := SaveCache(cache); err != nil {
		return err
	}

	return nil
}

func CompileSitemap(objects []*Object, dicts []*Dict) error {
	now := time.Now()

	sitemap := smg.NewSitemap(true)
	sitemap.SetName("sitemap")
	sitemap.SetHostname("https://baka-invade.org")
	sitemap.SetOutputPath("./../../docs")
	sitemap.SetLastMod(&now)

	sitemap.Add(&smg.SitemapLoc{
		Loc: "about/",
	})
	sitemap.Add(&smg.SitemapLoc{
		Loc: "index",
	})
	sitemap.Add(&smg.SitemapLoc{
		Loc: "search/",
	})
	sitemap.Add(&smg.SitemapLoc{
		Loc: "dict/search/",
	})

	for _, v := range objects {
		sitemap.Add(&smg.SitemapLoc{
			Loc: fmt.Sprintf("%s/", v.Code),
		})
	}
	for _, v := range dicts {
		sitemap.Add(&smg.SitemapLoc{
			Loc: fmt.Sprintf("dict/%s/", v.Code),
		})
	}

	filenames, err := sitemap.Save()
	if err != nil {
		return err
	}
	for _, v := range filenames {
		log.Printf("Sitemap Saved: %s", v)
	}
	return nil
}

// CompileObjectCovers
func CompileObjectCovers(cache *Cache, objects []*Object) error {
	for _, v := range objects {
		if !IsChecksumChanged(cache, v.Code, FileTypeObjectLogo) {
			continue
		}
		log.Printf("Object Logo Changed: %s", v.Code)
		// TODO: copy logo

		// TODO: Conflict with CompileObjects first time if cover not exists
		if err := DrawObjectCover(v); err != nil {
			return err
		}
		UpdateCache(cache, v.Code, FileTypeObjectLogo)
	}
	return nil
}

func CompileDicts(cache *Cache, dicts []*Dict) (bool, error) {
	tpl, err := template.New("dict-single.html").Funcs(template.FuncMap{
		"html": HTML,
	}).ParseFiles(
		"./tmpl/_disclaimer.html",
		"./tmpl/dict-single.html",
	)
	if err != nil {
		return false, err
	}

	var isChanged bool
	for _, v := range dicts {
		if !IsChecksumChanged(cache, v.Code, FileTypeDict) {
			continue
		}
		log.Printf("Dict Changed: %s", v.Code)

		dir := fmt.Sprintf("./../../docs/dict/%s", v.Code)
		if err := os.MkdirAll(dir, 0644); err != nil {
			return false, err
		}
		file, err := os.Create(fmt.Sprintf("%s/index.html", dir))
		if err != nil {
			return false, err
		}
		defer file.Close()

		var tplBuf bytes.Buffer
		if err = tpl.Execute(&tplBuf, v); err != nil {
			return false, err
		}
		b, err := Minify("text/html", tplBuf.Bytes())
		if err != nil {
			return false, err
		}
		if _, err := file.Write(b); err != nil {
			return false, err
		}
		if err := DrawDictCover(v); err != nil {
			return false, err
		}
		isChanged = true

		//
		UpdateCache(cache, v.Code, FileTypeDict)
	}
	return isChanged, nil
}

func CompileObjects(cache *Cache, objects []*Object) (bool, error) {
	tpl, err := template.New("wiki-single.html").Funcs(template.FuncMap{
		"isDangerousIncome":    IsDangerousIncome,
		"isDangerousDirection": IsDangerousDirection,
		"isDangerousInvasion":  IsDangerousInvasion,
		"isDangerousOwner":     IsDangerousOwner,
		"minus":                Minus,
		"html":                 HTML,
	}).ParseFiles(
		"./tmpl/_disclaimer.html",
		"./tmpl/wiki-single.html",
	)
	if err != nil {
		return false, err
	}

	var isChanged bool
	for _, v := range objects {
		if !IsChecksumChanged(cache, v.Code, FileTypeObject) {
			continue
		}
		log.Printf("Object Changed: %s", v.Code)

		dir := fmt.Sprintf("./../../docs/%s", v.Code)
		if err := os.MkdirAll(dir, 0644); err != nil {
			return false, err
		}
		file, err := os.Create(fmt.Sprintf("%s/index.html", dir))
		if err != nil {
			return false, err
		}
		defer file.Close()

		// TODO: necessary?
		if err := CopyLogo(v); err != nil {
			return false, err
		}

		var tplBuf bytes.Buffer
		if err = tpl.Execute(&tplBuf, v); err != nil {
			return false, err
		}
		b, err := Minify("text/html", tplBuf.Bytes())
		if err != nil {
			return false, err
		}
		if _, err := file.Write(b); err != nil {
			return false, err
		}
		if err := DrawObjectCover(v); err != nil {
			return false, err
		}
		isChanged = true

		//
		UpdateCache(cache, v.Code, FileTypeObjectLogo)
		UpdateCache(cache, v.Code, FileTypeObject)
	}
	return isChanged, nil
}

func CopyLogo(v *Object) error {
	dir := fmt.Sprintf("./../../docs/%s", v.Code)
	src := fmt.Sprintf("./../../database/wiki_logos/%s.jpg", v.Code)
	dst := fmt.Sprintf("%s/%s.jpg", dir, v.Code)

	err := Copy(src, dst)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return Copy("./assets/images/logo.jpg", fmt.Sprintf("%s/%s.jpg", dir, v.Code))
	}
	return err
}

func HTML(v string) template.HTML {
	return template.HTML(v)
}
