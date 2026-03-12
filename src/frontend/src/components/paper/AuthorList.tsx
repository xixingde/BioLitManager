// 作者列表组件
import React, { useState } from 'react';
import { List, Button, Select, Input, Space, Popconfirm, message } from 'antd';
import { PlusOutlined, MinusCircleOutlined, UpOutlined, DownOutlined } from '@ant-design/icons';
import type { Author } from '../../types/paper';
import { authService } from '../../services/authService';

const { Option } = Select;

interface AuthorListProps {
  mode: 'create' | 'edit' | 'view';
  value?: Author[];
  onChange?: (authors: Author[]) => void;
}

// 作者类型选项
const authorTypeOptions = [
  { label: '第一作者', value: 'first_author' },
  { label: '共同第一作者', value: 'co_first_author' },
  { label: '通讯作者', value: 'corresponding_author' },
  { label: '普通作者', value: 'author' },
];

const AuthorList: React.FC<AuthorListProps> = ({ mode, value = [], onChange }) => {
  const [users, setUsers] = useState<any[]>([]);
  const [loadingUsers, setLoadingUsers] = useState(false);

  // 加载用户列表
  const loadUsers = async (keyword: string = '') => {
    setLoadingUsers(true);
    try {
      const response = await authService.searchUsers({ keyword, page: 1, size: 50 });
      setUsers(response.data.list || []);
    } catch (error) {
      message.error('加载用户列表失败');
    } finally {
      setLoadingUsers(false);
    }
  };

  // 添加作者
  const addAuthor = () => {
    const newAuthor: Author = {
      id: Date.now(),
      name: '',
      author_type: 'author',
      rank: value.length + 1,
      department: '',
    };
    onChange?.([...value, newAuthor]);
  };

  // 删除作者
  const removeAuthor = (index: number) => {
    const newAuthors = value.filter((_, i) => i !== index);
    // 重新排序
    newAuthors.forEach((author, i) => {
      author.rank = i + 1;
    });
    onChange?.(newAuthors);
  };

  // 更新作者信息
  const updateAuthor = (index: number, field: keyof Author, fieldValue: any) => {
    const newAuthors = [...value];
    newAuthors[index] = { ...newAuthors[index], [field]: fieldValue };
    onChange?.(newAuthors);
  };

  // 上移作者
  const moveUp = (index: number) => {
    if (index === 0) return;
    const newAuthors = [...value];
    [newAuthors[index - 1], newAuthors[index]] = [newAuthors[index], newAuthors[index - 1]];
    // 重新排序
    newAuthors.forEach((author, i) => {
      author.rank = i + 1;
    });
    onChange?.(newAuthors);
  };

  // 下移作者
  const moveDown = (index: number) => {
    if (index === value.length - 1) return;
    const newAuthors = [...value];
    [newAuthors[index], newAuthors[index + 1]] = [newAuthors[index + 1], newAuthors[index]];
    // 重新排序
    newAuthors.forEach((author, i) => {
      author.rank = i + 1;
    });
    onChange?.(newAuthors);
  };

  // 检查作者类型是否互斥
  const checkAuthorTypeConflict = (currentIndex: number, selectedType: string) => {
    const conflictTypes: Record<string, string[]> = {
      first_author: ['co_first_author', 'corresponding_author'],
      co_first_author: ['first_author'],
      corresponding_author: ['first_author'],
    };

    const conflictingTypes = conflictTypes[selectedType] || [];
    const hasConflict = value.some((author, i) =>
      i !== currentIndex && conflictingTypes.includes(author.author_type)
    );

    if (hasConflict) {
      message.warning('作者类型冲突: 不能同时存在第一作者/共同第一作者/通讯作者');
    }
  };

  return (
    <div>
      <div style={{ marginBottom: 12, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <span>作者列表 (按排名排序)</span>
        {mode !== 'view' && (
          <Button type="primary" size="small" icon={<PlusOutlined />} onClick={addAuthor}>
            添加作者
          </Button>
        )}
      </div>

      <List
        dataSource={value}
        renderItem={(author, index) => (
          <List.Item
            key={author.id || index}
            style={{
              border: '1px solid #f0f0f0',
              borderRadius: 4,
              padding: 12,
              marginBottom: 8,
              backgroundColor: '#fafafa',
            }}
          >
            <div style={{ width: '100%' }}>
              <div style={{ marginBottom: 8, fontWeight: 'bold' }}>
                第 {index + 1} 作者
              </div>
              <Space direction="vertical" style={{ width: '100%' }} size="small">
                <Space size="small" style={{ width: '100%' }}>
                  <Input
                    placeholder="作者姓名"
                    value={author.name}
                    onChange={e => updateAuthor(index, 'name', e.target.value)}
                    disabled={mode === 'view'}
                    style={{ flex: 1 }}
                  />
                  <Select
                    placeholder="作者类型"
                    value={author.author_type}
                    onChange={val => {
                      checkAuthorTypeConflict(index, val);
                      updateAuthor(index, 'author_type', val);
                    }}
                    disabled={mode === 'view'}
                    style={{ width: 150 }}
                  >
                    {authorTypeOptions.map(option => (
                      <Option key={option.value} value={option.value}>
                        {option.label}
                      </Option>
                    ))}
                  </Select>
                  {mode !== 'view' && (
                    <Space>
                      <Button
                        type="text"
                        size="small"
                        icon={<UpOutlined />}
                        onClick={() => moveUp(index)}
                        disabled={index === 0}
                      />
                      <Button
                        type="text"
                        size="small"
                        icon={<DownOutlined />}
                        onClick={() => moveDown(index)}
                        disabled={index === value.length - 1}
                      />
                      <Popconfirm
                        title="确定删除该作者吗?"
                        onConfirm={() => removeAuthor(index)}
                        okText="确定"
                        cancelText="取消"
                      >
                        <Button type="text" size="small" danger icon={<MinusCircleOutlined />} />
                      </Popconfirm>
                    </Space>
                  )}
                </Space>

                <Select
                  placeholder="选择关联用户(可选)"
                  value={author.user_id}
                  onSearch={loadUsers}
                  onChange={val => updateAuthor(index, 'user_id', val)}
                  disabled={mode === 'view'}
                  showSearch
                  filterOption={false}
                  loading={loadingUsers}
                  allowClear
                  style={{ width: '100%' }}
                >
                  {users.map(user => (
                    <Option key={user.id} value={user.id}>
                      {user.name} ({user.username})
                    </Option>
                  ))}
                </Select>

                <Input
                  placeholder="单位"
                  value={author.department}
                  onChange={e => updateAuthor(index, 'department', e.target.value)}
                  disabled={mode === 'view'}
                />
              </Space>
            </div>
          </List.Item>
        )}
      />
    </div>
  );
};

export default AuthorList;
