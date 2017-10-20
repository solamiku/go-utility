package crypto

import "crypto/cipher"

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("dec_ecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("dec_ecb: output smaller than input")
	}

	for len(src) > 0 {
		// Write to the dst
		x.b.Encrypt(dst[:x.blockSize], src[:x.blockSize])

		// Move to the next block
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("dec_ecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("dec_ecb: output smaller than input")
	}

	if len(src) == 0 {
		return
	}

	for len(src) > 0 {
		// Write to the dst
		x.b.Decrypt(dst[:x.blockSize], src[:x.blockSize])

		// Move to the next block
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
