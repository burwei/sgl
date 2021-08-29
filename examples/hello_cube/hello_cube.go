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

func main() {
	runtime.LockOSThread()
	window := sgl.InitGlfwAndOpenGL(windowWidth, windowHeight)
	defer glfw.Terminate()

	vp := sgl.NewViewpoint(windowWidth, windowHeight)

	cube := sgl.SimpleModel{}
	cube.SetProgram(sgl.NewSimpleVShader(), sgl.NewSimpleFShader(1, 0.3, 0.3))
	cube.SetMatrixes(&vp)
	cube.SetVao(sgl.NewSimpleCube(20))

	sgl.SetBasicGlobalConfigs()

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
