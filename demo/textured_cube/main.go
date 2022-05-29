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
	window := sgl.Init(width, height, title)
	defer sgl.Terminate()

	vp := sgl.NewViewpoint(width, height)

	cube := objects.NewTexCubeObj()
	cube.BuildProgramFromFile("./objects/tex_cube_obj.vert", "./objects/tex_cube_obj.frag")
	cube.BindProgramVar(objects.TexCubeObjProgVar{
		TextureSrc: "wood.png",
		Vp:         &vp,
	})
	cube.BuildVaoFromVertices(sgl.NewUniTexCube(200))
	cube.SetModelPos(0, 0, 0)

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
