package client

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
)

// loadKeyFromFile carga una clave de un archivo.
func loadKeyFromFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// encryptAES encripta un mensaje utilizando AES en modo CBC.
func encryptAES(key []byte, message string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	paddedMessage := []byte(message)
	padding := blockSize - len(paddedMessage)%blockSize
	paddedMessage = append(paddedMessage, bytes.Repeat([]byte{byte(padding)}, padding)...)

	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(paddedMessage))
	mode.CryptBlocks(ciphertext, paddedMessage)

	ciphertextWithIV := append(iv, ciphertext...)
	return base64.URLEncoding.EncodeToString(ciphertextWithIV), nil
}

// decryptAES desencripta un mensaje cifrado con la clave AES utilizando AES en modo CBC.
func decryptAES(key []byte, ciphertext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertextWithIV, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	iv := ciphertextWithIV[:blockSize]
	ciphertextBytes := ciphertextWithIV[blockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertextBytes))
	mode.CryptBlocks(plaintext, ciphertextBytes)

	padding := plaintext[len(plaintext)-1]
	plaintext = plaintext[:len(plaintext)-int(padding)]

	return string(plaintext), nil
}
