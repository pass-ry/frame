// AES-256

package aes

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateKey()
			t.Logf("%s GenerateKey() = %v", tt.name, got)
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		key       string
		plaintext string
	}
	tests := []struct {
		name           string
		args           args
		wantCiphertext string
		wantErr        bool
	}{
		{
			name: "1",
			args: args{
				key:       "myverystrongpasswordo32bitlength",
				plaintext: "Hello World",
			},
			wantCiphertext: "3ed419795150ede4820f19d104d1654a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCiphertext, err := Encrypt(tt.args.key, tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt(%s) error = %v, wantErr %v", tt.args.plaintext, err, tt.wantErr)
				return
			}
			if gotCiphertext != tt.wantCiphertext {
				t.Errorf("Encrypt(%s) = %v, want %v", tt.args.plaintext, gotCiphertext, tt.wantCiphertext)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		key        string
		ciphertext string
	}
	tests := []struct {
		name          string
		args          args
		wantPlaintext string
		wantErr       bool
	}{
		{
			name: "1",
			args: args{
				key:        "myverystrongpasswordo32bitlength",
				ciphertext: "3ed419795150ede4820f19d104d1654a",
			},
			wantPlaintext: "Hello World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPlaintext, err := Decrypt(tt.args.key, tt.args.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPlaintext != tt.wantPlaintext {
				t.Errorf("Decrypt() = %v, want %v", gotPlaintext, tt.wantPlaintext)
			}
		})
	}
}

func TestIntergration(t *testing.T) {}

func BenchmarkEncrypt(b *testing.B) {}

func BenchmarkDecrypt(b *testing.B) {}
