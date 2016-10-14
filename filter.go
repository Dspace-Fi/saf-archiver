package filter

import (
	"regexp"
	"strings"
)

type Filter func(string) string

var Filters = map[string]Filter{
	"uef.isolang":       dcLanguageIso,
	"uef.peerreview":    eprintStatus,
	"uef.type":          eprintType,
	"uef.doi":           doi,
	"uef.openaire-type": openAireType,
}

// TODO: - provide generic comparison from tables
// TODO: - read comparison tables from files

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

// Try to map from SoleCRIS to ePrintTypes
func eprintType(s string) string {

	// Journal Article
	if s == "Ammatilliset aikakauslehtiartikkelit" ||
		s == "Muut aikakauslehtiartikkelit" ||
		s == "Tieteelliset aikakauslehtiartikkelit" {
		return "http://purl.org/eprint/type/JournalArticle"
	}

	// Book Item
	if s == "Artikkelit tieteellisissä kokoomateoksissa" ||
		s == "Artikkelit muissa kokoomateoksissa" {
		return "http://purl.org/eprint/type/BookItem"
	}

	// Book
	if s == "Ammatilliset kirjat" ||
		s == "Tieteelliset kirjat" ||
		s == "Toimitetut ammatilliset kirjat / lehden erikoisnumerot" ||
		s == "Toimitetut  tieteelliset kirjat / lehden erikoisnumerot" ||
		s == "Yleistajuiset kirjat" {
		return "http://purl.org/eprint/type/Book"
	}

	// Thesis
	if s == "Väitöskirjat" {
		return "http://purl.org/eprint/type/Thesis"
	}

	return s
}

// Try to map from SoleCRIS to types used in OpenAIRE
// See: https://www.kiwi.fi/display/Julkaisuarkistopalvelut/OpenAiren+vaatimat+muutokset
func openAireType(s string) string {

	// Article
	if s == "Tieteelliset aikakauslehtiartikkelit" ||
		s == "Muut aikakauslehtiartikkelit" ||
		s == "Ammatilliset aikakauslehtiartikkelit" ||
		s == "Artikkelit tieteellisissä kokoomateoksissa" ||
		s == "Artikkelit muissa kokoomateoksissa" {
		return "article"
	}

	// Types of theses
	if s == "Väitöskirjat" {
		return "doctoralThesis"
	}

	if s == "Lisensiaatintutkimukset" {
		return "other" // No equivalent in OpenAIRE
	}

	if s == "Pro gradu -tutkielmat tai vastaavat" {
		return "masterThesis" // SoleCRIS might include bachelor's theses here?
	}

	// Book
	if s == "Ammatilliset kirjat" ||
		s == "Tieteelliset kirjat" ||
		s == "Toimitetut ammatilliset kirjat / lehden erikoisnumerot" ||
		s == "Toimitetut  tieteelliset kirjat / lehden erikoisnumerot" ||
		s == "Yleistajuiset kirjat" {
		return "book"
	}

	return s
}

// Try to convert doi's to http://doi.org/xxx -format
// This filter is woefully underspecified and specific to
// UEF SoleCRIS conventions.
// See i.e. for http://stackoverflow.com/questions/27910/finding-a-doi-in-a-document-or-page
// http://blog.crossref.org/2015/08/doi-regular-expressions.html
// for guidance in writing better logic
func doi(s string) string {

	lc := strings.ToLower(s) // use lower case for all but simple tests

	if s == "-" { // means no doi in UEF input convention
		return ""
	}

	if strings.HasPrefix(lc, "http://") { // if it's http, leave it as it is
		return s
	}

	if strings.HasPrefix(lc, "doi:") { // add http-prefix, trust that the rest is correct doi
		return "http://doi.org/" + s[4:]
	}

	// try doi-regexp (does work for majority of cases, but not for all of them
	if matched, _ := regexp.MatchString(`^10.\d{4,9}/[-._;()/:A-Za-z0-9]+$`, s); matched {
		return "http://doi.org/" + s
	}

	// otherwise don't do anything
	return s

}
