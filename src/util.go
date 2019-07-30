// Go ddfplus API
//
// Copyright 2019 Barchart.com, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the GNU license
// available at https://github.com/barchart/go-ddfplus-api/blob/master/LICENSE.
package ddf

import (
	"fmt"
	"strconv"
)

func ParseFloat(s string, bc string) (float64, error) {
	if len(s) == 0 {
		return 0.0, fmt.Errorf("zero length string")
	}

	var sign = 1.0
	if s[0] == '-' {
		s = s[1:]
		sign = -1.0
	}

	switch bc {
	case "2": // 8ths

		n, err := strconv.ParseFloat(s[0:len(s)-1], 64)
		if err != nil {
			return 0.0, err
		}

		d, err := strconv.ParseFloat(s[len(s)-1:], 64)
		if err != nil {
			return 0.0, err
		}

		return sign * (n + d/8), nil

	case "3": // 16ths
		if len(s) < 2 {
			return 0.0, fmt.Errorf("Invalid length %d", len(s))
		}

		n, err := strconv.ParseFloat(s[0:len(s)-2], 64)
		if err != nil {
			return 0.0, err
		}

		d, err := strconv.ParseFloat(s[len(s)-2:], 64)
		if err != nil {
			return 0.0, err
		}

		return sign * (n + d/16), nil

	case "4": // 32nds
		if len(s) < 2 {
			return 0.0, fmt.Errorf("Invalid length %d", len(s))
		}

		n, err := strconv.ParseFloat(s[0:len(s)-2], 64)
		if err != nil {
			return 0.0, err
		}

		d, err := strconv.ParseFloat(s[len(s)-2:], 64)
		if err != nil {
			return 0.0, err
		}

		return sign * (n + d/32), nil

	case "5": // 64th
		if len(s) < 2 {
			return 0.0, fmt.Errorf("Invalid length %d", len(s))
		}

		n, err := strconv.ParseFloat(s[0:len(s)-2], 64)
		if err != nil {
			return 0.0, err
		}

		d, err := strconv.ParseFloat(s[len(s)-2:], 64)
		if err != nil {
			return 0.0, err
		}

		return sign * (n + d/64), nil

	case "6": // 128ths
		if len(s) < 3 {
			return 0.0, fmt.Errorf("Invalid length %d", len(s))
		}
		n, err := strconv.ParseFloat(s[0:len(s)-3], 64)
		if err != nil {
			return 0.0, err
		}

		d, err := strconv.ParseFloat(s[len(s)-3:], 64)
		if err != nil {
			return 0.0, err
		}

		return sign * (n + d/128), nil

	case "7": // 256ths
		if len(s) < 3 {
			return 0.0, fmt.Errorf("Invalid length %d", len(s))
		}

		n, err := strconv.ParseFloat(s[0:len(s)-3], 64)
		if err != nil {
			return 0.0, err
		}

		d, err := strconv.ParseFloat(s[len(s)-3:], 64)
		if err != nil {
			return 0.0, err
		}

		return sign * (n + d/256), nil

	case "8": // 0s
		f, err := strconv.Atoi(s)
		if err != nil {
			return 0.0, err
		}
		return sign * float64(f) / 1.0, nil

	case "9": // 1/10
		f, err := strconv.Atoi(s)
		if err != nil {
			return 0.0, err
		}

		return sign * float64(f) / 10.0, nil

	case "A": // 1/100
		f, err := strconv.Atoi(s)
		if err != nil {
			return 0.0, err
		}

		return sign * float64(f) / 100.0, nil

	case "B": // 1/1000
		f, err := strconv.Atoi(s)
		if err != nil {
			return 0.0, err
		}

		return sign * float64(f) / 1000.0, nil

	case "C":
		f, err := strconv.Atoi(s)
		if err != nil {
			return 0.0, err
		}

		return sign * float64(f) / 10000.0, nil

	case "D":
		f, err := strconv.Atoi(s)
		if err != nil {
			return 0.0, err
		}

		return sign * float64(f) / 100000.0, nil

	case "E":
		f, err := strconv.Atoi(s)
		if err != nil {
			return 0.0, err
		}

		return sign * float64(f) / 1000000.0, nil

	case "F":
		f, err := strconv.Atoi(s)
		if err != nil {
			return 0.0, err
		}

		return sign * float64(f) / 10000000.0, nil
	}

	return 0.0, nil
}
