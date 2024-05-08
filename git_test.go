package main

import "testing"

func TestBuildCommitCommand(t *testing.T) {
	msg := "test"
	body := "test body"
	expected := "git commit -m test -m 'test body'"
	_, got := buildCommitCommand(msg, body, nil)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestMatchTicketNumber(t *testing.T) {
	cases := []struct {
		Desc string
		msg  string
		want bool
	}{
		{"it should match with 1 digit", "TICKET-1", true},
		{"it should match with 2 digits", "TICKET-12", true},
		{"it should match with 3 digits", "TICKET-123", true},
		{"it should match with 4 digits", "TICKET-1234", true},
		{"it should match with 5 digits", "TICKET-12345", true},
		{"it should match with 6 digits", "TICKET-123456", true},
	}

	for _, tc := range cases {
		t.Run(tc.Desc, func(t *testing.T) {
			got := matchTicketNumber("TICKET", tc.msg)
			assertEqualBools(t, tc.want, got)
		})
	}
}

func assertEqualStrings(t testing.TB, expected, got string) {
	t.Helper()
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func assertEqualBools(t testing.TB, expected, got bool) {
	t.Helper()
	if got != expected {
		t.Errorf("expected %t, got %t", expected, got)
	}
}
