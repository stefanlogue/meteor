package config

import (
	"fmt"
	"strings"
)

func ConvertTemplate(t string) (string, error) {
	if !strings.Contains(t, "@type") || !strings.Contains(t, "@message") {
		return t, fmt.Errorf("template must contain @type and @message")
	}
	t = strings.Replace(t, ":", "{{if .IsBreakingChange}}!{{end}}:", 1)
	t = strings.ReplaceAll(t, "@type", "{{.Type}}")
	t = strings.ReplaceAll(t, "(@scope)", "{{if .Scope}}({{.Scope}}){{end}}")
	t = strings.ReplaceAll(t, "@ticket", "{{.TicketNumber}}")
	t = strings.ReplaceAll(t, "@message", "{{.Message}}")
	return t, nil
}
