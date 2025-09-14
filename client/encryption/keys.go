package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
)

const KeySize = 2048

func GetUserKeyDir(username string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".cli-chat-app", "keys", username)
	return dir, nil
}

func SaveKeyPair(username string) error {
	dir, err := GetUserKeyDir(username)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		return err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	pubPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	pubFile := filepath.Join(dir, "public.pem")
	if err := os.WriteFile(pubFile, pem.EncodeToMemory(pubPem), 0644); err != nil {
		return err
	}

	return nil
}

func LoadCurrentUserPrivateKey(username string) (*rsa.PrivateKey, error) {
	dir, err := GetUserKeyDir(username)
	if err != nil {
		return nil, err
	}

	privFile := filepath.Join(dir, "private.pem")
	privBytes, err := os.ReadFile(privFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key format")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadCurrentUserPublicKey(username string) (*rsa.PublicKey, error) {
	dir, err := GetUserKeyDir(username)
	if err != nil {
		return nil, err
	}

	pubFile := filepath.Join(dir, "public.pem")
	pubBytes, err := os.ReadFile(pubFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key format")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if pubKey, ok := pub.(*rsa.PublicKey); ok {
		return pubKey, nil
	}

	return nil, errors.New("no RSA public key found")
}

func LoadRecipientPublicKey(recipientUsername string) (*rsa.PublicKey, error) {

	//make api call to get public key of recipient from server

	//temporary solution : loading from local system
	return LoadCurrentUserPublicKey(recipientUsername)

	// return nil, errors.New("recipient public key not found, implement server API call to fetch it")
}
