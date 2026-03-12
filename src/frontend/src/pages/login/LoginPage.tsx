import React, { useState } from 'react';
import { Form, Input, Button, Card, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import './LoginPage.css';

const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const [loading, setLoading] = useState(false);

  const from = (location.state as any)?.from?.pathname || '/';

  const validatePassword = (_: any, value: string) => {
    if (!value) {
      return Promise.reject('请输入密码');
    }
    if (value.length < 8) {
      return Promise.reject('密码长度至少8位');
    }
    if (!/[A-Z]/.test(value)) {
      return Promise.reject('密码需包含大写字母');
    }
    if (!/[a-z]/.test(value)) {
      return Promise.reject('密码需包含小写字母');
    }
    if (!/[0-9]/.test(value)) {
      return Promise.reject('密码需包含数字');
    }
    if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(value)) {
      return Promise.reject('密码需包含特殊字符');
    }
    return Promise.resolve();
  };

  const handleLogin = async (values: { username: string; password: string }) => {
    setLoading(true);
    try {
      await login(values.username, values.password);
      message.success('登录成功');
      navigate(from, { replace: true });
    } catch (error: any) {
      const errorMsg = error?.message || '登录失败，请检查用户名和密码';
      if (errorMsg.includes('账户锁定')) {
        message.error(errorMsg);
      } else if (errorMsg.includes('账户已禁用')) {
        message.error(errorMsg);
      } else {
        message.error(errorMsg);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <Card className="login-card" title="文献管理系统" bordered={false}>
        <Form name="login" onFinish={handleLogin} autoComplete="off" size="large">
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ validator: validatePassword }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              登录
            </Button>
          </Form.Item>
        </Form>

        <div className="login-tips">
          <p>默认管理员账户: admin / Admin@123</p>
        </div>
      </Card>
    </div>
  );
};

export default LoginPage;
