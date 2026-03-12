// 统计分析页面
import { useState, useEffect, useCallback } from 'react';
import { Card, Row, Col, Statistic, Tabs, Spin, Select, Table, message, Space, Typography } from 'antd';
import { BarChart, LineChart, PieChart, TableChart } from '../../components/common/StatsCharts';
import ExportButton from '../../components/common/ExportButton';
import { statsService } from '../../services/statsService';
import { useAuth } from '../../hooks/useAuth';
import type { BasicStats, AuthorStats, ProjectStats, DepartmentStats, YearlyStats, JournalStats } from '../../types/statistics';

const { Title, Text } = Typography;

// 收录类型映射
const INDEXING_TYPE_MAP: Record<string, string> = {
  SCI: 'SCI',
  EI: 'EI',
  CI: 'CI',
  DI: 'DI',
  CORE: 'CORE',
};

const StatsPage: React.FC = () => {
  const { hasPermission } = useAuth();
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('basic');

  // 基础统计数据
  const [basicStats, setBasicStats] = useState<BasicStats | null>(null);
  // 年度统计数据
  const [yearlyStats, setYearlyStats] = useState<YearlyStats | null>(null);
  // 期刊统计数据
  const [journalStats, setJournalStats] = useState<JournalStats | null>(null);

  // 按作者统计
  const [authorStats, setAuthorStats] = useState<AuthorStats | null>(null);
  const [selectedAuthorId, setSelectedAuthorId] = useState<number | null>(null);

  // 按课题统计
  const [projectStats, setProjectStats] = useState<ProjectStats | null>(null);
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);

  // 按单位统计
  const [departmentStats, setDepartmentStats] = useState<DepartmentStats | null>(null);
  const [selectedDepartment, setSelectedDepartment] = useState<string>('');

  // 是否有导出权限
  const canExport = hasPermission('stats:export');

  // 加载基础统计数据
  const fetchBasicStats = useCallback(async () => {
    try {
      const data = await statsService.getBasicStats();
      setBasicStats(data);
    } catch (error) {
      console.error('获取基础统计失败:', error);
      message.error('获取基础统计数据失败');
    }
  }, []);

  // 加载年度统计数据
  const fetchYearlyStats = useCallback(async () => {
    try {
      const data = await statsService.getYearlyStats();
      setYearlyStats(data);
    } catch (error) {
      console.error('获取年度统计失败:', error);
      message.error('获取年度统计数据失败');
    }
  }, []);

  // 加载期刊统计数据
  const fetchJournalStats = useCallback(async () => {
    try {
      const data = await statsService.getJournalStats();
      setJournalStats(data);
    } catch (error) {
      console.error('获取期刊统计失败:', error);
      message.error('获取期刊统计数据失败');
    }
  }, []);

  // 加载作者统计数据
  const fetchAuthorStats = useCallback(async (authorId: number) => {
    if (!authorId) return;
    try {
      const data = await statsService.getAuthorStats(authorId);
      setAuthorStats(data);
    } catch (error) {
      console.error('获取作者统计失败:', error);
      message.error('获取作者统计数据失败');
    }
  }, []);

  // 加载课题统计数据
  const fetchProjectStats = useCallback(async (projectId: number) => {
    if (!projectId) return;
    try {
      const data = await statsService.getProjectStats(projectId);
      setProjectStats(data);
    } catch (error) {
      console.error('获取课题统计失败:', error);
      message.error('获取课题统计数据失败');
    }
  }, []);

  // 加载单位统计数据
  const fetchDepartmentStats = useCallback(async (department: string) => {
    if (!department) return;
    try {
      const data = await statsService.getDepartmentStats(department);
      setDepartmentStats(data);
    } catch (error) {
      console.error('获取单位统计失败:', error);
      message.error('获取单位统计数据失败');
    }
  }, []);

  // 初始化加载数据
  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      await Promise.all([fetchBasicStats(), fetchYearlyStats(), fetchJournalStats()]);
      setLoading(false);
    };
    loadData();
  }, [fetchBasicStats, fetchYearlyStats, fetchJournalStats]);

  // Tab 切换处理
  const handleTabChange = (key: string) => {
    setActiveTab(key);
  };

  // 处理作者选择变化
  const handleAuthorChange = (authorId: number) => {
    setSelectedAuthorId(authorId);
    fetchAuthorStats(authorId);
  };

  // 处理课题选择变化
  const handleProjectChange = (projectId: number) => {
    setSelectedProjectId(projectId);
    fetchProjectStats(projectId);
  };

  // 处理单位选择变化
  const handleDepartmentChange = (department: string) => {
    setSelectedDepartment(department);
    fetchDepartmentStats(department);
  };

  // 将 yearlyCounts 转换为图表数据格式
  const getYearlyChartData = () => {
    if (!basicStats?.yearlyCounts) return [];
    return Object.entries(basicStats.yearlyCounts)
      .map(([year, count]) => ({ name: year, value: count }))
      .sort((a, b) => Number(a.name) - Number(b.name));
  };

  // 将 typeCounts 转换为图表数据格式
  const getTypeChartData = () => {
    if (!basicStats?.typeCounts) return [];
    return Object.entries(basicStats.typeCounts).map(([type, count]) => ({
      name: INDEXING_TYPE_MAP[type] || type,
      value: count,
    }));
  };

  // 将 journalCounts 转换为图表数据格式
  const getJournalChartData = () => {
    if (!basicStats?.journalCounts) return [];
    return basicStats.journalCounts.slice(0, 10).map((item) => ({
      name: item.journal,
      value: item.count,
    }));
  };

  // 加载状态
  if (loading) {
    return (
      <div style={{ padding: '24px', minHeight: '100vh', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  // 基础统计 Tab 内容
  const renderBasicStats = () => (
    <div>
      {/* 基础统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="论文总数"
              value={basicStats?.totalPapers || 0}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="平均影响因子"
              value={basicStats?.avgImpactFactor || 0}
              precision={2}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总引用次数"
              value={basicStats?.totalCitations || 0}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总他引次数"
              value={basicStats?.totalSelfCitations || 0}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 图表区域 */}
      <Row gutter={[16, 16]}>
        {/* 年度发表数量柱状图 */}
        <Col xs={24} lg={12}>
          <BarChart
            title="各年份发表数量"
            data={getYearlyChartData()}
            height={350}
          />
        </Col>

        {/* 收录类型分布饼图 */}
        <Col xs={24} lg={12}>
          <PieChart
            title="收录类型分布"
            data={getTypeChartData()}
            height={350}
          />
        </Col>

        {/* 年度趋势折线图 */}
        <Col xs={24}>
          <LineChart
            title="年度发表趋势"
            data={getYearlyChartData()}
            smooth
            area
            height={350}
          />
        </Col>

        {/* 期刊分布饼图 */}
        <Col xs={24} lg={12}>
          <PieChart
            title="期刊分布 (Top 10)"
            data={getJournalChartData()}
            height={350}
          />
        </Col>

        {/* 期刊详细表格 */}
        <Col xs={24} lg={12}>
          <TableChart
            title="期刊详细数据"
            columns={[
              { title: '期刊名称', dataIndex: 'journal', key: 'journal', width: '60%' },
              { title: '论文数量', dataIndex: 'count', key: 'count', width: '40%', sorter: (a: any, b: any) => a.count - b.count },
            ]}
            data={basicStats?.journalCounts || []}
            pagination
            pageSize={10}
          />
        </Col>
      </Row>
    </div>
  );

  // 按作者统计 Tab 内容
  const renderAuthorStats = () => (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <Space direction="vertical" style={{ width: '100%' }}>
          <Text>选择作者查看统计数据：</Text>
          <Select
            style={{ width: 300 }}
            placeholder="请选择作者"
            showSearch
            allowClear
            onChange={handleAuthorChange}
            options={[
              { value: 1, label: '张三' },
              { value: 2, label: '李四' },
              { value: 3, label: '王五' },
            ]}
          />
        </Space>
      </Card>

      {authorStats && (
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="论文总数"
                value={authorStats.paperCount}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="第一作者数量"
                value={authorStats.firstAuthorCount}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="通讯作者数量"
                value={authorStats.correspondingAuthorCount}
                valueStyle={{ color: '#fa8c16' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="平均影响因子"
                value={authorStats.avgImpactFactor}
                precision={2}
                valueStyle={{ color: '#722ed1' }}
              />
            </Card>
          </Col>
          <Col xs={24}>
            <Card>
              <Statistic
                title="总引用次数"
                value={authorStats.totalCitations}
                valueStyle={{ color: '#13c2c2' }}
              />
            </Card>
          </Col>
        </Row>
      )}
    </div>
  );

  // 按课题统计 Tab 内容
  const renderProjectStats = () => (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <Space direction="vertical" style={{ width: '100%' }}>
          <Text>选择课题查看统计数据：</Text>
          <Select
            style={{ width: 300 }}
            placeholder="请选择课题"
            showSearch
            allowClear
            onChange={handleProjectChange}
            options={[
              { value: 1, label: '课题A' },
              { value: 2, label: '课题B' },
              { value: 3, label: '课题C' },
            ]}
          />
        </Space>
      </Card>

      {projectStats && (
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="论文总数"
                value={projectStats.paperCount}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="高影响因子论文 (IF≥5)"
                value={projectStats.highImpactPaperCount}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="SCI论文数量"
                value={projectStats.sciPaperCount}
                valueStyle={{ color: '#fa8c16' }}
              />
            </Card>
          </Col>
        </Row>
      )}
    </div>
  );

  // 按单位统计 Tab 内容
  const renderDepartmentStats = () => (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <Space direction="vertical" style={{ width: '100%' }}>
          <Text>选择单位查看统计数据：</Text>
          <Select
            style={{ width: 300 }}
            placeholder="请选择单位"
            showSearch
            allowClear
            onChange={handleDepartmentChange}
            options={[
              { value: '计算机学院', label: '计算机学院' },
              { value: '信息学院', label: '信息学院' },
              { value: '软件学院', label: '软件学院' },
            ]}
          />
        </Space>
      </Card>

      {departmentStats && (
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} lg={8}>
            <Card>
              <Statistic
                title="论文总数"
                value={departmentStats.paperCount}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={8}>
            <Card>
              <Statistic
                title="总影响因子"
                value={departmentStats.totalImpactFactor}
                precision={2}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={8}>
            <Card>
              <Statistic
                title="总引用次数"
                value={departmentStats.totalCitations}
                valueStyle={{ color: '#fa8c16' }}
              />
            </Card>
          </Col>
        </Row>
      )}
    </div>
  );

  // Tab 项配置
  const tabItems = [
    {
      key: 'basic',
      label: '基础统计',
      children: renderBasicStats(),
    },
    {
      key: 'author',
      label: '按作者统计',
      children: renderAuthorStats(),
    },
    {
      key: 'project',
      label: '按课题统计',
      children: renderProjectStats(),
    },
    {
      key: 'department',
      label: '按单位统计',
      children: renderDepartmentStats(),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      {/* 页面标题和导出按钮 */}
      <Card style={{ marginBottom: 24 }}>
        <Row justify="space-between" align="middle">
          <Col>
            <Title level={4} style={{ margin: 0 }}>
              统计分析
            </Title>
          </Col>
          <Col>
            {canExport && (
              <ExportButton
                type="stats"
                statsType={activeTab}
              />
            )}
          </Col>
        </Row>
      </Card>

      {/* 统计维度切换 */}
      <Card>
        <Tabs
          activeKey={activeTab}
          onChange={handleTabChange}
          items={tabItems}
        />
      </Card>
    </div>
  );
};

export default StatsPage;
