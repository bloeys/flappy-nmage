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
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/veandco/go-sdl2/sdl"
)

//TODO: Sounds!

const (
	pipeXSpacing float32 = 6
	pipeYSpacing float32 = 6
)

type GameState int

const (
	GameState_Playing GameState = iota
	GameState_Lost
)

var _ engine.Game = &Game{}

type Game struct {
	win       *engine.Window
	imguiInfo nmageimgui.ImguiInfo
	gameState GameState

	BaseObj *quads.BoxCollider2D
}

var (
	game *Game

	simpleMat, digitMat *materials.Material

	birdSprite *quads.Quad
	birdBoxCol *quads.BoxCollider2D

	backgroundSprite *quads.Quad

	gravity float32 = -9.81 * 4

	pipes      []Pipe
	pipesSpeed = gglm.NewVec3(-5, 0, 0)

	birdVelocity         = gglm.NewVec3(0, 0, 0)
	dragAmount   float32 = 0.9
	jumpForce    float32 = 30

	font imgui.Font

	score                uint
	lastTouchedMiddleCol *quads.BoxCollider2D

	//Digit quads
	digitQuads []*quads.Quad
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

	game = &Game{
		win:       window,
		imguiInfo: nmageimgui.NewImGUI(),
		gameState: GameState_Playing,
	}

	engine.Run(game, game.win, game.imguiInfo)
}

