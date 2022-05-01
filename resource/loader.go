package resource

import (
	"fmt"
	"image"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// Loader is used to load and cache game resources like images and audio files.
type Loader struct {
	// OpenAssetFunc is used to open an asset resource identified by its path.
	// The returned resource will be closed after it will be loaded.
	OpenAssetFunc func(path string) io.ReadCloser

	wavDecoder wavDecoder
	oggDecoder oggDecoder

	images map[string]*ebiten.Image
	wavs   map[string]*wav.Stream
	oggs   map[string]*vorbis.Stream
}

type wavDecoder interface {
	DecodeWAV(r io.Reader) (*wav.Stream, error)
}

type oggDecoder interface {
	DecodeOGG(r io.Reader) (*vorbis.Stream, error)
}

func NewLoader(wd wavDecoder, od oggDecoder) *Loader {
	return &Loader{
		images:     make(map[string]*ebiten.Image),
		wavs:       make(map[string]*wav.Stream),
		oggs:       make(map[string]*vorbis.Stream),
		wavDecoder: wd,
		oggDecoder: od,
	}
}

func (l *Loader) PreloadImage(path string) {
	l.LoadImage(path)
}

func (l *Loader) PreloadWAV(path string) {
	l.LoadWAV(path)
}

func (l *Loader) PreloadOGG(path string) {
	l.LoadOGG(path)
}

func (l *Loader) LoadWAV(path string) *wav.Stream {
	stream, ok := l.wavs[path]
	if !ok {
		r := l.OpenAssetFunc(path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q wav reader: %v", path, err))
			}
		}()
		var err error
		stream, err = l.wavDecoder.DecodeWAV(r)
		if err != nil {
			panic(fmt.Sprintf("decode %q wav: %v", path, err))
		}
		l.wavs[path] = stream
	}
	return stream
}

func (l *Loader) LoadOGG(path string) *vorbis.Stream {
	stream, ok := l.oggs[path]
	if !ok {
		r := l.OpenAssetFunc(path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q ogg reader: %v", path, err))
			}
		}()
		var err error
		stream, err = l.oggDecoder.DecodeOGG(r)
		if err != nil {
			panic(fmt.Sprintf("decode %q ogg: %v", path, err))
		}
		l.oggs[path] = stream
	}
	return stream
}

func (l *Loader) LoadImage(path string) *ebiten.Image {
	img, ok := l.images[path]
	if !ok {
		r := l.OpenAssetFunc(path)
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
		l.images[path] = img
	}
	return img
}
