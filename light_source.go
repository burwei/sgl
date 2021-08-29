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

func NewLightSrc(
	ambientStrength float32,
	specularStrength float32,
	shininess float32,
	x float32,
	y float32,
	z float32,
	r float32,
	g float32,
	b float32,
) SimpleLightSrc {
	ls := SimpleLightSrc{}
	ls.AmbientStrength = ambientStrength
	ls.SpecularStrength = specularStrength
	ls.Shininess = shininess
	ls.Pos = mgl32.Vec3{x, y, z}
	ls.Color = mgl32.Vec3{r, g, b}
	return ls
}
