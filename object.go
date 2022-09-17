package sgl

import (
	"fmt"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object interface {
	// GetProgram gets the program of the object.
	GetProgram() uint32

	// SetProgram sets the program of the object using an existing program.
	// It compiles the sharder program without binding the program variables.
	SetProgram(program uint32)

	// GetProgram gets the program variables of the object.
	// The return value should be a "progVar".
	// "progVar" is the customizes struct that contains all the uniform
	// variables that will be used in the shader program of a certain object.
	GetProgVar() interface{}

	// SetProgram sets the program variables of the object.
	// The binding of uniform variables to the program should be set here.
	SetProgVar(progVar interface{})

	// GetVertices gets the vertices of the object.
	GetVertices() *[]float32

	// SetVertices sets the vertices of the object.
	// VAO, VBO and EBO should be set here if needed.
	SetVertices(vertices *[]float32)

	// GetModel gets the model of the object.
	GetModel() mgl32.Mat4

	// SetModel sets the model of the object.
	SetModel(model mgl32.Mat4)

	// Render refreshes uniform variables and draw the object.
	// All references of the variables that would change the
	// object's states in the main loop (i.e. uniform variables)
	// should have already been prepared when calling BindProgramVar().
	Render()
}

// BaseObjVar is the customizes struct that contains all the uniform
// variables
type BaseObjVar struct {
	Vp *Viewpoint
	Ls *LightSrc
	Mt *Material
}

// BaseObj is an Object that keeps the basic info of a Object.
// It'll render a wireframe object.
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

	// ProgVar is the customizes struct that contains all the uniform
	// variables that will be used in the shader program of a certain object.
	// Developer should implement their own ProgVar and put it here to shadow
	// this BaseObjVar type of ProgVar
	ProgVar BaseObjVar
}

// NewBaseObj return a default BaseObj
func NewBaseObj() Object {
	// The recommended way to make shader program is importing
	// the .vert and .frag files. In this way the editor/IDE can
	// check the syntax for us.
	// Like this:
	// 		obj.Program = sgl.MakeProgramFromFile(verPath, fragPath)

	// But if you want the shader codes to be more portable without
	// considering the file paths related to the executed file,
	// put shader codes into a string also is a way.
	obj := &BaseObj{}
	obj.SetProgram(MakeProgram(getBaseObjVS(), getBaseObjFS()))

	return obj
}

func (obj *BaseObj) GetProgram() uint32 {
	return obj.Program
}

func (obj *BaseObj) SetProgram(program uint32) {
	obj.Program = program
}

func (obj *BaseObj) GetProgVar() interface{} {
	return obj.ProgVar
}

func (obj *BaseObj) SetProgVar(progVar interface{}) {
	if pv, ok := progVar.(BaseObjVar); ok {
		obj.ProgVar = pv
	} else {
		panic("progVar is not a BaseObjVar")
	}

	obj.Uniform = map[string]int32{}

	obj.Uniform["project"] = gl.GetUniformLocation(obj.Program, gl.Str("projection\x00"))
	obj.Uniform["camera"] = gl.GetUniformLocation(obj.Program, gl.Str("camera\x00"))
	obj.Uniform["model"] = gl.GetUniformLocation(obj.Program, gl.Str("model\x00"))
	gl.BindFragDataLocation(obj.Program, 0, gl.Str("outputColor\x00"))
}

func (obj *BaseObj) GetVertices() *[]float32 {
	return obj.Vertices
}

func (obj *BaseObj) SetVertices(vertices *[]float32) {
	obj.Vertices = vertices

	var vao uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(*vertices)*4, // 4 is the size of float32
		gl.Ptr(*vertices),
		gl.STATIC_DRAW,
	)

	vertAttrib := uint32(0) // 0 is the index of variable "vert" defined in vShader

	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(
		vertAttrib,
		3,
		gl.FLOAT,
		false,
		3*4, // 4 is the size of float32, and there're 5 floats per vertex in the vertex array.
		0,
	)

	obj.Vao = vao
}

func (obj *BaseObj) GetModel() mgl32.Mat4 {
	return obj.Model
}

func (obj *BaseObj) SetModel(model mgl32.Mat4) {
	obj.Model = model
}

func (obj *BaseObj) Render() {
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	gl.UseProgram(obj.Program)
	gl.UniformMatrix4fv(obj.Uniform["project"], 1, false, &(obj.ProgVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.Uniform["camera"], 1, false, &(obj.ProgVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.Uniform["model"], 1, false, &obj.Model[0])
	gl.BindVertexArray(obj.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*obj.Vertices)/3)) // 6: X,Y,Z
}

// getBaseObjVS returns the vertex shader of BaseObj
func getBaseObjVS() string {
	return fmt.Sprintf(
		`
		#version 330
		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;
		layout (location = 0) in vec3 vert;
		void main() {
			gl_Position = projection * camera * model * vec4(vert, 1);
		}
		%v`,
		"\x00",
	)
}

// getBaseObjFS returns the fragment shader of BaseObj
// outputColor is vec4(red, green, blue, alpha) which usually are
// uniform varialbes, but we set it a constant here to create a simpliest
// base object.
// We set the color gray so that it'll be visible in both black and white
// background.
func getBaseObjFS() string {
	return fmt.Sprintf(
		`
		#version 330
		out vec4 outputColor;
		void main() {
			outputColor = vec4(0.5, 0.5, 0.5, 1.0);
		}
		%v`,
		"\x00",
	)
}
