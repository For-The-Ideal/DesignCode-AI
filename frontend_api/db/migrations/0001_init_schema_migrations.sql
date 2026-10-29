-- ══════════════════════════════════════════════════════════════════════
--  0001 基线：迁移记账表
--
--  这张表由部署脚本在跑迁移前自动创建（IF NOT EXISTS），
--  这里显式声明一次，保证手工执行迁移时同样有表可用。
--
--  迁移规范（违反会导致回滚失效）：
--    · 只做加法：加表、加列、加索引；删列/改名必须拆成「先双写 → 再删」两期
--    · 必须是幂等的：用 IF NOT EXISTS / IF EXISTS 包住
--    · 文件一旦发布过就永不修改（checksum 会变，脚本会告警并跳过）
--    · 文件名格式 NNNN_描述.sql，按字典序执行
-- ══════════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS schema_migrations (
  version    VARCHAR(191) PRIMARY KEY,
  checksum   CHAR(64)    NOT NULL,
  applied_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
