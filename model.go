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

type Model interface {
	SetProgram()
	SetMatrixes()
	SetVao()
	Render()
}

type SimpleModel struct {
	Program      uint32
	Vao          uint32
	Vertices     *[]float32
	Model        mgl32.Mat4
	ModelUniform int32
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
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(*vertices)*4, // 4 is the size of float32
		gl.Ptr(*vertices),
		gl.STATIC_DRAW,
	)

	vertAttrib := uint32(0) // 0 is the index of variable "vert" defined in GLSL
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(
		vertAttrib,
		3,
		gl.FLOAT,
		false,
		3*4, // 4 is the size of float32, and there're 5 floats per vertex in the vertex array.
		0,
	)
	m.Vao = vao
}

func (m *SimpleModel) Render(vp *SimpleViewPoint) {
	gl.UseProgram(m.Program)
	gl.UniformMatrix4fv(vp.ProjectionUniform, 1, false, &vp.Projection[0])
	gl.UniformMatrix4fv(vp.CameraUniform, 1, false, &vp.Camera[0])
	gl.UniformMatrix4fv(m.ModelUniform, 1, false, &m.Model[0])
	gl.BindVertexArray(m.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*m.Vertices)/3)) // 3: X,Y,Z
}

type UniTexModel struct {
	Program        uint32
	Vao            uint32
	Vertices       *[]float32
	Model          mgl32.Mat4
	ModelUniform   int32
	Texture        uint32
	TextureUniform int32
}

