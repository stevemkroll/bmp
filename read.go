package bmp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// Read reads the image and returns a bmp struct or an error
func Read(r io.Reader) (b *BMP, err error) {

	// read bytes
	allBmpBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return b, err
	}

	// check length
	if len(allBmpBytes) < 54 {
		return nil, ErrShort
	}

	// check format signature
	if bytes.Compare(allBmpBytes[0:2], MagicBytes) != 0 {
		return nil, ErrBadMagic
	}

	// check header size
	header := allBmpBytes[14:18]
	headerSize := 0

	for i := range header {
		headerSize += int(header[i])
	}

	switch headerSize {
	case 12:
		fmt.Printf("%+v\n", "BITMAPCOREHEADER / OS21XBITMAPHEADER")
		break
	case 16:
	case 64:
		fmt.Printf("%+v\n", "OS22XBITMAPHEADER")
		break
	case 40:
		fmt.Printf("%+v\n", "BITMAPINFOHEADER")
		break
	case 52:
		fmt.Printf("%+v\n", "BITMAPV2INFOHEADER")
		break
	case 56:
		fmt.Printf("%+v\n", "BITMAPV3INFOHEADER")
		break
	case 108:
		fmt.Printf("%+v\n", "BITMAPV4HEADER")
		break
	case 124:
		fmt.Printf("%+v\n", "BITMAPV5HEADER")
		break
	default:
		return nil, ErrGeneric
	}

	// // since we have good header we proceed to pick it apart, checking
	// // for sanity/consistency as we go
	// copy(bmpFileHdr.Magic[0:2], bmpBytes[0:2])
	// bmpFileHdr.Size = Uint32FromLSBytes(bmpBytes[2:6])
	// // check expected file length against what we actually read in
	// if bmpFileHdr.Size != uint32(len(allBmpBytes)) {
	// 	fmt.Printf("header says file contains %d bytes, actually read %d\n",
	// 		bmpFileHdr.Size, len(allBmpBytes))
	// 	return nil, ErrShort
	// }
	// if bytes.Compare(bmpBytes[6:10], []byte{0, 0, 0, 0}) != 0 {
	// 	log.Printf("bmp: nonzero bytes in reserved area\n")
	// 	return nil, ErrBadHeader
	// }
	// // for whatever reason using mapOffset didn't work out
	// // I don't know what the offset is from.
	// // Apparently not from the start of the file.  In order to find the
	// // bitmap we captured the size of the colortable and added it to the end of
	// // the DIBHeader.  That worked.
	// mapOffset := Uint32FromLSBytes(bmpBytes[10:14])
	// bmpFileHdr.OffsetBits = mapOffset
	// bmpInfoHdr.Width = Int32FromLSBytes(bmpBytes[18:22])
	// // // verbose.Printf("width(%d)\n", int(bmpInfoHdr.Width))
	// // I think neg width is always an err
	// if bmpInfoHdr.Width <= 0 {
	// 	fmt.Printf("bmp: width <= 0 ; found %d\n", bmpInfoHdr.Width)
	// 	return nil, ErrGeneric
	// }
	// bmpInfoHdr.Height = Int32FromLSBytes(bmpBytes[22:26])
	// // // verbose.Printf("height(%d)\n", int(bmpInfoHdr.Height))
	// // OS2 can have inverted map so neg height is not an error
	// // if bmpInfoHdr.Height < 0 {
	// // 	// // verbose.Printf("top->down pixel order found (normal is bottom->up)\n")
	// // }
	// // planes is always 1
	// bmpInfoHdr.Planes = Uint16FromLSBytes(bmpBytes[26:28])
	// if bmpInfoHdr.Planes != 1 {
	// 	// log.Printf("bmp: Bad number of planes, must be 1 but found %d\n", bmpInfoHdr.biPlanes)
	// 	return nil, ErrBadHeader
	// }
	// bmpInfoHdr.BitsPerPixel = Uint16FromLSBytes(bmpBytes[28:30])

	// switch bmpInfoHdr.BitsPerPixel {
	// case 1:
	// 	// working
	// case 2:
	// 	fmt.Printf("2 bit per pixel not supported (Windows CE only)\n")
	// 	return nil, Err02NotSupported
	// case 4:
	// 	// working
	// case 8:
	// 	// working
	// case 16:
	// 	fmt.Printf("16 bit per pixel not supported\n")
	// 	return nil, Err16NotSupported
	// case 24:
	// 	// working
	// case 32:
	// 	// fmt.Printf("32 bit per pixel not supported\n")
	// 	return nil, Err32NotSupported
	// default:
	// 	// fmt.Printf("bmp: bad number of bits per pixel, must be 1/2/4/8/16/24/32 but got %d\n", bmpInfoHdr.BitsPerPixel)
	// 	return nil, ErrBadHeader
	// }

	// bmpInfoHdr.Compression = Uint32FromLSBytes(bmpBytes[30:34])

	// switch bmpInfoHdr.Compression {
	// case 0:
	// 	// uncompressed - working
	// case 1:
	// 	// RLE-8 - working
	// case 2:
	// 	// RLE-8 - testing now
	// case 3:
	// 	// fmt.Printf("bmp: Compression.BitFields is not handled\n")
	// 	return nil, ErrGeneric
	// default:
	// 	// fmt.Printf("bmp: compression value is not recognized - found (%d)\n", bmpInfoHdr.Compression)
	// 	return nil, ErrGeneric
	// }

	// bmpInfoHdr.ImageSize = Uint32FromLSBytes(bmpBytes[34:38])
	// bmpInfoHdr.XPixelsPerMeter = Int32FromLSBytes(bmpBytes[38:42])
	// bmpInfoHdr.YPixelsPerMeter = Int32FromLSBytes(bmpBytes[42:46])
	// bmpInfoHdr.Colors = Uint32FromLSBytes(bmpBytes[46:50])
	// bmpInfoHdr.Important = Uint32FromLSBytes(bmpBytes[50:54])

	// numQuads := uint32((mapOffset - (bmpInfoHdr.HdrSize + 14)) >> 2) // /= 4
	// bmpBytes = make([]byte, numQuads*4)
	// copy(bmpBytes, allBmpBytes[bmpInfoHdr.HdrSize+14:bmpInfoHdr.HdrSize+14+numQuads*4])

	// bf.ColorTable = make(color.Palette, numQuads)
	// switch bmpInfoHdr.BitsPerPixel {
	// case 1, 2, 4, 8:
	// 	for i := range bf.ColorTable {
	// 		if uint32(i) >= numQuads {
	// 			break
	// 		}
	// 		// BMP images are stored in BGR order rather than RGB order.
	// 		// Every 4th byte is padding  (bmp source was padded with zero)
	// 		bf.ColorTable[i] = color.RGBA{bmpBytes[4*i+2], bmpBytes[4*i+1], bmpBytes[4*i+0], 0xFF}
	// 	}
	// case 16, 24, 32: // color table is empty
	// }
	// mapSize := bmpInfoHdr.Width * bmpInfoHdr.Height
	// switch bmpInfoHdr.BitsPerPixel {
	// case 1:
	// 	mapSize >>= 3
	// case 4:
	// 	mapSize >>= 1
	// case 8:
	// 	// mapSize = mapSize
	// case 16:
	// 	// not implemented
	// 	return nil, ErrCantHappen
	// case 24:
	// 	mapSize *= 3
	// case 32:
	// 	// not implemented
	// 	return nil, ErrCantHappen
	// }
	// if bmpInfoHdr.ImageSize != 0 {
	// 	bf.ImageData = make([]byte, bmpInfoHdr.ImageSize)
	// 	copy(bf.ImageData, allBmpBytes[bmpInfoHdr.HdrSize+14+numQuads*4:])
	// 	n := len(bf.ImageData)
	// 	if uint32(n) != bmpInfoHdr.ImageSize {
	// 		log.Printf("bmp: bad copy - expected %d bytes got %d\n", bmpInfoHdr.ImageSize, n)
	// 		return nil, ErrShort
	// 	}
	// } else {
	// 	bf.ImageData = make([]byte, mapSize)
	// 	copy(bf.ImageData, allBmpBytes[bmpInfoHdr.HdrSize+14+numQuads*4:])
	// 	n := len(bf.ImageData)
	// 	if int32(n) != mapSize {
	// 		log.Printf("bmp: bad copy - expected %d bytes got %d\n", bmpInfoHdr.ImageSize, n)
	// 		return nil, ErrShort
	// 	}

	// }

	// // copy our loose header elements into the struct we're returning
	// // copy has to occur after they're fully built out
	// bf.FileHeader = bmpFileHdr
	// bf.DIBHeader = bmpInfoHdr

	// return &bf, err
	return b, nil
}
