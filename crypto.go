package bankmuscatpg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"strings"
)

func getAES256GCMEncrypted(plaintext, key string) (string, error) {
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	nonceHex := hex.EncodeToString(nonce)
	ciphertextHex := hex.EncodeToString(ciphertext)

	return nonceHex + ciphertextHex, nil
}

func stringToMap(requestData string) (map[string]interface{}, error) {
	requestMap := make(map[string]interface{})

	keyValuePairs := strings.Split(requestData, "&")

	for _, keyValue := range keyValuePairs {
		parts := strings.SplitN(keyValue, "=", 2)
		key, err := url.QueryUnescape(parts[0])
		if err != nil {
			return nil, err
		}
		value, err := url.QueryUnescape(parts[1])
		if err != nil {
			return nil, err
		}
		requestMap[key] = value
	}

	return requestMap, nil
}

func (Bmpg *BankMuscatPG) DecryptAES256GCM(encryptedTextHex string) (map[string]interface{}, error) {
	encryptedData, err := hex.DecodeString(encryptedTextHex)
	if err != nil {
		return make(map[string]interface{}), err
	}

	key, err := hex.DecodeString(Bmpg.WorkingKey)
	if err != nil {
		return make(map[string]interface{}), err
	}

	nonce := encryptedData[:12]
	ciphertext := encryptedData[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return make(map[string]interface{}), err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return make(map[string]interface{}), err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return make(map[string]interface{}), err
	}

	decryptedMap, err := stringToMap(string(plaintext))
	if err != nil {
		return make(map[string]interface{}), err
	}

	return decryptedMap, nil
}
