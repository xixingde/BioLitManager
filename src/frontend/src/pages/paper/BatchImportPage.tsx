// 批量导入页面
import React, { useState } from 'react';
import { Card, Breadcrumb, Button, Upload, message, Steps, Table, Result, Space, Alert } from 'antd';
import { HomeOutlined, FileTextOutlined, DownloadOutlined, UploadOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import type { UploadFile } from 'antd';
import { paperService } from '../../services/paperService';

const { Step } = Steps;
const { Dragger } = Upload;

interface ImportError {
  row: number;
  field: string;
  message: string;
}

const BatchImportPage: React.FC = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [importing, setImporting] = useState(false);
  const [importResult, setImportResult] = useState<{
    success: number;
    failed: number;
    errors: string[];
  } | null>(null);

  // 下载导入模板
  const handleDownloadTemplate = async () => {
    try {
      const response = await paperService.downloadImportTemplate();
      // 创建blob链接并下载文件
      const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = '论文导入模板.xlsx';
      link.click();
      URL.revokeObjectURL(link.href);
      message.success('模板下载成功');
    } catch (error) {
      message.error('模板下载失败');
    }
  };

  // 文件上传前校验
  const beforeUpload = (file: File): boolean => {
    const isExcel =
      file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
      file.type === 'application/vnd.ms-excel' ||
      file.name.endsWith('.xlsx') ||
      file.name.endsWith('.xls') ||
      file.name.endsWith('.csv');

    if (!isExcel) {
      message.error('只能上传 Excel 或 CSV 文件');
      return false;
    }

    const isLt10M = file.size / 1024 / 1024 < 10;
    if (!isLt10M) {
      message.error('文件大小不能超过 10MB');
      return false;
    }

    return true;
  };

  // 处理文件选择
  const handleFileChange = (info: any) => {
    setFileList(info.fileList.slice(-1)); // 只保留最新的一个文件
    if (info.file.status === 'done') {
      message.success('文件上传成功');
      setCurrentStep(1);
    }
  };

  // 开始导入
  const handleStartImport = async () => {
    if (fileList.length === 0) {
      message.error('请先选择要导入的文件');
      return;
    }

    setImporting(true);
    try {
      const file = fileList[0].originFileObj as File;
      const response = await paperService.batchImport(file);
      setImportResult(response.data.data);
      setCurrentStep(2);
      message.success(`导入完成: 成功 ${response.data.data.success} 条, 失败 ${response.data.data.failed} 条`);
    } catch (error) {
      message.error('导入失败');
      setImportResult({
        success: 0,
        failed: 1,
        errors: ['导入过程中发生错误']
      });
      setCurrentStep(2);
    } finally {
      setImporting(false);
    }
  };

  // 重新导入
  const handleReImport = () => {
    setFileList([]);
    setImportResult(null);
    setCurrentStep(0);
  };

  // 返回列表
  const handleBackToList = () => {
    window.location.href = '/papers';
  };

  // 错误信息表格列
  const errorColumns = [
    {
      title: '错误信息',
      dataIndex: 'error',
      key: 'error',
      render: (error: string) => (
        <span style={{ color: '#ff4d4f' }}>
          {error}
        </span>
      )
    }
  ];

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
        <Breadcrumb.Item>批量导入</Breadcrumb.Item>
      </Breadcrumb>

      <Card title="批量导入论文">
        <Steps current={currentStep} style={{ marginBottom: 24 }}>
          <Step title="上传文件" description="下载模板并填写数据" />
          <Step title="确认导入" description="检查文件并开始导入" />
          <Step title="完成" description="查看导入结果" />
        </Steps>

        {/* 步骤1: 上传文件 */}
        {currentStep === 0 && (
          <div>
            <Alert
              message="导入说明"
              description={
                <div>
                  <p>1. 请先下载导入模板,按照模板格式填写数据</p>
                  <p>2. 论文标题、DOI、期刊ID为必填字段</p>
                  <p>3. 单次最多导入 1000 条数据</p>
                  <p>4. 支持 Excel (.xlsx, .xls) 和 CSV 格式</p>
                </div>
              }
              type="info"
              showIcon
              style={{ marginBottom: 16 }}
            />

            <Space style={{ marginBottom: 16 }}>
              <Button type="primary" icon={<DownloadOutlined />} onClick={handleDownloadTemplate}>
                下载导入模板
              </Button>
            </Space>

            <Dragger
              fileList={fileList}
              beforeUpload={beforeUpload}
              onChange={handleFileChange}
              onRemove={() => {
                setFileList([]);
              }}
              accept=".xlsx,.xls,.csv"
              multiple={false}
            >
              <p className="ant-upload-drag-icon">
                <UploadOutlined />
              </p>
              <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
              <p className="ant-upload-hint">支持 Excel 和 CSV 文件,文件大小不超过 10MB</p>
            </Dragger>
          </div>
        )}

        {/* 步骤2: 确认导入 */}
        {currentStep === 1 && (
          <div>
            <Alert
              message="确认导入"
              description={
                <div>
                  <p>已选择文件: {fileList[0]?.name}</p>
                  <p>文件大小: {(fileList[0]?.size || 0) / 1024 / 1024} MB</p>
                  <p>点击下方按钮开始导入,导入过程中请不要关闭页面</p>
                </div>
              }
              type="warning"
              showIcon
              style={{ marginBottom: 16 }}
            />

            <Space>
              <Button type="primary" loading={importing} onClick={handleStartImport}>
                开始导入
              </Button>
              <Button onClick={() => setCurrentStep(0)}>
                重新选择文件
              </Button>
            </Space>
          </div>
        )}

        {/* 步骤3: 完成 */}
        {currentStep === 2 && importResult && (
          <div>
            {importResult.failed === 0 ? (
              <Result
                icon={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
                title="导入成功"
                subTitle={`成功导入 ${importResult.success} 条数据`}
                extra={[
                  <Button type="primary" key="list" onClick={handleBackToList}>
                    查看论文列表
                  </Button>,
                  <Button key="reimport" onClick={handleReImport}>
                    继续导入
                  </Button>,
                ]}
              />
            ) : (
              <Result
                icon={<CloseCircleOutlined style={{ color: '#ff4d4f' }} />}
                title={`导入完成,部分数据导入失败`}
                subTitle={`成功: ${importResult.success} 条, 失败: ${importResult.failed} 条`}
                extra={[
                  <Button type="primary" key="list" onClick={handleBackToList}>
                    查看论文列表
                  </Button>,
                  <Button key="reimport" onClick={handleReImport}>
                    继续导入
                  </Button>,
                ]}
              >
                {importResult.errors && importResult.errors.length > 0 && (
                  <div style={{ marginTop: 24 }}>
                    <h4>错误详情:</h4>
                    <Table
                      columns={errorColumns}
                      dataSource={importResult.errors.map((error, index) => ({
                        key: index,
                        error
                      }))}
                      pagination={false}
                      size="small"
                    />
                  </div>
                )}
              </Result>
            )}
          </div>
        )}
      </Card>
    </div>
  );
};

export default BatchImportPage;
