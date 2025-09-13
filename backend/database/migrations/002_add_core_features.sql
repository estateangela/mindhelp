-- 新增核心功能資料表
-- 創建時間: 2024-01-02
-- 描述: 新增專家文章、心理測驗、評論、收藏、通知、聊天會話等功能

-- 創建專家文章表
CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(200) NOT NULL,
    author VARCHAR(100) NOT NULL,
    author_title VARCHAR(100),
    publish_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    summary VARCHAR(500),
    content TEXT NOT NULL,
    tags TEXT, -- JSON array stored as text
    image_url VARCHAR(255),
    is_published BOOLEAN DEFAULT TRUE,
    view_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 創建心理測驗表
CREATE TABLE IF NOT EXISTS quizzes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50), -- anxiety, depression, etc.
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 創建測驗題目表
CREATE TABLE IF NOT EXISTS quiz_questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id UUID NOT NULL,
    question TEXT NOT NULL,
    options TEXT, -- JSON array stored as text
    order_num INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE
);

-- 創建測驗提交表
CREATE TABLE IF NOT EXISTS quiz_submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    quiz_id UUID NOT NULL,
    answers TEXT NOT NULL, -- JSON object stored as text
    score INTEGER NOT NULL,
    result TEXT,
    completed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE
);

-- 創建評論表
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    resource_id UUID NOT NULL, -- Location ID
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment VARCHAR(1000),
    is_helpful INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES locations(id) ON DELETE CASCADE,
    UNIQUE(user_id, resource_id) -- 每個使用者只能對同一資源評論一次
);

-- 創建收藏表
CREATE TABLE IF NOT EXISTS bookmarks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    resource_type VARCHAR(20) NOT NULL, -- article, location
    article_id UUID,
    location_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE,
    UNIQUE(user_id, resource_type, article_id), -- 防止重複收藏文章
    UNIQUE(user_id, resource_type, location_id) -- 防止重複收藏地點
);

-- 創建通知表
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL, -- NEW_ARTICLE, PROMOTION, SYSTEM
    title VARCHAR(200) NOT NULL,
    body VARCHAR(500),
    is_read BOOLEAN DEFAULT FALSE,
    payload TEXT, -- JSON stored as text
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 創建聊天會話表
CREATE TABLE IF NOT EXISTS chat_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    title VARCHAR(200),
    first_message_snippet VARCHAR(100),
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 創建使用者設定表
CREATE TABLE IF NOT EXISTS user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    notify_new_article BOOLEAN DEFAULT TRUE,
    notify_promotions BOOLEAN DEFAULT FALSE,
    notify_system_updates BOOLEAN DEFAULT TRUE,
    push_token VARCHAR(255),
    platform VARCHAR(10), -- ios, android
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 創建應用程式配置表
CREATE TABLE IF NOT EXISTS app_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(50) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    description VARCHAR(200),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 為現有的 chat_messages 表添加 session_id 欄位
ALTER TABLE chat_messages 
ADD COLUMN IF NOT EXISTS session_id UUID;

-- 添加外鍵約束
ALTER TABLE chat_messages 
ADD CONSTRAINT fk_chat_messages_session_id 
FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE SET NULL;

-- 創建索引
CREATE INDEX IF NOT EXISTS idx_articles_title ON articles(title);
CREATE INDEX IF NOT EXISTS idx_articles_author ON articles(author);
CREATE INDEX IF NOT EXISTS idx_articles_category ON articles USING gin(to_tsvector('english', title || ' ' || content));
CREATE INDEX IF NOT EXISTS idx_articles_publish_date ON articles(publish_date);
CREATE INDEX IF NOT EXISTS idx_articles_is_published ON articles(is_published);
CREATE INDEX IF NOT EXISTS idx_articles_deleted_at ON articles(deleted_at);

CREATE INDEX IF NOT EXISTS idx_quizzes_category ON quizzes(category);
CREATE INDEX IF NOT EXISTS idx_quizzes_is_active ON quizzes(is_active);
CREATE INDEX IF NOT EXISTS idx_quizzes_deleted_at ON quizzes(deleted_at);

