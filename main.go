package main

import (
	"math/rand"
	"time"

	"github.com/bloeys/flappy-nmage/quads"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/engine"
	"github.com/bloeys/nmage/input"
	"github.com/bloeys/nmage/materials"
	"github.com/bloeys/nmage/renderer/rend3dgl"
	"github.com/bloeys/nmage/timing"
	nmageimgui "github.com/bloeys/nmage/ui/imgui"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	pipeXSpacing float32 = 6
	pipeYSpacing float32 = 6
)

var (
	simpleMat *materials.Material

	birdSprite       *quads.Quad
	backgroundSprite *quads.Quad

	gravity float32 = -9.81 * 4

	pipes      []Pipe
	pipesSpeed = gglm.NewVec3(-5, 0, 0)

	birdVelocity         = gglm.NewVec3(0, 0, 0)
	dragAmount   float32 = 0.9
	jumpForce    float32 = 30
)

func main() {

	rand.Seed(time.Now().Unix())

	//Init engine
	err := engine.Init()
	if err != nil {
		panic("Failed to init nMage. Err:" + err.Error())
	}

	//Create window
	window, err := engine.CreateOpenGLWindowCentered("nMage", 1280, 720, engine.WindowFlags_RESIZABLE, rend3dgl.NewRend3DGL())
	if err != nil {
		panic("Failed to create window. Err: " + err.Error())
	}
	defer window.Destroy()

	engine.SetVSync(true)

	window.SDLWin.SetTitle("Flappy Bird")

	g := &Game{
		shouldRun: true,
		win:       window,
		imguiInfo: nmageimgui.NewImGUI(),
	}

	engine.Run(g)
}

var _ engine.Game = &Game{}

type Game struct {
	shouldRun bool
	win       *engine.Window
	imguiInfo nmageimgui.ImguiInfo
}

func (g *Game) Init() {

	var err error

	//Load assets
	birdSprite, err = quads.NewQuad("bird", "./res/textures/yellowbird-midflap.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}

	backgroundSprite, err = quads.NewQuad("background", "./res/textures/background-day.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}

	simpleMat = materials.NewMaterial("simpleMat", "./res/shaders/simple")

	//Setup camera
	camPos := gglm.NewVec3(0, 0, -10)
	// camForward := gglm.NewVec3(0, 0, 1)
	targetPos := gglm.NewVec3(0, 0, 0)

	viewMat := gglm.LookAt(camPos, targetPos, gglm.NewVec3(0, 1, 0))
	simpleMat.SetUnifMat4("viewMat", &viewMat.Mat4)

	projMat := gglm.Ortho(-10, 10, 10, -10, 0.1, 500)
	simpleMat.SetUnifMat4("projMat", &projMat.Mat4)

	//Set positions
	birdSprite.ScaleReadWrite().Set(0.75, 1, 1)
	birdSprite.PosReadWrite().Set(-5, 0, 2)
}

func (g *Game) Start() {

	backgroundSprite.Entity.ScaleReadWrite().Set(10, 10, 1)

	//Pipes
	pos := gglm.NewVec3(12, 0, 0)
	for i := 0; i < 10; i++ {

		randY := rand.Float32() * 5
		if rand.Float32() > 0.5 {
			randY *= -1
		}

		p := NewPipe(true)
		p.PosReadWrite().Add(gglm.NewVec3(pos.X(), 10+randY+pipeYSpacing*0.5, 1))
		p.ScaleReadWrite().Set(1, 10, 1)
		pipes = append(pipes, p)

		p = NewPipe(false)
		p.PosReadWrite().Add(gglm.NewVec3(pos.X(), -10+randY-pipeYSpacing*0.5, 1))
		p.ScaleReadWrite().Set(1, 10, 1)
		pipes = append(pipes, p)

		pos.SetX(pos.X() + pipeXSpacing)
	}
}

func (g *Game) FrameStart() {

}

func (g *Game) Update() {

	if input.IsQuitClicked() || input.KeyClicked(sdl.K_ESCAPE) {
		g.shouldRun = false
	}

	//Move the bird
	if birdVelocity.Y() > 0 {
		birdVelocity.SetY(birdVelocity.Y() * dragAmount)
	}

	birdVelocity.SetY(birdVelocity.Y() + gravity*timing.DT())
	if input.KeyClicked(sdl.K_SPACE) {
		birdVelocity.SetY(jumpForce)
	}

	positionDelta := *birdVelocity
	positionDelta.Scale(timing.DT())
	birdSprite.PosReadWrite().Add(&positionDelta)

	//Move the pipes
	spd := *pipesSpeed
	spd.Scale(timing.DT())
	for i := 0; i < len(pipes); i++ {
		pipes[i].PosReadWrite().Add(&spd)
	}

	//Find the pipe with largest X pos
	var largestXPosPipe *Pipe = &pipes[0]
	for i := 1; i < len(pipes); i++ {
		if pipes[i].PosRead().X() > largestXPosPipe.PosRead().X() {
			largestXPosPipe = &pipes[i]
		}
	}

	//Pipe logic
	for i := 0; i < len(pipes); i += 2 {

		xPos := pipes[i].PosRead().X()
		if xPos > -12 {
			continue
		}

		randY := rand.Float32() * 5
		if rand.Float32() > 0.5 {
			randY *= -1
		}

		upperPipePos := pipes[i].PosReadWrite()
		lowerPipePos := pipes[i+1].PosReadWrite()

		newX := largestXPosPipe.PosRead().X() + pipeXSpacing
		upperPipePos.SetX(newX)
		lowerPipePos.SetX(newX)

		upperPipePos.SetY(10 + randY + pipeYSpacing*0.5)
		lowerPipePos.SetY(-10 + randY - pipeYSpacing*0.5)
	}
}

func (g *Game) Render() {

	//Draw background
	simpleMat.DiffuseTex = backgroundSprite.Tex.TexID
	simpleMat.SetAttribute(backgroundSprite.Mesh.Buf)
	simpleMat.Bind()
	g.win.Rend.Draw(backgroundSprite.Mesh, backgroundSprite.ModelMat(), simpleMat)

	//Draw pipe
	simpleMat.DiffuseTex = pipes[0].Tex.TexID
	simpleMat.SetAttribute(pipes[0].Mesh.Buf)
	simpleMat.Bind()
	for i := 0; i < len(pipes); i++ {
		g.win.Rend.Draw(pipes[0].Mesh, pipes[i].ModelMat(), simpleMat)
	}

	//Draw bird
	simpleMat.DiffuseTex = birdSprite.Tex.TexID
	simpleMat.SetAttribute(birdSprite.Mesh.Buf)
	simpleMat.Bind()
	g.win.Rend.Draw(birdSprite.Mesh, birdSprite.ModelMat(), simpleMat)
}

func (g *Game) FrameEnd() {

}

func (g *Game) ShouldRun() bool {
	return g.shouldRun
}

func (g *Game) GetWindow() *engine.Window {
	return g.win
}

func (g *Game) GetImGUI() nmageimgui.ImguiInfo {
	return g.imguiInfo
}

func (g *Game) Deinit() {

}
