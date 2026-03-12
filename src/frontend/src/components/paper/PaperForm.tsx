// 论文表单组件
import React from 'react';
import { Form, Input, Select, DatePicker, InputNumber, message } from 'antd';
import type { FormInstance } from 'antd';
import dayjs from 'dayjs';
import { paperService } from '../../services/paperService';
import { journalService } from '../../services/journalService';
import AuthorList from './AuthorList';
import ProjectSelector from './ProjectSelector';
import FileUpload from './FileUpload';
import type { Paper, PaperForm, Author } from '../../types/paper';

const { TextArea } = Input;
const { Option } = Select;

interface PaperFormProps {
  mode: 'create' | 'edit' | 'view';
  initialValues?: Partial<Paper>;
  onSubmit?: (values: PaperForm) => Promise<void>;
  onSaveDraft?: (values: PaperForm) => Promise<void>;
  paperId?: number;
  form?: FormInstance;
}

const PaperForm: React.FC<PaperFormProps> = ({
  mode,
  initialValues,
  onSubmit,
  onSaveDraft,
  paperId,
  form
}) => {
  const [formInstance] = Form.useForm();
  const currentForm = form || formInstance;
  const [journals, setJournals] = React.useState<any[]>([]);
  const [loadingJournals, setLoadingJournals] = React.useState(false);
  const [duplicateCheck, setDuplicateCheck] = React.useState<{ count: number; papers: Paper[] } | null>(null);

  // 加载期刊列表
  React.useEffect(() => {
    loadJournals();
  }, []);

  // 初始化表单值
  React.useEffect(() => {
    if (initialValues) {
      currentForm.setFieldsValue({
        ...initialValues,
        publish_date: initialValues.publish_date ? dayjs(initialValues.publish_date) : undefined,
      });
    }
  }, [initialValues, currentForm]);

  const loadJournals = async () => {
    setLoadingJournals(true);
    try {
      const response = await journalService.listJournals({ page: 1, size: 1000 });
      setJournals(response.data.data.list || []);
    } catch (error) {
      message.error('加载期刊列表失败');
    } finally {
      setLoadingJournals(false);
    }
  };

  // 检查论文重复
  const checkDuplicate = async () => {
    const title = currentForm.getFieldValue('title');
    const doi = currentForm.getFieldValue('doi');

    if (!title && !doi) {
      return;
    }

    try {
      const response = await paperService.checkDuplicate({ title, doi });
      setDuplicateCheck(response.data);
      if (response.data.count > 0) {
        message.warning(`检测到 ${response.data.count} 篇可能重复的论文`);
      } else {
        message.success('未检测到重复论文');
      }
    } catch (error) {
      message.error('重复检查失败');
    }
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await currentForm.validateFields();
      const formData: PaperForm = {
        ...values,
        publish_date: values.publish_date ? values.publish_date.format('YYYY-MM-DD') : undefined,
      };
      await onSubmit?.(formData);
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  // 保存草稿
  const handleSaveDraft = async () => {
    try {
      const values = currentForm.getFieldsValue();
      const formData: PaperForm = {
        ...values,
        publish_date: values.publish_date ? values.publish_date.format('YYYY-MM-DD') : undefined,
      };
      await onSaveDraft?.(formData);
    } catch (error) {
      console.error('保存草稿失败:', error);
    }
  };

  return (
    <Form
      form={currentForm}
      layout="vertical"
      disabled={mode === 'view'}
    >
      <Form.Item
        label="论文标题"
        name="title"
        rules={[{ required: true, message: '请输入论文标题' }]}
      >
        <Input placeholder="请输入论文标题" maxLength={500} showCount onBlur={checkDuplicate} />
      </Form.Item>

      <Form.Item label="摘要" name="abstract">
        <TextArea placeholder="请输入论文摘要" rows={6} maxLength={2000} showCount />
      </Form.Item>

      <Form.Item
        label="期刊"
        name="journal_id"
        rules={[{ required: true, message: '请选择期刊' }]}
      >
        <Select
          placeholder="请选择期刊"
          loading={loadingJournals}
          showSearch
          filterOption={(input, option) =>
            (option?.children as string)?.toLowerCase().includes(input.toLowerCase())
          }
        >
          {journals.map(journal => (
            <Option key={journal.id} value={journal.id}>
              {journal.short_name} ({journal.full_name})
            </Option>
          ))}
        </Select>
      </Form.Item>

      <Form.Item label="DOI" name="doi" rules={[{ required: true, message: '请输入DOI' }]}>
        <Input placeholder="请输入DOI" onBlur={checkDuplicate} />
      </Form.Item>

      <Form.Item label="影响因子" name="impact_factor" rules={[{ required: true, message: '请输入影响因子' }]}>
        <InputNumber placeholder="请输入影响因子" min={0} step={0.001} style={{ width: '100%' }} />
      </Form.Item>

      <Form.Item label="出版日期" name="publish_date">
        <DatePicker placeholder="请选择出版日期" style={{ width: '100%' }} />
      </Form.Item>

      <Form.Item label="作者信息" name="authors">
        <AuthorList mode={mode} />
      </Form.Item>

      <Form.Item label="关联课题" name="projects">
        <ProjectSelector mode={mode} />
      </Form.Item>

      {paperId && (
        <Form.Item label="附件上传">
          <FileUpload paperId={paperId} mode={mode} />
        </Form.Item>
      )}

      {mode !== 'view' && (
        <div style={{ marginTop: 24 }}>
          <Form.Item>
            <div style={{ display: 'flex', gap: 12 }}>
              {onSubmit && (
                <button type="button" onClick={handleSubmit} style={{ padding: '6px 15px' }}>
                  提交审核
                </button>
              )}
              {onSaveDraft && (
                <button type="button" onClick={handleSaveDraft} style={{ padding: '6px 15px' }}>
                  保存草稿
                </button>
              )}
            </div>
          </Form.Item>
        </div>
      )}
    </Form>
  );
};

export default PaperForm;
