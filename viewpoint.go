package simplegl

import "github.com/go-gl/mathgl/mgl32"

type SimpleViewPoint struct {
	Projection        mgl32.Mat4
	ProjectionUniform int32
	Fovy              float32
	Aspect            float32
	Near              float32
	Far               float32
	Camera            mgl32.Mat4
	CameraUniform     int32
	Eye               mgl32.Vec3
	EyePosUniform     int32
	Center            mgl32.Vec3
	Top               mgl32.Vec3
}

func NewViewpoint(width int, height int) SimpleViewPoint {
	vp := SimpleViewPoint{}
	vp.Fovy = mgl32.DegToRad(45.0)
	vp.Aspect = float32(width) / float32(height)
	vp.Near = 0.1
	vp.Far = 200
	vp.Eye = mgl32.Vec3{0, 0, 100}
	vp.Center = mgl32.Vec3{0, 0, 0}
	vp.Top = mgl32.Vec3{0, 1, 0}
	return vp
}
