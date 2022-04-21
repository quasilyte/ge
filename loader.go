package ge

import (
	"fmt"
	"image"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type Loader struct {
	// OpenAssetFunc is used to open an asset resource identified by its path.
	// The returned resource will be closed after it will be loaded.
	OpenAssetFunc func(path string) io.ReadCloser

	audio *Audio

	context *Context
	images  map[string]*ebiten.Image
	wavs    map[string]*wav.Stream
}

func NewLoader(ctx *Context) *Loader {
	return &Loader{
		context: ctx,
		images:  make(map[string]*ebiten.Image),
		wavs:    make(map[string]*wav.Stream),
	}
}

func (l *Loader) PreloadImage(path string) {
	l.LoadImage(path)
}

func (l *Loader) PreloadWAV(path string) {
	l.LoadWAV(path)
}

func (l *Loader) LoadWAV(path string) *wav.Stream {
	stream, ok := l.wavs[path]
	if !ok {
		r := l.OpenAssetFunc(path)
		defer func() {
			if err := r.Close(); err != nil {
				l.context.OnCriticalError(fmt.Errorf("closing %q wav reader: %w", path, err))
			}
		}()
		var err error
		stream, err = wav.Decode(l.audio.audioContext, r)
		if err != nil {
			panic(fmt.Sprintf("decode %q wav: %v", path, err))
		}
		l.wavs[path] = stream
	}
	return stream
}

func (l *Loader) LoadImage(path string) *ebiten.Image {
	img, ok := l.images[path]
	if !ok {
		r := l.OpenAssetFunc(path)
		defer func() {
			if err := r.Close(); err != nil {
				l.context.OnCriticalError(fmt.Errorf("closing %q image reader: %w", path, err))
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
