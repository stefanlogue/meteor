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

func TestCheckBoardMatchesBranch(t *testing.T) {
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
		{"it should match different case", "ticket-1234", true},
		{"it should match different format", "fix-for-TICKET-1234", true},
		{"it should not match with no digits", "TICKET-", false},
	}

	for _, tc := range cases {
		t.Run(tc.Desc, func(t *testing.T) {
			got := checkBoardMatchesBranch("TICKET", tc.msg)
			assertEqualBools(t, tc.want, got)
		})
	}
}

func TestGetTicketNumberFromString(t *testing.T) {
	cases := []struct {
		Desc string
		msg  string
		sub  string
		want string
	}{
		{"it should return the ticket number", "TICKET-1234", "TICKET", "TICKET-1234"},
		{"it should return when ticket is not at the beginning", "this is a TICKET-1234", "TICKET", "TICKET-1234"},
	}

	for _, tc := range cases {
		t.Run(tc.Desc, func(t *testing.T) {
			got := getTicketNumberFromString(tc.msg, tc.sub)
			assertEqualStrings(t, tc.want, got)
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
