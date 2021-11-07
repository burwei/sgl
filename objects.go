package simplegl

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object interface {
	// SetProgramVar sets the program's constant and uniform variables.
	SetProgramVar(interface{})

	// PrepareProgram prepares shader program and uniform variables.
	PrepareProgram()

	// SetVertices sets vao and vbo.
	SetVertices(*[]float32)

	// GetModel gets the model of the object.
	GetModel() mgl32.Mat4

	// SetModel sets the model of the object.
	SetModel(mgl32.Mat4)

	// Render refreshs uniform variables and draw the object.
	// All references of the variables that would change the
	// object's states in the main loop (i.e. uniform variables)
	// should have already been prepared when calling SetProgramVar(),
	// thus we won't do any extra work here to get the best performance.
	Render()
}

type BasicNoLightObj struct {
	program  uint32
	vao      uint32
	vertices *[]float32
	model    mgl32.Mat4
	uni      map[string]int32
	progVar  BasicNoLightObjProgVar
}

type BasicNoLightObjProgVar struct {
	Red   float32
	Green float32
	Blue  float32
	Vp    *SimpleViewPoint
}

func (obj *BasicNoLightObj) SetProgramVar(progVar interface{}) {
	if pv, ok := progVar.(BasicNoLightObjProgVar); ok {
		obj.progVar = pv
	} else {
		panic("progVar is not a BasicNoLightObjProgVar")
	}
}

