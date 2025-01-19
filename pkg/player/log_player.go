package player

import "log"

type LogPlayer struct{}

func NewLogPlayer() *LogPlayer {
	return &LogPlayer{}
}

func (p *LogPlayer) Play(s string) {
	log.Printf("playing: %s", s)
}

func (p *LogPlayer) SetSpeed(_ int) {}
func (p *LogPlayer) SetPitch(_ int) {}
