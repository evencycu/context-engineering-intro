/**
 * Service Layer - Hybrid API with Mock Fallback
 * 
 * This module provides a unified API service that:
 * 1. Tries to call the real backend API first
 * 2. Falls back to mock data if backend is unavailable
 * 3. Can be configured to use mock-only mode for development
 */

import { apiService } from './apiService';
import { mockApiService } from './mockApiService';
import type { 
  UserProfile, 
  NotificationLog, 
  SendMessageRequest, 
  Project, 
  AudienceList, 
  SystemHealth, 
  BillingRecord, 
  Company, 
  BotInstance, 
  TransactionRecord, 
  ADGroup, 
  TeamChannel, 
  ChatGroup, 
  Template 
} from '../types';

// Configuration
const USE_MOCK_ONLY = false; // Set to true to always use mock data
const ENABLE_FALLBACK = true; // Set to false to disable fallback to mock

// Check if backend is available
let backendAvailable: boolean | null = null;

async function checkBackendHealth(): Promise<boolean> {
  try {
    const response = await fetch('/health', { 
      method: 'GET',
      signal: AbortSignal.timeout(3000), // 3 second timeout
    });
    return response.ok;
  } catch {
    return false;
  }
}

// Wrapper function that handles fallback logic
async function withFallback<T>(
  realCall: () => Promise<T>,
  mockCall: () => Promise<T>,
  operationName: string
): Promise<T> {
  if (USE_MOCK_ONLY) {
    console.log(`[${operationName}] Using mock data (mock-only mode)`);
    return mockCall();
  }

  try {
    const result = await realCall();
    if (backendAvailable !== true) {
      backendAvailable = true;
      console.log('[API] Backend is available');
    }
    return result;
  } catch (error) {
    console.warn(`[${operationName}] API call failed:`, error);
    
    if (ENABLE_FALLBACK) {
      console.log(`[${operationName}] Falling back to mock data`);
      backendAvailable = false;
      return mockCall();
    }
    
    throw error;
  }
}

// ============================================
// Unified Service Export
// ============================================

