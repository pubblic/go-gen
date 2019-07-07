package gen

import (
	"bytes"
	"io"
	"sync"

	"golang.org/x/exp/rand"
)

type RuneWriter interface {
	WriteRune(rune) (int, error)
}

type Writer interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
	RuneWriter
}

type Gen interface {
	// Gen generates an actual byte sequence. the provided Writer must not
	// return any error to the caller.
	//
	// Translation to Korea
	// Gen은 실제 데이터를 생성합니다. 제공되는 Writer는 어떤 쓰기에
	// 대해서도 실패하지 않습니다.
	Gen(*rand.Rand, Writer)
}

var bufferFree = sync.Pool{
	New: func() interface{} {
		buf := new(bytes.Buffer)
		buf.Grow(4096)
		return buf
	},
}

func New(source rand.Source, list ...Gen) string {
	buf := bufferFree.Get().(*bytes.Buffer)
	buf.Reset()
	rnd := rand.New(source)
	for _, x := range list {
		x.Gen(rnd, buf)
	}
	ret := buf.String()
	bufferFree.Put(buf)
	return ret
}

type Repeat struct {
	G Gen
	N int
}

func (v Repeat) Gen(r *rand.Rand, w Writer) {
	for i := 0; i < v.N; i++ {
		v.Gen(r, w)
	}
}

type RuneTable []rune

func (v RuneTable) Pick(r *rand.Rand) rune {
	i := r.Intn(len(v))
	return v[i]
}

func (v RuneTable) Gen(r *rand.Rand, w Writer) {
	ru := v.Pick(r)
	w.WriteRune(ru)
}

type ByteTable string

func (v ByteTable) Pick(r *rand.Rand) byte {
	i := r.Intn(len(v))
	return v[i]
}

func (v ByteTable) Gen(r *rand.Rand, w Writer) {
	b := v.Pick(r)
	w.WriteByte(b)
}

type LowerAlphabet int

func (v LowerAlphabet) Gen(r *rand.Rand, w Writer) {
	const some ByteTable = "abcdefghijklmnopqrstuvwxyz"
	n := int(v)
	for i := 0; i < n; i++ {
		some.Gen(r, w)
	}
}

type UpperAlphabet int

func (v UpperAlphabet) Gen(r *rand.Rand, w Writer) {
	const some ByteTable = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	n := int(v)
	for i := 0; i < n; i++ {
		some.Gen(r, w)
	}
}

type Alphabet int

func (v Alphabet) Gen(r *rand.Rand, w Writer) {
	const some ByteTable = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	n := int(v)
	for i := 0; i < n; i++ {
		some.Gen(r, w)
	}
}

type AlphaNumeric int

func (v AlphaNumeric) Gen(r *rand.Rand, w Writer) {
	const some ByteTable = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := int(v)
	for i := 0; i < n; i++ {
		some.Gen(r, w)
	}
}

type Number int

func (v Number) Gen(r *rand.Rand, w Writer) {
	const some ByteTable = "0123456789"
	n := int(v)
	for i := 0; i < n; i++ {
		some.Gen(r, w)
	}
}

type Chain []Gen

func (v Chain) Gen(r *rand.Rand, w Writer) {
	for _, x := range v {
		x.Gen(r, w)
	}
}

type Pick []Gen

func (v Pick) Gen(r *rand.Rand, w Writer) {
	i := r.Intn(len(v))
	v[i].Gen(r, w)
}

type Bytes []byte

func (v Bytes) Gen(r *rand.Rand, w Writer) {
	w.Write(v)
}

type String string

func (v String) Gen(r *rand.Rand, w Writer) {
	w.WriteString(string(v))
}

type Shuffle []Gen

func (v Shuffle) Gen(r *rand.Rand, w Writer) {
	i := len(v) - 1
	for ; i > 0; i-- {
		j := r.Intn(i + 1)
		v[i], v[j] = v[j], v[i]
	}
	Chain(v).Gen(r, w)
}

type Option struct {
	G Gen
}

func (v Option) Gen(r *rand.Rand, w Writer) {
	if parity(getByte(r)) {
		v.G.Gen(r, w)
	}
}

type May struct {
	G Gen
	N int // N is a numerator.
	D int // D is a denominator
}

func (v May) Gen(r *rand.Rand, w Writer) {
	if r.Intn(v.D) < v.N {
		v.G.Gen(r, w)
	}
}

func getByte(r *rand.Rand) byte {
	var buf [1]byte
	r.Read(buf[:])
	return buf[0]
}

func parity(c byte) bool {
	b := (c >> 0 & 1) ^
		(c >> 1 & 1) ^
		(c >> 2 & 1) ^
		(c >> 3 & 1) ^
		(c >> 4 & 1) ^
		(c >> 5 & 1) ^
		(c >> 6 & 1) ^
		(c >> 7 & 1)
	return b == 1
}
