package coding

import (
	"testing"
)

func TestJSON(t *testing.T) {
	var e Encoding = JSON{}
	var d Decoding = JSON{}

	if e.Produces() != "application/json" {
		t.Fatalf("unexpected, got: %v", e.Produces())
	}

	if d.Accepts() != "application/json, application/vnd.api+json" {
		t.Fatalf("unexpected, got: %v", d.Accepts())
	}
}
