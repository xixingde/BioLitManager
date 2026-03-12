// 归档管理页面
import { useState, useEffect } from 'react';
import { 
  Table, Button, Select, Space, Tag, message, Card, Row, Col, 
  Modal, Descriptions, Divider, Form, Input, Drawer 
} from 'antd';
import { EyeOutlined, EyeInvisibleOutlined, EditOutlined, FilterOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { archiveService, ArchiveListParams, ModifyRequest } from '../../services/archiveService';
import type { Paper, ArchiveInfo } from '../../types/paper';
import { useNavigate } from 'react-router-dom';

const { Option } = Select;
const { TextArea } = Input;

// 归档论文类型（扩展Paper类型，包含归档信息）
interface ArchivePaper extends Paper {
  archive_info?: ArchiveInfo;
}

// 课题选项（模拟数据，实际应从项目服务获取）
const projectOptions = [
  { value: 'project1', label: '国家重点研发计划' },
  { value: 'project2', label: '国家自然科学基金' },
  { value: 'project3', label: '省部级项目' },
  { value: 'project4', label: '横向课题' },
];

// 收录类型选项
const paperTypeOptions = [
  { value: 'sci', label: 'SCI' },
  { value: 'ei', label: 'EI' },
  { value: 'cssci', label: 'CSSCI' },
  { value: 'core', label: '核心期刊' },
];

// 生成年份选项
const currentYear = new Date().getFullYear();
const yearOptions = Array.from({ length: 10 }, (_, i) => ({
  value: currentYear - i,
  label: `${currentYear - i}年`
}));

const ArchivePage: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [archives, setArchives] = useState<ArchivePaper[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [size] = useState(20);
  
  // 筛选条件
  const [year, setYear] = useState<number | undefined>();
  const [paperType, setPaperType] = useState<string | undefined>();
  const [projectCode, setProjectCode] = useState<string | undefined>();
  const [author, setAuthor] = useState<string | undefined>();
  
  // 详情抽屉
  const [detailVisible, setDetailVisible] = useState(false);
  const [selectedPaper, setSelectedPaper] = useState<ArchivePaper | null>(null);
  
  // 修改申请弹窗
  const [modifyModalVisible, setModifyModalVisible] = useState(false);
  const [modifyForm] = Form.useForm();

  // 获取归档列表
  const fetchArchives = async () => {
    setLoading(true);
    try {
      const params: ArchiveListParams = {
        year,
        paperType,
        projectCode,
        author,
        page,
        pageSize: size
      };
      const response = await archiveService.getArchiveList(params);
      setArchives(response.data.list as unknown as ArchivePaper[]);
      setTotal(response.data.total);
    } catch (error) {
      message.error('获取归档列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchArchives();
  }, [page, year, paperType, projectCode, author]);

  // 查看详情
  const handleViewDetail = async (record: ArchivePaper) => {
    try {
      const response = await archiveService.getArchiveByPaperId(record.id);
      setSelectedPaper({ ...record, ...response.data } as ArchivePaper);
      setDetailVisible(true);
    } catch (error) {
      message.error('获取归档详情失败');
    }
  };

  // 隐藏归档论文
  const handleHide = async (paperId: number) => {
    try {
      await archiveService.hideArchive(paperId);
      message.success('隐藏成功');
      fetchArchives();
    } catch (error) {
      message.error('隐藏失败');
    }
  };

  // 提交修改申请
  const handleModifySubmit = async () => {
    try {
      const values = await modifyForm.validateFields();
      if (selectedPaper) {
        const requestData: ModifyRequest = {
          requestType: values.requestType,
          requestReason: values.requestReason,
          requestData: values.requestData ? JSON.parse(values.requestData) : undefined
        };
        await archiveService.submitModifyRequest(selectedPaper.id, requestData);
        message.success('修改申请提交成功');
        setModifyModalVisible(false);
        modifyForm.resetFields();
      }
    } catch (error) {
      message.error('提交修改申请失败');
    }
  };

  // 打开修改申请弹窗
  const handleModifyRequest = (record: ArchivePaper) => {
    setSelectedPaper(record);
    setModifyModalVisible(true);
  };

  // 状态标签颜色
  const statusColors: Record<string, string> = {
    'draft': 'default',
    '待业务审核': 'blue',
    '待政工审核': 'processing',
    '审核通过': 'success',
    '驳回': 'error'
  };

  // 表格列定义
  const columns: ColumnsType<ArchivePaper> = [
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
      title: '作者列表',
      dataIndex: 'authors',
      key: 'authors',
      width: 200,
      render: (authors?: ArchivePaper['authors']) => {
        if (!authors || authors.length === 0) return '-';
        return authors.map(a => a.name).join(', ');
      }
    },
    {
      title: '期刊',
      dataIndex: ['journal', 'short_name'],
      key: 'journal',
      width: 120
    },
    {
      title: '归档时间',
      dataIndex: ['archive_info', 'archive_time'],
      key: 'archive_time',
      width: 180,
      render: (time: string) => time ? new Date(time).toLocaleString('zh-CN') : '-'
    },
    {
      title: '归档人',
      dataIndex: ['archive_info', 'archive_user', 'name'],
      key: 'archive_user',
      width: 100
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={statusColors[status]}>{status}</Tag>
      )
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          <Button 
            type="link" 
            icon={<EyeOutlined />} 
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          <Button 
            type="link" 
            icon={<EditOutlined />} 
            onClick={() => handleModifyRequest(record)}
          >
            修改申请
          </Button>
          <Button 
            type="link" 
            danger 
            icon={<EyeInvisibleOutlined />} 
            onClick={() => handleHide(record.id)}
          >
            隐藏
          </Button>
        </Space>
      )
    }
  ];

  return (
    <div>
      {/* 筛选区域 */}
      <Card style={{ marginBottom: 16 }}>
        <Row gutter={16}>
          <Col span={6}>
            <Space direction="vertical" style={{ width: '100%' }}>
              <span>年份</span>
              <Select
                placeholder="选择年份"
                allowClear
                style={{ width: '100%' }}
                value={year}
                onChange={setYear}
              >
                {yearOptions.map(opt => (
                  <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                ))}
              </Select>
            </Space>
          </Col>
          <Col span={6}>
            <Space direction="vertical" style={{ width: '100%' }}>
              <span>收录类型</span>
              <Select
                placeholder="选择收录类型"
                allowClear
                style={{ width: '100%' }}
                value={paperType}
                onChange={setPaperType}
              >
                {paperTypeOptions.map(opt => (
                  <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                ))}
              </Select>
            </Space>
          </Col>
          <Col span={6}>
            <Space direction="vertical" style={{ width: '100%' }}>
              <span>课题</span>
              <Select
                placeholder="选择课题"
                allowClear
                style={{ width: '100%' }}
                value={projectCode}
                onChange={setProjectCode}
              >
                {projectOptions.map(opt => (
                  <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                ))}
              </Select>
            </Space>
          </Col>
          <Col span={6}>
            <Space direction="vertical" style={{ width: '100%' }}>
              <span>作者</span>
              <Input
                placeholder="输入作者姓名"
                value={author}
                onChange={(e) => setAuthor(e.target.value || undefined)}
                allowClear
              />
            </Space>
          </Col>
        </Row>
        <Divider style={{ margin: '16px 0' }} />
        <Space>
          <Button 
            icon={<FilterOutlined />} 
            onClick={() => {
              setPage(1);
              fetchArchives();
            }}
          >
            应用筛选
          </Button>
          <Button onClick={() => {
            setYear(undefined);
            setPaperType(undefined);
            setProjectCode(undefined);
            setAuthor(undefined);
            setPage(1);
          }}>
            重置
          </Button>
        </Space>
      </Card>

      {/* 归档列表 */}
      <Card title={`归档论文列表 (共 ${total} 条)`}>
        <Table
          columns={columns}
          dataSource={archives}
          rowKey="id"
          loading={loading}
          pagination={{
            current: page,
            pageSize: size,
            total,
            showSizeChanger: false,
            showQuickJumper: true,
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
            onChange: (pageNum) => setPage(pageNum)
          }}
        />
      </Card>

      {/* 详情抽屉 */}
      <Drawer
        title="归档论文详情"
        placement="right"
        width={600}
        open={detailVisible}
        onClose={() => setDetailVisible(false)}
      >
        {selectedPaper && (
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="论文ID">{selectedPaper.id}</Descriptions.Item>
            <Descriptions.Item label="论文标题">{selectedPaper.title}</Descriptions.Item>
            <Descriptions.Item label="作者列表">
              {selectedPaper.authors?.map(a => a.name).join(', ')}
            </Descriptions.Item>
            <Descriptions.Item label="期刊">
              {selectedPaper.journal?.full_name} ({selectedPaper.journal?.short_name})
            </Descriptions.Item>
            <Descriptions.Item label="DOI">{selectedPaper.doi}</Descriptions.Item>
            <Descriptions.Item label="影响因子">{selectedPaper.impact_factor}</Descriptions.Item>
            <Descriptions.Item label="出版日期">
              {selectedPaper.publish_date ? new Date(selectedPaper.publish_date).toLocaleDateString('zh-CN') : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="审核状态">
              <Tag color={statusColors[selectedPaper.status]}>{selectedPaper.status}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="课题">
              {selectedPaper.projects?.map(p => p.name).join(', ')}
            </Descriptions.Item>
            <Descriptions.Item label="归档时间">
              {selectedPaper.archive_info?.archive_time 
                ? new Date(selectedPaper.archive_info.archive_time).toLocaleString('zh-CN') 
                : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="归档人">
              {selectedPaper.archive_info?.archive_user?.name || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="归档原因">
              {selectedPaper.archive_info?.archive_reason || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="摘要">{selectedPaper.abstract || '-'}</Descriptions.Item>
            <Descriptions.Item label="附件">
              {selectedPaper.attachments?.map(a => (
                <div key={a.id}>
                  <a>{a.file_name}</a> ({(a.file_size / 1024).toFixed(2)} KB)
                </div>
              )) || '无'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>

      {/* 修改申请弹窗 */}
      <Modal
        title="提交修改申请"
        open={modifyModalVisible}
        onOk={handleModifySubmit}
        onCancel={() => {
          setModifyModalVisible(false);
          modifyForm.resetFields();
        }}
        okText="提交"
        cancelText="取消"
      >
        <Form form={modifyForm} layout="vertical">
          <Form.Item
            name="requestType"
            label="申请类型"
            rules={[{ required: true, message: '请选择申请类型' }]}
          >
            <Select placeholder="选择申请类型">
              <Option value="update_info">信息更新</Option>
              <Option value="add_attachment">添加附件</Option>
              <Option value="remove_attachment">删除附件</Option>
              <Option value="other">其他</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="requestReason"
            label="申请原因"
            rules={[{ required: true, message: '请输入申请原因' }]}
          >
            <TextArea rows={4} placeholder="请详细描述申请原因" />
          </Form.Item>
          <Form.Item
            name="requestData"
            label="修改内容 (JSON格式，可选)"
          >
            <TextArea rows={3} placeholder='{"field": "value"}' />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ArchivePage;
