
import { MOCK_USERS, generateMockLogs, MOCK_AUDIENCE_LISTS, MOCK_BILLING, MOCK_COMPANIES, MOCK_BOTS, MOCK_PROJECTS, MOCK_TRANSACTIONS, MOCK_AD_GROUPS, MOCK_CHANNELS, MOCK_CHAT_GROUPS, MOCK_TEMPLATES, MOCK_GLOBAL_TEMPLATES } from '../constants';
import { UserProfile, NotificationLog, SendMessageRequest, MessageStatus, Project, AudienceList, SystemHealth, BillingRecord, Company, BotInstance, TransactionRecord, ADGroup, TeamChannel, ChatGroup, Template } from '../types';

// Simulate API delay
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// Mutable in-memory store for the session
let currentUsers = JSON.parse(JSON.stringify(MOCK_USERS)) as UserProfile[];
let currentLists = JSON.parse(JSON.stringify(MOCK_AUDIENCE_LISTS)) as AudienceList[];
let currentBots = JSON.parse(JSON.stringify(MOCK_BOTS)) as BotInstance[];
let currentADGroups = JSON.parse(JSON.stringify(MOCK_AD_GROUPS)) as ADGroup[];
let currentChatGroups = JSON.parse(JSON.stringify(MOCK_CHAT_GROUPS)) as ChatGroup[];
let currentTemplates = JSON.parse(JSON.stringify(MOCK_TEMPLATES)) as Template[];
let currentGlobalTemplates = JSON.parse(JSON.stringify(MOCK_GLOBAL_TEMPLATES)) as Template[];

