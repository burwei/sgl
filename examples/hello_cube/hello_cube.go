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

	cube := sgl.BasicObject{}
	cube.PrepareProgram(1, 0.3, 0.3)
	cube.SetUniforms(&vp)
	cube.SetVertices(sgl.NewCube(200))

	sgl.BeforeMainLoop()
	for !window.ShouldClose() {
		sgl.BeforeDrawing()

		// rotate
		cube.Model = mgl32.Rotate3DY(math.Pi / 6).Mat4()

		// Render
		cube.Render(&vp)

		sgl.AfterDrawing(window)
	}
}
