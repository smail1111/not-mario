package objects

func (g *Game) NewCoin(pos FVector) Entity {
	coin := g.NewEntity("coin", pos, 8.0, 14.0, "spin", FVector{0.0, -2.0})

	g.Coins[coin.ID] = coin

	return coin
}

type Particle struct {
	Entity
	Movement FVector
	LifeSpan int
	frame    int
}

func (p *Particle) Update() (alive bool) {
	p.frame++

	p.Pos.X += p.Movement.X

	p.Pos.Y += p.Movement.Y

	p.Anim.Update()

	return p.frame < p.LifeSpan
}
