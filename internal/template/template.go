package template

import (
	"html/template"
	"log"
	"os"
)

var HtmlTemplate *template.Template

func SetupMetricTemplate() {
	bytes, err := os.ReadFile("internal/template/index.html")

	if err != nil {
		log.Fatal(err)
	}

	HtmlTemplate, err = template.New("").Parse(string(bytes))

	if err != nil {
		log.Fatal(err)
	}
}
