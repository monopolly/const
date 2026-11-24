package main

import (
	"fmt"
	"strings"
)

var helps Help

type Help struct {
	funcs [][2]string // сгенерированные функции: сигнатура, описание
}

func (a *Help) init() (b *Help) {
	a.funcs = nil
	return a
}

// doc добавляет сгенерированную функцию в справку файла
func (a *Help) doc(sig, desc string) {
	a.funcs = append(a.funcs, [2]string{sig, desc})
}

// опции первой строки-комментария: //name=Status type=int #swift
var structOptions = [][2]string{
	{"name=Status", "Named type for constants: type Status int. Package funcs get type prefix: StatusIndex(), StatusValid(), so many lists can live in one package"},
	{"type=string", "Constants type: int (default), int8..int64, uint..uint64, string"},
	{"#simple", "Constants only, no functions"},
	{"#up", "Uppercase golang names: NEW instead of New"},
	{"#index", "Additional int constants StatusIndexNew, StatusIndexActive... (from 1)"},
	{"#swift", "Create .swift enum next to .go file"},
}

// опции полей: new // title{New} #active code{401}
var fieldOptions = [][2]string{
	{"any.string", "Value can be any string: sensor.alarm.* > SensorAlarmAny, sensor.> > SensorAll, $SYS.REQ.PING > SysReqPing"},
	{"\"a b//c\"", "Quoted value if it has spaces or //"},
	{"_", "Skip value (reserves iota position)"},
	{"iota=100", "Start new const block from this value (int types)"},
	{"name{}", "Replace string value. Ex: art // name{artwork}"},
	{"go{}", "Replace golang const name. Ex: art // go{ArtList}"},
	{"title{}", "Title() for field. Ex: title{Real Estate}. Golang name if not set"},
	{"desc{}", "Desc() for field. Ex: desc{Asset is on sale}. Desc() is generated only if any field has desc{}"},
	{"list{}", "List() []string for field. Ex: list{geoid, country}"},
	{"#tag", "Group by tag. Ex: #virtual > IsVirtual() bool, Virtual() bool, StatusVirtualIndexes()"},
	{"@tag", "Same as #tag"},
	{"key{int}", "Custom func. Ex: code{401} > Code() int"},
	{"key{text}", "Custom func. Ex: label{Nice code} > Label() string"},
	{"key{bool}", "Custom func. Ex: paid{true} > Paid() bool"},
	{"@key=value", "Same as key{value}, value without spaces. Ex: @code=401"},
}

func (a *Help) render() (res []byte) {

	var list []string
	list = append(list, "/*")
	list = append(list, "")
	list = append(list, "Constants Generator")
	list = append(list, "Help to auto-generated correct constants for lists")
	list = append(list, "Martin Prestone (c) 2024")
	list = append(list, "github.com/monopolly")
	list = append(list, "")

	list = append(list, "Struct options (first // line)")
	for _, x := range structOptions {
		list = append(list, Status(x[0], x[1]))
	}
	list = append(list, "")

	list = append(list, "Field options (after // on field line)")
	for _, x := range fieldOptions {
		list = append(list, Status(x[0], x[1]))
	}
	list = append(list, "")
	list = append(list, "Custom funcs: key is capitalized (code{} > Code()), type is int if all values are numbers,")
	list = append(list, "bool if all values are true/false, otherwise string. Missing values return 0 / \"\" / false.")
	list = append(list, "")
	list = append(list, "Golang names: words split by any non letter/digit, * > Any, > > All, {id} > Id,")
	list = append(list, "UPPER words kept (BEMT) unless whole value is upper ($SYS.REQ > SysReq), #up > SENSOR_ALARM_ANY.")

	if len(a.funcs) > 0 {
		list = append(list, "")
		list = append(list, "Generated functions")
		width := 0
		for _, x := range a.funcs {
			width = max(width, len(x[0]))
		}
		for _, x := range a.funcs {
			list = append(list, fmt.Sprintf("%-*s  %s", width, x[0], x[1]))
		}
	}

	// add \t
	for p := range list {
		if list[p] != "" {
			list[p] = "\t" + list[p]
		}
	}
	list = append(list, "")
	list = append(list, "*/\n")

	// render
	res = []byte(strings.Join(list, "\n"))
	return
}

// var oldhelp = `
// /*
// 	Struct options:

// 	[] - make tuple json ex: [1,"name"]
// 	! - New() instead NewItem()
// 	# - create fields tuple index ID = 0, Name = 1 etc, parse []interface...
// 	no - no init function New()
// 	up — title fields ex: ProfileIndexID
// 	map - generate json map for tuple
// 	sql:name - generate sql table
// 	ts:name - generate typescript model
// 	swift - generate swift model

// 	Field option:

// 	GO:
// 	+time - Time() function for int
// 	-func - no field func
// 	++ — add multiply func for int
// 	++ — simple func for boolean ConfirmTrue()/ConfirmFalse()
// 	+raw — write raw json, must be []byte
// 	go:"type" - replace field type go:"[]*News"

// 	JS:
// 	js:name - custom json name ex: js:money.count or js:company.report.name

// 	SQL:
// 	$$ - sql type
// 	sql:"..." - sql text add to field sql:"primary key"
// 	-insert - don't add field to insertSQL
// 	-sql - don't add field in sql table structure

// 	Swift:
// 	@! - field required
// 	@@type - add custom field type
// 	@name - swift file model: @@newsfeed, @@news

// */
// `

func Status(name string, v interface{}) string {
	return fmt.Sprintf("%-14s%v", name, v)
}
