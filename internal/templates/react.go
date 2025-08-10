package templates

// ReactListTemplate generates a React list component with CRUD operations
const ReactListTemplate = `import React, { useState, useEffect } from 'react';
import { 
  Table, 
  Button, 
  Modal, 
  Form, 
  Input, 
  Space, 
  Popconfirm, 
  message,
  Card,
  Row,
  Col,
  Typography
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { {{.Names.PascalCase}}, {{.Names.PascalCase}}Filter } from '../../types/{{.Names.CamelCase}}';
import { use{{.Names.PascalCase}} } from '../../hooks/use{{.Names.PascalCase}}';
import {{.Names.PascalCase}}Form from './{{.Names.PascalCase}}Form';
import {{.Names.PascalCase}}Detail from './{{.Names.PascalCase}}Detail';

const { Title } = Typography;

interface {{.Names.PascalCase}}ListProps {
  title?: string;
  showActions?: boolean;
  selectable?: boolean;
  onSelect?: ({{.Names.CamelCase}}: {{.Names.PascalCase}}) => void;
}

const {{.Names.PascalCase}}List: React.FC<{{.Names.PascalCase}}ListProps> = ({
  title = "{{.DisplayName}} Management",
  showActions = true,
  selectable = false,
  onSelect
}) => {
  const [filter, setFilter] = useState<{{.Names.PascalCase}}Filter>({
    page: 1,
    page_size: 10,
    search: '',
    sort: 'created_at',
    order: 'DESC'
  });
  
  const [isFormModalVisible, setIsFormModalVisible] = useState(false);
  const [isDetailModalVisible, setIsDetailModalVisible] = useState(false);
  const [selectedRecord, setSelectedRecord] = useState<{{.Names.PascalCase}} | null>(null);
  const [formMode, setFormMode] = useState<'create' | 'edit'>('create');

  const {
    data: {{.Names.CamelPlural}},
    total,
    loading,
    error,
    fetchAll,
    create,
    update,
    remove
  } = use{{.Names.PascalCase}}();

  useEffect(() => {
    fetchAll(filter);
  }, [filter, fetchAll]);

  const handleCreate = () => {
    setSelectedRecord(null);
    setFormMode('create');
    setIsFormModalVisible(true);
  };

  const handleEdit = (record: {{.Names.PascalCase}}) => {
    setSelectedRecord(record);
    setFormMode('edit');
    setIsFormModalVisible(true);
  };

  const handleView = (record: {{.Names.PascalCase}}) => {
    setSelectedRecord(record);
    setIsDetailModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await remove(id);
      message.success('{{.DisplayName}} deleted successfully');
      fetchAll(filter);
    } catch (error) {
      message.error('Failed to delete {{.Names.Singular}}');
    }
  };

  const handleFormSubmit = async (values: any) => {
    try {
      if (formMode === 'create') {
        await create(values);
        message.success('{{.DisplayName}} created successfully');
      } else if (selectedRecord) {
        await update(selectedRecord.id, values);
        message.success('{{.DisplayName}} updated successfully');
      }
      setIsFormModalVisible(false);
      fetchAll(filter);
    } catch (error) {
      message.error(` + "`Failed to ${formMode} {{.Names.Singular}}`" + `);
    }
  };

  const handleSearch = (value: string) => {
    setFilter(prev => ({ ...prev, search: value, page: 1 }));
  };

  const handleTableChange = (pagination: any, filters: any, sorter: any) => {
    setFilter(prev => ({
      ...prev,
      page: pagination.current,
      page_size: pagination.pageSize,
      sort: sorter.field || 'created_at',
      order: sorter.order === 'ascend' ? 'ASC' : 'DESC'
    }));
  };

  const columns = [
{{- range .Fields}}
{{- if not .ReadOnly}}
    {
      title: '{{.DisplayName}}',
      dataIndex: '{{.Names.CamelCase}}',
      key: '{{.Names.CamelCase}}',
      sorter: true,
{{- if eq .Type "string"}}
      render: (text: string) => text || '-',
{{- else if eq .Type "boolean"}}
      render: (value: boolean) => value ? 'Yes' : 'No',
{{- else if eq .Type "date"}}
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
{{- else}}
      render: (value: any) => value?.toString() || '-',
{{- end}}
    },
{{- end}}
{{- end}}
    {
      title: 'Created',
      dataIndex: 'created_at',
      key: 'created_at',
      sorter: true,
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    ...(showActions ? [{
      title: 'Actions',
      key: 'actions',
      width: 200,
      render: (_: any, record: {{.Names.PascalCase}}) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleView(record)}
          >
            View
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            Edit
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this {{.Names.Singular}}?"
            onConfirm={() => handleDelete(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
            >
              Delete
            </Button>
          </Popconfirm>
          {selectable && (
            <Button
              type="primary"
              size="small"
              onClick={() => onSelect?.(record)}
            >
              Select
            </Button>
          )}
        </Space>
      ),
    }] : [])
  ];

  if (error) {
    return (
      <Card>
        <Typography.Text type="danger">
          Error loading {{.Names.Plural}}: {error.message}
        </Typography.Text>
      </Card>
    );
  }

  return (
    <div className="{{.Names.KebabCase}}-list">
      <Card>
        <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
          <Col>
            <Title level={3}>{title}</Title>
          </Col>
          <Col>
            <Space>
              <Input.Search
                placeholder="Search {{.Names.Plural}}..."
                allowClear
                onSearch={handleSearch}
                style={{ width: 300 }}
              />
              {showActions && (
                <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={handleCreate}
                >
                  Add {{.DisplayName}}
                </Button>
              )}
            </Space>
          </Col>
        </Row>

        <Table
          columns={columns}
          dataSource={{{.Names.CamelPlural}}}
          loading={loading}
          rowKey="id"
          pagination={{
            current: filter.page,
            pageSize: filter.page_size,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => 
              ` + "`${range[0]}-${range[1]} of ${total} {{.Names.Plural}}`" + `,
          }}
          onChange={handleTableChange}
          scroll={{ x: 'max-content' }}
        />
      </Card>

      <Modal
        title={` + "`${formMode === 'create' ? 'Add' : 'Edit'} {{.DisplayName}}`" + `}
        open={isFormModalVisible}
        onCancel={() => setIsFormModalVisible(false)}
        footer={null}
        width={600}
        destroyOnClose
      >
        <{{.Names.PascalCase}}Form
          initialValues={selectedRecord}
          onSubmit={handleFormSubmit}
          onCancel={() => setIsFormModalVisible(false)}
          loading={loading}
        />
      </Modal>

      <Modal
        title="{{.DisplayName}} Details"
        open={isDetailModalVisible}
        onCancel={() => setIsDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setIsDetailModalVisible(false)}>
            Close
          </Button>,
        ]}
        width={800}
      >
        {selectedRecord && (
          <{{.Names.PascalCase}}Detail {{.Names.CamelCase}}={selectedRecord} />
        )}
      </Modal>
    </div>
  );
};

export default {{.Names.PascalCase}}List;`

