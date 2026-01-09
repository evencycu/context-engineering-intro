
import { Project, Template, NotificationLog, MessageStatus, UserProfile, MessagePriority, AudienceList, BillingRecord, Company, BotInstance, TransactionRecord, ADGroup, TeamChannel, ChatGroup } from './types';

export const MOCK_PROJECTS: Project[] = [
  { id: 'proj_fin_01', name: 'Finance Core Systems', apiKey: 'sk_live_fin_123', department: 'Finance', companyId: 'comp_01' },
  { id: 'proj_hr_02', name: 'HR Portal Notifications', apiKey: 'sk_live_hr_456', department: 'Human Resources', companyId: 'comp_01' },
  { id: 'proj_sec_03', name: 'InfoSec Alerts', apiKey: 'sk_live_sec_789', department: 'Security', companyId: 'comp_01' },
  { id: 'proj_nw_01', name: 'Sales Notifications', apiKey: 'sk_live_nw_001', department: 'Sales', companyId: 'comp_02' },
];

export const MOCK_TEMPLATES: Template[] = [
  {
    id: 'tpl_alert',
    projectId: 'proj_fin_01',
    name: 'Critical System Alert',
    description: 'Use for high priority outages or maintenance.',
    variables: ['title', 'severity', 'description', 'eta'],
    defaultJsonStructure: JSON.stringify({
      type: "AdaptiveCard",
      version: "1.4",
      body: [
        { type: "TextBlock", text: "{title}", size: "Large", weight: "Bolder", color: "{severity}" },
        { type: "TextBlock", text: "{description}", wrap: true },
        { type: "FactSet", facts: [{ title: "ETA", value: "{eta}" }] }
      ]
    }, null, 2)
  },
  {
    id: 'tpl_holiday',
    projectId: 'proj_hr_02',
    name: 'Holiday Announcement',
    description: 'General announcement for public holidays.',
    variables: ['holidayName', 'date', 'greeting'],
    defaultJsonStructure: JSON.stringify({
      type: "AdaptiveCard",
      version: "1.4",
      body: [
        { type: "TextBlock", text: "Holiday Notice: {holidayName}", size: "Large", weight: "Bolder" },
        { type: "TextBlock", text: "Date: {date}", isSubtle: true },
        { type: "TextBlock", text: "{greeting}", wrap: true }
      ]
    }, null, 2)
  }
];

export const MOCK_GLOBAL_TEMPLATES: Template[] = [
  {
    id: 'global_tpl_01',
    projectId: 'global',
    name: 'Standard Security Incident',
    description: 'ISO 27001 compliant security incident report format.',
    variables: ['incidentId', 'impactLevel', 'instructions'],
    defaultJsonStructure: JSON.stringify({
      type: "AdaptiveCard",
      version: "1.4",
      msteams: { width: "Full" },
      body: [
        { type: "TextBlock", text: "SECURITY INCIDENT #{incidentId}", size: "Large", weight: "Bolder", color: "Attention" },
        { type: "TextBlock", text: "Impact Level: {impactLevel}", weight: "Bolder" },
        { type: "TextBlock", text: "{instructions}", wrap: true }
      ],
      actions: [
        { type: "Action.OpenUrl", title: "View Incident", url: "https://security.corp/incidents/{incidentId}" }
      ]
    }, null, 2)
  },
  {
    id: 'global_tpl_02',
    projectId: 'global',
    name: 'Company Townhall Invite',
    description: 'Standard invitation card for quarterly townhalls.',
    variables: ['topic', 'speaker', 'meetingLink'],
    defaultJsonStructure: JSON.stringify({
      type: "AdaptiveCard",
      version: "1.4",
      body: [
        { type: "TextBlock", text: "Townhall: {topic}", size: "Medium", weight: "Bolder", color: "Accent" },
        { type: "TextBlock", text: "Speaker: {speaker}", isSubtle: true },
      ],
      actions: [
        { type: "Action.OpenUrl", title: "Join Meeting", url: "{meetingLink}" }
      ]
    }, null, 2)
  }
];

