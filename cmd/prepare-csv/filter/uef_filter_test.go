package filter

import "testing"

func TestDoi(t *testing.T) {

	var tests = [][]string{
		{"doi:10.1000/xxxx", "http://doi.org/10.1000/xxxx"},
		{"http://doi.org/10.1000/xxxx", "http://doi.org/10.1000/xxxx"},
		{"10.1000/xxxx", "http://doi.org/10.1000/xxxx"},
	}

	for _, pair := range tests {
		res := uefDoi(pair[0])
		if res != pair[1] {
			t.Error("doi: ", pair[0],
				"expected: ", pair[1],
				"got: ", res)
		}
	}
}
