package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

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

// Value — текст после key до первого пробела: Value("a iota=5 b", "iota=") > "5"
func Value(raw, key string) string {
	_, after, ok := strings.Cut(raw, key)
	if !ok {
		return ""
	}
	if i := strings.IndexFunc(after, unicode.IsSpace); i > -1 {
		return after[:i]
	}
	return after
}

// Between — текст между startTag и endTag: Between("title{News}", "title{", "}") > "News"
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

func formatLineSwift(k string, v any, lens int) (res string) {
	return fmt.Sprintf("%-*s%v", lens, k, v)
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}
