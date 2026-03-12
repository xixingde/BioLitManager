// 论文录入页面
import React from 'react';
import { Card, Breadcrumb, message } from 'antd';
import { HomeOutlined, FileTextOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import PaperForm from '../../components/paper/PaperForm';
import { paperService } from '../../services/paperService';
import type { PaperForm as PaperFormData } from '../../types/paper';

const PaperCreatePage: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = React.useState(false);
  const [paperId, setPaperId] = React.useState<number | undefined>();

  // 提交审核
  const handleSubmit = async (values: PaperFormData) => {
    setLoading(true);
    try {
      // 如果已经创建了论文(保存了草稿),则更新并提交
      if (paperId) {
        await paperService.updatePaper(paperId, values);
        await paperService.submitForReview(paperId);
      } else {
        // 创建新论文并直接提交
        const response = await paperService.createPaper(values);
        await paperService.submitForReview(response.data.data.id);
      }
      message.success('论文提交审核成功');
      navigate('/papers');
    } catch (error) {
      message.error('提交审核失败');
    } finally {
      setLoading(false);
    }
  };

  // 保存草稿
  const handleSaveDraft = async (values: PaperFormData) => {
    setLoading(true);
    try {
      if (paperId) {
        // 更新现有草稿
        await paperService.saveDraft(paperId, values);
        message.success('草稿保存成功');
      } else {
        // 创建新草稿
        const response = await paperService.createPaper(values);
        setPaperId(response.data.data.id);
        message.success('草稿保存成功');
      }
    } catch (error) {
      message.error('保存草稿失败');
    } finally {
      setLoading(false);
    }
  };

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
        <Breadcrumb.Item>录入论文</Breadcrumb.Item>
      </Breadcrumb>

      <Card title="录入论文" loading={loading}>
        <PaperForm
          mode="create"
          onSubmit={handleSubmit}
          onSaveDraft={handleSaveDraft}
          paperId={paperId}
        />
      </Card>
    </div>
  );
};

export default PaperCreatePage;
