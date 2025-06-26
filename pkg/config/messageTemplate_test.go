package config

import "testing"

func TestConvertTemplate(t *testing.T) {
	validCases := []struct {
		name  string
		input string
		want  string
	}{
		{"adds breaking change marker", "@type: @message", "{{.Type}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"},
		{"converts template", "@type(@scope): @message", "{{.Type}}{{if .Scope}}({{.Scope}}){{end}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"},
		{"converts without scope", "@type: @message", "{{.Type}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"},
	}
	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ConvertTemplate(tc.input)
			assertEqual(t, tc.want, got)
			assertIsNotError(t, err)
		})
	}

	errorCases := []struct {
		name  string
		input string
	}{
		{"invalid template", "@type({{@scope}}: @message"},
		{"must contain @type and @message", "@scope: "},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ConvertTemplate(tc.input)
			assertIsError(t, err)
		})
	}
}
