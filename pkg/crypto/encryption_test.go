package crypto_test

import (
	"testing"

	"github.com/podoru/spinner-podoru/pkg/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	encryptor, err := crypto.NewEncryptor("test-secret-key")
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	plaintext := []byte("hello world")

	ciphertext, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("expected %s, got %s", plaintext, decrypted)
	}
}

func TestEncryptDecryptString(t *testing.T) {
	encryptor, err := crypto.NewEncryptor("test-secret-key")
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	plaintext := "sensitive data"

	ciphertext, err := encryptor.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if ciphertext == plaintext {
		t.Error("ciphertext should not equal plaintext")
	}

	decrypted, err := encryptor.DecryptString(ciphertext)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("expected %s, got %s", plaintext, decrypted)
	}
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	encryptor, err := crypto.NewEncryptor("test-secret-key")
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	_, err = encryptor.Decrypt([]byte("short"))
	if err != crypto.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestHashPassword(t *testing.T) {
	password := "mypassword123"

	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("hashing failed: %v", err)
	}

	if hash == password {
		t.Error("hash should not equal password")
	}

	if len(hash) < 20 {
		t.Error("hash seems too short")
	}
}

func TestCheckPassword_Valid(t *testing.T) {
	password := "mypassword123"

	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("hashing failed: %v", err)
	}

	if !crypto.CheckPassword(password, hash) {
		t.Error("expected password to match hash")
	}
}

func TestCheckPassword_Invalid(t *testing.T) {
	password := "mypassword123"

	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("hashing failed: %v", err)
	}

	if crypto.CheckPassword("wrongpassword", hash) {
		t.Error("expected wrong password to not match hash")
	}
}

func TestGenerateRandomString(t *testing.T) {
	length := 16

	str1, err := crypto.GenerateRandomString(length)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	if len(str1) != length {
		t.Errorf("expected length %d, got %d", length, len(str1))
	}

	str2, err := crypto.GenerateRandomString(length)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	if str1 == str2 {
		t.Error("random strings should be different")
	}
}

func TestHashToken(t *testing.T) {
	token := "my-secret-token"

	hash1 := crypto.HashToken(token)
	hash2 := crypto.HashToken(token)

	if hash1 != hash2 {
		t.Error("same token should produce same hash")
	}

	if hash1 == token {
		t.Error("hash should not equal token")
	}

	differentHash := crypto.HashToken("different-token")
	if hash1 == differentHash {
		t.Error("different tokens should produce different hashes")
	}
}

func TestDifferentKeys(t *testing.T) {
	encryptor1, _ := crypto.NewEncryptor("key1")
	encryptor2, _ := crypto.NewEncryptor("key2")

	plaintext := []byte("secret message")

	ciphertext, err := encryptor1.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	_, err = encryptor2.Decrypt(ciphertext)
	if err == nil {
		t.Error("decryption with different key should fail")
	}
}
