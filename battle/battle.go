package battle

type Answer interface {
	Next() bool
	F() int
	X() int
	Y() int
	C() int
	H() []byte
}

type Battle interface {
	Begin() Answer
	Click(x, y int, p []byte) Answer
	Reset()
}

type battle struct {
	sizes []int
	level int
}

func New(level int, sizes ...int) Battle {
	return &battle{
		sizes: sizes,
		level: level,
	}
}

func (b *battle) begin() *game {
	g := &game{}
	g.initialize(b.level, b.sizes...)
	return g
}

func (b *battle) unpack(p []byte) *game {
	g := &game{}
	p = g.decompress(p)
	return g
}

func (b *battle) Begin() Answer { return b.begin().Field() }

func (b *battle) Click(x, y int, p []byte) Answer { return b.unpack(p).Click(x, y) }

func (b *battle) Reset() {}
