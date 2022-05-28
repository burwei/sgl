package sgl

func NewPlane(l float32) *[]float32 {
	return &[]float32{
		//  X, Y, Z
		-l / 2, 0, -l / 2,
		-l / 2, 0, l / 2,
		l / 2, 0, -l / 2,
		l / 2, 0, -l / 2,
		-l / 2, 0, l / 2,
		l / 2, 0, l / 2,
	}
}

// NewSimpleCube will return vertices of a cube with side length l.
func NewCube(l float32) *[]float32 {
	return &[]float32{
		//  X, Y, Z
		// Bottom (-Y)
		-l / 2, -l / 2, -l / 2,
		l / 2, -l / 2, -l / 2,
		-l / 2, -l / 2, l / 2,
		l / 2, -l / 2, -l / 2,
		l / 2, -l / 2, l / 2,
		-l / 2, -l / 2, l / 2,

		// Top (+Y)
		-l / 2, l / 2, -l / 2,
		-l / 2, l / 2, l / 2,
		l / 2, l / 2, -l / 2,
		l / 2, l / 2, -l / 2,
		-l / 2, l / 2, l / 2,
		l / 2, l / 2, l / 2,

		// Front (+Z)
		-l / 2, -l / 2, l / 2,
		l / 2, -l / 2, l / 2,
		-l / 2, l / 2, l / 2,
		l / 2, -l / 2, l / 2,
		l / 2, l / 2, l / 2,
		-l / 2, l / 2, l / 2,

		// Back (-Z)
		-l / 2, -l / 2, -l / 2,
		-l / 2, l / 2, -l / 2,
		l / 2, -l / 2, -l / 2,
		l / 2, -l / 2, -l / 2,
		-l / 2, l / 2, -l / 2,
		l / 2, l / 2, -l / 2,

		// l/2eft (-X)
		-l / 2, -l / 2, l / 2,
		-l / 2, l / 2, -l / 2,
		-l / 2, -l / 2, -l / 2,
		-l / 2, -l / 2, l / 2,
		-l / 2, l / 2, l / 2,
		-l / 2, l / 2, -l / 2,

		// Right (+X)
		l / 2, -l / 2, l / 2,
		l / 2, -l / 2, -l / 2,
		l / 2, l / 2, -l / 2,
		l / 2, -l / 2, l / 2,
		l / 2, l / 2, -l / 2,
		l / 2, l / 2, l / 2,
	}
}

// NewUniTexCube will return vertices of a cube with side length l.
// The vertices contains custom data (texture vector)
func NewUniTexCube(l float32) *[]float32 {
	return &[]float32{
		//  X, Y, Z, U, V
		// Bottom (-Y)
		-l / 2, -l / 2, -l / 2, 0.0, 0.0,
		l / 2, -l / 2, -l / 2, 1.0, 0.0,
		-l / 2, -l / 2, l / 2, 0.0, 1.0,
		l / 2, -l / 2, -l / 2, 1.0, 0.0,
		l / 2, -l / 2, l / 2, 1.0, 1.0,
		-l / 2, -l / 2, l / 2, 0.0, 1.0,

		// Top (+Y)
		-l / 2, l / 2, -l / 2, 0.0, 0.0,
		-l / 2, l / 2, l / 2, 0.0, 1.0,
		l / 2, l / 2, -l / 2, 1.0, 0.0,
		l / 2, l / 2, -l / 2, 1.0, 0.0,
		-l / 2, l / 2, l / 2, 0.0, 1.0,
		l / 2, l / 2, l / 2, 1.0, 1.0,

		// Front (+Z)
		-l / 2, -l / 2, l / 2, 1.0, 0.0,
		l / 2, -l / 2, l / 2, 0.0, 0.0,
		-l / 2, l / 2, l / 2, 1.0, 1.0,
		l / 2, -l / 2, l / 2, 0.0, 0.0,
		l / 2, l / 2, l / 2, 0.0, 1.0,
		-l / 2, l / 2, l / 2, 1.0, 1.0,

		// Back (-Z)
		-l / 2, -l / 2, -l / 2, 0.0, 0.0,
		-l / 2, l / 2, -l / 2, 0.0, 1.0,
		l / 2, -l / 2, -l / 2, 1.0, 0.0,
		l / 2, -l / 2, -l / 2, 1.0, 0.0,
		-l / 2, l / 2, -l / 2, 0.0, 1.0,
		l / 2, l / 2, -l / 2, 1.0, 1.0,

		// l/2eft (-X)
		-l / 2, -l / 2, l / 2, 0.0, 1.0,
		-l / 2, l / 2, -l / 2, 1.0, 0.0,
		-l / 2, -l / 2, -l / 2, 0.0, 0.0,
		-l / 2, -l / 2, l / 2, 0.0, 1.0,
		-l / 2, l / 2, l / 2, 1.0, 1.0,
		-l / 2, l / 2, -l / 2, 1.0, 0.0,

		// Right (+X)
		l / 2, -l / 2, l / 2, 1.0, 1.0,
		l / 2, -l / 2, -l / 2, 1.0, 0.0,
		l / 2, l / 2, -l / 2, 0.0, 0.0,
		l / 2, -l / 2, l / 2, 1.0, 1.0,
		l / 2, l / 2, -l / 2, 0.0, 0.0,
		l / 2, l / 2, l / 2, 0.0, 1.0,
	}
}
