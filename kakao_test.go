package kakaogo_test

import (
	"testing"

	kakaogo "github.com/yms2772/kakao.go"
)

func TestNew(t *testing.T) {
	if _, err := kakaogo.New("example@email.com", "pAsSwORd"); err != nil {
		t.Errorf("%+v\n", err)
	}
}
