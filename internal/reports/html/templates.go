package html

import (
	"html/template"
)

func resolveTemplates() *template.Template {
	templates, _ := template.New("").Funcs(
		template.FuncMap{
			"inc":            increment(),
			"passes":         checkStatus(),
			"ratio":          calculateRatio(),
			"formatDateTime": formatDateTime(),
			"formatDate":     formatDate(),
			"formatTime":     formatTime(),
			"toHumanTime":    toHumanTime(),
		}).ParseFS(templateFiles, "templates/*.tmpl")

	return templates
}
