package battle

type shooter struct {
	u uint8
	f func() (int, int, bool)
}

func (s shooter) shot() (int, int, bool) {
	return s.f()
}

func (s shooter) id() uint8 {
	return s.u
}

type game struct {
	fields [2]field
	hits   []point
	kill   uint8
	ship   map[uint8]uint8
	deck   uint8
	hard   uint8
	shooter
}

func (g *game) initialize(hard uint8, sizes ...uint8) {
	g.fields[0].initialize(sizes...)
	g.fields[1].initialize(sizes...)
	g.shooter = shooter{f: g.random}
	g.ship = make(map[uint8]uint8)
	for _, size := range sizes {
		g.ship[size]++
		g.deck += size
	}
	g.hard = hard
}

func (g *game) Field() Answer {
	var points []point
	for n := range &g.fields {
		for y := range &g.fields[n] {
			for x := range &g.fields[n][y] {
				if n == 1 && g.fields[n].raw(x, y) < 2 {
					continue
				}
				points = append(points, g.fields[n].point(n, x, y))
			}
		}
	}
	return &answer{points: points, state: g.state()}
}

type answer struct {
	points []point
	state  []byte
	point
}

func (a *answer) H() []byte {
	if len(a.points) == 0 {
		return a.state
	}
	return nil
}

func (a *answer) Next() bool {
	if len(a.points) > 0 {
		a.point, a.points = a.points[0], a.points[1:]
		return true
	}
	return false
}

func (g *game) Click(x int, y int) Answer {
	points, hit := g.fields[1].shot(1, x, y)
	if points == nil {
		return nil
	}
	if !hit {
		points = append(points, g.answer()...)
	}
	return &answer{points: points, state: g.state()}
}

func (g *game) answer() []point {
	var points []point
	for {
		x, y, ok := g.shot()
		if !ok {
			break
		}
		shots, hit := g.fields[0].shot(0, x, y)
		points = append(points, shots...)
		if !hit {
			break
		}
		g.add(shots...)
	}
	return points
}

func (g *game) random() (x int, y int, ok bool) {
	if g.kill > 0 {
		g.ship[g.kill]--
		g.kill = 0
	}
	switch g.hard {
	case 1:
		x, y, ok = g.fields[0].weight(0, g.ship).XYZ()
	default:
		x, y, ok = g.fields[0].random(0).XYZ()
	}
	if ok {
		g.shooter = shooter{u: 1, f: g.right}
	}
	g.hits = g.hits[:0]
	return
}

func (g *game) right() (int, int, bool) {
	x, y := g.xy()
	x++
	return g.next(x, y, shooter{u: 2, f: g.left})
}

func (g *game) left() (int, int, bool) {
	x, y := g.xy()
	x--
	return g.next(x, y, shooter{u: 3, f: g.down})
}

func (g *game) down() (int, int, bool) {
	x, y := g.xy()
	y++
	return g.next(x, y, shooter{u: 4, f: g.up})
}

func (g *game) up() (int, int, bool) {
	x, y := g.xy()
	y--
	return g.next(x, y, shooter{u: 0, f: g.random})
}

func (g *game) next(x, y int, shooter shooter) (int, int, bool) {
	if g.fields[0].target(x, y) {
		return x, y, true
	}
	if len(g.hits) > 0 {
		g.hits = g.hits[:1]
	}
	g.shooter = shooter
	return g.shot()
}

func (g *game) add(shots ...point) {
	g.hits = append(g.hits, shots[len(shots)-1])
	g.kill++
	g.deck--
}

func (g *game) xy() (int, int) {
	if len(g.hits) == 0 {
		return -1, -1
	}
	return g.hits[len(g.hits)-1].XY()
}

func (g *game) end() bool {
	return g.deck == 0
}

func (g *game) alive() bool {
	return g.deck > 0
}

func (g *game) state() []byte { return g.compress([]byte{}) }

func (g *game) compress(b []byte) []byte {
	b = g.fields[0].compress(b)
	b = g.fields[1].compress(b)

	var m uint8
	for k := range g.ship {
		m = max(m, k)
	}

	b = append(b, g.shooter.id(), g.kill, g.deck, g.hard, m)
	for m > 0 {
		b = append(b, g.ship[m])
		m--
	}

	b = append(b, byte(len(g.hits)))
	for _, p := range g.hits {
		b = p.compress(b)
	}

	return b
}

func (g *game) decompress(b []byte) []byte {
	if len(b) < fieldSize*2 {
		return nil
	}

	b = g.fields[0].decompress(b)
	b = g.fields[1].decompress(b)

	if len(b) < 5 {
		return nil
	}

	var u, m, h uint8
	u, g.kill, g.deck, g.hard, m = b[0], b[1], b[2], b[3], b[4]
	b = b[5:]

	if len(b) < int(m) {
		return nil
	}

	g.ship = make(map[uint8]uint8, m)
	for m > 0 {
		g.ship[m] = b[0]
		b = b[1:]
		m--
	}

	if len(b) < 1 {
		return nil
	}

	var p point
	h = b[0]
	b = b[1:]

	if len(b) < int(h<<1) {
		return nil
	}

	g.hits = make([]point, 0, h)
	for h > 0 {
		h--
		b = p.decompress(b)
		g.hits = append(g.hits, p)
	}

	if len(b) > 0 {
		return nil
	}

	switch u {
	case 0:
		g.shooter = shooter{u: 0, f: g.random}
	case 1:
		g.shooter = shooter{u: 1, f: g.right}
	case 2:
		g.shooter = shooter{u: 2, f: g.left}
	case 3:
		g.shooter = shooter{u: 3, f: g.down}
	case 4:
		g.shooter = shooter{u: 4, f: g.up}
	default:
		return nil
	}

	return b
}
