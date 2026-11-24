package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// isNamed — есть ли отдельный тип (name=Status)
func isNamed() bool {
	return settings.IndexName != ""
}

// valueType — тип констант: Status или int/string
func valueType() string {
	if isNamed() {
		return settings.IndexName
	}
	return indexType()
}

// declared — все объявленные имена, чтобы не сгенерировать дубликат (ошибка компиляции)
var declared = map[string]string{}

// declare регистрирует имя пакетного уровня или метода
func declare(name, what string) {
	key := name
	if strings.HasPrefix(what, "method") {
		key = "method:" + name
	}
	if prev, ok := declared[key]; ok {
		log.Fatalf("name conflict: %s %q clashes with %s (rename one with go{Name})", what, name, prev)
	}
	declared[key] = what
}

// accessor — заголовок функции, возвращающей значение по константе.
// с name= это метод (a Status) Title(), без него — функция Title(a int).
// desc попадает в справку сгенерированного файла
func accessor(name, ret, desc string) string {
	if isNamed() {
		declare(name, "method "+name)
		helps.doc(fmt.Sprintf("(%s) %s() %s", settings.IndexName, name, ret), desc)
		return fmt.Sprintf("func (a %s) %s() %s {", settings.IndexName, name, ret)
	}
	declare(name, "func "+name)
	helps.doc(fmt.Sprintf("%s(%s) %s", name, indexType(), ret), desc)
	return fmt.Sprintf("func %s(a %s) %s {", name, indexType(), ret)
}

// function — заголовок функции пакетного уровня: name(args) ret.
// с name=Status имя получает префикс типа (StatusIndex, StatusValid...),
// чтобы в одном пакете можно было держать несколько списков
func function(name, args, ret, desc string) string {
	name = settings.IndexName + name
	declare(name, "func "+name)
	helps.doc(fmt.Sprintf("%s(%s) %s", name, args, ret), desc)
	return fmt.Sprintf("func %s(%s) %s {", name, args, ret)
}

// activeFields — поля без "_"
func activeFields() (res []*field) {
	for _, x := range fields {
		if x.Skip {
			continue
		}
		res = append(res, x)
	}
	return
}

func goNames(list []*field) []string {
	var res []string
	for _, x := range list {
		res = append(res, x.Go)
	}
	return res
}

// название полей для json
func generateGo() []byte {

	//список
	var list []string

	settings.IndexType = indexType()

	for _, x := range activeFields() {
		declare(x.Go, "const "+x.Go)
	}

	if isNamed() {
		declare(settings.IndexName, "type "+settings.IndexName)
		list = append(list, fmt.Sprintf("type %s %s", settings.IndexName, settings.IndexType))
	}

	switch isIntType() {
	case true:
		list = append(list, generateGoIntType()...)
	default:
		list = append(list, generateGoStringType()...)
	}

	switch settings.Simple {
	case true:
	default:
		list = append(list, generateGoString()...)
		list = append(list, generateGoInt()...)
		list = append(list, generateGoValidTypeKey()...)
		list = append(list, generateGoValidKey()...)
		list = append(list, generateGoList()...)
		list = append(list, generateTags()...)
		list = append(list, generateList()...)
		list = append(list, generateIndex()...)
		list = append(list, generateGoIsKey()...)
		list = append(list, generateFunc()...)
		list = append(list, generateCustomFunc()...)
	}

	return []byte(strings.Join(list, "\n"))
}

func comment(x *field) string {
	if x.Comment == "" {
		return ""
	}
	return " //" + x.Comment
}

func generateGoIntType() []string {
	var list []string

	list = append(list, "const (")
	for pos, x := range fields {

		if pos == 0 {
			switch x.Iota > 0 {
			case true:
				list = append(list, fmt.Sprintf("%s = %s(iota) + %d%s", x.Go, x.Type, x.Iota, comment(x)))
			default:
				list = append(list, fmt.Sprintf("%s = %s(iota)%s", x.Go, x.Type, comment(x)))
			}
			continue
		}

		switch x.Iota > 0 {
		case true:
			list = append(list, ")\n")
			list = append(list, "const (")
			list = append(list, fmt.Sprintf("%s = %s(iota) + %d%s", x.Go, x.Type, x.Iota, comment(x)))
		default:
			list = append(list, fmt.Sprintf("%s%s", x.Go, comment(x)))
		}
	}

	list = append(list, ")")

	return list
}

func generateGoStringType() []string {
	var list []string

	list = append(list, "const (")
	for _, x := range fields {
		switch isNamed() {
		case true:
			list = append(list, fmt.Sprintf(`%s = %s(%s)%s`, x.Go, x.Type, strconv.Quote(x.Name), comment(x)))
		case false:
			list = append(list, fmt.Sprintf(`%s = %s%s`, x.Go, strconv.Quote(x.Name), comment(x)))
		}
	}

	list = append(list, ")")

	return list
}

