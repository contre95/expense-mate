package health

import "testing"

func TestHealthChecker(t *testing.T) {
	h := NewService()
	resp := h.Ping()
	if resp != "pong" {
		t.Errorf(" Response should be 'pong' got %s instead.", resp)
	}
}
