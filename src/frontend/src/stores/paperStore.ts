// 论文状态管理
import { create } from 'zustand';
import type { Paper, PaperListResponse, PaperForm } from '../types/paper';
import { paperService } from '../services/paperService';

interface PaperState {
  // 论文列表
  papers: Paper[];
  total: number;
  page: number;
  size: number;
  loading: boolean;

  // 当前论文
  currentPaper: Paper | null;

  // 操作方法
  fetchPapers: (params?: { page?: number; size?: number; status?: string; keyword?: string }) => Promise<void>;
  fetchPaper: (id: number) => Promise<void>;
  createPaper: (data: PaperForm) => Promise<number>;
  updatePaper: (id: number, data: Partial<PaperForm>) => Promise<void>;
  deletePaper: (id: number) => Promise<void>;
  submitForReview: (id: number) => Promise<void>;
  saveDraft: (id: number, data: Partial<PaperForm>) => Promise<void>;
  checkDuplicate: (title: string, doi: string) => Promise<{ count: number; papers: Paper[] }>;
  batchImport: (file: File) => Promise<{ success: number; failed: number; errors: string[] }>;
  setCurrentPaper: (paper: Paper | null) => void;
  reset: () => void;
}

const initialState = {
  papers: [],
  total: 0,
  page: 1,
  size: 10,
  loading: false,
  currentPaper: null,
};

export const usePaperStore = create<PaperState>((set, get) => ({
  ...initialState,

  // 获取论文列表
  fetchPapers: async (params = {}) => {
    set({ loading: true });
    try {
      const { page = 1, size = 10, status, keyword } = params;
      const response = await paperService.getPapers({ page, size, status, keyword });
      set({
        papers: response.data.data.list || [],
        total: response.data.data.total || 0,
        page,
        size,
        loading: false,
      });
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  // 获取论文详情
  fetchPaper: async (id: number) => {
    set({ loading: true });
    try {
      const response = await paperService.getPaper(id);
      set({
        currentPaper: response.data.data,
        loading: false,
      });
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  // 创建论文
  createPaper: async (data: PaperForm) => {
    set({ loading: true });
    try {
      const response = await paperService.createPaper(data);
      set({ loading: false });
      return response.data.data.id;
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  // 更新论文
  updatePaper: async (id: number, data: Partial<PaperForm>) => {
    set({ loading: true });
    try {
      await paperService.updatePaper(id, data);
      set({ loading: false });
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  // 删除论文
  deletePaper: async (id: number) => {
    set({ loading: true });
    try {
      await paperService.deletePaper(id);
      // 从列表中移除
      const { papers } = get();
      set({
        papers: papers.filter(p => p.id !== id),
        loading: false,
      });
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  // 提交审核
  submitForReview: async (id: number) => {
    set({ loading: true });
    try {
      await paperService.submitForReview(id);
      // 更新列表中的状态
      const { papers } = get();
      set({
        papers: papers.map(p =>
          p.id === id ? { ...p, status: '待业务审核' } : p
        ),
        loading: false,
      });
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  // 保存草稿
  saveDraft: async (id: number, data: Partial<PaperForm>) => {
    set({ loading: true });
    try {
      await paperService.saveDraft(id, data);
      // 更新当前论文
      if (get().currentPaper?.id === id) {
        await get().fetchPaper(id);
      }
      set({ loading: false });
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  // 检查重复
  checkDuplicate: async (title: string, doi: string) => {
    try {
      const response = await paperService.checkDuplicate({ title, doi });
      return response.data.data;
    } catch (error) {
      throw error;
    }
  },

  // 批量导入
  batchImport: async (file: File) => {
    set({ loading: true });
    try {
      const response = await paperService.batchImport(file);
      set({ loading: false });
      return response.data.data;
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  // 设置当前论文
  setCurrentPaper: (paper: Paper | null) => {
    set({ currentPaper: paper });
  },

  // 重置状态
  reset: () => {
    set(initialState);
  },
}));
