// 审核表单组件
import React from 'react';
import { Form, Radio, Input, Button, Space } from 'antd';
import type { FormInstance } from 'antd';

const { TextArea } = Input;

interface ReviewFormProps {
  onSubmit?: (values: { result: '通过' | '驳回'; comment: string }) => Promise<void>;
  onCancel?: () => void;
  loading?: boolean;
  form?: FormInstance;
  reviewType: 'business' | 'political';
}

const ReviewForm: React.FC<ReviewFormProps> = ({
  onSubmit,
  onCancel,
  loading = false,
  form,
  reviewType
}) => {
  const [formInstance] = Form.useForm();
  const currentForm = form || formInstance;

  const handleSubmit = async () => {
    try {
      const values = await currentForm.validateFields();
      await onSubmit?.(values);
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  return (
    <Form
      form={currentForm}
      layout="vertical"
      initialValues={{
        result: '通过'
      }}
    >
      <Form.Item
        label="审核结果"
        name="result"
        rules={[{ required: true, message: '请选择审核结果' }]}
      >
        <Radio.Group>
          <Radio value="通过">通过</Radio>
          <Radio value="驳回">驳回</Radio>
        </Radio.Group>
      </Form.Item>

      <Form.Item
        label="审核意见"
        name="comment"
        rules={[
          {
            required: true,
            message: '请输入审核意见'
          },
          {
            min: 10,
            message: '审核意见至少10个字符'
          },
          {
            max: 500,
            message: '审核意见不能超过500个字符'
          }
        ]}
      >
        <TextArea
          placeholder={`请输入${reviewType === 'business' ? '业务' : '政工'}审核意见,至少10个字符`}
          rows={6}
          showCount
          maxLength={500}
        />
      </Form.Item>

      {onSubmit && onCancel && (
        <Form.Item>
          <Space>
            <Button type="primary" onClick={handleSubmit} loading={loading}>
              提交审核
            </Button>
            <Button onClick={onCancel}>
              取消
            </Button>
          </Space>
        </Form.Item>
      )}
    </Form>
  );
};

export default ReviewForm;
