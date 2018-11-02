package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Dialogue defines connversation flow
type Dialogue struct {
	Name       string `json:"name"`
	Avatar     rl.Texture2D
	AvatarFile string    `json:"avatar"`
	Text       string    `json:"text"`
	Choices    []*Choice `json:"choices"`
	Target     string    `json:"target"`
	Next       *Dialogue `json:"next"`
}

// Choice is a selection from dialogue branches
type Choice struct {
	Text string    `json:"text"`
	Next *Dialogue `json:"next"`
}

func initText(t *Dialogue) {
	if t.AvatarFile != "" {
		t.Avatar = GetTexture(fmt.Sprintf("assets/gfx/%s", t.AvatarFile))
	}

	if t.Next != nil {
		initText(t.Next)
	}

	if t.Choices != nil {
		for _, ch := range t.Choices {
			if ch.Next != nil {
				initText(ch.Next)
			}
		}
	}
}