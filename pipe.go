package main

import "github.com/bloeys/gglm/gglm"

type Transform struct {
	Pos   *gglm.Vec3
	Scale *gglm.Vec3
	Rot   *gglm.Quat

	IsDirty bool
	TrMat   *gglm.TrMat
}

type Pipe struct {
	ShouldRegen bool
	Pos         *gglm.Vec3
	OrigPos     *gglm.Vec3
	TrMat       *gglm.TrMat
}

func NewPipe(isTop bool) Pipe {

	var q *gglm.Quat
	if isTop {
		q = gglm.NewQuatEuler(gglm.NewVec3(0, 90*gglm.Deg2Rad, 180*gglm.Deg2Rad))
	} else {
		q = gglm.NewQuatEuler(gglm.NewVec3(0, 90*gglm.Deg2Rad, 0))
	}

	x := gglm.NewTrMatId()
	x.Rotate(q.Angle(), q.Axis())
	return Pipe{
		ShouldRegen: true,
		Pos:         &gglm.Vec3{},
		OrigPos:     &gglm.Vec3{},
		TrMat:       x,
	}
}
