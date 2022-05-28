package sgl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Group struct {
	objects      map[string]Object
	objectModels map[string]mgl32.Mat4
	groupModel   mgl32.Mat4
}

func NewGroup() Group {
	g := Group{}
	g.objects = map[string]Object{}
	g.objectModels = map[string]mgl32.Mat4{}
	g.groupModel = mgl32.Ident4()
	return g
}

func (g *Group) AddObject(name string, obj Object) {
	g.objects[name] = obj
	g.objectModels[name] = obj.GetModel()
}

func (g *Group) SetObjectModel(name string, newModel mgl32.Mat4) {
	g.objectModels[name] = newModel
}

func (g *Group) SetGroupModel(newModel mgl32.Mat4) {
	g.groupModel = newModel
}

func (g *Group) Render() {
	for name, obj := range g.objects {
		model := g.groupModel.Mul4(g.objectModels[name])
		obj.SetModel(model)
		obj.Render()
	}
}
