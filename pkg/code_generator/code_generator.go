package code_generator

import (
	"math/rand"
	"strconv"
	"time"
)

type Generator interface {
	GetCode() string
}

type CodeGenerator struct{}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (c *CodeGenerator) GetCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 1001
	max := 9999
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}
