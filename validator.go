// validator.go
package main

import (
    "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
    "crypto/rand"

    "fmt"
    "mime/multipart"
    "regexp"
    "strconv"
    "strings"
    "io"
)

// ValidateCreditCard checks if the provided credit card number is valid using the Luhn algorithm.
func ValidateCreditCard(cardNumber string) bool {
    // Remove spaces and dashes
    cardNumber = strings.ReplaceAll(cardNumber, " ", "")
    cardNumber = strings.ReplaceAll(cardNumber, "-", "")

    // Check if the card number contains only digits
    if !isNumeric(cardNumber) {
        return false
    }

    // Implementing Luhn Algorithm
    sum := 0
    alternate := false

    for i := len(cardNumber) - 1; i >= 0; i-- {
        n, _ := strconv.Atoi(string(cardNumber[i]))

        if alternate {
            n *= 2
            if n > 9 {
                n -= 9
            }
        }

        sum += n
        alternate = !alternate
    }

    return sum%10 == 0
}

// isNumeric checks if the input string contains only numeric characters.
func isNumeric(s string) bool {
    return regexp.MustCompile(`^[0-9]+$`).MatchString(s)
}

// GetCardType returns the type of credit card based on its number.
func GetCardType(cardNumber string) string {
    if !isNumeric(cardNumber) {
        return "Invalid"
    }

    switch {
    case regexp.MustCompile(`^4[0-9]{12}(?:[0-9]{3})?$`).MatchString(cardNumber):
        return `<strong>Visa</strong> <img src="/assets/visa_logo.png" alt="Visa Logo" style="height: 20px; vertical-align: middle;"/>`
    case regexp.MustCompile(`^5[1-5][0-9]{14}$`).MatchString(cardNumber):
        return `<strong>MasterCard</strong> <img src="/assets/mastercard_logo.png" alt="MasterCard Logo" style="height: 20px; vertical-align: middle;"/>`
    case regexp.MustCompile(`^3[47][0-9]{13}$`).MatchString(cardNumber):
        return `<strong>American Express</strong> <img src="/assets/american_express_logo.jpeg" alt="American Express Logo" style="height: 20px; vertical-align: middle;"/>`
    case regexp.MustCompile(`^6(?:011|5[0-9]{2})[0-9]{12}$`).MatchString(cardNumber):
        return `<strong>Discover</strong> <img src="/assets/discover_logo.png" alt="Discover Logo" style="height: 20px; vertical-align: middle;"/>`
    default:
        return "Unknown"
    }
}

// isValidImageType checks if the uploaded file is of a valid image type.
func isValidImageType(file *multipart.FileHeader) bool {
    allowedTypes := map[string]bool{
        "image/png":  true,
        "image/jpeg": true,
        "image/jpg": true,
    }

    return allowedTypes[file.Header.Get("Content-Type")]
}

// Encrypt encrypts plaintext using AES-GCM and returns a base64-encoded string
func Encrypt(plaintext []byte, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", fmt.Errorf("failed to generate nonce: %w", err)
    }

    // Seal encrypts the plaintext and appends the nonce at the front
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    // Return base64 encoded string
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertextBytes []byte, key []byte) (string, error) {
    // Create a new AES cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    // Create GCM (Galois/Counter Mode)
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertextBytes) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    // Extract nonce and actual ciphertext
    nonce, ciphertext := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]

    // Decrypting
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

    if err != nil {
        return "", fmt.Errorf("decryption failed: %w", err)
    }

    return string(plaintext), nil
}