export const mockApiService = {
  searchUsers: async (query: string): Promise<UserProfile[]> => {
    await delay(500); // Simulate network latency
    if (!query) return [];
    const lowerQuery = query.toLowerCase();
    return currentUsers.filter(u =>
      u.displayName.toLowerCase().includes(lowerQuery) ||
      u.userPrincipalName.toLowerCase().includes(lowerQuery) ||
      u.tags?.some(t => t.toLowerCase().includes(lowerQuery))
    );
  },

  getAllUsers: async (): Promise<UserProfile[]> => {
    await delay(600);
    return currentUsers;
  },

  getAllGroups: async (): Promise<ADGroup[]> => {
    await delay(600);
    return currentADGroups;
  },

  // --- Channel Management (Admin) ---

  getGroupChannels: async (groupId: string): Promise<TeamChannel[]> => {
    await delay(800);
    return MOCK_CHANNELS[groupId] || [];
  },

  syncChannels: async (groupId: string): Promise<void> => {
    await delay(1500);
    console.log(`[Sync] Channels for group ${groupId} synced.`);
  },

  getProject: async (projectId: string): Promise<Project | null> => {
    await delay(600);
    return MOCK_PROJECTS.find(p => p.id === projectId) || null;
  },

  getProjects: async (): Promise<Project[]> => {
    await delay(600);
    return MOCK_PROJECTS;
  },

  getProjectUsers: async (projectId: string): Promise<UserProfile[]> => {
    await delay(600);
    // In a real scenario, we would filter by project permissions.
    // For this mock, we return all users as "Available Project Users".
    return [...currentUsers];
  },

  sendMessage: async (request: SendMessageRequest): Promise<{ success: boolean; message: string }> => {
    console.log(`[Sending to Project: ${request.projectId}]`, request);
    await delay(1200);

    // Random failure simulation
    if (Math.random() > 0.98) {
      throw new Error('Upstream Bot Framework Timeout');
    }

    if (request.scheduledFor) {
      return { success: true, message: `Message scheduled for ${new Date(request.scheduledFor).toLocaleString()}` };
    }

    return { success: true, message: 'Message queued for delivery.' };
  },

  getLogs: async (projectId: string): Promise<NotificationLog[]> => {
    await delay(800);
    const allLogs = generateMockLogs();
    return allLogs.filter(log => log.projectId === projectId);
  },

  updateProject: async (projectId: string, data: Partial<Project>): Promise<void> => {
    await delay(1000);
    console.log(`[Update Project ${projectId}]`, data);
    return;
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
    await delay(1000);
    const newProject: Project = {
      id: `proj_${Math.random().toString(36).substring(2, 9)}`,
      name: data.notifyKey,
      apiKey: `sk_live_${Math.random().toString(36).substring(2, 15)}`,
      department: data.description,
      companyId: data.companyId,
    };
    MOCK_PROJECTS.push(newProject);
    return newProject;
  },

  createCompany: async (data: {
    name: string;
    contactEmail: string;
    contactPhone: string;
    address: string;
    billingEnabled?: boolean;
  }): Promise<Company> => {
    await delay(1000);
    const newCompany: Company = {
      id: `comp_${Math.random().toString(36).substring(2, 9)}`,
      name: data.name,
      adminEmail: data.contactEmail,
    };
    MOCK_COMPANIES.push(newCompany);
    return newCompany;
  },

  updateCompany: async (companyId: string, data: Partial<Company>): Promise<void> => {
    await delay(1000);
    const company = MOCK_COMPANIES.find(c => c.id === companyId);
    if (company) {
      if (data.name) company.name = data.name;
      if (data.adminEmail) company.adminEmail = data.adminEmail;
      if (data.contactPhone !== undefined) company.contactPhone = data.contactPhone;
      if (data.address !== undefined) company.address = data.address;
      if (data.billingEnabled !== undefined) company.billingEnabled = data.billingEnabled;
    }
  },

  regenerateApiKey: async (projectId: string): Promise<string> => {
    await delay(1500);
    const newKey = `sk_live_${Math.random().toString(36).substring(2, 15)}_${Date.now()}`;
    return newKey;
  },

  // --- Audience & Tagging APIs ---

  getAudienceLists: async (projectId: string): Promise<AudienceList[]> => {
    await delay(600);
    return currentLists.filter(l => l.projectId === projectId);
  },

  uploadAudienceList: async (projectId: string, name: string, count: number): Promise<AudienceList> => {
    await delay(1500);
    const newList: AudienceList = {
      id: `list_${Date.now()}`,
      projectId,
      name,
      type: 'Static',
      count,
      lastUpdated: new Date().toISOString().split('T')[0]
    };
    currentLists.push(newList);
    return newList;
  },

  addTagToUser: async (userId: string, tag: string): Promise<void> => {
    await delay(300);
    const userIndex = currentUsers.findIndex(u => u.id === userId);
    if (userIndex > -1) {
      const user = currentUsers[userIndex];
      const tags = user.tags || [];
      // Clean tag string
      const cleanTag = tag.trim();
      if (!tags.includes(cleanTag)) {
        currentUsers[userIndex] = { ...user, tags: [...tags, cleanTag] };
      }
    }
  },

  removeTagFromUser: async (userId: string, tag: string): Promise<void> => {
    await delay(300);
    const userIndex = currentUsers.findIndex(u => u.id === userId);
    if (userIndex > -1) {
      const user = currentUsers[userIndex];
      if (user.tags) {
        currentUsers[userIndex] = { ...user, tags: user.tags.filter(t => t !== tag) };
      }
    }
  },

  // --- Chat Group APIs ---

  getChatGroups: async (projectId: string): Promise<ChatGroup[]> => {
    await delay(600);
    return currentChatGroups.filter(c => c.projectId === projectId);
  },

  registerChatGroup: async (projectId: string, name: string, chatId: string): Promise<ChatGroup> => {
    await delay(1000);
    const newGroup: ChatGroup = {
      id: `cg_${Date.now()}`,
      projectId,
      name,
      chatId,
      registeredBy: 'current.user@demo.com', // Mock user
      createdAt: new Date().toISOString().split('T')[0]
    };
    currentChatGroups.push(newGroup);
    return newGroup;
  },

  // --- Template APIs ---

  getTemplates: async (projectId: string): Promise<Template[]> => {
    await delay(500);
    return currentTemplates.filter(t => t.projectId === projectId);
  },

  getGlobalTemplates: async (): Promise<Template[]> => {
    await delay(500);
    return currentGlobalTemplates;
  },

  createTemplate: async (data: Omit<Template, 'id'>): Promise<Template> => {
    await delay(800);
    const newTemplate = { ...data, id: `tpl_${Date.now()}` };
    if (data.projectId === 'global') {
      currentGlobalTemplates.push(newTemplate);
    } else {
      currentTemplates.push(newTemplate);
    }
    return newTemplate;
  },

  updateTemplate: async (id: string, data: Partial<Template>): Promise<void> => {
    await delay(800);
    const globalIdx = currentGlobalTemplates.findIndex(t => t.id === id);
    if (globalIdx > -1) {
      currentGlobalTemplates[globalIdx] = { ...currentGlobalTemplates[globalIdx], ...data };
      return;
    }
    const idx = currentTemplates.findIndex(t => t.id === id);
    if (idx > -1) {
      currentTemplates[idx] = { ...currentTemplates[idx], ...data };
    }
  },

  deleteTemplate: async (id: string): Promise<void> => {
    await delay(500);
    const globalIdx = currentGlobalTemplates.findIndex(t => t.id === id);
    if (globalIdx > -1) {
      currentGlobalTemplates.splice(globalIdx, 1);
      return;
    }
    const idx = currentTemplates.findIndex(t => t.id === id);
    if (idx > -1) {
      currentTemplates.splice(idx, 1);
    }
  },

  assignGlobalTemplateToProjects: async (templateId: string, projectIds: string[]): Promise<void> => {
    await delay(1200);
    const template = currentGlobalTemplates.find(t => t.id === templateId);
    if (!template) throw new Error("Template not found");

    projectIds.forEach(pid => {
      const newCopy: Template = {
        ...template,
        id: `tpl_copy_${Date.now()}_${Math.random().toString(36).substring(7)}`,
        projectId: pid,
        name: `${template.name} (Global Copy)`
      };
      currentTemplates.push(newCopy);
    });
  },

  // --- Admin APIs ---

  getSystemHealth: async (): Promise<SystemHealth> => {
    await delay(1000);
    return {
      cpuUsage: 45,
      memoryUsage: 62,
      queueLength: 124,
      apiLatencyMs: 230,
      aadSyncStatus: 'Healthy',
      lastSyncTime: new Date().toISOString()
    };
  },

  getCompany: async (companyId: string): Promise<Company | null> => {
    await delay(600);
    return MOCK_COMPANIES.find(c => c.id === companyId) || null;
  },

  getCompanies: async (): Promise<Company[]> => {
    await delay(600);
    return MOCK_COMPANIES;
  },

  getBots: async (): Promise<BotInstance[]> => {
    await delay(600);
    return currentBots;
  },

  registerBot: async (data: Omit<BotInstance, 'id' | 'status' | 'usageCount'>): Promise<BotInstance> => {
    await delay(1200);
    const newBot: BotInstance = {
      id: `bot_${Date.now()}`,
      ...data,
      status: 'Healthy',
      usageCount: 0
    };
    currentBots.push(newBot);
    return newBot;
  },

  getBillingRecords: async (): Promise<BillingRecord[]> => {
    await delay(800);
    return MOCK_BILLING;
  },

  getTransactionHistory: async (filter: { companyId?: string, projectId?: string, month?: string }): Promise<TransactionRecord[]> => {
    await delay(800);
    let results = MOCK_TRANSACTIONS;

    if (filter.companyId) results = results.filter(t => t.companyId === filter.companyId);
    if (filter.projectId) results = results.filter(t => t.projectId === filter.projectId);
    if (filter.month) results = results.filter(t => t.timestamp.startsWith(filter.month));

    return results;
  },

  performSystemAction: async (action: 'sync_aad' | 'restart_service' | 'block_ip'): Promise<void> => {
    await delay(2000);
    console.log(`[System Action] ${action} executed.`);
  }
};
