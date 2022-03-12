package main

import (
	"math/rand"

	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/assets"
	"github.com/bloeys/nmage/engine"
	"github.com/bloeys/nmage/input"
	"github.com/bloeys/nmage/materials"
	"github.com/bloeys/nmage/meshes"
	"github.com/bloeys/nmage/renderer/rend3dgl"
	"github.com/bloeys/nmage/timing"
	nmageimgui "github.com/bloeys/nmage/ui/imgui"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	simpleMat *materials.Material

	birdMesh  *meshes.Mesh
	birdTex   assets.Texture
	birdTrMat = gglm.NewTrMatId()

	pipeMesh *meshes.Mesh
	pipeTex  assets.Texture

	gravity float32 = -9.81 * 4

	pipes      []Pipe
	pipesSpeed = gglm.NewVec3(-5, 0, 0)

	birdVelocity         = gglm.NewVec3(0, 0, 0)
	dragAmount   float32 = 0.9
	jumpForce    float32 = 30
)

func main() {

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
	birdMesh, err = meshes.NewMesh("bird", "./res/models/bird.fbx", 0)
	if err != nil {
		panic("Failed to load mesh. Err:" + err.Error())
	}

	birdTex, err = assets.LoadPNGTexture("./res/textures/bird.png")
	if err != nil {
		panic("Failed to load texture. Err:" + err.Error())
	}

	pipeMesh, err = meshes.NewMesh("pipe", "./res/models/pipe.fbx", 0)
	if err != nil {
		panic("Failed to load mesh. Err:" + err.Error())
	}

	pipeTex, err = assets.LoadPNGTexture("./res/textures/pipe.png")
	if err != nil {
		panic("Failed to load texture. Err:" + err.Error())
	}

	simpleMat = materials.NewMaterial("simpleMat", "./res/shaders/simple")

	//Setup camera
	camPos := gglm.NewVec3(0, 0, -10)
	// camForward := gglm.NewVec3(0, 0, 1)
	targetPos := gglm.NewVec3(0, 0, 0)

	viewMat := gglm.LookAt(camPos, targetPos, gglm.NewVec3(0, 1, 0))
	simpleMat.SetUnifMat4("viewMat", &viewMat.Mat4)

	projMat := gglm.Ortho(-10, 10, 10, -10, 0.1, 500)
	// projMat := gglm.Perspective(45*gglm.Deg2Rad, float32(1280)/float32(720), 0.1, 500)
	simpleMat.SetUnifMat4("projMat", &projMat.Mat4)

	//Set positions
	birdTrMat.Scale(gglm.NewVec3(0.75, 1, 1))
	birdTrMat.Translate(gglm.NewVec3(-5, 0, 0))
	birdTrMat.Rotate(90*gglm.Deg2Rad, gglm.NewVec3(0, 1, 0))
}

func (g *Game) Start() {

	var pipeXSpacing float32 = 6
	pos := gglm.NewVec3(12, 0, 0)
	for i := 0; i < 10; i++ {

		p := NewPipe(true)
		p.Pos.Add(gglm.NewVec3(pos.X(), 10, 0))
		p.TrMat.Translate(gglm.NewVec3(pos.X(), 14, 0))
		p.TrMat.Scale(gglm.NewVec3(0, 10, 0))
		pipes = append(pipes, p)

		p = NewPipe(false)
		p.Pos.Add(gglm.NewVec3(pos.X(), -10, 0))
		p.TrMat.Translate(gglm.NewVec3(pos.X(), -14, 0))
		p.TrMat.Scale(gglm.NewVec3(0, 10, 0))
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
	birdTrMat.Translate(&positionDelta)

	//Move the pipes
	spd := *pipesSpeed
	spd.Scale(timing.DT())
	for i := 0; i < len(pipes); i++ {
		pipes[i].Pos.Add(&spd)
		pipes[i].TrMat.Translate(&spd)
	}

	//Pipe logic
	for i := 0; i < len(pipes); i += 2 {

		xPos := pipes[i].Pos.X()
		if !pipes[i].ShouldRegen {
			continue
		}

		if xPos < 10 || xPos > 12 {
			continue
		}

		pipes[i].ShouldRegen = false
		randY := rand.Float32() * 5
		if rand.Float32() > 0.5 {
			randY *= -1
		}

		pipes[i].Pos.SetY(pipes[i].Pos.Y() + randY)
		pipes[i+1].Pos.SetY(pipes[i].Pos.Y() + randY)

		trAmount := gglm.NewVec3(0, randY, 0)
		pipes[i].TrMat.Translate(trAmount)
		pipes[i+1].TrMat.Translate(trAmount)
	}
}

func (g *Game) Render() {

	//Draw bird
	simpleMat.DiffuseTex = birdTex.TexID
	simpleMat.SetAttribute(birdMesh.Buf)
	g.win.Rend.Draw(birdMesh, birdTrMat, simpleMat)

	//Draw pipe
	simpleMat.DiffuseTex = pipeTex.TexID
	simpleMat.SetAttribute(pipeMesh.Buf)
	simpleMat.Bind()
	for i := 0; i < len(pipes); i++ {
		g.win.Rend.Draw(pipeMesh, pipes[i].TrMat, simpleMat)
	}
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
