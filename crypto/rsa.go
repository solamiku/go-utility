package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

//generate ras key file with x509 standard
//return private key, public key, error
func GenerateRsaKey(bits int) (string, string, error) {
	//private key
	privatek, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	desStream := x509.MarshalPKCS1PrivateKey(privatek)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: desStream,
	}
	privateBuf := bytes.NewBuffer([]byte(""))
	err = pem.Encode(privateBuf, block)
	if err != nil {
		return "", "", err
	}
	//public key
	publick := &privatek.PublicKey
	derpkix, err := x509.MarshalPKIXPublicKey(publick)
	if err != nil {
		return "", "", err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derpkix,
	}
	publicBuf := bytes.NewBuffer([]byte(""))
	err = pem.Encode(publicBuf, block)
	if err != nil {
		return "", "", err
	}
	return privateBuf.String(), publicBuf.String(), nil
}

func EncryptRsa(ori []byte, publicKey string) ([]byte, error) {
	publicBuf := bytes.NewBufferString(publicKey)
	block, _ := pem.Decode(publicBuf.Bytes())
	if block == nil {
		return NIL_BYTES, CryptoError("decode public key nil")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return NIL_BYTES, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, ori)
}

func DecryptRsa(ciphertext []byte, privateKey string) ([]byte, error) {
	privateBuf := bytes.NewBufferString(privateKey)
	block, _ := pem.Decode(privateBuf.Bytes())
	if block == nil {
		return NIL_BYTES, CryptoError("decode public key nil")
	}
	priInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return NIL_BYTES, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priInterface, ciphertext)
}
