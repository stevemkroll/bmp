package bmp

import (
	"bufio"
	"image"
	"image/color"
	"io"
	"log"
)

func unwindRLE8(r io.Reader, b *BMP) ([]byte, error) {

	maxReadBytes := len(b.ImageData)
	rowWidth := b.DIBHeader.Width
	if (rowWidth % 4) != 0 {
		rowWidth++
	}
	pixMap := make([]byte, 0, b.DIBHeader.Height*rowWidth)

	br := bufio.NewReader(r)
	bytesRead := 0
	lineCt := 0

	for {
		if len(pixMap) == cap(pixMap) {
			break
		}
		if bytesRead >= maxReadBytes {
			break
		}
		numPix, err := br.ReadByte()
		if err != nil {

			log.Printf("bmp: bad read in RLE8\n")
			return nil, err
		} else {
			bytesRead++
		}

		pixVal, err := br.ReadByte()
		if err != nil {

			log.Printf("bmp: bad read in RLE8\n")
			return nil, err
		} else {
			bytesRead++
		}

		if numPix > 0 {
			for x := 0; x < int(numPix); x++ {
				pixMap = append(pixMap, pixVal)
			}
			continue
		} else {
			if inRangeByte(0, pixVal, 2) {
				switch pixVal {
				case 0:
					for {
						if (len(pixMap) % 4) == 0 {
							break
						}

						pixMap = append(pixMap, 0)
					}

					lineCt++

					continue
				case 1:

					lineCt++

					goto xit
				case 2:
					log.Printf("Delta value found but no handler available for it\n")
					return nil, ErrNoDelta
					deltax, err := br.ReadByte()
					deltay, err := br.ReadByte()

					deltax = deltax
					deltay = deltay
					err = err
					bytesRead += 2

				}
				log.Printf("can't happen\n")
				return nil, ErrCantHappen
			}
			numPix = pixVal

			for i := 0; i < int(numPix); i++ {
				pixVal, err := br.ReadByte()
				if err != nil {

					log.Printf("bmp: bad read in RLE8\n")
					return nil, err
				} else {
					bytesRead++
				}
				pixMap = append(pixMap, pixVal)
			}
			if (bytesRead % 2) != 0 {
				_, err := br.ReadByte()
				if err != nil {

					log.Printf("bmp: bad read in RLE8\n")
					return nil, err
				} else {
					bytesRead++
				}
			}
		}

	}
xit:

	if len(pixMap) != cap(pixMap) {

	}
	if bytesRead != len(b.ImageData) {

	}

	for {
		if len(pixMap) >= cap(pixMap) {
			break
		}

		pixMap = append(pixMap, 0x0)
	}
	b.ImageData = pixMap
	return pixMap, nil
}

// decodePaletted8 reads an 8 bit-per-pixel BMP image from r.
func decodePaletted8(r io.Reader, c image.Config, b *BMP) (image.Image, error) {
	rowWidth := b.DIBHeader.Width
	if (rowWidth % 4) != 0 {
		rowWidth++
	}
	maxBits := len(b.ImageData)
	_ = b.DIBHeader.Height * rowWidth
	paletted := image.NewPaletted(image.Rect(0, 0, c.Width, c.Height), c.ColorModel.(color.Palette))
	bytesRead := int32(0)
	tmp := make([]byte, 1)
	for y := c.Height - 1; y >= 0; y-- {
		if bytesRead >= int32(maxBits) {
			break
		}
		p := paletted.Pix[y*paletted.Stride : y*paletted.Stride+c.Width]
		n, err := r.Read(p)
		if err != nil {
			log.Printf("bmp: bad read in Pal8\n")
			return nil, err
		}
		if n != c.Width {
		}
		bytesRead += int32(c.Width)

		if bytesRead >= int32(maxBits) {
			break
		}

		for {
			if bytesRead >= int32(maxBits) {
				break
			}
			if (bytesRead % 4) == 0 {
				break
			}
			_, err := r.Read(tmp)
			if err != nil {
				log.Printf("bmp: bad read in Pal8\n")
				return nil, err
			}
			bytesRead++

		}
	}
	return paletted, nil
}
