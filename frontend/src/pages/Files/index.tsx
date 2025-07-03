/**
 * LYSS AI Platform 文件管理页面
 * 提供文件上传、管理和文档问答功能
 */

import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Upload,
  Button,
  Table,
  Space,
  Tag,
  Popconfirm,
  message,
  Modal,
  Input,
  Row,
  Col,
  Statistic,
  Progress,
  Tooltip,
  Badge,
  Alert,
  Divider,
} from 'antd';
import {
  UploadOutlined,
  FileOutlined,
  DeleteOutlined,
  DownloadOutlined,
  EyeOutlined,
  MessageOutlined,
  SearchOutlined,
  InboxOutlined,
  FilePdfOutlined,
  FileTextOutlined,
  FileWordOutlined,
  FileExcelOutlined,
  FileImageOutlined,
  CloudUploadOutlined,
  FolderOutlined,
} from '@ant-design/icons';
import { UploadProps } from 'antd/lib/upload';
import { useAuthStore } from '@/stores/auth';
import { api } from '@/utils/request';
import { API_ENDPOINTS } from '@/config/api';
import { UploadedFile, FileUploadRequest } from '@/types/api';
import { createStyles } from 'antd-style';
import { ColumnsType } from 'antd/es/table';

const { Title, Text } = Typography;
const { Dragger } = Upload;
const { Search } = Input;

const useStyles = createStyles(({ token }) => ({
  container: {
    padding: '24px',
  },
  uploadCard: {
    marginBottom: '24px',
  },
  statsCard: {
    marginBottom: '24px',
  },
  fileList: {
    '& .ant-table-tbody > tr > td': {
      padding: '12px 16px',
    },
  },
  fileIcon: {
    fontSize: '24px',
    marginRight: '8px',
  },
  fileName: {
    fontWeight: 500,
  },
  fileSize: {
    color: token.colorTextSecondary,
    fontSize: '12px',
  },
  uploadArea: {
    background: token.colorBgElevated,
    border: `2px dashed ${token.colorBorder}`,
    borderRadius: token.borderRadiusLG,
    padding: '40px',
    textAlign: 'center',
    '&.ant-upload-drag-hover': {
      border: `2px dashed ${token.colorPrimary}`,
      background: token.colorPrimaryBg,
    },
  },
  searchBar: {
    marginBottom: '16px',
  },
  actionButton: {
    marginLeft: '8px',
  },
}));

