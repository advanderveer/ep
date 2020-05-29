package coding

import (
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"
)

func TestTemplateExecError(t *testing.T) {
	tmpl := template.Must(template.New("root").Parse(`{{ .Bar }}`))

	w := httptest.NewRecorder()
	e := NewTemplate(tmpl).Encoder(w)

	err := e.Encode(output1{})
	if err == nil || !strings.Contains(err.Error(), "can't evaluate field Bar") {
		t.Fatalf("expected specific error, got: %v", err)
	}
}
