-- 创建 canvases 表
CREATE TABLE IF NOT EXISTS canvases (
    id UUID PRIMARY KEY,
    workspace_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'chat',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    model_id UUID,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建 messages 表
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    canvas_id UUID NOT NULL REFERENCES canvases(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES messages(id),
    role VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    metadata JSONB,
    token_count INTEGER,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建 attachments 表
CREATE TABLE IF NOT EXISTS attachments (
    id UUID PRIMARY KEY,
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    size INTEGER NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    url VARCHAR(1024) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_canvases_workspace_id ON canvases(workspace_id);
CREATE INDEX idx_canvases_created_by ON canvases(created_by);
CREATE INDEX idx_canvases_status ON canvases(status);
CREATE INDEX idx_messages_canvas_id ON messages(canvas_id);
CREATE INDEX idx_messages_parent_id ON messages(parent_id);
CREATE INDEX idx_messages_created_by ON messages(created_by);
CREATE INDEX idx_attachments_message_id ON attachments(message_id);
