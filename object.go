package sgl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Object interface {
	// BuildProgramFromFile will compile and set the shader program from
	// .vert file and .frag file.
	BuildProgramFromFile(vertPath string, fragPath string)

	// BindProgramVar bind the uniform variables to the shader program.
	// The progVar is the customizes struct that contains all the uniform
	// variables that will be used in the shader program of a certain object.
	BindProgramVar(progVar interface{})

	// BuildVaoFromVertices sets the vao and vbo from the vertices.
	BuildVaoFromVertices(*[]float32)

	// GetProgram gets the program of the object.
	GetProgram() uint32

	// SetProgram sets the program of the object.
	SetProgram(program uint32)

	// GetVertices get the vertices of the object.
	GetVertices() *[]float32

	// UpdateVertices update the value of vertices without changing the
	// address as well as vao and vbo.
	UpdateVertices(vertices *[]float32)

	// GetModel gets the model of the object.
	GetModel() mgl32.Mat4

	// SetModel sets the model of the object.
	SetModel(model mgl32.Mat4)

	// SetModelPos set the x, z, y positions of the object.
	SetModelPos(x float32, y float32, z float32)

	// Render refreshs uniform variables and draw the object.
	// All references of the variables that would change the
	// object's states in the main loop (i.e. uniform variables)
	// should have already been prepared when calling BindProgramVar().
	Render()
}

// BaseObj is an Object that keeps the basic info of a Object.
// These info will make an Object be able to be rendered.
type BaseObj struct {
	// Program is the shader program of the object.
	Program uint32

	// Vao stands for "Vertex Array Object", and it contains
	// one or more "Vertex Buffer Object" which is a memory
	// buffer that contains the data of vertices
	Vao uint32

	// Vertices are points that form the shape of the object.
	Vertices *[]float32

	// Model keeps translation/rotation info of the object.
	Model mgl32.Mat4

	// Uniform is the map of the name of uniform variables in
	// shader program and itself.
	Uniform map[string]int32
}

func (obj *BaseObj) BuildProgramFromFile(vertPath string, fragPath string) {
	obj.Program = MakeProgramFromFile(vertPath, fragPath)
}

func (obj *BaseObj) SetModelPos(x float32, y float32, z float32) {
	obj.Model = mgl32.Translate3D(x, y, z)
}

func (obj *BaseObj) GetProgram() uint32 {
	return obj.Program
}

func (obj *BaseObj) SetProgram(program uint32) {
	obj.Program = program
}

func (obj *BaseObj) GetVertices() *[]float32 {
	return obj.Vertices
}

func (obj *BaseObj) UpdateVertices(vertices *[]float32) {
	if obj.Vertices == nil {
		panic("UpdateVertices() needs the existed vertices, so it must be execute after SetVaoFromVertices()")
	} else {
		*(obj.Vertices) = *vertices
	}
}

func (obj *BaseObj) GetModel() mgl32.Mat4 {
	return obj.Model
}

func (obj *BaseObj) SetModel(model mgl32.Mat4) {
	obj.Model = model
}
