package objects

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/smail1111/mario/internal/utils"
)

type Game struct {
	TileMap TileMap

	Assets map[string]*ebiten.Image
	Anims  map[string]utils.Animation

	Goombas map[int]Goomba
	Koopas  map[int]Koopa
	Coins   map[int]Entity

	Particles map[int]Particle

	Mario Mario

	Flags map[int]Entity

	IDCounter func() (id int)

	Music *audio.Player

	Display       *ebiten.Image
	RenderDisplay *ebiten.Image

	Offset FVector

	world int
	stage int

	score int
	coins int
	timer int
	lives int
	frame int
	win   int
	init  bool
}

func (g *Game) Init() {
	if g.world == 0 {
		g.Assets = utils.MustLoadImages("images", map[string]*ebiten.Image{})
		for _, image := range g.Assets {
			utils.SetColorKey(image, 146, 144, 255)
			utils.SetColorKey(image, 92, 148, 252)
		}

		g.Anims = map[string]utils.Animation{
			"mario/small/idle":  utils.MustLoadAnimation("images/entities/mario/small/idle", "mario/small/idle", 1, false),
			"mario/small/run":   utils.MustLoadAnimation("images/entities/mario/small/run", "mario/small/run", marioDefaultAnimLength, true),
			"mario/small/jump":  utils.MustLoadAnimation("images/entities/mario/small/jump", "mario/small/jump", 1, false),
			"mario/small/turn":  utils.MustLoadAnimation("images/entities/mario/small/turn", "mario/small/turn", 1, false),
			"mario/small/swim":  utils.MustLoadAnimation("images/entities/mario/small/swim", "mario/small/swim", marioDefaultAnimLength, true),
			"mario/small/climb": utils.MustLoadAnimation("images/entities/mario/small/climb", "mario/small/climb", marioDefaultAnimLength, true),
			"mario/small/die":   utils.MustLoadAnimation("images/entities/mario/small/die", "mario/small/die", marioDeathAnimLength, false),

			"mario/big/idle":   utils.MustLoadAnimation("images/entities/mario/big/idle", "mario/big/idle", 1, false),
			"mario/big/run":    utils.MustLoadAnimation("images/entities/mario/big/run", "mario/big/run", marioDefaultAnimLength, true),
			"mario/big/jump":   utils.MustLoadAnimation("images/entities/mario/big/jump", "mario/big/jump", 1, false),
			"mario/big/turn":   utils.MustLoadAnimation("images/entities/mario/big/turn", "mario/big/turn", 1, false),
			"mario/big/swim":   utils.MustLoadAnimation("images/entities/mario/small/swim", "mario/big/swim", marioDefaultAnimLength, true),
			"mario/big/climb":  utils.MustLoadAnimation("images/entities/mario/big/climb", "mario/big/climb", marioDefaultAnimLength, true),
			"mario/big/crouch": utils.MustLoadAnimation("images/entities/mario/big/crouch", "mario/big/crouch", 1, false),
			"mario/big/shrink": utils.MustLoadAnimation("images/entities/mario/big/shrink", "mario/big/shrink", marioShrinkAnimLength, false),

			"goomba/idle":  utils.MustLoadAnimation("images/entities/goomba/idle", "goomba/idle", enemyDefaultAnimLength, true),
			"goomba/stomp": utils.MustLoadAnimation("images/entities/goomba/stomp", "goomba/stomp", goombaStompAnimLength, false),

			"koopa/walk":  utils.MustLoadAnimation("images/entities/koopa/walk", "koopa/walk", enemyDefaultAnimLength, true),
			"koopa/shell": utils.MustLoadAnimation("images/entities/koopa/shell", "koopa/shell", koopaShellAnimLength, false),

			"block/question/active": utils.MustLoadAnimation("images/entities/block/question/active", "block/question/active", 12, true),
			"block/question/hit":    utils.MustLoadAnimation("images/entities/block/question/hit", "block/question/hit", 6, false),

			"block/brick/active": utils.MustLoadAnimation("images/entities/block/brick/active", "block/brick/active", 1, false),
			"block/brick/break":  utils.MustLoadAnimation("images/entities/block/brick/break", "block/brick/break", 12, false),

			"coin/spin": utils.MustLoadAnimation("images/entities/coin/spin", "coin/spin", coinSpinAnimLength, true),

			"flag/idle": utils.MustLoadAnimation("images/entities/flag/idle", "flag/idle", 1, false),
		}
		for _, anim := range g.Anims {
			for _, image := range anim.Images {
				utils.SetColorKey(image, 146, 144, 255)
				utils.SetColorKey(image, 92, 148, 252)
			}
		}

		g.Display = ebiten.NewImage(320, 240)

		g.RenderDisplay = ebiten.NewImage(320, 240)

		g.world = 1
		g.stage = 1

		g.lives = startingLives
	}

	g.IDCounter = NewEntityIDCounter()

	g.Mario = g.NewMario(FVector{})

	g.TileMap = g.NewTileMap(tileSize, fmt.Sprintf("%d-%d", g.world, g.stage))

	g.Particles = map[int]Particle{}
	g.Goombas = map[int]Goomba{}
	g.Koopas = map[int]Koopa{}
	g.Coins = map[int]Entity{}
	g.Flags = map[int]Entity{}

	for _, tile := range g.TileMap.Extract(utils.NewSet(
		TileID{"spawner", 0},
		TileID{"spawner", 1},
		TileID{"spawner", 2},
		TileID{"spawner", 3},
		TileID{"spawner", 4},
	), false,
	) {
		switch tile.Variant {
		case 0:
			g.Mario.Pos = tile.Pos.FVector()
		case 1:
			g.NewGoomba(tile.Pos.FVector())
		case 2:
			g.NewKoopa(tile.Pos.FVector())
		case 3:
			g.NewCoin(tile.Pos.FVector())
		case 4:
			basePos := Vector{tile.Pos.X/16 + 2, tile.Pos.Y/16 + 10}

			g.TileMap.Tiles[basePos.String()] = Tile{
				Pos:     basePos,
				Type:    "block",
				Variant: 3,
			}

			flag := g.NewEntity("flag", FVector{float64(tile.Pos.X + 39), float64(tile.Pos.Y + 8)}, 2.0, 152.0, "idle", FVector{-15.0, 0.0})
			g.Flags[flag.ID] = flag
		}
	}

	for _, tile := range g.TileMap.Extract(utils.NewSet(
		TileID{"block", 1},
		TileID{"block", 2},
	), false) {
		coord := Vector{tile.Pos.X / tileSize, tile.Pos.Y / tileSize}

		switch tile.Variant {
		case 1:
			g.TileMap.Blocks[coord.String()] = Block{
				Entity: g.NewEntity("block/question", tile.Pos.FVector(), tileSize, tileSize, "active", FVector{}),
			}
		case 2:
			g.TileMap.Blocks[coord.String()] = Block{
				Entity: g.NewEntity("block/brick", tile.Pos.FVector(), tileSize, tileSize, "active", FVector{}),
			}
		}
	}

	g.Music = utils.PlayLoop("music")

	g.timer = startingTimer

	g.score = startingScore

	g.win = 0

	g.Offset = FVector{}

	g.init = true
}

