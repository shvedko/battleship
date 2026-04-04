package battle

import (
	"reflect"
	"testing"
)

func benchmark(b *testing.B, a int) {
	var n, c int
	for i := 0; i < b.N; i++ {
		var g game
		g.initialize(a, 4, 3, 3, 2, 2, 2)
		for g.alive() {
			p := g.answer()
			n += len(p)
			c++
		}
	}
	b.ReportMetric(float64(n)/float64(b.N), "shots/op")
	b.ReportMetric(float64(c)/float64(b.N), "moves/op")
}

func Benchmark_game(b *testing.B) {
	for i, n := range []string{"Random", "Weight"} {
		b.Run(n, func(b *testing.B) {
			benchmark(b, i)
		})
	}
}

func Test_game_compress_decompress(t *testing.T) {
	q := game{
		fields: [2]field{{
			[10]int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			[10]int{1, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			[10]int{1, 0, 1, 0, 1, 0, 0, 0, 1, 0},
			[10]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
			[10]int{0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			[10]int{0, 1, 0, 1, 0, 0, 1, 1, 1, 1},
			[10]int{0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			[10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int{0, 0, 1, 0, 0, 1, 1, 0, 1, 0},
			[10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, {
			[10]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[10]int{0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			[10]int{0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
			[10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[10]int{0, 1, 0, 0, 0, 0, 0, 1, 0, 0},
			[10]int{0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
			[10]int{1, 1, 1, 1, 0, 0, 0, 1, 0, 0},
			[10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[10]int{0, 1, 0, 0, 1, 0, 0, 0, 0, 0},
			[10]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0}}},
		hits:    []point{7330},
		kill:    1,
		ship:    map[int]int{1: 4, 2: 3, 3: 2, 4: 1},
		deck:    19,
		hard:    1,
		shooter: nil,
	}
	q.shooter = q.up

	b := make([]byte, 0, 128)
	b = q.compress(b)

	t.Log(b)

	e := []byte{0, 16, 0, 0, 0, 16, 16, 0, 0, 0, 16, 16, 16, 0, 16, 0, 0, 0, 0, 16, 0, 1, 0, 0, 0, 1, 1, 0, 17, 17, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 1, 16, 16, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 17, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 17, 17, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 16, 0, 0, 0, 0, 16, 0, 0, 4, 1, 19, 1, 4, 4, 3, 2, 1, 1, 28, 162}
	if !reflect.DeepEqual(b, e) {
		t.Errorf("%v != %v", b, e)
	}
	e = []byte{}

	p := game{}
	b = p.decompress(b)

	t.Log(b)
	t.Log(q)
	t.Log(p)

	if !reflect.DeepEqual(b, e) {
		t.Errorf("%v != %v", b, e)
	}

	i := q.shooterID()
	j := p.shooterID()
	q.shooter = nil
	p.shooter = nil

	if !reflect.DeepEqual(i, j) {
		t.Errorf("%v != %v", i, j)
	}
	if !reflect.DeepEqual(q, p) {
		t.Errorf("%v != %v", q, p)
	}
}
