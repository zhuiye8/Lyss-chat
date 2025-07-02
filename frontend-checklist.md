# 前端开发清单 (V2.1)

本清单概述了根据 V2.1 文档实现前端所需的步骤。

## 第一阶段：项目设置与依赖

- [ ] 初始化一个新的 Ant Design Pro 项目。
- [ ] 根据 `00_Introduction_and_Architecture.md` 中的版本号更新 `package.json`。
- [ ] 运行 `pnpm install` 安装所有依赖。
- [ ] 在 `src/app.tsx` 中为 React Query 配置 `QueryClientProvider`。
- [ ] 在 `src/requestErrorConfig.ts` 中设置全局请求配置，以处理 API 响应和认证令牌。

## 第二阶段：状态管理与服务

- [ ] 创建所需的 Zustand store 用于管理UI状态 (例如 `src/stores/chatStore.ts` 用于 `currentModelId`)。
- [ ] 在 `src/services/lyss/` 中创建 API 服务层，为所有后端端点定义函数。
- [ ] 在 `src/services/lyss/typings.d.ts` 中创建相应的 TypeScript 类型定义。

## 第三阶段：核心组件实现

- [ ] 实现 `ModelSelector` 组件 (`src/components/Business/ModelSelector/index.tsx`)。
    - [ ] 它必须从 `GET /api/v1/chat/available-models` 获取数据。
    - [ ] 它必须使用 `useChatStore` 来管理所选模型。
- [ ] 实现 `ChatArea` 组件 (`src/components/Business/ChatArea/index.tsx`)。
    - [ ] **关键**: 它必须使用 `@ant-design/x` 的钩子 (`useXChat`) 和原子组件 (`Chat`, `Bubble`, `Composer`) 来构建，而不是使用 `<ProChat>`。
    - [ ] `onSendMessage` 处理程序必须正确调用 `chatCompletion` 服务，并从 Zustand store 传递 `currentModelId`。
- [ ] 实现主聊天页面 (`src/pages/Chat/index.tsx`)，组装 `ModelSelector` 和 `ChatArea`。

## 第四阶段：供应商管理界面

- [ ] 实现 `ProviderForm` 组件 (`src/pages/Admin/Providers/components/ProviderForm.tsx`)。
    - [ ] 它必须根据所选的 `provider_type` 动态渲染表单字段。
    - [ ] 它必须将表单数据收集到一个 `config` 对象中以便提交。
- [ ] 在管理后台部分实现供应商列表和创建页面。

## 第五阶段：最终化与润色

- [ ] 实现面向用户的供应商管理页面 (用于 `PERSONAL` 范围的供应商)。
- [ ] 实现面向管理员的模型分发页面。
- [ ] 确保所有 API 调用都正确使用 React Query 的 `useQuery` 和 `useMutation` 钩子，以实现正确的状态管理。
- [ ] 完善UI，优雅地处理所有加载和错误状态。
- [ ] 对端到端的用户流程进行彻底测试。