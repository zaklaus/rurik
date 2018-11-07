package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	ry "github.com/gen2brain/raylib-go/raymath"
	"github.com/solarlune/resolv/resolv"
)

type player struct {
	Locked bool
}

// NewPlayer player
func (p *Object) NewPlayer() {
	aseData := GetAnimData("player")
	p.Ase = &aseData
	p.Texture = GetTexture("player.png")
	p.Size = []int32{p.Ase.FrameWidth, p.Ase.FrameHeight}
	p.Update = updatePlayer
	p.Draw = drawPlayer
	p.GetAABB = getSpriteAABB
	p.HandleCollision = handlePlayerCollision
	p.Facing = rl.NewVector2(1, 0)
	p.IsCollidable = true

	LocalPlayer = p

	playAnim(p, "StandE")
}

func updatePlayer(p *Object, dt float32) {
	p.Ase.Update(dt)

	var moveSpeed float32 = 120

	p.Movement.X = 0
	p.Movement.Y = 0

	if !p.Locked {
		p.Movement.X = GetAxis("horizontal")
		p.Movement.Y = GetAxis("vertical")
	}

	var tag string

	if ry.Vector2Length(p.Movement) > 0 {
		//ry.Vector2Normalize(&p.Movement)
		ry.Vector2Scale(&p.Movement, moveSpeed)

		p.Facing.X = p.Movement.X
		p.Facing.Y = p.Movement.Y
		ry.Vector2Normalize(&p.Facing)

		tag = "Walk"
	} else {
		tag = "Stand"
	}

	if p.Facing.Y > 0 {
		tag += "N"
	} else if p.Facing.Y < 0 {
		tag += "S"
	}

	if p.Facing.X > 0 {
		tag += "E"
	} else if p.Facing.X < 0 {
		tag += "W"
	}

	playAnim(p, tag)

	p.Movement.X *= dt
	p.Movement.Y *= dt

	resX, okX := CheckForCollision(p, int32(p.Movement.X), 0)
	resY, okY := CheckForCollision(p, 0, int32(p.Movement.Y))

	if okX {
		p.Movement.X = float32(resX.ResolveX)
	}

	if okY {
		p.Movement.Y = float32(resY.ResolveY)
	}

	p.Position.X += p.Movement.X
	p.Position.Y += p.Movement.Y
}

func drawPlayer(p *Object) {
	source := getSpriteRectangle(p)
	dest := getSpriteOrigin(p)

	if DebugMode {
		c := getSpriteAABB(p)
		rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
		drawTextCentered(p.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
	}

	rl.DrawTexturePro(p.Texture, source, dest, rl.Vector2{}, 0, SkyColor)
}

func handlePlayerCollision(res *resolv.Collision, p, other *Object) {
	fmt.Println("Collision has happened!")
}
