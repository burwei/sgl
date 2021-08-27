package simplegl

import (
	"fmt"
)

func NewSimpleVShader() string {
	return fmt.Sprintf(
		`
		#version 330

		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;

		in vec3 vert;

		void main() {
		gl_Position = projection * camera * model * vec4(vert, 1);
		}
		%v`,
		"\x00",
	)
}

func NewSimpleFShader(r float32, g float32, b float32) string {
	return fmt.Sprintf(
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
}

func NewTexVShader() string {
	return fmt.Sprintf(
		`
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
		%v`,
		"\x00",
	)
}

func NewTexFShader() string{
	return fmt.Sprintf(
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
}
