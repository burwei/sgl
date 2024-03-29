package sgl

import "github.com/go-gl/mathgl/mgl32"

type Viewpoint struct {
	Projection mgl32.Mat4
	Fovy       float32
	Aspect     float32
	Near       float32
	Far        float32
	Camera     mgl32.Mat4
	Eye        mgl32.Vec3
	Target     mgl32.Vec3
	Top        mgl32.Vec3
}

func NewViewpoint(width int, height int) Viewpoint {
	vp := Viewpoint{}
	vp.Fovy = mgl32.DegToRad(45.0)
	vp.Aspect = float32(width) / float32(height)
	vp.Near = 0.1
	vp.Far = 2000
	vp.Eye = mgl32.Vec3{0, 0, 1000}
	vp.Target = mgl32.Vec3{0, 0, 0}
	vp.Top = mgl32.Vec3{0, 1, 0}
	vp.Projection = mgl32.Perspective(vp.Fovy, vp.Aspect, vp.Near, vp.Far)
	vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
	return vp
}
