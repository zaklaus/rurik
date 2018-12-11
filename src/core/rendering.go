/*
 * @Author: V4 Games
 * @Date: 2018-12-11 02:40:45
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-11 03:06:34
 */

package core

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"madaraszd.net/zaklaus/rurik/src/system"
)

var (
	// WorldTexture represents the render target used by the game world
	WorldTexture system.RenderTarget

	// UITexture represents the render target used by the interface
	UITexture system.RenderTarget

	finalRenderTexture system.RenderTarget
	renderTextureQueue = []renderQueueEntry{}
)

type renderQueueEntry struct {
	Target system.RenderTarget
	FlipY  bool
}

// PushRenderTarget appends the render target to the queue to be processed by the compositor pipeline
func PushRenderTarget(tex system.RenderTarget, flipY bool) {
	renderTextureQueue = append(renderTextureQueue, renderQueueEntry{
		Target: tex,
		FlipY:  flipY,
	})
}

func renderGame() {
	rl.BeginDrawing()
	{ // Render the game world
		rl.BeginTextureMode(*WorldTexture)
		{
			drawProfiler.StartInvocation()
			{
				rl.ClearBackground(rl.Black)

				CurrentGameMode.Draw()
			}
			drawProfiler.StopInvocation()
		}
		rl.EndTextureMode()

		// Render all UI elements
		rl.BeginTextureMode(*UITexture)
		{
			rl.ClearBackground(rl.Blank)
			CurrentGameMode.DrawUI()
			DrawEditor()
		}
		rl.EndTextureMode()

		// Render all post-fx elements
		CurrentGameMode.PostDraw()

		// Blend results into one final texture
		rl.BeginTextureMode(*finalRenderTexture)
		{
			rl.DrawTexture(WorldTexture.Texture, 0, 0, rl.White)

			// process the render queue
			rl.BeginBlendMode(rl.BlendAdditive)
			{
				for _, r := range renderTextureQueue {
					v := r.Target
					height := float32(v.Texture.Height)
					if r.FlipY {
						height *= -1
					}
					rl.DrawTextureRec(
						v.Texture,
						rl.NewRectangle(0, 0, float32(v.Texture.Width), height),
						rl.Vector2{},
						rl.White,
					)
				}
			}
			rl.EndBlendMode()

			rl.BeginBlendMode(rl.BlendAlpha)
			{
				rl.DrawTexture(UITexture.Texture, 0, 0, rl.White)
			}
			rl.EndBlendMode()
		}
		rl.EndTextureMode()
	}
	rl.EndDrawing()

	// output final render texture onto the screen
	rl.DrawTexturePro(
		finalRenderTexture.Texture,
		rl.NewRectangle(0, 0, float32(system.ScreenWidth), float32(system.ScreenHeight)),
		rl.NewRectangle(0, 0, float32(system.WindowWidth), float32(system.WindowHeight)),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

	renderTextureQueue = []renderQueueEntry{}
}