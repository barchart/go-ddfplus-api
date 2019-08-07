package ddf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SymbolType int

const (
	Unknown      SymbolType = 0
	Future       SymbolType = 1
	FutureOption SymbolType = 2
)

type Symbol struct {
	Type    SymbolType
	Root    string
	Month   string
	Year    int
	Strike  int
	CallPut string
}

func ParseSymbol(s string) (Symbol, error) {
	var (
		symbol Symbol
		err    error
	)

	if len(s) < 1 {
		return symbol, fmt.Errorf("Invalid symbol length (0)")
	}

	switch ch := s[len(s)-1:][0]; {
	case ch >= '0' && ch <= '9':
		symbol.Type = Future
		yr := ""
	Loop:
		for i := len(s) - 1; i > 0; i-- {
			switch ch2 := s[i : i+1][0]; {
			case ch2 >= '0' && ch2 <= '9':
				yr = string(ch2) + yr
			default:
				symbol.Year, _ = strconv.Atoi(yr)
				symbol.Month = string(ch2)
				symbol.Root = strings.TrimSpace(s[0:i])
				break Loop // regular break jumps out of the switch
			}
		}
	default:
		// Futures option, and others ....
		var reFutOpt = regexp.MustCompile(`(?i)^(.{1,2})([A-Z])([0-9]{1,4})([A-Z])$`)
		arr := reFutOpt.FindStringSubmatch(s)
		if len(arr) == 5 { // original + 4 parts
			symbol.Type = FutureOption
			symbol.Root = arr[1]
			symbol.Month = arr[2]
			symbol.Strike, _ = strconv.Atoi(arr[3])
			symbol.CallPut = arr[4]
		} else {
			fmt.Println("???", s)
		}
	}

	return symbol, err
}
