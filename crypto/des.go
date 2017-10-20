package crypto

import (
	"crypto/cipher"
	"crypto/des"
)

//DES - CBC
func DesCBC(ori, key []byte, encrpt bool, paddingType ...int) ([]byte, error) {
	if encrpt {
		return desEncrypt(ori, key, DES_CBC, paddingType...)
	} else {
		return desDecrypt(ori, key, DES_CBC, paddingType...)
	}
}

//DES-ECB
func DesECB(ori, key []byte, encrpt bool, paddingType ...int) ([]byte, error) {
	if encrpt {
		return desEncrypt(ori, key, DES_ECB, paddingType...)
	} else {
		return desDecrypt(ori, key, DES_ECB, paddingType...)
	}
}

func desEncrypt(ori, key []byte, nType int, paddingType ...int) ([]byte, error) {
	nPaddingType := PADDING_ZERO
	if len(paddingType) > 0 {
		nPaddingType = paddingType[0]
	}
	block, err := des.NewCipher(key)
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

func desDecrypt(ori, key []byte, nType int, paddingType ...int) ([]byte, error) {
	nPaddingType := PADDING_ZERO
	if len(paddingType) > 0 {
		nPaddingType = paddingType[0]
	}
	block, err := des.NewCipher(key)
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
