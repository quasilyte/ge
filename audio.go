package ge

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type Audio struct {
	ctx          *Context
	audioContext *audio.Context
	wavPlayers   map[*wav.Stream]*audio.Player
}

func (a *Audio) init(ctx *Context) {
	a.ctx = ctx
	a.audioContext = audio.NewContext(32000)
	a.wavPlayers = make(map[*wav.Stream]*audio.Player)
}

func (a *Audio) PlayWAV(stream *wav.Stream, forced bool) {
	p, ok := a.wavPlayers[stream]
	if !ok {
		var err error
		p, err = audio.NewPlayer(a.audioContext, stream)
		if err != nil {
			a.ctx.OnCriticalError(err)
		}
		a.wavPlayers[stream] = p
	}
	if !forced && p.IsPlaying() {
		return
	}
	p.SetVolume(0.1)
	p.Rewind()
	p.Play()
}
