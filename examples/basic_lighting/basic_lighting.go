package main

import (
	_ "image/png"
	"math"
	"runtime"

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
	runtime.LockOSThread()
	window := sgl.InitGlfwAndOpenGL(width, height, title)
	defer glfw.Terminate()

	vp := sgl.NewViewpoint(width, height)
	ls := sgl.NewLightSrc()

	cube := sgl.BasicLightObject{}
	cube.Program = sgl.MakeProgram(
		sgl.NewBasicLightVShader(),
		sgl.NewBasicLightFShader(1, 0.3, 0.3),
	)
	cube.SetUniforms(&vp, &ls)
	cube.SetVertices(sgl.NewCube(20))

	angle := 0.0
	previousTime := glfw.GetTime()
	rotateY := mgl32.Rotate3DY(-math.Pi / 6).Mat4()

	sgl.BeforeMainLoop()
	for !window.ShouldClose() {
		sgl.BeforeDrawing()

		// make the cube rotate
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		cube.Model = rotateY.Mul4(
			mgl32.HomogRotate3D(float32(angle)/5, mgl32.Vec3{1, 0, 0}),
		)

		// Render
		cube.Render(&vp, &ls)

		sgl.AfterDrawing(window)
	}
}
