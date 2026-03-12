// 首页
import { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Spin, List, Tag, Avatar, Button, Alert, Space, Typography } from 'antd';
import { PlusOutlined, UploadOutlined, FileTextOutlined, CheckCircleOutlined, AuditOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';
import { userStore } from '../../stores/userStore';
import { getHomeData } from '../../services/statisticsService';
import type { Statistics, RecentPaper, PendingReviewTask } from '../../types/statistics';

const { Title, Text } = Typography;

// 状态标签颜色映射
const getStatusColor = (status: string): string => {
  const colorMap: Record<string, string> = {
    draft: 'default',
    '待业务审核': 'orange',
    '待政工审核': 'gold',
    '审核通过': 'green',
    '驳回': 'red',
  };
  return colorMap[status] || 'default';
};

// 状态标签文本映射
const getStatusText = (status: string): string => {
  const textMap: Record<string, string> = {
    draft: '草稿',
    '待业务审核': '待业务审核',
    '待政工审核': '待政工审核',
    '审核通过': '审核通过',
    '驳回': '驳回',
  };
  return textMap[status] || status;
};

const HomePage: React.FC = () => {
  const navigate = useNavigate();
  const { user } = userStore();
  
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [statistics, setStatistics] = useState<Statistics | null>(null);
  const [recentPapers, setRecentPapers] = useState<RecentPaper[]>([]);
  const [pendingReviews, setPendingReviews] = useState<PendingReviewTask[]>([]);

  // 获取首页数据
  const fetchHomeData = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getHomeData(user?.id);
      setStatistics(data.statistics);
      setRecentPapers(data.recentPapers);
      setPendingReviews(data.pendingReviews);
    } catch (err) {
      setError('获取数据失败，请稍后重试');
      console.error('获取首页数据失败:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHomeData();
  }, [user?.id]);

  // 判断用户角色是否显示审核入口
  const showBusinessReview = user?.role === 'business_reviewer' || user?.role === 'admin' || user?.role === 'super_admin';
  const showPoliticalReview = user?.role === 'political_reviewer' || user?.role === 'admin' || user?.role === 'super_admin';

  // 快捷操作列表
  const quickActions = [
    {
      key: 'new',
      title: '录入论文',
      icon: <PlusOutlined />,
      path: '/papers/create',
      description: '录入新的论文信息',
    },
    {
      key: 'import',
      title: '批量导入',
      icon: <UploadOutlined />,
      path: '/papers/batch-import',
      description: '批量导入论文数据',
    },
    {
      key: 'my',
      title: '我的论文',
      icon: <FileTextOutlined />,
      path: '/papers',
      description: '查看我提交的论文',
    },
    ...(showBusinessReview ? [{
      key: 'business-review',
      title: '业务审核',
      icon: <AuditOutlined />,
      path: '/reviews/business',
      description: '审核业务相关论文',
    }] : []),
    ...(showPoliticalReview ? [{
      key: 'political-review',
      title: '政工审核',
      icon: <CheckCircleOutlined />,
      path: '/reviews/political',
      description: '审核政工相关论文',
    }] : []),
  ];

  // 加载状态
  if (loading) {
    return (
      <div style={{ padding: '24px', minHeight: '100vh', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  // 错误状态
  if (error) {
    return (
      <div style={{ padding: '24px' }}>
        <Alert
          message="数据加载失败"
          description={error}
          type="error"
          showIcon
          action={
            <Button type="primary" onClick={fetchHomeData}>
              重试
            </Button>
          }
        />
      </div>
    );
  }

  return (
    <div style={{ padding: '24px' }}>
      {/* 欢迎区域 */}
      <Card style={{ marginBottom: 24 }}>
        <Row justify="space-between" align="middle">
          <Col>
            <Space align="center">
              <Avatar size={48} style={{ backgroundColor: '#1890ff' }}>
                {user?.name?.charAt(0) || 'U'}
              </Avatar>
              <div>
                <Title level={4} style={{ margin: 0 }}>
                  欢迎回来，{user?.name || user?.username || '用户'}
                </Title>
                <Text type="secondary">{dayjs().format('YYYY年MM月DD日 dddd')}</Text>
              </div>
            </Space>
          </Col>
        </Row>
      </Card>

      {/* 统计卡片区域 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="论文总数"
              value={statistics?.total || 0}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="待审核"
              value={(statistics?.pendingBusiness || 0) + (statistics?.pendingPolitical || 0)}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="已通过"
              value={statistics?.approved || 0}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="我的草稿"
              value={statistics?.myDraft || 0}
              valueStyle={{ color: '#8c8c8c' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 快捷操作和最近动态区域 */}
      <Row gutter={[16, 16]}>
        {/* 快捷操作入口 */}
        <Col xs={24} lg={12}>
          <Card title="快捷操作" bordered={false} style={{ height: '100%' }}>
            <List
              grid={{ gutter: 16, column: 2 }}
              dataSource={quickActions}
              renderItem={(item) => (
                <List.Item>
                  <Card
                    hoverable
                    onClick={() => navigate(item.path)}
                    style={{ textAlign: 'center' }}
                  >
                    <div style={{ fontSize: 32, marginBottom: 8, color: '#1890ff' }}>
                      {item.icon}
                    </div>
                    <div style={{ fontWeight: 500 }}>{item.title}</div>
                    <div style={{ fontSize: 12, color: '#8c8c8c' }}>{item.description}</div>
                  </Card>
                </List.Item>
              )}
            />
          </Card>
        </Col>

        {/* 最近动态 */}
        <Col xs={24} lg={12}>
          <Card 
            title="最近动态" 
            bordered={false} 
            extra={<a onClick={() => navigate('/papers')}>查看全部</a>}
            style={{ height: '100%' }}
          >
            {recentPapers.length > 0 ? (
              <List
                itemLayout="horizontal"
                dataSource={recentPapers}
                renderItem={(paper) => (
                  <List.Item 
                    style={{ cursor: 'pointer' }}
                    onClick={() => navigate(`/papers/${paper.id}`)}
                  >
                    <List.Item.Meta
                      avatar={
                        <Avatar style={{ backgroundColor: '#f0f0f0' }}>
                          <FileTextOutlined />
                        </Avatar>
                      }
                      title={
                        <Space>
                          <Text ellipsis style={{ maxWidth: 200, margin: 0 }}>{paper.title}</Text>
                          <Tag color={getStatusColor(paper.status)}>
                            {getStatusText(paper.status)}
                          </Tag>
                        </Space>
                      }
                      description={
                        <Text type="secondary">
                          {paper.submitter_name} · {dayjs(paper.created_at).format('MM-DD HH:mm')}
                        </Text>
                      }
                    />
                  </List.Item>
                )}
              />
            ) : (
              <div style={{ textAlign: 'center', padding: '40px 0', color: '#8c8c8c' }}>
                暂无最近动态
              </div>
            )}
          </Card>
        </Col>
      </Row>

      {/* 待审核任务（仅审核角色可见） */}
      {(showBusinessReview || showPoliticalReview) && pendingReviews.length > 0 && (
        <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
          <Col span={24}>
            <Card 
              title="待审核任务" 
              bordered={false}
              extra={
                <a onClick={() => navigate(showBusinessReview ? '/reviews/business' : '/reviews/political')}>
                  查看全部
                </a>
              }
            >
              <List
                itemLayout="horizontal"
                dataSource={pendingReviews.slice(0, 5)}
                renderItem={(task) => (
                  <List.Item 
                    style={{ cursor: 'pointer' }}
                    onClick={() => navigate(`/papers/${task.id}`)}
                  >
                    <List.Item.Meta
                      avatar={
                        <Avatar style={{ backgroundColor: task.review_type === 'business' ? '#fa8c16' : '#faad14' }}>
                          <AuditOutlined />
                        </Avatar>
                      }
                      title={
                        <Space>
                          <Text ellipsis style={{ maxWidth: 300, margin: 0 }}>{task.title}</Text>
                          <Tag color={task.review_type === 'business' ? 'orange' : 'gold'}>
                            {task.review_type === 'business' ? '业务审核' : '政工审核'}
                          </Tag>
                        </Space>
                      }
                      description={
                        <Text type="secondary">
                          {task.submitter_name} · 提交于 {dayjs(task.submit_time).format('MM-DD HH:mm')}
                        </Text>
                      }
                    />
                  </List.Item>
                )}
              />
            </Card>
          </Col>
        </Row>
      )}
    </div>
  );
};

export default HomePage;
