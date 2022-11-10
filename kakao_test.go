package kakaogo_test

import (
	"testing"

	"github.com/yms2772/kakaogo"
)

func TestNew(t *testing.T) {
	if _, err := kakaogo.New("example@email.com", "pAsSwORd"); err != nil {
		t.Errorf("%+v\n", err)
	}
}
