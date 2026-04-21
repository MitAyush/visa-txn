// Package config loads application settings from the environment.
package config

import "os"

type Sign string

const (
	SignDebit  Sign = "debit"
	SignCredit Sign = "credit"
)

type OperationType struct {
	Description string
	Sign        Sign
}

type Config struct {
	Port           string
	DBUrl          string
	OperationTypes map[int]OperationType
}

func Load() *Config {
	return &Config{
		Port:  getEnv("PORT", "8080"),
		DBUrl: getEnv("DB_URL", "storage/visa.db"),
		OperationTypes: map[int]OperationType{
			1: {Description: "Normal Purchase", Sign: SignDebit},
			2: {Description: "Purchase with installments", Sign: SignDebit},
			3: {Description: "Withdrawal", Sign: SignDebit},
			4: {Description: "Credit Voucher", Sign: SignCredit},
		},
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
