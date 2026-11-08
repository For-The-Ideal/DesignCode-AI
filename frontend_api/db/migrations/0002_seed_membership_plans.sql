-- ══════════════════════════════════════════════════════════════════════
--  0002 会员套餐 / 积分包初始数据
--
--  背景：membership_plans、credit_packages 由 AutoMigrate 建表，
--        但代码里从未灌过初始数据，导致：
--          GET  /api/v1/membership/plans        → plans / packages 全空
--          POST /api/v1/membership/upgrade      → 404 套餐不存在或已下架
--          POST /api/v1/membership/buy-credits  → 400 积分套餐不存在
--
--  幂等：membership_plans 用 level 唯一键 ON DUPLICATE KEY UPDATE；
--        credit_packages 无唯一键，按 name 做 NOT EXISTS 去重。
--
--  自动初始化：服务启动时 repository.MembershipRepository.SeedIfEmpty
--        会在表为空时灌入同一份数据；本文件用于生产库手工执行/核对。
--
--  level 必须与 internal/handler/membership_handler.go 的 planFeatures 对应：
--    0=免费版(100 次/月)  1=专业版(500 次/月)  2=旗舰版(2000 次/月)
-- ══════════════════════════════════════════════════════════════════════

INSERT INTO membership_plans
  (created_at, updated_at, deleted_at, level, name, credits_per_month, price, original_price, is_active, sort_order)
VALUES
  (NOW(3), NOW(3), NULL, 0, '免费版', 100,  0,   0,   1, 0),
  (NOW(3), NOW(3), NULL, 1, '专业版', 500,  29,  59,  1, 1),
  (NOW(3), NOW(3), NULL, 2, '旗舰版', 2000, 99,  199, 1, 2)
ON DUPLICATE KEY UPDATE
  deleted_at        = NULL,
  name              = VALUES(name),
  credits_per_month = VALUES(credits_per_month),
  price             = VALUES(price),
  original_price    = VALUES(original_price),
  is_active         = VALUES(is_active),
  sort_order        = VALUES(sort_order),
  updated_at        = NOW(3);

INSERT INTO credit_packages
  (created_at, updated_at, deleted_at, name, credits, price, original_price)
SELECT NOW(3), NOW(3), NULL, t.name, t.credits, t.price, t.original_price
FROM (
  SELECT '100 次加油包'  AS name, 100  AS credits, 9   AS price, 12  AS original_price
  UNION ALL
  SELECT '500 次加油包',        500,               39,          59
  UNION ALL
  SELECT '2000 次加油包',       2000,              128,         199
) AS t
WHERE NOT EXISTS (
  SELECT 1 FROM credit_packages p
  WHERE p.name = t.name AND p.deleted_at IS NULL
);
