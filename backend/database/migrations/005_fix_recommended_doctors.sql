-- 檢查 recommended_doctors 表結構
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'recommended_doctors' 
ORDER BY ordinal_position;

-- 如果 name 欄位不存在，新增它
ALTER TABLE recommended_doctors 
ADD COLUMN IF NOT EXISTS name VARCHAR(255);

-- 更新現有記錄的 name 欄位
UPDATE recommended_doctors 
SET name = COALESCE(name, '未命名醫師') 
WHERE name IS NULL OR name = '';

-- 設定 name 欄位為 NOT NULL
ALTER TABLE recommended_doctors 
ALTER COLUMN name SET NOT NULL;

-- 為 name 欄位建立索引
CREATE INDEX IF NOT EXISTS idx_recommended_doctors_name 
ON recommended_doctors (name);

-- 檢查修正後的結果
SELECT COUNT(*) as total_records, 
       COUNT(CASE WHEN name IS NOT NULL AND name != '' THEN 1 END) as records_with_name
FROM recommended_doctors;
