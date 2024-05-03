package main

import (
	"unicode/utf8"

	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/strutil"
)

func makeMove(m func(string, int) int) func(etk.CodeBuffer) etk.CodeBuffer {
	return func(buf etk.CodeBuffer) etk.CodeBuffer {
		buf.Dot = m(buf.Content, buf.Dot)
		return buf
	}
}

func moveDotLeft(buffer string, dot int) int {
	_, w := utf8.DecodeLastRuneInString(buffer[:dot])
	return dot - w
}

func moveDotRight(buffer string, dot int) int {
	_, w := utf8.DecodeRuneInString(buffer[dot:])
	return dot + w
}

func moveDotSOL(buffer string, dot int) int {
	return strutil.FindLastSOL(buffer[:dot])
}

func moveDotEOL(buffer string, dot int) int {
	return strutil.FindFirstEOL(buffer[dot:]) + dot
}
