package transform

import (
	"bytes"
	"strings"
	"testing"
)

func TestTransform(t *testing.T) {

	var b bytes.Buffer
	r := strings.NewReader(goodCsv)
	err := Transform(r, &b)
	if err != nil {
		t.Fatal("Sent good csv, did not expect error. Got:", err)
	}
}