func (obj *BasicNoLightObj) PrepareProgram() {
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
		uniform float red;
		uniform float green;
		uniform float blue;
		void main() {
			outputColor = vec4(red, green, blue, 1.0);
		}
		%v`,
		"\x00",
	)
	obj.program = MakeProgram(vShader, fShader)

	obj.model = mgl32.Ident4()
	obj.uni = map[string]int32{}
	obj.uni["project"] = gl.GetUniformLocation(obj.program, gl.Str("projection\x00"))
	obj.uni["camera"] = gl.GetUniformLocation(obj.program, gl.Str("camera\x00"))
	obj.uni["model"] = gl.GetUniformLocation(obj.program, gl.Str("model\x00"))
	obj.uni["red"] = gl.GetUniformLocation(obj.program, gl.Str("red\x00"))
	obj.uni["green"] = gl.GetUniformLocation(obj.program, gl.Str("green\x00"))
	obj.uni["blue"] = gl.GetUniformLocation(obj.program, gl.Str("blue\x00"))

	gl.UniformMatrix4fv(obj.uni["project"], 1, false, &(obj.progVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.uni["camera"], 1, false, &(obj.progVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.uni["model"], 1, false, &obj.model[0])
	gl.BindFragDataLocation(obj.program, 0, gl.Str("outputColor\x00"))
}

func (obj *BasicNoLightObj) SetVertices(vertices *[]float32) {
	obj.vertices = vertices

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
	obj.vao = vao
}

func (obj *BasicNoLightObj) GetModel() mgl32.Mat4 {
	return obj.model
}

func (obj *BasicNoLightObj) SetModel(newModel mgl32.Mat4) {
	obj.model = newModel
}

func (obj *BasicNoLightObj) Render() {
	gl.UseProgram(obj.program)
	gl.UniformMatrix4fv(obj.uni["project"], 1, false, &(obj.progVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.uni["camera"], 1, false, &(obj.progVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.uni["model"], 1, false, &obj.model[0])
	gl.Uniform1f(obj.uni["red"], obj.progVar.Red)
	gl.Uniform1f(obj.uni["green"], obj.progVar.Green)
	gl.Uniform1f(obj.uni["blue"], obj.progVar.Blue)
	gl.BindVertexArray(obj.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*obj.vertices)/3)) // 3: X,Y,Z
}

type BasicObj struct {
	program  uint32
	vao      uint32
	vertices *[]float32
	model    mgl32.Mat4
	uni      map[string]int32
	progVar  BasicObjProgVar
}

type BasicObjProgVar struct {
	Red   float32
	Green float32
	Blue  float32
	Vp    *SimpleViewPoint
	Ls    *SimpleLightSrc
}

func (obj *BasicObj) SetProgramVar(progVar interface{}) {
	if pv, ok := progVar.(BasicObjProgVar); ok {
		obj.progVar = pv
	} else {
		panic("progVar is not a BasicObjProgVar")
	}
}

func (obj *BasicObj) PrepareProgram() {
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
		
		uniform float red;
		uniform float green;
		uniform float blue;
		uniform float ambientStrength;
		uniform float specularStrength;
		uniform float shininess; 
		uniform vec3 lightPos; 
		uniform vec3 lightColor;
		uniform vec3 viewPos;

		void main()
		{
			vec3 objectColor = vec3(red, green, blue);

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
		"\x00",
	)
	obj.program = MakeProgram(vShader, fShader)

	obj.model = mgl32.Ident4()
	obj.uni = map[string]int32{}
	obj.uni["project"] = gl.GetUniformLocation(obj.program, gl.Str("projection\x00"))
	obj.uni["camera"] = gl.GetUniformLocation(obj.program, gl.Str("camera\x00"))
	obj.uni["model"] = gl.GetUniformLocation(obj.program, gl.Str("model\x00"))
	obj.uni["lightPos"] = gl.GetUniformLocation(obj.program, gl.Str("lightPos\x00"))
	obj.uni["lightColor"] = gl.GetUniformLocation(obj.program, gl.Str("lightColor\x00"))
	obj.uni["viewPos"] = gl.GetUniformLocation(obj.program, gl.Str("viewPos\x00"))
	obj.uni["red"] = gl.GetUniformLocation(obj.program, gl.Str("red\x00"))
	obj.uni["green"] = gl.GetUniformLocation(obj.program, gl.Str("green\x00"))
	obj.uni["blue"] = gl.GetUniformLocation(obj.program, gl.Str("blue\x00"))
	obj.uni["ambientStrength"] = gl.GetUniformLocation(obj.program, gl.Str("ambientStrength\x00"))
	obj.uni["specularStrength"] = gl.GetUniformLocation(obj.program, gl.Str("specularStrength\x00"))
	obj.uni["shininess"] = gl.GetUniformLocation(obj.program, gl.Str("shininess\x00"))

	gl.UniformMatrix4fv(obj.uni["project"], 1, false, &(obj.progVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.uni["camera"], 1, false, &(obj.progVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.uni["model"], 1, false, &obj.model[0])
	gl.Uniform3fv(obj.uni["lightPos"], 1, &(obj.progVar.Ls.Pos[0]))
	gl.Uniform3fv(obj.uni["lightColor"], 1, &(obj.progVar.Ls.Color[0]))
	gl.Uniform3fv(obj.uni["viewPos"], 1, &(obj.progVar.Vp.Eye[0]))
	gl.Uniform1f(obj.uni["red"], obj.progVar.Red)
	gl.Uniform1f(obj.uni["green"], obj.progVar.Green)
	gl.Uniform1f(obj.uni["blue"], obj.progVar.Blue)
	gl.Uniform1f(obj.uni["ambientStrength"], obj.progVar.Ls.AmbientStrength)
	gl.Uniform1f(obj.uni["specularStrength"], obj.progVar.Ls.SpecularStrength)
	gl.Uniform1f(obj.uni["shininess"], obj.progVar.Ls.Shininess)
	gl.BindFragDataLocation(obj.program, 0, gl.Str("outputColor\x00"))
}

func (obj *BasicObj) SetVertices(vertices *[]float32) {
	newVertices := obj.addNormal(*vertices)
	obj.vertices = &newVertices

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(*obj.vertices)*4, // 4 is the size of float32
		gl.Ptr(*obj.vertices),
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
	obj.vao = vao
}

func (obj *BasicObj) GetModel() mgl32.Mat4 {
	return obj.model
}

func (obj *BasicObj) SetModel(newModel mgl32.Mat4) {
	obj.model = newModel
}

func (obj *BasicObj) Render() {
	gl.UseProgram(obj.program)
	gl.UniformMatrix4fv(obj.uni["project"], 1, false, &(obj.progVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.uni["camera"], 1, false, &(obj.progVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.uni["model"], 1, false, &obj.model[0])
	gl.UniformMatrix3fv(obj.uni["lightPos"], 1, false, &(obj.progVar.Ls.Pos[0]))
	gl.UniformMatrix3fv(obj.uni["lightColor"], 1, false, &(obj.progVar.Ls.Color[0]))
	gl.UniformMatrix3fv(obj.uni["viewPos"], 1, false, &(obj.progVar.Vp.Eye[0]))
	gl.Uniform1f(obj.uni["red"], obj.progVar.Red)
	gl.Uniform1f(obj.uni["green"], obj.progVar.Green)
	gl.Uniform1f(obj.uni["blue"], obj.progVar.Blue)
	gl.Uniform1f(obj.uni["ambientStrength"], obj.progVar.Ls.AmbientStrength)
	gl.Uniform1f(obj.uni["specularStrength"], obj.progVar.Ls.SpecularStrength)
	gl.Uniform1f(obj.uni["shininess"], obj.progVar.Ls.Shininess)
	gl.BindVertexArray(obj.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*obj.vertices)/6)) // 6: X,Y,Z,NX,NY,NZ
}

func (obj *BasicObj) addNormal(vertices []float32) []float32 {
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

type BasicTexObj struct {
	program  uint32
	vao      uint32
	vertices *[]float32
	model    mgl32.Mat4
	uni      map[string]int32
	progVar  BasicTexObjProgVar
	texture  uint32
}

type BasicTexObjProgVar struct {
	TextureSrc string
	Vp         *SimpleViewPoint
}

func (obj *BasicTexObj) SetProgramVar(progVar interface{}) {
	if pv, ok := progVar.(BasicTexObjProgVar); ok {
		obj.progVar = pv
	} else {
		panic("progVar is not a BasicTexObjProgVar")
	}
	obj.setTexture()
}

func (obj *BasicTexObj) PrepareProgram() {
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
	obj.program = MakeProgram(vShader, fShader)

	obj.model = mgl32.Ident4()
	obj.uni = map[string]int32{}
	obj.uni["project"] = gl.GetUniformLocation(obj.program, gl.Str("projection\x00"))
	obj.uni["camera"] = gl.GetUniformLocation(obj.program, gl.Str("camera\x00"))
	obj.uni["model"] = gl.GetUniformLocation(obj.program, gl.Str("model\x00"))
	obj.uni["tex"] = gl.GetUniformLocation(obj.program, gl.Str("tex\x00"))

	gl.UniformMatrix4fv(obj.uni["project"], 1, false, &(obj.progVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.uni["camera"], 1, false, &(obj.progVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.uni["model"], 1, false, &obj.model[0])
	gl.Uniform1i(obj.uni["tex"], 0)
	gl.BindFragDataLocation(obj.program, 0, gl.Str("outputColor\x00"))
}

func (obj *BasicTexObj) SetVertices(vertices *[]float32) {
	obj.vertices = vertices

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

	obj.vao = vao
}

func (obj *BasicTexObj) GetModel() mgl32.Mat4 {
	return obj.model
}

func (obj *BasicTexObj) SetModel(newModel mgl32.Mat4) {
	obj.model = newModel
}

func (obj *BasicTexObj) Render() {
	gl.UseProgram(obj.program)
	gl.UniformMatrix4fv(obj.uni["project"], 1, false, &(obj.progVar.Vp.Projection[0]))
	gl.UniformMatrix4fv(obj.uni["camera"], 1, false, &(obj.progVar.Vp.Camera[0]))
	gl.UniformMatrix4fv(obj.uni["model"], 1, false, &obj.model[0])
	gl.BindVertexArray(obj.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, obj.texture)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*obj.vertices)/5)) // 6: X,Y,Z,U,V
}

func (obj *BasicTexObj) setTexture() {
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
