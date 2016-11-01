package generator

type Generator func([]string) string

var Generators = map[string]Generator{
	"uef.dc-citation": uefDcCitation,
}
