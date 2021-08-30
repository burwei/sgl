package simplegl

import (
	"fmt"
)

// This file contains several simple GLSL code
// to provides some simple vertex shaders and fragment shaders.

func NewBasicVShader() string {
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

func NewBasicFShader(r float32, g float32, b float32) string {
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

func NewBasicTexVShader() string {
	return fmt.Sprintf(
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
}

func NewBasicTexFShader() string {
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

func NewBasicLightVShader() string {
	// Todo: avoid inverse calculation
	return fmt.Sprintf(
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
}

func NewBasicLightFShader(r float32, g float32, b float32) string {
	return fmt.Sprintf(
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
}
