package objects

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/burwei/sgl"
	"github.com/go-gl/gl/all-core/gl"
)

type TexCubeObjVar struct {
	TextureSrc string
	Vp         *sgl.Viewpoint
}

type TexCubeObj struct {
	progVar TexCubeObjVar
	texture uint32
	sgl.BaseObj
}

func NewTexCubeObj() sgl.Object {
	obj := &TexCubeObj{}
	obj.SetProgram(sgl.MakeProgramFromFile("./objects/tex_cube_obj.vert", "./objects/tex_cube_obj.frag"))

	return obj
}

func (obj *TexCubeObj) SetProgVar(progVar interface{}) {
	if pv, ok := progVar.(TexCubeObjVar); ok {
		obj.progVar = pv
	} else {
		panic("progVar is not a TexCubeObjProgVar")
	}

	obj.setTexture()

	obj.Uniform = map[string]int32{}

	obj.Uniform["project"] = gl.GetUniformLocation(obj.Program, gl.Str("projection\x00"))
	obj.Uniform["camera"] = gl.GetUniformLocation(obj.Program, gl.Str("camera\x00"))
	obj.Uniform["model"] = gl.GetUniformLocation(obj.Program, gl.Str("model\x00"))
	obj.Uniform["tex"] = gl.GetUniformLocation(obj.Program, gl.Str("tex\x00"))
	gl.BindFragDataLocation(obj.Program, 0, gl.Str("outputColor\x00"))
}

func (obj *TexCubeObj) SetVertices(vertices *[]float32) {
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
		5*4, // 4 is the size of float32, and there're 5 floats per vertex in the vertex array.
		0,
	)

	texCoordAttrib := uint32(1) // 1 is the index of variable "vertTexCoord" defined in vShader
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(
		texCoordAttrib,
		2,
		gl.FLOAT,
		false,
		5*4,
		3*4, // use offset 3*4 because we use only last 2 floats of each vertex here.
	)

	obj.Vao = vao
}

func (obj *TexCubeObj) Render() {
	gl.UseProgram(obj.Program)
	gl.UniformMatrix4fv(obj.Uniform["project"], 1, false, &(obj.progVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.Uniform["camera"], 1, false, &(obj.progVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.Uniform["model"], 1, false, &obj.Model[0])
	gl.BindVertexArray(obj.Vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, obj.texture)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*obj.Vertices)/5)) // 5: X,Y,Z,U,V
}

func (obj *TexCubeObj) setTexture() {
	imgFile, err := os.Open(obj.progVar.TextureSrc)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic(fmt.Errorf("unsupported stride"))
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	obj.texture = texture
}
