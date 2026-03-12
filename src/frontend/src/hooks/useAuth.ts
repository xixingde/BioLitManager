import { useEffect } from 'react';
import { userStore } from '@/stores/userStore';
import type { Permission } from '@/types/user';
import { ROLE } from '@/utils/constants';

export const useAuth = () => {
  const user = userStore((state) => state.user);
  const token = userStore((state) => state.token);
  const permissions = userStore((state) => state.permissions);
  const isAuthenticated = userStore((state) => state.isAuthenticated);
  const login = userStore((state) => state.login);
  const logout = userStore((state) => state.logout);
  const updatePermissions = userStore((state) => state.updatePermissions);
  const initAuth = userStore((state) => state.initAuth);

  useEffect(() => {
    initAuth();
  }, [initAuth]);

  const hasPermission = (permission: Permission): boolean => {
    if (!user) return false;
    if (user.role === ROLE.SUPER_ADMIN) return true;
    return permissions.includes(permission);
  };

  const hasAnyPermission = (requiredPermissions: Permission[]): boolean => {
    if (!user) return false;
    if (user.role === ROLE.SUPER_ADMIN) return true;
    return requiredPermissions.some((permission) => permissions.includes(permission));
  };

  const hasAllPermissions = (requiredPermissions: Permission[]): boolean => {
    if (!user) return false;
    if (user.role === ROLE.SUPER_ADMIN) return true;
    return requiredPermissions.every((permission) => permissions.includes(permission));
  };

  const hasRole = (role: string): boolean => {
    return user?.role === role;
  };

  return {
    user,
    token,
    permissions,
    isAuthenticated,
    login,
    logout,
    updatePermissions,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    hasRole,
  };
};
