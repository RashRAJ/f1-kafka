package kafka

import (
	"strconv"
	"strings"
)

func ParseDriverNumber(driverNumberStr string) *int {
	if driverNumberStr != "" && driverNumberStr != "0" {
		if driverNum, err := strconv.Atoi(driverNumberStr); err == nil && driverNum != 0 {
			return &driverNum
		}
	}
	return nil
}

func ParseYears(yearsStr string) []int {
	parts := strings.Split(yearsStr, ",")
	years := make([]int, 0, len(parts))
	for _, part := range parts {
		if year, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			years = append(years, year)
		}
	}
	return years
}
