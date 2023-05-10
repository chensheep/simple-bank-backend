package util

const (
	USD = "USD"
	EUR = "EUR"
	TWD = "TWD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, TWD:
		return true
	}
	return false
}
