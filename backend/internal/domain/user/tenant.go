package user

import (
	"time"
)

// Tenant 表示租户实体
type Tenant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Domain    *string   `json:"domain,omitempty" db:"domain"`
	Status    string    `json:"status" db:"status"`
	MaxUsers  int       `json:"max_users" db:"max_users"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TenantStatus 表示租户状态
const (
	TenantStatusActive   = "active"
	TenantStatusInactive = "inactive"
)

// CreateTenantRequest 表示创建租户的请求
type CreateTenantRequest struct {
	Name     string  `json:"name" validate:"required"`
	Domain   *string `json:"domain,omitempty"`
	MaxUsers int     `json:"max_users" validate:"required,min=1"`
}

// UpdateTenantRequest 表示更新租户的请求
type UpdateTenantRequest struct {
	Name     *string `json:"name,omitempty"`
	Domain   *string `json:"domain,omitempty"`
	Status   *string `json:"status,omitempty"`
	MaxUsers *int    `json:"max_users,omitempty" validate:"omitempty,min=1"`
}

// TenantRepository 表示租户仓库接口
type TenantRepository interface {
	Create(tenant *Tenant) error
	GetByID(id string) (*Tenant, error)
	GetByDomain(domain string) (*Tenant, error)
	Update(tenant *Tenant) error
	Delete(id string) error
	List(offset, limit int) ([]*Tenant, int, error)
}

// TenantService 表示租户服务接口
type TenantService interface {
	Create(req *CreateTenantRequest) (*Tenant, error)
	GetByID(id string) (*Tenant, error)
	GetByDomain(domain string) (*Tenant, error)
	Update(id string, req *UpdateTenantRequest) (*Tenant, error)
	Delete(id string) error
	List(page, pageSize int) ([]*Tenant, int, error)
}
