# 指令：依据最终架构，执行 LYSS AI 平台开发文档 V2 重构

你好，首席解决方案架构师。我们已经完成了关键架构的讨论并达成最终共识。现在，请严格遵循以下最终确定的设计方案，对整套开发文档进行 **V2 版本的重构**，确保所有文档都精确、无歧义地反映此最终设计。

## 1. 最终核心架构：供应商“作用域” (Provider Scope)

这是本次重构的 **唯一真理和核心基石**。所有相关设计必须严格遵守此模型。系统中不再有“系统级/个人级”的口头区分，而是由供应商自身的 `scope` 属性决定其一切行为。

| 属性/行为 | `scope = ORGANIZATION` | `scope = PERSONAL` |
| :--- | :--- | :--- |
| **中文含义** | 组织供应商 | 个人供应商 |
| **创建权限** | 仅限 **管理员 (Admin)** 角色 | **任何已登录用户** |
| **资产归属** | 属于 **组织**，不因创建者离职而受影响 | 严格属于 **创建者个人** |
| **可见性** | 在“模型分发”后台对管理员可见 | 仅对创建者本人可见 |
| **核心能力** | **可被分发** 给组织内其他用户 | **不可被分发**，仅限自用 |

---

## 2. 技术栈升级指令

* **前端**:
    * **必须** 废弃 `@ant-design/pro-chat`，全面转向使用其官方推荐的继任者 `@ant-design/x`，要注意`@ant-design/x`和`@ant-design/pro-chat`，它只是一个AI 组件库。
    * **必须** 将 `antd` 及其相关依赖 (`@ant-design/...`) 升级至当前（2025年7月）最新的 **稳定版本**。请使用 Google Search 工具精确查证并指定版本号。
* **后端**:
    * **必须** 重新调研并更新所有后端核心依赖（如 `FastAPI`, `SQLAlchemy`, `Pydantic`, `Qdrant-client`, `mem0ai` 等）至当前最新的、经过社区验证的 **稳定版本**。
* **Mem0.ai**:
    * **必须** 在文档中明确，`mem0.ai` 服务自身的 LLM/Embedding Provider 是 **全局统一的**，由超级管理员通过 **环境变量** 进行配置，与上述的用户供应商体系完全分离。

---

## 3. 文档重构行动项 (Action Items)

请基于 `第1部分` 的最终架构和 `第2部分` 的技术栈指令，**彻底重写或更新** 以下文档的相关章节：

* **`00_Introduction_and_Architecture.md`**:
    * 更新“核心架构”部分，详细描述供应商的 `scope` 模型。
    * 更新“技术选型”章节，列出所有升级后的依赖及其确切版本号。

* **`02_Backend_Development/02_Database_Schema.md`**:
    * 在 `providers` 表的设计中，**必须** 增加一个 `scope` 字段 (类型为枚举: `ORGANIZATION`, `PERSONAL`)。
    * 明确 `owner_id` 字段为外键，关联到 `users` 表，用于指明创建者。
    * 使用 Mermaid.js 更新 ER 图以反映最新的表结构。

* **`01_Backend_Development/03_RBAC_Implementation.md`**:
    * 更新角色权限矩阵，明确：
        * `Admin`: 具备创建 `scope=ORGANIZATION` 和 `scope=PERSONAL` 供应商的权限。
        * `User`: 仅具备创建 `scope=PERSONAL` 供应商的权限。

* **`01_Backend_Development/01_API_Contract.md` (使用 OpenAPI 3.0)**:
    * **创建供应商 API (`POST /api/v1/providers`)**: 请求体中必须包含 `scope` 字段。API内部需要有权限检查逻辑：普通用户提交 `scope=ORGANIZATION` 的请求应被拒绝 (403 Forbidden)。
    * **模型分发源列表 API (`GET /api/v1/admin/distributable-models`)**: 此 API 的数据源 **必须只查询** `scope = 'ORGANIZATION'` 的供应商及其模型。
    * **个人可用模型列表 API (`GET /api/v1/chat/available-models`)**: 其业务逻辑描述必须为：返回一个聚合列表，包含 (1) 所有已授权给当前用户的 `ORGANIZATION` 模型 + (2) 当前用户自己创建的所有 `PERSONAL` 模型。

* **`01_Backend_Development/05_Mem0_Integration.md`**:
    * 按 `2.3` 指令，重写配置部分，提供清晰的 `.env.example` 示例。

* **全局审查**: 请审查所有其他文档（如图表、流程描述），确保其中不再有任何与旧模型相关的模糊描述，全部更新为 `scope` 模型。

---

## 4. 交付物 (Deliverables)

1.  完成上述所有文档的重构。
2.  创建你的记忆文件，确保能时刻追踪项目状态。
3.  在状态中明确注明：“**架构已最终确定，所有文档已按 v2-final 指令重构完毕，等待项目负责人最终确认后即可启动开发。**”

现在，请开始执行最终版的重构任务。