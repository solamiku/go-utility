package compress

import (
	"bytes"
	"testing"

	"github.com/solamiku/go-utility/crypto"
)

func Test_compress(t *testing.T) {
	ori := bytes.NewBufferString("test")
	dst, err := Zlib(ori.Bytes(), true)
	if err != nil {
		t.Fatal(err)
	}
	dori, err := Zlib(dst, false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ori.Bytes(), dori) {
		t.Errorf("%s not equal %s", ori, dori)
	}

	dst2 := ZlibBase64(ori.Bytes())
	dori2 := Base64Unzlib(dst2)
	if !bytes.Equal(ori.Bytes(), dori2) {
		t.Errorf("%s not equal %s", ori, dori)
	}

	dst3, err := GZip(ori.Bytes(), true)
	if err != nil {
		t.Fatal(err)
	}
	dori3, err := GZip(dst3, false)
	if !bytes.Equal(ori.Bytes(), dori3) {
		t.Errorf("%s not equal %s", ori, dori)
	}

	dst4 := Pack0Encode(ori.Bytes(), []byte("test1201"))
	dori4 := Pack0Decode(dst4, []byte("test1201"))
	if !bytes.Equal(dori4, ori.Bytes()) {
		t.Errorf("%s not equal %s", ori, dori4)
	}

	dst5 := []byte("KALoJ8dMthePlcoRQ/1f+g==")
	dst5 = Base64decode(dst5)
	r5, err := crypto.DesECB(dst5, []byte("45985633"), false, crypto.PADDING_PKCS5)
	t.Log(string(r5), err)
}
