package sgl

import (
	"fmt"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Object is the basic interface of every renderable object.
type Object interface {
	// GetProgram gets the program of the object.
	GetProgram() uint32

	// SetProgram sets the program of the object using an existing program.
	SetProgram(program uint32)

	// GetProgram gets the program variables of the object.
	// The return value should be a program variable struct.
	// Program variable struct is the customizes struct that contains all
	// the uniform variables that will be used in the shader program of a
	// certain object.
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
	// should have already been prepared when calling SetProgVar().
	Render()
}

// BaseObjVar is the program variable struct for BaseObj.
// Every Object struct will have it's own program variable struct.
type BaseObjVar struct {
	Vp *Viewpoint
}

// BaseObj is an Object that keeps the basic info of a Object which
// will render a wireframe object.
// BaseObj could be embedded into other Object struct to provide common
// methods like GetProgram(), GetProgVar(), GetVertices(), GetModel()
// and SetProgram(). So when we're developing a new Object struct, we only
// need to implement SetProgVar(), SetVertices() and Render().
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

// NewBaseObj return a BaseObj instance with its program.
// Each Object usually only be able to use one kind of program,
// so it's a nice practice to use NewXXX() to create an Object instance
// that contains the default program.
// However, we could always use SetProgram() to set the program if there
// is an existing program which is already been compiled.
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

// SimpleObjVar is the program variable struct for SimpleObj.
type SimpleObjVar struct {
	Red   float32
	Green float32
	Blue  float32
	Vp    *Viewpoint
	Ls    *LightSrc
	Mt    *Material
}

// SimpleObj is the Object struct that will render a mono color object which
// has certain material properties. The mono color object will reflect the light
// from a single light source.
type SimpleObj struct {
	progVar SimpleObjVar
	BaseObj
}

// NewSimpleObj returns a SimpleObj instance with its program.
func NewSimpleObj() Object {
	obj := &SimpleObj{}
	obj.SetProgram(MakeProgram(getSimpleObjVS(), getSimpleObjFS()))

	return obj
}

func (obj *SimpleObj) SetProgVar(progVar interface{}) {
	if pv, ok := progVar.(SimpleObjVar); ok {
		obj.progVar = pv
	} else {
		panic("progVar is not a SimpleObjVar")
	}

	obj.Uniform = map[string]int32{}

	obj.Uniform["project"] = gl.GetUniformLocation(obj.Program, gl.Str("projection\x00"))
	obj.Uniform["camera"] = gl.GetUniformLocation(obj.Program, gl.Str("camera\x00"))
	obj.Uniform["model"] = gl.GetUniformLocation(obj.Program, gl.Str("model\x00"))
	obj.Uniform["lightPos"] = gl.GetUniformLocation(obj.Program, gl.Str("lightPos\x00"))
	obj.Uniform["lightColor"] = gl.GetUniformLocation(obj.Program, gl.Str("lightColor\x00"))
	obj.Uniform["lightIntensity"] = gl.GetUniformLocation(obj.Program, gl.Str("lightIntensity\x00"))
	obj.Uniform["viewPos"] = gl.GetUniformLocation(obj.Program, gl.Str("viewPos\x00"))
	obj.Uniform["red"] = gl.GetUniformLocation(obj.Program, gl.Str("red\x00"))
	obj.Uniform["green"] = gl.GetUniformLocation(obj.Program, gl.Str("green\x00"))
	obj.Uniform["blue"] = gl.GetUniformLocation(obj.Program, gl.Str("blue\x00"))
	obj.Uniform["materialAmbient"] = gl.GetUniformLocation(obj.Program, gl.Str("materialAmbient\x00"))
	obj.Uniform["materialDiffuse"] = gl.GetUniformLocation(obj.Program, gl.Str("materialDiffuse\x00"))
	obj.Uniform["materialSpecular"] = gl.GetUniformLocation(obj.Program, gl.Str("materialSpecular\x00"))
	obj.Uniform["materialShininess"] = gl.GetUniformLocation(obj.Program, gl.Str("materialShininess\x00"))
	gl.BindFragDataLocation(obj.Program, 0, gl.Str("outputColor\x00"))
}

func (obj *SimpleObj) SetVertices(vertices *[]float32) {
	newVertices := AddNormal(*vertices)
	obj.Vertices = &newVertices

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(*obj.Vertices)*4, // 4 is the size of float32
		gl.Ptr(*obj.Vertices),
		gl.STATIC_DRAW,
	)

	vertAttrib := uint32(0) // 0 is the index of variable "aPos" defined in vShader
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(
		vertAttrib,
		3,
		gl.FLOAT,
		false,
		6*4, // 4 is the size of float32, and there are 6 floats per vertex in the vertex array.
		0,
	)
	normal := uint32(1) // 1 is the index of variable "aNormal" defined in vShader
	gl.EnableVertexAttribArray(normal)
	gl.VertexAttribPointerWithOffset(
		normal,
		3,
		gl.FLOAT,
		false,
		6*4, // 4 is the size of float32, and there are 6 floats per vertex in the vertex array.
		3*4, // use offset 3*4 because we use only last 3 floats of each vertex here.
	)
	obj.Vao = vao
}

