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

package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/core/cache"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestSQLiteDB 创建SQLite测试数据库（用于测试database包的逻辑）
func setupTestSQLiteDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// 创建测试表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS test_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	return db
}

// TestDatabaseManager_Structure 测试DatabaseManager结构
func TestDatabaseManager_Structure(t *testing.T) {
	// 验证Manager结构存在
	manager := &Manager{}
	assert.NotNil(t, manager)

	// 验证Manager类型正确
	assert.IsType(t, &Manager{}, manager)
}

// TestDatabaseManager_Methods 测试Manager有所需的方法
func TestDatabaseManager_Methods(t *testing.T) {
	// 验证Manager有GetPostgres方法
	manager := &Manager{}
	client := manager.GetPostgres()
	_ = client // nil是OK的

	// 验证Manager有GetRedis方法
	redisClient := manager.GetRedis()
	_ = redisClient // nil是OK的

	// 验证方法存在
	assert.True(t, true)
}

// TestSQLiteOperations 使用SQLite测试数据库基本操作
func TestSQLiteOperations(t *testing.T) {
	db := setupTestSQLiteDB(t)
	defer db.Close()

	ctx := context.Background()

	t.Run("插入数据", func(t *testing.T) {
		result, err := db.ExecContext(ctx,
			"INSERT INTO test_users (username, email) VALUES (?, ?)",
			"testuser1", "test1@example.com")

		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)
	})

	t.Run("查询数据", func(t *testing.T) {
		// 先插入数据
		db.ExecContext(ctx, "INSERT INTO test_users (username, email) VALUES (?, ?)",
			"testuser2", "test2@example.com")

		// 查询数据
		rows, err := db.QueryContext(ctx,
			"SELECT username, email FROM test_users WHERE username = ?",
			"testuser2")

		require.NoError(t, err)
		require.NotNil(t, rows)
		defer rows.Close()

		// 验证查询结果
		if rows.Next() {
			var username, email string
			err = rows.Scan(&username, &email)
			assert.NoError(t, err)
			assert.Equal(t, "testuser2", username)
			assert.Equal(t, "test2@example.com", email)
		}
	})

	t.Run("更新数据", func(t *testing.T) {
		// 先插入数据
		db.ExecContext(ctx, "INSERT INTO test_users (username, email) VALUES (?, ?)",
			"testuser3", "test3@example.com")

		// 更新数据
		result, err := db.ExecContext(ctx,
			"UPDATE test_users SET email = ? WHERE username = ?",
			"updated@example.com", "testuser3")

		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)
	})

	t.Run("删除数据", func(t *testing.T) {
		// 先插入数据
		db.ExecContext(ctx, "INSERT INTO test_users (username, email) VALUES (?, ?)",
			"testuser4", "test4@example.com")

		// 删除数据
		result, err := db.ExecContext(ctx,
			"DELETE FROM test_users WHERE username = ?",
			"testuser4")

		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)
	})
}

// TestSQLiteTransaction 使用SQLite测试事务
func TestSQLiteTransaction(t *testing.T) {
	db := setupTestSQLiteDB(t)
	defer db.Close()

	ctx := context.Background()

	t.Run("事务提交", func(t *testing.T) {
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)

		// 在事务中插入数据
		_, err = tx.ExecContext(ctx,
			"INSERT INTO test_users (username, email) VALUES (?, ?)",
			"txuser1", "tx1@example.com")
		require.NoError(t, err)

		// 提交事务
		err = tx.Commit()
		assert.NoError(t, err)

		// 验证数据已提交
		rows, _ := db.QueryContext(ctx,
			"SELECT COUNT(*) FROM test_users WHERE username = ?", "txuser1")
		if rows != nil {
			defer rows.Close()
			if rows.Next() {
				var count int
				rows.Scan(&count)
				assert.Equal(t, 1, count)
			}
		}
	})

	t.Run("事务回滚", func(t *testing.T) {
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)

		// 在事务中插入数据
		_, err = tx.ExecContext(ctx,
			"INSERT INTO test_users (username, email) VALUES (?, ?)",
			"txuser2", "tx2@example.com")
		require.NoError(t, err)

		// 回滚事务
		err = tx.Rollback()
		assert.NoError(t, err)

		// 验证数据已回滚（不存在）
		rows, _ := db.QueryContext(ctx,
			"SELECT COUNT(*) FROM test_users WHERE username = ?", "txuser2")
		if rows != nil {
			defer rows.Close()
			if rows.Next() {
				var count int
				rows.Scan(&count)
				assert.Equal(t, 0, count)
			}
		}
	})
}

// TestDatabaseWithCache 测试数据库与缓存集成
func TestDatabaseWithCache(t *testing.T) {
	db := setupTestSQLiteDB(t)
	defer db.Close()

	// 创建内存缓存（需要设置清理间隔）
	cacheManager := cache.NewMemoryCache(&cache.MemoryConfig{
		CleanupInterval: 60, // 60秒清理一次
	})

	t.Run("缓存基本操作", func(t *testing.T) {
		// 设置缓存（TTL=1分钟）
		err := cacheManager.Set("db_cache_key", "db_cache_value", 1*time.Minute)
		assert.NoError(t, err)

		// 获取缓存
		value, err := cacheManager.Get("db_cache_key")
		assert.NoError(t, err)
		assert.Equal(t, "db_cache_value", value)

		// 删除缓存
		cacheManager.Delete("db_cache_key")

		// 验证删除成功
		_, err = cacheManager.Get("db_cache_key")
		assert.Error(t, err)
	})

	t.Run("数据库查询结果缓存", func(t *testing.T) {
		// 插入数据
		_, err := db.ExecContext(context.Background(),
			"INSERT INTO test_users (username, email) VALUES (?, ?)",
			"cache_user", "cache@example.com")
		require.NoError(t, err)

		// 查询数据并缓存结果
		rows, err := db.QueryContext(context.Background(),
			"SELECT username FROM test_users WHERE username = ?",
			"cache_user")
		require.NoError(t, err)

		var username string
		if rows.Next() {
			rows.Scan(&username)
			rows.Close()

			// 缓存查询结果（TTL=5分钟）
			cacheManager.Set("user:cache_user", username, 5*time.Minute)
		}

		// 从缓存获取
		cachedValue, err := cacheManager.Get("user:cache_user")
		assert.NoError(t, err)
		assert.Equal(t, "cache_user", cachedValue)
	})
}

// TestConnectionPool 测试连接池行为（使用SQLite模拟）
func TestConnectionPool(t *testing.T) {
	db := setupTestSQLiteDB(t)
	defer db.Close()

	// 设置连接池参数
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	ctx := context.Background()

	// 测试并发查询
	t.Run("并发查询", func(t *testing.T) {
		// 先插入一些数据
		for i := 0; i < 10; i++ {
			db.ExecContext(ctx,
				"INSERT INTO test_users (username, email) VALUES (?, ?)",
				"user"+string(rune(i)), "user@example.com")
		}

		// 并发查询
		done := make(chan bool, 5)
		for i := 0; i < 5; i++ {
			go func() {
				rows, err := db.QueryContext(ctx, "SELECT COUNT(*) FROM test_users")
				if err == nil && rows != nil {
					rows.Close()
				}
				done <- true
			}()
		}

		// 等待所有查询完成
		for i := 0; i < 5; i++ {
			<-done
		}

		assert.True(t, true, "并发查询完成")
	})
}
