package simplegl

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Model interface{
	SetProgram()
	SetMatrixes()
	SetTexture()
	SetVao()
	Render()
}

type SimpleModel struct {
	Program        uint32
	Vao            uint32
	Vertices       *[]float32
	Model          mgl32.Mat4
	ModelUniform   int32
	Texture        uint32
	TextureUniform int32
}

func (m *SimpleModel) SetProgram(vertexShaderSource, fragmentShaderSource string) {
	vertexShader, err := m.compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := m.compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	gl.UseProgram(program)

	m.Program = program
}

func (m *SimpleModel) compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (m *SimpleModel) SetMatrixes(vp *SimpleViewPoint) {
	vp.Projection = mgl32.Perspective(vp.Fovy, vp.Aspect, vp.Near, vp.Far)
	vp.ProjectionUniform = gl.GetUniformLocation(m.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(vp.ProjectionUniform, 1, false, &vp.Projection[0])

	vp.Camera = mgl32.LookAtV(vp.Eye, vp.Center, vp.Top)
	vp.CameraUniform = gl.GetUniformLocation(m.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(vp.CameraUniform, 1, false, &vp.Camera[0])

	m.Model = mgl32.Ident4()
	m.ModelUniform = gl.GetUniformLocation(m.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(m.ModelUniform, 1, false, &m.Model[0])

	m.TextureUniform = gl.GetUniformLocation(m.Program, gl.Str("tex\x00"))
	gl.Uniform1i(m.TextureUniform, 0)

	gl.BindFragDataLocation(m.Program, 0, gl.Str("outputColor\x00"))
}

func (m *SimpleModel) SetVao(vertices *[]float32) {
	m.Vertices = vertices

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(*vertices)*4, gl.Ptr(*vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(m.Program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	texCoordAttrib := uint32(gl.GetAttribLocation(m.Program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	m.Vao = vao
}

func (m *SimpleModel) SetTexture(file string) {
	imgFile, err := os.Open(file)
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

	m.Texture = texture
}

func (m *SimpleModel) Render(vp *SimpleViewPoint) {
	gl.UseProgram(m.Program)
	gl.UniformMatrix4fv(vp.ProjectionUniform, 1, false, &vp.Projection[0])
	gl.UniformMatrix4fv(vp.CameraUniform, 1, false, &vp.Camera[0])
	gl.UniformMatrix4fv(m.ModelUniform, 1, false, &m.Model[0])
	gl.BindVertexArray(m.Vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, m.Texture)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*m.Vertices)/5)) // 5: X,Y,Z,U,V
}