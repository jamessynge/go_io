package goioutil

import (
	"unicode"
	"unicode/utf8"

	"github.com/golang/glog"
)

func MayBeTruncatedUtf8(buf []byte) bool {
	// Max byte length of a UTF-8 code point is 6 bytes, so it can only be
	// truncated if the length is 5 or less.
	length := len(buf)
	if length == 0 || length > 5 {
		glog.Infof("MayBeTruncatedUtf8: Invalid trailing length of %d", length)
		return false
	}
	// All bytes after the first must have bits matching "10xxxxxx".
	for i := 1; i < length; i++ {
		b := buf[i]
		if 0x80 != (0xC0 & b) {
			glog.Infof("MayBeTruncatedUtf8: invalid trailing byte[%d] %x", i, b)
			return false
		}
	}
	// How long "should" this code point's encoding be?
	var b byte = buf[0]
	size := 0
	for ; b >= 128; size++ {
		b <<= 1
	}
	glog.Infof("MayBeTruncatedUtf8: expected rune length %d, in %d trailing bytes", size, length)
	return length < size && size <= 6
}

// Do the bytes b contain only graphic ASCII or only graphic UTF-8, including
// white space characters.
func IsGraphicAsciiOrUtf8(b []byte) (isGraphicAscii, isGraphicUtf8 bool) {
	// Assume it is OK, until we prove otherwise.
	isGraphicAscii = true
	isGraphicUtf8 = true
	sawUtf8 := false
	sawAscii := false

	length := len(b)
	for i := 0; i < length; {
		if b[i] >= 128 {
			isGraphicAscii = false
			rest := b[i:]
			r, size := utf8.DecodeRune(rest)
			if r == utf8.RuneError {
				// If we're right at the end of the buffer, it is possible that the
				// last UTF-8 code point is truncated. If so, ignore the error.
				glog.Infof("IsGraphicAsciiOrUtf8: invalid rune at offset %d: %q", i, b[i:])
				if MayBeTruncatedUtf8(rest) {
					break
				}
				isGraphicUtf8 = false
				return
			}
			sawUtf8 = true
			if !unicode.IsGraphic(r) && !unicode.IsSpace(r) {
				glog.Infof("IsGraphicAsciiOrUtf8: multi-byte codepoint %X (%q) at offset %d is not graphic", r, b[i:i+size], i)
				isGraphicUtf8 = false
				return
			}
			i += size
			continue
		}
		sawAscii = true
		r := rune(b[i])
		if !unicode.IsGraphic(r) && !unicode.IsSpace(r) {
			glog.Infof("IsGraphicAsciiOrUtf8: byte %X (%q) at offset %d is not graphic", r, b[i:i+1], i)
			isGraphicAscii = false
			isGraphicUtf8 = false
			return
		}
		i++
	}

	if sawUtf8 || !sawAscii {
		isGraphicAscii = false
	} else if sawAscii {
		isGraphicUtf8 = false
	}
	return
}
