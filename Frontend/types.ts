
export interface Project {
  id: string;
  name: string; // Maps to backend notify_key
  apiKey: string;
  department: string; // Maps to backend description
  companyId: string;
  priority?: string; // Low (Broadcast), Normal, High (Alert)
  dailyLimit?: number; // Maps to backend daily_limit
  monthlyLimit?: number; // Maps to backend monthly_limit
}

export interface UserProfile {
  id: string;
  displayName: string;
  userPrincipalName: string; // Email
  jobTitle: string;
  tags?: string[];
}

export enum TemplateType {
  SYSTEM_ALERT = 'SYSTEM_ALERT',
  HOLIDAY_NOTICE = 'HOLIDAY_NOTICE',
  DEPLOYMENT_SUCCESS = 'DEPLOYMENT_SUCCESS',
}

export interface Template {
  id: string;
  projectId: string;
  name: string;
  description: string;
  variables: string[]; // e.g. ['title', 'message', 'severity']
  defaultJsonStructure: string;
}

export enum MessageStatus {
  SENT = 'SENT',
  FAILED = 'FAILED',
  PENDING = 'PENDING',
  SCHEDULED = 'SCHEDULED',
}

export enum MessagePriority {
  HIGH = 'High (Alert)',
  NORMAL = 'Normal',
  LOW = 'Low (Broadcast)',
}

export interface NotificationLog {
  id: string;
  projectId: string;
  timestamp: string;
  senderUpn: string;
  recipient: string; // Can be email, tag, or list name
  recipientType: 'User' | 'Channel' | 'Tag' | 'List';
  templateName: string | 'Custom JSON';
  status: MessageStatus;
  priority: MessagePriority;
  errorMessage?: string;
  contentPreview: string; // JSON string
  scheduledFor?: string;
}

export interface SendMessageRequest {
  projectId: string;
  recipient: string;
  recipientType: 'User' | 'Channel' | 'Tag' | 'List';
  content: string; // JSON string (pre-rendered content or direct content)
  isTemplate: boolean;
  templateId?: string;
  templateData?: Record<string, any>; // Template variable values for backend rendering (optional)
  priority: MessagePriority;
  scheduledFor?: string; // ISO String
}

// --- New Feature Types ---

export interface AudienceList {
  id: string;
  projectId: string;
  name: string;
  type: 'Static' | 'Dynamic'; // Uploaded vs Query
  count: number;
  lastUpdated: string;
}

export interface ChatGroup {
  id: string;
  projectId: string;
  name: string; // Friendly name used in the app
  chatId: string; // The 19:... thread ID from Teams
  registeredBy: string;
  createdAt: string;
}

export interface SystemHealth {
  cpuUsage: number;
  memoryUsage: number;
  queueLength: number;
  apiLatencyMs: number;
  aadSyncStatus: 'Healthy' | 'Syncing' | 'Error';
  lastSyncTime: string;
}

export interface BillingRecord {
  id: string;
  period: string; // "2023-10"
  companyName: string;
  projectName: string;
  messageCount: number;
  cost: number;
  currency: string;
}

export interface TransactionRecord {
  id: string;
  timestamp: string;
  companyId: string;
  companyName: string;
  projectId: string;
  projectName: string;
  type: 'Message_Sent' | 'Bot_Maintenance' | 'Subscription_Fee';
  quantity: number;
  amount: number;
  currency: string;
  status: 'Paid' | 'Pending' | 'Failed';
}

// --- Admin Types ---

export interface Company {
  id: string;
  name: string;
  adminEmail: string;
  contactPhone?: string;
  address?: string;
  billingEnabled?: boolean;
}

export interface BotInstance {
  id: string;
  displayName: string;
  appId: string; // Microsoft App ID (GUID)
  tenantId: string; // Azure Tenant ID (GUID)
  status: 'Healthy' | 'Degraded' | 'Offline';
  usageCount: number;
}

export interface ADGroup {
  id: string;
  displayName: string;
  mailNickname: string;
  description: string;
  groupTypes: string[]; // e.g. ["Unified"] for M365 groups
  memberCount: number;
}

export interface TeamChannel {
  id: string;
  displayName: string;
  membershipType: 'Standard' | 'Private' | 'Shared';
  description?: string;
}
