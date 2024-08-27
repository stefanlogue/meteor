package config

import "testing"

func TestConvertTemplate(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"adds breaking change marker", "@type: @message", "{{.Type}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"},
		{"converts template", "@type(@scope): @message", "{{.Type}}{{if .Scope}}({{.Scope}}){{end}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"},
		{"converts without scope", "@type: @message", "{{.Type}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := ConvertTemplate(tc.input)
			assertEqual(t, tc.want, got)
		})
	}
}
