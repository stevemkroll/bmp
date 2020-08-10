package bmp

import (
	"os"
	"testing"
)

func TestReadBMP(t *testing.T) {

	file, err := os.Open("images/8bit.bmp")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n%+v\n", file)

}
