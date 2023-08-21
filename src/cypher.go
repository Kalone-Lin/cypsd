package src

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
)

var key = "4321dfgsgf4512"

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

/**
Aes 加密
*/
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

/**
Aes解密
*/
func AesDecryptECB(encrypted []byte, key []byte) ([]byte, error) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted := make([]byte, len(encrypted))
	if len(encrypted)%cipher.BlockSize() != 0 {
		return []byte{}, fmt.Errorf("invalid encrypted content, invalid length")
	}
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	if trim < 0 {
		return []byte{}, fmt.Errorf("invalid encrypted content, invalid trim")
	}
	return decrypted[:trim], nil
}

/**
加密密码
*/
func EncodeRedisPassword(pwd string) (string, error) {
	if len(pwd) == 0 {
		return "", nil
	}
	key := []byte(key)
	data := []byte(pwd)
	endata := AesEncryptECB(data, key)
	return base64.StdEncoding.EncodeToString(endata), nil
}

/**
解密密码
*/
func DecodeRedisPassword(enpwd string) (string, error) {
	if len(enpwd) == 0 {
		return "", nil
	}
	key := []byte(key)
	data, err := base64.StdEncoding.DecodeString(enpwd)
	if err != nil {
		return "", err
	}
	dedata, err := AesDecryptECB(data, key)
	if err != nil {
		return "", err
	}
	return string(dedata), nil
}
