package battle

import (
	"reflect"
)

type shooter func() (int, int, bool)

func (s shooter) shot() (int, int, bool) {
	return s()
}

type game struct {
	fields [2]field
	hits   []point
	kill   int
	ship   map[int]int
	deck   int
	hard   int
	shooter
}

func (g *game) initialize(hard int, sizes ...int) {
	g.fields[0].initialize(sizes...)
	g.fields[1].initialize(sizes...)
	g.shooter = g.random
	g.ship = make(map[int]int)
	for _, size := range sizes {
		g.ship[size]++
		g.deck += size
	}
	g.hard = hard
}

func (g *game) Field() *answer {
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

func (g *game) Click(x int, y int) *answer {
	points, hit := g.fields[1].shot(1, x, y)
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
		g.shooter = g.right
	}
	g.hits = g.hits[:0]
	return
}

func (g *game) right() (int, int, bool) {
	x, y := g.xy()
	x++
	return g.next(x, y, g.left)
}

func (g *game) left() (int, int, bool) {
	x, y := g.xy()
	x--
	return g.next(x, y, g.down)
}

func (g *game) down() (int, int, bool) {
	x, y := g.xy()
	y++
	return g.next(x, y, g.up)
}

func (g *game) up() (int, int, bool) {
	x, y := g.xy()
	y--
	return g.next(x, y, g.random)
}

func (g *game) next(x, y int, s shooter) (int, int, bool) {
	if g.fields[0].target(x, y) {
		return x, y, true
	}
	if len(g.hits) > 0 {
		g.hits = g.hits[:1]
	}
	g.shooter = s
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

func (g *game) shooterID() byte {
	ptr := reflect.ValueOf(g.shooter).Pointer()
	switch ptr {
	case reflect.ValueOf(g.right).Pointer():
		return 1
	case reflect.ValueOf(g.left).Pointer():
		return 2
	case reflect.ValueOf(g.down).Pointer():
		return 3
	case reflect.ValueOf(g.up).Pointer():
		return 4
	default:
		return 0
	}
}

func (g *game) compress(b []byte) []byte {
	b = g.fields[0].compress(b)
	b = g.fields[1].compress(b)

	var m int
	for k := range g.ship {
		m = max(m, k)
	}

	b = append(b, g.shooterID(), byte(g.kill), byte(g.deck), byte(g.hard), byte(m))
	for i := 1; i <= m; i++ {
		b = append(b, byte(g.ship[i]))
	}

	b = append(b, byte(len(g.hits)))
	for _, p := range g.hits {
		b = p.compress(b)
	}

	return b
}

func (g *game) decompress(b []byte) []byte {
	b = g.fields[0].decompress(b)
	b = g.fields[1].decompress(b)

	var s, m, h int
	s, g.kill, g.deck, g.hard, m = int(b[0]), int(b[1]), int(b[2]), int(b[3]), int(b[4])
	b = b[5:]
	g.ship = make(map[int]int, m)
	for i := 1; i <= m; i++ {
		g.ship[i] = int(b[0])
		b = b[1:]
	}

	var p point
	h = int(b[0])
	b = b[1:]
	g.hits = make([]point, h)
	for i := 0; i < h; i++ {
		b = p.decompress(b)
		g.hits[i] = p
	}

	switch s {
	case 0:
		g.shooter = g.random
	case 1:
		g.shooter = g.right
	case 2:
		g.shooter = g.left
	case 3:
		g.shooter = g.down
	case 4:
		g.shooter = g.up
	}

	return b
}
