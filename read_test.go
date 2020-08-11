package bmp

import (
	"os"
	"testing"
)

var images = []string{
	"1bit.bmp",
	"1bitcolor.bmp",
	"4bit.bmp",
	"4bitcompressed.bmp",
	"8bit.bmp",
	"8bitcompressed.bmp",
	"8bitgray.bmp",
	"24bit.bmp",
}

func TestReadBMP(t *testing.T) {
	for i := range images {
		path := "images/" + images[i]
		file, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		_, err = Read(file)
		if err != nil {
			t.Fatal(err)
		}
	}
}
