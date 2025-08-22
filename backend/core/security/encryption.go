package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// TokenEncryption provides AES-256 encryption for sensitive tokens
type TokenEncryption struct {
	key []byte
}

// NewTokenEncryption creates a new token encryption instance
// The key must be 32 bytes for AES-256
func NewTokenEncryption(key []byte) (*TokenEncryption, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be exactly 32 bytes, got %d", len(key))
	}

	return &TokenEncryption{key: key}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func (e *TokenEncryption) Encrypt(plaintext string) ([]byte, error) {
	if plaintext == "" {
		return nil, fmt.Errorf("plaintext cannot be empty")
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func (e *TokenEncryption) Decrypt(ciphertext []byte) (string, error) {
	if len(ciphertext) == 0 {
		return "", fmt.Errorf("ciphertext cannot be empty")
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// EncryptTokens encrypts OAuth tokens for secure storage
func (e *TokenEncryption) EncryptTokens(accessToken, refreshToken string) ([]byte, []byte, error) {
	var encryptedAccess, encryptedRefresh []byte
	var err error

	if accessToken != "" {
		encryptedAccess, err = e.Encrypt(accessToken)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt access token: %w", err)
		}
	}

	if refreshToken != "" {
		encryptedRefresh, err = e.Encrypt(refreshToken)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt refresh token: %w", err)
		}
	}

	return encryptedAccess, encryptedRefresh, nil
}

// DecryptTokens decrypts OAuth tokens from secure storage
func (e *TokenEncryption) DecryptTokens(encryptedAccess, encryptedRefresh []byte) (string, string, error) {
	var accessToken, refreshToken string
	var err error

	if len(encryptedAccess) > 0 {
		accessToken, err = e.Decrypt(encryptedAccess)
		if err != nil {
			return "", "", fmt.Errorf("failed to decrypt access token: %w", err)
		}
	}

	if len(encryptedRefresh) > 0 {
		refreshToken, err = e.Decrypt(encryptedRefresh)
		if err != nil {
			return "", "", fmt.Errorf("failed to decrypt refresh token: %w", err)
		}
	}

	return accessToken, refreshToken, nil
}

// GenerateEncryptionKey generates a new 32-byte encryption key
func GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}
	return key, nil
}
