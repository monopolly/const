package main

import (
	"strings"
	"unicode"
)

// Имена констант из произвольных строк (type=string), например NATS subjects:
//
//	sensor.alarm.smoke                 > SensorAlarmSmoke
//	sensor.alarm.*                     > SensorAlarmAny
//	sensor.>                           > SensorAll
//	myapp.us-east.users.profile.update > MyappUsEastUsersProfileUpdate
//	obs.gnss.BEMT                      > ObsGnssBEMT
//	wasmbus.ctl.v1.{lattice}.scale     > WasmbusCtlV1LatticeScale
//	$SYS.REQ.SERVER.PING               > SysReqServerPing
//
// #up: SENSOR_ALARM_SMOKE, swift: sensorAlarmSmoke

// спецсимволы, которые сами по себе несут смысл
var tokenWords = map[rune]string{
	'*': "Any", // nats: один токен
	'>': "All", // nats: все оставшиеся токены
}

// nameWords разбивает строку на слова: разделитель — любой символ кроме буквы и цифры
func nameWords(s string) (words []string) {
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			words = append(words, string(cur))
			cur = cur[:0]
		}
	}
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			cur = append(cur, r)
		default:
			flush()
			if w, ok := tokenWords[r]; ok {
				words = append(words, w)
			}
		}
	}
	flush()
	return
}

func isUpperWord(w string) bool {
	hasLetter := false
	for _, r := range w {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}
	return hasLetter
}

// normalizeWords приводит регистр: BEMT сохраняется как аббревиатура,
// но если вся строка в верхнем регистре ($SYS.REQ.SERVER.PING) — слова пишутся как Sys, Req...
func normalizeWords(words []string) []string {
	allUpper := true
	for _, w := range words {
		if w == tokenWords['*'] || w == tokenWords['>'] {
			continue
		}
		if !isUpperWord(w) && strings.IndexFunc(w, unicode.IsLetter) >= 0 {
			allUpper = false
		}
	}

	res := make([]string, len(words))
	for i, w := range words {
		switch {
		case allUpper && isUpperWord(w):
			res[i] = Title(strings.ToLower(w))
		default:
			res[i] = Title(w)
		}
	}
	return res
}

// goName — CamelCase имя константы из строки
func goName(s string) string {
	words := normalizeWords(nameWords(s))
	name := strings.Join(words, "")
	return fixLeadingDigit(name)
}

// goUpperName — #up: SENSOR_ALARM_SMOKE
func goUpperName(s string) string {
	words := nameWords(s)
	for i, w := range words {
		words[i] = strings.ToUpper(w)
	}
	return fixLeadingDigit(strings.Join(words, "_"))
}

// swiftName — lowerCamel для case в Swift: sensorAlarmSmoke
func swiftName(s string) string {
	words := normalizeWords(nameWords(s))
	if len(words) == 0 {
		return ""
	}
	first := words[0]
	switch {
	case isUpperWord(first):
		first = strings.ToLower(first)
	default:
		first = lowerFirst(first)
	}
	words[0] = first
	return fixLeadingDigit(strings.Join(words, ""))
}

// идентификатор не может начинаться с цифры: 42.action > N42Action
func fixLeadingDigit(s string) string {
	if s != "" && unicode.IsDigit([]rune(s)[0]) {
		return "N" + s
	}
	return s
}
