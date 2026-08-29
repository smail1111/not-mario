package objects

import "github.com/smail1111/not-mario/internal/utils"

const (
	screenWidth = 320

	startingLives = 10
	startingTimer = 300
	startingScore = 0

	timerTickFrames = 30

	tileSize = 16

	gravity           = 0.2
	minimumVelocityY  = 0.01
	terminalVelocityY = 5.0

	normalizationX = 0.03

	marioTurnAcceletation = 0.01
	marioWalkAcceleration = 0.06
	marioMaxWalkSpeed     = 1.8
	marioRunSpeedMult     = 1.5

	marioJumpRefreshAmount   = 1
	marioJumpingAcceleration = -0.8
	minimumJumpingVelocity   = -4.1

	marioStompVelocity = -1.0

	marioDeathAnimLength   = 240
	marioShrinkAnimLength  = 25
	marioDefaultAnimLength = 6

	marioShrinkHideAmt = 12

	enemyDefaultAnimLength = 10

	coinCollectScore = 10
	enemyStompScore  = 100

	maxEnemyStompStreak = 6

	goombaSpeed           = 0.4
	goombaStompAnimLength = 30

	koopaWalkSpeed       = 0.5
	koopaShellMoveSpeed  = 3.0
	koopaShellAnimLength = 60

	blockHitAnimOffsetY = -5.0
	blockNormalizationY = 0.5

	coinSpinAnimLength   = 12
	coin1UPAmount        = 100
	coinParticleSpeedY   = -1.0
	coinParticleLifespan = 30

	flagSlideSpeed    = 1.0
	flagFlingVelocity = 2.0

	minMarioFlagPosition = -1.0

	timerScorePerTick = 10

	flingVelocityX = 0.0
	flingVelocityY = -3.0
)

const (
	_ = iota

	goombaStompAmt
	goombaDeathAmount

	koopaShellAmt
	koopaDeathAmt

	marioDeathAmt
)

const (
	left = iota
	right
	up
	down
)

var (
	directions     = utils.NewSet(left, right, up, down)
	collisionTypes = utils.NewSet("pipe", "block", "block/question", "block/brick")
)
