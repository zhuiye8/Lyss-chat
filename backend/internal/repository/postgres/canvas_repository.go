package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/pkg/db"
)

// CanvasRepository 表示画布仓库
type CanvasRepository struct {
	db *db.Postgres
}

// NewCanvasRepository 创建一个新的画布仓库
func NewCanvasRepository(db *db.Postgres) *CanvasRepository {
	return &CanvasRepository{
		db: db,
	}
}

// Create 创建一个新画布
func (r *CanvasRepository) Create(canvas *chat.Canvas) error {
	// 生成 UUID
	if canvas.ID == "" {
		canvas.ID = uuid.New().String()
	}

	// 设置时间戳
	now := time.Now()
	canvas.CreatedAt = now
	canvas.UpdatedAt = now

	query := `
		INSERT INTO canvases (id, workspace_id, title, description, type, status, model_id, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.DB.Exec(
		query,
		canvas.ID,
		canvas.WorkspaceID,
		canvas.Title,
		canvas.Description,
		canvas.Type,
		canvas.Status,
		canvas.ModelID,
		canvas.CreatedBy,
		canvas.CreatedAt,
		canvas.UpdatedAt,
	)

	return err
}

// GetByID 通过 ID 获取画布
func (r *CanvasRepository) GetByID(id string) (*chat.Canvas, error) {
	query := `
		SELECT id, workspace_id, title, description, type, status, model_id, created_by, created_at, updated_at
		FROM canvases
		WHERE id = $1
	`

	var canvas chat.Canvas
	err := r.db.DB.Get(&canvas, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("画布不存在: %w", err)
		}
		return nil, err
	}

	return &canvas, nil
}

// Update 更新画布
func (r *CanvasRepository) Update(canvas *chat.Canvas) error {
	// 更新时间戳
	canvas.UpdatedAt = time.Now()

	query := `
		UPDATE canvases
		SET title = $1, description = $2, status = $3, model_id = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.DB.Exec(
		query,
		canvas.Title,
		canvas.Description,
		canvas.Status,
		canvas.ModelID,
		canvas.UpdatedAt,
		canvas.ID,
	)

	return err
}

// Delete 删除画布
func (r *CanvasRepository) Delete(id string) error {
	query := `
		DELETE FROM canvases
		WHERE id = $1
	`

	_, err := r.db.DB.Exec(query, id)
	return err
}

// List 列出画布
func (r *CanvasRepository) List(workspaceID string, canvasType *string, offset, limit int) ([]*chat.Canvas, int, error) {
	// 构建查询条件
	whereClause := "WHERE workspace_id = $1"
	args := []interface{}{workspaceID}
	argIndex := 2

	if canvasType != nil {
		whereClause += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *canvasType)
		argIndex++
	}

	// 获取总数
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM canvases
		%s
	`, whereClause)

	var total int
	err := r.db.DB.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 获取画布列表
	query := fmt.Sprintf(`
		SELECT id, workspace_id, title, description, type, status, model_id, created_by, created_at, updated_at
		FROM canvases
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var canvases []*chat.Canvas
	err = r.db.DB.Select(&canvases, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return canvases, total, nil
}
