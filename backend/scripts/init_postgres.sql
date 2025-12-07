-- PostgreSQL 17 初始化脚本
-- Coze Studio Super 数据库初始化

-- 安装 pgvector 扩展（用于向量检索）
CREATE EXTENSION IF NOT EXISTS vector;

-- 创建数据库（如果不存在）
-- 注意：这个脚本应该在 postgres 数据库中执行
-- CREATE DATABASE coze_studio_db WITH ENCODING 'UTF8';

-- 连接到目标数据库后执行以下语句

-- 1. 安装必要的扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";     -- UUID 生成
CREATE EXTENSION IF NOT EXISTS "pg_trgm";       -- 文本相似度搜索
CREATE EXTENSION IF NOT EXISTS "btree_gin";     -- GIN 索引支持
CREATE EXTENSION IF NOT EXISTS "btree_gist";    -- GIST 索引支持

-- 2. 创建自定义函数（如需要）

-- 向量相似度搜索函数（使用 pgvector）
-- 已由 pgvector 扩展提供：<->, <#>, <=>

-- 3. 创建索引优化函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 4. 设置搜索路径
-- SET search_path TO public;

-- 5. 创建序列（GORM 会自动创建，这里仅作参考）
-- CREATE SEQUENCE IF NOT EXISTS users_id_seq;
-- CREATE SEQUENCE IF NOT EXISTS workspaces_id_seq;
-- ...

-- 6. 配置数据库参数（可选）
-- ALTER DATABASE coze_studio_db SET timezone TO 'Asia/Shanghai';
-- ALTER DATABASE coze_studio_db SET client_encoding TO 'UTF8';

-- 7. 创建只读用户（可选，用于报表查询）
-- CREATE USER coze_readonly WITH PASSWORD 'readonly_password';
-- GRANT CONNECT ON DATABASE coze_studio_db TO coze_readonly;
-- GRANT USAGE ON SCHEMA public TO coze_readonly;
-- GRANT SELECT ON ALL TABLES IN SCHEMA public TO coze_readonly;

-- 8. 注释
COMMENT ON DATABASE coze_studio_db IS 'Coze Studio Super 融合项目数据库';

-- 完成
SELECT 'PostgreSQL 17 初始化完成' AS status;

