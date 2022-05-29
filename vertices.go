package sgl

func AddNormal(vertices []float32) []float32 {
	newVertices := []float32{}
	if len(vertices)%9 != 0 {
		return vertices
	}
	// One vertex contains 3 float values, three vertices contain 9 float values
	// Three vertices construct a plane(triangle),
	// therefore we'll get one normal every three points
	for i := 0; i < len(vertices); i += 9 {
		// pt1 : ( vertices[i], vertices[i+1], vertices[i+2] )
		// pt2 : ( vertices[i+3], vertices[i+4], vertices[i+5] )
		// pt3 : ( vertices[i+6], vertices[i+7], vertices[i+8] )
		// vector: origin -> center of the triangle
		center := []float32{
			(vertices[i] + vertices[i+3] + vertices[i+6]) / 3,
			(vertices[i+1] + vertices[i+4] + vertices[i+7]) / 3,
			(vertices[i+2] + vertices[i+5] + vertices[i+8]) / 3,
		}
		// vector: pt1 -> pt2
		vec1 := []float32{
			vertices[i+3] - vertices[i],
			vertices[i+4] - vertices[i+1],
			vertices[i+5] - vertices[i+2],
		}
		// vector: pt2 -> pt3
		vec2 := []float32{
			vertices[i+6] - vertices[i+3],
			vertices[i+7] - vertices[i+4],
			vertices[i+8] - vertices[i+5],
		}
		// normal = vec1 x vec2 (cross product)
		normal := []float32{
			vec1[1]*vec2[2] - vec1[2]*vec2[1],
			vec1[2]*vec2[0] - vec1[0]*vec2[2],
			vec1[0]*vec2[1] - vec1[1]*vec2[0],
		}
		// check if normal . center (dot product) is negative
		dot := normal[0]*center[0] + normal[1]*center[1] + normal[2]*center[2]
		if dot < 0 {
			normal = []float32{
				-normal[0],
				-normal[1],
				-normal[2],
			}
		}
		// newPt1
		newVertices = append(newVertices, vertices[i])
		newVertices = append(newVertices, vertices[i+1])
		newVertices = append(newVertices, vertices[i+2])
		newVertices = append(newVertices, normal[0])
		newVertices = append(newVertices, normal[1])
		newVertices = append(newVertices, normal[2])
		// newPt2
		newVertices = append(newVertices, vertices[i+3])
		newVertices = append(newVertices, vertices[i+4])
		newVertices = append(newVertices, vertices[i+5])
		newVertices = append(newVertices, normal[0])
		newVertices = append(newVertices, normal[1])
		newVertices = append(newVertices, normal[2])
		// newPt3
		newVertices = append(newVertices, vertices[i+6])
		newVertices = append(newVertices, vertices[i+7])
		newVertices = append(newVertices, vertices[i+8])
		newVertices = append(newVertices, normal[0])
		newVertices = append(newVertices, normal[1])
		newVertices = append(newVertices, normal[2])
	}

	// The output vertices will contain 6 * n float values,
	// where n is the number of vertices and 6 is: x position,
	// y position, z position, normal vector x, normal vector y
	// normal vector z
	// i.e. x, y, z, nx, ny, nz
	return newVertices
}
