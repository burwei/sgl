package main

import (
	"math"

	"github.com/burwei/sgl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
	title  = "SimpleGL"
)

func main() {
	window := sgl.Init(width, height, title)
	defer sgl.Terminate()

	vp := sgl.NewViewpoint(width, height)
	ls := sgl.NewLightSrc()
	mt := sgl.NewMaterial()

	cube := sgl.NewSimpleObj()
	cube.SetProgVar(sgl.SimpleObjVar{
		Red:   1,
		Green: 0.3,
		Blue:  0.3,
		Vp:    &vp,
		Ls:    &ls,
		Mt:    &mt,
	})
	cube.SetVertices(sgl.NewCube(200))
	cube.SetModel(mgl32.Translate3D(0, 0, 0))

	angle := 0.0
	previousTime := glfw.GetTime()
	rotateY := mgl32.Rotate3DY(-math.Pi / 6).Mat4()

	sgl.BeforeMainLoop(window, &vp)
	for !window.ShouldClose() {
		sgl.BeforeDrawing()

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

		sgl.AfterDrawing(window)
	}
}
