package web

import (
	"html/template"
)

var (
	NotFoundTmpl = template.Must(template.New("").Parse(`Not Found`))
)

type (
	NotFoundOutput struct{}
)

func (h *handler) NotFound() (out NotFoundOutput) { return }

func (_ NotFoundOutput) Template() *template.Template { return NotFoundTmpl }

func (_ NotFoundOutput) Status() int { return 404 }
