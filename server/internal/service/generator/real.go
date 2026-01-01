package generator

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

const lengthPassword = 15

var passwordChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.")

type Real struct {
	r *rand.Rand
}

func NewReal() *Real {
	return &Real{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *Real) GenPassword() string {
	var password bytes.Buffer
	for range lengthPassword {
		char := passwordChars[g.r.Intn(len(passwordChars))]
		password.WriteRune(char)
	}
	return password.String()
}

func (g *Real) GenLogin() string {
	const n = 1000000
	number := g.r.Intn(n)
	return fmt.Sprintf("u%d", number)
}
