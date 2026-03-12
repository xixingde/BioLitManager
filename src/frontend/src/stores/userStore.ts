import { create } from 'zustand';
import type { UserInfo, Permission } from '@/types/user';
import { authService } from '@/services/authService';
import { TOKEN_KEY, USER_INFO_KEY } from '@/utils/constants';

interface UserState {
  user: UserInfo | null;
  token: string | null;
  permissions: Permission[];
  isAuthenticated: boolean;
}

interface UserActions {
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  updatePermissions: (permissions: Permission[]) => void;
  initAuth: () => void;
}

const loadFromStorage = () => {
  const token = localStorage.getItem(TOKEN_KEY);
  const userInfo = localStorage.getItem(USER_INFO_KEY);
  const user: UserInfo | null = userInfo ? JSON.parse(userInfo) : null;
  return { token, user };
};

export const userStore = create<UserState & UserActions>((set) => ({
  user: null,
  token: null,
  permissions: [],
  isAuthenticated: false,

  login: async (username: string, password: string) => {
    const response = await authService.login({ username, password });
    localStorage.setItem(TOKEN_KEY, response.token);
    localStorage.setItem(USER_INFO_KEY, JSON.stringify(response.user));
    set({
      user: response.user,
      token: response.token,
      isAuthenticated: true,
    });
  },

  logout: async () => {
    try {
      await authService.logout();
    } finally {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_INFO_KEY);
      set({
        user: null,
        token: null,
        permissions: [],
        isAuthenticated: false,
      });
    }
  },

  updatePermissions: (permissions: Permission[]) => {
    set({ permissions });
  },

  initAuth: () => {
    const { token, user } = loadFromStorage();
    set({
      user,
      token,
      isAuthenticated: !!token && !!user,
    });
  },
}));
