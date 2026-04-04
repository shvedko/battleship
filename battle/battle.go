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
	sizes []uint8
	level uint8
}

func New(level uint8, sizes ...uint8) Battle {
	return &battle{
		sizes: sizes,
		level: level,
	}
}

func (b *battle) new() *game {
	g := &game{}
	g.initialize(b.level, b.sizes...)
	return g
}

func (b *battle) get(p []byte) *game {
	g := &game{}
	p = g.decompress(p)
	if p == nil {
		return nil
	}
	return g
}

func (b *battle) Begin() Answer {
	g := b.new()
	if g == nil {
		return nil
	}
	return g.Field()
}

func (b *battle) Click(x, y int, p []byte) Answer {
	g := b.get(p)
	if g == nil {
		return nil
	}
	return g.Click(x, y)
}

func (b *battle) Reset() {}