func (g *Game) Update() error {
	if !g.init {
		g.Init()
	}

	g.frame++

	if g.win == 0 {
		if g.timer > 0 {
			if g.frame%timerTickFrames == 0 {
				g.timer--
			}
		} else {
			g.Mario.Death = marioDeathAmt
		}
	}

	var movement FVector

	if g.win == 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.Mario.Crouch()
		} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
			g.Mario.UnCrouch()
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			movement = g.Mario.Move(left)
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
			movement = g.Mario.Move(right)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.Mario.Jump()
		}

		alive := g.Mario.Update(movement)

		if !alive {
			g.init = false

			if g.lives == 0 {
				g.world = 0
			}
		}
	} else {
		switch g.win {
		case 1:
			g.Mario.Anim.Frame = 11

		case 2:
			if g.timer > 0 {
				g.timer--
				g.score += timerScorePerTick
			} else {
				g.Mario.Fling(flagFlingVelocity, 0.0, false, 0)
				g.win++

				go utils.PlaySound("stage_clear")
			}

		default:
			g.Mario.Update(g.Mario.Move(right))
			g.win++

			if g.win > 360 {
				g.init = false
			}
		}
	}

	for id, goomba := range g.Goombas {
		alive := goomba.Update(FVector{})

		if alive {
			g.Goombas[id] = goomba
		} else {
			defer delete(g.Goombas, id)
		}
	}

	for id, koopa := range g.Koopas {
		alive := koopa.Update(FVector{})

		if alive {
			g.Koopas[id] = koopa
		} else {
			defer delete(g.Koopas, id)
		}
	}

	for pos, block := range g.TileMap.Blocks {
		alive := block.Update()

		if alive {
			g.TileMap.Blocks[pos] = block
		} else {
			defer delete(g.TileMap.Blocks, pos)
		}
	}

	for id, coin := range g.Coins {
		coin.Anim.Update()

		if coin.Overlaps(g.Mario.FRect) {
			g.coins++
			g.score += coinCollectScore

			go utils.PlaySound("coin")

			defer delete(g.Coins, id)
		} else {
			g.Coins[id] = coin
		}
	}

	for _, flag := range g.Flags {
		if g.Mario.Overlaps(flag.FRect) {
			if g.Music.IsPlaying() {
				g.Mario.Pos.Y = min(g.Mario.Pos.Y, flag.Bottom()+minMarioFlagPosition-g.Mario.Height)
				g.Mario.Velocity = FVector{}
				g.Mario.SetAnimation("climb")
				g.Music.Close()
				g.win++
			}

			if g.Mario.Bottom() <= flag.Bottom()+minMarioFlagPosition {
				g.Mario.Pos.Y += flagSlideSpeed
			} else if g.win < 2 {
				g.win++
			}
		}
	}

	if g.coins >= coin1UPAmount {
		g.coins -= coin1UPAmount

		g.lives++
		go utils.PlaySound("1_up")
	}

	for id, particle := range g.Particles {
		alive := particle.Update()

		if alive {
			g.Particles[id] = particle
		} else {
			defer delete(g.Particles, id)
		}
	}

	if g.win == 0 {
		g.Offset.X = max(g.Offset.X, g.Mario.Pos.X-screenWidth/2)
	}

	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	g.Display.Clear()

	g.TileMap.Draw(g.Display, g.Assets, g.Offset)

	for _, flag := range g.Flags {
		flag.Draw(g.Display, g.Offset)
	}

	for _, particle := range g.Particles {
		particle.Draw(g.Display, g.Offset)
	}

	for _, coin := range g.Coins {
		coin.Draw(g.Display, g.Offset)
	}

	for _, goomba := range g.Goombas {
		goomba.Draw(g.Display, g.Offset)
	}

	for _, koopa := range g.Koopas {
		koopa.Draw(g.Display, g.Offset)
	}

	g.Mario.Draw(g.Display, g.Offset)

	g.DrawHud(g.Display)

	g.RenderDisplay.Fill(color.RGBA{92, 148, 252, 255})

	silhouete := utils.Silhouette(g.Display)

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(1.0, 1.0)

	g.RenderDisplay.DrawImage(silhouete, &opts)

	g.RenderDisplay.DrawImage(g.Display, nil)

	opts = ebiten.DrawImageOptions{}
	opts.GeoM.Scale(2.0, 2.0)

	screen.DrawImage(g.RenderDisplay, &opts)
}

func (g Game) DrawHud(screen *ebiten.Image) {
	const hudMarginX = screenWidth / 20

	ebitenutil.DebugPrintAt(screen, "SCORE", hudMarginX, 2)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.score), hudMarginX+15-len(fmt.Sprint(g.score))*3, 16)

	ebitenutil.DebugPrintAt(screen, "COINS", screenWidth/5+hudMarginX, 2)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.coins), screenWidth/5+hudMarginX+15-len(fmt.Sprint(g.coins))*3, 16)

	ebitenutil.DebugPrintAt(screen, "WORLD", screenWidth/5*2+hudMarginX, 2)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d-%d", g.world, g.stage), screenWidth/5*2+hudMarginX+15-len(fmt.Sprintf("%d-%d", g.world, g.stage))*3, 16)

	ebitenutil.DebugPrintAt(screen, "TIMER", screenWidth/5*3+hudMarginX, 2)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.timer), screenWidth/5*3+hudMarginX+15-len(fmt.Sprint(g.timer))*3, 16)

	ebitenutil.DebugPrintAt(screen, "LIVES", screenWidth/5*4+hudMarginX, 2)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.lives), screenWidth/5*4+hudMarginX+15-len(fmt.Sprint(g.lives))*3, 16)
}

func (g Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
