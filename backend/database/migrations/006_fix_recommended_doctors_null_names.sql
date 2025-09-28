-- 修復 recommended_doctors 表中 name 欄位的 null 值問題
-- 先將所有 null 的 name 設為預設值，然後再加上 NOT NULL 約束

BEGIN;

-- 1. 將所有 name 為 null 的記錄更新為預設值
UPDATE recommended_doctors 
SET name = COALESCE(name, '未知醫師') 
WHERE name IS NULL OR name = '';

-- 2. 確保 name 欄位不為空字串
UPDATE recommended_doctors 
SET name = '未知醫師' 
WHERE name = '';

-- 3. 重新應用 NOT NULL 約束
ALTER TABLE recommended_doctors 
ALTER COLUMN name SET NOT NULL;

COMMIT;
