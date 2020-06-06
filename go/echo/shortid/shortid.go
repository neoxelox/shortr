package shortid

import (
	"fmt"
	"math"
	"regexp"
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
	if matched, _ := regexp.MatchString("[0-9]", chr); matched {
		return int([]rune(chr)[0] - 48), nil
	} else if matched, _ := regexp.MatchString("[A-Z]", chr); matched {
		return int([]rune(chr)[0] - 55), nil
	} else if matched, _ := regexp.MatchString("[a-z]", chr); matched {
		return int([]rune(chr)[0] - 61), nil
	}
	return -1, fmt.Errorf("%s is not valid character", chr)
}

func dehydrate(dgt int) (string, error) {
	switch {
	case dgt < 10:
		return string(dgt + 48), nil
	case dgt <= 35:
		return string(dgt + 55), nil
	case dgt < 62:
		return string(dgt + 61), nil
	default:
		return "", fmt.Errorf("%d is not a valid integer", dgt)
	}
}
