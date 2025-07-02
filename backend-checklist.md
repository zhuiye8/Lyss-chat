# 后端开发清单 (V2.1)

本清单概述了根据 V2.1 文档实现后端所需的步骤。

## 第一阶段：项目脚手架与核心模型

- [ ] 初始化 Poetry 项目: `poetry init` & `poetry install`。
- [ ] 根据 `04_Provider_Plugin_System.md` 创建完整的目录结构。
- [ ] 实现 `app/core/config.py` 用于环境变量管理。
- [ ] 实现 `app/core/security.py` 用于加密/解密。
- [ ] 实现 `app/core/database.py` 用于 SQLAlchemy 设置。
- [ ] 根据 `02_Database_Schema.md` 在 `app/models/` 中创建所有 SQLAlchemy 模型，确保关联关系正确定义。
- [ ] 在 `app/schemas/` 中创建所有与模型对应的 Pydantic schema。

## 第二阶段：用户认证与 RBAC

- [ ] 实现 `app/services/user_manager.py` 以支持 `fastapi-users`。
- [ ] 实现 `app/api/deps.py` 以设置认证后端和依赖项。
- [ ] 实现 `app/core/permissions.py` 以提供 `RoleChecker` 依赖。
- [ ] 在 `app/main.py` 中挂载所有 `fastapi-users` 的路由。
- [ ] 创建 `initial_data.py` 脚本用于创建超级管理员。

## 第三阶段：供应商插件系统

- [ ] 在 `app/providers/base.py` 中实现 `LLMProvider` 抽象基类。
- [ ] �� `app/providers/registry.py` 中实现供应商注册表。
- [ ] 在 `app/providers/factory.py` 中实现供应商工厂。
- [ ] 至少实现一个具体的供应商（例如 `openai_provider.py`）来测试整个系统，包括其特定的 Pydantic 配置模型。
- [ ] 在 `app/main.py` 的 `startup` 事件中调用 `discover_providers()`。

## 第四阶段：核心 API 端点与服务

- [ ] 创建 `app/services/provider_service.py`，包含 `create_provider` 和 `delete_provider` 逻辑，并集成 scope 验证。
- [ ] 创建 `/api/v1/providers` 路由和端点，将逻辑委托给 `ProviderService`。
- [ ] 创建 `app/services/chat_service.py`。
- [ ] 根据 `06_Core_Business_Logic.md` 在 `ChatService` 中实现 `get_available_models_for_user` 方法。
- [ ] 在 `ChatService` 中实现聊天完成的权限检查逻辑。
- [ ] 创建 `/api/v1/chat` 路由及其相应的端点。
- [ ] 创建仅限管理员访问的端点 (`/api/v1/admin/...`) 用于模型分发。

## 第五阶段：Mem0 集成与最终化

- [ ] 根据文档实现 `app/services/memory_service.py`。
- [ ] 使用后台任务将 `MemoryService` 的调用集成到 `ChatService` 中。
- [ ] 实现 `UsageService` 用于成本计算，并将其集成到 `ChatService`。
- [ ] 为所有关键服务编写单元和集成测试。
- [ ] 完成用于生产部署�� `docker-compose.yml` 和 `Dockerfile`。
- [ ] 在 `.env.example` 中详细记录所有环境变量。