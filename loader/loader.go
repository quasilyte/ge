package loader

import (
	"fmt"
	"image"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type WAVDecoder interface {
	DecodeWAV(r io.Reader) (*wav.Stream, error)
}

type Cache struct {
	// OpenAssetFunc is used to open an asset resource identified by its path.
	// The returned resource will be closed after it will be loaded.
	OpenAssetFunc func(path string) io.ReadCloser

	WAVDecoder WAVDecoder

	images map[string]*ebiten.Image
	wavs   map[string]*wav.Stream
}

func NewCache() *Cache {
	return &Cache{
		images: make(map[string]*ebiten.Image),
		wavs:   make(map[string]*wav.Stream),
	}
}

func (c *Cache) PreloadImage(path string) {
	c.GetImage(path)
}

func (c *Cache) PreloadWAV(path string) {
	c.GetWAV(path)
}

func (c *Cache) GetWAV(path string) *wav.Stream {
	stream, ok := c.wavs[path]
	if !ok {
		r := c.OpenAssetFunc(path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q wav reader: %v", path, err))
			}
		}()
		var err error
		stream, err = c.WAVDecoder.DecodeWAV(r)
		if err != nil {
			panic(fmt.Sprintf("decode %q wav: %v", path, err))
		}
		c.wavs[path] = stream
	}
	return stream
}

func (c *Cache) GetImage(path string) *ebiten.Image {
	img, ok := c.images[path]
	if !ok {
		r := c.OpenAssetFunc(path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q image reader: %v", path, err))
			}
		}()
		rawImage, _, err := image.Decode(r)
		if err != nil {
			panic(fmt.Sprintf("decode %q image: %v", path, err))
		}
		img = ebiten.NewImageFromImage(rawImage)
		c.images[path] = img
	}
	return img
}
