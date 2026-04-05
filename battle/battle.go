package battle

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type Answer interface {
	Next() bool
	F() int
	X() int
	Y() int
	C() int
	H() []byte
}

type Battle interface {
	Begin([]byte) Answer
	Click(int, int, []byte, []byte) Answer
	Reset([]byte)
	Encryption([32]byte) error
}

type battle struct {
	sizes []uint8
	level uint8
	coder cipher.AEAD
}

func (b *battle) Encryption(key [32]byte) (err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return
	}
	b.coder, err = cipher.NewGCM(block)
	return
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

func (b *battle) decrypt(k []byte, p []byte) []byte {
	if b.coder != nil {
		z := b.coder.NonceSize()
		if len(p) < z {
			return nil
		}
		var q []byte
		q, p = p[:z], p[z:]
		var err error
		p, err = b.coder.Open(p[:0], q, p, k[:])
		if err != nil {
			return nil
		}
	}
	return p
}

type encryptor struct {
	coder cipher.AEAD
	key   []byte
	Answer
}

func (e *encryptor) H() []byte {
	p := e.Answer.H()
	if p == nil {
		return nil
	}
	z := e.coder.NonceSize()
	q := make([]byte, z, z+len(p))
	_, err := rand.Read(q[:z])
	if err != nil {
		return nil
	}
	return e.coder.Seal(q[:z], q, p, e.key)
}

func (b *battle) encrypt(k []byte, a Answer) Answer {
	if b.coder != nil {
		return &encryptor{Answer: a, coder: b.coder, key: k}
	}
	return a
}

func (b *battle) Begin(k []byte) Answer {
	g := b.new()
	if g == nil {
		return nil
	}
	return b.encrypt(k, g.Field())
}

func (b *battle) Click(x int, y int, p []byte, k []byte) Answer {
	p = b.decrypt(k, p)
	g := b.get(p)
	if g == nil {
		return nil
	}
	return b.encrypt(k, g.Click(x, y))
}

func (b *battle) Reset([]byte) {}
