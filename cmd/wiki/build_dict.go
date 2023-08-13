package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/russross/blackfriday/v2"
	"github.com/samber/lo"
)

// PrepareDicts
func PrepareDicts(dicts []*Dict) {
	for _, v := range dicts {
		for _, j := range v.Examples {
			j.Description = ParseDictContent(dicts, j.Description)
			j.Correct = ParseDictContent(dicts, j.Correct)
			j.Incorrect = ParseDictContent(dicts, j.Incorrect)
		}
		v.Description = ParseDictContent(dicts, v.Description)
		v.Code = strings.ReplaceAll(v.Word, " ", "_")
		v.WordTeardown = WordTeardown(v.Word, v.Bopomofo)
		v.ExampleStr = strings.Join(lo.Slice(lo.Map(v.Examples, func(v *DictExample, _ int) string {
			return v.Words[0]
		}), 0, 3), "、")
	}
}

// SimpleDicts
func SimpleDicts(dicts []*Dict) (simps []*SimpleDict) {
	for _, v := range dicts {
		word := v.Word
		if word == v.Code {
			word = ""
		}
		simps = append(simps, &SimpleDict{
			Code:       v.Code,
			Word:       word,
			ExampleStr: v.ExampleStr,
		})
	}
	return
}

// ParseDictContent
func ParseDictContent(dicts []*Dict, v string) string {
	v = string(blackfriday.Run([]byte(v)))

	// [[word]] -> .b-dict-word
	v = ReplaceAllStringSubmatchFunc(regexp.MustCompile(` ?\[\[(.*?)\]\] ?`), v, func(groups []string) string {
		if len(groups) == 0 {
			return ""
		}
		return fmt.Sprintf(`<a class="b-dict-word" href="./../%s">%s</a>`, groups[1], groups[1])
	})

	// {{word}} -> .b-dict-mark
	v = ReplaceAllStringSubmatchFunc(regexp.MustCompile(` ?\{\{(.*?)\}\} ?`), v, func(groups []string) string {
		if len(groups) == 0 {
			return ""
		}
		return fmt.Sprintf(`<span class="b-dict-mark">%s</span>`, groups[1])
	})
	return v
}

// WordTeardown
func WordTeardown(v string, bopomofo string) (teardowns []*DictWordTeardown) {
	foreignBuffer := ""
	bopomofos := strings.Split(bopomofo, " ")

	var bopomofoIndex int
	for _, r := range v {
		if unicode.Is(unicode.Han, r) {
			if foreignBuffer != "" {
				teardowns = append(teardowns, &DictWordTeardown{
					Character: foreignBuffer,
				})
				foreignBuffer = ""
			}
			teardown := &DictWordTeardown{
				Character: string(r),
			}
			currentBopomofo := bopomofos[bopomofoIndex]
			for _, j := range currentBopomofo {
				if string(j) == "ˊ" || string(j) == "ˇ" || string(j) == "ˋ" || string(j) == "˙" {
					continue
				}
				teardown.Bopomofo = append(teardown.Bopomofo, string(j))
			}
			if strings.Contains(currentBopomofo, "ˊ") {
				teardown.Accent = "ˊ"
			} else if strings.Contains(currentBopomofo, "ˇ") {
				teardown.Accent = "ˇ"
			} else if strings.Contains(currentBopomofo, "ˋ") {
				teardown.Accent = "ˋ"
			} else if strings.Contains(currentBopomofo, "˙") {
				teardown.Accent = "˙"
			}
			bopomofoIndex++

			teardowns = append(teardowns, teardown)
		} else {
			foreignBuffer += string(r)
		}
	}
	if foreignBuffer != "" {
		teardowns = append(teardowns, &DictWordTeardown{
			Character: foreignBuffer,
		})
	}
	return
}
