package corpus

import "math/rand"

type RingCorpus struct {
	phrases []string
	index   int
}

func NewRingCorpus(phrases ...string) *RingCorpus {
	return &RingCorpus{
		phrases: phrases,
		index:   0,
	}
}

func (c *RingCorpus) NextPhrase() string {
	if len(c.phrases) == 0 {
		return ""
	}

	result := c.phrases[c.index]
	c.index = (c.index + 1) % len(c.phrases)

	return result
}

type RandomCorpus struct {
	phrases []string
}

func NewRandomCorpus(phrases ...string) *RandomCorpus {
	return &RandomCorpus{
		phrases: phrases,
	}
}

func (c *RandomCorpus) NextPhrase() string {
	if len(c.phrases) == 0 {
		return ""
	}

	index := rand.Intn(len(c.phrases))

	return c.phrases[index]
}
