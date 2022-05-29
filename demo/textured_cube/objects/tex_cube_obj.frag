#version 330

in vec2 fragTexCoord;
out vec4 outputColor;

uniform sampler2D tex;

void main() {
    outputColor = texture(tex, fragTexCoord);
}