// ReactFormTemplate generates a React form component
const ReactFormTemplate = `import React, { useEffect } from 'react';
import { Form, Input, Button, Space, Row, Col{{range .Fields}}{{if eq .Type "boolean"}}, Switch{{end}}{{if eq .Type "number"}}, InputNumber{{end}}{{if eq .Type "date"}}, DatePicker{{end}}{{if .Options}}, Select{{end}}{{end}} } from 'antd';
import { {{.Names.PascalCase}} } from '../../types/{{.Names.CamelCase}}';
{{range .Fields}}{{if eq .Type "date"}}import dayjs from 'dayjs';{{break}}{{end}}{{end}}

const { TextArea } = Input;
{{range .Fields}}{{if .Options}}const { Option } = Select;{{break}}{{end}}{{end}}

interface {{.Names.PascalCase}}FormProps {
  initialValues?: Partial<{{.Names.PascalCase}}>;
  onSubmit: (values: any) => void;
  onCancel: () => void;
  loading?: boolean;
}

const {{.Names.PascalCase}}Form: React.FC<{{.Names.PascalCase}}FormProps> = ({
  initialValues,
  onSubmit,
  onCancel,
  loading = false
}) => {
  const [form] = Form.useForm();

  useEffect(() => {
    if (initialValues) {
      const formValues = { ...initialValues };
{{range .Fields}}
{{- if eq .Type "date"}}
      if (formValues.{{.Names.CamelCase}}) {
        formValues.{{.Names.CamelCase}} = dayjs(formValues.{{.Names.CamelCase}});
      }
{{- end}}
{{end}}
      form.setFieldsValue(formValues);
    }
  }, [initialValues, form]);

  const handleSubmit = (values: any) => {
{{range .Fields}}
{{- if eq .Type "date"}}
    if (values.{{.Names.CamelCase}}) {
      values.{{.Names.CamelCase}} = values.{{.Names.CamelCase}}.toISOString();
    }
{{- end}}
{{end}}
    onSubmit(values);
  };

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={handleSubmit}
      initialValues={initialValues}
    >
      <Row gutter={16}>
{{range .Fields}}
{{- if not .ReadOnly}}
        <Col span={{{if eq .Type "text"}}24{{else}}12{{end}}>
          <Form.Item
            label="{{.DisplayName}}"
            name="{{.Names.CamelCase}}"
            rules={[
{{- if .Required}}
              { required: true, message: 'Please enter {{.DisplayName}}' },
{{- end}}
{{- if eq .Type "email"}}
              { type: 'email', message: 'Please enter a valid email' },
{{- end}}
{{- if .Pattern}}
              { pattern: new RegExp('{{.Pattern}}'), message: '{{.DisplayName}} format is invalid' },
{{- end}}
{{- if .MinLength}}
              { min: {{.MinLength}}, message: '{{.DisplayName}} must be at least {{.MinLength}} characters' },
{{- end}}
{{- if .MaxLength}}
              { max: {{.MaxLength}}, message: '{{.DisplayName}} must be at most {{.MaxLength}} characters' },
{{- end}}
            ]}
          >
{{- if eq .Type "boolean"}}
            <Switch checkedChildren="Yes" unCheckedChildren="No" />
{{- else if eq .Type "number"}}
            <InputNumber 
              style={{ width: '100%' }}
              placeholder="Enter {{.DisplayName}}"
{{- if .Min}}
              min={{{.Min}}}
{{- end}}
{{- if .Max}}
              max={{{.Max}}}
{{- end}}
            />
{{- else if eq .Type "date"}}
            <DatePicker 
              style={{ width: '100%' }}
              placeholder="Select {{.DisplayName}}"
            />
{{- else if .Options}}
            <Select placeholder="Select {{.DisplayName}}">
{{- range .Options}}
              <Option value="{{.}}">{{.}}</Option>
{{- end}}
            </Select>
{{- else if eq .Type "text"}}
            <TextArea 
              rows={4}
              placeholder="Enter {{.DisplayName}}"
{{- if .MaxLength}}
              maxLength={{{.MaxLength}}}
              showCount
{{- end}}
            />
{{- else}}
            <Input 
              placeholder="Enter {{.DisplayName}}"
{{- if .MaxLength}}
              maxLength={{{.MaxLength}}}
{{- end}}
            />
{{- end}}
          </Form.Item>
        </Col>
{{- end}}
{{end}}
      </Row>

      <Form.Item>
        <Space>
          <Button type="primary" htmlType="submit" loading={loading}>
            {initialValues ? 'Update' : 'Create'}
          </Button>
          <Button onClick={onCancel}>
            Cancel
          </Button>
        </Space>
      </Form.Item>
    </Form>
  );
};

export default {{.Names.PascalCase}}Form;`

