# Credit Card Validator

## Overview

The Credit Card Validator is a Go-based web application that provides functionality to securely validate, encrypt, and decrypt credit card numbers. The application uses AES-GCM encryption to ensure that sensitive data is protected during transmission and storage.

## Features

- Validate credit card numbers using the Luhn algorithm.
- Encrypt credit card numbers for secure storage.
- Decrypt encrypted credit card numbers for verification.
- RESTful API endpoints for easy integration with other applications.

## Technologies Used

- Go (Golang)
- Echo web framework
- Go Teserract for OCR
- AES-GCM for encryption
- Go modules for dependency management
- Air for hot reloading during development

## Getting Started

### Prerequisites

Before running the application, ensure you have the following installed:

- Go 1.16 or later
- Git
- [Air](https://github.com/cosmtrek/air) for hot reloading

### Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/RoystonDAlmeida/credit-card-validator.git
   cd credit-card-validator
   ```

2. **Set up environment variables**:

   Create a `.env` file in the root of the project directory and add your encryption key:

   ```bash
   ENCRYPTION_KEY=
   ```

   Ensure that the key is exactly 16 bytes long for AES-128 encryption.

3. **Install dependencies**:

   Run the following command to install required packages:

   ```bash
   go mod tidy
   ```

4. **Install Air**:

   If you haven't installed Air yet, you can do so by running:

   ```bash
   go install github.com/cosmtrek/air@latest
   ```

### Running the Application

To start the application with hot reloading using Air, run:

```bash
air
```

The server will start on `http://localhost:8080`, and any changes you make to your code will automatically trigger a reload.

### API Endpoints

#### 1. Validate Credit Card Number

**Endpoint**: `/validate`

**Method**: `POST`

**Request Body**:
```json
    {
        "cardNumber": " "
    }
```

**Response**:
```json
    {
        "valid": true,
        "message": "Valid credit card number."
    }
```

#### 2. Encrypt Credit Card Number

**Endpoint**: `/encrypt`

**Method**: `POST`

**Request Body**:
```json
    {
        "cardNumber": " "
    }
```

**Response**:
```json
    {
        "encryptedCardNumber": "base64_encoded_encrypted_string"
    }
```

#### 3. Decrypt Credit Card Number

**Endpoint**: `/decrypt`

**Method**: `POST`

**Request Body**:
```json
    {
        "encryptedCardNumber": "base64_encoded_encrypted_string"
    }
```

**Response**:
```json
    {
        "decryptedCardNumber": " "
    }
```

## Security Considerations

- Always use HTTPS in production to protect sensitive data in transit.
- Rotate your encryption keys regularly and store them securely.
- Ensure that access to this service is restricted to authorized users only.

## Testing

You can run tests using Go's built-in testing framework. To execute tests, run:

```bash
    go test ./...
```

## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the Go community for their support and resources.
- Special thanks to contributors who help improve this project.
```

You can copy this entire block and paste it into your README.md file in your GitHub repository. Adjust any specific details as needed to fit your actual implementation or project structure.