export const service = {
  // Check if currently using mock
  isUsingMock: () => backendAvailable === false || USE_MOCK_ONLY,
  
  // Check backend health
  checkHealth: checkBackendHealth,

  // Directory APIs
  searchUsers: (query: string): Promise<UserProfile[]> => 
    withFallback(
      () => apiService.searchUsers(query),
      () => mockApiService.searchUsers(query),
      'searchUsers'
    ),

  getAllGroups: (): Promise<ADGroup[]> =>
    withFallback(
      () => apiService.getAllGroups(),
      () => mockApiService.getAllGroups(),
      'getAllGroups'
    ),

  getGroupChannels: (groupId: string): Promise<TeamChannel[]> =>
    withFallback(
      () => apiService.getGroupChannels(groupId),
      () => mockApiService.getGroupChannels(groupId),
      'getGroupChannels'
    ),

  // Project APIs
  getProjects: (): Promise<Project[]> =>
    withFallback(
      () => apiService.getProjects(),
      () => mockApiService.getProjects(),
      'getProjects'
    ),

  getProject: (projectId: string): Promise<Project | null> =>
    withFallback(
      () => apiService.getProject(projectId),
      () => mockApiService.getProject(projectId),
      'getProject'
    ),

  updateProject: (projectId: string, data: Partial<Project>): Promise<void> =>
    withFallback(
      () => apiService.updateProject(projectId, data),
      () => mockApiService.updateProject(projectId, data),
      'updateProject'
    ),

  createProject: (data: {
    companyId: string;
    notifyKey: string;
    description: string;
    dailyLimit?: number;
    monthlyLimit?: number;
    priority?: string;
    createdBy: string;
  }): Promise<Project> =>
    withFallback(
      () => apiService.createProject(data),
      () => mockApiService.createProject(data),
      'createProject'
    ),

  createCompany: (data: {
    name: string;
    contactEmail: string;
    contactPhone: string;
    address: string;
    billingEnabled?: boolean;
  }): Promise<Company> =>
    withFallback(
      () => apiService.createCompany(data),
      () => mockApiService.createCompany(data),
      'createCompany'
    ),

  updateCompany: (companyId: string, data: Partial<Company>): Promise<void> =>
    withFallback(
      () => apiService.updateCompany(companyId, data),
      () => mockApiService.updateCompany(companyId, data),
      'updateCompany'
    ),

  // Template APIs
  getTemplates: (projectId: string): Promise<Template[]> =>
    withFallback(
      () => apiService.getTemplates(projectId),
      () => mockApiService.getTemplates(projectId),
      'getTemplates'
    ),

  getGlobalTemplates: (): Promise<Template[]> =>
    withFallback(
      async () => {
        // Backend might not have global templates endpoint yet
        // For now, return empty array from real API
        return [];
      },
      () => mockApiService.getGlobalTemplates(),
      'getGlobalTemplates'
    ),

  createGlobalTemplate: (data: Omit<Template, 'id'>): Promise<Template> =>
    withFallback(
      async () => {
        // Backend might not have global templates endpoint yet
        // For now, always use mock (don't throw error, just use mock directly)
        return mockApiService.createTemplate({ ...data, projectId: 'global' });
      },
      () => mockApiService.createTemplate({ ...data, projectId: 'global' }),
      'createGlobalTemplate'
    ),

  updateGlobalTemplate: (id: string, data: Partial<Template>): Promise<void> =>
    withFallback(
      async () => {
        // Backend might not have global templates endpoint yet
        // For now, always use mock (don't throw error, just use mock directly)
        return mockApiService.updateTemplate(id, data);
      },
      () => mockApiService.updateTemplate(id, data),
      'updateGlobalTemplate'
    ),

  deleteGlobalTemplate: (id: string): Promise<void> =>
    withFallback(
      async () => {
        // Backend might not have global templates endpoint yet
        // For now, always use mock (don't throw error, just use mock directly)
        return mockApiService.deleteTemplate(id);
      },
      () => mockApiService.deleteTemplate(id),
      'deleteGlobalTemplate'
    ),

  createTemplate: (data: Omit<Template, 'id'>): Promise<Template> =>
    withFallback(
      () => apiService.createTemplate(data.projectId, data),
      () => mockApiService.createTemplate(data),
      'createTemplate'
    ),

  updateTemplate: (id: string, data: Partial<Template>): Promise<void> =>
    withFallback(
      async () => {
        if (!data.projectId) {
          throw new Error('Project ID is required to update template');
        }
        await apiService.updateTemplate(data.projectId, id, data);
      },
      () => mockApiService.updateTemplate(id, data),
      'updateTemplate'
    ),

  deleteTemplate: (id: string, projectId?: string): Promise<void> =>
    withFallback(
      async () => {
        if (!projectId) {
          throw new Error('Project ID is required to delete template');
        }
        await apiService.deleteTemplate(projectId, id);
      },
      () => mockApiService.deleteTemplate(id),
      'deleteTemplate'
    ),

  // Chat Group APIs
  getChatGroups: (projectId: string): Promise<ChatGroup[]> =>
    withFallback(
      () => apiService.getChatGroups(projectId),
      () => mockApiService.getChatGroups(projectId),
      'getChatGroups'
    ),

  registerChatGroup: (projectId: string, name: string, chatId: string): Promise<ChatGroup> =>
    withFallback(
      () => apiService.registerChatGroup(projectId, name, chatId),
      () => mockApiService.registerChatGroup(projectId, name, chatId),
      'registerChatGroup'
    ),

  // Notification APIs
  sendMessage: (request: SendMessageRequest): Promise<{ success: boolean; message: string }> =>
    withFallback(
      () => apiService.sendMessage(request),
      () => mockApiService.sendMessage(request),
      'sendMessage'
    ),

  getLogs: (projectId: string): Promise<NotificationLog[]> =>
    withFallback(
      () => apiService.getLogs(projectId),
      () => mockApiService.getLogs(projectId),
      'getLogs'
    ),

  // Admin APIs
  getSystemHealth: (): Promise<SystemHealth> =>
    withFallback(
      () => apiService.getSystemHealth(),
      () => mockApiService.getSystemHealth(),
      'getSystemHealth'
    ),

  getBillingRecords: (): Promise<BillingRecord[]> =>
    withFallback(
      () => apiService.getBillingRecords(),
      () => mockApiService.getBillingRecords(),
      'getBillingRecords'
    ),

  getCompanies: (): Promise<Company[]> =>
    withFallback(
      () => apiService.getCompanies(),
      () => mockApiService.getCompanies(),
      'getCompanies'
    ),

  getCompany: (companyId: string): Promise<Company | null> =>
    withFallback(
      () => apiService.getCompany(companyId),
      () => mockApiService.getCompany(companyId),
      'getCompany'
    ),

  getBots: (): Promise<BotInstance[]> =>
    withFallback(
      () => apiService.getBots(),
      () => mockApiService.getBots(),
      'getBots'
    ),

  registerBot: (data: Omit<BotInstance, 'id' | 'status' | 'usageCount'>): Promise<BotInstance> =>
    withFallback(
      () => apiService.registerBot(data),
      () => mockApiService.registerBot(data),
      'registerBot'
    ),

  getTransactionHistory: (filter: { companyId?: string; projectId?: string; month?: string }): Promise<TransactionRecord[]> =>
    withFallback(
      () => apiService.getTransactionHistory(filter),
      () => mockApiService.getTransactionHistory(filter),
      'getTransactionHistory'
    ),

  // Audience APIs (currently mock-only, backend needs implementation)
  getAudienceLists: (projectId: string): Promise<AudienceList[]> =>
    mockApiService.getAudienceLists(projectId),

  uploadAudienceList: (projectId: string, name: string, count: number): Promise<AudienceList> =>
    mockApiService.uploadAudienceList(projectId, name, count),

  // User tagging (currently mock-only)
  addTagToUser: (userId: string, tag: string): Promise<void> =>
    mockApiService.addTagToUser(userId, tag),

  removeTagFromUser: (userId: string, tag: string): Promise<void> =>
    mockApiService.removeTagFromUser(userId, tag),

  getAllUsers: (): Promise<UserProfile[]> =>
    mockApiService.getAllUsers(),

  getProjectUsers: (projectId: string): Promise<UserProfile[]> =>
    mockApiService.getProjectUsers(projectId),

  // Channel sync (currently mock-only)
  syncChannels: (groupId: string): Promise<void> =>
    mockApiService.syncChannels(groupId),

  // Global template assignment (currently mock-only)
  assignGlobalTemplateToProjects: (templateId: string, projectIds: string[]): Promise<void> =>
    mockApiService.assignGlobalTemplateToProjects(templateId, projectIds),

  // System actions (currently mock-only)
  performSystemAction: (action: 'sync_aad' | 'restart_service' | 'block_ip'): Promise<void> =>
    mockApiService.performSystemAction(action),

  // API Key regeneration (currently mock-only)
  regenerateApiKey: (projectId: string): Promise<string> =>
    mockApiService.regenerateApiKey(projectId),
};

// Re-export for direct access if needed
export { apiService } from './apiService';
export { mockApiService } from './mockApiService';

export default service;
