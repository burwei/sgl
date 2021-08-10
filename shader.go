package simplegl

import (
	"fmt"
)

const (
	SimpleVShader = `
		#version 330

		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;

		in vec3 vert;

		void main() {
		gl_Position = projection * camera * model * vec4(vert, 1);
		}
	` + "\x00"
	TexVShader = `
		#version 330

		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;

		in vec3 vert;
		in vec2 vertTexCoord;

		out vec2 fragTexCoord;

		void main() {
		fragTexCoord = vertTexCoord;
		gl_Position = projection * camera * model * vec4(vert, 1);
		}
	` + "\x00"
	TexFShader = `
		#version 330

		uniform sampler2D tex;

		in vec2 fragTexCoord;

		out vec4 outputColor;

		void main() {
		outputColor = texture(tex, fragTexCoord);
		}
	` + "\x00"
)

func NewSimpleFShader(r float32, g float32, b float32, a float32) string {
	source1 := `
		#version 330
		out vec4 outputColor;
		void main() {
		outputColor = 
	`
	source2 := fmt.Sprintf(
		"vec4(%.3f, %.3f, %.3f, %.3f);}",
		r,
		g,
		b,
		a,
	)
	return source1 + source2
}
