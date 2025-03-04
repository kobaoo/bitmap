package tools

// Convert integer to hexadecimal string (improved)
func toHex(n int) string {
	hexChars := "0123456789abcdef"
	if n == 0 {
		return "00" // Handle zero case
	}
	var res string
	for n > 0 {
		res = string(hexChars[n%16]) + res
		n = n / 16
	}
	if len(res) == 1 {
		res = "0" + res
	}
	return res
}

// Convert hexadecimal string to integer (decimal)
func toDec(hexStr string) uint32 {
	// Hexadecimal digits and their corresponding decimal values
	hexChars := "0123456789abcdef"
	hexMap := make(map[rune]int)
	for i, c := range hexChars {
		hexMap[c] = i
	}

	// Convert the hex string to decimal
	var result int
	for _, c := range hexStr {
		// Update the result with each character's corresponding decimal value
		result = result*16 + hexMap[c]
	}

	return uint32(result)
}
