package str

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// EmptyString ...
func EmptyString(text string) *string {
	if text == "" {
		return nil
	}
	return &text
}

// EmptyInt ...
func EmptyInt(number int) *int {
	if number == 0 {
		return nil
	}
	return &number
}

func StringToInt(data string) int {
	res, err := strconv.Atoi(data)
	if err != nil {
		res = 0
	}

	return res
}

func StringToBool(data string) bool {
	res, err := strconv.ParseBool(data)
	if err != nil {
		res = false
	}

	return res
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()

	return str
}

func IsActive(data string) *string {

	var isActive string
	res, err := strconv.ParseBool(data)

	if err != nil {
		isActive = ""
		return &isActive
	}

	if res {
		isActive = "and is_active = 'true'"
	} else {
		isActive = "and is_active = 'false'"
	}

	return &isActive
}

// Unique ...
func Unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
