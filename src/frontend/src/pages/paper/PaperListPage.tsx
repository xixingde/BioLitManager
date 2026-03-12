// 论文列表页
import { useState, useEffect } from 'react';
import { Table, Button, Input, Select, Space, Tag, message } from 'antd';
import { PlusOutlined, SearchOutlined, EyeOutlined, EditOutlined, DeleteOutlined, SendOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { paperService } from '../../services/paperService';
import type { Paper } from '../../types/paper';
import { useNavigate } from 'react-router-dom';

const { Search } = Input;
const { Option } = Select;

const PaperListPage: React.FC = () => {
  const navigate = useNavigate();
  const [papers, setPapers] = useState<Paper[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [size, setSize] = useState(10);
  const [loading, setLoading] = useState(false);
  const [keyword, setKeyword] = useState('');
  const [status, setStatus] = useState<string | undefined>();

  // 获取论文列表
  const fetchPapers = async () => {
    setLoading(true);
    try {
      const response = await paperService.getPapers({ page, size, keyword, status });
      setPapers(response.data.list);
      setTotal(response.data.total);
    } catch (error) {
      message.error('获取论文列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPapers();
  }, [page, size, keyword, status]);

  // 删除论文
  const handleDelete = async (id: number) => {
    try {
      await paperService.deletePaper(id);
      message.success('删除成功');
      fetchPapers();
    } catch (error) {
      message.error('删除失败');
    }
  };

  // 提交审核
  const handleSubmit = async (id: number) => {
    try {
      await paperService.submitForReview(id);
      message.success('提交审核成功');
      fetchPapers();
    } catch (error) {
      message.error('提交审核失败');
    }
  };

  // 状态标签颜色映射
  const statusColors: Record<string, string> = {
    'draft': 'default',
    '待业务审核': 'blue',
    '待政工审核': 'processing',
    '审核通过': 'success',
    '驳回': 'error'
  };

  const columns: ColumnsType<Paper> = [
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
      title: '期刊',
      dataIndex: ['journal', 'short_name'],
      key: 'journal',
      width: 150
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 120,
      render: (status: string) => (
        <Tag color={statusColors[status]}>{status}</Tag>
      )
    },
    {
      title: '提交时间',
      dataIndex: 'submit_time',
      key: 'submit_time',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN')
    },
    {
      title: '操作',
      key: 'action',
      width: 280,
      render: (_, record) => (
        <Space size="small">
          <Button type="link" icon={<EyeOutlined />} onClick={() => navigate(`/papers/${record.id}`)}>
            查看
          </Button>
          {record.status === 'draft' && (
            <>
              <Button type="link" icon={<EditOutlined />} onClick={() => navigate(`/papers/${record.id}/edit`)}>
                编辑
              </Button>
              <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
                删除
              </Button>
              <Button type="link" icon={<SendOutlined />} onClick={() => handleSubmit(record.id)}>
                提交审核
              </Button>
            </>
          )}
        </Space>
      )
    }
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Space>
          <Search
            placeholder="搜索论文标题"
            allowClear
            style={{ width: 300 }}
            onSearch={setKeyword}
          />
          <Select
            placeholder="选择状态"
            allowClear
            style={{ width: 150 }}
            onChange={setStatus}
          >
            <Option value="draft">草稿</Option>
            <Option value="待业务审核">待业务审核</Option>
            <Option value="待政工审核">待政工审核</Option>
            <Option value="审核通过">审核通过</Option>
            <Option value="驳回">驳回</Option>
          </Select>
        </Space>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/papers/create')}>
          新增论文
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={papers}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: size,
          total,
          onChange: (page, size) => {
            setPage(page);
            setSize(size);
          }
        }}
      />
    </div>
  );
};

export default PaperListPage;
