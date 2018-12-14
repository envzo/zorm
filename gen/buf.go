package gen

import (
	"bytes"
	"strings"
)

type Buf struct {
	*bytes.Buffer
}

func NewBuf() *Buf { return &Buf{Buffer: new(bytes.Buffer)} }

func (b *Buf) W(s ...string) *Buf {
	_, _ = b.WriteString(strings.Join(s, ""))
	return b
}

func (b *Buf) WL(s ...string) *Buf {
	return b.W(s...).Ln()
}

func (b *Buf) WL2(s ...string) *Buf {
	return b.W(s...).Ln2()
}

func (b *Buf) Tab() *Buf {
	_, _ = b.WriteString("	")
	return b
}

func (b *Buf) Spc() *Buf {
	_, _ = b.WriteString(" ")
	return b
}

func (b *Buf) Ln() *Buf {
	_, _ = b.WriteString("\n")
	return b
}

func (b *Buf) Ln2() *Buf {
	return b.Ln().Ln()
}