func (obj *SimpleObj) Render() {
	gl.UseProgram(obj.Program)
	gl.UniformMatrix4fv(obj.Uniform["project"], 1, false, &(obj.progVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.Uniform["camera"], 1, false, &(obj.progVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.Uniform["model"], 1, false, &obj.Model[0])
	gl.Uniform3fv(obj.Uniform["lightPos"], 1, &(obj.progVar.Ls.Pos[0]))
	gl.Uniform3fv(obj.Uniform["lightColor"], 1, &(obj.progVar.Ls.Color[0]))
	gl.Uniform3fv(obj.Uniform["viewPos"], 1, &(obj.progVar.Vp.Eye[0]))
	gl.Uniform1f(obj.Uniform["lightIntensity"], obj.progVar.Ls.Intensity)
	gl.Uniform1f(obj.Uniform["red"], obj.progVar.Red)
	gl.Uniform1f(obj.Uniform["green"], obj.progVar.Green)
	gl.Uniform1f(obj.Uniform["blue"], obj.progVar.Blue)
	gl.Uniform3fv(obj.Uniform["materialAmbient"], 1, &(obj.progVar.Mt.Ambient[0]))
	gl.Uniform3fv(obj.Uniform["materialDiffuse"], 1, &(obj.progVar.Mt.Diffuse[0]))
	gl.Uniform3fv(obj.Uniform["materialSpecular"], 1, &(obj.progVar.Mt.Specular[0]))
	gl.Uniform1f(obj.Uniform["materialShininess"], obj.progVar.Mt.Shininess)
	gl.BindVertexArray(obj.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*obj.Vertices)/6)) // 6: X,Y,Z,NX,NY,NZ
}

// getSimpleObjVS returns the vertex shader of SimpleObj
func getSimpleObjVS() string {
	return fmt.Sprintf(
		`
		#version 330

		layout(location = 0) in vec3 aPos;
		layout(location = 1) in vec3 aNormal;

		out vec3 FragPos;
		out vec3 Normal;

		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;

		void main() {
    		FragPos = vec3(model * vec4(aPos, 1.0));
    		Normal = mat3(transpose(inverse(model))) * aNormal;

    		gl_Position = projection * camera * vec4(FragPos, 1.0);
		}
		%v`,
		"\x00",
	)
}

// getSimpleObjFS returns the fragment shader of SimpleObj
// outputColor is vec4(red, green, blue, alpha) which usually are
// uniform varialbes, but we set it a constant here to create a simpliest
// base object.
// We set the color gray so that it'll be visible in both black and white
// background.
func getSimpleObjFS() string {
	return fmt.Sprintf(
		`
		#version 330
		out vec4 FragColor;

		in vec3 Normal;
		in vec3 FragPos;

		uniform vec3 viewPos;

		uniform float red;
		uniform float green;
		uniform float blue;

		uniform vec3 lightPos;
		uniform vec3 lightColor;
		uniform float lightIntensity;

		uniform vec3 materialAmbient;
		uniform vec3 materialDiffuse;
		uniform vec3 materialSpecular;
		uniform float materialShininess;

		void main() {
    		vec3 objectColor = vec3(red, green, blue);

			// ambient
    		vec3 ambient = lightColor * materialAmbient;

			// diffuse 
    		vec3 norm = normalize(Normal);
    		vec3 lightDir = normalize(lightPos - FragPos);
    		float diff = max(dot(norm, lightDir), 0.0);
    		vec3 diffuse = (lightIntensity * lightColor) * (diff * materialDiffuse);

			// specular
    		vec3 viewDir = normalize(viewPos - FragPos);
    		vec3 reflectDir = reflect(-lightDir, norm);
    		float spec = pow(max(dot(viewDir, reflectDir), 0.0), materialShininess);
    		vec3 specular = lightColor * (spec * materialSpecular);

    		vec3 result = (ambient + diffuse + specular) * objectColor;
    		FragColor = vec4(result, 1.0);
		}
		%v`,
		"\x00",
	)
}
