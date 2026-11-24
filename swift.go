package main

import (
	"fmt"
	"log"
	"strings"
)

/*

	enum logo: Int {
	   case me, link, zoom, skype, meet, teams, telegram, slack, youtube
	   var icon: Image {
	        get {
	           switch self{
	            case .me:         return Image("me").resizable()
	            case .link:       return Image("link").resizable()
	            case .zoom:       return Image("zoom").resizable()
	            case .skype:      return Image("skype").resizable()
	            case .meet:       return Image("meet").resizable()
	            case .teams:      return Image("teams").resizable()
	            case .telegram:   return Image("telegram").resizable()
	            case .slack:      return Image("slack").resizable()
	            case .youtube:    return Image("youtube").resizable()
	           }
	        }
	   }
	}
*/

// enum
func generateSwift() (v []byte) {

	var list []string
	list = append(list, "//enum list")

	var enumName string

	switch settings.IndexName != "" {
	case true:
		enumName = fmt.Sprintf("%s", settings.IndexName)
	default:
		enumName = fmt.Sprintf("%s", settings.Package)
	}

	var enumType string
	switch {
	case isIntType():
		enumType = "Int"
	case settings.IndexType == "string":
		enumType = "String"
	default:
		log.Fatalf("swift: unsupported enum type %q", settings.IndexType)
	}

	// "_" в Swift недопустим как case
	fields := activeFields()

	list = append(list, fmt.Sprintf("enum %s: %s {\n", enumName, enumType))

	// case
	var caselist []string //id
	var idlist []string
	var titlelist []string
	var desclist []string

	var max int
	for _, x := range fields {
		name := fmt.Sprintf("case .%s:", swiftCase(x))
		if max < len(name) {
			max = len(name)
		}
	}

	max = max + 2

	var descDefault string
	for _, x := range fields {

		var title string
		switch x.Title != "" {
		case true:
			title = x.Title
		case false:
			title = x.Go
		}

		name := swiftCase(x)
		switch enumType {
		case "Int":
			// rawValue задаем явно: после "_" в Go iota сдвигается, а в Swift — нет
			caselist = append(caselist, fmt.Sprintf("%s = %d", name, x.Index))
			idlist = append(idlist, "\t\t\t\t"+formatLineSwift(fmt.Sprintf("case .%s:", name), fmt.Sprintf(`return %d`, x.Index), max))
		case "String":
			caselist = append(caselist, fmt.Sprintf("%s = %s", name, swiftString(x.Name)))
			idlist = append(idlist, "\t\t\t\t"+formatLineSwift(fmt.Sprintf("case .%s:", name), "return "+swiftString(x.Name), max))
		}

		titlelist = append(titlelist, "\t\t\t\t"+formatLineSwift(fmt.Sprintf("case .%s:", name), "return "+swiftString(title), max))
		if x.Desc != "" {
			desclist = append(desclist, "\t\t\t\t"+formatLineSwift(fmt.Sprintf("case .%s:", name), "return "+swiftString(x.Desc), max))
		} else {
			descDefault = "\t\t\t\t" + formatLineSwift("default:", `return ""`, max)
		}

	}

	// case
	list = append(list, "\t//case")
	list = append(list, fmt.Sprintf("\tcase %s", strings.Join(caselist, ", ")))
	list = append(list, "")

	// id
	list = append(list, "\t//id")
	list = append(list, fmt.Sprintf("\tvar id: %s {", enumType))
	list = append(list, "\t\tget {")
	list = append(list, "\t\t\tswitch self {")
	list = append(list, fmt.Sprintf("%s", strings.Join(idlist, "\n")))
	list = append(list, "\t\t\t}")
	list = append(list, "\t\t}")
	list = append(list, "\t}")
	list = append(list, "\t")

	// title — всегда
	{
		list = append(list, "\t//title")
		list = append(list, "\tvar title: String {")
		list = append(list, "\t\tget {")
		list = append(list, "\t\t\tswitch self {")
		list = append(list, fmt.Sprintf("%s", strings.Join(titlelist, "\n")))
		list = append(list, "\t\t\t}")
		list = append(list, "\t\t}")
		list = append(list, "\t}")
		list = append(list, "\t")
	}

	// desc — только если есть desc{}
	if len(desclist) > 0 {
		list = append(list, "\t//desc")
		list = append(list, "\tvar desc: String {")
		list = append(list, "\t\tget {")
		list = append(list, "\t\t\tswitch self {")
		list = append(list, fmt.Sprintf("%s", strings.Join(desclist, "\n")))
		if descDefault != "" && len(desclist) < len(fields) {
			list = append(list, descDefault)
		}
		list = append(list, "\t\t\t}")
		list = append(list, "\t\t}")
		list = append(list, "\t}")
		list = append(list, "\t")
	}

	// tags

	order, tags := tagGroups()

	for _, tagname := range order {
		fieldlist := tags[tagname]

		var titlelist []string
		for _, x := range fieldlist {
			titlelist = append(titlelist, fmt.Sprintf("\t\t\t\tcase .%s: return true", swiftCase(x)))
		}
		if len(fieldlist) < len(fields) {
			titlelist = append(titlelist, "\t\t\t\tdefault: return false")
		}

		// list = append(list, "\t//tags")
		list = append(list, fmt.Sprintf("\tvar %s: Bool {", strings.ToLower(tagname)))
		list = append(list, "\t\tget {")
		list = append(list, "\t\t\tswitch self {")
		list = append(list, fmt.Sprintf("%s", strings.Join(titlelist, "\n")))
		list = append(list, "\t\t\t}")
		list = append(list, "\t\t}")
		list = append(list, "\t}")
		list = append(list, "\t")
	}

	// end
	list = append(list, "}\n")

	return []byte(strings.Join(list, "\n"))
}

// swiftCase — имя case в Swift (с учетом #up и ключевых слов)
func swiftCase(x *field) string {
	name := swiftName(x.Name)
	if settings.UpperCase {
		name = goUpperName(x.Name)
	}
	return swiftIdent(name)
}
