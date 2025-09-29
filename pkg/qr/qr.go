package qr

import (
	"bytes"

	qrcode "github.com/skip2/go-qrcode"
)

// PNG generates a PNG QR code for the given content at the given pixel size.
// Size is the image width/height in pixels. Reasonable values: 128–1024.
func PNG(content string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256
	}
	var buf bytes.Buffer
	png, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	png.DisableBorder = false
	if err := png.Write(size, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
