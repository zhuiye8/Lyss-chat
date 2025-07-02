# LYSS AI 平台 - RBAC 权限系统设计 (V2)

**版本**: 2.0
**最后更新**: 2025年7月2日

---

## 1. 核心理念

本平台的 V2 权限系统 (RBAC) 严格围绕**供应商作用域 (`scope`)** 和**用户角色 (`role`)** 这两个核心概念构建。

1.  **角色作为用户属性**: 用户的角色 (`user`, `admin`, `super_admin`) 存储在 `users` 表中，决定其基础操作权限。
2.  **Scope 决定资产行为**: 供应商的 `scope` (`ORGANIZATION`, `PERSONAL`) 决定了其归属、可见性和核心能力（能否被分发）。
3.  **API 层多维验证**: API 端点将同时验证用户角色和（在适用时）资源的 `scope`，以实现精确的权限控制。

## 2. 角色与权限矩阵 (V2)

| 功能模块 | 角色 | 允许的 `scope` | 核心权限描述 |
| :--- | :--- | :--- | :--- |
| **创建供应商** | `admin` | `ORGANIZATION`, `PERSONAL` | 管理员可以为组织创建可分发的供应商，也可以为自己创建个人供应商。 |
| | `user` | `PERSONAL` | 普通用户只能为自己创建个人供应商，不可分发。 |
| **查看供应商** | `admin` | `ORGANIZATION`, `PERSONAL` | 管���员可以看到所有组织供应商和自己创建的个人供应商。 |
| | `user` | `PERSONAL` | 普通用户只能看到自己创建的个人供应商。 |
| **模型分发** | `admin` | `ORGANIZATION` | 只有 `scope=ORGANIZATION` 的供应商下的模型才能被管理员看到并分发。 |
| | `user` | (不适用) | 普通用户没有分发权限。 |
| **使用模型** | `admin`, `user` | `ORGANIZATION`, `PERSONAL` | 任何用户都可以使用 (1) 被分发给他们的组织模型 + (2) 他们自己创建的个人模型。 |

## 3. 实现方案

权限检查的核心将位于 **API 服务层**，而不是仅仅依赖于路由层的 `Depends`。这允许我们实现更复杂的、与数据相关的权限逻辑。

### 3.1. 角色检查依赖 (无变化)

我们仍然使用 `RoleChecker` 来保护需要特定角色的端点。

```python
# app/core/permissions.py

class RoleChecker:
    def __init__(self, allowed_roles: list[str]):
        self.allowed_roles = allowed_roles

    def __call__(self, user: User = Depends(current_active_user)):
        if user.role not in self.allowed_roles:
            raise HTTPException(status.HTTP_403_FORBIDDEN, "Role not allowed")
        return user

require_admin = RoleChecker(["admin", "super_admin"])
```

### 3.2. 服务层的 Scope 验证 (核心变更)

真正的权限逻辑将在服务函数中实现，因为它需要访问数据库。

**示例: 创建供应商的服务逻辑**

```python
# app/services/provider_service.py

from app.models import User, Provider
from app.schemas.provider import ProviderCreate
from sqlalchemy.ext.asyncio import AsyncSession
from fastapi import HTTPException, status

class ProviderService:
    async def create_provider(self, db: AsyncSession, provider_data: ProviderCreate, current_user: User) -> Provider:
        
        # 核心权限检查：Scope Validation
        if provider_data.scope == 'ORGANIZATION' and current_user.role not in ['admin', 'super_admin']:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Only Admins can create ORGANIZATION providers."
            )

        # 创建逻辑
        new_provider = Provider(
            **provider_data.dict(),
            owner_id=current_user.id  # 明确所有者
        )
        db.add(new_provider)
        await db.commit()
        await db.refresh(new_provider)
        return new_provider

    async def delete_provider(self, db: AsyncSession, provider_id: uuid.UUID, current_user: User):
        provider = await db.get(Provider, provider_id)
        if not provider:
            raise HTTPException(status_code=404, detail="Provider not found")

        # 对象级权限检查：只有所有者或超管能删除
        if provider.owner_id != current_user.id and current_user.role != 'super_admin':
            raise HTTPException(status_code=403, detail="Not authorized to delete this provider")

        await db.delete(provider)
        await db.commit()
```

### 3.3. API 路由层

API 路由层现在更简洁，主要负责接收数据和调用服务。

```python
# app/api/v1/providers.py

from fastapi import APIRouter, Depends
from app.models.user import User
from app.api.deps import current_active_user
from app.services.provider_service import ProviderService # 假设服务被封装
from app.schemas.provider import ProviderCreate, ProviderRead

router = APIRouter()
provider_service = ProviderService() # 实例化服务

@router.post("/", response_model=ProviderRead)
async def create_provider(
    provider_data: ProviderCreate,
    current_user: User = Depends(current_active_user) # 任何登录用户都可以尝试调用
):
    # 将权限检查和业务逻辑委托给服务层
    return await provider_service.create_provider(
        db=db_session, # 假设db session已注入
        provider_data=provider_data,
        current_user=current_user
    )
```

这种设计将复杂的权限规则内聚到服务层，使得 API 路��更清晰，也更容易对核心业务规则进行单元测试。
