#version 330

layout(location = 0) in vec3 vert;
layout(location = 1) in vec2 vertTexCoord;

out vec2 fragTexCoord;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}