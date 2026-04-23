package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type Claims struct {
	UserID   int64  `json:"sub"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (m *TokenManager) Generate(user domain.User) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Exp:      now.Add(m.ttl).Unix(),
		Iat:      now.Unix(),
	}

	header, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}

	unsigned := encodeTokenPart(header) + "." + encodeTokenPart(payload)
	signature := m.sign(unsigned)
	return unsigned + "." + encodeTokenPart(signature), nil
}

func (m *TokenManager) Parse(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, fmt.Errorf("%w: malformed token", domain.ErrUnauthorized)
	}

	unsigned := parts[0] + "." + parts[1]
	signature, err := decodeTokenPart(parts[2])
	if err != nil {
		return Claims{}, fmt.Errorf("%w: invalid signature", domain.ErrUnauthorized)
	}

	expected := m.sign(unsigned)
	if !hmac.Equal(signature, expected) {
		return Claims{}, fmt.Errorf("%w: invalid signature", domain.ErrUnauthorized)
	}

	payload, err := decodeTokenPart(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("%w: invalid payload", domain.ErrUnauthorized)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, fmt.Errorf("%w: invalid claims", domain.ErrUnauthorized)
	}

	if claims.Exp <= time.Now().UTC().Unix() {
		return Claims{}, fmt.Errorf("%w: token expired", domain.ErrUnauthorized)
	}

	return claims, nil
}

func (m *TokenManager) sign(unsigned string) []byte {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(unsigned))
	return mac.Sum(nil)
}

func encodeTokenPart(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func decodeTokenPart(value string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(value)
}
