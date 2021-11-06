package simplegl

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)


type Uniforms map[string]int32

type Object interface {
	// Compile the program using vShader and fShader.
	// The parameters might not be the same for different objects
	// PrepareProgram(...)
	
	// Set uniform variables in GLSL.
	// The parameters might not be the same for different objects
	// SetUniforms(...)
	
	// Set vao and vbo. 
	// The parameters might not be the same for different objects
	// SetVertices(...)
	
	// Refresh uniform variables and draw the object.
	// All references of the variables that would change the
	// object's states in the main loop should already have been 
	// stored in SetUniforms(), so that we won't do extra work here 
	// to get the best performance.
	Render()
}

type BasicNoLightObject struct {
	Program  uint32
	Vao      uint32
	Vertices *[]float32
	Model    mgl32.Mat4
	Uni      *Uniforms
	Vp       *SimpleViewPoint
}

func (m *BasicNoLightObject) PrepareProgram(r float32, g float32, b float32) {
	vShader := fmt.Sprintf(
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
	fShader := fmt.Sprintf(
		`
		#version 330
		out vec4 outputColor;
		void main() {
			outputColor = vec4(%.3f, %.3f, %.3f, 1.0);
		}
		%v`,
		r,
		g,
		b,
		"\x00",
	)
	m.Program = MakeProgram(vShader, fShader)
}

func (m *BasicNoLightObject) SetUniforms(vp *SimpleViewPoint) {
	m.Uni = &Uniforms{}
	m.Vp = vp

	(*m.Uni)["project"] = gl.GetUniformLocation(m.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv((*m.Uni)["project"], 1, false, &(m.Vp.Projection[0]))

	vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
	(*m.Uni)["camera"] = gl.GetUniformLocation(m.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv((*m.Uni)["camera"], 1, false, &(m.Vp.Camera[0]))

	m.Model = mgl32.Ident4()
	(*m.Uni)["model"] = gl.GetUniformLocation(m.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv((*m.Uni)["model"], 1, false, &m.Model[0])

	gl.BindFragDataLocation(m.Program, 0, gl.Str("outputColor\x00"))
}

func (m *BasicNoLightObject) SetVertices(vertices *[]float32) {
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
	m.Vao = vao
}

func (m *BasicNoLightObject) Render() {
	gl.UseProgram(m.Program)
	gl.UniformMatrix4fv((*m.Uni)["project"], 1, false, &(m.Vp.Projection[0]))
	gl.UniformMatrix4fv((*m.Uni)["camera"], 1, false, &(m.Vp.Camera[0]))
	gl.UniformMatrix4fv((*m.Uni)["model"], 1, false, &m.Model[0])
	gl.BindVertexArray(m.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*m.Vertices)/3)) // 3: X,Y,Z
}

type BasicObject struct {
	Program  uint32
	Vao      uint32
	Vertices *[]float32
	Model    mgl32.Mat4
	Uni      *Uniforms
	Vp       *SimpleViewPoint
	Ls       *SimpleLightSrc
}

func (m *BasicObject) PrepareProgram(r float32, g float32, b float32) {
	vShader := fmt.Sprintf(
		`
		#version 330 core
		layout (location = 0) in vec3 aPos;
		layout (location = 1) in vec3 aNormal;

		out vec3 FragPos;
		out vec3 Normal;

		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;

		void main()
		{
			FragPos = vec3(model * vec4(aPos, 1.0));
			Normal = mat3(transpose(inverse(model))) * aNormal;   

			gl_Position = projection * camera * vec4(FragPos, 1.0);
		}
		%v`,
		"\x00",
	)
	fShader := fmt.Sprintf(
		`
		#version 330 core
		out vec4 FragColor;

		in vec3 Normal;  
		in vec3 FragPos;  
		
		uniform float ambientStrength;
		uniform float specularStrength;
		uniform float shininess; 
		uniform vec3 lightPos; 
		uniform vec3 lightColor;
		uniform vec3 viewPos;

		void main()
		{
			vec3 objectColor = vec3(%.3f, %.3f, %.3f);

			// ambient
			vec3 ambient = ambientStrength * lightColor;
				
			// diffuse 
			vec3 norm = normalize(Normal);
			vec3 lightDir = normalize(lightPos - FragPos);
			float diff = max(dot(norm, lightDir), 0.0);
			vec3 diffuse = diff * lightColor;
				
			// specular
			vec3 viewDir = normalize(viewPos - FragPos);
			vec3 reflectDir = reflect(-lightDir, norm);  
			float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);
			vec3 specular = specularStrength * spec * lightColor;  
				
			vec3 result = (ambient + diffuse + specular) * objectColor;
			FragColor = vec4(result, 1.0);
		} 
		%v`,
		r,
		g,
		b,
		"\x00",
	)
	m.Program = MakeProgram(vShader, fShader)
}

func (m *BasicObject) SetUniforms(vp *SimpleViewPoint, ls *SimpleLightSrc) {
	m.Uni = &Uniforms{}
	m.Vp = vp
	m.Ls = ls

	(*m.Uni)["project"] = gl.GetUniformLocation(m.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv((*m.Uni)["project"], 1, false, &(m.Vp.Projection[0]))

	(*m.Uni)["camera"] = gl.GetUniformLocation(m.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv((*m.Uni)["camera"], 1, false, &(m.Vp.Camera[0]))

	m.Model = mgl32.Ident4()
	(*m.Uni)["model"] = gl.GetUniformLocation(m.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv((*m.Uni)["model"], 1, false, &m.Model[0])

	(*m.Uni)["lightPos"] = gl.GetUniformLocation(m.Program, gl.Str("lightPos\x00"))
	gl.Uniform3fv((*m.Uni)["lightPos"], 1, &(m.Ls.Pos[0]))
	(*m.Uni)["lightColor"] = gl.GetUniformLocation(m.Program, gl.Str("lightColor\x00"))
	gl.Uniform3fv((*m.Uni)["lightColor"], 1, &(m.Ls.Color[0]))
	(*m.Uni)["viewPos"] = gl.GetUniformLocation(m.Program, gl.Str("viewPos\x00"))
	gl.Uniform3fv((*m.Uni)["viewPos"], 1, &(m.Vp.Eye[0]))
	(*m.Uni)["ambientStrength"] = gl.GetUniformLocation(m.Program, gl.Str("ambientStrength\x00"))
	gl.Uniform1f((*m.Uni)["ambientStrength"], m.Ls.AmbientStrength)
	(*m.Uni)["specularStrength"] = gl.GetUniformLocation(m.Program, gl.Str("specularStrength\x00"))
	gl.Uniform1f((*m.Uni)["specularStrength"], m.Ls.SpecularStrength)
	(*m.Uni)["shininess"] = gl.GetUniformLocation(m.Program, gl.Str("shininess\x00"))
	gl.Uniform1f((*m.Uni)["shininess"], m.Ls.Shininess)

	gl.BindFragDataLocation(m.Program, 0, gl.Str("outputColor\x00"))
}

func (m *BasicObject) SetVertices(vertices *[]float32) {
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

	vertAttrib := uint32(0) // 0 is the index of variable "aPos" defined in vShader
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(
		vertAttrib,
		3,
		gl.FLOAT,
		false,
		6*4, // 4 is the size of float32, and there're 6 floats per vertex in the vertex array.
		0,
	)
	normal := uint32(1) // 1 is the index of variable "aNormal" defined in vShader
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

func (m *BasicObject) Render() {
	gl.UseProgram(m.Program)
	gl.UniformMatrix4fv((*m.Uni)["project"], 1, false, &(m.Vp.Projection[0]))
	gl.UniformMatrix4fv((*m.Uni)["camera"], 1, false, &(m.Vp.Camera[0]))
	gl.UniformMatrix4fv((*m.Uni)["model"], 1, false, &m.Model[0])
	gl.UniformMatrix3fv((*m.Uni)["lightPos"], 1, false, &(m.Ls.Pos[0]))
	gl.UniformMatrix3fv((*m.Uni)["lightColor"], 1, false, &(m.Ls.Color[0]))
	gl.UniformMatrix3fv((*m.Uni)["viewPos"], 1, false, &(m.Vp.Eye[0]))
	gl.Uniform1f((*m.Uni)["ambientStrength"], m.Ls.AmbientStrength)
	gl.Uniform1f((*m.Uni)["specularStrength"], m.Ls.SpecularStrength)
	gl.Uniform1f((*m.Uni)["shininess"], m.Ls.Shininess)
	gl.BindVertexArray(m.Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*m.Vertices)/6)) // 6: X,Y,Z,NX,NY,NZ
}

func (m *BasicObject) addNormal(vertices []float32) []float32 {
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
		dot := normal[0]*center[0] + normal[1]*center[1] + normal[2]*center[2]
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

type BasicTexObject struct {
	Program  uint32
	Vao      uint32
	Vertices *[]float32
	Model    mgl32.Mat4
	Texture  uint32
	Uni      *Uniforms
	Vp       *SimpleViewPoint
}

func (m *BasicTexObject) PrepareProgram() {
	vShader := fmt.Sprintf(
		`
		#version 330

		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;

		layout (location = 0) in vec3 vert;
		layout (location = 1) in vec2 vertTexCoord;

		out vec2 fragTexCoord;

		void main() {
			fragTexCoord = vertTexCoord;
			gl_Position = projection * camera * model * vec4(vert, 1);
		}
		%v`,
		"\x00",
	)
	fShader := fmt.Sprintf(
		`
		#version 330

		uniform sampler2D tex;

		in vec2 fragTexCoord;

		out vec4 outputColor;

		void main() {
			outputColor = texture(tex, fragTexCoord);
		}
		%v`,
		"\x00",
	)
	m.Program = MakeProgram(vShader, fShader)
}

func (m *BasicTexObject) SetUniforms(vp *SimpleViewPoint) {
	m.Uni = &Uniforms{}
	m.Vp = vp

	vp.Projection = mgl32.Perspective(vp.Fovy, vp.Aspect, vp.Near, vp.Far)
	(*m.Uni)["project"] = gl.GetUniformLocation(m.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv((*m.Uni)["project"], 1, false, &(m.Vp.Projection[0]))

	vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
	(*m.Uni)["camera"] = gl.GetUniformLocation(m.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv((*m.Uni)["camera"], 1, false, &(m.Vp.Camera[0]))

	m.Model = mgl32.Ident4()
	(*m.Uni)["model"] = gl.GetUniformLocation(m.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv((*m.Uni)["model"], 1, false, &m.Model[0])

	(*m.Uni)["tex"] = gl.GetUniformLocation(m.Program, gl.Str("tex\x00"))
	gl.Uniform1i((*m.Uni)["tex"], 0)

	gl.BindFragDataLocation(m.Program, 0, gl.Str("outputColor\x00"))
}

func (m *BasicTexObject) SetVertices(vertices *[]float32) {
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

	m.Vao = vao
}

func (m *BasicTexObject) SetTexture(file string) {
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

func (m *BasicTexObject) Render() {
	gl.UseProgram(m.Program)
	gl.UniformMatrix4fv((*m.Uni)["project"], 1, false, &(m.Vp.Projection[0]))
	gl.UniformMatrix4fv((*m.Uni)["camera"], 1, false, &(m.Vp.Camera[0]))
	gl.UniformMatrix4fv((*m.Uni)["model"], 1, false, &m.Model[0])
	gl.BindVertexArray(m.Vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, m.Texture)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*m.Vertices)/5)) // 6: X,Y,Z,U,V
}
