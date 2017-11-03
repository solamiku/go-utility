package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func Test_des(t *testing.T) {
	test_descbc_zero(t)
	test_desecb_zero(t)
	test_aescbc_zero(t)
	test_aesecb_zero(t)
	test_rsa(t)
	test_genrandpwd(t)
}

func test_descbc_zero(t *testing.T) {
	ori := []byte("test")
	key := []byte("test1089")
	dst, err := DesCBC(ori, key, true)
	dstr := base64.StdEncoding.EncodeToString(dst)
	if dstr != "PJvKFe/jlLE=" {
		t.Errorf("ori:%s dst:%s err:%v", ori, dstr, err)
	}

	tori, err := DesCBC(dst, key, false)
	if string(tori) != "test" {
		t.Errorf("tori :", string(tori))
	}
}

func test_desecb_zero(t *testing.T) {
	ori := []byte("test")
	key := []byte("test1089")
	dst, err := DesECB(ori, key, true)
	dstr := base64.StdEncoding.EncodeToString(dst)
	if dstr != "yuicAaIikPU=" {
		t.Errorf("ori:%s dst:%s err:%v", ori, dstr, err)
	}

	tori, err := DesECB(dst, key, false)
	if string(tori) != "test" {
		t.Errorf("tori :", string(tori))
	}
}

func test_aescbc_zero(t *testing.T) {
	ori := []byte("test")
	key := []byte("test1089test1089")
	dst, err := AesCBC(ori, key, true)
	dstr := base64.StdEncoding.EncodeToString(dst)
	if dstr != "prDUTeCwfXzrQC4OFYDqEQ==" {
		t.Errorf("ori:%s dst:%s err:%v", ori, dstr, err)
	}

	tori, err := AesCBC(dst, key, false)
	if string(tori) != "test" {
		t.Errorf("tori :", string(tori))
	}
}

func test_aesecb_zero(t *testing.T) {
	ori := []byte("test")
	key := []byte("test1089test1089")
	dst, err := AesECB(ori, key, true)
	dstr := base64.StdEncoding.EncodeToString(dst)
	if dstr != "TkvhLfapVhIlBkrD+VmyRQ==" {
		t.Errorf("ori:%s dst:%s err:%v", ori, dstr, err)
	}

	tori, err := AesECB(dst, key, false)
	if string(tori) != "test" {
		t.Errorf("tori :", string(tori))
	}
}

func test_rsa(t *testing.T) {
	private, public, err := GenerateRsaKey(1024)
	if err != nil {
		t.Fatal("rsa err:%v", err)
	}
	//	t.Log(private)
	//	t.Log(public)

	ori := []byte("test")
	encrypt, err := EncryptRsa(ori, public)
	if err != nil {
		t.Fatal("encrypt rsa err:%v", err)
	}
	//	t.Logf("encrypt %s to %s", ori, base64.StdEncoding.EncodeToString(encrypt))

	dori, err := DecryptRsa(encrypt, private)
	if err != nil {
		t.Fatal("decrypt rsa err:%v", err)
	}
	if !bytes.Equal(dori, ori) {
		t.Errorf("rsa %s not equal %s", ori, dori)
	}
}

func test_genrandpwd(t *testing.T) {
	t.Log(string(GenRandPassword(10)))
	t.Log(string(GenRandPassword(10, true)))
}

func Benchmark_crypto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRandPassword(10)
	}
}
