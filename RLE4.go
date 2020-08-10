package bmp

import (
	"bufio"
	//"bytes"
	//"fmt"
	"image"
	"image/color"
	"io"
	"log"
	//"os"
)

func unPack2(b byte) [2]byte {
	var x [2]byte
	x[0] = (b >> 4) & 0x0f
	x[1] = b & 0xf
	return x
}

func unwindRLE4(r io.Reader, b *BMP) ([]byte, error) {
	maxReadBytes := len(b.ImageData)
	rowWidth := b.DIBHeader.Width
	if (rowWidth % 2) != 0 {
		rowWidth++
	}
	pixMap := make([]byte, 0, b.DIBHeader.Height*rowWidth/2)
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
			log.Printf("bmp: bad read in RLE4\n")
			return nil, err
		} else {
			bytesRead++
		}
		pixVal, err := br.ReadByte()
		if err != nil {
			log.Printf("bmp: bad read in RLE4\n")
			return nil, err
		} else {
			bytesRead++
		}
		if numPix > 0 {
			loopCt := numPix / 2
			loopXtra := numPix - (loopCt * 2)
			for x := 0; x < int(loopCt); x++ {
				pixMap = append(pixMap, pixVal)
			}
			if loopXtra != 0 {
				pixMap = append(pixMap, pixVal&0xf0)
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
					log.Printf("Delta value found but no delta handler available\n")
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
			loopCt := numPix / 2
			loopXtra := numPix - (loopCt * 2)
			for x := 0; x < int(loopCt); x++ {
				pixVal, err := br.ReadByte()
				if err != nil {
					log.Printf("bmp: bad read in RLE4\n")
					return nil, err
				} else {
					bytesRead++
				}
				pixMap = append(pixMap, pixVal)
			}
			if loopXtra != 0 {
				pixMap = append(pixMap, pixVal&0xf0)
			}

			if (bytesRead % 2) != 0 {
				_, err := br.ReadByte()
				if err != nil {

					log.Printf("bmp: bad read in RLE4\n")
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

// decodePaletted reads a 4 bit-per-pixel BMP image from r.
func decodePaletted4(r io.Reader, c image.Config, b *BMP) (image.Image, error) {
	maxBits := len(b.ImageData)
	// verbose.Printf("Entry to decodePaletted4\n")
	paletted := image.NewPaletted(image.Rect(0, 0, c.Width, c.Height), c.ColorModel.(color.Palette))
	br := bufio.NewReader(r)
	var bytesRead int
	// verbose.Printf("Height = %d   Width = %d  H*W = %d\n", c.Height, c.Width, c.Height*c.Width)
	// verbose.Printf("maxBits = %d\n", maxBits)
	// verbose.Printf("paletted.Stride(%d)\n", paletted.Stride)
	lastPix := c.Height * c.Width
	rowWidth := c.Width / 2
	rowXtra := c.Width - (rowWidth * 2)
	// N.B. BMP images are stored bottom-up rather than top-down, left to right
	for y := c.Height - 1; y >= 0; y-- {
		var pix2 byte
		var err error
		var start, finish int
		for x := 0; x < rowWidth*2; x += 2 {
			if bytesRead >= maxBits {
				break
			}
			pix2, err := br.ReadByte()
			if err != nil {
				log.Printf("bmp: bad read in Pal4\n")
				return nil, err
			}
			bytesRead++
			b := unPack2(pix2)
			start := x + (y * c.Width)
			finish := start + 2
			if finish > lastPix {
				finish = lastPix
			}
			if start > lastPix {
				start = lastPix
			}
			//// verbose.Printf("start(%d) finish(%d) byte[%d]  %v\n", start, finish, bytesRead, b)
			copy(paletted.Pix[start:finish], b[:])
		}
		// last byte of scanline may not have all bits used so piece it out
		if rowXtra != 0 {
			//// verbose.Printf("adding last pixel to line\n")
			pix2, err = br.ReadByte()
			if err != nil {
				log.Printf("bmp: bad read in Pal4\n")
				return nil, err
			}
			bytesRead++
			b := unPack2(pix2)
			start += 2
			finish = start + rowXtra
			if finish > lastPix {
				// verbose.Printf("LastPix\n")
				finish = lastPix
			}
			if start > lastPix {
				start = lastPix
			}
			//// verbose.Printf("+start(%d) finish(%d) byte[%d]  %v\n", start, finish, bytesRead, b)
			copy(paletted.Pix[start:finish], b[:rowXtra])
		}
		// scanlines are padded if necessary to multiple of uint32 (DWORD)
		for {
			if (bytesRead % 4) == 0 {
				break
			}
			pix2, err = br.ReadByte()
			if err != nil {
				log.Printf("bmp: bad read in Pal4\n")
				return nil, err
			}
			bytesRead++
			// verbose.Printf("byte[%d]\n", bytesRead)
		}

	}
	return paletted, nil
}
