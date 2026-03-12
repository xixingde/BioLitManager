// 统计图表组件
import React, { useRef, useCallback } from 'react';
import * as echarts from 'echarts';
import type { ECharts } from 'echarts';
import { Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { Card } from 'antd';

// 默认颜色配置
const DEFAULT_COLORS = [
  '#5470c6',
  '#91cc75',
  '#fac858',
  '#ee6666',
  '#73c0de',
  '#3ba272',
  '#fc8452',
  '#9a60b4',
  '#ea7ccc',
];

// 柱状图Props
export interface BarChartProps {
  title: string;
  data: Array<{ name: string; value: number }>;
  xAxisKey?: string;
  yAxisKey?: string;
  horizontal?: boolean;
  height?: number;
  colors?: string[];
}

// 折线图Props
export interface LineChartProps {
  title: string;
  data: Array<{ name: string; value: number }>;
  xAxisKey?: string;
  yAxisKey?: string;
  smooth?: boolean;
  area?: boolean;
  height?: number;
  colors?: string[];
}

// 饼图Props
export interface PieChartProps {
  title: string;
  data: Array<{ name: string; value: number }>;
  radius?: number;
  height?: number;
  colors?: string[];
}

// 表格Props
export interface TableChartProps {
  title: string;
  columns: ColumnsType<any>;
  data: any[];
  pagination?: boolean;
  pageSize?: number;
}

// 柱状图组件
const BarChart: React.FC<BarChartProps> = ({
  title,
  data,
  horizontal = false,
  height = 300,
  colors = DEFAULT_COLORS,
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<ECharts | null>(null);

  const initChart = useCallback(() => {
    if (!chartRef.current) return;

    chartInstance.current = echarts.init(chartRef.current);
    
    const option: echarts.EChartsOption = {
      title: {
        text: title,
        left: 'center',
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow',
        },
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: horizontal
        ? {
            type: 'value',
            boundaryGap: [0, 0.01],
          }
        : {
            type: 'category',
            data: data.map(item => item.name),
          },
      yAxis: horizontal
        ? {
            type: 'category',
            data: data.map(item => item.name),
          }
        : {
            type: 'value',
          },
      series: [
        {
          name: title,
          type: 'bar',
          data: data.map(item => item.value),
          itemStyle: {
            color: (params: any) => {
              return colors[params.dataIndex % colors.length];
            },
          },
          label: {
            show: true,
            position: horizontal ? 'right' : 'top',
          },
        },
      ],
    };

    chartInstance.current.setOption(option);
  }, [title, data, horizontal, colors]);

  React.useEffect(() => {
    initChart();

    const handleResize = () => {
      chartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chartInstance.current?.dispose();
    };
  }, [initChart]);

  // 数据变化时更新图表
  React.useEffect(() => {
    if (chartInstance.current) {
      chartInstance.current.setOption({
        xAxis: horizontal
          ? {
              type: 'value',
              boundaryGap: [0, 0.01],
            }
          : {
              type: 'category',
              data: data.map(item => item.name),
            },
        yAxis: horizontal
          ? {
              type: 'category',
              data: data.map(item => item.name),
            }
          : {
              type: 'value',
            },
        series: [
          {
            data: data.map(item => item.value),
          },
        ],
      });
    }
  }, [data, horizontal]);

  return (
    <Card style={{ marginBottom: 16 }}>
      <div ref={chartRef} style={{ height, width: '100%' }} />
    </Card>
  );
};

// 折线图组件
const LineChart: React.FC<LineChartProps> = ({
  title,
  data,
  smooth = false,
  area = false,
  height = 300,
  colors = DEFAULT_COLORS,
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<ECharts | null>(null);

  const initChart = useCallback(() => {
    if (!chartRef.current) return;

    chartInstance.current = echarts.init(chartRef.current);

    const option: echarts.EChartsOption = {
      title: {
        text: title,
        left: 'center',
      },
      tooltip: {
        trigger: 'axis',
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: data.map(item => item.name),
      },
      yAxis: {
        type: 'value',
      },
      series: [
        {
          name: title,
          type: 'line',
          smooth,
          data: data.map(item => item.value),
          itemStyle: {
            color: colors[0],
          },
          areaStyle: area
            ? {
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  { offset: 0, color: colors[0] },
                  { offset: 1, color: 'rgba(255, 255, 255, 0)' },
                ]),
              }
            : undefined,
          label: {
            show: true,
            position: 'top',
          },
        },
      ],
    };

    chartInstance.current.setOption(option);
  }, [title, data, smooth, area, colors]);

  React.useEffect(() => {
    initChart();

    const handleResize = () => {
      chartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chartInstance.current?.dispose();
    };
  }, [initChart]);

  React.useEffect(() => {
    if (chartInstance.current) {
      chartInstance.current.setOption({
        xAxis: {
          data: data.map(item => item.name),
        },
        series: [
          {
            data: data.map(item => item.value),
          },
        ],
      });
    }
  }, [data]);

  return (
    <Card style={{ marginBottom: 16 }}>
      <div ref={chartRef} style={{ height, width: '100%' }} />
    </Card>
  );
};

// 饼图组件
const PieChart: React.FC<PieChartProps> = ({
  title,
  data,
  radius = '60%',
  height = 300,
  colors = DEFAULT_COLORS,
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<ECharts | null>(null);

  const initChart = useCallback(() => {
    if (!chartRef.current) return;

    chartInstance.current = echarts.init(chartRef.current);

    const option: echarts.EChartsOption = {
      title: {
        text: title,
        left: 'center',
      },
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)',
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        top: 'middle',
      },
      series: [
        {
          name: title,
          type: 'pie',
          radius: typeof radius === 'number' ? `${radius}%` : radius,
          data: data.map((item, index) => ({
            value: item.value,
            name: item.name,
            itemStyle: {
              color: colors[index % colors.length],
            },
          })),
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)',
            },
          },
          label: {
            show: true,
            formatter: '{b}: {c} ({d}%)',
          },
        },
      ],
    };

    chartInstance.current.setOption(option);
  }, [title, data, radius, colors]);

  React.useEffect(() => {
    initChart();

    const handleResize = () => {
      chartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chartInstance.current?.dispose();
    };
  }, [initChart]);

  React.useEffect(() => {
    if (chartInstance.current) {
      chartInstance.current.setOption({
        series: [
          {
            data: data.map((item, index) => ({
              value: item.value,
              name: item.name,
              itemStyle: {
                color: colors[index % colors.length],
              },
            })),
          },
        ],
      });
    }
  }, [data, colors]);

  return (
    <Card style={{ marginBottom: 16 }}>
      <div ref={chartRef} style={{ height, width: '100%' }} />
    </Card>
  );
};

// 表格组件
const TableChart: React.FC<TableChartProps> = ({
  title,
  columns,
  data,
  pagination = true,
  pageSize = 10,
}) => {
  return (
    <Card title={title} style={{ marginBottom: 16 }}>
      <Table
        columns={columns}
        dataSource={data}
        pagination={
          pagination
            ? {
                pageSize,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total) => `共 ${total} 条`,
              }
            : false
        }
        scroll={{ x: 'max-content' }}
      />
    </Card>
  );
};

// 图表导出工具函数
export const exportChartToPng = (chartInstance: ECharts, filename: string): void => {
  const url = chartInstance.getDataURL({
    type: 'png',
    pixelRatio: 2,
    backgroundColor: '#fff',
  });
  
  const link = document.createElement('a');
  link.download = `${filename}.png`;
  link.href = url;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};

export { BarChart, LineChart, PieChart, TableChart };
