package chat

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// CanvasRepository 实现了 chat.CanvasRepository 接口
type CanvasRepository struct {
	db     *sqlx.DB
	logger *logger.Logger
}

// NewCanvasRepository 创建一个新的 CanvasRepository 实例
func NewCanvasRepository(db *sqlx.DB, logger *logger.Logger) *CanvasRepository {
	return &CanvasRepository{
		db:     db,
		logger: logger,
	}
}

// Create 创建一个新的画布
func (r *CanvasRepository) Create(canvas *chat.Canvas) error {
	query := `
		INSERT INTO canvases (
			id, workspace_id, title, description, type, status, model_id, created_by, created_at, updated_at
		) VALUES (
			:id, :workspace_id, :title, :description, :type, :status, :model_id, :created_by, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExec(query, canvas)
	if err != nil {
		r.logger.Error("创建画布失败", err)
		return fmt.Errorf("创建画布失败: %w", err)
	}

	return nil
}

// GetByID 根据 ID 获取画布
func (r *CanvasRepository) GetByID(id string) (*chat.Canvas, error) {
	query := `
		SELECT id, workspace_id, title, description, type, status, model_id, created_by, created_at, updated_at
		FROM canvases
		WHERE id = $1
	`

	var canvas chat.Canvas
	err := r.db.Get(&canvas, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("画布不存在: %s", id)
		}
		r.logger.Error("获取画布失败", err)
		return nil, fmt.Errorf("获取画布失败: %w", err)
	}

	return &canvas, nil
}

// Update 更新画布
func (r *CanvasRepository) Update(canvas *chat.Canvas) error {
	query := `
		UPDATE canvases
		SET title = :title,
			description = :description,
			status = :status,
			model_id = :model_id,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExec(query, canvas)
	if err != nil {
		r.logger.Error("更新画布失败", err)
		return fmt.Errorf("更新画布失败: %w", err)
	}

	return nil
}

// Delete 删除画布
func (r *CanvasRepository) Delete(id string) error {
	query := `DELETE FROM canvases WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("删除画布失败", err)
		return fmt.Errorf("删除画布失败: %w", err)
	}

	return nil
}

// List 列出工作区下的画布
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

	// 查询总数
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM canvases %s
	`, whereClause)

	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		r.logger.Error("获取画布总数失败", err)
		return nil, 0, fmt.Errorf("获取画布总数失败: %w", err)
	}

	// 查询数据
	query := fmt.Sprintf(`
		SELECT id, workspace_id, title, description, type, status, model_id, created_by, created_at, updated_at
		FROM canvases
		%s
		ORDER BY updated_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var canvases []*chat.Canvas
	err = r.db.Select(&canvases, query, args...)
	if err != nil {
		r.logger.Error("获取画布列表失败", err)
		return nil, 0, fmt.Errorf("获取画布列表失败: %w", err)
	}

	return canvases, total, nil
}
