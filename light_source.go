package simplegl

import "github.com/go-gl/mathgl/mgl32"

type SimpleLightSrc struct {
	Ambient        float32
	AmbientUniform int32
	Pos            mgl32.Vec3
	PosUniform     int32
	Color          mgl32.Vec3
	ColorUniform   int32
}

func NewLightSrc(
	ambient float32,
	x float32,
	y float32,
	z float32,
	r float32,
	g float32,
	b float32,
) SimpleLightSrc {
	ls := SimpleLightSrc{}
	ls.Ambient = ambient
	ls.Pos = mgl32.Vec3{x, y, z}
	ls.Color = mgl32.Vec3{r, g, b}
	return ls
}
