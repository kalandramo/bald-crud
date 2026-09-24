# 数据库迁移

本目录存放跨版本升级所需的 schema 迁移脚本与操作说明。

## 为什么需要迁移

`bald-crud` 的 `TenantID` mixin 在 **viewer/v0.2.0** 中把列类型由**数值**改为**字符串**：

| 后端 | v0.1.0（旧） | v0.2.0（新） |
|---|---|---|
| gorm（postgres/mysql/sqlite） | `int unsigned` | `varchar(64)` |
| clickhouse / doris | `UInt32` | `String` / `VARCHAR(64)` |
| mongodb | int32 / int64 | string |

**变更原因（安全修复）**：旧版通过 `strconv.ParseUint` 把租户 ID 桥接为数值，非数字租户 ID（如 `"t-default"`）**解析失败后静默留 0**；而 `IsPlatformContext()` 的判据是 `TenantIDValue == 0` —— 这类租户因此被误判为「平台视图」，`EnforceTenant` 放行，**租户隔离静默失效**。统一为 string 后无解析失败路径。

## 你需要迁移吗

**自检**（一条命令）：

```bash
grep -rn "bald-crud/.*mixin" --include="*.go" .
```

- **有输出** → 你的实体嵌入了 mixin，**需要迁移**（见下）
- **无输出** → 未使用 mixin，**无需迁移**

典型需要迁移的模型长这样：

```go
import "github.com/kalandramo/bald-crud/gorm/mixin"

type Order struct {
    mixin.TenantID        // ← 嵌入 mixin 即需要迁移
    ID string
}
```

**不需要迁移**的情形：模型自己声明 `TenantID string`（不依赖 mixin）。例如 bald-admin 的全部模型都是这种写法，其 `tenant_id` 列本就是字符串，**无需任何 schema 变更**。

## 数据前提

旧列是数值类型，**库中不可能存在非数字租户 ID**——能被写入的值必定是数字或 NULL。故本迁移是**纯类型转换**：`101` → `'101'`，无需值映射表。

若你的业务租户 ID 是 `"t-default"` 这类字符串，说明此前**从未成功写入过** `tenant_id`（被数值列拒绝或存为 0）。**执行迁移前请先核对存量数据的正确性**——迁移只改类型，不会凭空造出正确的租户归属。

## 各后端迁移步骤

### PostgreSQL

```bash
# 1. 先看现状
psql "$DSN" -c "SELECT table_name, data_type FROM information_schema.columns
                 WHERE column_name = 'tenant_id' AND table_schema = current_schema();"

# 2. 编辑脚本顶部的 tables 数组，列出你的实际表名
#    （每个嵌入 mixin 的模型对应一张表）

# 3. 备份（强烈建议）
pg_dump "$DSN" -t orders -t products > backup_$(date +%Y%m%d).sql

# 4. 执行
psql "$DSN" -f v0.2.0-tenant_id-to-string.sql

# 5. 验证：应全部为 character varying
psql "$DSN" -c "SELECT table_name, data_type FROM information_schema.columns
                 WHERE column_name = 'tenant_id' AND table_schema = current_schema();"
```

脚本用 `DO $$ ... $$` 块 + `information_schema` 守卫，**幂等**——重复执行会跳过已是 `varchar` 的列。

### MySQL / MariaDB

见 `v0.2.0-tenant_id-to-string.sql` 中 MySQL 段的注释（提供两条路径：prepared statement 幂等版，或直接逐表 `MODIFY COLUMN`）。

```bash
mysqldump -h HOST -P PORT -u USER -p DBNAME orders products > backup.sql
mysql -h HOST -P PORT -u USER -p DBNAME < v0.2.0-tenant_id-to-string.sql
```

### SQLite

> ⚠️ **破坏性操作**：SQLite 不支持 `ALTER COLUMN TYPE`，只能「建新表 → 拷数据 → 删旧表 → 换名」。**删除原表不可逆**。本步骤需人工逐步确认，脚本文件里刻意**不提供**这些语句。

**推荐路径**：SQLite 通常用于开发/测试，若数据可重建，**直接删库重跑**即可（`database.sql.migrate: true` 时 `AutoMigrate` 会自动建表）。

**仅当库中有不可再生数据时**，按以下步骤（每步确认无误再执行下一步）：

```bash
# 1. 备份（必须）
cp your.db your.db.bak

# 2. 查看原表完整结构（务必原样照抄到新表）
sqlite3 your.db ".schema orders"

# 3. 建新表：结构与原表一致，仅 tenant_id 改为 TEXT
sqlite3 your.db "CREATE TABLE orders_new (...);"

# 4. 拷数据（CAST 完成类型转换）
sqlite3 your.db "INSERT INTO orders_new SELECT ... CAST(tenant_id AS TEXT) ... FROM orders;"

# 5. 核对行数一致
sqlite3 your.db "SELECT (SELECT COUNT(*) FROM orders), (SELECT COUNT(*) FROM orders_new);"

# 6. 确认无误后，才进行换名（含删表，不可逆）
#    建议此步前再备份一次：cp your.db your.db.bak2
```

### ClickHouse / Doris

两者是 OLAP 引擎，改列类型语法与 OLTP 不同：

```sql
ALTER TABLE orders MODIFY COLUMN tenant_id String;        -- ClickHouse
ALTER TABLE orders MODIFY COLUMN tenant_id VARCHAR(64);   -- Doris
```

**注意**：若 `tenant_id` 参与**分桶键**（`DISTRIBUTED BY` / `BUCKETS`）或**分区键**，改类型需**重建表**。执行前先确认表定义（`SHOW CREATE TABLE orders`）。

### MongoDB

无 schema 约束，需逐文档转换类型：

```bash
mongosh "$URI" --eval '
  db.orders.updateMany(
    { tenant_id: { $type: ["int", "long", "double", "decimal"] } },
    [{ $set: { tenant_id: { $toString: "$tenant_id" } } }]
  );
'
```

`$type` 守卫保证幂等（已是 string 的文档不重复处理）。验证：

```bash
mongosh "$URI" --eval 'db.orders.countDocuments({ tenant_id: { $type: ["int", "long"] } })'  # 应为 0
```

## 迁移后验证清单

1. **列类型正确**：按上节各后端的验证语句确认
2. **数据无损**：行数与迁移前一致
3. **应用可读写**：启动服务，创建一条带租户的数据，确认写入的是字符串租户 ID
4. **隔离生效**：用租户 A 的凭据查询，确认看不到租户 B 的数据

## 回滚

**类型迁移不可自动回滚**——字符串转回数值会丢失非数字租户 ID（这正是本次修复要消除的问题）。

若必须回退，从**迁移前的备份**恢复：

```bash
psql "$DSN" < backup_YYYYMMDD.sql        # PostgreSQL
mysql ... < backup.sql                   # MySQL
cp your.db.bak your.db                   # SQLite
```

**代码侧回退**需同时降依赖：`go get github.com/kalandramo/bald-crud/viewer@v0.1.0`（注意：该版本存在前述租户隔离静默失效缺陷，不建议）。
