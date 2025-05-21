-- 删除索引
DROP INDEX IF EXISTS idx_canvases_workspace_id;
DROP INDEX IF EXISTS idx_canvases_created_by;
DROP INDEX IF EXISTS idx_canvases_status;
DROP INDEX IF EXISTS idx_messages_canvas_id;
DROP INDEX IF EXISTS idx_messages_parent_id;
DROP INDEX IF EXISTS idx_messages_created_by;
DROP INDEX IF EXISTS idx_attachments_message_id;

-- 删除表
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS canvases;