export const MOCK_USERS: UserProfile[] = [
  { id: 'u1', displayName: 'Sarah Connor', userPrincipalName: 'sarah.connor@finance.corp', jobTitle: 'Senior Analyst', tags: ['Finance', 'Manager'] },
  { id: 'u2', displayName: 'John Doe', userPrincipalName: 'john.doe@tech.corp', jobTitle: 'DevOps Engineer', tags: ['IT', 'OnCall'] },
  { id: 'u3', displayName: 'Jane Smith', userPrincipalName: 'jane.smith@hr.corp', jobTitle: 'HR Director', tags: ['HR', 'Executive'] },
  { id: 'u4', displayName: 'Michael Scott', userPrincipalName: 'm.scott@paper.corp', jobTitle: 'Regional Manager', tags: ['Management'] },
  { id: 'u5', displayName: 'Dwight Schrute', userPrincipalName: 'd.schrute@paper.corp', jobTitle: 'Assistant to the Regional Manager', tags: ['Sales', 'Safety'] },
  { id: 'u6', displayName: 'Pam Beesly', userPrincipalName: 'p.beesly@paper.corp', jobTitle: 'Office Administrator', tags: ['Admin'] },
];

export const MOCK_AD_GROUPS: ADGroup[] = [
  { id: 'g1', displayName: 'All Employees (HQ)', mailNickname: 'all-hq', description: 'All staff located in Headquarters', groupTypes: ['Unified'], memberCount: 450 },
  { id: 'g2', displayName: 'IT Administrators', mailNickname: 'it-admins', description: 'Global IT Admin Access', groupTypes: ['Security'], memberCount: 12 },
  { id: 'g3', displayName: 'Finance Dept', mailNickname: 'finance-dept', description: 'Finance Department Members', groupTypes: ['Unified'], memberCount: 35 },
  { id: 'g4', displayName: 'Project Alpha Team', mailNickname: 'proj-alpha', description: 'Cross-functional team for Project Alpha', groupTypes: ['Unified'], memberCount: 8 },
  { id: 'g5', displayName: 'Remote Workers', mailNickname: 'remote-vpn', description: 'VPN Access Group', groupTypes: ['Security'], memberCount: 120 },
];

export const MOCK_CHANNELS: Record<string, TeamChannel[]> = {
  'g1': [
    { id: 'c1', displayName: 'General', membershipType: 'Standard' },
    { id: 'c2', displayName: 'Announcements', membershipType: 'Standard', description: 'Company wide news' },
    { id: 'c3', displayName: 'Water Cooler', membershipType: 'Standard' }
  ],
  'g3': [
    { id: 'c4', displayName: 'General', membershipType: 'Standard' },
    { id: 'c5', displayName: 'Budget Planning', membershipType: 'Private', description: 'Restricted access' },
    { id: 'c6', displayName: 'Invoices', membershipType: 'Standard' }
  ],
  'g4': [
     { id: 'c7', displayName: 'General', membershipType: 'Standard' },
     { id: 'c8', displayName: 'DevOps Alerts', membershipType: 'Standard' }
  ]
};

export const MOCK_AUDIENCE_LISTS: AudienceList[] = [
  { id: 'list_01', projectId: 'proj_hr_02', name: 'All Employees (HQ)', type: 'Static', count: 450, lastUpdated: '2023-10-01' },
  { id: 'list_02', projectId: 'proj_fin_01', name: 'External Contractors', type: 'Static', count: 25, lastUpdated: '2023-10-15' },
];

export const MOCK_CHAT_GROUPS: ChatGroup[] = [
  { id: 'cg_1', projectId: 'proj_fin_01', name: 'Finance Leadership Sync', chatId: '19:82741928471@thread.v2', registeredBy: 'sarah.connor@finance.corp', createdAt: '2023-10-01' },
  { id: 'cg_2', projectId: 'proj_hr_02', name: 'Q4 Event Planning', chatId: '19:19283746551@thread.v2', registeredBy: 'jane.smith@hr.corp', createdAt: '2023-10-05' }
];

