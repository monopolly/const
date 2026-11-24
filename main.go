package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/monopolly/file"
)

var (
	fields             []*field //список всех полей
	helpFile           bytes.Buffer
	goFile             bytes.Buffer //для структуры
	fieldnameSQLErrors []string     //список полей которые нельзя использовать в sql
)

//go build -o $GOPATH/bin/jsons

func main() {

	// help init
	helps.init()

	//парсим входящие файл
	in := parsePath()

	//парсит файл и создает структуру
	settings.parse(in)

	// go
	render()

}

// indexfile
func render() {

	// package
	goFile.Write(header())

	// import
	// goFile.Write(importers())

	// go (до help: help перечисляет сгенерированные функции)
	body := golang()

	// help
	goFile.Write(helps.render())

	//go
	goFile.Write(body)

	// gofmt: заодно проверяет, что сгенерирован синтаксически корректный код
	out, err := format.Source(goFile.Bytes())
	if err != nil {
		log.Printf("generated go code is invalid: %v", err)
		out = goFile.Bytes()
	}
	file.Save(settings.file.golang, out)

	// swift
	if settings.Swift {
		file.Save(settings.file.swift, generateSwift())
	}

}

func parsePath() (in string) {
	files := os.Args
	if len(files) != 2 {
		log.Fatal("Const generate golang const list\n", "Example: ", files[0]+" status")
		return
	}

	name := files[1]

	b := file.OpenE(name)
	if b == nil {
		log.Fatalf("can't read %s", name)
	}

	in = string(b)

	name = strings.TrimSuffix(name, filepath.Ext(name))

	// не перезаписываем исходный файл результатом (status.go > status.go)
	for _, ext := range []string{".go", ".swift", ".js"} {
		if name+ext == files[1] {
			log.Fatalf("input file %s would be overwritten by output, rename it (e.g. without extension)", files[1])
		}
	}

	//in = "test.go"
	settings.file.path = path.Dir(name)
	settings.file.golang = fmt.Sprintf("%s.go", name)
	settings.file.js = fmt.Sprintf("%s.js", name)
	settings.file.swift = fmt.Sprintf("%s.swift", name)
	return
}

func header() []byte {
	return []byte(fmt.Sprintf("package %s\n", settings.Package))
}

var imports = map[string]bool{
	// "context": true,
	// "fmt":     true,
	// "reflect": true,
	// "strings": true,
	// "time":    true,
}

// "github.com/jackc/pgx/v5/pgxpool"
// "github.com/monopolly/cast"
// "github.com/monopolly/jsons"

func importers() []byte {

	list := []string{}
	for x := range imports {
		list = append(list, fmt.Sprintf(`"%s"`, x))
	}

	p := fmt.Sprintf(`
	import (
		%s
	)
	`, strings.Join(list, "\n"))

	return []byte(p)
}
