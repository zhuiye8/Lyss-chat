/**
 * LYSS AI Platform 权限控制配置
 * 基于用户角色的权限管理系统
 */

import { User } from '@/types/api';

export default function access(initialState: { currentUser?: User } | undefined) {
  const { currentUser } = initialState ?? {};
  
  return {
    // 管理员权限：admin 或 super_admin
    canAdmin: currentUser && ['admin', 'super_admin'].includes(currentUser.role),
    
    // 超级管理员权限：仅 super_admin
    canSuperAdmin: currentUser && currentUser.role === 'super_admin',
    
    // 普通用户权限：已登录即可
    canUser: !!currentUser,
    
    // 特定角色检查函数
    hasRole: (role: string | string[]) => {
      if (!currentUser) return false;
      if (Array.isArray(role)) {
        return role.includes(currentUser.role);
      }
      return currentUser.role === role;
    },
  };
}
