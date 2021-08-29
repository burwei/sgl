package simplegl

// NewSimpleCube will return vertices of a cube with side length l.
func NewSimpleCube(l float32) *[]float32{
	return &[]float32{
		//  X, Y, Z, U, V
		// Bottom (-Y)
		-l, -l, -l,
		l, -l, -l,
		-l, -l, l,
		l, -l, -l,
		l, -l, l,
		-l, -l, l,
	
		// Top (+Y)
		-l, l, -l,
		-l, l, l,
		l, l, -l,
		l, l, -l,
		-l, l, l,
		l, l, l,
	
		// Front (+Z)
		-l, -l, l,
		l, -l, l,
		-l, l, l,
		l, -l, l,
		l, l, l,
		-l, l, l,
	
		// Back (-Z)
		-l, -l, -l,
		-l, l, -l,
		l, -l, -l,
		l, -l, -l,
		-l, l, -l,
		l, l, -l,
	
		// Left (-X)
		-l, -l, l,
		-l, l, -l,
		-l, -l, -l,
		-l, -l, l,
		-l, l, l,
		-l, l, -l,
	
		// Right (+X)
		l, -l, l,
		l, -l, -l,
		l, l, -l,
		l, -l, l,
		l, l, -l,
		l, l, l,
	}
}

// NewUniTexCube will return vertices of a cube with side length l.
// The vertices contains custom data (texture vector)
func NewUniTexCube(l float32) *[]float32{
	return &[]float32{
		//  X, Y, Z, U, V
		// Bottom (-Y)
		-l, -l, -l, 0.0, 0.0,
		l, -l, -l, 1.0, 0.0,
		-l, -l, l, 0.0, 1.0,
		l, -l, -l, 1.0, 0.0,
		l, -l, l, 1.0, 1.0,
		-l, -l, l, 0.0, 1.0,
	
		// Top (+Y)
		-l, l, -l, 0.0, 0.0,
		-l, l, l, 0.0, 1.0,
		l, l, -l, 1.0, 0.0,
		l, l, -l, 1.0, 0.0,
		-l, l, l, 0.0, 1.0,
		l, l, l, 1.0, 1.0,
	
		// Front (+Z)
		-l, -l, l, 1.0, 0.0,
		l, -l, l, 0.0, 0.0,
		-l, l, l, 1.0, 1.0,
		l, -l, l, 0.0, 0.0,
		l, l, l, 0.0, 1.0,
		-l, l, l, 1.0, 1.0,
	
		// Back (-Z)
		-l, -l, -l, 0.0, 0.0,
		-l, l, -l, 0.0, 1.0,
		l, -l, -l, 1.0, 0.0,
		l, -l, -l, 1.0, 0.0,
		-l, l, -l, 0.0, 1.0,
		l, l, -l, 1.0, 1.0,
	
		// Left (-X)
		-l, -l, l, 0.0, 1.0,
		-l, l, -l, 1.0, 0.0,
		-l, -l, -l, 0.0, 0.0,
		-l, -l, l, 0.0, 1.0,
		-l, l, l, 1.0, 1.0,
		-l, l, -l, 1.0, 0.0,
	
		// Right (+X)
		l, -l, l, 1.0, 1.0,
		l, -l, -l, 1.0, 0.0,
		l, l, -l, 0.0, 0.0,
		l, -l, l, 1.0, 1.0,
		l, l, -l, 0.0, 0.0,
		l, l, l, 0.0, 1.0,
	}
}