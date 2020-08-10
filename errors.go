package bmp

import "errors"

// Errors
var (
	ErrGeneric          = errors.New("error")
	Err02NotSupported   = errors.New("error: 2 bit format not supported")
	Err16NotSupported   = errors.New("error: 16 bit format not supported")
	Err32NotSupported   = errors.New("error: 32 bit format not supported")
	ErrBadHeader        = errors.New("error: file header")
	ErrBadMagic         = errors.New("error: magic")
	ErrCantHappen       = errors.New("error: cant happen")
	ErrEmptyBitmap      = errors.New("error: empty")
	ErrNoDelta          = errors.New("error: no delta")
	ErrOS21NotSupported = errors.New("error: OS2 v1 format not supported")
	ErrOS22NotSupported = errors.New("error: OS2 v2 format not supported")
	ErrShort            = errors.New("error: file too short")
	ErrV4NotSupported   = errors.New("error: V4 format not supported")
	ErrV5NotSupported   = errors.New("error: V5 format not supported")
)
