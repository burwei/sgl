package main

import (
	_ "image/png"
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

	ls := sgl.NewLightSrc()
	mt := sgl.NewMaterial()

	cube1 := sgl.NewSimpleObj()
	cube1.SetProgVar(sgl.SimpleObjVar{
		Red:   1,
		Green: 0.3,
		Blue:  0.3,
		Vp:    vp,
		Ls:    ls,
		Mt:    mt,
	})
	cube1.SetVertices(sgl.NewCube(200))
	cube1.SetModel(mgl32.Translate3D(0, 0, 0))

	cube2 := &sgl.SimpleObj{}
	cube2.SetProgram(cube1.GetProgram())
	cube2.SetProgVar(sgl.SimpleObjVar{
		Red:   0.3,
		Green: 1.0,
		Blue:  0.3,
		Vp:    vp,
		Ls:    ls,
		Mt:    mt,
	})
	cube2.SetVertices(sgl.NewCube(100))
	cube2.SetModel(mgl32.Translate3D(0, 0, 0))

	group := sgl.NewGroup()
	group.AddObject("cube1", cube1)
	group.AddObject("cube2", cube2)

	angle := 0.0
	tr := 0.0
	dir := 1.0
	rotateY := mgl32.Rotate3DY(-math.Pi / 6).Mat4()
	previousTime := glfw.GetTime()

	sgl.MainLoop(window, func() {
		// calc states
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		if math.Abs(tr) >= 250 {
			dir *= -1
		}
		tr += elapsed * 50 * dir

		// make cube1 rotate around X axis
		group.SetObjectModel("cube1", rotateY.Mul4(
			mgl32.Rotate3DX(float32(angle)/5).Mat4(),
		))
		// make cube2 translate on X axis
		group.SetObjectModel("cube2",
			rotateY.Mul4(
				mgl32.Translate3D(float32(tr)*2, 0, 0),
			),
		)

		// Group movement 1: translate only
		group.SetGroupModel(
			mgl32.Translate3D(0, float32(tr), 0),
		)

		// Group movement 2: both rotate and translate
		// group.SetGroupModel(
		// 	mgl32.Translate3D(0, float32(tr), 0).Mul4(
		// 		mgl32.Rotate3DY(float32(angle) / 5).Mat4(),
		// 	),
		// )

		// Render
		group.Render()
	})
}
