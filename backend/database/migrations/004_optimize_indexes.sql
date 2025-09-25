-- 優化查詢效能的索引
-- 針對 counseling_centers 表

-- 為搜索欄位建立複合索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_counseling_centers_search 
ON counseling_centers USING gin(to_tsvector('chinese', COALESCE(name, '') || ' ' || COALESCE(address, '')));

-- 為 name 欄位建立索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_counseling_centers_name 
ON counseling_centers (name) WHERE name IS NOT NULL;

-- 為 address 欄位建立索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_counseling_centers_address 
ON counseling_centers (address) WHERE address IS NOT NULL;

-- 為 online_counseling 欄位建立索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_counseling_centers_online 
ON counseling_centers (online_counseling);

-- 為 deleted_at 欄位建立索引（軟刪除）
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_counseling_centers_deleted_at 
ON counseling_centers (deleted_at);

-- 修正 NULL 值問題
UPDATE counseling_centers 
SET name = '未命名機構' 
WHERE name IS NULL OR name = '';

-- 為 recommended_doctors 表建立索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_recommended_doctors_name 
ON recommended_doctors (name) WHERE name IS NOT NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_recommended_doctors_deleted_at 
ON recommended_doctors (deleted_at);

-- 修正 recommended_doctors 的 NULL 值
UPDATE recommended_doctors 
SET name = '未命名醫師' 
WHERE name IS NULL OR name = '';

UPDATE recommended_doctors 
SET description = '暫無描述' 
WHERE description IS NULL OR description = '';

-- 為 counselors 表建立索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_counselors_work_location 
ON counselors (work_location) WHERE work_location IS NOT NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_counselors_deleted_at 
ON counselors (deleted_at);