func generateGoList() []string {
	var list []string

	list = append(list, function("Indexes", "", "[]"+valueType(), "All values"))
	list = append(list, fmt.Sprintf("return []%s{%s}", valueType(), strings.Join(goNames(activeFields()), ", ")))
	list = append(list, "}")

	return list
}

// switch по константам с возвратом строки; поля с пустым значением не попадают в switch
func stringSwitch(name, desc string, value func(x *field) string) []string {
	var list []string
	list = append(list, "\n")
	list = append(list, accessor(name, "string", desc))
	list = append(list, "switch a {")
	for _, x := range activeFields() {
		v := value(x)
		if v == "" {
			continue
		}
		list = append(list, fmt.Sprintf(`case %s: return %s`, x.Go, strconv.Quote(v)))
	}
	list = append(list, `default: return ""`)
	list = append(list, "}")
	list = append(list, "}")
	return list
}

func anyField(f func(x *field) bool) bool {
	for _, x := range activeFields() {
		if f(x) {
			return true
		}
	}
	return false
}

// String() и Title() всегда; Desc() — если есть хоть один desc{}
func generateGoString() []string {

	var list []string
	list = append(list, stringSwitch("String", "String value", func(x *field) string { return x.Name })...)

	list = append(list, stringSwitch("Title", "title{}, golang name if not set", func(x *field) string { return x.Title })...)

	if anyField(func(x *field) bool { return x.Desc != "" }) {
		list = append(list, stringSwitch("Desc", "desc{}", func(x *field) string { return x.Desc })...)
	}
	list = append(list, "\n")

	return list

}

// get struct value with function
func generateGoInt() []string {

	if !isIntType() || !isNamed() {
		return nil
	}
	var list []string
	list = append(list, "\n")
	list = append(list, accessor("Int", "int", "Value as int"))
	list = append(list, "return int(a)")
	list = append(list, "}")
	list = append(list, "\n")

	return list

}

// get struct value with function
func generateGoValidKey() []string {

	var list []string
	list = append(list, "\n")
	list = append(list, function("Valid", "v "+settings.IndexType, "bool", "Is raw "+settings.IndexType+" a known value"))
	list = append(list, fmt.Sprintf("switch %s(v) {", valueType()))
	list = append(list, fmt.Sprintf(`case %s: return true`, strings.Join(goNames(activeFields()), ", ")))
	list = append(list, `default: return false`)
	list = append(list, "}")
	list = append(list, "}")
	list = append(list, "\n")

	return list

}

// get struct value with function
func generateGoValidTypeKey() []string {
	if !isNamed() {
		return nil
	}

	var list []string
	list = append(list, "\n")
	list = append(list, function("ValidIndexType", "v "+settings.IndexName, "bool", "Is value known"))
	list = append(list, "switch v {")
	list = append(list, fmt.Sprintf(`case %s: return true`, strings.Join(goNames(activeFields()), ", ")))
	list = append(list, `default: return false`)
	list = append(list, "}")
	list = append(list, "}")
	list = append(list, "\n")

	list = append(list, "\n")
	list = append(list, accessor("Valid", "bool", "Is value known"))
	list = append(list, fmt.Sprintf("return %sValidIndexType(a)", settings.IndexName))
	list = append(list, "}")
	list = append(list, "\n")

	return list

}

// tagGroups — теги в порядке появления в файле (map дает случайный порядок и нестабильный вывод)
func tagGroups() (order []string, tags map[string][]*field) {
	tags = map[string][]*field{}
	for _, f := range activeFields() {
		for _, tag := range f.TagOrder {
			if _, ok := tags[tag]; !ok {
				order = append(order, tag)
			}
			tags[tag] = append(tags[tag], f)
		}
	}
	return
}

// boolSwitch — true для перечисленных констант
func boolSwitch(header string, list []*field) []string {
	var res []string
	res = append(res, header)
	res = append(res, "switch a {")
	res = append(res, fmt.Sprintf("case %s: return true", strings.Join(goNames(list), ", ")))
	res = append(res, "default: return false")
	res = append(res, "}")
	res = append(res, "}")
	return res
}

func tagListFunc(name, tag string, list []*field) []string {
	var res []string
	res = append(res, function(name, "", "[]"+valueType(), "Values with #"+lowerFirst(tag)))
	res = append(res, fmt.Sprintf(`return []%s{%s}`, valueType(), strings.Join(goNames(list), ", ")))
	res = append(res, "}\n")
	return res
}

// Status: StatusVirtualIndexes() + (a Status) IsVirtual()
// без name: VirtualList() + IsVirtual(a int)
func generateTags() []string {
	order, tags := tagGroups()
	if len(order) == 0 {
		return nil
	}

	var list []string
	list = append(list, "\n")
	for _, tag := range order {
		switch isNamed() {
		case true:
			list = append(list, tagListFunc(tag+"Indexes", tag, tags[tag])...)
		default:
			list = append(list, tagListFunc(tag+"List", tag, tags[tag])...)
		}
		list = append(list, "\n")
		list = append(list, boolSwitch(accessor("Is"+tag, "bool", "Has #"+lowerFirst(tag)), tags[tag])...)
		list = append(list, "\n")
	}

	return list
}

