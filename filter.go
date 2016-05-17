package filter

import "strings"

type Filter func(string) string

var Filters = map[string]Filter{
	"uef.isolang":    dcLanguageIso,
	"uef.peerreview": eprintStatus,
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

// Peer review status from input "0" or "1"
func eprintStatus(s string) string {
	if s == "0" {
		return "http://purl.org/eprint/status/NonPeerReviewed"
	} else if s == "1" {
		return "http://purl.org/eprint/status/PeerReviewed"
	}

	return s
}
