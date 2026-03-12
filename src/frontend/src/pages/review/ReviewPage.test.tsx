import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import { message } from 'antd';
import ReviewPage from '../ReviewPage';
import * as reviewService from '../../services/reviewService';

// Mock services
jest.mock('../../services/reviewService');
jest.mock('antd', () => ({
  ...jest.requireActual('antd'),
  message: {
    success: jest.fn(),
    error: jest.fn(),
    warning: jest.fn(),
  },
}));

const mockPendingPapers = [
  {
    id: 1,
    title: '人工智能在生物学中的应用研究',
    status: '待业务审核',
    submitter_name: '张三',
    submit_time: '2024-01-15T10:00:00Z',
    days_since_submit: 1,
  },
  {
    id: 2,
    title: '机器学习在基因组学中的应用',
    status: '待业务审核',
    submitter_name: '李四',
    submit_time: '2024-01-14T10:00:00Z',
    days_since_submit: 2,
  },
];

const mockReviewLogs = [
  {
    id: 1,
    paper_id: 1,
    review_type: '业务审核',
    result: '通过',
    comment: '论文内容完整，格式符合要求',
    reviewer: {
      id: 2,
      username: 'reviewer1',
      name: '审核员1',
      role: '业务审核员',
    },
    review_time: '2024-01-16T10:00:00Z',
  },
];

