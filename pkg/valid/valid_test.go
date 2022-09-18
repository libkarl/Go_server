package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelativePath(t *testing.T) {
	testCases := []struct {
		label string
		path  string
		ok    bool
	}{
		{
			label: "valid relative file",
			path:  "this.txt",
			ok:    true,
		},
		{
			label: "valid relative path",
			path:  "here/this.txt",
			ok:    true,
		},
		{
			label: "absolute path",
			path:  "/this.txt",
			ok:    false,
		},
		{
			label: "unsupported relative path",
			path:  "./this.txt",
			ok:    false,
		},
		{
			label: "path too short",
			path:  "xy",
			ok:    false,
		},
		{
			label: "hidden path",
			path:  ".this.txt",
			ok:    false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.ok, RelativePath(tc.path), tc.label)
	}
}