// ReactDetailTemplate generates a React detail view component
const ReactDetailTemplate = `import React from 'react';
import { Descriptions, Card, Tag, Typography } from 'antd';
import { {{.Names.PascalCase}} } from '../../types/{{.Names.CamelCase}}';

const { Title } = Typography;

interface {{.Names.PascalCase}}DetailProps {
  {{.Names.CamelCase}}: {{.Names.PascalCase}};
}

const {{.Names.PascalCase}}Detail: React.FC<{{.Names.PascalCase}}DetailProps> = ({ {{.Names.CamelCase}} }) => {
  return (
    <div className="{{.Names.KebabCase}}-detail">
      <Card>
        <Title level={4}>{{.DisplayName}} Information</Title>
        <Descriptions column={2} bordered>
          <Descriptions.Item label="ID">
            {{{.Names.CamelCase}}.id}
          </Descriptions.Item>
{{range .Fields}}
          <Descriptions.Item label="{{.DisplayName}}">
{{- if eq .Type "boolean"}}
            <Tag color={{{$.Names.CamelCase}}.{{.Names.CamelCase}} ? 'green' : 'red'}>
              {{{$.Names.CamelCase}}.{{.Names.CamelCase}} ? 'Yes' : 'No'}
            </Tag>
{{- else if eq .Type "date"}}
            {{{$.Names.CamelCase}}.{{.Names.CamelCase}} ? new Date({{$.Names.CamelCase}}.{{.Names.CamelCase}}).toLocaleString() : '-'}
{{- else if eq .Type "email"}}
            <a href={` + "`mailto:${{{$.Names.CamelCase}}.{{.Names.CamelCase}}}`" + `}>
              {{{$.Names.CamelCase}}.{{.Names.CamelCase}} || '-'}
            </a>
{{- else if eq .Type "url"}}
            <a href={{{$.Names.CamelCase}}.{{.Names.CamelCase}}} target="_blank" rel="noopener noreferrer">
              {{{$.Names.CamelCase}}.{{.Names.CamelCase}} || '-'}
            </a>
{{- else}}
            {{{$.Names.CamelCase}}.{{.Names.CamelCase}}?.toString() || '-'}
{{- end}}
          </Descriptions.Item>
{{end}}
          <Descriptions.Item label="Created At">
            {new Date({{.Names.CamelCase}}.created_at).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Updated At">
            {new Date({{.Names.CamelCase}}.updated_at).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
};

export default {{.Names.PascalCase}}Detail;`

