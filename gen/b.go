package gen

import (
	"bytes"
	"strings"
)

type B struct {
	*bytes.Buffer
}

func NewB() *B { return &B{Buffer: new(bytes.Buffer)} }

func (b *B) W(s ...string) *B {
	_, _ = b.WriteString(strings.Join(s, ""))
	return b
}

func (b *B) WL(s ...string) *B {
	return b.W(s...).Ln()
}

func (b *B) WL2(s ...string) *B {
	return b.W(s...).Ln2()
}

func (b *B) Tab() *B {
	_, _ = b.WriteString("	")
	return b
}

func (b *B) Spc() *B {
	_, _ = b.WriteString(" ")
	return b
}

func (b *B) Ln() *B {
	_, _ = b.WriteString("\n")
	return b
}

func (b *B) Ln2() *B {
	return b.Ln().Ln()
}
