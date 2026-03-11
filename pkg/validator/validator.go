package validator

import (
	"regexp"
	"strconv"
)

func ValidateCPF(cpf string) bool {
	cpf = regexp.MustCompile(`[^0-9]`).ReplaceAllString(cpf, "")
	if len(cpf) != 11 {
		return false
	}

	if cpf == "00000000000" || cpf == "11111111111" || cpf == "22222222222" ||
		cpf == "33333333333" || cpf == "44444444444" || cpf == "55555555555" ||
		cpf == "66666666666" || cpf == "77777777777" || cpf == "88888888888" ||
		cpf == "99999999999" {
		return false
	}

	var sum int
	for i := 0; i < 9; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (10 - i)
	}
	digit1 := (sum * 10) % 11
	if digit1 == 10 {
		digit1 = 0
	}

	sum = 0
	for i := 0; i < 10; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (11 - i)
	}
	digit2 := (sum * 10) % 11
	if digit2 == 10 {
		digit2 = 0
	}

	return strconv.Itoa(digit1) == string(cpf[9]) && strconv.Itoa(digit2) == string(cpf[10])
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ValidatePhone(phone string) bool {
	phone = regexp.MustCompile(`[^0-9]`).ReplaceAllString(phone, "")
	return len(phone) >= 10 && len(phone) <= 11
}