// ReactTypesTemplate generates TypeScript type definitions
const ReactTypesTemplate = `// {{.DisplayName}} TypeScript definitions
// Generated by ViberCode CLI

export interface {{.Names.PascalCase}} {
  id: string;
  created_at: string;
  updated_at: string;
{{range .Fields}}
  {{.Names.CamelCase}}{{if not .Required}}?{{end}}: {{.TypeScriptType}};
{{end}}
}

export interface {{.Names.PascalCase}}Request {
{{range .Fields}}
{{- if not .ReadOnly}}
  {{.Names.CamelCase}}{{if not .Required}}?{{end}}: {{.TypeScriptType}};
{{- end}}
{{end}}
}

export interface {{.Names.PascalCase}}Response {
  id: string;
  created_at: string;
  updated_at: string;
{{range .Fields}}
  {{.Names.CamelCase}}{{if not .Required}}?{{end}}: {{.TypeScriptType}};
{{end}}
}

export interface {{.Names.PascalCase}}Filter {
  page?: number;
  page_size?: number;
  sort?: string;
  order?: 'ASC' | 'DESC';
  search?: string;
{{range .Fields}}
{{- if .Filterable}}
  {{.Names.CamelCase}}?: {{.TypeScriptType}};
{{- end}}
{{end}}
}

export interface {{.Names.PascalCase}}ListResponse {
  data: {{.Names.PascalCase}}Response[];
  total: number;
  page: number;
  page_size: number;
}

// Form validation types
export interface {{.Names.PascalCase}}FormData extends Omit<{{.Names.PascalCase}}Request, 'id' | 'created_at' | 'updated_at'> {}

// API error types
export interface {{.Names.PascalCase}}APIError {
  message: string;
  code?: string;
  field?: string;
}

// Hook state types
export interface Use{{.Names.PascalCase}}State {
  data: {{.Names.PascalCase}}[];
  total: number;
  loading: boolean;
  error: {{.Names.PascalCase}}APIError | null;
}

export interface Use{{.Names.PascalCase}}Actions {
  fetchAll: (filter?: {{.Names.PascalCase}}Filter) => Promise<void>;
  fetchOne: (id: string) => Promise<{{.Names.PascalCase}} | null>;
  create: (data: {{.Names.PascalCase}}Request) => Promise<{{.Names.PascalCase}}>;
  update: (id: string, data: {{.Names.PascalCase}}Request) => Promise<{{.Names.PascalCase}}>;
  remove: (id: string) => Promise<void>;
  clearError: () => void;
}`

