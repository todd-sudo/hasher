package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"hasher/internal/config"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// savePrivateKeyToFile - dump private key to file
func savePrivateKeyToFile(privateKey *rsa.PrivateKey, privatePemFileName string) {
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create(privatePemFileName)
	if err != nil {
		fmt.Printf("error when create private.pem: %s \n", err)
		os.Exit(1)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		fmt.Printf("error when encode private pem: %s \n", err)
		os.Exit(1)
	}
}

func uploadPrivateKey(privatePemFileName string) (*rsa.PrivateKey, error) {

	privatePem, err := os.ReadFile(privatePemFileName)
	if err != nil {
		log.Println("Ошибка чтения private.pem", err)
		return nil, err
	}
	block, _ := pem.Decode(privatePem)

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("Ошибка парсинга private key", err)
		return nil, err
	}
	return privKey, nil
}

func generateKeys(cfg config.Config) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, cfg.SizeRSAKey)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func rsaEncrypt(secretMessage string, key rsa.PublicKey) (string, error) {

	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha512.New(), rng, &key, []byte(secretMessage), label)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func rsaDecrypt(cipherText string, privKey rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha512.New(), rng, &privKey, ct, label)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return string(plaintext), nil
}

func hashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func compareHashAndPassword(hashedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
