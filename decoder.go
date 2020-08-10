package bmp

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"log"
)

// Decode decodes the image
func Decode(r io.Reader) (img image.Image, err error) {

	b, err := Read(r)
	if err != nil {
		return img, err
	}

	var c image.Config
	switch b.DIBHeader.BitsPerPixel {
	case 1, 2, 4, 8:
		c = image.Config{ColorModel: b.ColorTable, Width: int(b.DIBHeader.Width), Height: int(b.DIBHeader.Height)}
	case 16:
		return nil, Err16NotSupported
	case 24:
		c = image.Config{ColorModel: color.RGBAModel, Width: int(b.DIBHeader.Width), Height: int(b.DIBHeader.Height)}
	case 32:
		return nil, Err32NotSupported
	default:
		return nil, ErrCantHappen
	}

	if len(b.ImageData) <= 0 {
		return nil, ErrEmptyBitmap
	}

	nr := bytes.NewReader(b.ImageData)

	switch b.DIBHeader.BitsPerPixel {
	case 1:
		img, err = decodePaletted1(nr, c, b)
	case 4:
		if b.DIBHeader.Compression == CompressionRLE4 {
			pixbufadr, err := unwindRLE4(nr, b)
			if err != nil {
				log.Printf("bmp: bad read from RLE4\n")
				return nil, err
			}
			nr = bytes.NewReader(pixbufadr)
		}
		img, err = decodePaletted4(nr, c, b)
	case 8:
		if b.DIBHeader.Compression == CompressionRLE8 {
			pixbufadr, err := unwindRLE8(nr, b)
			if err != nil {
				log.Printf("bmp: bad read from RLE8\n")
				return nil, err
			}
			nr = bytes.NewReader(pixbufadr)
		}
		img, err = decodePaletted8(nr, c, b)
	case 24:
		img, err = decodeRGBA(nr, c)
	default:
		log.Printf("bmp: can't happen Decode\n") // only 1/4/8/24 allowed by earlier logic
		return nil, ErrCantHappen
	}
	return img, err
}

// DecodeConfig returns the decoder config
func DecodeConfig(r io.Reader) (config image.Config, err error) {

	bf, err := Read(r) // not efficient but simple wins : read bitmap just to git header info

	switch bf.DIBHeader.BitsPerPixel {
	case 1, 2, 4, 8:
		return image.Config{ColorModel: bf.ColorTable, Width: int(bf.DIBHeader.Width), Height: int(bf.DIBHeader.Height)}, nil
	case 16:
		// fmt.Printf("16 bit per pixel not supported\n")
		return config, Err16NotSupported
	case 24:
		// // verbose.Printf("24 colormodel=%v\n", color.RGBAModel)
		return image.Config{ColorModel: color.RGBAModel, Width: int(bf.DIBHeader.Width), Height: int(bf.DIBHeader.Height)}, nil
	case 32:
		// fmt.Printf("32 bit per pixel not supported\n")
		return config, Err32NotSupported
	default:
		// log.Printf("bmp: can't happen DecodeConfig\n") // only 1/4/8/24 allowed by earlier logic
		var noConfig image.Config // return an empty config struct with err
		return noConfig, ErrCantHappen
	}
}
