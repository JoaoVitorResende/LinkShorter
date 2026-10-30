package codegen

import "math/rand/v2"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Generator struct {
	length int
}

func New(length int) *Generator {
	if length <= 0 {
		length = 6
	}
	return &Generator{length: length}
}

func (g *Generator) Generate() string {
	b := make([]byte, g.length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}