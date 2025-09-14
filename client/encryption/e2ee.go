package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

// EncryptWithPublicKey encrypts a plain text message using recipient's public key.
// It returns base64-encoded ciphertext.
func encryptWithPublicKey(pub *rsa.PublicKey, message string) (string, error) {
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(message))
	if err != nil {
		return "", err
	}	
	// Base64 encode for safe transmission as string
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// DecryptWithPrivateKey decrypts base64-encoded ciphertext using recipient's private key.
func decryptWithPrivateKey(priv *rsa.PrivateKey, base64Cipher string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(base64Cipher)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, priv, cipherBytes)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}

func EncryptMessage(username, recipientUsername, message string) (string, error) {
	pubKey, err := LoadRecipientPublicKey(recipientUsername)
	if err != nil {
		return "", err
	}
	return encryptWithPublicKey(pubKey, message)
}

func DecryptMessage(username, base64Cipher string) (string, error) {
	privKey, err := LoadCurrentUserPrivateKey(username)
	if err != nil {
		return "", err
	}
	return decryptWithPrivateKey(privKey, base64Cipher);
}