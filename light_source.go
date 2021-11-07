package simplegl

import "github.com/go-gl/mathgl/mgl32"

type SimpleLightSrc struct {
	AmbientStrength         float32
	AmbientStrengthUniform  int32
	SpecularStrength        float32
	SpecularStrengthUniform int32
	Shininess               float32
	ShininessUniform        int32
	Pos                     mgl32.Vec3
	PosUniform              int32
	Color                   mgl32.Vec3
	ColorUniform            int32
}

func NewLightSrc() SimpleLightSrc {
	ls := SimpleLightSrc{}
	ls.AmbientStrength = 0.3
	ls.SpecularStrength = 1
	ls.Shininess = 8
	ls.Pos = mgl32.Vec3{-500, 0, 1000}
	ls.Color = mgl32.Vec3{1, 1, 1}
	return ls
}
