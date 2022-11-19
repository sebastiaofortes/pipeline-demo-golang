package hello

import (
	"testing"
)

func TestHelloGitHub(t *testing.T) {
	result := HelloGitHub()
	if result != "Hello GitHub" {
		t.Errorf("HelloGitHub() = %s; want Hello GitHub", result)
	}
}

func TestHelloFail(t *testing.T) {
	result := HelloGitHub()
	if result != "Error" {
		t.Errorf("HelloGitHub() = %s; want Hello GitHub", result)
	}
}