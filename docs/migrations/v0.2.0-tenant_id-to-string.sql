-- ============================================================================
-- tenant_id 列类型迁移：数值 → 字符串
-- ============================================================================
--
-- 适用版本：bald-crud viewer/v0.2.0（含 gorm/clickhouse/doris/mongodb v0.2.0）
--
-- ## 为什么需要迁移
--
-- viewer/v0.1.0 的 mixin 把 tenant_id 定义为数值列（gorm: int unsigned；
-- clickhouse/doris: UInt32），而 bald 生态的租户 ID 是**字符串**
-- （contextx / authn / jwt / audit 全为 string，典型值如 "t-default"、"platform"）。
--
-- 旧版通过 strconv.ParseUint 桥接，非数字租户 ID **解析失败后静默留 0**；
-- 而 IsPlatformContext() 的判据是 TenantIDValue == 0 —— 于是这类租户被误判为
-- 「平台视图」，EnforceTenant 放行，**租户隔离静默失效**（安全缺陷）。
--
-- v0.2.0 把类型统一为 string，消除了解析失败路径。若你的库中 tenant_id 列
-- 仍是数值类型，升级代码后需执行本迁移——否则 ORM 写入字符串会失败或截断。
--
-- ## 何时需要执行
--
-- **仅当**你的实体嵌入了 bald-crud 的 TenantID mixin 时：
--
--     // gorm
--     import "github.com/kalandramo/bald-crud/gorm/mixin"
--     type Order struct {
--         mixin.TenantID        // ← 嵌入即需要
--         ID string
--     }
--
-- 若你的模型自己声明 `TenantID string`（如 bald-admin 的做法），**无需迁移**。
-- 自检：`grep -rn "bald-crud/.*mixin" --include="*.go" .`
--
-- ## 数据前提（重要）
--
-- 旧列是数值类型，因此**库中不可能存在非数字租户 ID**——能被写入的值
-- 必定是数字（或 NULL）。故本迁移是**纯类型转换**，不需要值映射表：
-- 数字 101 → 字符串 '101'。
--
-- 若你的业务租户 ID 是 "t-default" 这类字符串，说明你此前**从未成功写入过**
-- tenant_id（写入时被数值列拒绝或存为 0）——请先核对存量数据的正确性，
-- 再决定是否需要按业务规则回填。
--
-- ## 幂等性
--
-- 三个方言的脚本都用「先查类型再改」的守卫，重复执行无副作用。
-- ============================================================================


-- ============================================================================
-- PostgreSQL
-- ============================================================================
-- 使用：psql "$DSN" -f migrate_tenant_id_to_string.sql
--     或 psql -h HOST -p PORT -U USER -d DBNAME -f migrate_tenant_id_to_string.sql

BEGIN;

-- 通用：把指定表的 tenant_id 由数值改为 varchar(64)。
-- information_schema 守卫保证幂等（已是 varchar 则跳过）。
DO $$
DECLARE
    tbl text;
    cur_type text;
    -- 按需增删你的业务表名。mixin 嵌入的每张表都要列出。
    tables text[] := ARRAY[
        'orders',        -- 示例：替换为你的实际表名
        'products'
    ];
BEGIN
    FOREACH tbl IN ARRAY tables LOOP
        -- 跳过不存在的表（多环境共用脚本时的容错）
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = current_schema() AND table_name = tbl
        ) THEN
            RAISE NOTICE '跳过：表 % 不存在', tbl;
            CONTINUE;
        END IF;

        SELECT data_type INTO cur_type
        FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = tbl
          AND column_name = 'tenant_id';

        IF cur_type IS NULL THEN
            RAISE NOTICE '跳过：表 % 无 tenant_id 列', tbl;
            CONTINUE;
        END IF;

        IF cur_type = 'character varying' THEN
            RAISE NOTICE '跳过：表 %.tenant_id 已是 varchar', tbl;
            CONTINUE;
        END IF;

        -- USING 子句完成数值 → 文本转换（NULL 保持 NULL）
        EXECUTE format(
            'ALTER TABLE %I ALTER COLUMN tenant_id TYPE varchar(64) USING tenant_id::text',
            tbl
        );
        RAISE NOTICE '已迁移：%.tenant_id (% → varchar)', tbl, cur_type;
    END LOOP;
