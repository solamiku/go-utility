package crypto

import (
	"bytes"
	"math/rand"
	"time"
)

type CryptoError string

func (ce CryptoError) Error() string {
	return string(ce)
}

const (
	PADDING_ZERO  = 0
	PADDING_PKCS5 = 1

	DES_CBC = 0
	DES_ECB = 1

	AES_CBC = 0
	AES_ECB = 1
)

var NIL_BYTES = []byte("")

// generate simple rand password
func GenRandPassword(l int, symbol ...bool) []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if len(symbol) > 0 && symbol[0] {
		str += "!@#$%^&*()_+"
	}
	bytes := []byte(str)
	result := make([]byte, 0, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}

/*
	padding way
*/
//zero padding
func zeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(cipherText, padtext...)
}

//end remove padding 0
func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

//use number which just need to padding in end，if not need to padding
//force to padding blockSize's blockSize number
func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//remove pkcs5 padding
func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// remove last byte  unpadding times
	unpadding := int(origData[length-1])
	end := length - unpadding
	if end > len(origData) || end < 0 {
		return []byte("")
	}
	return origData[:end]
}

//cut out bytes as given's length
func unPaddingByLength(origData []byte, length int) []byte {
	olen := len(origData)
	if length > olen {
		panic("dec_ecb : unpadding length is longger than origin data")
	}
	return origData[:(olen - length)]
}
