// main.go
package main

import (
    "fmt"
    "net/http"
    "os"
    "io"
    "log"
    "regexp"
    "bytes"
    "encoding/json"

    "github.com/labstack/echo/v4"
    "github.com/otiai10/gosseract/v2"
    "github.com/joho/godotenv"  // For accessing environment variables
)

// DecryptRequest represents the request payload for decryption
type DecryptRequest struct {
    Ciphertext string `json:"ciphertext"`
}

// DecryptResponse represents the response payload for decryption
type DecryptResponse struct {
    DecryptedCardNumber string `json:"decryptedCardNumber"`
}

func main() {
    // Load .env file
    err := godotenv.Load()
    if err!=nil {
        log.Fatal("Error loading .env file")
    }

    key := os.Getenv("ENCRYPTION_KEY")  // Write your ENCRYPTION_KEY in ,env file(should be 16 characters for AES encryption)

    e := echo.New()

    // Serve static files from the "assets" directory
    e.Static("/assets", "assets")

    // Serve styles files from the "styles" directory
    e.Static("/styles","styles")

    // Serve script files from the "script" directory
    e.Static("/script", "script") 

    // Serve the index.html file
    e.GET("/", func(c echo.Context) error {
        return c.File("index.html")
    })

    // Handle image upload and text extraction
    e.POST("/upload", func(c echo.Context) error {
        file, err := c.FormFile("image")
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"message": "No image uploaded"})
        }

        fmt.Printf("Content Type: %s\n", file.Header.Get("Content-Type"))

        // Validate file type
        if !isValidImageType(file) {
            return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid file type. Only .png, .jpeg, or .jpg are allowed."})
        }

        src, err := file.Open()
        if err != nil {
            return c.String(http.StatusInternalServerError, "Failed to open image")
        }
        defer src.Close()

        // Create a temporary file
        tempFile, err := os.CreateTemp("", "*.jpeg") // Default to PNG; rename later if needed

        if err != nil {
            return c.String(http.StatusInternalServerError, "Failed to create temp file")
        }
        defer os.Remove(tempFile.Name()) // Clean up temp file after processing

        // Copy uploaded image to temporary file
        if _, err := io.Copy(tempFile, src); err != nil {
            return c.String(http.StatusInternalServerError, "Failed to save temp file")
        }

        // Use Tesseract to extract text from the image
        client := gosseract.NewClient()
        defer client.Close()

        client.SetImage(tempFile.Name())
        text, err := client.Text()

        if err != nil {
            return c.String(http.StatusInternalServerError, "Failed to extract text from image")
        }

        // Extract card number using regex (basic example)
        var cardNumber string

        // Adjusted regex for Visa, MasterCard, American Express, and Discover credit card numbers
        // Define the individual patterns
        american_express_regex :=`^3[47][0-9]{2}\s[0-9]{6}\s[0-9]{5}$`
        visa_regex := `^4[0-9]{3}\s[0-9]{4}\s[0-9]{4}\s[0-9]{4}$`
        mastercard_regex := `^5[1-5]{3}\s[0-9]{4}\s[0-9]{4}\s[0-9]{4}$`
        discover_regex := `6(?:011|5[0-9]{2})\s[0-9]{4}\s[0-9]{4}\s[0-9]{4}$`

        // Combine the patterns using alternation
	    combinedRegex := fmt.Sprintf("(%s)|(%s)|(%s)|(%s)", american_express_regex, visa_regex, mastercard_regex, discover_regex)
        re := regexp.MustCompile(combinedRegex)
        matches := re.FindStringSubmatch(text)

        if len(matches) > 0 {
            cardNumber = matches[0]
        }

        // Encrypt the card number using the key
        encryptedCardNumber, err := Encrypt([]byte(cardNumber), []byte(key))
        if err != nil {
            fmt.Println("Failed to encrypt card number:", err)
            return c.String(http.StatusInternalServerError, "Encryption failed")
        }

        return c.JSON(http.StatusOK, map[string]string{"encryptedCardNumber": encryptedCardNumber})
    })

    // Validate credit card number
    e.GET("/validate", func(c echo.Context) error {
        cardNumber := c.QueryParam("cardNumber")
        
        isValid := ValidateCreditCard(cardNumber)  // Call function from validator.go
        cardType := GetCardType(cardNumber)         // Call function from validator.go

        var result string
        if isValid {
            result = "<span style='color: green;'>The credit card number is valid.</span>"
        } else {
            result = "<span style='color: red;'>The credit card number is invalid.</span>"
        }

        return c.HTML(http.StatusOK, `<div class="result">` + result + `</div><div class="card-type">Card Type: ` + cardType + `</div>`)
    })

    
    e.POST("/decrypt", func(c echo.Context) error {
        // Read the request body
        body, err := io.ReadAll(c.Request().Body)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to read request body"})
        }

        // Print the raw request body for debugging
        //fmt.Println("Request Body:", string(body))

        // Restore the request body so it can be used again
        c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

        var req DecryptRequest

        // Bind the request body to the DecryptRequest struct
        if err := json.Unmarshal(body, &req); err != nil { // Use json.Unmarshal directly
            fmt.Printf("Error binding request: %v\n", err) // Log binding error
            return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
        }

        // Log the extracted ciphertext
        //fmt.Println("Ciphertext:", req.Ciphertext)

        // Decrypt the card number
        decryptedCardNumber, err := Decrypt([]byte(req.Ciphertext), []byte(os.Getenv("ENCRYPTION_KEY")))
    
        // Log the error for debugging purposes
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err) // Log detailed error
            return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Decryption failed"})
        }
    
        // Send response back to client
        response := DecryptResponse{DecryptedCardNumber: decryptedCardNumber}
        return c.JSON(http.StatusOK, response)
    })

    e.Logger.Fatal(e.Start(":8080"))
}
