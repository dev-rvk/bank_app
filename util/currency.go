package util

const (
 USD = "USD"
 INR = "INR"
 CAD = "CAD"
 YEN = "YEN"
)

func IsSuppourtedCurrency (currency string) bool {
	switch currency {
	case USD, INR, CAD, YEN:
		return true
	}
	return false
}