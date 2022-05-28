# SimpleGL
SimpleGL is a simple Go wrapper for modern OpenGL.   
It's a pure Go repo and is fully compatible with go-gl ecosystem.  

<img src="https://imgur.com/JX65X3U.gif" width="100%">


SimpleGL uses the packages below:  
 - [go-gl](https://github.com/go-gl/gl)
 - [glfw](https://github.com/go-gl/glfw)
 - [mgl32](https://github.com/go-gl/mathgl)

SimpleGL provides Object, Group, Viewpoint, LightSource, some common shapes and some routine functions to make modern OpenGL development more easily and fast.  
It could be seen as a lightweight wrapper just to simplify the OpenGL routines and organize the code, so developers can get rid of those verbose routines and focus on shaders, vertices and business logics.  

## Installation
```
go get github.com/burwei/sgl
```

## Quick Start
Let's get start with the hello cube program. It shows a rotating cube.  
This program is the modified version of [go-gl/example/gl41core-cube](https://github.com/go-gl/example/tree/master/gl41core-cube) example.  
```
package main

import (
	"math"

	"github.com/burwei/sgl"
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

The good news is, developers with Go programming skills and zero OpenGL knowledge might still able to learn some basic OpenGL directly from SimpleGL, since it organizes the code to make it easier to use and understand.

The usage introduction contains the contents below:  
 - OpenGL program struture
 - Object
 - Shape
 - Viewpoint & Coordinate system
 - LightSource & Material
 - Group
 - STL

### OpenGL Program structure
Modern OpenGL program can be roughly divided into two parts, CPU program and GPU programs.  

The CPU program contains two parts, setup and main loop.  
In setup part we call ```sgl.Init()```, which will lock the current thread and init OpenGL and GLFW. GLFW is the library that handles the graphics output and device input, such as window, keyboard, mouse, joystick and so on. Variable assignment, input callback settings and all things we should prepared before starting the main loop will be in the setup part.  
The main loop is the ```for !window.ShouldClose() {}``` loop. We render the objects in main loop. Before and after the rendering, we call ```sgl.BeforeDrawing()``` and ```sgl.AfterDrawing()``` to clean, swap buffers and poll events.  

The GPU programs contains also two parts, Program Object and shaders. Program Object is used in render operation and it's also the "Program" that SimpleGL refers to when calling APIs like ```sgl.Object.PrepareProgram()``` and ```slg.ObjectSetProgramVar()``` and so on. SimpleGL sees each Program Object as a final all-in-one program for each sgl.Object, so all varialbes of shaders attached to the Program Object are also seen as the variables of the "Program". Shaders are written in GLSL, and are used to determine how to draw the vertices.   
One Program Object can combine multiple shaders to do the rendering job, but we only attach a vertex shader and a fragment shader on it in SimpleGL (so far). Vertex shader calculates the positions of vertices and fragment shader calculates the colors of fragments.   

The above is just a simplified introduction. To know more about how OpenGL works, see [OpenGL rendering pipeline overview](https://www.khronos.org/opengl/wiki/Rendering_Pipeline_Overview).  

### Object
sgl.Object is an interface that represents a object with a specific Project Object that can be render on the window after it gets the program variables and vertex array it needs.  

sgl.Object + program variables + vertex array = visuable object  

sgl.BasicObj is the object with basic lighting, and it's able to draw any shape (any vertex array that contains 3*n float32 values). Developers can implement their own sgl.Object to create some cool objects. By implement sgl.Object, the object could be more easiy to use and be able to move together as a group.

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
Shapes are described by vertex arrays, which are 1-D float32 arrays. The most basic vertex arrays are those who use 3 float32 values to represent a vertex's X,Y,Z position. Sometimes vertex array will contains some meta data such as the direction of the texture.  

For instance, sgl.NewCube() is a vertex array that use 3 float32 to represent a vertex and form a cube with a user-defined side length.

```
cube.SetVertices(sgl.NewCube(200))
```


### Viewpoint & Coordinate system
sgl.Viewpoint provides a default camera (eye) position on (X, Y, Z) = (0, 0, 1000) and default target position on (X, Y, Z) = (0, 0, 0). The default top direction of the camera is positive Y and the default projection is perspective projection.   

Before we read the code, we should understand the position of the camera as well as the coordinate systems.  

<img src="https://imgur.com/9XwCWA1.png" width="80%">

There are four coordinate systems here:  
 1. local coordinate
 2. world-space coordinate
 3. view-space coordinate
 4. clip-space coordinate

 <img src="https://imgur.com/fw0Uao4.png" width="80%">

The usage of sgl.Viewpoint is simple. Just new one with the width and height of the window. Although there's no strong restrictions, all sgl.Object should contain a sgl.Viewpoint to make the object appear in view-space and clip-space coordinate correctly.  

```
vp := sgl.NewViewpoint(width, height)

cube := sgl.BasicObj{}
cube.SetProgramVar(sgl.BasicObjProgVar{
	Red:   1,
	Green: 0.3,
	Blue:  0.3,
	Vp:    &vp,
	Ls:    &ls,
	Mt:    &mt,
})
```


### LightSource & Material
sgl.LightSource and agl.Material provides a default light source and default material. These two are essential for those sgl.Object that render the lighting effect.  

sgl.LightSource contains 3 attributes: light position, light color and light intensity. All of them are easy to understand.  

sgl.Material contains 4 attributes: ambient, diffuse, specular and shininess. Ambient determines what color does the material reflects under ambient lighting; diffuse determines what color does the material reflects under diffuse lighting; specular determines the color of the material's specular highligh; and shininess determines the scattering/radius of the specular highlight.

```
ls := sgl.NewLightSrc()
mt := sgl.Material{
	Ambient: mgl32.Vec3{0.1, 0.1, 0.1},
	Diffuse: mgl32.Vec3{0.6, 0.6, 0.6},
	Specular: mgl32.Vec3{1.5, 1.5, 1.5},
	Shininess: 24,
}

newCube := sgl.BasicObj{}
newCube.SetProgramVar(sgl.BasicObjProgVar{
	Red:   1,
	Green: 0.3,
	Blue:  0.3,
	Vp:    &vp,
	Ls:    &ls,
	Mt:    &mt,
})
```

### Group
sgl.Group collects mutiple sgl.Object and make them move together like a bigger object. Besides making sgl.Object move together, sgl.Group can also move any collected sgl.Object individually.  

```
// before main loop
group := sgl.NewGroup()
group.AddObject("cube1", &cube1)

// in main loop
group.SetObjectModel("cube1", rotateY.Mul4(
	mgl32.Rotate3DX(float32(angle)/5).Mat4(),
))
group.SetGroupModel(
	mgl32.Translate3D(0, float32(tr), 0).Mul4(
		mgl32.Rotate3DY(float32(angle)/5).Mat4(),
	),
)
group.Render()
```

### STL
STL is a common file format for 3D models.  
SimpleGL also provides some APIs to read STL files and turn them into vertex arrays.  

The below shows how to read a STL file and create a object with it.   
The STL file is download from a free STL platform called cults3d.com, and the link is [here](https://cults3d.com/en/3d-model/game/iron-man-bust_by-max7th-kimjh).  
```
// free stl source: https://cults3d.com/en/3d-model/game/iron-man-bust_by-max7th-kimjh
// read binary STL file and shift it to the center
stlVertices := sgl.ReadBinaryStlFile("ironman_bust_max7th_bin.stl")
stl := sgl.BasicObj{}
stl.SetProgramVar(sgl.BasicObjProgVar{
	Red:   1,
	Green: 0.3,
	Blue:  0.3,
	Vp:    &vp,
	Ls:    &ls,
	Mt:    &mt,
})
stl.PrepareProgram(true)
stl.SetVertices(&stlVertices)

// read binary STL file without shifting
stlVertices := sgl.ReadBinaryStlFileRaw("ironman_bust_max7th_bin.stl")

// read binary STL file and shift it to a specific center
stlVertices := sgl.ReadBinaryStlFileWithCenter("ironman_bust_max7th_bin.stl", 50, 100, 20)
```
result:  
<img src="https://imgur.com/M2sSHD8.gif" width="60%">

## Examples
For more examples, see the example folder.