// Status: StatusVirtualList() + (a Status) Virtual()
// без name эти функции совпадают с generateTags, поэтому не генерируются
func generateFunc() []string {
	if !isNamed() {
		return nil
	}

	order, tags := tagGroups()
	if len(order) == 0 {
		return nil
	}

	var list []string
	list = append(list, "\n")
	for _, tag := range order {
		list = append(list, tagListFunc(tag+"List", tag, tags[tag])...)
		list = append(list, "\n")
		list = append(list, boolSwitch(accessor(tag, "bool", "Has #"+lowerFirst(tag)), tags[tag])...)
		list = append(list, "\n")
	}

	return list
}

// пользовательские функции:
// code{401}      > Code() int
// label{Nice}    > Label() string
// paid{true}     > Paid() bool
// @code=401      > Code() int
func generateCustomFunc() (list []string) {

	var order []string
	tags := map[string][]*field{}
	tagTypes := map[string]string{}

	for _, f := range activeFields() {
		for _, tag := range f.FuncOrder {
			v := f.CustomFunc[tag]
			t, ok := tagTypes[tag]
			switch {
			case !ok:
				order = append(order, tag)
				tagTypes[tag] = v.Type
			case t != v.Type:
				// разные типы значений (401 и abc) — функция строковая
				tagTypes[tag] = "string"
			}
			tags[tag] = append(tags[tag], f)
		}
	}

	if len(order) == 0 {
		return nil
	}

	list = append(list, "\n")
	for _, tag := range order {

		tagType := tagTypes[tag] // int / string / bool
		goType := tagType

		var cases []string
		for _, x := range tags[tag] {
			v := x.CustomFunc[tag]
			switch tagType {
			case "string":
				cases = append(cases, fmt.Sprintf(`case %s: return %s`, x.Go, strconv.Quote(v.String)))
			case "int":
				cases = append(cases, fmt.Sprintf(`case %s: return %d`, x.Go, v.Int))
			case "bool":
				if v.Bool {
					cases = append(cases, x.Go)
				}
			}
		}

		desc := fmt.Sprintf("%s{%s}", lowerFirst(tag), tagType)
		list = append(list, accessor(tag, goType, desc))

		switch tagType {
		case "bool":
			// для bool достаточно перечислить значения с true
			switch len(cases) {
			case 0:
				list = append(list, "return false")
			default:
				list = append(list, "switch a {")
				list = append(list, fmt.Sprintf("case %s: return true", strings.Join(cases, ", ")))
				list = append(list, "default: return false")
				list = append(list, "}")
			}
		case "int":
			list = append(list, "switch a {")
			list = append(list, cases...)
			list = append(list, `default: return 0`)
			list = append(list, "}")
		default:
			list = append(list, "switch a {")
			list = append(list, cases...)
			list = append(list, `default: return ""`)
			list = append(list, "}")
		}
		list = append(list, "}")
		list = append(list, "\n")
	}

	return
}

func generateList() []string {

	var withList []*field
	for _, x := range activeFields() {
		if len(x.List) > 0 {
			withList = append(withList, x)
		}
	}

	if len(withList) == 0 {
		return nil
	}

	var list []string
	list = append(list, "\n")
	list = append(list, accessor("List", "[]string", "list{}"))
	list = append(list, "switch a {")
	for _, x := range withList {
		var quoted []string
		for _, v := range x.List {
			quoted = append(quoted, strconv.Quote(v))
		}
		list = append(list, fmt.Sprintf("case %s: return []string{%s}", x.Go, strings.Join(quoted, ", ")))
	}
	list = append(list, "default: return nil")
	list = append(list, "}")
	list = append(list, "}")

	list = append(list, "\n")

	return list
}

func generateIndex() []string {

	if !settings.Index {
		return nil
	}

	var list []string
	list = append(list, "\n")

	list = append(list, "const (")
	list = append(list, "_ = int(iota)")

	for _, x := range activeFields() {
		name := settings.IndexName + "Index" + x.Go
		declare(name, "const "+name)
		list = append(list, name)
	}

	list = append(list, ")")
	helps.doc(settings.IndexName+"Index<Name> const", "#index: int position of value, from 1")

	return list
}

// Index("new") > int(New)
func generateGoIsKey() []string {
	if !isIntType() {
		return nil
	}

	var list []string
	list = append(list, "\n")
	list = append(list, function("Index", "v string", "int", "String value to int, 0 if unknown"))
	list = append(list, "switch v {")
	for _, x := range activeFields() {
		list = append(list, fmt.Sprintf(`case %s: return int(%s)`, strconv.Quote(x.Name), x.Go))
	}
	list = append(list, `default: return 0`)
	list = append(list, "}")
	list = append(list, "}")
	list = append(list, "\n")

	return list

}
