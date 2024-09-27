package configparser

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/google/uuid"
)

type Service struct {
	secret string
}

func New(secret string) (*Service, error) {
	if len(secret) != 32 {
		return nil, errors.New("invalid key len")
	}

	return &Service{
		secret: secret,
	}, nil
}

func (s *Service) ParseConfig(cfgLine string) (uuid.UUID, error) {
	c, err := aes.NewCipher([]byte(s.secret))
	if err != nil {
		return uuid.UUID{}, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return uuid.UUID{}, err
	}

	decoded, err := hex.DecodeString(cfgLine)
	if err != nil {
		return uuid.UUID{}, err
	}

	nonceSize := gcm.NonceSize()

	if len(decoded) < nonceSize {
		return uuid.UUID{}, errors.New("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]

	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(string(decrypted))
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (s *Service) CreateConfig(userID uuid.UUID) (string, error) {
	c, err := aes.NewCipher([]byte(s.secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(userID.String()), nil)

	return hex.EncodeToString(encrypted), nil
}
