package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// Refresh
func Refresh(c *cli.Context) error {
	if err := Copy("./tmpl/assets", "./../../docs/assets"); err != nil {
		log.Fatalln(err)
	}
	if err := os.WriteFile("./../../docs/CNAME", []byte("baka-invade.org"), 0644); err != nil {
		return err
	}
	if err := BuildTmpl("index", "/"); err != nil {
		return err
	}
	if err := BuildTmpl("about", "/about"); err != nil {
		return err
	}
	if err := BuildTmpl("wiki-search", "/search"); err != nil {
		return err
	}
	if err := BuildTmpl("dict-search", "/dict/search"); err != nil {
		return err
	}
	relationBytes, err := json.Marshal(relations)
	if err != nil {
		return err
	}
	if err := os.WriteFile("./../../docs/assets/scripts/relation_data.js", []byte(fmt.Sprintf("var relations = %s;", string(relationBytes))), 0644); err != nil {
		return err
	}
	if err := os.WriteFile("./../../docs/robots.txt", []byte("Sitemap: https://baka-invade.org/sitemap.xml.gz\nUser-agent: *\nDisallow:"), 0644); err != nil {
		return err
	}
	log.Println("Template Refreshed!")
	return nil
}
