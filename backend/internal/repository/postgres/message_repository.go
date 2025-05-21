package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/pkg/db"
)

// MessageRepository 表示消息仓库
type MessageRepository struct {
	db *db.Postgres
}

// NewMessageRepository 创建一个新的消息仓库
func NewMessageRepository(db *db.Postgres) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

// Create 创建一个新消息
func (r *MessageRepository) Create(message *chat.Message) error {
	// 生成 UUID
	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	// 设置时间戳
	message.CreatedAt = time.Now()

	// 处理元数据
	var metadata []byte
	var err error
	if message.Metadata != nil {
		metadata = message.Metadata
	} else {
		metadata = []byte("{}")
	}

	query := `
		INSERT INTO messages (id, canvas_id, parent_id, role, content, metadata, token_count, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.DB.Exec(
		query,
		message.ID,
		message.CanvasID,
		message.ParentID,
		message.Role,
		message.Content,
		metadata,
		message.TokenCount,
		message.CreatedBy,
		message.CreatedAt,
	)

	return err
}

// GetByID 通过 ID 获取消息
func (r *MessageRepository) GetByID(id string) (*chat.Message, error) {
	query := `
		SELECT id, canvas_id, parent_id, role, content, metadata, token_count, created_by, created_at
		FROM messages
		WHERE id = $1
	`

	var message chat.Message
	err := r.db.DB.Get(&message, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("消息不存在: %w", err)
		}
		return nil, err
	}

	return &message, nil
}

// GetByCanvasID 获取画布的所有消息
func (r *MessageRepository) GetByCanvasID(canvasID string, offset, limit int) ([]*chat.Message, int, error) {
	// 获取总数
	countQuery := `
		SELECT COUNT(*)
		FROM messages
		WHERE canvas_id = $1
	`

	var total int
	err := r.db.DB.Get(&total, countQuery, canvasID)
	if err != nil {
		return nil, 0, err
	}

	// 获取消息列表
	query := `
		SELECT id, canvas_id, parent_id, role, content, metadata, token_count, created_by, created_at
		FROM messages
		WHERE canvas_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	var messages []*chat.Message
	err = r.db.DB.Select(&messages, query, canvasID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetConversation 获取对话历史
func (r *MessageRepository) GetConversation(messageID string, limit int) ([]*chat.Message, error) {
	// 首先获取当前消息
	currentMessage, err := r.GetByID(messageID)
	if err != nil {
		return nil, err
	}

	// 获取画布 ID
	canvasID := currentMessage.CanvasID

	// 获取对话历史
	query := `
		WITH RECURSIVE conversation AS (
			-- 基本情况：当前消息
			SELECT id, canvas_id, parent_id, role, content, metadata, token_count, created_by, created_at, 0 as depth
			FROM messages
			WHERE id = $1
			
			UNION ALL
			
			-- 递归情况：父消息
			SELECT m.id, m.canvas_id, m.parent_id, m.role, m.content, m.metadata, m.token_count, m.created_by, m.created_at, c.depth + 1
			FROM messages m
			JOIN conversation c ON m.id = c.parent_id
			WHERE m.canvas_id = $2
		)
		SELECT id, canvas_id, parent_id, role, content, metadata, token_count, created_by, created_at
		FROM conversation
		ORDER BY depth DESC, created_at ASC
		LIMIT $3
	`

	var messages []*chat.Message
	err = r.db.DB.Select(&messages, query, messageID, canvasID, limit)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// CreateBatch 批量创建消息
func (r *MessageRepository) CreateBatch(messages []*chat.Message) error {
	tx, err := r.db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO messages (id, canvas_id, parent_id, role, content, metadata, token_count, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	for _, message := range messages {
		// 生成 UUID
		if message.ID == "" {
			message.ID = uuid.New().String()
		}

		// 设置时间戳
		if message.CreatedAt.IsZero() {
			message.CreatedAt = time.Now()
		}

		// 处理元数据
		var metadata []byte
		if message.Metadata != nil {
			metadata = message.Metadata
		} else {
			metadata = []byte("{}")
		}

		_, err = tx.Exec(
			query,
			message.ID,
			message.CanvasID,
			message.ParentID,
			message.Role,
			message.Content,
			metadata,
			message.TokenCount,
			message.CreatedBy,
			message.CreatedAt,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
