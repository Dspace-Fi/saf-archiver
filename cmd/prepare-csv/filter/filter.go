package filter

// Filter is a function that takes a string and
// returns a string. It doesn't actually do
// filtering by transforming
type Filter func(string) string

//
// Registered filters
//
var Filters = map[string]Filter{
	"uef.isolang":       uefDcLanguageIso,
	"uef.peerreview":    uefEprintStatus,
	"uef.type":          uefEprintType,
	"uef.doi":           uefDoi,
	"uef.openaire-type": uefOpenAireType,
}
