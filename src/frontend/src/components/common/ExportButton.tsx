// 导出按钮组件
import React, { useState, useEffect } from 'react';
import { Button, Dropdown, Modal, Checkbox, Spin, message } from 'antd';
import { DownloadOutlined, DownOutlined, FileExcelOutlined, FilePdfOutlined, FileWordOutlined } from '@ant-design/icons';
import { exportService, ExportPapersParams } from '../../services/exportService';
import { useAuth } from '../../hooks/useAuth';
import type { ExportField } from '../../types/statistics';
import { PERMISSION } from '../../utils/constants';

export interface ExportButtonProps {
  /** 导出类型 */
  type: 'papers' | 'paper' | 'stats';
  /** 单篇论文ID */
  paperId?: number;
  /** 论文ID列表（查询结果导出） */
  paperIds?: number[];
  /** 搜索参数 */
  searchParams?: Record<string, any>;
  /** 统计类型 */
  statsType?: string;
  /** 导出成功回调 */
  onSuccess?: () => void;
}

const ExportButton: React.FC<ExportButtonProps> = ({
  type,
  paperId,
  paperIds,
  searchParams,
  statsType,
  onSuccess,
}) => {
  const { hasPermission, user } = useAuth();
  const [loading, setLoading] = useState(false);
  const [fieldsModalVisible, setFieldsModalVisible] = useState(false);
  const [exportFields, setExportFields] = useState<ExportField[]>([]);
  const [selectedFields, setSelectedFields] = useState<string[]>([]);

  // 加载导出字段
  useEffect(() => {
    if (fieldsModalVisible && type === 'papers') {
      loadExportFields();
    }
  }, [fieldsModalVisible, type]);

  const loadExportFields = async () => {
    try {
      const response = await exportService.getExportFields();
      const fields = response.data || [];
      setExportFields(fields);
      // 默认选中常用字段
      const defaultSelected = fields
        .filter((f: ExportField) => f.selected)
        .map((f: ExportField) => f.key);
      setSelectedFields(defaultSelected);
    } catch (error) {
      console.error('加载导出字段失败:', error);
      message.error('加载导出字段失败');
    }
  };

  // 检查权限
  const checkPermission = (): boolean => {
    if (type === 'stats') {
      return hasPermission(PERMISSION.STATS_EXPORT);
    }
    if (type === 'paper') {
      // 单篇论文导出需要 paper:export 权限
      return hasPermission(PERMISSION.PAPER_EXPORT);
    }
    // papers 导出需要 paper:export 权限
    return hasPermission(PERMISSION.PAPER_EXPORT);
  };

  // 导出查询结果
  const handleExportPapers = async () => {
    if (selectedFields.length === 0) {
      message.warning('请至少选择一个导出字段');
      return;
    }

    setLoading(true);
    try {
      const params: ExportPapersParams = {
        searchParams,
        fields: selectedFields,
      };
      const response = await exportService.exportPapers(params);
      const { filePath } = response.data;
      
      // 下载文件
      const filename = `论文导出_${new Date().toLocaleDateString()}.xlsx`;
      exportService.downloadFile(filePath, filename);
      
      message.success('导出成功');
      onSuccess?.();
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败，请重试');
    } finally {
      setLoading(false);
      setFieldsModalVisible(false);
    }
  };

  // 导出单篇论文
  const handleExportPaper = async (format: 'pdf' | 'word') => {
    if (!paperId) {
      message.error('论文ID不能为空');
      return;
    }

    setLoading(true);
    try {
      const response = await exportService.exportPaper(paperId, format);
      const { filePath } = response.data;
      
      // 获取文件扩展名
      const ext = format === 'pdf' ? '.pdf' : '.docx';
      const filename = `论文_${paperId}${ext}`;
      exportService.downloadFile(filePath, filename);
      
      message.success('导出成功');
      onSuccess?.();
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 导出统计结果
  const handleExportStats = async (format: 'excel' | 'pdf') => {
    if (!statsType) {
      message.error('统计类型不能为空');
      return;
    }

    setLoading(true);
    try {
      const response = await exportService.exportStats(statsType, format);
      const { filePath } = response.data;
      
      // 获取文件扩展名
      const ext = format === 'excel' ? '.xlsx' : '.pdf';
      const filename = `统计_${statsType}_${new Date().toLocaleDateString()}${ext}`;
      exportService.downloadFile(filePath, filename);
      
      message.success('导出成功');
      onSuccess?.();
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 字段选择确认
  const handleFieldsConfirm = () => {
    handleExportPapers();
  };

  // 字段选择变化
  const handleFieldChange = (checkedValues: string[]) => {
    setSelectedFields(checkedValues);
  };

  // 全选/取消全选
  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedFields(exportFields.map((f: ExportField) => f.key));
    } else {
      setSelectedFields([]);
    }
  };

  // 渲染导出菜单项
  const renderMenuItems = () => {
    const items: { key: string; label: React.ReactNode; onClick: () => void }[] = [];

    if (type === 'papers') {
      // 查询结果导出 - Excel格式
      items.push({
        key: 'excel',
        label: (
          <span>
            <FileExcelOutlined style={{ marginRight: 8 }} />
            导出为 Excel
          </span>
        ),
        onClick: () => setFieldsModalVisible(true),
      });
    } else if (type === 'paper') {
      // 单篇论文导出 - PDF 和 Word
      items.push({
        key: 'pdf',
        label: (
          <span>
            <FilePdfOutlined style={{ marginRight: 8 }} />
            导出为 PDF
          </span>
        ),
        onClick: () => handleExportPaper('pdf'),
      });
      items.push({
        key: 'word',
        label: (
          <span>
            <FileWordOutlined style={{ marginRight: 8 }} />
            导出为 Word
          </span>
        ),
        onClick: () => handleExportPaper('word'),
      });
    } else if (type === 'stats') {
      // 统计结果导出 - Excel 和 PDF
      items.push({
        key: 'excel',
        label: (
          <span>
            <FileExcelOutlined style={{ marginRight: 8 }} />
            导出为 Excel
          </span>
        ),
        onClick: () => handleExportStats('excel'),
      });
      items.push({
        key: 'pdf',
        label: (
          <span>
            <FilePdfOutlined style={{ marginRight: 8 }} />
            导出为 PDF
          </span>
        ),
        onClick: () => handleExportStats('pdf'),
      });
    }

    return items;
  };

  // 权限检查
  if (!checkPermission()) {
    return null;
  }

  const menuItems = renderMenuItems();

  return (
    <>
      <Spin spinning={loading}>
        <Dropdown
          menu={{ items: menuItems }}
          trigger={['click']}
          disabled={loading}
        >
          <Button type="primary">
            <DownloadOutlined />
            导出 <DownOutlined />
          </Button>
        </Dropdown>
      </Spin>

      {/* 字段选择弹窗 */}
      <Modal
        title="选择导出字段"
        open={fieldsModalVisible}
        onOk={handleFieldsConfirm}
        onCancel={() => setFieldsModalVisible(false)}
        okText="确认导出"
        cancelText="取消"
        width={500}
        destroyOnClose
      >
        <div style={{ marginBottom: 16 }}>
          <Checkbox
            checked={selectedFields.length === exportFields.length}
            indeterminate={selectedFields.length > 0 && selectedFields.length < exportFields.length}
            onChange={(e) => handleSelectAll(e.target.checked)}
          >
            全选
          </Checkbox>
        </div>
        <Checkbox.Group
          value={selectedFields}
          onChange={handleFieldChange}
          style={{ width: '100%' }}
        >
          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            {exportFields.map((field: ExportField) => (
              <Checkbox key={field.key} value={field.key}>
                {field.label}
              </Checkbox>
            ))}
          </div>
        </Checkbox.Group>
      </Modal>
    </>
  );
};

export default ExportButton;
