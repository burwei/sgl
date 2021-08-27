package simplegl

// NewSimpleCube will return a cube with side length l.
func NewSimpleCube(l float32) *[]float32{
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