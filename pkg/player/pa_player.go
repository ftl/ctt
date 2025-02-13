package player

import (
	"context"
	"fmt"
	"log"

	"github.com/ftl/digimodes/cw"
	"github.com/ftl/patrix/pa"
)

type PAPlayer struct {
	speed      int
	farnsworth int
	pitch      int
}

func NewPAPlayer(speed, pitch int) *PAPlayer {
	return &PAPlayer{
		speed:      speed,
		farnsworth: 0,
		pitch:      pitch,
	}
}

func (p *PAPlayer) Play(text string) {
	ctx := context.Background()

	modulator := cw.NewModulator(float64(p.pitch), p.speed)
	modulator.SetFarnsworthWPM(p.farnsworth)
	modulator.AbortWhenDone(ctx.Done())

	oscillator, err := pa.NewOscillator()
	if err != nil {
		log.Fatal(err)
	}

	oscillator.Modulator = modulator
	oscillator.Start()

	_, err = fmt.Fprintln(modulator, text)
	if err != nil {
		log.Fatal(err)
	}

	oscillator.Stop(ctx)
	modulator.Close()
	oscillator.Close()
}

func (p *PAPlayer) SetSpeed(speed int) {
	p.speed = speed
}

func (p *PAPlayer) SetFarnsworth(fwpm int) {
	p.farnsworth = fwpm
}

func (p *PAPlayer) SetPitch(pitch int) {
	p.pitch = pitch
}
