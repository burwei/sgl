package sgl

import "github.com/go-gl/mathgl/mgl32"

type LightSrc struct {
	Pos       mgl32.Vec3
	Color     mgl32.Vec3
	Intensity float32
}

func NewLightSrc() LightSrc {
	ls := LightSrc{}
	ls.Pos = mgl32.Vec3{-500, 0, 1000}
	ls.Color = mgl32.Vec3{1, 1, 1}
	ls.Intensity = 1
	return ls
}
