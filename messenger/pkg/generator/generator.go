package generator

import (
	"math"
	"math/rand/v2"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateID() int64 {
	return rand.Int64N(math.MaxInt64) + 1
}
