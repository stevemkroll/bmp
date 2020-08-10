package bmp

// Compression Level Variables
var (
	CompressionRGB            uint32 = 0
	CompressionRLE8           uint32 = 1
	CompressionRLE4           uint32 = 2
	CompressionBITFIELDS      uint32 = 3
	CompressionJPEG           uint32 = 4
	CompressionPNG            uint32 = 5
	CompressionALPHABITFIELDS uint32 = 6
)

// MagicBytes represents the format signature
var MagicBytes = []byte{66, 77}
