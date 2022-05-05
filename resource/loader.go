package resource

import (
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Loader is used to load and cache game resources like images and audio files.
type Loader struct {
	// OpenAssetFunc is used to open an asset resource identified by its path.
	// The returned resource will be closed after it will be loaded.
	OpenAssetFunc func(path string) io.ReadCloser

	ImageRegistry ImageRegistry
	AudioRegistry AudioRegistry
	FontRegistry  FontRegistry

	wavDecoder wavDecoder
	oggDecoder oggDecoder

	images map[ImageID]*ebiten.Image
	wavs   map[AudioID]*wav.Stream
	oggs   map[AudioID]*vorbis.Stream
	fonts  map[FontID]font.Face
}

type wavDecoder interface {
	DecodeWAV(r io.Reader) (*wav.Stream, error)
}

type oggDecoder interface {
	DecodeOGG(r io.Reader) (*vorbis.Stream, error)
}

func NewLoader(wd wavDecoder, od oggDecoder) *Loader {
	l := &Loader{
		images:     make(map[ImageID]*ebiten.Image),
		wavs:       make(map[AudioID]*wav.Stream),
		oggs:       make(map[AudioID]*vorbis.Stream),
		fonts:      make(map[FontID]font.Face),
		wavDecoder: wd,
		oggDecoder: od,
	}
	l.AudioRegistry.mapping = make(map[AudioID]Audio)
	l.ImageRegistry.mapping = make(map[ImageID]Image)
	l.FontRegistry.mapping = make(map[FontID]Font)
	return l
}

func (l *Loader) PreloadImage(id ImageID) {
	l.LoadImage(id)
}

func (l *Loader) PreloadAudio(id AudioID) {
	audioInfo := l.GetAudioInfo(id)
	if strings.HasSuffix(audioInfo.Path, ".ogg") {
		l.LoadOGG(id)
	} else {
		l.LoadWAV(id)
	}
}

func (l *Loader) PreloadWAV(id AudioID) {
	l.LoadWAV(id)
}

func (l *Loader) PreloadOGG(id AudioID) {
	l.LoadOGG(id)
}

func (l *Loader) PreloadFont(id FontID) {
	l.LoadFont(id)
}

func (l *Loader) LoadWAV(id AudioID) *wav.Stream {
	stream, ok := l.wavs[id]
	if !ok {
		wavInfo := l.GetAudioInfo(id)
		r := l.OpenAssetFunc(wavInfo.Path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q wav reader: %v", wavInfo.Path, err))
			}
		}()
		var err error
		stream, err = l.wavDecoder.DecodeWAV(r)
		if err != nil {
			panic(fmt.Sprintf("decode %q wav: %v", wavInfo.Path, err))
		}
		l.wavs[id] = stream
	}
	return stream
}

func (l *Loader) GetAudioInfo(id AudioID) Audio {
	info, ok := l.AudioRegistry.mapping[id]
	if !ok {
		panic(fmt.Sprintf("unregistered audio with id=%d", id))
	}
	return info
}

func (l *Loader) LoadOGG(id AudioID) *vorbis.Stream {
	stream, ok := l.oggs[id]
	if !ok {
		oggInfo := l.GetAudioInfo(id)
		r := l.OpenAssetFunc(oggInfo.Path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q ogg reader: %v", oggInfo.Path, err))
			}
		}()
		var err error
		stream, err = l.oggDecoder.DecodeOGG(r)
		if err != nil {
			panic(fmt.Sprintf("decode %q ogg: %v", oggInfo.Path, err))
		}
		l.oggs[id] = stream
	}
	return stream
}

func (l *Loader) LoadFont(id FontID) font.Face {
	ff, ok := l.fonts[id]
	if !ok {
		fontInfo, ok := l.FontRegistry.mapping[id]
		if !ok {
			panic(fmt.Sprintf("unregistered font with id=%d", id))
		}
		r := l.OpenAssetFunc(fontInfo.Path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q font reader: %v", fontInfo.Path, err))
			}
		}()
		fontData, err := io.ReadAll(r)
		if err != nil {
			panic(fmt.Sprintf("reading %q data: %v", fontInfo.Path, err))
		}
		tt, err := opentype.Parse(fontData)
		if err != nil {
			panic(fmt.Sprintf("parsing %q font: %v", fontInfo.Path, err))
		}
		face, err := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    float64(fontInfo.Size),
			DPI:     96,
			Hinting: font.HintingFull,
		})
		if err != nil {
			panic(fmt.Sprintf("creating a font face for %q: %v", fontInfo.Path, err))
		}
		l.fonts[id] = face
	}
	return ff
}

func (l *Loader) LoadImage(id ImageID) *ebiten.Image {
	img, ok := l.images[id]
	if !ok {
		imageInfo, ok := l.ImageRegistry.mapping[id]
		if !ok {
			panic(fmt.Sprintf("unregistered image with id=%d", id))
		}
		r := l.OpenAssetFunc(imageInfo.Path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q image reader: %v", imageInfo.Path, err))
			}
		}()
		rawImage, _, err := image.Decode(r)
		if err != nil {
			panic(fmt.Sprintf("decode %q image: %v", imageInfo.Path, err))
		}
		img = ebiten.NewImageFromImage(rawImage)
		l.images[id] = img
	}
	return img
}
