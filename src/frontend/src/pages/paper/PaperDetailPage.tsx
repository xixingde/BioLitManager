// 论文详情页面
import React, { useEffect, useState } from 'react';
import { Card, Breadcrumb, Descriptions, Tag, Button, Space, message, Tabs, Table, Modal } from 'antd';
import { HomeOutlined, FileTextOutlined, EditOutlined, ArrowLeftOutlined, SendOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { paperService } from '../../services/paperService';
import { reviewService } from '../../services/reviewService';
import type { Paper, ReviewLog } from '../../types/paper';
import type { ReviewLog as ReviewLogType } from '../../types/review';

const { TabPane } = Tabs;

const PaperDetailPage: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState(false);
  const [paper, setPaper] = useState<Paper | null>(null);
  const [reviewLogs, setReviewLogs] = useState<ReviewLogType[]>([]);
  const [submitModalVisible, setSubmitModalVisible] = useState(false);

  // 加载论文详情
  const loadPaperDetail = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const response = await paperService.getPaper(Number(id));
      setPaper(response.data.data);
      // 加载审核记录
      await loadReviewLogs(Number(id));
    } catch (error) {
      message.error('加载论文详情失败');
    } finally {
      setLoading(false);
    }
  };

  // 加载审核记录
  const loadReviewLogs = async (paperId: number) => {
    try {
      const response = await reviewService.getReviewLogs(paperId);
      setReviewLogs(response.data.data);
    } catch (error) {
      console.error('加载审核记录失败:', error);
    }
  };

  useEffect(() => {
    loadPaperDetail();
  }, [id]);

  // 编辑论文
  const handleEdit = () => {
    navigate(`/papers/${id}/edit`);
  };

  // 提交审核
  const handleSubmitReview = async () => {
    if (!id) return;
    try {
      await paperService.submitForReview(Number(id));
      message.success('提交审核成功');
      setSubmitModalVisible(false);
      loadPaperDetail();
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

  // 审核记录表格列
  const reviewColumns = [
    {
      title: '审核类型',
      dataIndex: 'review_type',
      key: 'review_type',
      render: (type: string) => (
        <Tag color={type === 'business' ? 'blue' : 'green'}>
          {type === 'business' ? '业务审核' : '政工审核'}
        </Tag>
      ),
      width: 120
    },
    {
      title: '审核结果',
      dataIndex: 'result',
      key: 'result',
      render: (result: string) => (
        <Tag color={result === 'approved' ? 'success' : 'error'}>
          {result === 'approved' ? '通过' : '驳回'}
        </Tag>
      ),
      width: 120
    },
    {
      title: '审核意见',
      dataIndex: 'comment',
      key: 'comment',
      ellipsis: true
    },
    {
      title: '审核人',
      dataIndex: ['reviewer', 'name'],
      key: 'reviewer',
      width: 120
    },
    {
      title: '审核时间',
      dataIndex: 'review_time',
      key: 'review_time',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN')
    }
  ];

  if (!paper) {
    return <div>加载中...</div>;
  }

  return (
    <div>
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item>
          <HomeOutlined />
        </Breadcrumb.Item>
        <Breadcrumb.Item>
          <FileTextOutlined />
          论文管理
        </Breadcrumb.Item>
        <Breadcrumb.Item>论文详情</Breadcrumb.Item>
      </Breadcrumb>

      <Card
        title={
          <Space>
            <Button
              type="text"
              icon={<ArrowLeftOutlined />}
              onClick={() => navigate('/papers')}
            >
              返回
            </Button>
            <span>论文详情</span>
          </Space>
        }
        extra={
          <Space>
            {paper.status === 'draft' && (
              <>
                <Button type="primary" icon={<EditOutlined />} onClick={handleEdit}>
                  编辑
                </Button>
                <Button type="primary" icon={<SendOutlined />} onClick={() => setSubmitModalVisible(true)}>
                  提交审核
                </Button>
              </>
            )}
          </Space>
        }
        loading={loading}
      >
        <Tabs defaultActiveKey="1">
          <TabPane tab="基本信息" key="1">
            <Descriptions bordered column={2}>
              <Descriptions.Item label="论文标题" span={2}>
                {paper.title}
              </Descriptions.Item>
              <Descriptions.Item label="摘要" span={2}>
                {paper.abstract}
              </Descriptions.Item>
              <Descriptions.Item label="期刊">
                {paper.journal?.short_name}
              </Descriptions.Item>
              <Descriptions.Item label="影响因子">
                {paper.impact_factor}
              </Descriptions.Item>
              <Descriptions.Item label="DOI">
                {paper.doi}
              </Descriptions.Item>
              <Descriptions.Item label="出版日期">
                {paper.publish_date ? new Date(paper.publish_date).toLocaleDateString('zh-CN') : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={statusColors[paper.status]}>{paper.status}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="提交人">
                {paper.submitter?.name}
              </Descriptions.Item>
              <Descriptions.Item label="提交时间">
                {new Date(paper.submit_time).toLocaleString('zh-CN')}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {new Date(paper.created_at).toLocaleString('zh-CN')}
              </Descriptions.Item>

              {/* 作者信息 */}
              {paper.authors && paper.authors.length > 0 && (
                <Descriptions.Item label="作者列表" span={2}>
                  <Space direction="vertical" style={{ width: '100%' }}>
                    {paper.authors.map((author, index) => (
                      <div key={author.id || index}>
                        <Tag color="blue">{author.author_type}</Tag>
                        {index + 1}. {author.name}
                        {author.department && ` - ${author.department}`}
                      </div>
                    ))}
                  </Space>
                </Descriptions.Item>
              )}

              {/* 课题信息 */}
              {paper.projects && paper.projects.length > 0 && (
                <Descriptions.Item label="关联课题" span={2}>
                  <Space wrap>
                    {paper.projects.map(project => (
                      <Tag key={project.id} color="green">
                        {project.name} ({project.code})
                      </Tag>
                    ))}
                  </Space>
                </Descriptions.Item>
              )}

              {/* 附件信息 */}
              {paper.attachments && paper.attachments.length > 0 && (
                <Descriptions.Item label="附件列表" span={2}>
                  <Space direction="vertical" style={{ width: '100%' }}>
                    {paper.attachments.map(attachment => (
                      <Tag key={attachment.id} color="orange">
                        {attachment.file_type}: {attachment.file_name} ({Math.round(attachment.file_size / 1024)} KB)
                      </Tag>
                    ))}
                  </Space>
                </Descriptions.Item>
              )}
            </Descriptions>
          </TabPane>

          <TabPane tab="审核记录" key="2">
            <Table
              columns={reviewColumns}
              dataSource={reviewLogs}
              rowKey="id"
              pagination={false}
            />
          </TabPane>
        </Tabs>
      </Card>

      {/* 提交审核确认模态框 */}
      <Modal
        title="提交审核确认"
        open={submitModalVisible}
        onOk={handleSubmitReview}
        onCancel={() => setSubmitModalVisible(false)}
        okText="确认提交"
        cancelText="取消"
      >
        <p>确认提交该论文进行审核吗?</p>
        <p>提交后论文将进入"待业务审核"状态,需要业务审核人员进行审核。</p>
      </Modal>
    </div>
  );
};

export default PaperDetailPage;
