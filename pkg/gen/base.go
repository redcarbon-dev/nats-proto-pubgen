package gen

import (
	"bytes"
	_ "embed"
	"go/format"
	"text/template"
)

func RenderBaseFile() (string, error) {
	tmpl, err := template.New("base").Parse(base)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, nil)
	if err != nil {
		return "", err
	}

	formattedOut, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedOut), nil
}

//go:embed base.txt
var base string
