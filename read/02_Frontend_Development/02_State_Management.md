# LYSS AI 平台 - 前端状态管理方案

**版本**: 1.0
**最后更新**: 2025年7月2日

---

## 1. 状态管理哲学

在现代复杂的单页应用 (SPA) 中，清晰、高效的状态管理是成功的关键。我们拒绝单一的、庞大的全局 Redux-like store，而是采用一种更现代、更精细化的组合策略，将不同类型的状态交由最适合的工具来管理。

**核心原则**:

*   **区分状态类型**: 严格区分“**服务端状态**”和“**客户端状态**”。
*   **服务端状态 (Server State)**:
    *   **定义**: 存储在服务器上的数据，前端通过API进行读写。例如：用户列表、模型信息、对话历史等。它是异步的、远程的，且不由前端直接控制。
    *   **工具**: **React Query (`@tanstack/react-query`)**。
*   **客户端状态 (Client State)**:
    *   **定义**: 仅存在于浏览器中的数据，用于控制UI的交互和表现。例如：侧边栏是否折叠、当前选中的模型ID、主题模式（暗/亮）等。
    *   **工具**: **Zustand**。

## 2. 工具选型与 Rationale

| 工具 | 推荐版本 | 用途 | 设计思考 |
| :--- | :--- | :--- | :--- |
| **React Query** | `5.51.1` | 服务端状态���理 | 它是管理异步数据的行业标准。它极大地简化了数据获取、缓存、同步和更新的复杂逻辑，让我们无需手动编写 `useEffect` 和 `useState` 的组合来处理 loading, error, success 状态。其内置的缓存和后台重新获取机制能显著提升应用性能和数据一致性。 |
| **Zustand** | `4.5.4` | 客户端状态管理 | 它是一个极简、快速、可扩展的状态管理库。API直观，几乎没有学习成本。它不依赖 React Context，避免了不必要的重渲染，性能更优。非常适合管理那些跨组件共享的、与UI交互相关的少量全局状态。 |

## 3. React Query: 管理服务端状态

React Query 将是我们应用中数据获取的基石。

### 3.1. 配置 `QueryClientProvider`

在应用的根组件 (`app.tsx`) 中，我们需要包裹 `QueryClientProvider`。

```tsx
// src/app.tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      refetchOnWindowFocus: false, // 可选：根据需求调整
    },
  },
});

const App = () => (
  <QueryClientProvider client={queryClient}>
    {/* ... Rest of the app, including Ant Design Pro layout ... */}
  </QueryClientProvider>
);
```

### 3.2. 使用 `useQuery` 获取数据

在任何需要获取服务端数据的组件中，使用 `useQuery`。

**示例**: 获取当前用户信息。

```tsx
// src/components/RightContent/AvatarDropdown.tsx
import { useQuery } from '@tanstack/react-query';
import { getCurrentUser } from '@/services/lyss/userAPI';

const AvatarDropdown = () => {
  const { data: currentUser, isLoading } = useQuery({
    queryKey: ['currentUser'], // 唯一的缓存键
    queryFn: getCurrentUser,   // 异步获取函数
  });

  if (isLoading) {
    return <Spin />;
  }

  return (
    // ... JSX to display user avatar and name
    <span>{currentUser?.name}</span>
  );
};
```

### 3.3. 使用 `useMutation` 修改数据

对于创建、更新、删除等操作，使用 `useMutation`。它能优雅地处理加载状态，并在成功后自动使相关查询失效，触发数据重新获取。

**示例**: 更新供应商配置。

```tsx
// src/pages/Admin/Providers/components/ProviderForm.tsx
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateProvider } from '@/services/lyss/providerAPI';

const ProviderForm = ({ provider }) => {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: updateProvider,
    onSuccess: () => {
      // 当更新成功时，使所有与 providers 相关的查询失效
      queryClient.invalidateQueries({ queryKey: ['providers'] });
      message.success('Provider updated successfully!');
    },
    onError: () => {
      message.error('Failed to update provider.');
    }
  });

  const onFinish = (values) => {
    mutation.mutate({ id: provider.id, ...values });
  };

  return (
    <Form onFinish={onFinish} disabled={mutation.isPending}>
      {/* ... Form fields ... */}
      <Button type="primary" htmlType="submit" loading={mutation.isPending}>
        Submit
      </Button>
    </Form>
  );
};
```

## 4. Zustand: 管理客户端状态

我们将为不同的业务领域创建独立的 store，以保持代码的组织性和可维护性。

### 4.1. 创建 Store

**示例**: 创建一个管理聊天相关UI状态的 store。

```tsx
// src/stores/chatStore.ts
import { create } from 'zustand';

interface ChatState {
  currentModel: string | null; // 当前选中的模型ID
  isSidebarCollapsed: boolean;
  setCurrentModel: (modelId: string) => void;
  toggleSidebar: () => void;
}

export const useChatStore = create<ChatState>((set) => ({
  currentModel: null,
  isSidebarCollapsed: false,
  setCurrentModel: (modelId) => set({ currentModel: modelId }),
  toggleSidebar: () => set((state) => ({ isSidebarCollapsed: !state.isSidebarCollapsed })),
}));
```

### 4.2. 在组件中使用 Store

在组件中，可以像使用普通 hook 一样使用 store。Zustand 会自动处理重渲染，只有当 store 中被访问的部分发生变化时，组件才会更新。

```tsx
// src/components/Business/ModelSelector/index.tsx
import { useChatStore } from '@/stores/chatStore';

const ModelSelector = () => {
  // 只订阅 currentModel 和 setCurrentModel 的变化
  const { currentModel, setCurrentModel } = useChatStore();
  
  // ...
};

// src/components/Layout/SideBar.tsx
import { useChatStore } from '@/stores/chatStore';

const SideBar = () => {
  // 只订阅 isSidebarCollapsed 和 toggleSidebar 的变化
  const { isSidebarCollapsed, toggleSidebar } = useChatStore();

  // ...
};
```

通过这种组合策略，我们为 LYSS AI 平台的前端构建了一个健壮、高效且易于维护的状态管理系统。
