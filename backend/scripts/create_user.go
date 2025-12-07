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

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("COZE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 连接数据库
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		viper.GetString("database.host"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
		viper.GetInt("database.port"),
		viper.GetString("database.sslmode"),
		viper.GetString("database.timezone"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect database: %v\n", err)
		os.Exit(1)
	}

	// 获取参数
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run create_user.go <email> <password> [name]")
		fmt.Println("Example: go run create_user.go tim.dong@hotmail.com password123")
		os.Exit(1)
	}

	email := os.Args[1]
	password := os.Args[2]
	name := ""
	if len(os.Args) > 3 {
		name = os.Args[3]
	} else {
		// 如果没有提供名称，使用邮箱前缀
		name = strings.Split(email, "@")[0]
	}

	// 检查用户是否已存在
	var count int64
	db.Table("users").Where("email = ?", email).Count(&count)
	if count > 0 {
		fmt.Printf("User with email %s already exists\n", email)
		os.Exit(1)
	}

	// 生成密码哈希
	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Printf("Failed to hash password: %v\n", err)
		os.Exit(1)
	}

	// 生成用户ID（简单使用时间戳，实际应该使用ID生成器）
	userID := time.Now().UnixNano() / 1000000 // 毫秒时间戳
	now := time.Now().UnixMilli()

	// 创建唯一名称（从邮箱生成）
	uniqueName := strings.Split(email, "@")[0]
	// 确保唯一名称唯一
	var uniqueNameCount int64
	db.Table("users").Where("unique_name = ?", uniqueName).Count(&uniqueNameCount)
	if uniqueNameCount > 0 {
		uniqueName = fmt.Sprintf("%s%d", uniqueName, now%10000)
	}

	// 创建用户
	user := map[string]interface{}{
		"id":            userID,
		"email":         email,
		"username":      name,
		"name":          name,
		"unique_name":   uniqueName,
		"password_hash": hashedPassword,
		"icon_uri":      "coze-studio/default_icon/user_default_icon.png",
		"user_verified": false,
		"created_at_ms": now,
		"updated_at_ms": now,
	}

	result := db.Table("users").Create(user)
	if result.Error != nil {
		fmt.Printf("Failed to create user: %v\n", result.Error)
		os.Exit(1)
	}

	fmt.Printf("User created successfully!\n")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("User ID: %d\n", userID)
}

// hashPassword 使用 Argon2id 算法生成密码哈希
func hashPassword(password string) (string, error) {
	type argon2Params struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}

	defaultArgon2Params := argon2Params{
		memory:      64 * 1024, // 64MB
		iterations:  3,
		parallelism: 4,
		saltLength:  16,
		keyLength:  32,
	}

	p := defaultArgon2Params

	// 生成随机 salt
	salt := make([]byte, p.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// 使用 Argon2id 算法计算哈希值
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// 编码为 base64 格式
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encoded, nil
}

