package main

import "testing"

func TestBuildCommitCommand(t *testing.T) {
	msg := "test"
	body := "test body"
	expected := "commit -m test test body"
	got := buildCommitCommand(msg, body)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
