package addresscheck

import (
	"encoding/hex"
	"fmt"
	"strings"
)

var coinCheckMap = map[string]func(string) error{
	"ironfish": func(address string) error {
		const ironfishAddrLen = 64
		if len(address) != ironfishAddrLen {
			return fmt.Errorf("invalid address")
		}
		if _, err := hex.DecodeString(address); err != nil {
			return err
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

func ValidateAddress(targetCoinName, address string) error {
	coinName := getCoinName(targetCoinName)
	if coinName == "" {
		return nil
	}
	return coinCheckMap[coinName](address)
}
