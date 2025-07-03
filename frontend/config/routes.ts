/**
 * LYSS AI Platform 路由配置
 * 基于后端API设计的完整路由体系
 * 包含用户认证、对话界面、管理后台等功能模块
 */
export default [
  // 用户认证路由 (无布局)
  {
    path: '/user',
    layout: false,
    routes: [
      {
        name: 'login',
        path: '/user/login',
        component: './User/Login',
      },
      {
        name: 'register', 
        path: '/user/register',
        component: './User/Register',
      },
    ],
  },
  
  // 主应用路由
  {
    path: '/chat',
    name: 'chat',
    icon: 'message',
    component: './Chat',
  },
  
  // 文件管理
  {
    path: '/files',
    name: 'files',
    icon: 'folder',
    component: './Files',
  },
  
  // 成本监控
  {
    path: '/analytics',
    name: 'analytics', 
    icon: 'lineChart',
    component: './Analytics',
  },
  
  // 管理员功能
  {
    path: '/admin',
    name: 'admin',
    icon: 'crown',
    access: 'canAdmin',
    routes: [
      {
        path: '/admin',
        redirect: '/admin/providers',
      },
      {
        path: '/admin/providers',
        name: 'providers',
        component: './Admin/Providers',
      },
      {
        path: '/admin/models', 
        name: 'models',
        component: './Admin/Models',
      },
      {
        path: '/admin/users',
        name: 'users',
        component: './Admin/Users',
      },
    ],
  },
  
  // 个人设置
  {
    path: '/settings',
    name: 'settings',
    icon: 'setting',
    component: './Settings',
  },
  
  // 默认跳转到对话页面
  {
    path: '/',
    redirect: '/chat',
  },
  
  // 404页面
  {
    path: '*',
    layout: false,
    component: './404',
  },
];
