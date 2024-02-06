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
