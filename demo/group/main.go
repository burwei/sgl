package main

import (
	_ "image/png"
	"math"

	"github.com/burwei/sgl"
	"github.com/burwei/sgl/demo/group/objects"
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

	cube1 := objects.NewSimpleObj()
	cube1.BuildProgramFromFile("./objects/simpleobj.vert", "./objects/simpleobj.frag")
	cube1.BindProgramVar(objects.SimpleObjProgVar{
		Red:   1,
		Green: 0.3,
		Blue:  0.3,
		Vp:    &vp,
		Ls:    &ls,
		Mt:    &mt,
	})
	cube1.BuildVaoFromVertices(sgl.NewCube(200))
	cube1.SetModelPos(0, 0, 0)

	cube2 := objects.NewTexCubeObj()
	cube2.BuildProgramFromFile("./objects/tex_cube_obj.vert", "./objects/tex_cube_obj.frag")
	cube2.BindProgramVar(objects.TexCubeObjProgVar{
		TextureSrc: "wood.png",
		Vp:         &vp,
	})
	cube2.BuildVaoFromVertices(sgl.NewUniTexCube(100))

	group := sgl.NewGroup()
	group.AddObject("cube1", cube1)
	group.AddObject("cube2", cube2)

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
		// 	mgl32.Translate3D(0, float32(tr), 0).Mul4(
		// 		mgl32.Rotate3DY(float32(angle) / 5).Mat4(),
		// 	),
		// )

		// Render
		group.Render()

		sgl.AfterDrawing(window)
	}
}
