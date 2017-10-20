package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

//AES - CBC
func AesCBC(ori, key []byte, encrpt bool, paddingType ...int) ([]byte, error) {
	if encrpt {
		return aesEncrypt(ori, key, AES_CBC, paddingType...)
	} else {
		return aesDecrypt(ori, key, AES_CBC, paddingType...)
	}
}

//AES-ECB
func AesECB(ori, key []byte, encrpt bool, paddingType ...int) ([]byte, error) {
	if encrpt {
		return aesEncrypt(ori, key, AES_ECB, paddingType...)
	} else {
		return aesDecrypt(ori, key, AES_ECB, paddingType...)
	}
}

func aesEncrypt(ori, key []byte, nType int, paddingType ...int) ([]byte, error) {
	nPaddingType := PADDING_ZERO
	if len(paddingType) > 0 {
		nPaddingType = paddingType[0]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return NIL_BYTES, err
	}
	switch nPaddingType {
	case PADDING_ZERO:
		ori = zeroPadding(ori, block.BlockSize())
	case PADDING_PKCS5:
		ori = pKCS5Padding(ori, block.BlockSize())
	default:
		return NIL_BYTES, CryptoError("no padding way.")
	}
	var blockMode cipher.BlockMode
	switch nType {
	case DES_CBC:
		blockMode = cipher.NewCBCEncrypter(block, key)
	case DES_ECB:
		blockMode = NewECBEncrypter(block)
	}
	if blockMode == nil {
		return NIL_BYTES, CryptoError("no block mode.")
	}

	crypted := make([]byte, len(ori))
	blockMode.CryptBlocks(crypted, ori)
	return crypted, nil
}

func aesDecrypt(ori, key []byte, nType int, paddingType ...int) ([]byte, error) {
	nPaddingType := PADDING_ZERO
	if len(paddingType) > 0 {
		nPaddingType = paddingType[0]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return NIL_BYTES, err
	}
	var blockMode cipher.BlockMode
	switch nType {
	case DES_CBC:
		blockMode = cipher.NewCBCDecrypter(block, key)
	case DES_ECB:
		blockMode = NewECBDecrypter(block)
	}
	if blockMode == nil {
		return NIL_BYTES, CryptoError("no block mode.")
	}
	oridata := make([]byte, len(ori))
	blockMode.CryptBlocks(oridata, ori)
	switch nPaddingType {
	case PADDING_ZERO:
		oridata = zeroUnPadding(oridata)
	case PADDING_PKCS5:
		oridata = pKCS5UnPadding(oridata)
	default:
		return NIL_BYTES, CryptoError("no padding way.")
	}
	return oridata, nil
}
