package auth

import "testing"

func Test_generateToken(t *testing.T) {
	token := generateToken()
	if token == "" {
		t.Errorf("expected generated token returns non-empty string; got empty")
	}
	if len(token) < 8 {
		t.Errorf("expected generated token length equal or larger than 8; got %d", len(token))
	}
	if generateToken() == generateToken() {
		t.Errorf("expected generated tokens are random")
	}
}
