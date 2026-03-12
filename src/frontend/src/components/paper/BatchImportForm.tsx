// 批量导入表单组件
import React, { useState } from 'react';
import { Upload, Button, Alert, Table, Space, message } from 'antd';
import { UploadOutlined, DownloadOutlined } from '@ant-design/icons';
import type { UploadFile, UploadProps } from 'antd';
import { paperService } from '../../services/paperService';

interface BatchImportFormProps {
  onImportComplete?: () => void;
}

interface ImportError {
  row: number;
  field: string;
  error: string;
}

interface ImportResult {
  success: number;
  failed: number;
  errors: ImportError[];
}

const BatchImportForm: React.FC<BatchImportFormProps> = ({ onImportComplete }) => {
  const [file, setFile] = useState<UploadFile | null>(null);
  const [importResult, setImportResult] = useState<ImportResult | null>(null);
  const [isImporting, setIsImporting] = useState(false);

  // 下载导入模板
  const handleDownloadTemplate = async () => {
    try {
      const response = await paperService.downloadImportTemplate();
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', '论文导入模板.xlsx');
      document.body.appendChild(link);
      link.click();
      link.remove();
      message.success('模板下载成功');
    } catch (error) {
      message.error('模板下载失败');
    }
  };

  // 文件上传前校验
  const beforeUpload: UploadProps['beforeUpload'] = (uploadFile) => {
    const isExcel = uploadFile.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' || 
                    uploadFile.type === 'application/vnd.ms-excel';
    const isLt10M = uploadFile.size / 1024 / 1024 < 10;

    if (!isExcel) {
      message.error('只能上传 Excel 文件 (.xlsx, .xls)');
      return false;
    }

    if (!isLt10M) {
      message.error('文件大小不能超过 10MB');
      return false;
    }

    setFile(uploadFile);
    setImportResult(null);
    return false;
  };

  // 开始导入
  const handleImport = async () => {
    if (!file) {
      message.warning('请先选择要导入的文件');
      return;
    }

    setIsImporting(true);
    setImportResult(null);

    try {
      const fileObj = file.originFileObj || (file as any).raw;
      const response = await paperService.batchImport(fileObj);
      
      const result: ImportResult = {
        success: response.data.success || 0,
        failed: response.data.failed || 0,
        errors: []
      };

      // 解析错误信息
      if (response.data.errors && Array.isArray(response.data.errors)) {
        response.data.errors.forEach((errorStr: string, index: number) => {
          // 尝试解析错误格式：行号|字段|错误原因
          const parts = errorStr.split('|');
          if (parts.length >= 3) {
            result.errors.push({
              row: parseInt(parts[0]) || index + 1,
              field: parts[1],
              error: parts[2]
            });
          } else {
            result.errors.push({
              row: index + 1,
              field: '未知',
              error: errorStr
            });
          }
        });
      }

      setImportResult(result);

      if (result.success > 0) {
        message.success(`导入成功：${result.success} 条`);
      }

      if (result.failed > 0) {
        message.warning(`导入失败：${result.failed} 条`);
      }

      // 调用完成回调
      if (onImportComplete) {
        onImportComplete();
      }
    } catch (error: any) {
      const errorMsg = error.response?.data?.message || '导入失败，请检查文件格式';
      message.error(errorMsg);
      setImportResult({
        success: 0,
        failed: 1,
        errors: [{
          row: 0,
          field: '文件',
          error: errorMsg
        }]
      });
    } finally {
      setIsImporting(false);
    }
  };

  // 错误表格列定义
  const errorColumns = [
    {
      title: '行号',
      dataIndex: 'row',
      key: 'row',
      width: 80,
    },
    {
      title: '错误字段',
      dataIndex: 'field',
      key: 'field',
      width: 120,
    },
    {
      title: '错误原因',
      dataIndex: 'error',
      key: 'error',
    },
  ];

  // 重新导入（清除当前结果）
  const handleReimport = () => {
    setFile(null);
    setImportResult(null);
  };

  return (
    <div>
      {/* 下载模板按钮 */}
      <div style={{ marginBottom: 16 }}>
        <Button 
          type="primary" 
          icon={<DownloadOutlined />} 
          onClick={handleDownloadTemplate}
        >
          下载导入模板
        </Button>
      </div>

      {/* 文件上传 */}
      <Upload
        fileList={file ? [file] : []}
        beforeUpload={beforeUpload}
        onRemove={() => {
          setFile(null);
          setImportResult(null);
        }}
        accept=".xlsx,.xls"
        maxCount={1}
        disabled={isImporting}
      >
        <Button icon={<UploadOutlined />} disabled={isImporting}>
          选择 Excel 文件
        </Button>
      </Upload>

      {/* 导入按钮 */}
      <div style={{ marginTop: 16 }}>
        <Space>
          <Button 
            type="primary" 
            onClick={handleImport} 
            loading={isImporting}
            disabled={!file || isImporting}
          >
            {isImporting ? '导入中...' : '开始导入'}
          </Button>
          {importResult && (
            <Button onClick={handleReimport}>
              重新导入
            </Button>
          )}
        </Space>
      </div>

      {/* 导入结果展示 */}
      {importResult && (
        <div style={{ marginTop: 24 }}>
          {/* 结果摘要 */}
          {importResult.success > 0 && (
            <Alert
              message="导入成功"
              description={`成功导入 ${importResult.success} 条记录`}
              type="success"
              showIcon
              style={{ marginBottom: 16 }}
            />
          )}
          
          {importResult.failed > 0 && (
            <Alert
              message="导入失败"
              description={`有 ${importResult.failed} 条记录导入失败，请查看下方错误详情`}
              type="error"
              showIcon
              style={{ marginBottom: 16 }}
            />
          )}

          {/* 错误详情表格 */}
          {importResult.errors.length > 0 && (
            <div>
              <h4>错误详情</h4>
              <Table
                dataSource={importResult.errors.map((err, index) => ({ ...err, key: index }))}
                columns={errorColumns}
                pagination={false}
                size="small"
                scroll={{ y: 300 }}
              />
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default BatchImportForm;
