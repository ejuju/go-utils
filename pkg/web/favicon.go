package web

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
)

func generateMonochromeImage(size int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func ServeSimpleFavicon(c color.Color) http.HandlerFunc {
	// Generate 16x16 pixel monochrome image and encode as PNG
	buf := &bytes.Buffer{}
	err := png.Encode(buf, generateMonochromeImage(16, c))
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}
