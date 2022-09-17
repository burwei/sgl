package objects

import (
	"github.com/burwei/sgl"
	"github.com/go-gl/gl/all-core/gl"
)

type SimpleObj struct {
	progVar SimpleObjProgVar
	sgl.BaseObj
}

type SimpleObjProgVar struct {
	Red   float32
	Green float32
	Blue  float32
	Vp    *sgl.Viewpoint
	Ls    *sgl.LightSrc
	Mt    *sgl.Material
}

func NewSimpleObj() sgl.Object {
	obj := &SimpleObj{}
	obj.SetProgram(sgl.MakeProgramFromFile("./objects/simpleobj.vert", "./objects/simpleobj.frag"))

	return obj
}

func (obj *SimpleObj) SetProgVar(progVar interface{}) {
	if pv, ok := progVar.(SimpleObjProgVar); ok {
		obj.progVar = pv
	} else {
		panic("progVar is not a SimpleObjProgVar")
	}

	obj.Program = sgl.MakeProgramFromFile("./objects/simpleobj.vert", "./objects/simpleobj.frag")

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
	newVertices := sgl.AddNormal(*vertices)
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
