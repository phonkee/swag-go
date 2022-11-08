package swag

import "testing"

func TestPathJoin(t *testing.T) {
	for _, item := range []struct {
		left     string
		right    string
		expected string
	}{
		//{"", "", ""},
		{"", "v1", "/v1"},
		{"/", "/", "/"},
		{"/", "/v1", "/v1"},
		{"/v1", "/", "/v1/"},
		{"/v1", "/v2", "/v1/v2"},
		{"v1", "v2", "/v1/v2"},
		{"/v1/", "/v2", "/v1/v2"},
		{"/v1/", "/v2/", "/v1/v2/"},
		{"/v1", "/v2", "/v1/v2"},
	} {
		if actual := pathJoin(item.left, item.right); actual != item.expected {
			t.Errorf("pathJoin(%q, %q) = %q, expected %q", item.left, item.right, actual, item.expected)
		}
	}

}
