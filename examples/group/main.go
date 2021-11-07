package main

import (
	_ "image/png"
	"math"

	sgl "github.com/burwei/simplegl"
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

	cube1 := sgl.BasicObj{}
	cube1.SetProgramVar(sgl.BasicObjProgVar{
		Red:   1,
		Green: 0.3,
		Blue:  0.3,
		Vp:    &vp,
		Ls:    &ls,
	})
	cube1.PrepareProgram()
	cube1.SetVertices(sgl.NewCube(200))

	cube2 := sgl.BasicTexObj{}
	cube2.SetProgramVar(sgl.BasicTexObjProgVar{
		TextureSrc: "wood.png",
		Vp:         &vp,
	})
	cube2.PrepareProgram()
	cube2.SetVertices(sgl.NewUniTexCube(100))

	group := sgl.NewGroup()
	group.AddObject("cube1", &cube1)
	group.AddObject("cube2", &cube2)

	angle := 0.0
	tr := 0.0
	dir := 1.0
	rotateY := mgl32.Rotate3DY(-math.Pi / 6).Mat4()
	previousTime := glfw.GetTime()

	sgl.BeforeMainLoop(window, &vp)
	for !window.ShouldClose() {
		sgl.BeforeDrawing()

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
		// 	mgl32.Rotate3DY(float32(angle) / 5).Mat4().Mul4(
		// 		mgl32.Translate3D(0, float32(tr), 0),
		// 	),
		// )

		// Render
		group.Render()

		sgl.AfterDrawing(window)
	}
}
