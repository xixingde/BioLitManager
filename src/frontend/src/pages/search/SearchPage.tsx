// 搜索页面 - 多维度组合查询
import { useState, useEffect, useMemo } from 'react';
import { 
  Table, Input, Select, Space, Button, DatePicker, Tag, 
  message, Card, Row, Col, Divider, Empty 
} from 'antd';
import { SearchOutlined, EyeOutlined, PlusOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { searchService } from '../../services/searchService';
import type { PaperSearchResult, QueryCondition, LogicType, AuthorTypeFilter, SortField, SortOrder } from '../../types/paper';
import { useNavigate } from 'react-router-dom';
import ExportButton from '@/components/common/ExportButton';

const { RangePicker } = DatePicker;
const { Option } = Select;

// 字段选项
const fieldOptions = [
  { value: 'id', label: '论文ID' },
  { value: 'title', label: '标题' },
  { value: 'author_name', label: '作者' },
  { value: 'journal_name', label: '期刊' },
  { value: 'doi', label: 'DOI' },
  { value: 'project_code', label: '课题编号' },
  { value: 'project_name', label: '课题名称' },
];

// 作者类型选项
const authorTypeOptions = [
  { value: 'all', label: '全部' },
  { value: 'first_author', label: '第一作者' },
  { value: 'co_first_author', label: '共同第一作者' },
  { value: 'corresponding_author', label: '通讯作者' },
];

// 排序字段选项
const sortFieldOptions = [
  { value: 'publish_date', label: '出版日期' },
  { value: 'impact_factor', label: '影响因子' },
  { value: 'created_at', label: '创建时间' },
];

// 单个查询条件组件
interface ConditionItemProps {
  condition: QueryCondition;
  index: number;
  onChange: (index: number, condition: QueryCondition) => void;
  onRemove: (index: number) => void;
}

const ConditionItem: React.FC<ConditionItemProps> = ({ condition, index, onChange, onRemove }) => {
  return (
    <Card size="small" style={{ marginBottom: 8 }}>
      <Space align="center" style={{ width: '100%' }}>
        <span style={{ color: '#999' }}>条件 {index + 1}</span>
        <Select
          value={condition.field}
          style={{ width: 120 }}
          onChange={(value) => onChange(index, { ...condition, field: value })}
        >
          {fieldOptions.map(opt => (
            <Option key={opt.value} value={opt.value}>{opt.label}</Option>
          ))}
        </Select>
        <Select
          value={condition.operator}
          style={{ width: 100 }}
          onChange={(value) => onChange(index, { ...condition, operator: value })}
        >
          <Option value="eq">等于</Option>
          <Option value="ne">不等于</Option>
          <Option value="like">包含</Option>
          <Option value="gt">大于</Option>
          <Option value="gte">大于等于</Option>
          <Option value="lt">小于</Option>
          <Option value="lte">小于等于</Option>
        </Select>
        <Input
          placeholder="请输入值"
          style={{ width: 200 }}
          value={condition.value as string}
          onChange={(e) => onChange(index, { ...condition, value: e.target.value })}
        />
        <Button type="text" danger onClick={() => onRemove(index)}>删除</Button>
      </Space>
    </Card>
  );
};

const SearchPage: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<PaperSearchResult[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [size] = useState(20);
  
  // 搜索条件
  const [conditions, setConditions] = useState<QueryCondition[]>([]);
  const [logic, setLogic] = useState<LogicType>('AND');
  const [authorTypeFilter, setAuthorTypeFilter] = useState<AuthorTypeFilter>('all');
  
  // 排序
  const [sortField, setSortField] = useState<SortField>('publish_date');
  const [sortOrder, setSortOrder] = useState<SortOrder>('desc');
  
  // 简单搜索条件
  const [simpleKeyword, setSimpleKeyword] = useState('');
  const [dateRange, setDateRange] = useState<[string, string] | null>(null);
  const [impactFactorRange, setImpactFactorRange] = useState<[number, number] | null>(null);

  // 添加查询条件
  const addCondition = () => {
    setConditions([...conditions, { field: 'title', operator: 'like', value: '' }]);
  };

  // 修改查询条件
  const handleConditionChange = (index: number, condition: QueryCondition) => {
    const newConditions = [...conditions];
    newConditions[index] = condition;
    setConditions(newConditions);
  };

  // 删除查询条件
  const handleRemoveCondition = (index: number) => {
    setConditions(conditions.filter((_, i) => i !== index));
  };

  // 执行搜索
  const handleSearch = async (pageNum: number = 1) => {
    setLoading(true);
    try {
      // 构建查询条件
      const query: QueryCondition[] = [...conditions];
      
      // 添加简单搜索条件
      if (simpleKeyword) {
        query.push({ field: 'title', operator: 'like', value: simpleKeyword });
      }
      
      // 添加日期范围条件
      if (dateRange) {
        query.push({ field: 'publish_date', operator: 'between', value: dateRange });
      }
      
      // 添加影响因子范围条件
      if (impactFactorRange) {
        query.push({ field: 'impact_factor', operator: 'gte', value: impactFactorRange[0] });
        query.push({ field: 'impact_factor', operator: 'lte', value: impactFactorRange[1] });
      }
      
      const response = await searchService.advancedSearch({
        query: query.length > 0 ? query : undefined,
        logic,
        author_type_filter: authorTypeFilter,
        pagination: {
          page: pageNum,
          size
        },
        sort: {
          field: sortField,
          order: sortOrder
        }
      });
      
      setResults(response.data.list);
      setTotal(response.data.total);
      setPage(pageNum);
    } catch (error) {
      message.error('搜索失败');
    } finally {
      setLoading(false);
    }
  };

  // 重置搜索条件
  const handleReset = () => {
    setConditions([]);
    setSimpleKeyword('');
    setDateRange(null);
    setImpactFactorRange(null);
    setAuthorTypeFilter('all');
    setSortField('publish_date');
    setSortOrder('desc');
    setResults([]);
    setTotal(0);
    setPage(1);
  };

  // 构建当前搜索参数（用于导出）
  const searchParams = useMemo(() => {
    const query: QueryCondition[] = [...conditions];
    
    // 添加简单搜索条件
    if (simpleKeyword) {
      query.push({ field: 'title', operator: 'like', value: simpleKeyword });
    }
    
    // 添加日期范围条件
    if (dateRange) {
      query.push({ field: 'publish_date', operator: 'between', value: dateRange });
    }
    
    // 添加影响因子范围条件
    if (impactFactorRange) {
      query.push({ field: 'impact_factor', operator: 'gte', value: impactFactorRange[0] });
      query.push({ field: 'impact_factor', operator: 'lte', value: impactFactorRange[1] });
    }
    
    return {
      query: query.length > 0 ? query : undefined,
      logic,
      author_type_filter: authorTypeFilter,
      sort: {
        field: sortField,
        order: sortOrder
      }
    };
  }, [conditions, simpleKeyword, dateRange, impactFactorRange, logic, authorTypeFilter, sortField, sortOrder]);

  // 表格列定义
  const columns: ColumnsType<PaperSearchResult> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80
    },
    {
      title: '论文标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
      render: (text: string) => <a onClick={() => navigate(`/papers/${text}`)}>{text}</a>
    },
    {
      title: '作者列表',
      dataIndex: 'authors',
      key: 'authors',
      width: 200,
      render: (authors?: PaperSearchResult['authors']) => {
        if (!authors || authors.length === 0) return '-';
        return authors.map(a => a.name).join(', ');
      }
    },
    {
      title: '期刊',
      dataIndex: ['journal', 'short_name'],
      key: 'journal',
      width: 120
    },
    {
      title: '出版日期',
      dataIndex: 'publish_date',
      key: 'publish_date',
      width: 120,
      render: (date: string) => date ? new Date(date).toLocaleDateString('zh-CN') : '-'
    },
    {
      title: '影响因子',
      dataIndex: 'impact_factor',
      key: 'impact_factor',
      width: 100,
      sorter: true
    },
    {
      title: '审核状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          'draft': 'default',
          '待业务审核': 'blue',
          '待政工审核': 'processing',
          '审核通过': 'success',
          '驳回': 'error'
        };
        return <Tag color={colorMap[status] || 'default'}>{status}</Tag>;
      }
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      render: (_, record) => (
        <Button 
          type="link" 
          icon={<EyeOutlined />} 
          onClick={() => navigate(`/papers/${record.id}`)}
        >
          查看
        </Button>
      )
    }
  ];

  return (
    <div>
      <Card title="高级搜索" style={{ marginBottom: 16 }}>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          {/* 简单搜索区域 */}
          <Row gutter={16}>
            <Col span={8}>
              <Input
                placeholder="输入关键词搜索标题"
                prefix={<SearchOutlined />}
                value={simpleKeyword}
                onChange={(e) => setSimpleKeyword(e.target.value)}
                allowClear
              />
            </Col>
            <Col span={8}>
              <RangePicker
                style={{ width: '100%' }}
                onChange={(dates) => {
                  if (dates) {
                    setDateRange([
                      dates[0]?.format('YYYY-MM-DD') || '',
                      dates[1]?.format('YYYY-MM-DD') || ''
                    ]);
                  } else {
                    setDateRange(null);
                  }
                }}
              />
            </Col>
            <Col span={8}>
              <Space>
                <Input
                  placeholder="影响因子最小值"
                  type="number"
                  style={{ width: 120 }}
                  onChange={(e) => setImpactFactorRange(
                    impactFactorRange ? [Number(e.target.value), impactFactorRange[1]] : [Number(e.target.value), 0]
                  )}
                />
                <span>-</span>
                <Input
                  placeholder="影响因子最大值"
                  type="number"
                  style={{ width: 120 }}
                  onChange={(e) => setImpactFactorRange(
                    impactFactorRange ? [impactFactorRange[0], Number(e.target.value)] : [0, Number(e.target.value)]
                  )}
                />
              </Space>
            </Col>
          </Row>
          
          <Divider style={{ margin: '8px 0' }} />
          
          {/* 逻辑组合条件 */}
          <div>
            <Space>
              <span>逻辑组合：</span>
              <Select value={logic} onChange={setLogic} style={{ width: 100 }}>
                <Option value="AND">且 (AND)</Option>
                <Option value="OR">或 (OR)</Option>
                <Option value="NOT">非 (NOT)</Option>
              </Select>
              <Button type="dashed" onClick={addCondition}>添加条件</Button>
            </Space>
          </div>
          
          {/* 条件列表 */}
          {conditions.map((condition, index) => (
            <ConditionItem
              key={index}
              condition={condition}
              index={index}
              onChange={handleConditionChange}
              onRemove={handleRemoveCondition}
            />
          ))}
          
          <Divider style={{ margin: '8px 0' }} />
          
          {/* 作者类型筛选和排序 */}
          <Row gutter={16}>
            <Col span={8}>
              <Space>
                <span>作者类型：</span>
                <Select 
                  value={authorTypeFilter} 
                  onChange={setAuthorTypeFilter}
                  style={{ width: 150 }}
                >
                  {authorTypeOptions.map(opt => (
                    <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                  ))}
                </Select>
              </Space>
            </Col>
            <Col span={8}>
              <Space>
                <span>排序：</span>
                <Select 
                  value={sortField} 
                  onChange={setSortField}
                  style={{ width: 120 }}
                >
                  {sortFieldOptions.map(opt => (
                    <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                  ))}
                </Select>
                <Select 
                  value={sortOrder} 
                  onChange={setSortOrder}
                  style={{ width: 80 }}
                >
                  <Option value="desc">降序</Option>
                  <Option value="asc">升序</Option>
                </Select>
              </Space>
            </Col>
            <Col span={8} style={{ textAlign: 'right' }}>
              <Space>
                <Button onClick={handleReset}>重置</Button>
                <Button type="primary" icon={<SearchOutlined />} onClick={() => handleSearch(1)}>
                  搜索
                </Button>
              </Space>
            </Col>
          </Row>
        </Space>
      </Card>
      
      {/* 搜索结果 */}
      <Card 
        title={`搜索结果 (共 ${total} 条)`}
        extra={
          results.length > 0 ? (
            <ExportButton
              type="papers"
              searchParams={searchParams}
            />
          ) : null
        }
      >
        {results.length === 0 && !loading ? (
          <Empty description="未找到符合条件的论文" />
        ) : (
          <Table
            columns={columns}
            dataSource={results}
            rowKey="id"
            loading={loading}
            pagination={{
              current: page,
              pageSize: size,
              total,
              showSizeChanger: false,
              showQuickJumper: true,
              showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
              onChange: (pageNum) => handleSearch(pageNum)
            }}
            onRow={(record) => ({
              onClick: () => navigate(`/papers/${record.id}`),
              style: { cursor: 'pointer' }
            })}
          />
        )}
      </Card>
    </div>
  );
};

export default SearchPage;
