package quads

import (
	"github.com/bloeys/flappy-nmage/entity"
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
