package addresscheck

import (
	"fmt"
	"regexp"
	"strings"
)

var coinCheckMap = map[string]func(string) error{
	"ironfish": func(address string) error {
		const ironfishAddrLen = 64
		if len(address) != ironfishAddrLen {
			return fmt.Errorf("invalid address")
		}
		if !isHexString(address) {
			return fmt.Errorf("invalid address")
		}
		return nil
	},
}

func getCoinName(targetCoinName string) string {
	for coinName := range coinCheckMap {
		contains := strings.Contains(targetCoinName, coinName)
		if contains {
			return coinName
		}
	}
	return ""
}

func isHexString(str string) bool {
	match, _ := regexp.MatchString("^[0-9a-fA-F]+$", str)
	return match
}

func ValidateAddress(targetCoinName, address string) error {
	coinName := getCoinName(targetCoinName)
	return coinCheckMap[coinName](address)
}
