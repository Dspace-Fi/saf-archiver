package main

import "testing"

func TestEscapeEntities(t *testing.T) {

	var tests = [][]string{
		{"aa", "aa"},
		{"a & b", "a &amp; b"},
	}

	for _, pair := range tests {
		res := escapeEntities(pair[0])
		if res != pair[1] {
			t.Error("escapeEntities: ", pair[0],
				"expected", pair[1],
				"got:", res)
		}
	}
}

type dcvalueTestpair struct {
	header   string
	data     string
	expected []DCValue
}

type dcvalueNilTest struct {
	header string
	data   string
}

func TestMakeDCValues(t *testing.T) {

	// result cannot be nil
	var tests = []dcvalueTestpair{
		{"dc.title", "Frank Furter",
			[]DCValue{DCValue{Element: "title", Qualifier: "none", Value: "Frank Furter"}}},
		{"dc.title", "Frank Furter||Rey Dorado",
			[]DCValue{
				DCValue{Element: "title", Qualifier: "none", Value: "Frank Furter"},
				DCValue{Element: "title", Qualifier: "none", Value: "Rey Dorado"}}},
	}

	for _, pair := range tests {

		res := makeDCValues(pair.header, pair.data)
		if res == nil { // check nil
			t.Error("MakeDCValues:", pair.header, pair.data,
				"Expected:", pair.expected,
				"Got:", res,
			)
		} else { // check each slice member
			if len(res) != len(pair.expected) {
				t.Error("MakeDCValues: ", pair.header, pair.data,
					"Expected:", pair.expected,
					"Got:", res,
				)
			}
			for i, r := range res {
				if r != pair.expected[i] {
					t.Error("MakeDCValues:", pair.header, pair.data,
						"Expected:", pair.expected[i],
						"Got:", r,
					)
				}
			}

		}
	}

	nilTests := []dcvalueNilTest{
		{header: "invalid", data: "dummy"},
		{header: "dc.title", data: ""},
	}

	for _, tt := range nilTests {
		res := makeDCValues(tt.header, tt.data)
		if res != nil {
			t.Error("MakeDCValues", tt.header, tt.data,
				"Expected nil, got: ", res)
		}
	}
}
