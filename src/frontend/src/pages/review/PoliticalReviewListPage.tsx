// 政工审核列表页
import { useState, useEffect } from 'react';
import { Table, Button, Tag, message } from 'antd';
import { EyeOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { reviewService } from '../../services/reviewService';
import type { PendingReview } from '../../types/review';
import { useNavigate } from 'react-router-dom';

const PoliticalReviewListPage: React.FC = () => {
  const navigate = useNavigate();
  const [papers, setPapers] = useState<PendingReview[]>([]);
  const [loading, setLoading] = useState(false);

  // 获取待审核列表
  const fetchPendingReviews = async () => {
    setLoading(true);
    try {
      const response = await reviewService.getPendingPoliticalReviews();
      setPapers(response.data);
    } catch (error) {
      message.error('获取待审核列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPendingReviews();
    // 定时刷新
    const interval = setInterval(fetchPendingReviews, 5 * 60 * 1000);
    return () => clearInterval(interval);
  }, []);

  const columns: ColumnsType<PendingReview> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80
    },
    {
      title: '论文标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true
    },
    {
      title: '提交人',
      dataIndex: 'submitter_name',
      key: 'submitter_name',
      width: 120
    },
    {
      title: '提交时间',
      dataIndex: 'submit_time',
      key: 'submit_time',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN')
    },
    {
      title: '距提交天数',
      dataIndex: 'days_since_submit',
      key: 'days_since_submit',
      width: 120,
      render: (days: number) => (
        <Tag color={days >= 2 ? 'error' : 'success'}>{days} 天</Tag>
      )
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 120
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_, record) => (
        <Button type="link" icon={<EyeOutlined />} onClick={() => navigate(`/reviews/political/${record.id}`)}>
          去审核
        </Button>
      )
    }
  ];

  return (
    <div>
      <h2>待政工审核列表</h2>
      <Table
        columns={columns}
        dataSource={papers}
        rowKey="id"
        loading={loading}
        pagination={false}
      />
    </div>
  );
};

export default PoliticalReviewListPage;
