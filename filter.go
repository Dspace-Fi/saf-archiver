package filter

import "strings"

type Filter func(string) string

var Filters = map[string]Filter{
	"uef.isolang": dcLanguageIso,
}

// try to transform strings into dc.language.iso fields
// using ISO 639-1 (two character codes)
// http://www.infoterm.info/standardization/iso_639_1_2002.php
func dcLanguageIso(s string) string {
	t := strings.ToLower(s)

	if t == "suomi" {
		return "FI"
	}

	if t == "ruotsi" {
		return "SV"
	}

	if t == "englanti" {
		return "EN"
	}

	if t == "eesti, viro" {
		return "ET"
	}

	if t == "portugali" {
		return "PT"
	}

	if t == "espanja" {
		return "ES"
	}

	if t == "venäjä" {
		return "RU"
	}

	return s
}
