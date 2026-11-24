package main

import (
	"fmt"
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
func generateJS() (v []byte) {

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
	switch settings.IndexType {
	case "int", "int64", "uint", "uint64", "float64":
		enumType = "Int"
	case "string":
		enumType = "String"
	default:
		panic("enum type")
	}

	list = append(list, fmt.Sprintf("enum %s: %s {\n", enumName, enumType))

	// case
	var caselist []string //id
	var idlist []string
	var titlelist []string
	var desclist []string

	var max int
	for _, x := range fields {
		var name string
		name = x.Name
		if settings.UpperCase {
			name = strings.ToUpper(name)
		}
		name = fmt.Sprintf("case .%s:", name)
		if max < len(name) {
			max = len(name)
			fmt.Println(max)
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

		name := x.Name
		if settings.UpperCase {
			name = strings.ToUpper(name)
		}
		switch enumType {
		case "Int":
			switch x.Iota {
			case 0:
				caselist = append(caselist, name)
			default:
				caselist = append(caselist, fmt.Sprintf("%s = %d", name, x.Iota))
			}
			idlist = append(idlist, "\t\t\t\t"+formatLineSwift(fmt.Sprintf("case .%s:", name), fmt.Sprintf(`return %d`, x.Index), max))
		case "String":
			caselist = append(caselist, name)
			idlist = append(idlist, "\t\t\t\t"+formatLineSwift(fmt.Sprintf("case .%s:", name), fmt.Sprintf(`return "%s"`, x.Name), max))
		}

		titlelist = append(titlelist, "\t\t\t\t"+formatLineSwift(fmt.Sprintf("case .%s:", name), fmt.Sprintf(`return "%s"`, title), max))
		if x.Desc != "" {
			desclist = append(desclist, "\t\t\t\t"+formatLineSwift(fmt.Sprintf("case .%s:", name), fmt.Sprintf(`return "%s"`, x.Desc), max))
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

	// title
	list = append(list, "\t//title")
	list = append(list, "\tvar title: String {")
	list = append(list, "\t\tget {")
	list = append(list, "\t\t\tswitch self {")
	list = append(list, fmt.Sprintf("%s", strings.Join(titlelist, "\n")))
	list = append(list, "\t\t\t}")
	list = append(list, "\t\t}")
	list = append(list, "\t}")
	list = append(list, "\t")

	// desc
	list = append(list, "\t//desc")
	list = append(list, "\tvar desc: String {")
	list = append(list, "\t\tget {")
	list = append(list, "\t\t\tswitch self {")
	list = append(list, fmt.Sprintf("%s", strings.Join(desclist, "\n")))
	if descDefault != "" {
		list = append(list, descDefault)
	}
	list = append(list, "\t\t\t}")
	list = append(list, "\t\t}")
	list = append(list, "\t}")
	list = append(list, "\t")

	// tags

	tags := map[string][]*field{}
	for _, field := range fields {
		for tag := range field.Tags {
			tags[tag] = append(tags[tag], field)
		}
	}
	// if len(tags) == 0 {
	// 	return nil
	// }

	for tagname, fieldlist := range tags {

		var titlelist []string
		for _, x := range fieldlist {
			name := x.Name
			if settings.UpperCase {
				name = strings.ToUpper(name)
			}
			titlelist = append(titlelist, fmt.Sprintf("\t\t\t\tcase .%s: return true", name))
		}
		titlelist = append(titlelist, "\t\t\t\tdefault: return false")

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
