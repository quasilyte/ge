package resource

import (
	"fmt"
	"image"
	"io"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Loader is used to load and cache game resources like images and audio files.
type Loader struct {
	// OpenAssetFunc is used to open an asset resource identified by its path.
	// The returned resource will be closed after it will be loaded.
	OpenAssetFunc func(path string) io.ReadCloser

	ImageRegistry  ImageRegistry
	AudioRegistry  AudioRegistry
	FontRegistry   FontRegistry
	ShaderRegistry ShaderRegistry
	RawRegistry    RawRegistry

	wavDecoder wavDecoder
	oggDecoder oggDecoder

	images  map[ImageID]Image
	shaders map[ShaderID]*ebiten.Shader
	wavs    map[AudioID][]byte
	oggs    map[AudioID]*vorbis.Stream
	fonts   map[FontID]font.Face
	raws    map[RawID][]byte
}

type wavDecoder interface {
	DecodeWAV(r io.Reader) (*wav.Stream, error)
}

type oggDecoder interface {
	DecodeOGG(r io.Reader) (*vorbis.Stream, error)
}

func NewLoader(wd wavDecoder, od oggDecoder) *Loader {
	l := &Loader{
		images:     make(map[ImageID]Image),
		shaders:    make(map[ShaderID]*ebiten.Shader),
		wavs:       make(map[AudioID][]byte),
		oggs:       make(map[AudioID]*vorbis.Stream),
		fonts:      make(map[FontID]font.Face),
		raws:       make(map[RawID][]byte),
		wavDecoder: wd,
		oggDecoder: od,
	}
	l.AudioRegistry.mapping = make(map[AudioID]Audio)
	l.ImageRegistry.mapping = make(map[ImageID]ImageInfo)
	l.ShaderRegistry.mapping = make(map[ShaderID]ShaderInfo)
	l.FontRegistry.mapping = make(map[FontID]Font)
	l.RawRegistry.mapping = make(map[RawID]Raw)
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

func (l *Loader) PreloadShader(id ShaderID) {
	l.LoadShader(id)
}

func (l *Loader) PreloadRaw(id RawID) {
	l.LoadRaw(id)
}

func (l *Loader) LoadWAV(id AudioID) []byte {
	data, ok := l.wavs[id]
	if !ok {
		wavInfo := l.GetAudioInfo(id)
		r := l.OpenAssetFunc(wavInfo.Path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q wav reader: %v", wavInfo.Path, err))
			}
		}()
		stream, err := l.wavDecoder.DecodeWAV(r)
		if err != nil {
			panic(fmt.Sprintf("decode %q wav: %v", wavInfo.Path, err))
		}
		data, err := io.ReadAll(stream)
		if err != nil {
			panic(fmt.Sprintf("read %q wav: %v", wavInfo.Path, err))
		}
		l.wavs[id] = data
	}
	return data
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
		if fontInfo.LineSpacing != 0 && fontInfo.LineSpacing != 1 {
			h := float64(face.Metrics().Height.Round()) * fontInfo.LineSpacing
			face = text.FaceWithLineHeight(face, math.Round(h))
		}
		l.fonts[id] = face
	}
	return ff
}

func (l *Loader) LoadImage(id ImageID) Image {
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
		data := ebiten.NewImageFromImage(rawImage)
		img = Image{
			ID:                 id,
			Data:               data,
			DefaultFrameWidth:  imageInfo.FrameWidth,
			DefaultFrameHeight: imageInfo.FrameHeight,
		}
		l.images[id] = img
	}
	return img
}

func (l *Loader) LoadShader(id ShaderID) *ebiten.Shader {
	shader, ok := l.shaders[id]
	if !ok {
		shaderInfo, ok := l.ShaderRegistry.mapping[id]
		if !ok {
			panic(fmt.Sprintf("unregistered shader with id=%d", id))
		}
		r := l.OpenAssetFunc(shaderInfo.Path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q shader reader: %v", shaderInfo.Path, err))
			}
		}()
		data, err := io.ReadAll(r)
		if err != nil {
			panic(fmt.Sprintf("read %q shader: %v", shaderInfo.Path, err))
		}
		rawShader, err := ebiten.NewShader(data)
		if err != nil {
			panic(fmt.Sprintf("compile %q shader: %v", shaderInfo.Path, err))
		}
		shader = rawShader
		l.shaders[id] = shader
	}
	return shader
}

func (l *Loader) LoadRaw(id RawID) []byte {
	data, ok := l.raws[id]
	if !ok {
		rawInfo, ok := l.RawRegistry.mapping[id]
		if !ok {
			panic(fmt.Sprintf("unregistered raw with id=%d", id))
		}
		r := l.OpenAssetFunc(rawInfo.Path)
		defer func() {
			if err := r.Close(); err != nil {
				panic(fmt.Sprintf("closing %q raw reader: %v", rawInfo.Path, err))
			}
		}()
		var err error
		data, err = io.ReadAll(r)
		if err != nil {
			panic(fmt.Sprintf("read %q raw: %v", rawInfo.Path, err))
		}
		l.raws[id] = data
	}
	return data
}
