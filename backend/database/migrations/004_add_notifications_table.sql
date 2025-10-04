-- 創建通知表
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- 索引
    CONSTRAINT fk_notifications_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 創建索引
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_deleted_at ON notifications(deleted_at);

-- 添加註釋
COMMENT ON TABLE notifications IS '用戶通知表';
COMMENT ON COLUMN notifications.id IS '通知唯一標識符';
COMMENT ON COLUMN notifications.user_id IS '用戶ID';
COMMENT ON COLUMN notifications.title IS '通知標題';
COMMENT ON COLUMN notifications.content IS '通知內容';
COMMENT ON COLUMN notifications.type IS '通知類型：hourly_reminder, weekly_bulletin, system等';
COMMENT ON COLUMN notifications.is_read IS '是否已讀';
COMMENT ON COLUMN notifications.created_at IS '創建時間';
COMMENT ON COLUMN notifications.updated_at IS '更新時間';
COMMENT ON COLUMN notifications.deleted_at IS '軟刪除時間';
