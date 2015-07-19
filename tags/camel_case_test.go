package tags

import "testing"

func TestCamelCase(t *testing.T) {
	var tests = []struct {
		in, out string
	}{
		{"TestOne", "testOne"},
		{"testTwo", "testTwo"},
		{"ID", "id"},
		{"IDTest", "idTest"},
		{"ItemID", "itemID"},
	}

	for i := range tests {
		if s := TitleToCamel(tests[i].in); s != tests[i].out {
			t.Errorf("expected: %s, got: %s", tests[i].out, s)
		}
	}

}
