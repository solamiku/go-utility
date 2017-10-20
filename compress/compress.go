package compress

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"io"

	"github.com/solamiku/go-utility/crypto"
)

var NIL_BYTES = []byte("")

func Zlib(src []byte, compress bool) ([]byte, error) {
	if compress {
		return zlibCompress(src)
	} else {
		return zlibUncompress(src)
	}
}

// compress with zlib
func zlibCompress(src []byte) ([]byte, error) {
	var bf bytes.Buffer
	w, err := zlib.NewWriterLevel(&bf, zlib.BestCompression)
	if err != nil {
		return NIL_BYTES, err
	}
	w.Write(src)
	w.Close()
	return bf.Bytes(), nil
}

// uncompress with zlib
func zlibUncompress(src []byte) ([]byte, error) {
	br := bytes.NewReader(src)
	var out bytes.Buffer
	r, err := zlib.NewReader(br)
	if err != nil {
		return NIL_BYTES, err
	}
	io.Copy(&out, r)
	r.Close()
	return out.Bytes(), nil
}

func GZip(src []byte, compress bool) ([]byte, error) {
	if compress {
		return gZipCompress(src)
	} else {
		return gZipUncompress(src)
	}
}

func gZipCompress(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(src)
	w.Close()
	return buf.Bytes(), nil
}

func gZipUncompress(src []byte) ([]byte, error) {
	b := bytes.NewReader(src)
	var buf bytes.Buffer
	r, err := gzip.NewReader(b)
	if err != nil {
		return NIL_BYTES, err
	}
	io.Copy(&buf, r)
	return buf.Bytes(), nil
}

func Base64encode(src []byte) []byte {
	var b bytes.Buffer
	w := base64.NewEncoder(base64.URLEncoding, &b)
	w.Write(src)
	w.Close()
	return b.Bytes()
}

func Base64decode(src []byte) []byte {
	b := bytes.NewReader(src)
	var buf bytes.Buffer
	r := base64.NewDecoder(base64.URLEncoding, b)
	io.Copy(&buf, r)
	return buf.Bytes()
}

//zlib compress->base64 encode
func ZlibBase64(src []byte) []byte {
	b, _ := zlibCompress(src)
	return Base64encode(b)
}

//base64 encode->zlib uncompress
func Base64Unzlib(src []byte) []byte {
	b := Base64decode(src)
	d, _ := zlibUncompress(b)
	return d
}

// package 0 protocol pack
// zlib compress > des encrypt > base64 encode
func Pack0Encode(src []byte, pwd []byte) []byte {
	s1, _ := Zlib(src, true)
	if len(pwd) == 8 {
		s1, _ = crypto.DesCBC(s1, pwd, true, crypto.PADDING_PKCS5)
	}
	return Base64encode(s1)
}

// package 0 protocol unpack
// base64 decode > des dencrypt > zilib uncompress
func Pack0Decode(src []byte, pwd []byte) []byte {
	s1 := Base64decode(src)
	if len(s1) == 0 || len(s1)%8 != 0 {
		return []byte{}
	}
	if len(pwd) == 8 {
		s1, _ = crypto.DesCBC(s1, pwd, false, crypto.PADDING_PKCS5)
	}
	s1, _ = Zlib(s1, false)
	return s1
}
