package chat

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// MessageRepository 实现了 chat.MessageRepository 接口
type MessageRepository struct {
	db     *sqlx.DB
	logger *logger.Logger
}

// NewMessageRepository 创建一个新的 MessageRepository 实例
func NewMessageRepository(db *sqlx.DB, logger *logger.Logger) *MessageRepository {
	return &MessageRepository{
		db:     db,
		logger: logger,
	}
}

// Create 创建一个新的消息
func (r *MessageRepository) Create(message *chat.Message) error {
	// 将 metadata 转换为 JSON 字符串
	var metadataJSON sql.NullString
	if message.Metadata != nil {
		metadataBytes, err := json.Marshal(message.Metadata)
		if err != nil {
			r.logger.Error("序列化消息元数据失败", err)
			return fmt.Errorf("序列化消息元数据失败: %w", err)
		}
		metadataJSON = sql.NullString{
			String: string(metadataBytes),
			Valid:  true,
		}
	}

	query := `
		INSERT INTO messages (
			id, canvas_id, parent_id, role, content, metadata, created_at
		) VALUES (
			:id, :canvas_id, :parent_id, :role, :content, :metadata, :created_at
		)
	`

	// 创建一个包含 SQL 兼容字段的匿名结构体
	params := struct {
		ID        string         `db:"id"`
		CanvasID  string         `db:"canvas_id"`
		ParentID  sql.NullString `db:"parent_id"`
		Role      string         `db:"role"`
		Content   string         `db:"content"`
		Metadata  sql.NullString `db:"metadata"`
		CreatedAt string         `db:"created_at"`
	}{
		ID:        message.ID,
		CanvasID:  message.CanvasID,
		Role:      message.Role,
		Content:   message.Content,
		Metadata:  metadataJSON,
		CreatedAt: message.CreatedAt,
	}

	// 处理可能为空的 ParentID
	if message.ParentID != "" {
		params.ParentID = sql.NullString{
			String: message.ParentID,
			Valid:  true,
		}
	}

	_, err := r.db.NamedExec(query, params)
	if err != nil {
		r.logger.Error("创建消息失败", err)
		return fmt.Errorf("创建消息失败: %w", err)
	}

	return nil
}

// GetByID 根据 ID 获取消息
func (r *MessageRepository) GetByID(id string) (*chat.Message, error) {
	query := `
		SELECT id, canvas_id, parent_id, role, content, metadata, created_at
		FROM messages
		WHERE id = $1
	`

	var result struct {
		ID        string         `db:"id"`
		CanvasID  string         `db:"canvas_id"`
		ParentID  sql.NullString `db:"parent_id"`
		Role      string         `db:"role"`
		Content   string         `db:"content"`
		Metadata  sql.NullString `db:"metadata"`
		CreatedAt string         `db:"created_at"`
	}

	err := r.db.Get(&result, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("消息不存在: %s", id)
		}
		r.logger.Error("获取消息失败", err)
		return nil, fmt.Errorf("获取消息失败: %w", err)
	}

	// 构建消息对象
	message := &chat.Message{
		ID:        result.ID,
		CanvasID:  result.CanvasID,
		Role:      result.Role,
		Content:   result.Content,
		CreatedAt: result.CreatedAt,
	}

	// 处理可能为空的字段
	if result.ParentID.Valid {
		message.ParentID = result.ParentID.String
	}

	// 解析元数据
	if result.Metadata.Valid {
		var metadata map[string]interface{}
		err = json.Unmarshal([]byte(result.Metadata.String), &metadata)
		if err != nil {
			r.logger.Error("解析消息元数据失败", err)
			return nil, fmt.Errorf("解析消息元数据失败: %w", err)
		}
		message.Metadata = metadata
	}

	return message, nil
}

// ListByCanvasID 获取画布下的消息列表
func (r *MessageRepository) ListByCanvasID(canvasID string, offset, limit int) ([]*chat.Message, int, error) {
	// 查询总数
	countQuery := `SELECT COUNT(*) FROM messages WHERE canvas_id = $1`
	var total int
	err := r.db.Get(&total, countQuery, canvasID)
	if err != nil {
		r.logger.Error("获取消息总数失败", err)
		return nil, 0, fmt.Errorf("获取消息总数失败: %w", err)
	}

	// 查询数据
	query := `
		SELECT id, canvas_id, parent_id, role, content, metadata, created_at
		FROM messages
		WHERE canvas_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Queryx(query, canvasID, limit, offset)
	if err != nil {
		r.logger.Error("获取消息列表失败", err)
		return nil, 0, fmt.Errorf("获取消息列表失败: %w", err)
	}
	defer rows.Close()

	var messages []*chat.Message
	for rows.Next() {
		var result struct {
			ID        string         `db:"id"`
			CanvasID  string         `db:"canvas_id"`
			ParentID  sql.NullString `db:"parent_id"`
			Role      string         `db:"role"`
			Content   string         `db:"content"`
			Metadata  sql.NullString `db:"metadata"`
			CreatedAt string         `db:"created_at"`
		}

		err := rows.StructScan(&result)
		if err != nil {
			r.logger.Error("扫描消息数据失败", err)
			return nil, 0, fmt.Errorf("扫描消息数据失败: %w", err)
		}

		// 构建消息对象
		message := &chat.Message{
			ID:        result.ID,
			CanvasID:  result.CanvasID,
			Role:      result.Role,
			Content:   result.Content,
			CreatedAt: result.CreatedAt,
		}

		// 处理可能为空的字段
		if result.ParentID.Valid {
			message.ParentID = result.ParentID.String
		}

		// 解析元数据
		if result.Metadata.Valid {
			var metadata map[string]interface{}
			err = json.Unmarshal([]byte(result.Metadata.String), &metadata)
			if err != nil {
				r.logger.Error("解析消息元数据失败", err)
				return nil, 0, fmt.Errorf("解析消息元数据失败: %w", err)
			}
			message.Metadata = metadata
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("迭代消息行失败", err)
		return nil, 0, fmt.Errorf("迭代消息行失败: %w", err)
	}

	return messages, total, nil
}
