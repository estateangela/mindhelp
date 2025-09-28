-- 新增諮商師相關資料表
-- 創建時間: 2024-01-03
-- 描述: 新增諮商師、諮商所、推薦醫師等資料表

-- 創建諮商師表
CREATE TABLE IF NOT EXISTS counselors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    license_number VARCHAR(50) UNIQUE NOT NULL,
    gender VARCHAR(10),
    specialties TEXT,
    language_skills TEXT,
    work_location VARCHAR(255),
    work_unit VARCHAR(255),
    institution_code VARCHAR(50),
    psychology_school VARCHAR(255),
    treatment_methods TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 創建諮商所表
CREATE TABLE IF NOT EXISTS counseling_centers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address VARCHAR(500),
    phone VARCHAR(50),
    online_counseling BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 創建推薦醫師表
CREATE TABLE IF NOT EXISTS recommended_doctors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    experience_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 創建索引
CREATE INDEX IF NOT EXISTS idx_counselors_name ON counselors(name);
CREATE INDEX IF NOT EXISTS idx_counselors_license_number ON counselors(license_number);
CREATE INDEX IF NOT EXISTS idx_counselors_work_location ON counselors(work_location);
CREATE INDEX IF NOT EXISTS idx_counselors_work_unit ON counselors(work_unit);
CREATE INDEX IF NOT EXISTS idx_counselors_institution_code ON counselors(institution_code);
CREATE INDEX IF NOT EXISTS idx_counselors_deleted_at ON counselors(deleted_at);

CREATE INDEX IF NOT EXISTS idx_counseling_centers_name ON counseling_centers(name);
CREATE INDEX IF NOT EXISTS idx_counseling_centers_online_counseling ON counseling_centers(online_counseling);
CREATE INDEX IF NOT EXISTS idx_counseling_centers_deleted_at ON counseling_centers(deleted_at);

CREATE INDEX IF NOT EXISTS idx_recommended_doctors_name ON recommended_doctors(name);
CREATE INDEX IF NOT EXISTS idx_recommended_doctors_experience_count ON recommended_doctors(experience_count);
CREATE INDEX IF NOT EXISTS idx_recommended_doctors_deleted_at ON recommended_doctors(deleted_at);

-- 創建觸發器來處理 updated_at 欄位
CREATE TRIGGER update_counselors_updated_at BEFORE UPDATE ON counselors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_counseling_centers_updated_at BEFORE UPDATE ON counseling_centers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_recommended_doctors_updated_at BEFORE UPDATE ON recommended_doctors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
