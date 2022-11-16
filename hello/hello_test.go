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

func TestHelloGitHub2(t *testing.T) {
	result := HelloGitHub()
	if result != "Hello GitHub" {
		t.Errorf("HelloGitHub() = %s; want Hello GitHub", result)
	}
}