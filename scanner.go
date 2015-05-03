package eval

import (
	"unicode/utf8"
)

type scanner struct {
	src []byte // source
	// current char and offset; nxOffset = offset + len of char
	ch       rune // current character
	offset   int  // character offset
	nxOffset int  // next offset
	// line and char counters for error location
	linePos int // current line for error reporting
	charPos int // current character for error reporting
}

func newScanner(src []byte) scanner {
	s := scanner{src: src, ch: ' ', offset: 0, nxOffset: 0, linePos: 1, charPos: 0}
	s.next()
	return s
}

func (s *scanner) next() {
	if s.nxOffset < len(s.src) {
		s.offset = s.nxOffset
		if s.ch == '\n' {
			s.linePos += 1
			s.charPos = 0
		}
		r, w := utf8.DecodeRune(s.src[s.nxOffset:])
		if r == utf8.RuneError && w == 1 {
			panic(s.newError("illegal utf-8 character"))
		}
		s.nxOffset += w
		s.ch = r
		s.charPos += 1
	} else {
		s.offset = len(s.src)
		s.ch = -1 // end of expression
	}
}

func (s *scanner) newError(msg string) ScannerError {
	return ScannerError{message: msg, line: s.linePos, position: s.charPos}
}

func (s *scanner) scan() (token, string) {
	s.skipWhitespace()
	switch {
	case s.ch == '{':
		t, l := s.scanCurlyIdentifier()
		return t, l
	case isLetter(s.ch):
		t, l := s.scanIdentifier()
		return t, l
	case '0' <= s.ch && s.ch <= '9':
		return s.scanNumber()
	case s.ch == -1:
		return EOE, ""
	case s.ch == '"' || s.ch == '\'':
		return s.scanString()
	default:
		ch := s.ch
		s.next() // look ahead
		for _, ops := range tokens {
			if ch == ops.seq[0] && (len(ops.seq) == 1 || s.ch == ops.seq[1]) {
				if len(ops.seq) == 2 {
					s.next()
				}
				return ops.tok, ""
			}
		}
	}
	return BAD, string(s.ch)
}

func (s *scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (s *scanner) scanIdentifier() (token, string) {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) || s.ch == '_' || s.ch == '.' {
		s.next()
	}
	return IDENT, string(s.src[offs:s.offset])
}

func (s *scanner) scanCurlyIdentifier() (token, string) {
	s.next() // eat opening {
	offs := s.offset
	for prev := ' '; s.ch != '}' && prev != '\\'; s.next() {
		prev = s.ch
	}
	s.next() // eat closing
	return IDENT, string(s.src[offs : s.offset-1])
}

func (s *scanner) scanNumber() (token, string) {
	offs := s.offset
	for isDigit(s.ch) {
		s.next()
	}
	if s.ch == '.' {
		s.next()
		for isDigit(s.ch) {
			s.next()
		}
	}
	if s.ch == 'e' || s.ch == 'E' {
		s.next()
		if s.ch == '-' || s.ch == '+' {
			s.next()
		}
		for isDigit(s.ch) {
			s.next()
		}
	}
	return NUMBER, string(s.src[offs:s.offset])
}

func (s *scanner) scanString() (token, string) {
	delim := s.ch
	s.next() // eat delim
	offs := s.offset
	for s.ch != delim {
		if s.ch == '\\' {
			s.next()
		}
		s.next()
	}
	s.next() // eat delim
	return STRING, string(s.src[offs : s.offset-1])
}
