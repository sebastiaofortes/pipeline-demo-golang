package hello

import (
	"testing"
)

func TestHelloGitHub(t *testing.T) {
	result := HelloGitHub()
	if result != "Hello GitHub Actions" {
		t.Errorf("HelloGitHub() = %s; want Hello GitHub", result)
	}
}