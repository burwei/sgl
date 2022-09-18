package main

import (
	_ "image/png"
	"math"

	"github.com/burwei/sgl"
	"github.com/burwei/sgl/demo/textured_cube/objects"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
	title  = "SimpleGL"
)

func main() {
	window, vp := sgl.Init(width, height, title)
	defer sgl.Terminate()

	cube := objects.NewTexCubeObj()
	cube.SetProgVar(objects.TexCubeObjVar{
		TextureSrc: "wood.png",
		Vp:         vp,
	})
	cube.SetVertices(sgl.NewUniTexCube(200))
	cube.SetModel(mgl32.Translate3D(0, 0, 0))

	angle := 0.0
	previousTime := glfw.GetTime()
	rotateY := mgl32.Rotate3DY(-math.Pi / 6).Mat4()

	sgl.MainLoop(window, func() {
		// make the cube rotate
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		cube.SetModel(rotateY.Mul4(
			mgl32.Rotate3DX(float32(angle) / 5).Mat4(),
		))

		// Render
		cube.Render()
	})
}
