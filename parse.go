package main

import (
	"log"
	"strconv"
	"strings"
)

var settings Settings

type Settings struct {
	Package   string // package tester
	IndexName string // name=Status
	IndexType string // type=int
	Swift     bool   // #swift
	Index     bool   // #index
	UpperCase bool   // #up
	Simple    bool   // #simple
}

func (a *Settings) parse(in string) {

	for _, line := range strings.Split(in, "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// package name

		if p := strings.Fields(line); len(p) > 1 && p[0] == "package" {
			settings.Package = p[1]
			continue
		}

		// name=Status type=int
		if strings.HasPrefix(line, "//") {

			typeName := Value(line, "name=")
			if typeName != "" {
				settings.IndexName = typeName
				// fmt.Println("name", settings.IndexName)
			}
			indexType := Value(line, "type=")
			if indexType != "" {
				settings.IndexType = indexType
				// fmt.Println("type", settings.IndexType)
			}

			if strings.Contains(line, "#swift") {
				settings.Swift = true
			}

			if strings.Contains(line, "#index") {
				settings.Index = true
			}

			if strings.Contains(line, "#up") {
				settings.UpperCase = true
			}

			if strings.Contains(line, "#simple") {
				settings.Simple = true
			}
			continue
		}

		// fmt.Println(line)

		// fields
		name, comment := splitFieldLine(line)

		r := new(field)
		r.Name = name

		switch {
		case r.Name == "_":
			r.Go = "_"
		case settings.UpperCase:
			r.Go = goUpperName(r.Name)
		default:
			r.Go = goName(r.Name)
		}

		if settings.IndexName != "" {
			r.Type = settings.IndexName
		} else {
			r.Type = indexType()
		}

		r.Comment = strings.TrimSpace(comment)
		r.parse()
		if !isIdent(r.Go) {
			log.Fatalf("can't make golang name from %q (got %q), set it with go{Name}", name, r.Go)
		}
		fields = append(fields, r)
	}

	if a.Package == "" {
		log.Fatal("missing package line, e.g. package status")
	}
	if len(fields) == 0 {
		log.Fatal("no values in list")
	}

	// тип по умолчанию — int (до подсчета индексов, иначе они остаются нулевыми)
	settings.IndexType = indexType()

	if !isIntType() {
		return
	}

	// значения совпадают с тем, что даст iota в сгенерированном const-блоке:
	// первый элемент = 0, iota=N начинает новый блок со значения N
	count := -1
	for _, x := range fields {
		switch x.Iota > 0 {
		case true:
			count = x.Iota
		default:
			count++
		}
		x.Index = count
	}

}

var intTypes = map[string]bool{
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"byte": true, "rune": true,
}

func indexType() string {
	if settings.IndexType == "" {
		return "int"
	}
	return settings.IndexType
}

func isIntType() bool {
	return intTypes[indexType()]
}

// splitFieldLine делит строку поля на имя и комментарий.
// имя — первое слово (любые символы кроме пробела: sensor.alarm.*, $SYS.REQ.>)
// или строка в кавычках, если нужны пробелы или //: "http://x.com/a b" // title{...}
func splitFieldLine(line string) (name, comment string) {
	if strings.HasPrefix(line, `"`) || strings.HasPrefix(line, "`") {
		q, err := strconv.QuotedPrefix(line)
		if err != nil {
			log.Fatalf("bad quoted name: %s", line)
		}
		name, _ = strconv.Unquote(q)
		rest := strings.TrimSpace(line[len(q):])
		if rest != "" && !strings.HasPrefix(rest, "//") {
			log.Fatalf("expected // after quoted name: %s", line)
		}
		return name, rest
	}

	head := line
	if i := strings.Index(line, "//"); i > -1 {
		head, comment = line[:i], line[i:]
	}
	p := strings.Fields(head)
	if len(p) == 0 {
		log.Fatalf("field name is empty: %q", line)
	}
	return p[0], comment
}
