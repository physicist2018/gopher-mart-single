package luhn

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	ErrLuhnCheck             = errors.New("invalid number")
	ErrLunhCheckHasNotDigits = errors.New("number contains non-digit characters")
)

// luhnCheckDigit вычисляет контрольную цифру для переданной последовательности чисел
func checkDigit(numbers []int) int {
	sum := 0
	alt := false

	// Обрабатываем числа справа налево
	for i := len(numbers) - 1; i >= 0; i-- {
		n := numbers[i]
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}

	// Вычисляем контрольную цифру
	return (10 - (sum % 10)) % 10
}

// generateLuhnNumber генерирует случайное число с проверочной цифрой Луна и возвращает строковое представление
func GenerateLuhnNumber(length int) string {
	if length < 2 {
		return ""
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Генерируем случайную последовательность длины length - 1
	numbers := make([]int, length-1)
	for i := 0; i < length-1; i++ {
		numbers[i] = rnd.Intn(10)
	}

	// Вычисляем и добавляем контрольную цифру
	checkDigit := checkDigit(numbers)
	numbers = append(numbers, checkDigit)

	// Преобразуем числа в строковое представление
	var sb strings.Builder
	for _, n := range numbers {
		sb.WriteString(strconv.Itoa(n))
	}

	return sb.String()
}

// validateLuhnNumber проверяет, корректно ли число по алгоритму Луна
func ValidateLuhnNumber(number string) bool {
	// Преобразуем строку в срез целых чисел
	numbers := make([]int, len(number))
	for i, char := range number {
		num, err := strconv.Atoi(string(char))
		if err != nil {
			return false // Невалидный символ в строке
		}
		numbers[i] = num
	}

	// Проверяем, соответствует ли последняя цифра контрольной цифре
	cDigit := numbers[len(numbers)-1]
	return cDigit == checkDigit(numbers[:len(numbers)-1])
}
