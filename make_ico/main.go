package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	_ "image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	f, err := os.Open("Icon.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	sizes := []uint{256, 128, 64, 48, 32, 16}

	type icoEntry struct {
		data   []byte
		width  byte
		height byte
	}

	var entries []icoEntry
	for _, sz := range sizes {
		resized := resize.Resize(sz, sz, img, resize.Lanczos3)
		var buf bytes.Buffer
		if err := png.Encode(&buf, resized); err != nil {
			panic(err)
		}
		w := byte(sz)
		h := byte(sz)
		if sz == 256 {
			w = 0
			h = 0
		}
		entries = append(entries, icoEntry{data: buf.Bytes(), width: w, height: h})
	}

	out, err := os.Create("Icon.ico")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// ICONDIR
	binary.Write(out, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(out, binary.LittleEndian, uint16(1)) // image type: 1=ico
	binary.Write(out, binary.LittleEndian, uint16(len(entries)))

	// First entry offset = 6 (header) + 16 (per entry) * count
	offset := uint32(6 + 16*len(entries))
	for _, e := range entries {
		out.Write([]byte{e.width, e.height, 0, 0}) // w, h, colorCount, reserved
		binary.Write(out, binary.LittleEndian, uint16(1))        // planes
		binary.Write(out, binary.LittleEndian, uint16(32))       // bpp
		binary.Write(out, binary.LittleEndian, uint32(len(e.data))) // size
		binary.Write(out, binary.LittleEndian, offset)            // offset
		offset += uint32(len(e.data))
	}

	for _, e := range entries {
		out.Write(e.data)
	}
}
