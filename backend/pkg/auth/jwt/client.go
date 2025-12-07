/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Client JWT客户端封装
type Client struct {
	secretKey []byte
	method    jwt.SigningMethod
}

// Config JWT配置
type Config struct {
	SecretKey string `yaml:"secret_key" json:"secret_key"`
	Method    string `yaml:"method" json:"method"` // HS256, HS384, HS512
}

// Claims JWT声明
type Claims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// NewClient 创建JWT客户端
func NewClient(config *Config) (*Client, error) {
	var method jwt.SigningMethod
	switch config.Method {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	default:
		method = jwt.SigningMethodHS256
	}

	return &Client{
		secretKey: []byte(config.SecretKey),
		method:    method,
	}, nil
}

// GenerateToken 生成JWT令牌
func (c *Client) GenerateToken(claims *Claims, expiration time.Duration) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(c.method, claims)
	return token.SignedString(c.secretKey)
}

// ParseToken 解析JWT令牌
func (c *Client) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return c.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateToken 验证JWT令牌
func (c *Client) ValidateToken(tokenString string) error {
	_, err := c.ParseToken(tokenString)
	return err
}

// RefreshToken 刷新JWT令牌
func (c *Client) RefreshToken(tokenString string, expiration time.Duration) (string, error) {
	claims, err := c.ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to parse token for refresh: %w", err)
	}

	// 更新过期时间
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(c.method, claims)
	return token.SignedString(c.secretKey)
}

// ExtractClaims 提取声明（不验证签名）
func (c *Client) ExtractClaims(tokenString string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

