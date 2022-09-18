package sgl

import "github.com/go-gl/mathgl/mgl32"

type Material struct {
	Ambient   mgl32.Vec3
	Diffuse   mgl32.Vec3
	Specular  mgl32.Vec3
	Shininess float32
}

func NewMaterial() *Material {
	m := Material{}

	m.Ambient = mgl32.Vec3{0.3, 0.3, 0.3}
	m.Diffuse = mgl32.Vec3{1, 1, 1}
	m.Specular = mgl32.Vec3{1, 1, 1}
	m.Shininess = 8

	return &m
}
