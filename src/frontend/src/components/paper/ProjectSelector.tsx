// 课题选择器组件
import React, { useState } from 'react';
import { Select, Button, Input, Modal, Form, message, Space, Tag } from 'antd';
import { PlusOutlined, SearchOutlined } from '@ant-design/icons';
import { projectService } from '../../services/projectService';
import type { Project } from '../../types/paper';

const { Option } = Select;

interface ProjectSelectorProps {
  mode: 'create' | 'edit' | 'view';
  value?: number[];
  onChange?: (projectIds: number[]) => void;
}

const ProjectSelector: React.FC<ProjectSelectorProps> = ({ mode, value = [], onChange }) => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [form] = Form.useForm();

  // 加载课题列表
  const loadProjects = async (keyword: string = '') => {
    setLoading(true);
    try {
      const response = await projectService.getProjects({ page: 1, size: 50, name: keyword });
      setProjects(response.data.data.list || []);
    } catch (error) {
      message.error('加载课题列表失败');
    } finally {
      setLoading(false);
    }
  };

  React.useEffect(() => {
    loadProjects();
  }, []);

  // 搜索课题
  const handleSearch = (value: string) => {
    setSearchKeyword(value);
    loadProjects(value);
  };

  // 选择课题
  const handleChange = (selectedValues: number[]) => {
    onChange?.(selectedValues);
  };

  // 创建新课题
  const handleCreateProject = async () => {
    try {
      const values = await form.validateFields();
      await projectService.createProject(values);
      message.success('创建课题成功');
      setCreateModalVisible(false);
      form.resetFields();
      // 重新加载课题列表
      loadProjects();
    } catch (error) {
      message.error('创建课题失败');
    }
  };

  // 获取选中的课题信息
  const getSelectedProjects = (): Project[] => {
    return projects.filter(p => value.includes(p.id));
  };

  return (
    <div>
      {mode !== 'view' && (
        <Select
          mode="multiple"
          placeholder="请选择课题"
          value={value}
          onChange={handleChange}
          onSearch={handleSearch}
          filterOption={false}
          loading={loading}
          showSearch
          allowClear
          style={{ width: '100%' }}
          dropdownRender={menu => (
            <>
              {menu}
              <div style={{ padding: 8, borderTop: '1px solid #f0f0f0' }}>
                <Button
                  type="link"
                  icon={<PlusOutlined />}
                  onClick={() => setCreateModalVisible(true)}
                  style={{ width: '100%' }}
                >
                  手动录入新课题
                </Button>
              </div>
            </>
          )}
        >
          {projects.map(project => (
            <Option key={project.id} value={project.id}>
              {project.name} ({project.code})
            </Option>
          ))}
        </Select>
      )}

      {mode === 'view' && value.length > 0 && (
        <Space direction="vertical" style={{ width: '100%' }}>
          {getSelectedProjects().map(project => (
            <Tag key={project.id}>
              {project.name} ({project.code})
            </Tag>
          ))}
        </Space>
      )}

      {/* 创建新课题的模态框 */}
      <Modal
        title="创建新课题"
        open={createModalVisible}
        onOk={handleCreateProject}
        onCancel={() => {
          setCreateModalVisible(false);
          form.resetFields();
        }}
        okText="创建"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label="课题名称"
            name="name"
            rules={[{ required: true, message: '请输入课题名称' }]}
          >
            <Input placeholder="请输入课题名称" maxLength={200} />
          </Form.Item>

          <Form.Item
            label="课题编号"
            name="code"
            rules={[{ required: true, message: '请输入课题编号' }]}
          >
            <Input placeholder="请输入课题编号" maxLength={100} />
          </Form.Item>

          <Form.Item
            label="项目类型"
            name="project_type"
            rules={[{ required: true, message: '请选择项目类型' }]}
          >
            <Select placeholder="请选择项目类型">
              <Option value="纵向">纵向</Option>
              <Option value="横向">横向</Option>
            </Select>
          </Form.Item>

          <Form.Item label="项目来源" name="source">
            <Input placeholder="请输入项目来源" maxLength={200} />
          </Form.Item>

          <Form.Item
            label="项目级别"
            name="level"
            rules={[{ required: true, message: '请选择项目级别' }]}
          >
            <Select placeholder="请选择项目级别">
              <Option value="国家级">国家级</Option>
              <Option value="省部级">省部级</Option>
              <Option value="市级">市级</Option>
              <Option value="其他">其他</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ProjectSelector;
