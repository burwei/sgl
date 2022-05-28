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

	// free stl source: https://cults3d.com/en/3d-model/game/iron-man-bust_by-max7th-kimjh
	stlVertices := sgl.ReadBinaryStlFile("ironman_bust_max7th_bin.stl")
	stl := sgl.BasicObj{}
	stl.SetProgramVar(sgl.BasicObjProgVar{
		Red:   1,
		Green: 0.3,
		Blue:  0.3,
		Vp:    &vp,
		Ls:    &ls,
		Mt:    &mt,
	})
	stl.PrepareProgram(true)
	stl.SetVertices(&stlVertices)

	angle := 0.0
	previousTime := glfw.GetTime()
	rotateX := mgl32.Rotate3DX(-math.Pi / 2).Mat4()

	sgl.BeforeMainLoop(window, &vp)
	for !window.ShouldClose() {
		sgl.BeforeDrawing()

		// make the cube rotate
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		stl.SetModel(rotateX.Mul4(
			mgl32.Rotate3DZ(float32(angle) / 3).Mat4(),
		))

		// Render
		stl.Render()

		sgl.AfterDrawing(window)
	}
}