END $$;

COMMIT;

-- 验证：应全部显示 character varying
-- SELECT table_name, data_type FROM information_schema.columns
--  WHERE column_name = 'tenant_id' AND table_schema = current_schema();


-- ============================================================================
-- MySQL / MariaDB
-- ============================================================================
-- 使用：mysql -h HOST -P PORT -u USER -p DBNAME < migrate_tenant_id_to_string.sql
--
-- 注意：MySQL 无 DO $$ 块，幂等守卫用 prepared statement + information_schema。
-- 把下面两条 SET 的表名替换为你的实际表。

-- SET @tbl = 'orders';
-- SET @sql = (
--   SELECT IF(
--     (SELECT DATA_TYPE FROM information_schema.COLUMNS
--       WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = @tbl
--         AND COLUMN_NAME = 'tenant_id') IN ('varchar', 'char', 'text'),
--     'SELECT ''跳过：已是字符串类型'' AS note',
--     CONCAT('ALTER TABLE `', @tbl, '` MODIFY COLUMN tenant_id VARCHAR(64)')
--   )
-- );
-- PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
--
-- 多表时重复上述 4 行（改 @tbl）。或直接逐表执行（幂等由 MODIFY 保证——
-- MySQL 的 MODIFY 重复执行同样类型是无害的）：
--
-- ALTER TABLE orders   MODIFY COLUMN tenant_id VARCHAR(64);
-- ALTER TABLE products MODIFY COLUMN tenant_id VARCHAR(64);


-- ============================================================================
-- SQLite
-- ============================================================================
-- SQLite 不支持 ALTER COLUMN TYPE，必须重建表——该操作**破坏性**（含删除原表），
-- 故本文件**不提供**可执行语句，避免 `sqlite3 db < 本文件` 或复制粘贴时误删表。
--
-- 完整步骤见：docs/migrations/README.md §SQLite 重建步骤（需人工逐步确认）
--
-- 推荐替代路径：SQLite 通常用于开发/测试，若数据可重建，直接删库重跑
-- AutoMigrate 即可（`database.sql.migrate: true` 会自动建表）。
-- 仅当 SQLite 中有不可再生的数据时，才需按文档执行重建。


-- ============================================================================
-- ClickHouse / Doris
-- ============================================================================
-- 两者都是 OLAP 引擎，改列类型的语法与 OLTP 不同，且通常按分区/分桶重建。
--
-- ClickHouse：ALTER TABLE ... MODIFY COLUMN tenant_id String
-- Doris：      ALTER TABLE ... MODIFY COLUMN tenant_id VARCHAR(64)
--
-- 注意：
--   - 若 tenant_id 参与分桶键（DISTRIBUTED BY / BUCKETS），改类型需重建表；
--   - 若为分区键，同理。请先确认表定义再执行。
--
-- ALTER TABLE orders MODIFY COLUMN tenant_id String;   -- ClickHouse
-- ALTER TABLE orders MODIFY COLUMN tenant_id VARCHAR(64); -- Doris


-- ============================================================================
-- MongoDB
-- ============================================================================
-- 无 schema 约束，需逐文档更新类型（数值 → 字符串）。
--
-- 使用：mongosh "$URI" --eval "$(cat migrate_tenant_id_to_string.js)"
--      或把下面内容存为 .js 后 mongosh "$URI" migrate_tenant_id_to_string.js
--
-- // 对每张集合执行；$type 守卫保证幂等（已是 string 的文档不重复处理）
-- db.orders.updateMany(
--   { tenant_id: { $type: ["int", "long", "double", "decimal"] } },
--   [{ $set: { tenant_id: { $toString: "$tenant_id" } } }]
-- );
--
-- 验证：
-- db.orders.countDocuments({ tenant_id: { $type: ["int", "long"] } })  // 应为 0
