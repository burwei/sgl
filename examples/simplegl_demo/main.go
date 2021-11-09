package main

import (
	"fmt"
	"math"

	sgl "github.com/burwei/simplegl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
	title  = "SimpleGL"
)

func main() {
	window := sgl.Init(width, height, title)
	defer sgl.Terminate()

	vp := sgl.NewViewpoint(width, height)
	ls := sgl.NewLightSrc()
	mt := sgl.Material{
		Ambient: mgl32.Vec3{0.1, 0.1, 0.1},
		Diffuse: mgl32.Vec3{0.6, 0.6, 0.6},
		Specular: mgl32.Vec3{1.5, 1.5, 1.5},
		Shininess: 24,
	}

	emptyCube := sgl.BasicObj{}
	emptyCube.PrepareProgram(false)

	pos := [][]float32{
		// S
		{-300, 200, 0},
		{-320, 200, 0},
		{-340, 200, 0},
		{-360, 200, 0},
		{-380, 200, 0},
		{-380, 180, 0},
		{-380, 160, 0},
		{-380, 140, 0},
		{-380, 120, 0},
		{-380, 100, 0},
		{-360, 100, 0},
		{-340, 100, 0},
		{-320, 100, 0},
		{-300, 100, 0},
		{-300, 80, 0},
		{-300, 60, 0},
		{-300, 40, 0},
		{-300, 20, 0},
		{-300, 0, 0},
		{-320, 0, 0},
		{-340, 0, 0},
		{-360, 0, 0},
		{-380, 0, 0},
		// i
		{-260, 140, 0},
		{-260, 100, 0},
		{-260, 80, 0},
		{-260, 60, 0},
		{-260, 40, 0},
		{-260, 20, 0},
		{-260, 0, 0},
		// m
		{-220, 100, 0},
		{-220, 80, 0},
		{-220, 60, 0},
		{-220, 40, 0},
		{-220, 20, 0},
		{-220, 0, 0},
		{-220, 100, 0},
		{-200, 100, 0},
		{-180, 100, 0},
		{-180, 100, 0},
		{-180, 80, 0},
		{-180, 60, 0},
		{-180, 40, 0},
		{-180, 20, 0},
		{-180, 0, 0},
		{-310, 0, 0},
		{-160, 100, 0},
		{-140, 100, 0},
		{-140, 80, 0},
		{-140, 60, 0},
		{-140, 40, 0},
		{-140, 20, 0},
		{-140, 0, 0},
		// p
		{-100, 100, 0},
		{-100, 80, 0},
		{-100, 60, 0},
		{-100, 40, 0},
		{-100, 20, 0},
		{-100, 0, 0},
		{-100, -20, 0},
		{-100, -40, 0},
		{-100, -60, 0},
		{-100, 100, 0},
		{-80, 100, 0},
		{-60, 100, 0},
		{-40, 100, 0},
		{-20, 100, 0},
		{-20, 80, 0},
		{-20, 60, 0},
		{-20, 40, 0},
		{-20, 20, 0},
		{-20, 0, 0},
		{-40, 0, 0},
		{-60, 0, 0},
		{-80, 0, 0},
		{-100, 0, 0},
		// l
		{20, 200, 0},
		{20, 180, 0},
		{20, 160, 0},
		{20, 140, 0},
		{20, 120, 0},
		{20, 100, 0},
		{20, 80, 0},
		{20, 60, 0},
		{20, 40, 0},
		{20, 20, 0},
		{20, 0, 0},
		// e
		{60, 100, 0},
		{60, 80, 0},
		{60, 60, 0},
		{60, 40, 0},
		{60, 20, 0},
		{60, 0, 0},
		{80, 100, 0},
		{100, 100, 0},
		{120, 100, 0},
		{140, 100, 0},
		{140, 80, 0},
		{140, 60, 0},
		{140, 50, 0},
		{120, 50, 0},
		{100, 50, 0},
		{80, 50, 0},
		{140, 0, 0},
		{120, 0, 0},
		{100, 0, 0},
		{80, 0, 0},
		// G
		{260, 200, 0},
		{240, 200, 0},
		{220, 200, 0},
		{200, 200, 0},
		{180, 200, 0},
		{180, 180, 0},
		{180, 160, 0},
		{180, 140, 0},
		{180, 120, 0},
		{180, 100, 0},
		{180, 80, 0},
		{180, 60, 0},
		{180, 40, 0},
		{180, 20, 0},
		{180, 0, 0},
		{200, 0, 0},
		{220, 0, 0},
		{240, 0, 0},
		{260, 0, 0},
		{260, 20, 0},
		{260, 40, 0},
		{260, 60, 0},
		{260, 80, 0},
		{260, 100, 0},
		{240, 100, 0},
		{230, 100, 0},
		// L
		{300, 200, 0},
		{300, 180, 0},
		{300, 160, 0},
		{300, 140, 0},
		{300, 120, 0},
		{300, 100, 0},
		{300, 80, 0},
		{300, 60, 0},
		{300, 40, 0},
		{300, 20, 0},
		{300, 0, 0},
		{320, 0, 0},
		{340, 0, 0},
		{360, 0, 0},
		{380, 0, 0},
	}
	group := sgl.NewGroup()
	for i, v := range pos {
		newCube := sgl.BasicObj{}
		newCube.SetProgramVar(sgl.BasicObjProgVar{
			Red:   1,
			Green: 0.3,
			Blue:  0.3,
			Vp:    &vp,
			Ls:    &ls,
			Mt:    &mt,
		})
		newCube.BindProgramVar(emptyCube.GetProgram())
		newCube.SetVertices(sgl.NewCube(20))
		newCube.SetModel(
			mgl32.Translate3D(v[0], v[1], v[2]),
		)
		group.AddObject(fmt.Sprintf("cube%v", i), &newCube)
	}

	angle := 0.0
	previousTime := glfw.GetTime()
	speedConst := 0.7

	sgl.BeforeMainLoop(window, &vp)
	for !window.ShouldClose() {
		sgl.BeforeDrawing()

		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed * speedConst
		if math.Sin(angle) < 0 {
			speedConst += 0.15
		} else {
			speedConst = 0.7
		}

		group.SetGroupModel(
			mgl32.Translate3D(0, -80, 180).Mul4(
				mgl32.Rotate3DY(float32(angle) - math.Pi/8).Mat4(),
			),
		)
		group.Render()

		sgl.AfterDrawing(window)
	}
}
