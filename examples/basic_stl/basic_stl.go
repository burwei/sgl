package main

import (
	"fmt"
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

	// free stl source: https://cults3d.com/en/3d-model/game/iron-man-bust_by-max7th-kimjh
	stlVertices := sgl.ReadBinaryStlFile("ironman_bust_max7th_bin.stl", 0, 10000)
	fmt.Println(len(stlVertices))
	stl := sgl.BasicLightObject{}
	stl.PrepareProgram(1, 0.3, 0.3)
	stl.SetUniforms(&vp, &ls)
	stl.SetVertices(&stlVertices)

	angle := 0.0
	previousTime := glfw.GetTime()
	stl.Model = mgl32.Rotate3DX(-math.Pi/2).Mat4()
	// rotateY := mgl32.Rotate3DY(-math.Pi / 6).Mat4()

	sgl.BeforeMainLoop(window, &vp)
	for !window.ShouldClose() {
		sgl.BeforeDrawing()

		// make the cube rotate
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		// stl.Model = rotateY.Mul4(
		// 	mgl32.HomogRotate3D(float32(angle)/5, mgl32.Vec3{1, 0, 0}),
		// )

		// Render
		stl.Render(&vp, &ls)

		sgl.AfterDrawing(window)
	}
}
