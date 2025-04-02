package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FormatBalance(balance uint64) string {
	return fmt.Sprintf("%.2f", float64(balance)/100)
}

var amountPattern = regexp.MustCompile(`^\d+\.\d{2}$`)

func ParseAmount(amount string) (uint64, error) {
	amount = strings.TrimSpace(amount)

	if !amountPattern.MatchString(amount) {
		return 0, fmt.Errorf("amount must be in format '0.00' with exactly 2 decimal places")
	}

	parts := strings.Split(amount, ".")
	dollars, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount format: %s", amount)
	}

	cents, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid cents format: %s", amount)
	}

	return dollars*100 + cents, nil
}
