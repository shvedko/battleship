package battle

import (
	"math"
	"reflect"
	"sync/atomic"
	"testing"
)

func gameOne(h uint8) (n int, m int) {
	var g game
	g.initialize(h, 4, 3, 3, 2, 2, 2)
	for g.alive() {
		p := g.answer()
		for i := range p {
			if p[i].F() == 0 && p[i].C() != fieldOpen {
				n++
			}
		}
		m++
	}
	return
}

func Benchmark_game(b *testing.B) {
	for i, t := range map[uint8]string{0: "Random", 1: "Weight"} {
		b.Run(t, func(b *testing.B) {
			var n, m, s, z atomic.Int64
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					dn, dm := gameOne(i)
					n.Add(int64(dn))
					m.Add(int64(dm))
					s.Add(int64(dm * dm))
					z.Add(int64(dn * dn))
				}
			})
			b.ReportMetric(float64(n.Load())/float64(b.N), "shots/op")
			b.ReportMetric(float64(m.Load())/float64(b.N), "moves/op")
			b.ReportMetric(math.Sqrt(float64(s.Load())/float64(b.N)-float64(m.Load())/float64(b.N)*float64(m.Load())/float64(b.N)), "σ(moves)")
			b.ReportMetric(math.Sqrt(float64(z.Load())/float64(b.N)-float64(n.Load())/float64(b.N)*float64(n.Load())/float64(b.N)), "σ(shots)")
		})
	}
}

func gameVersus(h1, h2 uint8, swap bool) (n1 int, n2 int) {
	var g1, g2 game
	if swap {
		h1, h2 = h2, h1
	}
	g1.initialize(h1, 4, 3, 3, 2, 2, 2)
	g2.initialize(h2, 4, 3, 3, 2, 2, 2)
	g2.fields[0] = g1.fields[0]
	for {
		if g1.answer(); g1.end() {
			n1++
			break
		}
		if g2.answer(); g2.end() {
			n2++
			break
		}
	}
	if swap {
		return n2, n1
	}
	return n1, n2
}

func Benchmark_game_Random_vs_Weight(b *testing.B) {
	var n1, n2 atomic.Int64
	b.RunParallel(func(pb *testing.PB) {
		var i uint
		for pb.Next() {
			d1, d2 := gameVersus(0, 1, i&1 == 0)
			n1.Add(int64(d1))
			n2.Add(int64(d2))
			i++
		}
	})
	b.ReportMetric(float64(n1.Load())/float64(b.N)*100, "random")
	b.ReportMetric(float64(n2.Load())/float64(b.N)*100, "weight")
}

func Test_game_compress_decompress(t *testing.T) {
	q := game{
		fields: [2]field{{
			{0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			{1, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			{1, 0, 1, 0, 1, 0, 0, 0, 1, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
			{0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 1, 0, 1, 0, 0, 1, 1, 1, 1},
			{0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0, 1, 1, 0, 1, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, {
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{0, 1, 0, 0, 0, 0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
			{1, 1, 1, 1, 0, 0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{0, 1, 0, 0, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0, 0}}},
		hits: []point{7330, 1550},
		kill: 1,
		ship: map[uint8]uint8{1: 4, 2: 3, 3: 2, 4: 1},
		deck: 19,
		hard: 1,
	}
	q.shooter = shooter{4, q.up}

	b := make([]byte, 0, 128)
	b = q.compress(b)

	e := []byte{0, 16, 0, 0, 0, 16, 16, 0, 0, 0, 16, 16, 16, 0, 16, 0, 0, 0, 0, 16, 0, 1, 0, 0, 0, 1, 1, 0, 17, 17, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 1, 16, 16, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 17, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 17, 17, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 16, 0, 0, 0, 0, 16, 0, 0, 4, 1, 19, 1, 4, 1, 2, 3, 4, 2, 28, 162, 6, 14}
	if !reflect.DeepEqual(b, e) {
		t.Errorf("%v != %v", b, e)
	}
	e = []byte{}

	p := game{}
	b = p.decompress(b)

	if !reflect.DeepEqual(b, e) {
		t.Errorf("%v != %v", b, e)
	}

	i := q.shooter.id()
	j := p.shooter.id()
	q.shooter = shooter{}
	p.shooter = shooter{}

	if !reflect.DeepEqual(i, j) {
		t.Errorf("%v != %v", i, j)
	}
	if !reflect.DeepEqual(q, p) {
		t.Errorf("%v != %v", q, p)
	}
}
