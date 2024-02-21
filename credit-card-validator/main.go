package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/validate", validateHandler)
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		CreditCardNumber string `json:"credit_card_number"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	cardNumber := strings.Replace(requestData.CreditCardNumber, " ", "", -1)
	isValid := validateCreditCard(cardNumber)

	cardNetwork := identifyCardNetwork(cardNumber)
	responseData := struct {
		IsValid     bool   `json:"is_valid"`
		CardNetwork string `json:"card_network,omitempty"`
	}{
		IsValid:     isValid,
		CardNetwork: cardNetwork,
	}

	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func validateCreditCard(cardNumber string) bool {
	cardNumber = strings.Replace(cardNumber, " ", "", -1)
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}
	sum := 0
	isSecond := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(cardNumber[i]))
		if err != nil {
			return false
		}
		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isSecond = !isSecond
	}
	return sum%10 == 0
}

func identifyCardNetwork(cardNumber string) string {
	visaPattern := regexp.MustCompile("^4[0-9]{12}(?:[0-9]{3})?$")
	mastercardPattern := regexp.MustCompile("^5[1-5][0-9]{14}$")
	americanExpressPattern := regexp.MustCompile("^3[47][0-9]{13}$")

	if visaPattern.MatchString(cardNumber) {
		return "Visa"
	} else if mastercardPattern.MatchString(cardNumber) {
		return "MasterCard"
	} else if americanExpressPattern.MatchString(cardNumber) {
		return "American Express"
	}

	return "Unknown"
}
