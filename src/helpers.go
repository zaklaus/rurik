package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/solarlune/resolv/resolv"
)

const (
	// FrustumSafeMargin safe margin to be considered safe to render off-screen
	FrustumSafeMargin = 32.0
)

func rayRectangleInt32ToResolv(i rl.RectangleInt32) *resolv.Rectangle {
	return &resolv.Rectangle{
		BasicShape: resolv.BasicShape{
			X:           i.X,
			Y:           i.Y,
			Collideable: true,
		},
		W: i.Width,
		H: i.Height,
	}
}

func drawTextCentered(text string, posX, posY, fontSize int32, color rl.Color) {
	if fontSize < 10 {
		fontSize = 10
	}

	rl.DrawText(text, posX-rl.MeasureText(text, fontSize)/2, posY, fontSize, color)
}

func vector2Lerp(v1, v2 rl.Vector2, amount float32) (result rl.Vector2) {
	result.X = v1.X + amount*(v2.X-v1.X)
	result.Y = v1.Y + amount*(v2.Y-v1.Y)

	return result
}

func scalarLerp(v1, v2 float32, amount float32) (result float32) {
	result = v1 + amount*(v2-v1)

	return result
}

func lerpColor(a, b rl.Vector3, t float64) rl.Vector3 {
	return raymath.Vector3Lerp(a, b, float32(t))
}

func getColorFromHex(hex string) (rl.Vector3, error) {
	if hex == "" {
		return rl.Vector3{}, fmt.Errorf("hex not specified")
	}

	c, err := colorful.Hex("#" + hex[3:])

	if err != nil {
		return rl.Vector3{}, err
	}

	d := rl.NewVector3(
		float32(c.R),
		float32(c.G),
		float32(c.B),
	)

	return d, nil
}

func vec3ToColor(a rl.Vector3) rl.Color {
	return rl.NewColor(
		uint8(a.X*255),
		uint8(a.Y*255),
		uint8(a.Z*255),
		255,
	)
}

func isMouseInRectangle(x, y, x2, y2 int32) bool {
	x2 = x + x2
	y2 = y + y2

	mo := rl.GetMousePosition()
	m := [2]int32{
		int32(mo.X) / ScaleRatio,
		int32(mo.Y) / ScaleRatio,
	}

	if m[0] > x && m[0] < x2 &&
		m[1] > y && m[1] < y2 {
		return true
	}

	return false
}

func getSpriteAABB(o *Object) rl.RectangleInt32 {
	return rl.RectangleInt32{
		X:      int32(o.Position.X) - int32(float32(o.Ase.FrameWidth/2)) + int32(float32(o.Ase.FrameWidth/4)),
		Y:      int32(o.Position.Y),
		Width:  o.Ase.FrameWidth / 2,
		Height: o.Ase.FrameHeight / 2,
	}
}

func playAnim(p *Object, animName string) {
	if p.Ase.GetAnimation(animName) != nil {
		p.Ase.Play(animName)
	} else {
		//log.Println("Animation name:", animName, "not found!")
	}
}

func getSpriteRectangle(o *Object) rl.Rectangle {
	sourceX, sourceY := o.Ase.GetFrameXY()
	return rl.NewRectangle(float32(sourceX), float32(sourceY), float32(o.Ase.FrameWidth), float32(o.Ase.FrameHeight))
}

func getSpriteOrigin(o *Object) rl.Rectangle {
	return rl.NewRectangle(o.Position.X-float32(o.Ase.FrameWidth/2), o.Position.Y-float32(o.Ase.FrameHeight/2), float32(o.Ase.FrameWidth), float32(o.Ase.FrameHeight))
}

func isPointWithinRectangle(p rl.Vector2, r rl.Rectangle) bool {
	if p.X > r.X && p.X < (r.X+r.Width) &&
		p.Y > r.Y && p.Y < (r.Y+r.Height) {
		return true
	}

	return false
}

func isPointWithinFrustum(p rl.Vector2) bool {
	camOffset := rl.Vector2{
		X: float32(int(MainCamera.Position.X*MainCamera.Zoom - screenW/2)),
		Y: float32(int(MainCamera.Position.Y*MainCamera.Zoom - screenH/2)),
	}

	cam := rl.Rectangle{
		X:      camOffset.X - FrustumSafeMargin,
		Y:      camOffset.Y - FrustumSafeMargin,
		Width:  screenW + FrustumSafeMargin*2,
		Height: screenH + FrustumSafeMargin*2,
	}

	return isPointWithinRectangle(p, cam)
}
