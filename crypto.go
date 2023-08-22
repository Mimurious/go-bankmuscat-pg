package bankmuscatpg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
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
