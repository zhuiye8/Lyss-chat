package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/your-org/lyss-chat-backend/internal/domain/user"
	"github.com/your-org/lyss-chat-backend/pkg/db"
)

// UserRepository 表示用户仓库
type UserRepository struct {
	db *db.Postgres
}

// NewUserRepository 创建一个新的用户仓库
func NewUserRepository(db *db.Postgres) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create 创建一个新用户
func (r *UserRepository) Create(user *user.User) error {
	// 生成 UUID
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	query := `
		INSERT INTO users (id, tenant_id, email, password, name, avatar_url, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	_, err := r.db.DB.Exec(
		query,
		user.ID,
		user.TenantID,
		user.Email,
		user.Password,
		user.Name,
		user.AvatarURL,
		user.Status,
	)

	return err
}

// GetByID 通过 ID 获取用户
func (r *UserRepository) GetByID(id string) (*user.User, error) {
	query := `
		SELECT id, tenant_id, email, password, name, avatar_url, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u user.User
	err := r.db.DB.Get(&u, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("用户不存在: %w", err)
		}
		return nil, err
	}

	return &u, nil
}

// GetByEmail 通过邮箱获取用户
func (r *UserRepository) GetByEmail(email, tenantID string) (*user.User, error) {
	query := `
		SELECT id, tenant_id, email, password, name, avatar_url, status, created_at, updated_at
		FROM users
		WHERE email = $1 AND tenant_id = $2
	`

	var u user.User
	err := r.db.DB.Get(&u, query, email, tenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("用户不存在: %w", err)
		}
		return nil, err
	}

	return &u, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *user.User) error {
	query := `
		UPDATE users
		SET name = $1, avatar_url = $2, status = $3, updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.DB.Exec(
		query,
		user.Name,
		user.AvatarURL,
		user.Status,
		user.ID,
	)

	return err
}

// Delete 删除用户
func (r *UserRepository) Delete(id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.DB.Exec(query, id)
	return err
}

// List 列出用户
func (r *UserRepository) List(tenantID string, offset, limit int) ([]*user.User, int, error) {
	// 获取总数
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE tenant_id = $1
	`

	var total int
	err := r.db.DB.Get(&total, countQuery, tenantID)
	if err != nil {
		return nil, 0, err
	}

	// 获取用户列表
	query := `
		SELECT id, tenant_id, email, password, name, avatar_url, status, created_at, updated_at
		FROM users
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var users []*user.User
	err = r.db.DB.Select(&users, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
