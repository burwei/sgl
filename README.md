# SimpleGL
SimpleGL is a simple Go wrapper for modern OpenGL.   
It's a pure Go repo and is fully compatible with go-gl ecosystem.  

<img src="https://imgur.com/0WeWWPV.gif" width="100%">


SimpleGL uses the packages below:  
 - [go-gl](https://github.com/go-gl/gl)
 - [glfw](https://github.com/go-gl/glfw)
 - [mgl32](https://github.com/go-gl/mathgl)

SimpleGL provides Object, Group, Viewpoint, LightSource, some common shapes and some routine functions to make modern OpenGL development more easily, and fast.  
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
	mt := sgl.NewMaterial()

	cube := sgl.BasicObj{}
	cube.SetProgramVar(sgl.BasicObjProgVar{
		Red:   1,
		Green: 0.3,
		Blue:  0.3,
		Vp:    &vp,
		Ls:    &ls,
		Mt:    &mt,
	})
	cube.PrepareProgram(true)
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

## Usage
To use SimpleGL, developers should know how to develop modern OpenGL.  
[LeanOpenGL.com](https://learnopengl.com/Getting-started/OpenGL) is a good place to get started if one is not so familiar with OpenGL.  

However, developers with Go programming skills might still able to learn some basic OpenGL directly from SimpleGL, since it organizes the code to make it easier to use and understand.

The usage introduction contains the contents below:  
 - OpenGL program struture
 - Object
 - Shape
 - Viewpoint & Coordinate system
 - LightSource & Material
 - Group

### OpenGL Program structure
Modern OpenGL program can be roughly divided into two parts, CPU program and GPU program (shader program).  

The CPU program contains two parts, setup and main loop.  
In setup part we call ```sgl.Init()```, which will lock the current thread and init OpenGL and GLFW. GLFW is the library that handles the graph output and device input, such as window, keyboard, mouse, joystick and so on.  
Variable assignment, input callback settings and all things we should prepared before starting the main loop will be in the setup part.  
The main loop is ```for !window.ShouldClose() {}``` loop. In main loop part we render the objects. Before and after the rendering, we call ```sgl.BeforeDrawing()``` and ```sgl.AfterDrawing()``` to clean, swap buffers and poll events.  

The GPU program (shader program) is written in GLSL. One program object can contain multiple shaders, but usually we use the program object that contains one vertex shader and one fragment shader. Vertex shader calculates the positions of vertices and fragment shader calculates the colors of fragments. The GPU program is prepared in ```Object.PrepareProgram()```, and the variables of the GPU program are updated when calling ```Object.Render()```.  

To know more about how OpenGL works, see [OpenGL rendering pipeline overview](https://www.khronos.org/opengl/wiki/Rendering_Pipeline_Overview).  

### Object
sgl.Object is an interface that represents a object with specific vertex shader and fragment shader that can be render on the window after it gets the program variables and vertex array it needs.  

sgl.Object + program variables + vertex array = visuable object  

sgl.BasicObj is the object with basic lighting, and it's able to draw any shape (any 3*n vertex array). Developers can implement their own sgl.Object to create some cool objects. By implement sgl.Object, the object could be more easiy to use and be able to make a group to move together.

```
// create a sgl.Object
cube := sgl.BasicObj{}

// set shader program variables
cube.SetProgramVar(sgl.BasicObjProgVar{
	Red:   1,
	Green: 0.3,
	Blue:  0.3,
	Vp:    &vp,
	Ls:    &ls,
	Mt:    &mt,
})

// produce the shader program and bind the program variables with it
cube.PrepareProgram(true)

// set the vetex array
cube.SetVertices(sgl.NewCube(200))

// render the object (in main loop)
cube.Render()
```

### Shape
Shape is determined by vertex array. The most basic vertex array contains 3*n vertices, where 3 is X,Y,Z position in order and n is the number of the vertices.  


## Examples
For more examples, see the example folder.