// ReactHooksTemplate generates React hooks for API operations
const ReactHooksTemplate = `import { useState, useCallback } from 'react';
import { 
  {{.Names.PascalCase}}, 
  {{.Names.PascalCase}}Request, 
  {{.Names.PascalCase}}Filter, 
  {{.Names.PascalCase}}ListResponse,
  {{.Names.PascalCase}}APIError,
  Use{{.Names.PascalCase}}State,
  Use{{.Names.PascalCase}}Actions
} from '../types/{{.Names.CamelCase}}';

const API_BASE_URL = process.env.REACT_APP_API_URL || '{{.ApiBaseUrl}}';

export const use{{.Names.PascalCase}} = (): Use{{.Names.PascalCase}}State & Use{{.Names.PascalCase}}Actions => {
  const [data, setData] = useState<{{.Names.PascalCase}}[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<{{.Names.PascalCase}}APIError | null>(null);

  const handleError = useCallback((err: any) => {
    const errorMessage = err.response?.data?.message || err.message || 'An error occurred';
    setError({ message: errorMessage, code: err.response?.status?.toString() });
    console.error('{{.Names.PascalCase}} API Error:', err);
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const fetchAll = useCallback(async (filter: {{.Names.PascalCase}}Filter = {}) => {
    setLoading(true);
    setError(null);
    
    try {
      const queryParams = new URLSearchParams();
      Object.entries(filter).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          queryParams.append(key, value.toString());
        }
      });

      const response = await fetch(` + "`${API_BASE_URL}/{{.Names.KebabPlural}}?${queryParams}`" + `, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(` + "`HTTP error! status: ${response.status}`" + `);
      }

      const result: {{.Names.PascalCase}}ListResponse = await response.json();
      setData(result.data || []);
      setTotal(result.total || 0);
    } catch (err) {
      handleError(err);
      setData([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [handleError]);

  const fetchOne = useCallback(async (id: string): Promise<{{.Names.PascalCase}} | null> => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(` + "`${API_BASE_URL}/{{.Names.KebabPlural}}/${id}`" + `, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(` + "`HTTP error! status: ${response.status}`" + `);
      }

      const {{.Names.CamelCase}}: {{.Names.PascalCase}} = await response.json();
      return {{.Names.CamelCase}};
    } catch (err) {
      handleError(err);
      return null;
    } finally {
      setLoading(false);
    }
  }, [handleError]);

  const create = useCallback(async (requestData: {{.Names.PascalCase}}Request): Promise<{{.Names.PascalCase}}> => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(` + "`${API_BASE_URL}/{{.Names.KebabPlural}}`" + `, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error(` + "`HTTP error! status: ${response.status}`" + `);
      }

      const new{{.Names.PascalCase}}: {{.Names.PascalCase}} = await response.json();
      setData(prev => [new{{.Names.PascalCase}}, ...prev]);
      setTotal(prev => prev + 1);
      return new{{.Names.PascalCase}};
    } catch (err) {
      handleError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [handleError]);

  const update = useCallback(async (id: string, requestData: {{.Names.PascalCase}}Request): Promise<{{.Names.PascalCase}}> => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(` + "`${API_BASE_URL}/{{.Names.KebabPlural}}/${id}`" + `, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error(` + "`HTTP error! status: ${response.status}`" + `);
      }

      const updated{{.Names.PascalCase}}: {{.Names.PascalCase}} = await response.json();
      setData(prev => prev.map(item => item.id === id ? updated{{.Names.PascalCase}} : item));
      return updated{{.Names.PascalCase}};
    } catch (err) {
      handleError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [handleError]);

  const remove = useCallback(async (id: string): Promise<void> => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(` + "`${API_BASE_URL}/{{.Names.KebabPlural}}/${id}`" + `, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(` + "`HTTP error! status: ${response.status}`" + `);
      }

      setData(prev => prev.filter(item => item.id !== id));
      setTotal(prev => prev - 1);
    } catch (err) {
      handleError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [handleError]);

  return {
    // State
    data,
    total,
    loading,
    error,
    
    // Actions
    fetchAll,
    fetchOne,
    create,
    update,
    remove,
    clearError,
  };
};`