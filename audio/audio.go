package audio

import (
	"fmt"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/quasilyte/ge/loader"
)

type RegistryMap map[ResourceID]ResourceInfo

type ResourceID uint32

type ResourceInfo struct {
	Path string

	// Volume adjust how loud this sound will be.
	// The default value of 0 means "unadjusted".
	// Value greated than 0 increases the volume, negative values decrease it.
	// This setting accepts values in [-1, 1] range, where -1 mutes the sound
	// while 1 makes it as loud as possible.
	Volume float64
}

type AudioRegistry struct {
	mapping RegistryMap
}

func (r *AudioRegistry) Set(id ResourceID, info ResourceInfo) {
	r.mapping[id] = info
}

func (r *AudioRegistry) Assign(m RegistryMap) {
	r.mapping = m
}

type System struct {
	Registry AudioRegistry

	cache *loader.Cache

	audioContext *audio.Context
	wavResources map[ResourceID]*audioResource
}

type audioResource struct {
	player *audio.Player
	volume float64
}

func (sys *System) Init(c *loader.Cache) {
	sys.cache = c
	sys.audioContext = audio.NewContext(32000)
	sys.Registry.mapping = make(RegistryMap)
	sys.wavResources = make(map[ResourceID]*audioResource)
}

func (sys *System) DecodeWAV(r io.Reader) (*wav.Stream, error) {
	return wav.Decode(sys.audioContext, r)
}

func (sys *System) PlayWAV(id ResourceID) {
	resource, ok := sys.wavResources[id]
	if !ok {
		wavInfo, ok := sys.Registry.mapping[id]
		if !ok {
			panic(fmt.Sprintf("unregistered WAV with id=%d", id))
		}
		stream := sys.cache.GetWAV(wavInfo.Path)
		player, err := audio.NewPlayer(sys.audioContext, stream)
		if err != nil {
			panic(err.Error())
		}
		volume := (wavInfo.Volume / 2) + 0.5
		resource = &audioResource{
			player: player,
			volume: volume,
		}
		sys.wavResources[id] = resource
	}
	resource.player.SetVolume(resource.volume)
	resource.player.Rewind()
	resource.player.Play()
}
