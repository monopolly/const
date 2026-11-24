package main

import (
	"log"
	"strconv"
	"strings"
)

type field struct {
	Index   int    //for iota
	Name    string //new
	Go      string //New
	Title   string //New & Slim
	Type    string //int
	Comment string //nice touch
	Iota    int
	Tags    map[string]bool //#private #visible
	List    []string        //list{one. two}
	Skip    bool
	Desc    string

	CustomFunc map[string]KV //@code=401 или code{401} > func (a Type) Code() int > 401
	TagOrder   []string      // порядок тегов как в файле
	FuncOrder  []string      // порядок custom функций как в файле
}

func (a *field) parse() {

	a.Tags = make(map[string]bool)
	a.CustomFunc = make(map[string]KV)

	// init comments: убираем только ведущий //, иначе ломаются url{https://...}
	a.Comment = strings.TrimSpace(strings.TrimPrefix(a.Comment, "//"))
	if a.Name == "_" {
		a.Skip = true
	}

	// iota
	a.Iota, _ = strconv.Atoi(Value(stripBraces(a.Comment), "iota="))

	newname := strings.TrimSpace(Between(a.Comment, "name{", "}"))
	if newname != "" {
		a.Name = newname
	}
	newgoname := strings.TrimSpace(Between(a.Comment, "go{", "}"))
	if newgoname != "" {
		a.Go = newgoname
	}

	a.Title = strings.TrimSpace(Between(a.Comment, "title{", "}"))
	if a.Title == "" {
		a.Title = a.Go
	}

	a.Desc = strings.TrimSpace(Between(a.Comment, "desc{", "}"))

	a.parseTags()
	a.parseFunc()
	a.parseList()
	a.parseCustom()

}

// tags
func (a *field) parseTags() {
	for _, x := range strings.Fields(stripBraces(a.Comment)) {
		if !strings.HasPrefix(x, "#") {
			continue
		}
		a.addTag(strings.TrimPrefix(x, "#"))
	}
}

func (a *field) addTag(tag string) {
	tag = Title(strings.TrimSpace(tag))
	if !isIdent(tag) {
		log.Printf("%s: skip invalid tag %q", a.Name, tag)
		return
	}
	if !a.Tags[tag] {
		a.Tags[tag] = true
		a.TagOrder = append(a.TagOrder, tag)
	}
}

func (a *field) addFunc(name string, v KV) {
	if !isIdent(name) {
		log.Printf("%s: skip invalid func name %q", a.Name, name)
		return
	}
	if _, ok := a.CustomFunc[name]; !ok {
		a.FuncOrder = append(a.FuncOrder, name)
	}
	a.CustomFunc[name] = v
}

// tags
func (a *field) parseFunc() {
	for _, x := range strings.Fields(stripBraces(a.Comment)) {
		if !strings.HasPrefix(x, "@") {
			continue
		}
		name, value, ok := strings.Cut(strings.TrimPrefix(x, "@"), "=")
		switch ok {
		case false:
			a.addTag(name)
		default:
			a.addFunc(Title(name), parseValue(value))
		}
	}
}

// list
func (a *field) parseList() {

	for _, x := range strings.Split(Between(a.Comment, "list{", "}"), ",") {
		x = strings.TrimSpace(x)
		if x == "" {
			continue
		}
		a.List = append(a.List, x)
	}

}

type KV struct {
	Type   string // "int", "bool" или "string"
	Bool   bool   // значение (если Type == "bool")
	String string // строковое представление
	Int    int    // числовое значение (если Type == "int")
}

// parseCustom парсит строку вида "title{Nice touch} code{401} CodeName{Some}"
func (a *field) parseCustom() {
	input := a.Comment
	for {
		open := strings.Index(input, "{")
		if open == -1 {
			return
		}
		end := strings.Index(input[open:], "}")
		if end == -1 {
			return
		}
		end += open

		// ключ — последнее слово перед '{' (до него может быть "#active iota=1 ...")
		var key string
		if p := strings.Fields(input[:open]); len(p) > 0 && !strings.HasSuffix(input[:open], " ") {
			key = p[len(p)-1]
		}
		value := input[open+1 : end]
		input = input[end+1:]

		switch key {
		case "", "go", "desc", "title", "list", "name":
			continue
		}
		if !isIdent(key) {
			continue
		}

		a.addFunc(Title(key), parseValue(value))
	}
}

// parseValue определяет тип и парсит значение
func parseValue(s string) KV {
	s = strings.TrimSpace(s)

	// Пробуем парсить как int
	if val, err := strconv.Atoi(s); err == nil {
		return KV{
			Type:   "int",
			String: s,
			Int:    val,
		}
	}

	switch s {
	case "true", "false":
		return KV{Type: "bool", String: s, Bool: s == "true"}
	}

	// Иначе это строка
	return KV{
		Type:   "string",
		String: s,
		Int:    0,
	}
}
