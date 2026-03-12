// 审核页面
import React, { useEffect, useState } from 'react';
import { Card, Breadcrumb, Descriptions, Tag, Tabs, message, Modal } from 'antd';
import { HomeOutlined, CheckCircleOutlined, ArrowLeftOutlined, FileTextOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { paperService } from '../../services/paperService';
import { reviewService } from '../../services/reviewService';
import ReviewForm from '../../components/paper/ReviewForm';
import FileUpload from '../../components/paper/FileUpload';
import type { Paper } from '../../types/paper';
import type { ReviewLog as ReviewLogType } from '../../types/review';

const { TabPane } = Tabs;

const ReviewPage: React.FC = () => {
  const navigate = useNavigate();
  const { type, paperId } = useParams<{ type: string; paperId: string }>();
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [paper, setPaper] = useState<Paper | null>(null);
  const [reviewLogs, setReviewLogs] = useState<ReviewLogType[]>([]);

  // 审核类型
  const reviewType = type === 'business' ? 'business' : 'political';

  // 加载论文详情
  const loadPaperDetail = async () => {
    if (!paperId) return;
    setLoading(true);
    try {
      const response = await paperService.getPaper(Number(paperId));
      setPaper(response.data.data);
      // 加载审核记录
      await loadReviewLogs(Number(paperId));
    } catch (error) {
      message.error('加载论文详情失败');
    } finally {
      setLoading(false);
    }
  };

  // 加载审核记录
  const loadReviewLogs = async (id: number) => {
    try {
      const response = await reviewService.getReviewLogs(id);
      setReviewLogs(response.data.data);
    } catch (error) {
      console.error('加载审核记录失败:', error);
    }
  };

  useEffect(() => {
    loadPaperDetail();
  }, [paperId]);

  // 提交审核
  const handleSubmitReview = async (values: { result: '通过' | '驳回'; comment: string }) => {
    if (!paperId) return;
    setSubmitting(true);
    try {
      if (reviewType === 'business') {
        await reviewService.businessReview(Number(paperId), values);
      } else {
        await reviewService.politicalReview(Number(paperId), values);
      }
      message.success('审核提交成功');
      // 返回审核列表
      navigate(`/reviews/${type}`);
    } catch (error) {
      message.error('审核提交失败');
    } finally {
      setSubmitting(false);
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
          <CheckCircleOutlined />
          审核管理
        </Breadcrumb.Item>
        <Breadcrumb.Item>{reviewType === 'business' ? '业务审核' : '政工审核'}</Breadcrumb.Item>
        <Breadcrumb.Item>审核论文</Breadcrumb.Item>
      </Breadcrumb>

      <Card
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <Button
              type="text"
              icon={<ArrowLeftOutlined />}
              onClick={() => navigate(`/reviews/${type}`)}
            >
              返回
            </Button>
            <span>
              {reviewType === 'business' ? '业务审核' : '政工审核'} - {paper.title}
            </span>
          </div>
        }
        loading={loading}
      >
        <Tabs defaultActiveKey="1">
          {/* 论文信息 */}
          <TabPane tab="论文信息" key="1">
            <Descriptions bordered column={2}>
              <Descriptions.Item label="论文标题" span={2}>
                {paper.title}
              </Descriptions.Item>
              <Descriptions.Item label="摘要" span={2}>
                {paper.abstract}
              </Descriptions.Item>
              <Descriptions.Item label="期刊">
                {paper.journal?.short_name} ({paper.journal?.full_name})
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
              <Descriptions.Item label="当前状态">
                <Tag color={statusColors[paper.status]}>{paper.status}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="提交人">
                {paper.submitter?.name}
              </Descriptions.Item>
              <Descriptions.Item label="提交时间">
                {new Date(paper.submit_time).toLocaleString('zh-CN')}
              </Descriptions.Item>
              <Descriptions.Item label="距提交天数">
                {Math.floor((Date.now() - new Date(paper.submit_time).getTime()) / (1000 * 60 * 60 * 24))} 天
              </Descriptions.Item>

              {/* 作者信息 */}
              {paper.authors && paper.authors.length > 0 && (
                <Descriptions.Item label="作者列表" span={2}>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    {paper.authors.map((author, index) => (
                      <div key={author.id || index}>
                        <Tag color="blue">{author.author_type}</Tag>
                        {index + 1}. {author.name}
                        {author.department && ` - ${author.department}`}
                      </div>
                    ))}
                  </div>
                </Descriptions.Item>
              )}

              {/* 课题信息 */}
              {paper.projects && paper.projects.length > 0 && (
                <Descriptions.Item label="关联课题" span={2}>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                    {paper.projects.map(project => (
                      <Tag key={project.id} color="green">
                        {project.name} ({project.code})
                      </Tag>
                    ))}
                  </div>
                </Descriptions.Item>
              )}
            </Descriptions>
          </TabPane>

          {/* 附件查看 */}
          <TabPane tab="附件查看" key="2">
            <FileUpload paperId={Number(paperId)} mode="view" />
          </TabPane>

          {/* 历史审核记录 */}
          <TabPane tab="历史审核" key="3">
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              {reviewLogs.length === 0 ? (
                <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
                  暂无审核记录
                </div>
              ) : (
                reviewLogs.map(log => (
                  <Card key={log.id} size="small">
                    <Descriptions size="small" column={2}>
                      <Descriptions.Item label="审核类型">
                        <Tag color={log.review_type === 'business' ? 'blue' : 'green'}>
                          {log.review_type === 'business' ? '业务审核' : '政工审核'}
                        </Tag>
                      </Descriptions.Item>
                      <Descriptions.Item label="审核结果">
                        <Tag color={log.result === 'approved' ? 'success' : 'error'}>
                          {log.result === 'approved' ? '通过' : '驳回'}
                        </Tag>
                      </Descriptions.Item>
                      <Descriptions.Item label="审核人">
                        {log.reviewer?.name}
                      </Descriptions.Item>
                      <Descriptions.Item label="审核时间">
                        {new Date(log.review_time).toLocaleString('zh-CN')}
                      </Descriptions.Item>
                      <Descriptions.Item label="审核意见" span={2}>
                        {log.comment}
                      </Descriptions.Item>
                    </Descriptions>
                  </Card>
                ))
              )}
            </div>
          </TabPane>

          {/* 提交审核 */}
          <TabPane tab="提交审核" key="4">
            <Card title={reviewType === 'business' ? '业务审核' : '政工审核'}>
              <Alert
                message="审核说明"
                description={
                  <div>
                    <p>• 通过审核后,论文将进入下一审核阶段</p>
                    <p>• 驳回后,论文将返回草稿状态,提交人需要修改后重新提交</p>
                    <p>• 请认真填写审核意见,确保意见清晰、明确</p>
                  </div>
                }
                type="info"
                showIcon
                style={{ marginBottom: 24 }}
              />
              <ReviewForm
                reviewType={reviewType}
                onSubmit={handleSubmitReview}
                loading={submitting}
              />
            </Card>
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default ReviewPage;
