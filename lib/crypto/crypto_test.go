package crypto

import "testing"

func TestHash(t *testing.T) {
	m := Hash([]byte("1"), []byte("0"))
	if m == nil {
		t.Error()
	}

}