const renderWithRouter = (component: React.ReactNode) => {
  return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('ReviewPage Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (reviewService.getPendingBusinessReviews as jest.Mock).mockResolvedValue({
      data: mockPendingPapers,
    });
    (reviewService.getReviewLogs as jest.Mock).mockResolvedValue({
      data: mockReviewLogs,
    });
  });

  // 测试1: 业务审核列表展示
  it('should display business review list correctly', async () => {
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
      expect(screen.getByText('机器学习在基因组学中的应用')).toBeInTheDocument();
      expect(screen.getByText('张三')).toBeInTheDocument();
      expect(screen.getByText('李四')).toBeInTheDocument();
    });
  });

  // 测试2: 政工审核列表展示
  it('should display political review list correctly', async () => {
    const politicalPapers = [
      {
        id: 3,
        title: '深度学习在蛋白质结构预测中的应用',
        status: '待政工审核',
        submitter_name: '王五',
        submit_time: '2024-01-13T10:00:00Z',
        days_since_submit: 3,
      },
    ];

    (reviewService.getPendingPoliticalReviews as jest.Mock).mockResolvedValue({
      data: politicalPapers,
    });

    renderWithRouter(<ReviewPage reviewType="political" />);

    await waitFor(() => {
      expect(screen.getByText('深度学习在蛋白质结构预测中的应用')).toBeInTheDocument();
      expect(screen.getByText('王五')).toBeInTheDocument();
    });
  });

  // 测试3: 审核操作 - 通过
  it('should approve paper successfully', async () => {
    (reviewService.businessReview as jest.Mock).mockResolvedValue({});
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 点击审核按钮
    const reviewButton = screen.getAllByRole('button', { name: /审核/i })[0];
    fireEvent.click(reviewButton);

    // 在审核对话框中点击通过
    await waitFor(() => {
      const approveButton = screen.getByRole('button', { name: /通过/i });
      fireEvent.click(approveButton);

      // 确认审核
      const confirmButton = screen.getByRole('button', { name: /确定/i });
      fireEvent.click(confirmButton);
    });

    // 验证API调用
    await waitFor(() => {
      expect(reviewService.businessReview).toHaveBeenCalledWith(
        1,
        expect.objectContaining({
          result: '通过',
        })
      );
      expect(message.success).toHaveBeenCalledWith('审核成功');
    });
  });

  // 测试4: 审核操作 - 驳回
  it('should reject paper successfully', async () => {
    const user = userEvent.setup();
    (reviewService.businessReview as jest.Mock).mockResolvedValue({});
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 点击审核按钮
    const reviewButton = screen.getAllByRole('button', { name: /审核/i })[0];
    fireEvent.click(reviewButton);

    // 在审核对话框中选择驳回并输入意见
    await waitFor(async () => {
      const rejectRadio = screen.getByRole('radio', { name: /驳回/i });
      await user.click(rejectRadio);

      // 输入审核意见
      const commentInput = screen.getByPlaceholderText(/请输入审核意见/i);
      await user.type(commentInput, '论文格式不符合要求，请重新整理');

      // 确认审核
      const confirmButton = screen.getByRole('button', { name: /确定/i });
      await user.click(confirmButton);
    });

    // 验证API调用
    await waitFor(() => {
      expect(reviewService.businessReview).toHaveBeenCalledWith(
        1,
        expect.objectContaining({
          result: '驳回',
          comment: '论文格式不符合要求，请重新整理',
        })
      );
      expect(message.success).toHaveBeenCalledWith('审核成功');
    });
  });

  // 测试5: 审核操作 - 取消
  it('should cancel review operation', async () => {
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 点击审核按钮
    const reviewButton = screen.getAllByRole('button', { name: /审核/i })[0];
    fireEvent.click(reviewButton);

    // 取消审核
    const cancelButton = await screen.findByRole('button', { name: /取消/i });
    fireEvent.click(cancelButton);

    // 验证对话框关闭
    await waitFor(() => {
      expect(screen.queryByPlaceholderText(/请输入审核意见/i)).not.toBeInTheDocument();
    });
  });

  // 测试6: 查看审核记录
  it('should display review logs correctly', async () => {
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 点击查看审核记录按钮
    const viewLogsButton = screen.getAllByRole('button', { name: /查看记录/i })[0];
    fireEvent.click(viewLogsButton);

    // 验证审核记录展示
    await waitFor(() => {
      expect(screen.getByText('业务审核')).toBeInTheDocument();
      expect(screen.getByText('通过')).toBeInTheDocument();
      expect(screen.getByText('论文内容完整，格式符合要求')).toBeInTheDocument();
      expect(screen.getByText('审核员1')).toBeInTheDocument();
    });
  });

  // 测试7: 空审核列表
  it('should display empty state when no pending reviews', async () => {
    (reviewService.getPendingBusinessReviews as jest.Mock).mockResolvedValue({
      data: [],
    });

    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.queryByText('人工智能在生物学中的应用研究')).not.toBeInTheDocument();
      expect(screen.getByText(/暂无待审核论文/i) || screen.getByText(/没有找到待审核的论文/i)).toBeInTheDocument();
    });
  });

  // 测试8: 审核超时提醒
  it('should display timeout warning for pending reviews', async () => {
    const timeoutPapers = [
      {
        id: 1,
        title: '超时待审核论文',
        status: '待业务审核',
        submitter_name: '张三',
        submit_time: '2024-01-10T10:00:00Z',
        days_since_submit: 6, // 超过5天
      },
    ];

    (reviewService.getPendingBusinessReviews as jest.Mock).mockResolvedValue({
      data: timeoutPapers,
    });

    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('超时待审核论文')).toBeInTheDocument();
      // 验证超时警告标签
      expect(screen.getByText(/超时/i) || screen.getByText(/6天/i)).toBeInTheDocument();
    });
  });

  // 测试9: 权限控制 - 非审核员访问
  it('should handle unauthorized access', async () => {
    (reviewService.getPendingBusinessReviews as jest.Mock).mockRejectedValue(
      new Error('Unauthorized')
    );

    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(message.error).toHaveBeenCalledWith('权限不足');
    });
  });

  // 测试10: 批量审核功能
  it('should support batch review operation', async () => {
    (reviewService.businessReview as jest.Mock).mockResolvedValue({});
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 选中多个待审核论文
    const checkboxes = screen.getAllByRole('checkbox');
    checkboxes.forEach((checkbox) => {
      if (checkbox !== checkboxes[0]) { // 跳过全选框
        fireEvent.click(checkbox);
      }
    });

    // 点击批量审核按钮
    const batchReviewButton = screen.getByRole('button', { name: /批量审核/i });
    fireEvent.click(batchReviewButton);

    // 批量通过
    await waitFor(() => {
      const approveButton = screen.getByRole('button', { name: /通过/i });
      fireEvent.click(approveButton);

      const confirmButton = screen.getByRole('button', { name: /确定/i });
      fireEvent.click(confirmButton);
    });

    // 验证批量审核
    await waitFor(() => {
      expect(reviewService.businessReview).toHaveBeenCalledTimes(2);
      expect(message.success).toHaveBeenCalledWith('批量审核成功');
    });
  });

  // 测试11: 搜索待审核论文
  it('should search pending reviews', async () => {
    const user = userEvent.setup();
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 搜索功能实现（模拟）
    const searchInput = screen.getByPlaceholderText(/搜索论文/i);
    await user.type(searchInput, '机器学习');

    // 验证搜索结果
    await waitFor(() => {
      expect(screen.getByText('机器学习在基因组学中的应用')).toBeInTheDocument();
      expect(screen.queryByText('人工智能在生物学中的应用研究')).not.toBeInTheDocument();
    });
  });

  // 测试12: 审核历史时间线
  it('should display review timeline', async () => {
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 点击查看审核记录
    const viewLogsButton = screen.getAllByRole('button', { name: /查看记录/i })[0];
    fireEvent.click(viewLogsButton);

    // 验证时间线展示
    await waitFor(() => {
      expect(screen.getByText('2024-01-16')).toBeInTheDocument(); // 审核日期
      expect(screen.getByText('业务审核')).toBeInTheDocument();
    });
  });

  // 测试13: 驳回后的重新提交
  it('should handle rejected paper resubmission', async () => {
    const rejectedPapers = [
      {
        id: 1,
        title: '被驳回的论文',
        status: '驳回',
        submitter_name: '张三',
        submit_time: '2024-01-15T10:00:00Z',
        days_since_submit: 1,
      },
    ];

    (reviewService.getPendingBusinessReviews as jest.Mock).mockResolvedValue({
      data: rejectedPapers,
    });

    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('被驳回的论文')).toBeInTheDocument();
      // 验证驳回状态的标签
      expect(screen.getByText('驳回')).toBeInTheDocument();
    });
  });

  // 测试14: 审核意见必填校验
  it('should validate comment when rejecting', async () => {
    (reviewService.businessReview as jest.Mock).mockResolvedValue({});
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 点击审核按钮
    const reviewButton = screen.getAllByRole('button', { name: /审核/i })[0];
    fireEvent.click(reviewButton);

    // 选择驳回但不输入意见
    await waitFor(() => {
      const rejectRadio = screen.getByRole('radio', { name: /驳回/i });
      fireEvent.click(rejectRadio);

      // 尝试直接提交
      const confirmButton = screen.getByRole('button', { name: /确定/i });
      fireEvent.click(confirmButton);
    });

    // 验证错误提示
    await waitFor(() => {
      expect(screen.getByText(/请输入审核意见/i) || screen.getByText(/审核意见不能为空/i)).toBeInTheDocument();
      expect(reviewService.businessReview).not.toHaveBeenCalled();
    });
  });

  // 测试15: 刷新审核列表
  it('should refresh review list', async () => {
    renderWithRouter(<ReviewPage reviewType="business" />);

    await waitFor(() => {
      expect(screen.getByText('人工智能在生物学中的应用研究')).toBeInTheDocument();
    });

    // 点击刷新按钮
    const refreshButton = screen.getByRole('button', { name: /刷新/i });
    fireEvent.click(refreshButton);

    // 验证重新加载数据
    await waitFor(() => {
      expect(reviewService.getPendingBusinessReviews).toHaveBeenCalledTimes(2);
    });
  });
});
