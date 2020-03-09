/*
	Package des3 supports password encrypt/decrypt
	It depends on 3DES
	Support by 张松
*/

package des3

import (
	"testing"
)

func Test_All(t *testing.T) {
	Setup("http://192.168.2.70:20007", false)
	for _, plaintext := range []string{"abc", "cba", "Abc", "cBa", "122a", "asdf!!!xx++"} {
		ciphertext, err := Encrypt(plaintext)
		if err != nil {
			t.Fatal(err)
		}
		plaintextResult, err := Decrypt(ciphertext)
		if err != nil {
			t.Fatal(err)
		}
		if plaintext != plaintextResult || len(plaintextResult) == 0 {
			t.Fatal(plaintext, plaintextResult, "aren't same")
		}
	}
}

func Test_AllWithCache(t *testing.T) {
	Setup("http://192.168.2.70:20007", true)
	for _, plaintext := range []string{"abc", "cba", "Abc", "cBa", "122a", "asdf!!!xx++"} {
		ciphertext, err := Encrypt(plaintext)
		if err != nil {
			t.Fatal(err)
		}
		plaintextResult, err := Decrypt(ciphertext)
		if err != nil {
			t.Fatal(err)
		}
		if plaintext != plaintextResult || len(plaintextResult) == 0 {
			t.Fatal(plaintext, plaintextResult, "aren't same")
		}
	}
	for _, plaintext := range []string{"abc", "cba", "Abc", "cBa", "122a", "asdf!!!xx++"} {
		ciphertext, err := Encrypt(plaintext)
		if err != nil {
			t.Fatal(err)
		}
		plaintextResult, err := Decrypt(ciphertext)
		if err != nil {
			t.Fatal(err)
		}
		if plaintext != plaintextResult || len(plaintextResult) == 0 {
			t.Fatal(plaintext, plaintextResult, "aren't same")
		}
	}
}

func Test_decrypt(t *testing.T) {
	gotPlaintext, err := decrypt("A942B79235FFBAF299D9D3D928F45C28")
	if err != nil {
		t.Errorf("decrypt() error = %v", err)
		return
	}
	t.Logf("decrypt() = %s ", gotPlaintext)
}
