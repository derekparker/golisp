// Copyright 2014 SteelSeries ApS.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements a basic LISP interpretor for embedding in a go program for scripting.
// This file implements the tokenizer.

package golisp

import (
	"fmt"
	"unicode"
)

const (
	ILLEGAL = iota
	SYMBOL
	NUMBER
	HEXNUMBER
	FLOAT
	STRING
	QUOTE
	BACKQUOTE
	COMMA
	COMMAAT
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	PERIOD
	TRUE
	FALSE
	COMMENT
	EOF
)

type Tokenizer struct {
	LookaheadToken int
	LookaheadLit   string
	Source         string
	Position       int
}

func NewTokenizer(src string) *Tokenizer {
	t := &Tokenizer{Source: src}
	t.ConsumeToken()
	return t
}

func (self *Tokenizer) NextToken() (token int, lit string) {
	return self.LookaheadToken, self.LookaheadLit
}

func (self *Tokenizer) isSymbolCharacter(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '*' || ch == '-' || ch == '?' || ch == '!' || ch == '_' || ch == '>'
}

func (self *Tokenizer) readSymbol() (token int, lit string) {
	start := self.Position
	for !self.isEof() && self.isSymbolCharacter(rune(self.Source[self.Position])) {
		self.Position++
	}
	return SYMBOL, self.Source[start:self.Position]
}

func isHexChar(ch rune) bool {
	switch ch {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F':
		return true
	default:
		return false
	}
}

func (self *Tokenizer) readNumber() (token int, lit string) {
	start := self.Position
	isHex := false
	isFloat := false
	sawDecimal := false
	for !self.isEof() {
		ch := rune(self.Source[self.Position])
		if ch == '.' && !sawDecimal {
			isFloat = true
			sawDecimal = true
			self.Position++
		} else if (start == self.Position) && ch == '-' {
			self.Position++
		} else if unicode.IsNumber(ch) {
			self.Position++
		} else if (start == self.Position-1) && ch == 'x' {
			isHex = true
			self.Position++
		} else if isHex && isHexChar(ch) {
			self.Position++
		} else {
			break
		}
	}

	lit = self.Source[start:self.Position]
	if isHex {
		token = HEXNUMBER
	} else if isFloat {
		token = FLOAT
	} else {
		token = NUMBER
	}
	return
}

func (self *Tokenizer) readString() (token int, lit string) {
	buffer := make([]rune, 0, 10)
	self.Position++
	for !self.isEof() && rune(self.Source[self.Position]) != '"' {
		if rune(self.Source[self.Position]) == '\\' {
			self.Position++
		}
		buffer = append(buffer, rune(self.Source[self.Position]))
		self.Position++
	}
	if self.isEof() {
		return EOF, ""
	}
	self.Position++
	return STRING, string(buffer)
}

func (self *Tokenizer) isEof() bool {
	return self.Position >= len(self.Source)
}

func (self *Tokenizer) isAlmostEof() bool {
	return self.Position == len(self.Source)-1
}

func (self *Tokenizer) readNextToken() (token int, lit string) {
	if self.isEof() {
		return EOF, ""
	}
	for unicode.IsSpace(rune(self.Source[self.Position])) {
		self.Position++
		if self.isEof() {
			return EOF, ""
		}
	}
	currentChar := rune(self.Source[self.Position])
	var nextChar rune
	if !self.isAlmostEof() {
		nextChar = rune(self.Source[self.Position+1])
	}
	if unicode.IsLetter(currentChar) || currentChar == '_' {
		return self.readSymbol()
	} else if unicode.IsNumber(currentChar) {
		return self.readNumber()
	} else if currentChar == '-' && unicode.IsNumber(nextChar) {
		return self.readNumber()
	} else if currentChar == '"' {
		return self.readString()
	} else if currentChar == '\'' {
		self.Position++
		return QUOTE, "'"
	} else if currentChar == '`' {
		self.Position++
		return BACKQUOTE, "`"
	} else if currentChar == ',' && nextChar == '@' {
		self.Position += 2
		return COMMAAT, ",@"
	} else if currentChar == ',' {
		self.Position++
		return COMMA, ","
	} else if currentChar == '(' {
		self.Position++
		return LPAREN, "("
	} else if currentChar == ')' {
		self.Position++
		return RPAREN, ")"
	} else if currentChar == '[' {
		self.Position++
		return LBRACKET, "["
	} else if currentChar == ']' {
		self.Position++
		return RBRACKET, "]"
	} else if currentChar == '.' {
		self.Position++
		return PERIOD, "."
	} else if currentChar == '-' && nextChar == '>' {
		self.Position += 2
		return SYMBOL, "->"
	} else if currentChar == '=' && nextChar == '>' {
		self.Position += 2
		return SYMBOL, "=>"
	} else if currentChar == '+' {
		self.Position++
		return SYMBOL, "+"
	} else if currentChar == '-' {
		self.Position++
		return SYMBOL, "-"
	} else if currentChar == '*' {
		self.Position++
		return SYMBOL, "*"
	} else if currentChar == '/' {
		self.Position++
		return SYMBOL, "/"
	} else if currentChar == '%' {
		self.Position++
		return SYMBOL, "%"
	} else if currentChar == '<' && nextChar == '=' {
		self.Position += 2
		return SYMBOL, "<="
	} else if currentChar == '<' {
		self.Position++
		return SYMBOL, "<"
	} else if currentChar == '>' && nextChar == '=' {
		self.Position += 2
		return SYMBOL, ">="
	} else if currentChar == '>' {
		self.Position++
		return SYMBOL, ">"
	} else if currentChar == '=' && nextChar == '=' {
		self.Position += 2
		return SYMBOL, "=="
	} else if currentChar == '=' {
		self.Position++
		return SYMBOL, "="
	} else if currentChar == '!' && nextChar == '=' {
		self.Position += 2
		return SYMBOL, "!="
	} else if currentChar == '!' {
		self.Position++
		return SYMBOL, "!"
	} else if currentChar == '#' {
		self.Position += 2
		if nextChar == 't' {
			return TRUE, "#t"
		} else {
			return FALSE, "#f"
		}
	} else if currentChar == ';' {
		start := self.Position
		for {
			if self.isEof() {
				return COMMENT, self.Source[start:]
			} else if self.Source[self.Position] == '\n' {
				return COMMENT, self.Source[start:self.Position]
			}
			self.Position++
		}
	} else {
		return ILLEGAL, fmt.Sprintf("%c", currentChar)
	}
}

func (self *Tokenizer) ConsumeToken() {
	self.LookaheadToken, self.LookaheadLit = self.readNextToken()
	if self.LookaheadToken == COMMENT { // skip comments
		self.ConsumeToken()
	}
}
