// 布局组件
import React, { useState } from 'react';
import { Layout, Menu, Avatar, Dropdown, Button, theme } from 'antd';
import {
  FileTextOutlined,
  CheckCircleOutlined,
  BarChartOutlined,
  TeamOutlined,
  SettingOutlined,
  UserOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UploadOutlined
} from '@ant-design/icons';
import { useNavigate, Outlet, useLocation } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

const { Header, Sider, Content } = Layout;

const MainLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout, hasRole } = useAuth();
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  // 菜单项配置
  const menuItems = [
    {
      key: '/papers',
      icon: <FileTextOutlined />,
      label: '论文管理',
      children: [
        {
          key: '/papers',
          label: '论文列表',
        },
        {
          key: '/papers/create',
          label: '录入论文',
        },
        {
          key: '/papers/batch-import',
          label: '批量导入',
          icon: <UploadOutlined />,
        },
      ],
    },
    ...(hasRole('business_reviewer') || hasRole('political_reviewer') ? [
      {
        key: '/reviews',
        icon: <CheckCircleOutlined />,
        label: '审核管理',
        children: [
          ...(hasRole('business_reviewer') ? [
            {
              key: '/reviews/business',
              label: '业务审核',
            },
          ] : []),
          ...(hasRole('political_reviewer') ? [
            {
              key: '/reviews/political',
              label: '政工审核',
            },
          ] : []),
        ],
      },
    ] : []),
    {
      key: '/statistics',
      icon: <BarChartOutlined />,
      label: '统计分析',
    },
    ...(hasRole('admin') ? [
      {
        key: '/system',
        icon: <SettingOutlined />,
        label: '系统管理',
        children: [
          {
            key: '/system/users',
            label: '用户管理',
            icon: <TeamOutlined />,
          },
          {
            key: '/system/projects',
            label: '课题管理',
          },
          {
            key: '/system/journals',
            label: '期刊管理',
          },
          {
            key: '/system/config',
            label: '系统配置',
          },
        ],
      },
    ] : []),
  ];

  // 用户下拉菜单
  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人信息',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: () => logout(),
    },
  ];

  // 获取当前选中的菜单key
  const getSelectedKey = (): string => {
    const path = location.pathname;
    // 如果是详情页或编辑页,返回到列表页
    if (path.startsWith('/papers/') && path !== '/papers/create' && path !== '/papers/batch-import') {
      return '/papers';
    }
    if (path.startsWith('/reviews/business/')) {
      return '/reviews/business';
    }
    if (path.startsWith('/reviews/political/')) {
      return '/reviews/political';
    }
    return path;
  };

  // 获取展开的菜单key
  const getOpenKeys = (): string[] => {
    const path = location.pathname;
    if (path.startsWith('/papers')) {
      return ['/papers'];
    }
    if (path.startsWith('/reviews')) {
      return ['/reviews'];
    }
    if (path.startsWith('/statistics')) {
      return ['/statistics'];
    }
    if (path.startsWith('/system')) {
      return ['/system'];
    }
    return [];
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed}>
        <div style={{
          height: 64,
          margin: 16,
          background: 'rgba(255, 255, 255, 0.1)',
          borderRadius: 6,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: '#fff',
          fontSize: collapsed ? 14 : 18,
          fontWeight: 'bold',
        }}>
          {collapsed ? 'BioLit' : 'BioLit Manager'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[getSelectedKey()]}
          defaultOpenKeys={getOpenKeys()}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            padding: '0 16px',
            background: colorBgContainer,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            borderBottom: '1px solid #f0f0f0',
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: '16px',
              width: 64,
              height: 64,
            }}
          />
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              cursor: 'pointer',
              padding: '8px 16px',
              borderRadius: 4,
              transition: 'background-color 0.3s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#f5f5f5';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = 'transparent';
            }}
            >
              <Avatar icon={<UserOutlined />} />
              <span style={{ fontSize: 14 }}>
                {user?.name || user?.username || '用户'}
              </span>
            </div>
          </Dropdown>
        </Header>
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            minHeight: 280,
            background: colorBgContainer,
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout;
