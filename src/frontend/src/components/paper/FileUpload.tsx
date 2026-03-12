// 文件上传组件
import React, { useState } from 'react';
import { Upload, Button, List, Tag, Space, message, Modal, Progress } from 'antd';
import { UploadOutlined, DownloadOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import type { UploadFile, UploadProps } from 'antd';
import { fileService } from '../../services/fileService';

// 文件类型选项
const fileTypeOptions = [
  { label: '全文', value: '全文' },
  { label: '首页', value: '首页' },
  { label: '期刊封面', value: '期刊封面' },
  { label: '审批件', value: '审批件' },
];

interface FileUploadProps {
  paperId: number;
  mode: 'create' | 'edit' | 'view';
}

const FileUpload: React.FC<FileUploadProps> = ({ paperId, mode }) => {
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [attachments, setAttachments] = useState<any[]>([]);
  const [uploading, setUploading] = useState(false);
  const [previewVisible, setPreviewVisible] = useState(false);
  const [previewFile, setPreviewFile] = useState<any>(null);

  // 加载附件列表
  React.useEffect(() => {
    loadAttachments();
  }, [paperId]);

  const loadAttachments = async () => {
    try {
      const response = await paperService.getPaper(paperId);
      setAttachments(response.data.data.attachments || []);
    } catch (error) {
      console.error('加载附件列表失败:', error);
    }
  };

  // 文件类型选择
  const getFileType = (file: UploadFile): string => {
    return file.type || '全文';
  };

  // 上传前校验
  const beforeUpload = (file: File): boolean => {
    const isValidType = ['application/pdf', 'image/jpeg', 'image/png'].includes(file.type);
    if (!isValidType) {
      message.error('只能上传 PDF、JPG、PNG 格式的文件');
      return false;
    }
    const isLt100M = file.size / 1024 / 1024 < 100;
    if (!isLt100M) {
      message.error('文件大小不能超过 100MB');
      return false;
    }
    return true;
  };

  // 自定义上传
  const customRequest: UploadProps['customRequest'] = async options => {
    const { file, onProgress, onSuccess, onError } = options;
    setUploading(true);

    try {
      // 模拟上传进度
      const uploadFile = file as File;
      const fileType = getFileType({ name: uploadFile.name } as UploadFile);

      const formData = new FormData();
      formData.append('file', uploadFile);
      formData.append('paper_id', paperId.toString());
      formData.append('file_type', fileType);

      const xhr = new XMLHttpRequest();

      xhr.upload.onprogress = (e: ProgressEvent) => {
        if (e.lengthComputable) {
          const percent = Math.round((e.loaded / e.total) * 100);
          onProgress?.({ percent }, file as UploadFile);
        }
      };

      xhr.onload = () => {
        if (xhr.status === 200) {
          message.success('文件上传成功');
          onSuccess?.(xhr.response, file as UploadFile);
          // 重新加载附件列表
          loadAttachments();
        } else {
          message.error('文件上传失败');
          onError?.(new Error('上传失败'), file as UploadFile);
        }
        setUploading(false);
      };

      xhr.onerror = () => {
        message.error('文件上传失败');
        onError?.(new Error('上传失败'), file as UploadFile);
        setUploading(false);
      };

      xhr.open('POST', '/api/files/upload');
      xhr.setRequestHeader('Authorization', `Bearer ${localStorage.getItem('biolit_token')}`);
      xhr.send(formData);
    } catch (error) {
      message.error('文件上传失败');
      onError?.(error as Error, file as UploadFile);
      setUploading(false);
    }
  };

  // 删除附件
  const handleDelete = async (id: number) => {
    try {
      await fileService.deleteFile(id);
      message.success('删除成功');
      loadAttachments();
    } catch (error) {
      message.error('删除失败');
    }
  };

  // 预览文件
  const handlePreview = (attachment: any) => {
    setPreviewFile(attachment);
    setPreviewVisible(true);
  };

  // 下载文件
  const handleDownload = async (attachment: any) => {
    try {
      const response = await fileService.downloadFile(attachment.id);
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', attachment.file_name);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (error) {
      message.error('下载失败');
    }
  };

  // 格式化文件大小
  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  return (
    <div>
      {mode !== 'view' && (
        <Upload
          fileList={fileList}
          beforeUpload={beforeUpload}
          customRequest={customRequest}
          onRemove={file => {
            setFileList(fileList.filter(f => f.uid !== file.uid));
          }}
          disabled={uploading}
        >
          <Button icon={<UploadOutlined />} loading={uploading} disabled={uploading}>
            上传文件
          </Button>
        </Upload>
      )}

      <div style={{ marginTop: 16 }}>
        <h4>已上传文件</h4>
        <List
          dataSource={attachments}
          renderItem={item => (
            <List.Item
              key={item.id}
              actions={
                mode !== 'view'
                  ? [
                      <Button
                        type="link"
                        icon={<EyeOutlined />}
                        onClick={() => handlePreview(item)}
                      >
                        预览
                      </Button>,
                      <Button
                        type="link"
                        icon={<DownloadOutlined />}
                        onClick={() => handleDownload(item)}
                      >
                        下载
                      </Button>,
                      <Button
                        type="link"
                        danger
                        icon={<DeleteOutlined />}
                        onClick={() => handleDelete(item.id)}
                      >
                        删除
                      </Button>,
                    ]
                  : [
                      <Button
                        type="link"
                        icon={<DownloadOutlined />}
                        onClick={() => handleDownload(item)}
                      >
                        下载
                      </Button>,
                    ]
              }
            >
              <List.Item.Meta
                title={
                  <Space>
                    <span>{item.file_name}</span>
                    <Tag color="blue">{item.file_type}</Tag>
                  </Space>
                }
                description={
                  <Space>
                    <span>大小: {formatFileSize(item.file_size)}</span>
                    <span>|</span>
                    <span>上传时间: {new Date(item.created_at).toLocaleString('zh-CN')}</span>
                  </Space>
                }
              />
            </List.Item>
          )}
        />
      </div>

      {/* 预览模态框 */}
      <Modal
        title="文件预览"
        open={previewVisible}
        onCancel={() => setPreviewVisible(false)}
        footer={null}
        width={800}
      >
        {previewFile && previewFile.mime_type.startsWith('image') ? (
          <img alt={previewFile.file_name} style={{ width: '100%' }} src={`/api/files/${previewFile.id}/preview`} />
        ) : (
          <div style={{ textAlign: 'center', padding: 40 }}>
            该文件类型不支持预览,请下载后查看
          </div>
        )}
      </Modal>
    </div>
  );
};

export default FileUpload;
