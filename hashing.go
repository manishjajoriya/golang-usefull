package auth

import (
	"NoRethink/internal/config"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func Hash(password string, argonCfg config.Argon2idConfig) string {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonCfg.Iteration,
		argonCfg.MemoryMB*1024,
		argonCfg.Threads,
		argonCfg.Key,
	)

	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonCfg.MemoryMB*1024,
		argonCfg.Iteration,
		argonCfg.Threads,
		saltB64,
		hashB64,
	)

	return encoded
}

func Verify(password, encodedHash string) (valid bool, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var memory, iterations uint32
	var parallelism uint8

	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(hash, expectedHash) == 1, nil
}

func Sha256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
