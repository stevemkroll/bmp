package bmp

import "image/color"

// BMP defines the bitmap file structure
type BMP struct {
	FileHeader Header
	DIBHeader  DIB
	ColorTable color.Palette
	ImageData  []byte
}

// DIB defines the device independent bitmap header
type DIB struct {
	HdrSize         uint32
	Width           int32
	Height          int32
	Planes          uint16
	BitsPerPixel    uint16
	Compression     uint32
	ImageSize       uint32
	XPixelsPerMeter int32
	YPixelsPerMeter int32
	Colors          uint32
	Important       uint32
}

// Header defines the file header
type Header struct {
	Magic      [2]byte
	Size       uint32
	Reserved1  uint16
	Reserved2  uint16
	OffsetBits uint32
}

// Info defines the info structure
type Info struct {
	Header DIB
	Colors []uint32
}
