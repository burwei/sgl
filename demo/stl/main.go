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
	window, vp := sgl.Init(width, height, title)
	defer sgl.Terminate()

	// free stl source: https://cults3d.com/en/3d-model/game/iron-man-bust_by-max7th-kimjh
	stlVertices := sgl.ReadBinaryStlFile("ironman_bust_max7th_bin.stl")
	stl := sgl.NewBaseObj()
	stl.SetProgVar(sgl.BaseObjVar{Vp: vp})
	stl.SetVertices(&stlVertices)
	stl.SetModel(mgl32.Translate3D(0, 0, 0))

	angle := 0.0
	previousTime := glfw.GetTime()
	rotateX := mgl32.Rotate3DX(-math.Pi / 2).Mat4()

	sgl.MainLoop(window, func() {
		// make the object rotate
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		stl.SetModel(rotateX.Mul4(
			mgl32.Rotate3DZ(float32(angle) / 3).Mat4(),
		))

		// Render
		stl.Render()
	})
}