func (g *Game) Init() {

	var err error

	//Load assets
	birdSprite, err = quads.NewQuad("bird", "./res/textures/yellowbird-midflap.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}

	birdBoxCol = quads.NewBoxCollider2D(1, 1)

	backgroundSprite, err = quads.NewQuad("background", "./res/textures/background-day.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}

	g.BaseObj = quads.NewBoxCollider2D(1, 1)
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}

	tempQ, err := quads.NewQuad("base", "./res/textures/base.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	g.BaseObj.Quad = *tempQ

	simpleMat = materials.NewMaterial("simpleMat", "./res/shaders/simple")
	digitMat = materials.NewMaterial("digitMat", "./res/shaders/simple")

	//Setup camera
	camPos := gglm.NewVec3(0, 0, -10)
	// camForward := gglm.NewVec3(0, 0, 1)
	targetPos := gglm.NewVec3(0, 0, 0)

	viewMat := gglm.LookAt(camPos, targetPos, gglm.NewVec3(0, 1, 0))
	simpleMat.SetUnifMat4("viewMat", &viewMat.Mat4)
	digitMat.SetUnifMat4("viewMat", &viewMat.Mat4)

	projMat := gglm.Ortho(-10, 10, 10, -10, 0.1, 500)
	simpleMat.SetUnifMat4("projMat", &projMat.Mat4)
	digitMat.SetUnifMat4("projMat", &projMat.Mat4)

	//Fonts
	font = g.imguiInfo.AddFontTTF("./res/fonts/courier-prime.regular.ttf", 64, nil, nil)

	//Set positions
	birdSprite.ScaleReadWrite().Set(1.5, 2, 1)
	birdSprite.PosReadWrite().Set(-5, 0, 2)

	backgroundSprite.Entity.ScaleReadWrite().Set(20, 20, 1)

	g.BaseObj.Entity.ScaleReadWrite().Set(20, 5, 1)
	g.BaseObj.Entity.PosReadWrite().Set(0, -11, 1.1)

	g.LoadDigitQuads()
	g.InitPipes()
}

func (g *Game) LoadDigitQuads() {

	digitQuads = make([]*quads.Quad, 0, 10)

	digit, err := quads.NewQuad("background", "./res/textures/0.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/1.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/2.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/3.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/4.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/5.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/6.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/7.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/8.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)

	digit, err = quads.NewQuad("background", "./res/textures/9.png")
	if err != nil {
		panic("Failed to load sprite. Err:" + err.Error())
	}
	digitQuads = append(digitQuads, digit)
}

func (g *Game) InitPipes() {

	pos := gglm.NewVec3(12, 0, 0)
	for i := 0; i < 10; i++ {

		randY := rand.Float32() * 4
		if rand.Float32() > 0.5 {
			randY *= -1
		}

		//Top pipe and its collider
		p := NewPipe(true)
		p.PosReadWrite().Add(gglm.NewVec3(pos.X(), 10+randY+pipeYSpacing*0.5, 1))
		p.ScaleReadWrite().Set(2, 20, 1)

		pipePos := p.PosRead()
		p.Col.PosReadWrite().Set(pipePos.X(), pipePos.Y(), pipePos.Z()+0.1)
		*p.Col.ScaleReadWrite() = *p.ScaleRead()

		*p.MiddleCol.PosReadWrite() = *pipePos
		p.MiddleCol.ScaleReadWrite().Set(0.1, 30, 1)
		pipes = append(pipes, p)

		//Bottom pipe and its collider
		p = NewPipe(false)
		p.PosReadWrite().Add(gglm.NewVec3(pos.X(), -10+randY-pipeYSpacing*0.5, 1))
		p.ScaleReadWrite().Set(2, 20, 1)

		pipePos = p.PosRead()
		p.Col.PosReadWrite().Set(pipePos.X(), pipePos.Y(), pipePos.Z()+0.1)
		*p.Col.ScaleReadWrite() = *p.ScaleRead()

		pipes = append(pipes, p)

		pos.SetX(pos.X() + pipeXSpacing)
	}
}

func (g *Game) Update() {

	if input.IsQuitClicked() || input.KeyClicked(sdl.K_ESCAPE) {
		engine.Quit()
	}

	switch g.gameState {
	case GameState_Playing:
		g.Playing()
	case GameState_Lost:
		g.Lost()
	}
}

func (g *Game) Playing() {
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

	//Move bird collider
	birdPos := birdSprite.PosRead()
	birdScale := birdSprite.ScaleRead()
	birdBoxCol.PosReadWrite().Set(birdPos.X(), birdPos.Y(), birdPos.Z()+0.1)
	*birdBoxCol.ScaleReadWrite() = *birdScale

	//Move the pipes
	spd := *pipesSpeed
	spd.Scale(timing.DT())
	for i := 0; i < len(pipes); i++ {
		pipes[i].PosReadWrite().Add(&spd)

		pipePos := pipes[i].PosRead()
		pipes[i].Col.PosReadWrite().Set(pipePos.X(), pipePos.Y(), pipePos.Z()+0.1)

		if pipes[i].IsTop {
			*pipes[i].MiddleCol.PosReadWrite() = *pipePos
		}
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

		pipePos := pipes[i].PosRead()
		pipes[i].Col.PosReadWrite().Set(pipePos.X(), pipePos.Y(), pipePos.Z()+0.1)

		pipePos = pipes[i+1].PosRead()
		pipes[i+1].Col.PosReadWrite().Set(pipePos.X(), pipePos.Y(), pipePos.Z()+0.1)
	}

	//Pipe collisions
	for i := 0; i < len(pipes); i++ {
		if !isColliding(birdBoxCol, pipes[i].Col) {
			continue
		}

		g.gameState = GameState_Lost
	}

	//Pipe score collisions
	for i := 0; i < len(pipes); i++ {

		if !pipes[i].IsTop || !isColliding(birdBoxCol, pipes[i].MiddleCol) {
			continue
		}

		//Don't score twice in a row with the same collider
		if pipes[i].MiddleCol == lastTouchedMiddleCol {
			continue
		}

		score++
		lastTouchedMiddleCol = pipes[i].MiddleCol
	}

	//Base collision
	if isColliding(birdBoxCol, g.BaseObj) {
		g.gameState = GameState_Lost
	}

	//Going too high
	if birdPos.Y() >= 12 {
		g.gameState = GameState_Lost
	}
}

func (g *Game) DrawScore() {

	score := score
	origScore := score

	digitCount := 0
	if score == 0 {
		digitCount = 1
	}

	for score > 0 {
		digitCount++
		score /= 10
	}

	score = origScore
	var xOffsetDelta float32 = 1.6
	xOffset := float32(digitCount-1) * xOffsetDelta * 0.5
	for i := 0; i < 1 || score > 0; i++ {

		digit := score % 10
		score /= 10

		q := digitQuads[digit]
		q.PosReadWrite().Set(xOffset, 9, 7)
		q.ScaleReadWrite().Set(1.5, 1.5, 1)
		xOffset -= xOffsetDelta

		digitMat.DiffuseTex = q.Tex.TexID
		digitMat.Bind()
		g.win.Rend.Draw(digitQuads[0].Mesh, q.ModelMat(), digitMat)
	}
}

func (g *Game) Lost() {

	open := true

	w, h := game.win.SDLWin.GetSize()

	lostText := "You Lost!\nPress Enter to restart..."
	textSize := imgui.CalcTextSize(lostText, false, 100000)

	imgui.SetNextWindowPos(imgui.Vec2{})
	imgui.SetNextWindowSize(imgui.Vec2{X: float32(w), Y: float32(h)})
	imgui.BeginV("lost",
		&open,
		imgui.WindowFlagsNoBackground|imgui.WindowFlagsNoCollapse|
			imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|
			imgui.WindowFlagsNoDecoration)

	imgui.PushFont(font)
	imgui.SetCursorPos(imgui.Vec2{X: (float32(w)/2 - textSize.X*2), Y: (float32(h)/2 - textSize.Y)})

	imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, W: 1})
	imgui.Text(lostText)
	imgui.PopStyleColor()

	imgui.PopFont()
	imgui.End()

	//Reset pipe positions and restart the game
	if input.KeyClicked(sdl.K_RETURN) || input.KeyClicked(sdl.K_RETURN2) {

		score = 0
		g.gameState = GameState_Playing

		//Reset bird state
		birdVelocity = gglm.NewVec3(0, 0, 0)
		birdSprite.PosReadWrite().Set(-5, 0, 2)

		//Reset pipe state
		pos := gglm.NewVec3(12, 0, 0)
		for i := 0; i < len(pipes); i += 2 {

			randY := rand.Float32() * 5
			if rand.Float32() > 0.5 {
				randY *= -1
			}

			//Top pipe and its collider
			p := &pipes[i]
			*p.PosReadWrite() = *gglm.NewVec3(pos.X(), 10+randY+pipeYSpacing*0.5, 1)

			pipePos := p.PosRead()
			p.Col.PosReadWrite().Set(pipePos.X(), pipePos.Y(), pipePos.Z()+0.1)
			*p.MiddleCol.PosReadWrite() = *pipePos

			//Bottom pipe and its collider
			p = &pipes[i+1]
			*p.PosReadWrite() = *gglm.NewVec3(pos.X(), -10+randY-pipeYSpacing*0.5, 1)

			pipePos = p.PosRead()
			p.Col.PosReadWrite().Set(pipePos.X(), pipePos.Y(), pipePos.Z()+0.1)
			pos.SetX(pos.X() + pipeXSpacing)
		}
	}
}

func isColliding(a, b *quads.BoxCollider2D) bool {

	aBotLeft := a.BotLeft()
	aTopRight := a.TopRight()

	bBotLeft := b.BotLeft()
	bTopRight := b.TopRight()

	return aBotLeft.X() < bTopRight.X() &&
		bBotLeft.X() < aTopRight.X() &&

		bBotLeft.Y() < aTopRight.Y() &&
		aBotLeft.Y() < bTopRight.Y()
}

func (g *Game) Render() {

	//Draw background
	simpleMat.DiffuseTex = backgroundSprite.Tex.TexID
	simpleMat.Bind()
	g.win.Rend.Draw(backgroundSprite.Mesh, backgroundSprite.ModelMat(), simpleMat)

	//Draw base
	simpleMat.DiffuseTex = g.BaseObj.Tex.TexID
	simpleMat.Bind()
	g.win.Rend.Draw(g.BaseObj.Mesh, g.BaseObj.ModelMat(), simpleMat)

	//Draw pipe
	simpleMat.DiffuseTex = pipes[0].Tex.TexID
	simpleMat.Bind()
	for i := 0; i < len(pipes); i++ {
		g.win.Rend.Draw(pipes[0].Mesh, pipes[i].ModelMat(), simpleMat)
	}

	// //Draw pipe colliders
	// simpleMat.DiffuseTex = pipes[0].Col.Tex.TexID
	// simpleMat.Bind()

	// //First pipe collider is red
	// simpleMat.SetUnifVec3("tintColor", gglm.NewVec3(1, 0, 0))
	// g.win.Rend.Draw(pipes[0].Col.Mesh, pipes[0].Col.ModelMat(), simpleMat)

	// simpleMat.SetUnifVec3("tintColor", gglm.NewVec3(0, 0, 1))
	// for i := 1; i < len(pipes); i++ {
	// 	g.win.Rend.Draw(pipes[0].Col.Mesh, pipes[i].Col.ModelMat(), simpleMat)
	// }
	// simpleMat.SetUnifVec3("tintColor", gglm.NewVec3(1, 1, 1))

	//Draw bird
	simpleMat.DiffuseTex = birdSprite.Tex.TexID
	simpleMat.Bind()
	g.win.Rend.Draw(birdSprite.Mesh, birdSprite.ModelMat(), simpleMat)

	// //Draw bird collider
	// simpleMat.SetUnifVec3("tintColor", gglm.NewVec3(0, 1, 0))
	// simpleMat.DiffuseTex = birdBoxCol.Tex.TexID
	// simpleMat.Bind()
	// g.win.Rend.Draw(birdBoxCol.Mesh, birdBoxCol.ModelMat(), simpleMat)
	// simpleMat.SetUnifVec3("tintColor", gglm.NewVec3(1, 1, 1))

	// simpleMat.SetUnifVec3("tintColor", gglm.NewVec3(1, 0, 0))
	// simpleMat.DiffuseTex = col1.Tex.TexID
	// simpleMat.Bind()
	// g.win.Rend.Draw(col1.Mesh, col1.ModelMat(), simpleMat)
	// g.win.Rend.Draw(col1.Mesh, col2.ModelMat(), simpleMat)
	// simpleMat.SetUnifVec3("tintColor", gglm.NewVec3(1, 1, 1))

	g.DrawScore()
}

func (g *Game) FrameEnd() {

}

func (g *Game) DeInit() {

}
