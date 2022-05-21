package entity

import "github.com/bloeys/gglm/gglm"

type Entity struct {
	pos   *gglm.Vec3
	scale *gglm.Vec3
	rot   *gglm.Quat

	IsDirty  bool
	modelMat *gglm.TrMat
}

func (e *Entity) PosReadWrite() *gglm.Vec3 {
	e.IsDirty = true
	return e.pos
}

func (e *Entity) PosRead() *gglm.Vec3 {
	return e.pos.Clone()
}

func (e *Entity) ScaleReadWrite() *gglm.Vec3 {
	e.IsDirty = true
	return e.scale
}

func (e *Entity) ScaleRead() *gglm.Vec3 {
	return e.scale.Clone()
}

func (e *Entity) RotReadWrite() *gglm.Quat {
	e.IsDirty = true
	return e.rot
}

func (e *Entity) RotRead() *gglm.Quat {
	var q gglm.Quat = *e.rot
	return &q
}

func (e *Entity) ModelMat() *gglm.TrMat {

	if !e.IsDirty {
		return e.modelMat
	}

	e.IsDirty = false
	e.modelMat = gglm.NewTrMatId().Scale(e.scale).Translate(e.pos).Rotate(e.rot.Angle(), e.rot.Axis())

	return e.modelMat
}

func NewEntity() Entity {
	return Entity{
		pos:      gglm.NewVec3(0, 0, 0),
		scale:    gglm.NewVec3(1, 1, 1),
		rot:      gglm.NewQuatId(),
		modelMat: gglm.NewTrMatId(),
	}
}
