package quads

import (
	"github.com/bloeys/flappy-nmage/entity"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/assets"
	"github.com/bloeys/nmage/meshes"
)

type Quad struct {
	entity.Entity

	Mesh *meshes.Mesh
	Tex  assets.Texture
}

func NewQuad(spriteName, spritePath string) (*Quad, error) {

	quadMesh, err := meshes.NewMesh(spriteName, "./res/models/quad.obj", 0)
	if err != nil {
		return nil, err
	}

	spriteTex, err := assets.LoadPNGTexture(spritePath)
	if err != nil {
		return nil, err
	}

	return &Quad{
		Entity: entity.NewEntity(),
		Mesh:   quadMesh,
		Tex:    spriteTex,
	}, nil
}

type BoxCollider2D struct {
	Quad

	W float32
	H float32
}

func (b *BoxCollider2D) BotLeft() *gglm.Vec2 {
	pos := b.PosRead()
	scale := b.ScaleRead()
	return gglm.NewVec2(pos.X()-b.W*scale.X()*0.5, pos.Y()-b.H*scale.Y()*0.5)
}

func (b *BoxCollider2D) TopRight() *gglm.Vec2 {
	pos := b.PosRead()
	scale := b.ScaleRead()
	return gglm.NewVec2(pos.X()+b.W*scale.X()*0.5, pos.Y()+b.H*scale.Y()*0.5)
}

func NewBoxCollider2D(w, h float32) *BoxCollider2D {

	q, err := NewQuad("boxCollider2D", "./res/textures/white-outline.png")
	if err != nil {
		panic("Failed to create box collider quad. Err: " + err.Error())
	}

	return &BoxCollider2D{
		Quad: *q,
		W:    w,
		H:    h,
	}
}
