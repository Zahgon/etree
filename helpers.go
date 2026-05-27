// Copyright 2015-2019 Brett Vickers.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package etree

import (
	"io"
)

type stack[E any] struct {
	data []E
}

func (s *stack[E]) empty() bool { _ = "STUB: not implemented"; return false }

func (s *stack[E]) push(value E) { _ = "STUB: not implemented"; return }

func (s *stack[E]) pop() E { _ = "STUB: not implemented"; return *new(E) }

func (s *stack[E]) peek() E { _ = "STUB: not implemented"; return *new(E) }

type queue[E any] struct {
	data       []E
	head, tail int
}

func (f *queue[E]) add(value E) { _ = "STUB: not implemented"; return }

func (f *queue[E]) remove() E { _ = "STUB: not implemented"; return *new(E) }

func (f *queue[E]) len() int { _ = "STUB: not implemented"; return 0 }

func (f *queue[E]) grow() { _ = "STUB: not implemented"; return }

// xmlReader provides the interface by which an XML byte stream is
// processed and decoded.
type xmlReader interface {
	Bytes() int64
	Read(p []byte) (n int, err error)
}

// xmlSimpleReader implements a proxy reader that counts the number of
// bytes read from its encapsulated reader.
type xmlSimpleReader struct {
	r     io.Reader
	bytes int64
}

func newXmlSimpleReader(r io.Reader) xmlReader { _ = "STUB: not implemented"; return *new(xmlReader) }

func (xr *xmlSimpleReader) Bytes() int64 { _ = "STUB: not implemented"; return 0 }

func (xr *xmlSimpleReader) Read(p []byte) (n int, err error) {
	_ = "STUB: not implemented"
	return 0, nil
}

// xmlPeekReader implements a proxy reader that counts the number of
// bytes read from its encapsulated reader. It also allows the caller to
// "peek" at the previous portions of the buffer after they have been
// parsed.
type xmlPeekReader struct {
	r          io.Reader
	bytes      int64  // total bytes read by the Read function
	buf        []byte // internal read buffer
	bufSize    int    // total bytes used in the read buffer
	bufOffset  int64  // total bytes read when buf was last filled
	window     []byte // current read buffer window
	peekBuf    []byte // buffer used to store data to be peeked at later
	peekOffset int64  // total read offset of the start of the peek buffer
}

func newXmlPeekReader(r io.Reader) *xmlPeekReader { _ = "STUB: not implemented"; return nil }

func (xr *xmlPeekReader) Bytes() int64 { _ = "STUB: not implemented"; return 0 }

func (xr *xmlPeekReader) Read(p []byte) (n int, err error) {
	_ = "STUB: not implemented"
	return 0, nil
}

func (xr *xmlPeekReader) PeekPrepare(offset int64, maxLen int) { _ = "STUB: not implemented"; return }

func (xr *xmlPeekReader) PeekFinalize() []byte { _ = "STUB: not implemented"; return nil }

func (xr *xmlPeekReader) fill() error { _ = "STUB: not implemented"; return nil }

func (xr *xmlPeekReader) updatePeekBuf() { _ = "STUB: not implemented"; return }

// xmlWriter implements a proxy writer that counts the number of
// bytes written by its encapsulated writer.
type xmlWriter struct {
	w     io.Writer
	bytes int64
}

func newXmlWriter(w io.Writer) *xmlWriter { _ = "STUB: not implemented"; return nil }

func (xw *xmlWriter) Write(p []byte) (n int, err error) { _ = "STUB: not implemented"; return 0, nil }

// isWhitespace returns true if the byte slice contains only
// whitespace characters.
func isWhitespace(s string) bool { _ = "STUB: not implemented"; return false }

// spaceMatch returns true if namespace a is the empty string
// or if namespace a equals namespace b.
func spaceMatch(a, b string) bool { _ = "STUB: not implemented"; return false }

// spaceDecompose breaks a namespace:tag identifier at the ':'
// and returns the two parts.
func spaceDecompose(str string) (space, key string) { _ = "STUB: not implemented"; return "", "" }

// Strings used by indentCRLF and indentLF
const (
	indentSpaces = "\r\n                                                                "
	indentTabs   = "\r\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
)

// indentCRLF returns a CRLF newline followed by n copies of the first
// non-CRLF character in the source string.
func indentCRLF(n int, source string) string { _ = "STUB: not implemented"; return "" }

// indentLF returns a LF newline followed by n copies of the first non-LF
// character in the source string.
func indentLF(n int, source string) string { _ = "STUB: not implemented"; return "" }

// nextIndex returns the index of the next occurrence of byte ch in s,
// starting from offset.  It returns -1 if the byte is not found.
func nextIndex(s string, ch byte, offset int) int { _ = "STUB: not implemented"; return 0 }

// isInteger returns true if the string s contains an integer.
func isInteger(s string) bool { _ = "STUB: not implemented"; return false }

type escapeMode byte

const (
	escapeNormal escapeMode = iota
	escapeCanonicalText
	escapeCanonicalAttr
)

// escapeString writes an escaped version of a string to the writer.
func escapeString(w Writer, s string, m escapeMode) { _ = "STUB: not implemented"; return }

func isInCharacterRange(r rune) bool { _ = "STUB: not implemented"; return false }
