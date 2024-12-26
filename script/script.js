// script/script.js
function validateCardNumber() {
    const cardNumber = document.getElementById('cardNumber').value;
    const resultDiv = document.getElementById('result');
    const cardTypeDiv = document.getElementById('cardType');
    
    if (!/^[0-9]*$/.test(cardNumber)) {
        resultDiv.innerHTML = "Please enter only numeric characters.";
        cardTypeDiv.innerHTML = "";
        return;
    }

    fetch('/validate?cardNumber=' + encodeURIComponent(cardNumber))
        .then(response => response.text())
        .then(data => {
            const parser = new DOMParser();
            const doc = parser.parseFromString(data, 'text/html');
            resultDiv.innerHTML = doc.querySelector('.result').innerHTML || "No result found.";
            cardTypeDiv.innerHTML = doc.querySelector('.card-type').innerHTML || "";
        });
}

async function uploadImage() {
    const fileInput = document.getElementById('cardImage');
    const formData = new FormData();
    
    formData.append('image', fileInput.files[0]);

    const response = await fetch('/upload', {
        method: 'POST',
        body: formData
    });
    
    if (!response.ok) {
        const errorText = await response.json();
        alert(errorText.message);
        return;
    }
    
    const data = await response.json();

    if (data.encryptedCardNumber) {

        const decryptedCardNumber = await decryptCardNumber(data.encryptedCardNumber);

        // Clean up and populate the card number field
        const cleanedCardNumber = decryptedCardNumber.replaceAll(' ', '');
        document.getElementById('cardNumber').value = cleanedCardNumber; // Populate the card number field
        validateCardNumber(); // Validate the extracted number
    } else {
        alert("Failed to extract card information.");
    }
}

// Function to call the decryption endpoint on the server
async function decryptCardNumber(encryptedCardNumber) {
    const response = await fetch('/decrypt', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ "encryptedCardNumber":encryptedCardNumber }),
    });

 
    if (!response.ok) {
        const errorText = await response.text(); // Capture error message for debugging
        console.error('Decryption failed:', errorText);
        throw new Error('Decryption failed');
    }
 
    const data = await response.json();
    return data.decryptedCardNumber;
 }
