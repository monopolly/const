package main

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

func existList(list []string, index string) bool {
	return slices.Contains(list, index)
}

func Title(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// isIdent проверяет, что строка — валидный идентификатор Go/Swift
func isIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		switch {
		case r == '_' || unicode.IsLetter(r):
		case i > 0 && unicode.IsDigit(r):
		default:
			return false
		}
	}
	return true
}

// stripBraces убирает содержимое {...}, чтобы #tag и @func внутри desc{} / title{} не считались опциями
func stripBraces(s string) string {
	var b strings.Builder
	depth := 0
	for _, r := range s {
		switch {
		case r == '{':
			depth++
			b.WriteRune(' ')
		case r == '}' && depth > 0:
			depth--
			b.WriteRune(' ')
		case depth > 0:
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// swiftString экранирует строку для Swift-литерала
func swiftString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

// swiftIdent экранирует ключевые слова Swift обратными кавычками
func swiftIdent(s string) string {
	if swiftKeywords[s] {
		return "`" + s + "`"
	}
	return s
}

var swiftKeywords = map[string]bool{
	"associatedtype": true, "class": true, "deinit": true, "enum": true, "extension": true,
	"fileprivate": true, "func": true, "import": true, "init": true, "inout": true,
	"internal": true, "let": true, "open": true, "operator": true, "private": true,
	"precedencegroup": true, "protocol": true, "public": true, "rethrows": true,
	"static": true, "struct": true, "subscript": true, "typealias": true, "var": true,
	"break": true, "case": true, "catch": true, "continue": true, "default": true,
	"defer": true, "do": true, "else": true, "fallthrough": true, "for": true,
	"guard": true, "if": true, "in": true, "repeat": true, "return": true, "throw": true,
	"switch": true, "where": true, "while": true, "as": true, "false": true, "is": true,
	"nil": true, "self": true, "Self": true, "super": true, "throws": true, "true": true,
	"try": true, "Any": true, "Type": true,
}

// any text from start until \t\n\space
func Value(raw string, key string, endtag ...string) (value string) {
	_, after, ok := strings.Cut(raw, key)
	if !ok {
		return
	}
	for _, x := range after {
		switch {
		case unicode.IsSpace(x):
			if endtag != nil {
				value = value + string(x)
			} else {
				return
			}
		default:
			if endtag != nil && string(x) == endtag[0] {
				return
			}
			value = value + string(x)
		}
	}

	return
}

// any {{tag value}}
func Between(raw string, startTag, endTag string) (value string) {
	indexStart := strings.Index(raw, startTag)
	if indexStart == -1 {
		return
	}

	from := indexStart + len(startTag)
	indexEnd := strings.Index(raw[from:], endTag)
	if indexEnd == -1 {
		return
	}
	to := from + indexEnd //+ len(endTag)

	return raw[from:to]
}

// any {{tag value}}
func Betweens(raw string, startTag, endTag string) (value []string) {
	un := map[string]bool{}
	for {
		indexStart := strings.Index(raw, startTag)
		if indexStart == -1 {
			return
		}

		from := indexStart + len(startTag)
		indexEnd := strings.Index(raw[from:], endTag)
		if indexEnd == -1 {
			return
		}
		to := from + indexEnd //+ len(endTag)

		v := raw[from:to]
		if un[v] {
			raw = raw[to:]
			continue
		}
		un[v] = true
		value = append(value, v)
		raw = raw[to:]
	}

}

func formatline(k string, v interface{}) (res string) {
	switch len(k) / 10 {
	case 0:
		return fmt.Sprintf("%-15s %v", k, v)
	case 1:
		return fmt.Sprintf("%-15s %v", k, v)
	case 2:
		return fmt.Sprintf("%-25s %v", k, v)
	case 3:
		return fmt.Sprintf("%-35s %v", k, v)
	default:
		return fmt.Sprintf("%-45s %v", k, v)
	}
}

func formatLineSwift(k string, v interface{}, lens int) (res string) {
	return fmt.Sprintf("%-*s%v", lens, k, v)
}

func sqlword(s string) bool {
	return sqlkeywords[s]
}

// нельзя использовать в качестве sql полей
var sqlkeywords = map[string]bool{
	"all":               true,
	"analyse":           true,
	"analyze":           true,
	"and":               true,
	"any":               true,
	"array":             true,
	"as":                true,
	"asc":               true,
	"asymmetric":        true,
	"authorization":     true,
	"binary":            true,
	"both":              true,
	"case":              true,
	"cast":              true,
	"check":             true,
	"collate":           true,
	"collation":         true,
	"column":            true,
	"concurrently":      true,
	"constraint":        true,
	"create":            true,
	"cross":             true,
	"current_catalog":   true,
	"current_date":      true,
	"current_role":      true,
	"current_schema":    true,
	"current_time":      true,
	"current_timestamp": true,
	"current_user":      true,
	"day":               true,
	"default":           true,
	"deferrable":        true,
	"desc":              true,
	"distinct":          true,
	"do":                true,
	"else":              true,
	"end":               true,
	"except":            true,
	"false":             true,
	"fetch":             true,
	"filter":            true,
	"for":               true,
	"foreign":           true,
	"freeze":            true,
	"from":              true,
	"full":              true,
	"grant":             true,
	"group":             true,
	"having":            true,
	"hour":              true,
	"ilike":             true,
	"in":                true,
	"initially":         true,
	"inner":             true,
	"intersect":         true,
	"into":              true,
	"is":                true,
	"isnull":            true,
	"join":              true,
	"lateral":           true,
	"leading":           true,
	"left":              true,
	"like":              true,
	"limit":             true,
	"localtime":         true,
	"localtimestamp":    true,
	"minute":            true,
	"month":             true,
	"natural":           true,
	"not":               true,
	"notnull":           true,
	"null":              true,
	"offset":            true,
	"on":                true,
	"only":              true,
	"or":                true,
	"order":             true,
	"outer":             true,
	"over":              true,
	"overlaps":          true,
	"placing":           true,
	"primary":           true,
	"references":        true,
	"returning":         true,
	"right":             true,
	"second":            true,
	"select":            true,
	"session_user":      true,
	"similar":           true,
	"some":              true,
	"symmetric":         true,
	"table":             true,
	"tablesample":       true,
	"then":              true,
	"to":                true,
	"trailing":          true,
	"true":              true,
	"union":             true,
	"unique":            true,
	"user":              true,
	"using":             true,
	"variadic":          true,
	"varying":           true,
	"verbose":           true,
	"when":              true,
	"where":             true,
	"window":            true,
	"with":              true,
	"within":            true,
	"without":           true,
	"year":              true,
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}
