package buffer // import "github.com/tdewolff/buffer"

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/test"
)

func TestBufferPool(t *testing.T) {
	z := &bufferPool{}

	lorem := []byte("Lorem ipsum")
	dolor := []byte("dolor sit amet")
	consectetur := []byte("consectetur adipiscing elit")

	// set lorem as first buffer and get new dolor buffer
	b := z.swap(lorem, len(dolor))
	assert.Equal(t, 0, len(b))
	assert.Equal(t, len(dolor), cap(b))
	b = append(b, dolor...)

	// free first buffer so it will be reused
	z.free(len(lorem))
	b = z.swap(b, len(lorem))
	b = b[:len(lorem)]
	assert.Equal(t, lorem, b)

	b = z.swap(b, len(consectetur))
	b = append(b, consectetur...)

	// free in advance to reuse the same buffer
	z.free(len(dolor) + len(lorem) + len(consectetur))
	assert.Equal(t, 0, z.head)
	b = z.swap(b, len(consectetur))
	b = b[:len(consectetur)]
	assert.Equal(t, consectetur, b)

	// free in advance but request larger buffer
	z.free(len(consectetur))
	b = z.swap(b, len(consectetur)+1)
	b = append(b, consectetur...)
	b = append(b, '.')
	assert.Equal(t, len(consectetur)+1, cap(b))
}

func TestLexer(t *testing.T) {
	s := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.`
	z := NewLexer(bytes.NewBufferString(s))

	assert.Equal(t, io.EOF, z.err, "buffer must be fully in memory")
	assert.Equal(t, nil, z.Err(), "buffer is at EOF but must not return EOF until we reach that")
	assert.Equal(t, 0, z.Pos(), "buffer must start at position 0")
	assert.Equal(t, byte('L'), z.Peek(0), "first character must be 'L'")
	assert.Equal(t, byte('o'), z.Peek(1), "second character must be 'o'")

	z.Move(1)
	assert.Equal(t, byte('o'), z.Peek(0), "must be 'o' at position 1")
	assert.Equal(t, byte('r'), z.Peek(1), "must be 'r' at position 1")
	z.Rewind(6)
	assert.Equal(t, byte('i'), z.Peek(0), "must be 'i' at position 6")
	assert.Equal(t, byte('p'), z.Peek(1), "must be 'p' at position 7")

	assert.Equal(t, []byte("Lorem "), z.Lexeme(), "buffered string must now read 'Lorem ' when at position 6")
	assert.Equal(t, []byte("Lorem "), z.Shift(), "shift must return the buffered string")
	assert.Equal(t, len("Lorem "), z.ShiftLen(), "shifted length must equal last shift")
	assert.Equal(t, 0, z.Pos(), "after shifting position must be 0")
	assert.Equal(t, byte('i'), z.Peek(0), "must be 'i' at position 0 after shifting")
	assert.Equal(t, byte('p'), z.Peek(1), "must be 'p' at position 1 after shifting")
	assert.Nil(t, z.Err(), "error must be nil at this point")

	z.Move(len(s) - len("Lorem ") - 1)
	assert.Nil(t, z.Err(), "error must be nil just before the end of the buffer")
	z.Skip()
	assert.Equal(t, 0, z.Pos(), "after skipping position must be 0")
	z.Move(1)
	assert.Equal(t, io.EOF, z.Err(), "error must be EOF when past the buffer")
	z.Move(-1)
	assert.Nil(t, z.Err(), "error must be nil just before the end of the buffer, even when it has been past the buffer")
	z.Free(0) // has already been tested
}

func TestLexerShift(t *testing.T) {
	s := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.`
	z := NewLexerSize(test.NewPlainReader(bytes.NewBufferString(s)), 5)

	z.Move(len("Lorem "))
	assert.Equal(t, []byte("Lorem "), z.Shift(), "shift must return the buffered string")
	assert.Equal(t, len("Lorem "), z.ShiftLen(), "shifted length must equal last shift")

}

func TestLexerSmall(t *testing.T) {
	s := `abcdefghijklm`
	z := NewLexerSize(test.NewPlainReader(bytes.NewBufferString(s)), 4)
	assert.Equal(t, "i", string(z.Peek(8)), "first character must be 'i' at position 8")

	z = NewLexerSize(test.NewPlainReader(bytes.NewBufferString(s)), 4)
	assert.Equal(t, "m", string(z.Peek(12)), "first character must be 'm' at position 12")

	z = NewLexerSize(test.NewPlainReader(bytes.NewBufferString(s)), 0)
	assert.Equal(t, "e", string(z.Peek(4)), "first character must be '4' at position 4")

	z = NewLexerSize(test.NewPlainReader(bytes.NewBufferString(s)), 13)
	assert.Equal(t, byte(0), z.Peek(13), "thirteenth character must yield error")
}

func TestLexerRunes(t *testing.T) {
	z := NewLexer(bytes.NewBufferString("aæ†\U00100000"))
	r, n := z.PeekRune(0)
	assert.Equal(t, 1, n, "first character must be length 1")
	assert.Equal(t, 'a', r, "first character must be rune 'a'")
	r, n = z.PeekRune(1)
	assert.Equal(t, 2, n, "second character must be length 2")
	assert.Equal(t, 'æ', r, "second character must be rune 'æ'")
	r, n = z.PeekRune(3)
	assert.Equal(t, 3, n, "fourth character must be length 3")
	assert.Equal(t, '†', r, "fourth character must be rune '†'")
	r, n = z.PeekRune(6)
	assert.Equal(t, 4, n, "seventh character must be length 4")
	assert.Equal(t, '\U00100000', r, "seventh character must be rune '\U00100000'")
}

func TestLexerZeroLen(t *testing.T) {
	z := NewLexer(test.NewPlainReader(bytes.NewBufferString("")))
	assert.Equal(t, byte(0), z.Peek(0), "first character must yield error")
}

func TestLexerEmptyReader(t *testing.T) {
	z := NewLexer(test.NewEmptyReader())
	assert.Equal(t, byte(0), z.Peek(0), "first character must yield error")
	assert.Equal(t, io.EOF, z.Err(), "error must be EOF")
	assert.Equal(t, byte(0), z.Peek(0), "second peek must also yield error")
}
