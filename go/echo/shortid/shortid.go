package shortid

import (
	"fmt"
	"math"
)

// Encode transforms a number to a compressed base62 string representation
func Encode(number int) (string, error) {
	if number == 0 {
		return "0", nil
	}

	str := ""
	for number > 0 {
		digit := number % 62
		chr, err := dehydrate(digit)
		if err != nil {
			return "", err
		}
		str = chr + str
		number = int(number / 62)
	}
	return str, nil
}

// Decode transforms a compressed base62 string representation to a number
func Decode(id string) (int, error) {
	num := 0
	lid := len(id) - 1
	for i := lid; i >= 0; i-- {
		dgt, err := hydrate(string(id[i]))
		if err != nil {
			return -1, err
		}
		num += dgt * int(math.Pow(62, float64(lid-i)))
	}
	return num, nil
}

func hydrate(chr string) (int, error) {
	switch {
	case chr >= "0" && chr <= "9":
		return int([]rune(chr)[0] - 48), nil
	case chr >= "A" && chr <= "Z":
		return int([]rune(chr)[0] - 55), nil
	case chr >= "a" && chr <= "z":
		return int([]rune(chr)[0] - 61), nil
	default:
		return -1, fmt.Errorf("%s is not valid character", chr)
	}
}

func dehydrate(dgt int) (string, error) {
	switch {
	case dgt < 10:
		return string(rune(dgt + 48)), nil
	case dgt <= 35:
		return string(rune(dgt + 55)), nil
	case dgt < 62:
		return string(rune(dgt + 61)), nil
	default:
		return "", fmt.Errorf("%d is not a valid integer", dgt)
	}
}
