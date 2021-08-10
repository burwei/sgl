package main

import (
	_ "image/png"
	"runtime"

	sgl "github.com/burwei/simplegl"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom (-Y)
	-10.0, -10.0, -10.0, 0.0, 0.0,
	10.0, -10.0, -10.0, 1.0, 0.0,
	-10.0, -10.0, 10.0, 0.0, 1.0,
	10.0, -10.0, -10.0, 1.0, 0.0,
	10.0, -10.0, 10.0, 1.0, 1.0,
	-10.0, -10.0, 10.0, 0.0, 1.0,

	// Top (+Y)
	-10.0, 10.0, -10.0, 0.0, 0.0,
	-10.0, 10.0, 10.0, 0.0, 1.0,
	10.0, 10.0, -10.0, 1.0, 0.0,
	10.0, 10.0, -10.0, 1.0, 0.0,
	-10.0, 10.0, 10.0, 0.0, 1.0,
	10.0, 10.0, 10.0, 1.0, 1.0,

	// Front (+Z)
	-10.0, -10.0, 10.0, 1.0, 0.0,
	10.0, -10.0, 10.0, 0.0, 0.0,
	-10.0, 10.0, 10.0, 1.0, 1.0,
	10.0, -10.0, 10.0, 0.0, 0.0,
	10.0, 10.0, 10.0, 0.0, 1.0,
	-10.0, 10.0, 10.0, 1.0, 1.0,

	// Back (-Z)
	-10.0, -10.0, -10.0, 0.0, 0.0,
	-10.0, 10.0, -10.0, 0.0, 1.0,
	10.0, -10.0, -10.0, 1.0, 0.0,
	10.0, -10.0, -10.0, 1.0, 0.0,
	-10.0, 10.0, -10.0, 0.0, 1.0,
	10.0, 10.0, -10.0, 1.0, 1.0,

	// Left (-X)
	-10.0, -10.0, 10.0, 0.0, 1.0,
	-10.0, 10.0, -10.0, 1.0, 0.0,
	-10.0, -10.0, -10.0, 0.0, 0.0,
	-10.0, -10.0, 10.0, 0.0, 1.0,
	-10.0, 10.0, 10.0, 1.0, 1.0,
	-10.0, 10.0, -10.0, 1.0, 0.0,

	// Right (+X)
	10.0, -10.0, 10.0, 1.0, 1.0,
	10.0, -10.0, -10.0, 1.0, 0.0,
	10.0, 10.0, -10.0, 0.0, 0.0,
	10.0, -10.0, 10.0, 1.0, 1.0,
	10.0, 10.0, -10.0, 0.0, 0.0,
	10.0, 10.0, 10.0, 0.0, 1.0,
}

func main() {
	runtime.LockOSThread()
	window := sgl.InitGlfwAndOpenGL(windowWidth, windowHeight)
	defer glfw.Terminate()

	vp := sgl.NewViewPoint(windowWidth, windowHeight)

	cube := sgl.SimpleModel{}
	cube.SetProgram(sgl.TexVShader, sgl.TexFShader)
	cube.SetMatrixes(&vp)
	cube.SetTexture("wood.png") //source: https://unsplash.com/photos/mI-QcAP95Ok
	cube.SetVao(&cubeVertices)

	sgl.InitGlobalSettings()

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		// Clear before redraw
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// --- Drawing starts ---
		// update variables
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		cube.Model = mgl32.HomogRotate3D(float32(angle)/5, mgl32.Vec3{1, 0, 0})

		// Render
		cube.Render(&vp)
		// --- Drawing ends ---

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
