package resource

import (
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type Audio struct {
	Path string

	// Volume adjust how loud this sound will be.
	// The default value of 0 means "unadjusted".
	// Value greated than 0 increases the volume, negative values decrease it.
	// This setting accepts values in [-1, 1] range, where -1 mutes the sound
	// while 1 makes it as loud as possible.
	Volume float64
}

type AudioID int

type AudioRegistry struct {
	mapping map[AudioID]Audio
}

func (r *AudioRegistry) Set(id AudioID, info Audio) {
	r.mapping[id] = info
}

type AudioSystem struct {
	loader *Loader

	currentMusic *audioResource

	audioContext *audio.Context
	resources    map[AudioID]*audioResource

	currentQueueSound *audioResource
	soundQueue        []AudioID
}

type audioResource struct {
	player *audio.Player
	volume float64
}

func (sys *AudioSystem) Init(l *Loader) {
	sys.loader = l
	sys.audioContext = audio.NewContext(44100)
	sys.resources = make(map[AudioID]*audioResource)
	sys.soundQueue = make([]AudioID, 0, 4)

	// Audio player factory has lazy initialization that may lead
	// to a ~0.2s delay before the first sound can be played.
	// To avoid that delay, we force that factory to initialize
	// right now, before the game is started.
	dummy := sys.audioContext.NewPlayerFromBytes(nil)
	dummy.Rewind()
}

func (sys *AudioSystem) Update() {
	if sys.currentQueueSound == nil {
		if len(sys.soundQueue) == 0 {
			// Nothing to play in the queue.
			return
		}
		// Do a dequeue.
		sys.currentQueueSound = sys.playSound(sys.soundQueue[0])
		for i, id := range sys.soundQueue[1:] {
			sys.soundQueue[i] = id
		}
		sys.soundQueue = sys.soundQueue[:len(sys.soundQueue)-1]
		return
	}
	if !sys.currentQueueSound.player.IsPlaying() {
		// Finished playing the current enqueued sound.
		sys.currentQueueSound = nil
	}
}

func (sys *AudioSystem) DecodeWAV(r io.Reader) (*wav.Stream, error) {
	return wav.Decode(sys.audioContext, r)
}

func (sys *AudioSystem) DecodeOGG(r io.Reader) (*vorbis.Stream, error) {
	return vorbis.Decode(sys.audioContext, r)
}

func (sys *AudioSystem) getOGGResource(id AudioID) *audioResource {
	resource, ok := sys.resources[id]
	if ok {
		return resource
	}
	stream := sys.loader.LoadOGG(id)
	oggInfo := sys.loader.GetAudioInfo(id)
	loopedStream := audio.NewInfiniteLoop(stream, stream.Length())
	player, err := sys.audioContext.NewPlayer(loopedStream)
	if err != nil {
		panic(err.Error())
	}
	volume := (oggInfo.Volume / 2) + 0.5
	resource = &audioResource{
		player: player,
		volume: volume,
	}
	sys.resources[id] = resource
	return resource
}

func (sys *AudioSystem) PauseCurrentMusic() {
	if sys.currentMusic == nil {
		return
	}
	sys.currentMusic.player.Pause()
}

func (sys *AudioSystem) ContinueCurrentMusic() {
	if sys.currentMusic == nil || sys.currentMusic.player.IsPlaying() {
		return
	}
	sys.currentMusic.player.SetVolume(sys.currentMusic.volume)
	sys.currentMusic.player.Play()
}

func (sys *AudioSystem) ContinueMusic(id AudioID) {
	sys.continueMusic(sys.getOGGResource(id))
}

func (sys *AudioSystem) continueMusic(res *audioResource) {
	if res.player.IsPlaying() {
		return
	}
	sys.currentMusic = res
	res.player.SetVolume(res.volume)
	res.player.Play()
}

func (sys *AudioSystem) PlayMusic(id AudioID) {
	res := sys.getOGGResource(id)
	if sys.currentMusic != nil && res.player == sys.currentMusic.player && res.player.IsPlaying() {
		return
	}
	sys.currentMusic = res
	res.player.SetVolume(res.volume)
	res.player.Rewind()
	res.player.Play()
}

func (sys *AudioSystem) ResetQueue() {
	sys.soundQueue = sys.soundQueue[:0]
}

func (sys *AudioSystem) EnqueueSound(id AudioID) {
	sys.soundQueue = append(sys.soundQueue, id)
}

func (sys *AudioSystem) PlaySound(id AudioID) {
	sys.playSound(id)
}

func (sys *AudioSystem) playSound(id AudioID) *audioResource {
	resource, ok := sys.resources[id]
	if !ok {
		stream := sys.loader.LoadWAV(id)
		wavInfo := sys.loader.GetAudioInfo(id)
		player, err := sys.audioContext.NewPlayer(stream)
		if err != nil {
			panic(err.Error())
		}
		volume := (wavInfo.Volume / 2) + 0.5
		resource = &audioResource{
			player: player,
			volume: volume,
		}
		sys.resources[id] = resource
	}
	resource.player.SetVolume(resource.volume)
	resource.player.Rewind()
	resource.player.Play()
	return resource
}
