-- 添加推薦醫師表
-- 創建時間: 2024-09-28
-- 描述: 創建推薦醫師資料表

-- 創建推薦醫師表
CREATE TABLE IF NOT EXISTS recommended_doctors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    specialty VARCHAR(255),
    description TEXT,
    hospital VARCHAR(255),
    location VARCHAR(255),
    contact VARCHAR(255),
    experience_years INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 創建索引
CREATE INDEX IF NOT EXISTS idx_recommended_doctors_name ON recommended_doctors(name);
CREATE INDEX IF NOT EXISTS idx_recommended_doctors_specialty ON recommended_doctors(specialty);
CREATE INDEX IF NOT EXISTS idx_recommended_doctors_is_active ON recommended_doctors(is_active);
CREATE INDEX IF NOT EXISTS idx_recommended_doctors_deleted_at ON recommended_doctors(deleted_at);

-- 創建觸發器
CREATE TRIGGER update_recommended_doctors_updated_at BEFORE UPDATE ON recommended_doctors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
