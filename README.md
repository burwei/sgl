# SimpleGL
SimpleGL is a simple Go wrapper for modern OpenGL.   
It's a pure Go repo and is fully compatible with go-gl ecosystem.   
SimpleGL uses the packages below:
 - [go-gl](https://github.com/go-gl/gl)
 - [glfw](https://github.com/go-gl/glfw)
 - [mgl32](https://github.com/go-gl/mathgl)

SimpleGL provides Object, Group, ViewPoint, LightSource, some common shapes and some routine functions to make modern OpenGL development more easily, and fast.  
It could be seen as a lightweight wrapper just to simplify the OpenGL routines and organize the code, so developers can get rid of those verbose routines and focus on shaders, vertices and business logics.  

## Installation
```
go get github.com/burwei/simplegl
```

## Quick Start
Let's get start with the hello cube program. It shows a rotating cube.  
This program is the modified version of [go-gl/example/gl41core-cube](https://github.com/go-gl/example/tree/master/gl41core-cube) example.  
```
package main

import (
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

	cube := sgl.BasicObj{}
	cube.SetProgramVar(sgl.BasicObjProgVar{
		Red:   1,
		Green: 0.3,
		Blue:  0.3,
		Vp:    &vp,
		Ls:    &ls,
	})
	cube.PrepareProgram()
	cube.SetVertices(sgl.NewCube(200))

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
```
result:  
<img src="https://imgur.com/adbC9dE.gif" width="60%">


## Examples
For more examples, see the example folder.