CREATE INDEX IF NOT EXISTS idx_quiz_questions_quiz_id ON quiz_questions(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_order ON quiz_questions(quiz_id, order_num);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_deleted_at ON quiz_questions(deleted_at);

CREATE INDEX IF NOT EXISTS idx_quiz_submissions_user_id ON quiz_submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_quiz_submissions_quiz_id ON quiz_submissions(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_submissions_completed_at ON quiz_submissions(completed_at);
CREATE INDEX IF NOT EXISTS idx_quiz_submissions_deleted_at ON quiz_submissions(deleted_at);

CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_resource_id ON reviews(resource_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
CREATE INDEX IF NOT EXISTS idx_reviews_deleted_at ON reviews(deleted_at);

CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_resource_type ON bookmarks(resource_type);
CREATE INDEX IF NOT EXISTS idx_bookmarks_article_id ON bookmarks(article_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_location_id ON bookmarks(location_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_deleted_at ON bookmarks(deleted_at);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_deleted_at ON notifications(deleted_at);

CREATE INDEX IF NOT EXISTS idx_chat_sessions_user_id ON chat_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_last_updated_at ON chat_sessions(last_updated_at);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_is_active ON chat_sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_deleted_at ON chat_sessions(deleted_at);

CREATE INDEX IF NOT EXISTS idx_user_settings_user_id ON user_settings(user_id);
CREATE INDEX IF NOT EXISTS idx_user_settings_deleted_at ON user_settings(deleted_at);

CREATE INDEX IF NOT EXISTS idx_app_configs_key ON app_configs(key);
CREATE INDEX IF NOT EXISTS idx_app_configs_is_active ON app_configs(is_active);
CREATE INDEX IF NOT EXISTS idx_app_configs_deleted_at ON app_configs(deleted_at);

CREATE INDEX IF NOT EXISTS idx_chat_messages_session_id ON chat_messages(session_id);

-- 更新觸發器來處理新表的 updated_at 欄位
CREATE TRIGGER update_articles_updated_at BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_quizzes_updated_at BEFORE UPDATE ON quizzes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_quiz_questions_updated_at BEFORE UPDATE ON quiz_questions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_quiz_submissions_updated_at BEFORE UPDATE ON quiz_submissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookmarks_updated_at BEFORE UPDATE ON bookmarks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_chat_sessions_updated_at BEFORE UPDATE ON chat_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_app_configs_updated_at BEFORE UPDATE ON app_configs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 插入初始配置資料
INSERT INTO app_configs (key, value, description, is_active) VALUES 
('resource_types', '[{"key": "clinic", "displayName": "身心科診所"}, {"key": "counseling_center", "displayName": "心理諮商所"}, {"key": "free_service", "displayName": "免費諮詢服務"}, {"key": "clinical_psychology", "displayName": "臨床心理科"}]', '資源類型配置', true),
('specialties', '[{"key": "CBT", "displayName": "認知行為治療"}, {"key": "ADHD", "displayName": "注意力不足過動症"}, {"key": "anxiety", "displayName": "焦慮症"}, {"key": "depression", "displayName": "憂鬱症"}, {"key": "trauma", "displayName": "創傷治療"}, {"key": "family_therapy", "displayName": "家庭治療"}, {"key": "child_psychology", "displayName": "兒童心理"}]', '專業領域配置', true),
('features', '{"enableReviews": true, "enableTherapistProfiles": false, "enableGroupChat": false, "enableVideoConsult": false}', '功能開關配置', true)
ON CONFLICT (key) DO NOTHING;

-- 插入範例心理測驗資料
INSERT INTO quizzes (title, description, category, is_active) VALUES 
('GAD-7 焦慮自評量表', '廣泛性焦慮症七項量表，用於評估焦慮程度', 'anxiety', true),
('PHQ-9 憂鬱自評量表', '病人健康問卷九項量表，用於評估憂鬱程度', 'depression', true),
('PSS 壓力感知量表', '感知壓力量表，評估個人對壓力的感知程度', 'stress', true)
ON CONFLICT DO NOTHING;

-- 插入範例文章資料
INSERT INTO articles (title, author, author_title, summary, content, tags, is_published) VALUES 
('如何幫助焦慮的朋友', '王心理師', '臨床心理師', '當身邊的朋友出現焦慮症狀時，我們可以如何提供適當的支持和幫助', 
'<h2>理解焦慮</h2><p>焦慮是一種正常的情緒反應，但當它變得過度或持續時，就可能影響日常生活...</p>
<h2>如何提供支持</h2><p>1. 聆聽而不批判<br>2. 鼓勵尋求專業協助<br>3. 陪伴就醫...</p>', 
'["焦慮", "朋友", "支持", "陪伴"]', true),

('認識憂鬱症', '李醫師', '精神科醫師', '深入了解憂鬱症的症狀、成因以及治療方式', 
'<h2>什麼是憂鬱症</h2><p>憂鬱症是一種常見的精神疾病，影響著全球數億人...</p>
<h2>症狀識別</h2><p>持續的低落情緒、失去興趣、睡眠問題...</p>', 
'["憂鬱症", "症狀", "治療", "心理健康"]', true),

('壓力管理的基本技巧', '陳諮商師', '諮商心理師', '學習有效的壓力管理技巧，提升生活品質', 
'<h2>認識壓力</h2><p>壓力是現代生活中不可避免的一部分...</p>
<h2>管理技巧</h2><p>1. 深呼吸練習<br>2. 正念冥想<br>3. 時間管理...</p>', 
'["壓力管理", "放鬆技巧", "正念", "生活品質"]', true)
ON CONFLICT DO NOTHING;
