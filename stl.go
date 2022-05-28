package sgl

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

type head struct {
	Header [80]byte
	TriNum uint32
}

type triangle struct {
	Normal [3]float32
	Vert1  [3]float32
	Vert2  [3]float32
	Vert3  [3]float32
	Count  uint16
}

func ReadBinaryStlFile(file string) []float32 {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	h := head{}
	err = binary.Read(f, binary.LittleEndian, &h)
	if err != nil {
		return nil
	}
	triangles := make([]triangle, h.TriNum)
	err = binary.Read(f, binary.LittleEndian, triangles)
	if err != nil {
		return nil
	}
	vertices := []float32{}
	min := -1 * math.MaxFloat64
	max := math.MaxFloat64
	maxX := min
	maxY := min
	maxZ := min
	minX := max
	minY := max
	minZ := max
	for i := 0; i < int(h.TriNum); i++ {
		// get max
		maxX = math.Max(maxX, float64(triangles[i].Vert1[0]))
		maxX = math.Max(maxX, float64(triangles[i].Vert2[0]))
		maxX = math.Max(maxX, float64(triangles[i].Vert3[0]))
		maxY = math.Max(maxY, float64(triangles[i].Vert1[1]))
		maxY = math.Max(maxY, float64(triangles[i].Vert2[1]))
		maxY = math.Max(maxY, float64(triangles[i].Vert3[1]))
		maxZ = math.Max(maxZ, float64(triangles[i].Vert1[2]))
		maxZ = math.Max(maxZ, float64(triangles[i].Vert2[2]))
		maxZ = math.Max(maxZ, float64(triangles[i].Vert3[2]))
		// get min
		minX = math.Min(minX, float64(triangles[i].Vert1[0]))
		minX = math.Min(minX, float64(triangles[i].Vert2[0]))
		minX = math.Min(minX, float64(triangles[i].Vert3[0]))
		minY = math.Min(minY, float64(triangles[i].Vert1[1]))
		minY = math.Min(minY, float64(triangles[i].Vert2[1]))
		minY = math.Min(minY, float64(triangles[i].Vert3[1]))
		minZ = math.Min(minZ, float64(triangles[i].Vert1[2]))
		minZ = math.Min(minZ, float64(triangles[i].Vert2[2]))
		minZ = math.Min(minZ, float64(triangles[i].Vert3[2]))
	}
	centerX := float32((minX + maxX) / 2)
	centerY := float32((minY + maxY) / 2)
	centerZ := float32((minZ + maxZ) / 2)
	for i := 0; i < int(h.TriNum); i++ {
		// vertex 1
		vertices = append(vertices, triangles[i].Vert1[0]-centerX)
		vertices = append(vertices, triangles[i].Vert1[1]-centerY)
		vertices = append(vertices, triangles[i].Vert1[2]-centerZ)
		// vertex 2
		vertices = append(vertices, triangles[i].Vert2[0]-centerX)
		vertices = append(vertices, triangles[i].Vert2[1]-centerY)
		vertices = append(vertices, triangles[i].Vert2[2]-centerZ)
		// vertex 3
		vertices = append(vertices, triangles[i].Vert3[0]-centerX)
		vertices = append(vertices, triangles[i].Vert3[1]-centerY)
		vertices = append(vertices, triangles[i].Vert3[2]-centerZ)
		// fmt.Printf("vert1: ( %.2f, %.2f, %.2f )\n",triangles[i].Vert1[0],triangles[i].Vert1[1],triangles[i].Vert1[2])
		// fmt.Printf("vert2: ( %.2f, %.2f, %.2f )\n",triangles[i].Vert2[0],triangles[i].Vert2[1],triangles[i].Vert2[2])
		// fmt.Printf("vert3: ( %.2f, %.2f, %.2f )\n",triangles[i].Vert3[0],triangles[i].Vert3[1],triangles[i].Vert3[2])
		// fmt.Printf("normal: ( %.2f, %.2f, %.2f )\n",triangles[i].Normal[0],triangles[i].Normal[1],triangles[i].Normal[2])
	}
	return vertices
}

func ReadBinaryStlFileRaw(file string) []float32 {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	h := head{}
	err = binary.Read(f, binary.LittleEndian, &h)
	if err != nil {
		return nil
	}
	triangles := make([]triangle, h.TriNum)
	err = binary.Read(f, binary.LittleEndian, triangles)
	if err != nil {
		return nil
	}
	vertices := []float32{}
	for i := 0; i < int(h.TriNum); i++ {
		// vertex 1
		vertices = append(vertices, triangles[i].Vert1[0])
		vertices = append(vertices, triangles[i].Vert1[1])
		vertices = append(vertices, triangles[i].Vert1[2])
		// vertex 2
		vertices = append(vertices, triangles[i].Vert2[0])
		vertices = append(vertices, triangles[i].Vert2[1])
		vertices = append(vertices, triangles[i].Vert2[2])
		// vertex 3
		vertices = append(vertices, triangles[i].Vert3[0])
		vertices = append(vertices, triangles[i].Vert3[1])
		vertices = append(vertices, triangles[i].Vert3[2])
	}
	return vertices
}

func ReadBinaryStlFileWithCenter(
	file string,
	centerX float32,
	centerY float32,
	centerZ float32,
) []float32 {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	h := head{}
	err = binary.Read(f, binary.LittleEndian, &h)
	if err != nil {
		return nil
	}
	triangles := make([]triangle, h.TriNum)
	err = binary.Read(f, binary.LittleEndian, triangles)
	if err != nil {
		return nil
	}
	vertices := []float32{}
	for i := 0; i < int(h.TriNum); i++ {
		// vertex 1
		vertices = append(vertices, triangles[i].Vert1[0]-centerX)
		vertices = append(vertices, triangles[i].Vert1[1]-centerY)
		vertices = append(vertices, triangles[i].Vert1[2]-centerZ)
		// vertex 2
		vertices = append(vertices, triangles[i].Vert2[0]-centerX)
		vertices = append(vertices, triangles[i].Vert2[1]-centerY)
		vertices = append(vertices, triangles[i].Vert2[2]-centerZ)
		// vertex 3
		vertices = append(vertices, triangles[i].Vert3[0]-centerX)
		vertices = append(vertices, triangles[i].Vert3[1]-centerY)
		vertices = append(vertices, triangles[i].Vert3[2]-centerZ)
	}
	return vertices
}
