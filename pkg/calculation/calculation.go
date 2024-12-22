package calculation

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func Calc(expression string) (float64, error) {
	if len(expression) == 0 {
		return 0, errors.New("неизвестный символ")
	}

	openCount := 0
	closeCount := 0
	for _, r := range expression {
		if r == '(' {
			openCount++
		} else if r == ')' {
			closeCount++
		}
	}
	if openCount != closeCount {
		return 0, errors.New("ошибочка в количестве скобок")
	}

	for i := 0; i < len(expression); i++ {
		if expression[i] == 'e' || expression[i] == 'E' {
			return 0, errors.New("переполнение")
		}
	}

	for i := 0; i < len(expression); i++ {
		if expression[i] >= '0' && expression[i] <= '9' {
			j := i
			for j < len(expression) && (expression[j] >= '0' && expression[j] <= '9' || expression[j] == '.') {
				j++
			}
			if j-i > 15 {
				return 0, errors.New("переполнение")
			}
			i = j
		}
	}

	number, egor := podschet(expression)
	if egor != nil {
		return 0, egor
	}

	num, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return 0, errors.New("неизвестный символ")
	}

	if math.IsInf(num, 0) || math.IsNaN(num) {
		return 0, errors.New("переполнение")
	}

	return num, nil
}

func skobki(str string, ind int) (string, int) {
	run := []rune(str)
	sl := ""
	depth := 0

	for i := ind; i < len(run); i++ {
		if run[i] == 40 {
			depth++
			if depth == 1 {
				continue
			}
		}
		if run[i] == 41 {
			depth--
			if depth == 0 {
				result, err := podschet(sl)
				if err != nil {
					return "", i - ind
				}
				return result, i - ind
			}
		}
		if depth > 0 {
			sl += string(run[i])
		}
	}

	return sl, len(run) - ind
}

func podschet(str string) (string, error) {
	run := []rune(str)
	q := 0
	slice := make([]string, 1)
	col := 0

	for i := 0; i < len(run); i++ {
		if q > 0 {
			q--
			continue
		}
		if run[i] == 40 {
			result, offset := skobki(str, i)
			if result != "" {
				slice[col] += result
			}
			q += offset
			continue
		} else if run[i] >= 48 && run[i] <= 57 || run[i] == 46 {
			slice[col] += string(run[i])
			continue
		} else if run[i] == 43 || run[i] == 45 || run[i] == 42 || run[i] == 47 {
			if i == 0 || i == len(run)-1 {
				return "", errors.New("неизвестный символ")
			}
			slice = append(slice, string(run[i]))
			slice = append(slice, "")
			col += 2
			continue
		} else if run[i] == 41 {
			continue
		} else {
			return "", errors.New("неизвестный символ")
		}
	}

	for i := 0; i < len(slice); i++ {
		if slice[i] == "*" {
			num1, err1 := strconv.ParseFloat(slice[i-1], 64)
			num2, err2 := strconv.ParseFloat(slice[i+1], 64)
			if err1 != nil || err2 != nil {
				return "", errors.New("неизвестный символ")
			}
			result := num1 * num2
			if math.IsInf(result, 0) || math.IsNaN(result) {
				return "", errors.New("переполнение")
			}
			slice[i] = fmt.Sprint(result)
			slice = remove(slice, i+1)
			slice = remove(slice, i-1)
			i--
		} else if slice[i] == "/" {
			num1, err1 := strconv.ParseFloat(slice[i-1], 64)
			num2, err2 := strconv.ParseFloat(slice[i+1], 64)
			if err1 != nil || err2 != nil {
				return "", errors.New("неизвестный символ")
			}
			if num2 == 0 {
				return "", errors.New("на ноль делить нельзя")
			}
			result := num1 / num2
			if math.IsInf(result, 0) || math.IsNaN(result) {
				return "", errors.New("переполнение")
			}
			slice[i] = fmt.Sprint(result)
			slice = remove(slice, i+1)
			slice = remove(slice, i-1)
			i--
		}
	}

	for i := 0; i < len(slice); i++ {
		if slice[i] == "+" {
			num1, err1 := strconv.ParseFloat(slice[i-1], 64)
			num2, err2 := strconv.ParseFloat(slice[i+1], 64)
			if err1 != nil || err2 != nil {
				return "", errors.New("неизвестный символ")
			}
			result := num1 + num2
			if math.IsInf(result, 0) || math.IsNaN(result) {
				return "", errors.New("переполнение")
			}
			slice[i] = fmt.Sprint(result)
			slice = remove(slice, i+1)
			slice = remove(slice, i-1)
			i--
		} else if slice[i] == "-" {
			num1, err1 := strconv.ParseFloat(slice[i-1], 64)
			num2, err2 := strconv.ParseFloat(slice[i+1], 64)
			if err1 != nil || err2 != nil {
				return "", errors.New("неизвестный символ")
			}
			result := num1 - num2
			if math.IsInf(result, 0) || math.IsNaN(result) {
				return "", errors.New("переполнение")
			}
			slice[i] = fmt.Sprint(result)
			slice = remove(slice, i+1)
			slice = remove(slice, i-1)
			i--
		}
	}

	return slice[0], nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