const Files: React.FC = () => {
  const { styles } = useStyles();
  const { user } = useAuthStore();
  const [files, setFiles] = useState<UploadedFile[]>([]);
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [searchText, setSearchText] = useState('');
  const [selectedFile, setSelectedFile] = useState<UploadedFile | null>(null);
  const [previewModal, setPreviewModal] = useState(false);
  const [chatModal, setChatModal] = useState(false);
  const [stats, setStats] = useState({
    totalFiles: 0,
    totalSize: 0,
    supportedTypes: ['PDF', 'TXT', 'DOCX', 'XLSX', 'PNG', 'JPG', 'JPEG'],
  });

  // 获取文件列表
  const fetchFiles = async () => {
    setLoading(true);
    try {
      const response = await api.get<UploadedFile[]>(API_ENDPOINTS.FILES.LIST);
      setFiles(response);
      setStats(prev => ({
        ...prev,
        totalFiles: response.length,
        totalSize: response.reduce((sum, file) => sum + file.size, 0),
      }));
    } catch (error) {
      console.error('获取文件列表失败:', error);
      message.error('获取文件列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchFiles();
  }, []);

  // 文件上传配置
  const uploadProps: UploadProps = {
    name: 'file',
    multiple: true,
    maxCount: 10,
    accept: '.pdf,.txt,.docx,.xlsx,.png,.jpg,.jpeg',
    showUploadList: false,
    beforeUpload: (file) => {
      // 检查文件大小（限制为50MB）
      const isLt50M = file.size / 1024 / 1024 < 50;
      if (!isLt50M) {
        message.error('文件大小不能超过50MB');
        return false;
      }
      return true;
    },
    customRequest: async ({ file, onSuccess, onError }) => {
      try {
        setUploading(true);
        const formData = new FormData();
        formData.append('file', file as File);
        
        const response = await api.post(API_ENDPOINTS.FILES.UPLOAD, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });
        
        message.success(`${(file as File).name} 上传成功`);
        onSuccess?.(response);
        fetchFiles(); // 刷新文件列表
      } catch (error: any) {
        console.error('文件上传失败:', error);
        message.error(`${(file as File).name} 上传失败`);
        onError?.(error);
      } finally {
        setUploading(false);
      }
    },
  };

  // 删除文件
  const handleDeleteFile = async (fileId: string) => {
    try {
      await api.delete(`${API_ENDPOINTS.FILES.DELETE}/${fileId}`);
      message.success('文件删除成功');
      fetchFiles();
    } catch (error) {
      console.error('删除文件失败:', error);
      message.error('删除文件失败');
    }
  };

  // 下载文件
  const handleDownloadFile = async (file: UploadedFile) => {
    try {
      const response = await api.get(`${API_ENDPOINTS.FILES.DOWNLOAD}/${file.id}`, {
        responseType: 'blob',
      });
      
      const url = window.URL.createObjectURL(new Blob([response]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', file.filename);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('下载文件失败:', error);
      message.error('下载文件失败');
    }
  };

  // 文件预览
  const handlePreviewFile = (file: UploadedFile) => {
    setSelectedFile(file);
    setPreviewModal(true);
  };

  // 文档问答
  const handleChatWithFile = (file: UploadedFile) => {
    setSelectedFile(file);
    setChatModal(true);
  };

  // 获取文件图标
  const getFileIcon = (filename: string) => {
    const ext = filename.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf':
        return <FilePdfOutlined className={styles.fileIcon} style={{ color: '#ff4d4f' }} />;
      case 'txt':
        return <FileTextOutlined className={styles.fileIcon} style={{ color: '#1890ff' }} />;
      case 'docx':
      case 'doc':
        return <FileWordOutlined className={styles.fileIcon} style={{ color: '#1890ff' }} />;
      case 'xlsx':
      case 'xls':
        return <FileExcelOutlined className={styles.fileIcon} style={{ color: '#52c41a' }} />;
      case 'png':
      case 'jpg':
      case 'jpeg':
        return <FileImageOutlined className={styles.fileIcon} style={{ color: '#722ed1' }} />;
      default:
        return <FileOutlined className={styles.fileIcon} />;
    }
  };

  // 格式化文件大小
  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // 表格列定义
  const columns: ColumnsType<UploadedFile> = [
    {
      title: '文件名',
      dataIndex: 'filename',
      key: 'filename',
      render: (text: string, record: UploadedFile) => (
        <div style={{ display: 'flex', alignItems: 'center' }}>
          {getFileIcon(text)}
          <div>
            <div className={styles.fileName}>{text}</div>
            <div className={styles.fileSize}>{formatFileSize(record.size)}</div>
          </div>
        </div>
      ),
      sorter: (a, b) => a.filename.localeCompare(b.filename),
    },
    {
      title: '类型',
      dataIndex: 'content_type',
      key: 'content_type',
      render: (text: string) => <Tag color="blue">{text}</Tag>,
      width: 150,
    },
    {
      title: '上传时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => new Date(text).toLocaleString(),
      sorter: (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
      width: 180,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'processed' ? 'success' : 'processing'}>
          {status === 'processed' ? '已处理' : '处理中'}
        </Tag>
      ),
      width: 100,
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record: UploadedFile) => (
        <Space size="small">
          <Tooltip title="预览">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => handlePreviewFile(record)}
            />
          </Tooltip>
          <Tooltip title="文档问答">
            <Button
              type="text"
              icon={<MessageOutlined />}
              onClick={() => handleChatWithFile(record)}
              disabled={record.status !== 'processed'}
            />
          </Tooltip>
          <Tooltip title="下载">
            <Button
              type="text"
              icon={<DownloadOutlined />}
              onClick={() => handleDownloadFile(record)}
            />
          </Tooltip>
          <Popconfirm
            title="确定删除此文件吗？"
            onConfirm={() => handleDeleteFile(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button type="text" danger icon={<DeleteOutlined />} />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
      width: 150,
    },
  ];

  // 过滤文件列表
  const filteredFiles = files.filter(file =>
    file.filename.toLowerCase().includes(searchText.toLowerCase())
  );

  return (
    <div className={styles.container}>
      {/* 页面标题 */}
      <Title level={2} style={{ marginBottom: '16px' }}>
        <FolderOutlined /> 文件管理
      </Title>
      <Text type="secondary">
        上传、管理文件，并与文档进行智能问答
      </Text>

      {/* 统计信息 */}
      <Card className={styles.statsCard}>
        <Row gutter={16}>
          <Col span={6}>
            <Statistic title="总文件数" value={stats.totalFiles} />
          </Col>
          <Col span={6}>
            <Statistic title="总大小" value={formatFileSize(stats.totalSize)} />
          </Col>
          <Col span={6}>
            <Statistic title="支持格式" value={stats.supportedTypes.length} />
          </Col>
          <Col span={6}>
            <div>
              <Text strong>已处理文件</Text>
              <div style={{ marginTop: '8px' }}>
                <Progress 
                  percent={Math.round((files.filter(f => f.status === 'processed').length / Math.max(files.length, 1)) * 100)}
                  size="small"
                />
              </div>
            </div>
          </Col>
        </Row>
      </Card>

      {/* 文件上传区域 */}
      <Card title="文件上传" className={styles.uploadCard}>
        <Alert
          message="支持的文件格式"
          description={`PDF, TXT, DOCX, XLSX, PNG, JPG, JPEG（最大50MB）`}
          type="info"
          showIcon
          style={{ marginBottom: '16px' }}
        />
        
        <Dragger {...uploadProps} className={styles.uploadArea}>
          <p className="ant-upload-drag-icon">
            <CloudUploadOutlined style={{ fontSize: '48px', color: '#1890ff' }} />
          </p>
          <p className="ant-upload-text">
            点击或拖拽文件到此区域上传
          </p>
          <p className="ant-upload-hint">
            支持单个或批量上传，严格限制文件大小不超过50MB
          </p>
        </Dragger>

        {uploading && (
          <div style={{ marginTop: '16px', textAlign: 'center' }}>
            <Progress 
              percent={100} 
              status="active" 
              format={() => '上传中...'}
            />
          </div>
        )}
      </Card>

      {/* 文件列表 */}
      <Card title="文件列表">
        <div className={styles.searchBar}>
          <Search
            placeholder="搜索文件名..."
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 300 }}
            prefix={<SearchOutlined />}
          />
        </div>

        <Table
          columns={columns}
          dataSource={filteredFiles}
          rowKey="id"
          loading={loading}
          pagination={{
            total: filteredFiles.length,
            pageSize: 10,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
          }}
          className={styles.fileList}
        />
      </Card>

      {/* 文件预览模态框 */}
      <Modal
        title="文件预览"
        open={previewModal}
        onCancel={() => setPreviewModal(false)}
        footer={null}
        width={800}
        style={{ top: 20 }}
      >
        {selectedFile && (
          <div style={{ padding: '16px' }}>
            <div style={{ marginBottom: '16px' }}>
              {getFileIcon(selectedFile.filename)}
              <Text strong>{selectedFile.filename}</Text>
              <Text type="secondary" style={{ marginLeft: '16px' }}>
                {formatFileSize(selectedFile.size)}
              </Text>
            </div>
            <Divider />
            <div style={{ textAlign: 'center', padding: '40px' }}>
              <Text type="secondary">
                文件预览功能开发中...
              </Text>
            </div>
          </div>
        )}
      </Modal>

      {/* 文档问答模态框 */}
      <Modal
        title="文档问答"
        open={chatModal}
        onCancel={() => setChatModal(false)}
        footer={null}
        width={800}
        style={{ top: 20 }}
      >
        {selectedFile && (
          <div style={{ padding: '16px' }}>
            <div style={{ marginBottom: '16px' }}>
              <Badge status="success" text="已处理" />
              <Text strong style={{ marginLeft: '8px' }}>
                {selectedFile.filename}
              </Text>
            </div>
            <Divider />
            <div style={{ textAlign: 'center', padding: '40px' }}>
              <Text type="secondary">
                文档问答功能开发中...
              </Text>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default Files;