func (m *UniTexModel) SetProgram(vertexShaderSource, fragmentShaderSource string) {
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

func (m *UniTexModel) compileShader(source string, shaderType uint32) (uint32, error) {
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

func (m *UniTexModel) SetMatrixes(vp *SimpleViewPoint) {
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

func (m *UniTexModel) SetVao(vertices *[]float32) {
	m.Vertices = vertices

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

	vertAttrib := uint32(0) // 0 is the index of variable "vert" defined in GLSL
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(
		vertAttrib,
		3,
		gl.FLOAT,
		false,
		5*4, // 4 is the size of float32, and there're 5 floats per vertex in the vertex array.
		0,
	)

	texCoordAttrib := uint32(1) // 1 is the index of variable "vertTexCoord" defined in GLSL
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(
		texCoordAttrib,
		2,
		gl.FLOAT,
		false,
		5*4,
		3*4, // use offset 3*4 because we use only last 2 floats of each vertex here.
	)

	m.Vao = vao
}

func (m *UniTexModel) SetTexture(file string) {
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

func (m *UniTexModel) Render(vp *SimpleViewPoint) {
	gl.UseProgram(m.Program)
	gl.UniformMatrix4fv(vp.ProjectionUniform, 1, false, &vp.Projection[0])
	gl.UniformMatrix4fv(vp.CameraUniform, 1, false, &vp.Camera[0])
	gl.UniformMatrix4fv(m.ModelUniform, 1, false, &m.Model[0])
	gl.BindVertexArray(m.Vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, m.Texture)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*m.Vertices)/5)) // 6: X,Y,Z,U,V
}

type SimpleLightModel struct {
	Program      uint32
	Vao          uint32
	Vertices     *[]float32
	Model        mgl32.Mat4
	ModelUniform int32
}

func (m *SimpleLightModel) SetProgram(vertexShaderSource, fragmentShaderSource string) {
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

func (m *SimpleLightModel) compileShader(source string, shaderType uint32) (uint32, error) {
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

func (m *SimpleLightModel) SetMatrixes(vp *SimpleViewPoint, ls *SimpleLightSrc) {
	vp.Projection = mgl32.Perspective(vp.Fovy, vp.Aspect, vp.Near, vp.Far)
	vp.ProjectionUniform = gl.GetUniformLocation(m.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(vp.ProjectionUniform, 1, false, &vp.Projection[0])

	vp.Camera = mgl32.LookAtV(vp.Eye, vp.Center, vp.Top)
	vp.CameraUniform = gl.GetUniformLocation(m.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(vp.CameraUniform, 1, false, &vp.Camera[0])

	m.Model = mgl32.Ident4()
	m.ModelUniform = gl.GetUniformLocation(m.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(m.ModelUniform, 1, false, &m.Model[0])

	ls.PosUniform = gl.GetUniformLocation(m.Program, gl.Str("lightPos\x00"))
	gl.Uniform3fv(ls.PosUniform, 1, &ls.Pos[0])
	ls.ColorUniform = gl.GetUniformLocation(m.Program, gl.Str("lightColor\x00"))
	gl.Uniform3fv(ls.ColorUniform, 1, &ls.Color[0])
	vp.EyePosUniform = gl.GetUniformLocation(m.Program, gl.Str("viewPos\x00"))
	gl.Uniform3fv(vp.EyePosUniform, 1, &vp.Eye[0])
	ls.AmbientStrengthUniform = gl.GetUniformLocation(m.Program, gl.Str("ambientStrength\x00"))
	gl.Uniform1f(ls.AmbientStrengthUniform, ls.AmbientStrength)
	ls.SpecularStrengthUniform = gl.GetUniformLocation(m.Program, gl.Str("specularStrength\x00"))
	gl.Uniform1f(ls.ColorUniform, ls.SpecularStrength)
	ls.ShininessUniform = gl.GetUniformLocation(m.Program, gl.Str("shininess\x00"))
	gl.Uniform1f(ls.ShininessUniform, ls.Shininess)
	
	gl.BindFragDataLocation(m.Program, 0, gl.Str("outputColor\x00"))
}

func (m *SimpleLightModel) SetVao(vertices *[]float32) {
	newVertices := m.addNormal(*vertices)
	m.Vertices = &newVertices

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(*m.Vertices)*4, // 4 is the size of float32
		gl.Ptr(*m.Vertices),
		gl.STATIC_DRAW,
	)

	vertAttrib := uint32(0) // 0 is the index of variable "aPos" defined in GLSL
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(
		vertAttrib,
		3,
		gl.FLOAT,
		false,
		6*4, // 4 is the size of float32, and there're 6 floats per vertex in the vertex array.
		0,
	)
	normal := uint32(1) // 1 is the index of variable "aNormal" defined in GLSL
	gl.EnableVertexAttribArray(normal)
	gl.VertexAttribPointerWithOffset(
		normal,
		3,
		gl.FLOAT,
		false,
		6*4, // 4 is the size of float32, and there're 3 floats per vertex in the vertex array.
		3*4, // use offset 3*4 because we use only last 3 floats of each vertex here.
	)
	m.Vao = vao
}

func (m *SimpleLightModel) Render(vp *SimpleViewPoint, ls *SimpleLightSrc) {
	gl.UseProgram(m.Program)
	gl.UniformMatrix4fv(vp.ProjectionUniform, 1, false, &vp.Projection[0])
	gl.UniformMatrix4fv(vp.CameraUniform, 1, false, &vp.Camera[0])
	gl.UniformMatrix4fv(m.ModelUniform, 1, false, &m.Model[0])
	gl.UniformMatrix3fv(ls.PosUniform, 1, false, &ls.Pos[0])
	gl.UniformMatrix3fv(ls.ColorUniform, 1, false, &ls.Color[0])
	gl.UniformMatrix3fv(vp.EyePosUniform, 1, false, &vp.Eye[0])	
	gl.Uniform1f(ls.AmbientStrengthUniform, ls.AmbientStrength)
	gl.Uniform1f(ls.SpecularStrengthUniform, ls.SpecularStrength)
	gl.Uniform1f(ls.ShininessUniform, ls.Shininess)
	gl.BindVertexArray(m.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*m.Vertices)/6)) // 6: X,Y,Z,NX,NY,NZ
}

func (m *SimpleLightModel) addNormal(vertices []float32) []float32 {
	newVertices := []float32{}
	if len(vertices)%9 != 0 {
		return vertices
	}
	// One vertex contains 3 float values, three vertices contain 9 float values
	// Three vertices construct a plane(triangle), 
	// therefore we'll get one normal every three points
	for i := 0; i < len(vertices); i += 9 {
		// pt1 : ( vertices[i], vertices[i+1], vertices[i+2] )
		// pt2 : ( vertices[i+3], vertices[i+4], vertices[i+5] )
		// pt3 : ( vertices[i+6], vertices[i+7], vertices[i+8] )
		// vector: origin -> center of the triangle
		center := []float32{
			(vertices[i] + vertices[i+3] + vertices[i+6]) / 3,
			(vertices[i+1] + vertices[i+4] + vertices[i+7]) / 3,
			(vertices[i+2] + vertices[i+5] + vertices[i+8]) / 3,
		}
		// vector: pt1 -> pt2
		vec1 := []float32{
			vertices[i+3] - vertices[i],
			vertices[i+4] - vertices[i+1],
			vertices[i+5] - vertices[i+2],
		}
		// vector: pt2 -> pt3
		vec2 := []float32{
			vertices[i+6] - vertices[i+3],
			vertices[i+7] - vertices[i+4],
			vertices[i+8] - vertices[i+5],
		}
		// normal = vec1 x vec2 (cross product)
		normal := []float32{
			vec1[1]*vec2[2] - vec1[2]*vec2[1],
			vec1[2]*vec2[0] - vec1[0]*vec2[2],
			vec1[0]*vec2[1] - vec1[1]*vec2[0],
		}
		// check if normal . center (dot product) is negative
		dot := normal[0]*center[0]+normal[1]*center[1]+normal[2]*center[2]
		if dot < 0 {
			normal = []float32{
				-normal[0],
				-normal[1],
				-normal[2],
			}
		}
		// newPt1
		newVertices = append(newVertices, vertices[i])
		newVertices = append(newVertices, vertices[i+1])
		newVertices = append(newVertices, vertices[i+2])
		newVertices = append(newVertices, normal[0])
		newVertices = append(newVertices, normal[1])
		newVertices = append(newVertices, normal[2])
		// newPt2
		newVertices = append(newVertices, vertices[i+3])
		newVertices = append(newVertices, vertices[i+4])
		newVertices = append(newVertices, vertices[i+5])
		newVertices = append(newVertices, normal[0])
		newVertices = append(newVertices, normal[1])
		newVertices = append(newVertices, normal[2])
		// newPt3
		newVertices = append(newVertices, vertices[i+6])
		newVertices = append(newVertices, vertices[i+7])
		newVertices = append(newVertices, vertices[i+8])
		newVertices = append(newVertices, normal[0])
		newVertices = append(newVertices, normal[1])
		newVertices = append(newVertices, normal[2])
	}
	return newVertices
}
