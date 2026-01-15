/**
 * API Service - Real backend API calls
 * Base URL is configured via Vite proxy to avoid CORS issues
 * 
 * Backend Response Formats:
 * 1. Using response.Success(): { code: 200, message: "...", data: {...} }
 * 2. Using direct gin.H: { success: true, data: {...} } or { data: {...}, pagination: {...} }
 */

import {
  UserProfile,
  NotificationLog,
  SendMessageRequest,
  Project,
  SystemHealth,
  BillingRecord,
  Company,
  BotInstance,
  TransactionRecord,
  ADGroup,
  TeamChannel,
  ChatGroup,
  Template,
  MessageStatus,
  MessagePriority,
  SyncStatus
} from '../types';

// API Base URLs - uses Vite proxy in development
// External API (public): /api/v1 - for notify, provision, messages
// Internal API: /internal/v1 - for admin features (projects, bots, billing, etc.)
const API_EXTERNAL = '/api/v1';
const API_INTERNAL = '/internal/v1';

// ============================================
// Backend Response Types (actual format from Go)
// ============================================

// Backend Project model (from libs/models/models.go)
interface BackendProject {
  id: string;
  company_id: string;
  notify_key: string;
  description: string;
  status: string;
  daily_limit: number;
  monthly_limit: number;
  priority: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

// Backend Company model
interface BackendCompany {
  id: string;
  name: string;
  contact_email: string;
  contact_phone: string;
  address: string;
  status: string;
  billing_enabled: boolean;
  created_at: string;
  updated_at: string;
}

// Backend Notification model
interface BackendNotification {
  id: string;
  project_id: string;
  sender_id?: string;
  message_type: string;
  content: string;
  mentions?: string[];
  attachments?: any[];
  priority: string;
  status: string;  // pending, processing, enqueued, sent, failed, cancelled
  metadata?: Record<string, any>;
  targets?: string[];
  error_message?: string;
  sent_at?: string;
  created_at: string;
  updated_at: string;
}

// Backend TeamsBot model
interface BackendTeamsBot {
  id: string;
  type: string;
  company_id?: string;
  name: string;
  description?: string;
  app_id: string;
  tenant_id?: string;
  status?: string;  // active, inactive, suspended, maintenance
  webhook_url?: string;
  capabilities?: Record<string, any>;
  rate_limit_per_minute?: number;
  max_concurrent_requests?: number;
  created_at: string;
  updated_at: string;
}

// Backend UsageRecord
interface BackendUsageRecord {
  id: string;
  company_id: string;
  project_id: string;
  user_id?: string;
  record_type: string;  // notification, attachment, mention, adaptive_card
  quantity: number;
  unit_price: number;
  total_amount: number;
  billing_period: string;
  reference_id?: string;
  metadata?: Record<string, any>;
  created_at: string;
}

// Backend SystemHealth
interface BackendSystemHealth {
  status: string;
  timestamp: string;
  components: Record<string, {
    status: string;
    message?: string;
    timestamp: string;
    latency_ms?: number;
  }>;
  overall: {
    score: number;
    grade: string;
    message: string;
  };
}

// Backend Template model (matches models.Template)
interface BackendTemplate {
  id: string;
  project_id: string;
  name: string;
  description: string;
  variables?: Array<{
    key: string;
    label: string;
    type: string;
    options?: string[];
  }>;  // JSONBTemplateVariables - array of TemplateVariable objects
  default_json_structure: string;  // The Adaptive Card JSON structure
  created_at: string;
  updated_at: string;
}

// ============================================
// Generic Response Wrappers
// ============================================

// Format 1: response.Success() wrapper
interface ResponseSuccessWrapper<T> {
  code: number;
  message: string;
  data?: T;
  error?: string;
  details?: any;
}

// Format 2: Direct gin.H wrapper  
interface DirectResponseWrapper<T> {
  success?: boolean;
  data?: T;
  error?: string;
  pagination?: {
    total: number;
    limit: number;
    offset: number;
    pages: number;
  };
}

// Combined type that handles both formats
type BackendResponse<T> = ResponseSuccessWrapper<T> | DirectResponseWrapper<T>;

// Helper to extract data from either response format
function extractData<T>(response: BackendResponse<T>): T | undefined {
  // Handle response.Success() format: { code, message, data }
  if ('code' in response && 'data' in response) {
    return response.data;
  }
  // Handle direct gin.H format: { success, data } or { data, pagination }
  if ('data' in response) {
    return response.data;
  }
  // If response is already the data type (some endpoints return data directly)
  return response as T;
}

// Helper function for API calls
async function apiCall<T>(
  endpoint: string,
  options: RequestInit = {},
  useInternal: boolean = true
): Promise<BackendResponse<T>> {
  const base = useInternal ? API_INTERNAL : API_EXTERNAL;
  const url = `${base}${endpoint}`;

  const defaultHeaders: HeadersInit = {
    'Content-Type': 'application/json',
  };

  const response = await fetch(url, {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  });

  if (!response.ok) {
    // Try to extract error message from response
    let errorMessage = `API Error: ${response.status}`;
    try {
      const errorData = await response.json();
      // Handle different error response formats
      if (errorData.error) {
        errorMessage = errorData.error;
      } else if (errorData.message) {
        errorMessage = errorData.message;
      } else if (errorData.details) {
        errorMessage = typeof errorData.details === 'string'
          ? errorData.details
          : errorData.details.toString();
      } else if (typeof errorData === 'string') {
        errorMessage = errorData;
      }
    } catch {
      // If JSON parsing fails, use status text
      errorMessage = response.statusText || `HTTP ${response.status}`;
    }
    throw new Error(errorMessage);
  }

  return response.json();
}

// ============================================
// DTO Transformers (Backend → Frontend)
// ============================================

function transformProject(backend: BackendProject): Project {
  // Map backend priority to frontend priority format
  const priorityMap: Record<string, string> = {
    'low': 'Low (Broadcast)',
    'normal': 'Normal',
    'high': 'High (Alert)',
  };

  return {
    id: backend.id,
    name: backend.notify_key,  // notify_key is the project name/identifier
    apiKey: '', // API key is not returned in list, need separate call
    department: backend.description, // description maps to department
    companyId: backend.company_id,
    priority: priorityMap[backend.priority] || 'Normal',
    dailyLimit: backend.daily_limit,
    monthlyLimit: backend.monthly_limit,
  };
}

function transformCompany(backend: BackendCompany): Company {
  return {
    id: backend.id,
    name: backend.name,
    adminEmail: backend.contact_email,
    contactPhone: backend.contact_phone,
    address: backend.address,
    billingEnabled: backend.billing_enabled,
  };
}

function transformNotificationToLog(backend: BackendNotification): NotificationLog {
  // Map backend status to frontend MessageStatus
  const statusMap: Record<string, MessageStatus> = {
    'pending': MessageStatus.PENDING,
    'processing': MessageStatus.PENDING,
    'enqueued': MessageStatus.PENDING,
    'sent': MessageStatus.SENT,
    'failed': MessageStatus.FAILED,
    'cancelled': MessageStatus.FAILED,
  };

  // Map backend priority to frontend MessagePriority
  const priorityMap: Record<string, MessagePriority> = {
    'high': MessagePriority.HIGH,
    'normal': MessagePriority.NORMAL,
    'low': MessagePriority.LOW,
  };

  return {
    id: backend.id,
    projectId: backend.project_id,
    timestamp: backend.created_at,
    senderUpn: backend.sender_id || 'system',
    recipient: backend.targets?.[0] || 'all',
    recipientType: 'User', // Default, backend doesn't distinguish
    templateName: backend.message_type === 'adaptive_card' ? 'Adaptive Card' : 'Custom JSON',
    status: statusMap[backend.status] || MessageStatus.PENDING,
    priority: priorityMap[backend.priority] || MessagePriority.NORMAL,
    errorMessage: backend.error_message,
    contentPreview: backend.content?.substring(0, 100) || '',
    scheduledFor: backend.metadata?.scheduledFor,
  };
}

function transformBot(backend: BackendTeamsBot): BotInstance {
  const statusMap: Record<string, 'Healthy' | 'Degraded' | 'Offline'> = {
    'active': 'Healthy',
    'inactive': 'Offline',
    'suspended': 'Offline',
    'maintenance': 'Degraded',
  };

  return {
    id: backend.id,
    displayName: backend.name,
    appId: backend.app_id,
    tenantId: backend.tenant_id || '',
    status: statusMap[backend.status || 'inactive'] || 'Offline',
    usageCount: 0, // Backend doesn't track this
  };
}

function transformSystemHealth(backend: BackendSystemHealth): SystemHealth {
  // Extract component statuses
  const dbHealth = backend.components?.['database'];
  const redisHealth = backend.components?.['redis'];

  // Calculate aggregate metrics
  const isHealthy = backend.overall?.grade === 'A' || backend.overall?.grade === 'B';

  return {
    cpuUsage: 0, // Not provided by backend, would need /monitoring/performance
    memoryUsage: 0,
    queueLength: 0,
    apiLatencyMs: dbHealth?.latency_ms || 0,
    aadSyncStatus: isHealthy ? 'Healthy' : 'Error',
    lastSyncTime: backend.timestamp,
  };
}

function transformTemplate(backend: BackendTemplate): Template {
  // Convert backend variables (array of objects) to frontend (array of strings)
  // Backend uses TemplateVariable[] with {key, label, type, options}
  // Frontend expects string[] (just the keys)
  const variableKeys = backend.variables?.map(v => v.key || '') || [];

  return {
    id: backend.id,
    projectId: backend.project_id,
    name: backend.name,
    description: backend.description,
    variables: variableKeys,
    defaultJsonStructure: backend.default_json_structure, // Backend uses default_json_structure, not content
  };
}

// ============================================
// API Service Implementation
// ============================================

export const apiService = {
  // ============================================
  // Directory APIs
  // ============================================

  searchUsers: async (query: string): Promise<UserProfile[]> => {
    if (!query) return [];
    const response = await apiCall<UserProfile[]>(
      `/directory/users/search?q=${encodeURIComponent(query)}`
    );
    return extractData(response) || [];
  },

  getAllGroups: async (): Promise<ADGroup[]> => {
    const response = await apiCall<ADGroup[]>('/directory/groups');
    return extractData(response) || [];
  },

  getGroupChannels: async (groupId: string): Promise<TeamChannel[]> => {
    const response = await apiCall<TeamChannel[]>(
      `/directory/groups/${groupId}/channels`
    );
    return extractData(response) || [];
  },

  syncDirectory: async (): Promise<void> => {
    await apiCall<null>('/directory/sync', { method: 'POST' }, true);
  },

  getSyncStatus: async (): Promise<SyncStatus> => {
    const response = await apiCall<SyncStatus>('/directory/sync/status', {}, true);
    const data = extractData(response);
    if (!data) {
      throw new Error('Failed to get sync status');
    }
    return data;
  },

  // ============================================
  // Project APIs
  // ============================================

  getProjects: async (): Promise<Project[]> => {
    const response = await apiCall<BackendProject[]>('/projects');
    const backendProjects = extractData(response) || [];
    // Filter out global template project (UUID: 00000000-0000-0000-0000-000000000000)
    const GLOBAL_PROJECT_ID = '00000000-0000-0000-0000-000000000000';
    const filteredProjects = backendProjects.filter(p => p.id !== GLOBAL_PROJECT_ID);
    return filteredProjects.map(transformProject);
  },

  getProject: async (projectId: string): Promise<Project | null> => {
    const response = await apiCall<BackendProject>(`/projects/${projectId}`);
    const backend = extractData(response);
    return backend ? transformProject(backend) : null;
  },

  updateProject: async (projectId: string, data: Partial<Project>): Promise<void> => {
    // Transform frontend fields to backend fields
    // Backend expects: notifyKey, description, priority (camelCase)
    const backendData: any = {};
    if (data.name) backendData.notifyKey = data.name;
    if (data.department) backendData.description = data.department;
    // Priority can be updated if provided
    if (data.priority) {
      // Map frontend priority to backend priority
      const priorityMap: Record<string, string> = {
        'High (Alert)': 'high',
        'Normal': 'normal',
        'Low (Broadcast)': 'low',
      };
      backendData.priority = priorityMap[data.priority] || data.priority.toLowerCase();
    }

    await apiCall<BackendProject>(`/projects/${projectId}`, {
      method: 'PUT',
      body: JSON.stringify(backendData),
    });
  },

  createProject: async (data: {
    companyId: string;
    notifyKey: string;
    description: string;
    dailyLimit?: number;
    monthlyLimit?: number;
    priority?: string;
    createdBy: string;
  }): Promise<Project> => {
    // Backend expects camelCase: companyId, notifyKey, description, dailyLimit, monthlyLimit, priority, createdBy
    const requestBody = {
      companyId: data.companyId,
      notifyKey: data.notifyKey,
      description: data.description,
      dailyLimit: data.dailyLimit || 10000,
      monthlyLimit: data.monthlyLimit || 300000,
      priority: data.priority || 'normal',
      createdBy: data.createdBy,
    };

    const response = await apiCall<BackendProject>('/projects', {
      method: 'POST',
      body: JSON.stringify(requestBody),
    });

    const backend = extractData(response);
    if (!backend) {
      // Try to extract error message from response
      const errorMsg = (response as any).error || (response as any).message || 'Failed to create project';
      throw new Error(errorMsg);
    }
    return transformProject(backend);
  },

  createCompany: async (data: {
    name: string;
    contactEmail: string;
    contactPhone: string;
    address: string;
    billingEnabled?: boolean;
  }): Promise<Company> => {
    const response = await apiCall<BackendCompany>('/companies', {
      method: 'POST',
      body: JSON.stringify({
        name: data.name,
        contactEmail: data.contactEmail,
        contactPhone: data.contactPhone,
        address: data.address,
        billingEnabled: data.billingEnabled || false,
      }),
    });
    const backend = extractData(response);
    if (!backend) {
      throw new Error('Failed to create company');
    }
    return transformCompany(backend);
  },

  updateCompany: async (companyId: string, data: Partial<Company>): Promise<void> => {
    const backendData: any = {};
    if (data.name) backendData.name = data.name;
    if (data.adminEmail) backendData.contactEmail = data.adminEmail;
    if (data.contactPhone !== undefined) backendData.contactPhone = data.contactPhone;
    if (data.address !== undefined) backendData.address = data.address;
    if (data.billingEnabled !== undefined) backendData.billingEnabled = data.billingEnabled;

    await apiCall<BackendCompany>(`/companies/${companyId}`, {
      method: 'PUT',
      body: JSON.stringify(backendData),
    });
  },

  // ============================================
  // Template APIs
  // ============================================

  getTemplates: async (projectId: string): Promise<Template[]> => {
    // Backend uses 'id' as path parameter, not 'projectId'
    const response = await apiCall<BackendTemplate[]>(
      `/projects/${projectId}/templates`
    );
    const backendTemplates = extractData(response) || [];
    return backendTemplates.map(transformTemplate);
  },

  createTemplate: async (projectId: string, data: Omit<Template, 'id'>): Promise<Template> => {
    // Backend expects camelCase: name, description, defaultJsonStructure, variables (TemplateVariable[])
    const backendData = {
      name: data.name,
      description: data.description || '',
      defaultJsonStructure: data.defaultJsonStructure,
      variables: (data.variables || []).map(key => ({
        key: typeof key === 'string' ? key : key,
        label: typeof key === 'string' ? key : (key as any).label || key,
        type: typeof key === 'string' ? 'text' : (key as any).type || 'text'
      })),
    };

    // Backend uses 'id' as path parameter
    const response = await apiCall<BackendTemplate>(
      `/projects/${projectId}/templates`,
      {
        method: 'POST',
        body: JSON.stringify(backendData),
      }
    );
    const backend = extractData(response);
    if (!backend) {
      throw new Error('Failed to create template: No data returned from server');
    }
    return transformTemplate(backend);
  },

  updateTemplate: async (projectId: string, templateId: string, data: Partial<Template>): Promise<void> => {
    // Backend expects camelCase: name, description, defaultJsonStructure, variables
    const backendData: any = {};
    if (data.name) backendData.name = data.name;
    if (data.description) backendData.description = data.description;
    if (data.defaultJsonStructure) backendData.defaultJsonStructure = data.defaultJsonStructure;
    if (data.variables) {
      backendData.variables = data.variables.map(key => ({ key, label: key, type: 'text' }));
    }

    // Backend uses 'id' as path parameter for project
    await apiCall<BackendTemplate>(
      `/projects/${projectId}/templates/${templateId}`,
      {
        method: 'PUT',
        body: JSON.stringify(backendData),
      }
    );
  },

  deleteTemplate: async (projectId: string, templateId: string): Promise<void> => {
    // Backend uses 'id' as path parameter for project
    await apiCall<null>(
      `/projects/${projectId}/templates/${templateId}`,
      { method: 'DELETE' }
    );
  },

  // ============================================
  // Global Template APIs
  // ============================================

  getGlobalTemplates: async (): Promise<Template[]> => {
    const response = await apiCall<BackendTemplate[]>(
      `/global/templates`
    );
    const backendTemplates = extractData(response) || [];
    return backendTemplates.map(transformTemplate);
  },

  createGlobalTemplate: async (data: Omit<Template, 'id'>): Promise<Template> => {
    const backendData = {
      name: data.name,
      description: data.description,
      defaultJsonStructure: data.defaultJsonStructure,
      variables: data.variables.map(key => ({
        key,
        label: key,
        type: 'text'
      })),
    };

    const response = await apiCall<BackendTemplate>(
      `/global/templates`,
      {
        method: 'POST',
        body: JSON.stringify(backendData),
      }
    );
    const backend = extractData(response);
    return backend ? transformTemplate(backend) : {
      id: '',
      projectId: 'global',
      name: data.name,
      description: data.description,
      variables: data.variables,
      defaultJsonStructure: data.defaultJsonStructure,
    };
  },

  updateGlobalTemplate: async (templateId: string, data: Partial<Template>): Promise<void> => {
    const backendData: any = {};
    if (data.name) backendData.name = data.name;
    if (data.description) backendData.description = data.description;
    if (data.defaultJsonStructure) backendData.defaultJsonStructure = data.defaultJsonStructure;
    if (data.variables) {
      backendData.variables = data.variables.map(key => ({ key, label: key, type: 'text' }));
    }

    await apiCall<BackendTemplate>(
      `/global/templates/${templateId}`,
      {
        method: 'PUT',
        body: JSON.stringify(backendData),
      }
    );
  },

  deleteGlobalTemplate: async (templateId: string): Promise<void> => {
    await apiCall<null>(
      `/global/templates/${templateId}`,
      { method: 'DELETE' }
    );
  },

  // ============================================
  // Chat Group APIs
  // ============================================

  getChatGroups: async (projectId: string): Promise<ChatGroup[]> => {
    // Backend uses 'id' as path parameter, not 'projectId'
    const response = await apiCall<ChatGroup[]>(
      `/projects/${projectId}/chat-groups`
    );
    return extractData(response) || [];
  },

  registerChatGroup: async (projectId: string, name: string, chatId: string): Promise<ChatGroup> => {
    // Backend uses 'id' as path parameter, not 'projectId'
    const response = await apiCall<ChatGroup>(
      `/projects/${projectId}/chat-groups`,
      {
        method: 'POST',
        body: JSON.stringify({ name, chatId }),
      }
    );
    return extractData(response) || { id: '', projectId, name, chatId, registeredBy: '', createdAt: '' };
  },

  // ============================================
  // Notification / Messaging APIs  
  // ============================================

  /**
   * Send notification via internal API
   * Backend endpoint: POST /internal/v1/notifications
   * Alternative: POST /api/v1/notify (external, uses notifyKey)
   * 
   * For UI, we use internal API which accepts projectId (UUID)
   * Supports both direct content and template-based notifications
   */
  sendMessage: async (request: SendMessageRequest): Promise<{ success: boolean; message: string }> => {
    // Use internal API with projectId (UUID)
    // Map frontend priority to backend priority
    const priorityMap: Record<string, string> = {
      'High (Alert)': 'high',
      'Normal': 'normal',
      'Low (Broadcast)': 'low',
    };

    // Build request body based on whether using template or direct content
    // Backend handler expects: projectId (UUID), messageType, content, templateId (optional), templateData (optional)
    const requestBody: any = {
      projectId: request.projectId,
      priority: priorityMap[request.priority] || 'normal',
      targets: request.recipient ? [request.recipient] : ['all'],
    };

    // Handle template-based notification
    if (request.isTemplate && request.templateId) {
      // Backend service supports templateId and templateData
      // The service will render the template using templateData
      requestBody.templateId = request.templateId;
      requestBody.messageType = 'adaptive_card'; // Templates are always adaptive cards

      // If templateData is provided, use it for backend rendering
      // Otherwise, use pre-rendered content if available
      if (request.templateData) {
        requestBody.templateData = request.templateData;
        // Content is optional when templateData is provided (backend will render)
        if (request.content) {
          requestBody.content = request.content;
        }
      } else if (request.content) {
        // Pre-rendered content (client-side rendering)
        requestBody.content = request.content;
      }
    } else {
      // Direct content notification
      requestBody.messageType = request.isTemplate ? 'adaptive_card' : 'text';
      requestBody.content = request.content;
    }

    // Add scheduled time if provided
    if (request.scheduledFor) {
      requestBody.metadata = { scheduledFor: request.scheduledFor };
    }

    const response = await apiCall<{
      notificationId?: string;
      status?: string;
      message?: string;
      destinationsCount?: number;
    }>(
      '/notifications',
      {
        method: 'POST',
        body: JSON.stringify(requestBody),
      },
      true // Use internal API
    );

    // Handle different response formats
    const data = extractData(response);
    const isSuccess = data !== undefined || ('code' in response && response.code === 200);

    // Extract message from response
    let message = 'Message sent successfully';
    if (data && 'message' in data) {
      message = data.message as string;
    } else if ('message' in response) {
      message = (response as any).message;
    } else if (!isSuccess) {
      message = (response as any).error || 'Failed to send';
    }

    return {
      success: isSuccess,
      message,
    };
  },

  /**
   * Get notification logs for a project
   * Backend endpoint: GET /internal/v1/notifications/project/:projectId
   */
  getLogs: async (projectId: string): Promise<NotificationLog[]> => {
    const response = await apiCall<BackendNotification[]>(
      `/notifications/project/${projectId}`
    );
    const backendNotifications = extractData(response) || [];
    return backendNotifications.map(transformNotificationToLog);
  },

  // ============================================
  // Audience Lists APIs
  // ============================================

  getAudienceLists: async (projectId: string): Promise<AudienceList[]> => {
    const response = await apiCall<AudienceList[]>(
      `/projects/${projectId}/audience-lists`,
      {},
      true // Use internal API
    );
    return extractData(response) || [];
  },

  uploadAudienceList: async (projectId: string, name: string, count: number): Promise<AudienceList> => {
    const response = await apiCall<AudienceList>(
      `/projects/${projectId}/audience-lists`,
      {
        method: 'POST',
        body: JSON.stringify({
          name,
          type: 'Static',
          count,
        }),
      },
      true // Use internal API
    );
    const data = extractData(response);
    if (!data) {
      throw new Error('Failed to create audience list');
    }
    return data;
  },

  // ============================================
  // API Key Management APIs
  // ============================================

  regenerateApiKey: async (projectId: string): Promise<string> => {
    const response = await apiCall<{ apiKey: string }>(
      `/projects/${projectId}/regenerate-key`,
      {
        method: 'POST',
      },
      true // Use internal API
    );
    const data = extractData(response);
    if (!data || !data.apiKey) {
      throw new Error('Failed to regenerate API key');
    }
    return data.apiKey;
  },

  // ============================================
  // Admin / Monitoring APIs
  // ============================================

  getSystemHealth: async (): Promise<SystemHealth> => {
    const response = await apiCall<BackendSystemHealth>('/monitoring/health');
    const backend = extractData(response);

    if (backend) {
      return transformSystemHealth(backend);
    }

    return {
      cpuUsage: 0,
      memoryUsage: 0,
      queueLength: 0,
      apiLatencyMs: 0,
      aadSyncStatus: 'Error',
      lastSyncTime: new Date().toISOString(),
    };
  },

  /**
   * Get billing records (usage summary)
   * Backend endpoint: GET /internal/v1/billing/usage/summary
   */
  getBillingRecords: async (): Promise<BillingRecord[]> => {
    interface UsageSummaryData {
      total_notifications: number;
      total_cost: number;
      breakdown_by_project?: Array<{
        project_id: string;
        project_name?: string;
        company_name?: string;
        notification_count: number;
        total_cost: number;
      }>;
    }

    const response = await apiCall<UsageSummaryData>('/billing/usage/summary');
    const data = extractData(response);

    if (!data?.breakdown_by_project) {
      return [];
    }

    // Transform to BillingRecord format
    const currentMonth = new Date().toISOString().slice(0, 7);
    return data.breakdown_by_project.map((item, index) => ({
      id: `billing-${index}`,
      period: currentMonth,
      companyName: item.company_name || 'Unknown',
      projectName: item.project_name || item.project_id,
      messageCount: item.notification_count,
      cost: item.total_cost,
      currency: 'USD',
    }));
  },

  getCompanies: async (): Promise<Company[]> => {
    const response = await apiCall<BackendCompany[]>('/companies');
    const backendCompanies = extractData(response) || [];
    return backendCompanies.map(transformCompany);
  },

  getCompany: async (companyId: string): Promise<Company | null> => {
    const response = await apiCall<BackendCompany>(`/companies/${companyId}`);
    const backend = extractData(response);
    return backend ? transformCompany(backend) : null;
  },

  /**
   * Get bot instances
   * Backend endpoint: GET /internal/v1/bots/platform
   */
  getBots: async (): Promise<BotInstance[]> => {
    const response = await apiCall<BackendTeamsBot[]>('/bots/platform');
    const backendBots = extractData(response) || [];
    return backendBots.map(transformBot);
  },

  /**
   * Register a new bot
   * Backend endpoint: POST /internal/v1/bots/platform
   */
  registerBot: async (data: Omit<BotInstance, 'id' | 'status' | 'usageCount'>): Promise<BotInstance> => {
    const response = await apiCall<BackendTeamsBot>('/bots/platform', {
      method: 'POST',
      body: JSON.stringify({
        name: data.displayName,
        description: `Bot ${data.displayName}`,
        app_id: data.appId,
        app_password: 'placeholder', // Required by backend, should come from secure input
        tenant_id: data.tenantId,
        webhook_url: 'https://placeholder.com/webhook', // Required by backend
        type: 'platform',
        status: 'active',
      }),
    });

    const backend = extractData(response);
    if (backend) {
      return transformBot(backend);
    }

    return {
      id: '',
      displayName: data.displayName,
      appId: data.appId,
      tenantId: data.tenantId,
      status: 'Healthy',
      usageCount: 0,
    };
  },

  /**
   * Get transaction history (usage records)
   * Backend endpoint: GET /internal/v1/billing/usage
   */
  getTransactionHistory: async (filter: {
    companyId?: string;
    projectId?: string;
    month?: string
  }): Promise<TransactionRecord[]> => {
    const params = new URLSearchParams();
    if (filter.companyId) params.append('company_id', filter.companyId);
    if (filter.projectId) params.append('project_id', filter.projectId);
    // Backend uses start_date/end_date instead of month
    if (filter.month) {
      const [year, month] = filter.month.split('-');
      const lastDay = new Date(parseInt(year), parseInt(month), 0).getDate();
      params.append('start_date', `${filter.month}-01`);
      params.append('end_date', `${filter.month}-${lastDay}`);
    }

    const queryString = params.toString();

    interface UsageRecordsData {
      records: BackendUsageRecord[];
      total: number;
      page: number;
      page_size: number;
      total_pages: number;
    }

    const response = await apiCall<UsageRecordsData>(
      `/billing/usage${queryString ? `?${queryString}` : ''}`
    );

    const data = extractData(response);
    const records = data?.records || [];

    // Transform backend UsageRecord to frontend TransactionRecord
    return records.map((r: BackendUsageRecord) => ({
      id: r.id,
      timestamp: r.created_at,
      companyId: r.company_id,
      companyName: 'Unknown', // Would need to join with company data
      projectId: r.project_id,
      projectName: 'Unknown', // Would need to join with project data
      type: r.record_type === 'notification' ? 'Message_Sent' as const :
        r.record_type === 'attachment' ? 'Bot_Maintenance' as const :
          'Subscription_Fee' as const,
      quantity: r.quantity,
      amount: r.total_amount,
      currency: 'USD',
      status: 'Paid' as const,
    }));
  },
};

export default apiService;
