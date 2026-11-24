package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/monopolly/jsons"
)

var jsondemo = jsons.Create()

type field struct {
	Index    int    //for iota
	Name     string //new
	Go       string //New
	Title    string //New & Slim
	HasTitle bool   // title{} задан явно
	Type     string //int
	Value    string // = "name"
	Comment  string //nice touch
	Iota     int
	Tags     map[string]bool //#private #visible
	List     []string        //list{one. two}
	Skip     bool
	Desc     string

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
	a.HasTitle = a.Title != ""
	if a.Title == "" {
		a.Title = a.Go
	}

	a.Desc = strings.TrimSpace(Between(a.Comment, "desc{", "}"))

	a.parseTags()
	a.parseFunc()
	a.parseList()
	a.parseCustom()

}

// func (a *field) parseGolangOptions(v string) {

// 	imports["github.com/monopolly/jsons"] = true
// 	imports["github.com/monopolly/cast"] = true

// 	a.Go.Nofunc = strings.Contains(v, "nofunc")
// 	a.Go.Must = strings.Contains(v, "must")
// 	a.Go.UpperCase = strings.Contains(v, "up")
// 	a.Go.Type = Between(v, `type="`, `"`)
// 	a.Go.Title = Between(v, `title="`, `"`)

// 	name := Between(v, `name="`, `"`)
// 	if name != "" {
// 		a.Go.Name = name
// 	}

// 	// help
// 	helps.golang["nofunc"] = "Do not create any jsons functions for fiels"
// 	helps.golang["up"] = "Make uppercase for functions"
// 	helps.golang["must"] = "Create one validation function for all must fields"
// 	helps.golang[`type=""`] = `Replace golang type. Ex: type="[]*News"`
// 	helps.golang[`name=""`] = `Replace golang struct name. Ex: name="NewsList"`

// }

// // sql options
// func (a *field) parseSQLOptions(v string) {

// 	imports["github.com/jackc/pgx/v5/pgxpool"] = true

// 	a.SQL.NoInsert = strings.Contains(v, "noinsert")
// 	a.SQL.Skip = strings.Contains(v, "skip")
// 	a.SQL.Primary = strings.Contains(v, "primarykey")
// 	a.SQL.Inc = strings.Contains(v, "inc")
// 	a.SQL.Inc = strings.Contains(v, "inc")
// 	if a.SQL.Inc {
// 		a.SQL.NoInsert = true
// 	}
// 	a.SQL.Default.Simple = strings.Contains(v, "defaults")
// 	a.SQL.Default.Unixtime = strings.Contains(v, "unix")
// 	a.SQL.Index.Simple = strings.Contains(v, "index")
// 	// a.SQL.Index.Jsonb = strings.Contains(v, "gin")
// 	a.SQL.Replace = Between(v, `replace="`, `"`)
// 	a.SQL.Append = Between(v, `add="`, `"`)
// 	a.SQL.Unique = Betweens(v, `unique="`, `"`)
// 	a.SQL.Index.Group = Betweens(v, `idx="`, `"`)
// 	a.SQL.Index.Search = Betweens(v, `search="`, `"`)
// 	a.SQL.Type = Between(v, `type="`, `"`)
// 	name := Between(v, `name="`, `"`)
// 	if name != "" {
// 		a.SQL.Name = name
// 	}

// 	// help
// 	helps.sql["noinsert"] = "Do use field for insert function"
// 	helps.sql["index"] = "Create simple default index or gin for jsonb"
// 	// helps.sql["gin"] = "Create jsonb gin index"
// 	helps.sql["idx"] = `Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"`
// 	helps.sql["search"] = `Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');`
// 	helps.sql["defaults"] = "Add default value based on field type"
// 	helps.sql["unix"] = "Add default value: extract(epoch from now())"
// 	helps.sql["skip"] = "Skip field for sql queries"
// 	helps.sql[`replace=""`] = `Rewrite sql for field. Ex: replace="bigint primary key"`
// 	helps.sql[`type=""`] = `Rewrite sql type for field. Ex: type="jsonb"`
// 	helps.sql[`add=""`] = `Append sql for field. Ex: replace="primary key"`
// 	helps.sql[`unique=""`] = `Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"`
// }

// func (a *field) parseJSOptions(v string) {

// 	a.Json.Name = Between(v, `name="`, `"`)
// 	a.Json.Skip = strings.Contains(v, "skip")
// 	a.Json.Bool = strings.Contains(v, "bool")
// 	a.Json.Inc = strings.Contains(v, "inc")
// 	a.Json.Time = strings.Contains(v, "time")
// 	a.Json.Raw = strings.Contains(v, "raw")

// 	helps.js["skip"] = `Skip json field`
// 	helps.js["inc"] = "Add inc jsons function for numbers fields"
// 	helps.js["bool"] = "Add set jsons function for bool fields"
// 	helps.js["time"] = "Create convert jsons function for unixtime fields"
// 	helps.js["raw"] = "Set Raw json function inside field"
// 	helps.js[`name=""`] = `Replace json field name. Ex: name="sid"`

// }
// func (a *field) parseSwiftOptions(v string) {

// 	a.Swift.Files = Betweens(v, `file="`, `"`)
// 	a.Swift.Type = Between(v, `type="`, `"`)
// 	a.Swift.Skip = strings.Contains(v, "skip")
// 	a.Swift.Enum.Skip = strings.Contains(v, "skip")
// 	a.Swift.Must = strings.Contains(v, "must")

// 	// swift
// 	helps.swift["must"] = "Swift field with required values"
// 	helps.swift["skip"] = "Ignore field for Swift"
// 	helps.swift[`type=""`] = `Replace type for Swift model. Ex: type="string"`
// 	helps.swift[`file=""`] = `Create another swift file for this field Swift model. file="f1", file="f2"`

// }

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
