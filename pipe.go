package main

import (
	"github.com/bloeys/flappy-nmage/quads"
	"github.com/bloeys/gglm/gglm"
)

type Pipe struct {
	quads.Quad
	Col *quads.BoxCollider2D
	//MiddleCol is only defined for top pipes
	MiddleCol *quads.BoxCollider2D

	IsTop bool
}

func NewPipe(isTop bool) Pipe {

	pipeQuad, err := quads.NewQuad("pipe", "./res/textures/pipe-green.png")
	if err != nil {
		panic("Failed to create pipe quad. Err: " + err.Error())
	}

	var middleCol *quads.BoxCollider2D
	if isTop {
		*pipeQuad.RotReadWrite() = *gglm.NewQuatEuler(gglm.NewVec3(0, 0, 180*gglm.Deg2Rad))
		middleCol = quads.NewBoxCollider2D(1, 1)
	}

	return Pipe{
		Quad:      *pipeQuad,
		Col:       quads.NewBoxCollider2D(1, 1),
		MiddleCol: middleCol,
		IsTop:     isTop,
	}
}