export const MOCK_BILLING: BillingRecord[] = [
  { id: 'bill_01', period: '2023-09', companyName: 'Contoso Corp HQ', projectName: 'Finance Core Systems', messageCount: 15400, cost: 154.00, currency: 'USD' },
  { id: 'bill_02', period: '2023-09', companyName: 'Contoso Corp HQ', projectName: 'HR Portal Notifications', messageCount: 4500, cost: 45.00, currency: 'USD' },
  { id: 'bill_03', period: '2023-10', companyName: 'Contoso Corp HQ', projectName: 'Finance Core Systems', messageCount: 8200, cost: 82.00, currency: 'USD' },
  { id: 'bill_04', period: '2023-10', companyName: 'Northwind Traders', projectName: 'Sales Notifications', messageCount: 1200, cost: 12.00, currency: 'USD' },
];

export const MOCK_COMPANIES: Company[] = [
  { id: 'comp_01', name: 'Contoso Corp HQ', adminEmail: 'admin@contoso.com' },
  { id: 'comp_02', name: 'Northwind Traders', adminEmail: 'it@northwind.com' },
];

export const MOCK_BOTS: BotInstance[] = [
  { id: 'bot_01', displayName: 'Contoso Assistant', appId: 'a1b2c3d4-e5f6-7890-a1b2-c3d4e5f67890', tenantId: 'tenant-guid-1', status: 'Healthy', usageCount: 15420 },
  { id: 'bot_02', displayName: 'HR Helper', appId: 'b2c3d4e5-f6a1-8901-b2c3-d4e5f6a18901', tenantId: 'tenant-guid-1', status: 'Degraded', usageCount: 320 },
];

export const MOCK_TRANSACTIONS: TransactionRecord[] = [
  { id: 'tx_001', timestamp: '2023-10-25T10:00:00Z', companyId: 'comp_01', companyName: 'Contoso Corp HQ', projectId: 'proj_fin_01', projectName: 'Finance Core Systems', type: 'Message_Sent', quantity: 5000, amount: 50.00, currency: 'USD', status: 'Paid' },
  { id: 'tx_002', timestamp: '2023-10-25T14:30:00Z', companyId: 'comp_01', companyName: 'Contoso Corp HQ', projectId: 'proj_hr_02', projectName: 'HR Portal Notifications', type: 'Message_Sent', quantity: 200, amount: 2.00, currency: 'USD', status: 'Paid' },
  { id: 'tx_003', timestamp: '2023-10-26T09:15:00Z', companyId: 'comp_01', companyName: 'Contoso Corp HQ', projectId: 'proj_fin_01', projectName: 'Finance Core Systems', type: 'Bot_Maintenance', quantity: 1, amount: 25.00, currency: 'USD', status: 'Paid' },
  { id: 'tx_004', timestamp: '2023-10-26T11:00:00Z', companyId: 'comp_02', companyName: 'Northwind Traders', projectId: 'proj_nw_01', projectName: 'Sales Notifications', type: 'Message_Sent', quantity: 1200, amount: 12.00, currency: 'USD', status: 'Pending' },
  { id: 'tx_005', timestamp: '2023-10-27T08:00:00Z', companyId: 'comp_01', companyName: 'Contoso Corp HQ', projectId: 'proj_hr_02', projectName: 'HR Portal Notifications', type: 'Subscription_Fee', quantity: 1, amount: 100.00, currency: 'USD', status: 'Paid' },
];

export const generateMockLogs = (): NotificationLog[] => {
  return Array.from({ length: 15 }).map((_, i) => ({
    id: `log_${i}`,
    projectId: i % 2 === 0 ? 'proj_fin_01' : 'proj_hr_02',
    timestamp: new Date(Date.now() - i * 3600000).toISOString(),
    senderUpn: 'admin@corp.com',
    recipient: i % 3 === 0 ? 'General Channel' : 'sarah.connor@finance.corp',
    recipientType: i % 3 === 0 ? 'Channel' : 'User',
    templateName: i % 4 === 0 ? 'Custom JSON' : 'Critical System Alert',
    status: i === 0 ? MessageStatus.SCHEDULED : (i === 2 ? MessageStatus.FAILED : MessageStatus.SENT),
    priority: i % 5 === 0 ? MessagePriority.HIGH : MessagePriority.NORMAL,
    errorMessage: i === 2 ? 'Gateway Timeout (504)' : undefined,
    contentPreview: '{"type": "AdaptiveCard", ...}',
    scheduledFor: i === 0 ? new Date(Date.now() + 86400000).toISOString() : undefined
  }));
};
