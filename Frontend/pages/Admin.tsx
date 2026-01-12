
import React, { useEffect, useState } from 'react';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { service } from '../services';
import { SystemHealth, BillingRecord, Company, BotInstance, Project, TransactionRecord, UserProfile, ADGroup, TeamChannel, Template } from '../types';
import { Plus, RefreshCw, Server, ShieldAlert, Filter, Download, Building2, Briefcase, Users, Search, ChevronDown, ChevronRight, Hash, Lock, Globe, LayoutTemplate, Share2, Trash2, Edit2, Code, Variable } from 'lucide-react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

interface AdminProps {
  view: 'health' | 'organization' | 'aad' | 'bots' | 'billing' | 'ops' | 'templates';
}

export const Admin: React.FC<AdminProps> = ({ view }) => {
  const [health, setHealth] = useState<SystemHealth | null>(null);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [bots, setBots] = useState<BotInstance[]>([]);
  const [billing, setBilling] = useState<BillingRecord[]>([]);
  const [transactions, setTransactions] = useState<TransactionRecord[]>([]);
  const [globalTemplates, setGlobalTemplates] = useState<Template[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // Organization View State
  const [orgTab, setOrgTab] = useState<'companies' | 'projects'>('companies');

  // Directory View State
  const [dirTab, setDirTab] = useState<'users' | 'groups'>('users');
  const [directoryUsers, setDirectoryUsers] = useState<UserProfile[]>([]);
  const [directoryGroups, setDirectoryGroups] = useState<ADGroup[]>([]);
  const [dirSearch, setDirSearch] = useState('');

  // Channel Expansion State
  const [expandedGroupId, setExpandedGroupId] = useState<string | null>(null);
  const [groupChannels, setGroupChannels] = useState<TeamChannel[]>([]);
  const [loadingChannels, setLoadingChannels] = useState(false);

  // Billing Filters
  const [filterCompany, setFilterCompany] = useState('');
  const [filterProject, setFilterProject] = useState('');
  const [filterMonth, setFilterMonth] = useState('');

  // Bot Registration Modal State
  const [showBotModal, setShowBotModal] = useState(false);
  const [botForm, setBotForm] = useState({ displayName: '', appId: '', tenantId: '', clientSecret: '' });
  const [isRegistering, setIsRegistering] = useState(false);

  // Template Distribution Modal State
  const [showDistributeModal, setShowDistributeModal] = useState(false);
  const [selectedTemplateForDist, setSelectedTemplateForDist] = useState<Template | null>(null);
  const [selectedProjectsForDist, setSelectedProjectsForDist] = useState<string[]>([]);
  const [isDistributing, setIsDistributing] = useState(false);

  // Global Template Modal State
  const [showGlobalTemplateModal, setShowGlobalTemplateModal] = useState(false);
  const [editingGlobalTemplate, setEditingGlobalTemplate] = useState<Template | null>(null);
  const [globalTemplateForm, setGlobalTemplateForm] = useState<Partial<Template>>({
    name: '',
    description: '',
    variables: [],
    defaultJsonStructure: '{}'
  });
  const [globalTemplateVariableInput, setGlobalTemplateVariableInput] = useState('');
  const [isSavingGlobalTemplate, setIsSavingGlobalTemplate] = useState(false);

  // Company/Project Creation Modal State
  const [showCompanyModal, setShowCompanyModal] = useState(false);
  const [showProjectModal, setShowProjectModal] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [editingCompany, setEditingCompany] = useState<Company | null>(null);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [companyForm, setCompanyForm] = useState({
    name: '',
    contactEmail: '',
    contactPhone: '',
    address: '',
    billingEnabled: false,
  });
  const [projectForm, setProjectForm] = useState({
    companyId: '',
    notifyKey: '',
    description: '',
    dailyLimit: 10000,
    monthlyLimit: 300000,
    priority: 'normal' as 'low' | 'normal' | 'high',
  });

  useEffect(() => {
    loadData();
    // Reset search on view change
    setDirSearch('');
  }, [view]);

  // Load Transactions when filters change (if in billing view)
  useEffect(() => {
    if (view === 'billing') {
      loadTransactions();
    }
  }, [filterCompany, filterProject, filterMonth]);

  const loadData = async () => {
    setIsLoading(true);
    try {
      if (view === 'health') {
        const h = await service.getSystemHealth();
        setHealth(h);
      } else if (view === 'organization') {
        const [c, p] = await Promise.all([
          service.getCompanies(),
          service.getProjects()
        ]);
        setCompanies(c);
        setProjects(p);
      } else if (view === 'aad') {
        const [u, g] = await Promise.all([
          service.getAllUsers(),
          service.getAllGroups()
        ]);
        setDirectoryUsers(u);
        setDirectoryGroups(g);
      } else if (view === 'bots') {
        const b = await service.getBots();
        setBots(b);
      } else if (view === 'billing') {
        const [bil, c, p] = await Promise.all([
          service.getBillingRecords(),
          service.getCompanies(),
          service.getProjects()
        ]);
        setBilling(bil);
        setCompanies(c);
        setProjects(p);
        await loadTransactions();
      } else if (view === 'templates') {
        const [t, p] = await Promise.all([
          service.getGlobalTemplates(),
          service.getProjects()
        ]);
        setGlobalTemplates(t);
        setProjects(p);
      }
    } catch (e) {
      console.error(e);
    } finally {
      setIsLoading(false);
    }
  };

  const loadTransactions = async () => {
    const txs = await service.getTransactionHistory({
      companyId: filterCompany || undefined,
      projectId: filterProject || undefined,
      month: filterMonth || undefined
    });
    setTransactions(txs);
  };

  const handleSystemAction = async (action: 'sync_aad' | 'restart_service') => {
    if (!window.confirm(`Are you sure you want to ${action.replace('_', ' ').toUpperCase()}?`)) return;
    await service.performSystemAction(action);
    alert('Action executed successfully.');
  };

  const handleRegisterBot = async () => {
    setIsRegistering(true);
    try {
      await service.registerBot({
        displayName: botForm.displayName,
        appId: botForm.appId,
        tenantId: botForm.tenantId
      });
      setShowBotModal(false);
      setBotForm({ displayName: '', appId: '', tenantId: '', clientSecret: '' });
      loadData(); // Refresh list
    } catch (e) {
      alert('Failed to register bot');
    } finally {
      setIsRegistering(false);
    }
  };

  const handleExpandGroup = async (groupId: string) => {
    if (expandedGroupId === groupId) {
      setExpandedGroupId(null);
      setGroupChannels([]);
      return;
    }

    setExpandedGroupId(groupId);
    setLoadingChannels(true);
    try {
      const channels = await service.getGroupChannels(groupId);
      setGroupChannels(channels);
    } catch (e) {
      console.error("Failed to load channels", e);
    } finally {
      setLoadingChannels(false);
    }
  };

  const handleSyncChannels = async (groupId: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!window.confirm("Sync channels for this group? This may take a few seconds.")) return;
    try {
      await service.syncChannels(groupId);
      const channels = await service.getGroupChannels(groupId);
      if (expandedGroupId === groupId) {
        setGroupChannels(channels);
      }
      alert('Channels synced successfully.');
    } catch (e) {
      alert('Sync failed.');
    }
  };

  const openDistributeModal = (template: Template) => {
    setSelectedTemplateForDist(template);
    setSelectedProjectsForDist([]);
    setShowDistributeModal(true);
  };

  const handleDistribute = async () => {
    if (!selectedTemplateForDist || selectedProjectsForDist.length === 0) return;
    setIsDistributing(true);
    try {
      await service.assignGlobalTemplateToProjects(selectedTemplateForDist.id, selectedProjectsForDist);
      alert(`Template distributed to ${selectedProjectsForDist.length} projects successfully.`);
      setShowDistributeModal(false);
    } catch (e) {
      alert('Failed to distribute template.');
    } finally {
      setIsDistributing(false);
    }
  };

  // Global Template CRUD handlers
  const handleOpenCreateGlobalTemplate = () => {
    setEditingGlobalTemplate(null);
    setGlobalTemplateForm({
      name: '',
      description: '',
      variables: [],
      defaultJsonStructure: JSON.stringify({
        type: "AdaptiveCard",
        version: "1.4",
        body: [{ type: "TextBlock", text: "New Global Template" }]
      }, null, 2)
    });
    setGlobalTemplateVariableInput('');
    setShowGlobalTemplateModal(true);
  };

  const handleOpenEditGlobalTemplate = (t: Template) => {
    setEditingGlobalTemplate(t);
    setGlobalTemplateForm({
      name: t.name,
      description: t.description,
      variables: [...t.variables],
      defaultJsonStructure: t.defaultJsonStructure
    });
    setGlobalTemplateVariableInput('');
    setShowGlobalTemplateModal(true);
  };

  const handleSaveGlobalTemplate = async () => {
    if (!globalTemplateForm.name || !globalTemplateForm.defaultJsonStructure) {
      alert('Please fill in template name and JSON structure');
      return;
    }

    // Validate JSON structure
    try {
      JSON.parse(globalTemplateForm.defaultJsonStructure!);
    } catch (e) {
      alert('Invalid JSON structure. Please check your JSON syntax.');
      return;
    }

    setIsSavingGlobalTemplate(true);
    try {
      if (editingGlobalTemplate) {
        await service.updateGlobalTemplate(editingGlobalTemplate.id, globalTemplateForm);
      } else {
        await service.createGlobalTemplate({
          projectId: 'global', // Global templates don't belong to a project
          name: globalTemplateForm.name!,
          description: globalTemplateForm.description || '',
          variables: globalTemplateForm.variables || [],
          defaultJsonStructure: globalTemplateForm.defaultJsonStructure!
        });
      }
      setShowGlobalTemplateModal(false);
      // Reload global templates
      const templates = await service.getGlobalTemplates();
      setGlobalTemplates(templates);
    } catch (e: any) {
      console.error('Failed to save global template:', e);
      const errorMessage = e?.message || e?.toString() || 'Failed to save global template';
      alert(`Failed to save global template: ${errorMessage}`);
    } finally {
      setIsSavingGlobalTemplate(false);
    }
  };

  const handleDeleteGlobalTemplate = async (id: string) => {
    if (!window.confirm("Are you sure you want to delete this global template?")) return;
    try {
      await service.deleteGlobalTemplate(id);
      // Reload global templates
      const templates = await service.getGlobalTemplates();
      setGlobalTemplates(templates);
    } catch (e: any) {
      console.error('Failed to delete global template:', e);
      const errorMessage = e?.message || e?.toString() || 'Failed to delete global template';
      alert(`Failed to delete global template: ${errorMessage}`);
    }
  };

  const addGlobalTemplateVariable = () => {
    if (globalTemplateVariableInput && !globalTemplateForm.variables?.includes(globalTemplateVariableInput)) {
      setGlobalTemplateForm({
        ...globalTemplateForm,
        variables: [...(globalTemplateForm.variables || []), globalTemplateVariableInput]
      });
      setGlobalTemplateVariableInput('');
    }
  };

  const removeGlobalTemplateVariable = (v: string) => {
    setGlobalTemplateForm({
      ...globalTemplateForm,
      variables: globalTemplateForm.variables?.filter(item => item !== v)
    });
  };

  const handleOpenEditCompany = async (company: Company) => {
    setEditingCompany(company);
    try {
      // Fetch full company data for editing
      const fullCompany = await service.getCompany(company.id);
      if (fullCompany) {
        setCompanyForm({
          name: fullCompany.name,
          contactEmail: fullCompany.adminEmail,
          contactPhone: fullCompany.contactPhone || '',
          address: fullCompany.address || '',
          billingEnabled: fullCompany.billingEnabled || false,
        });
      } else {
        // Fallback to basic data
        setCompanyForm({
          name: company.name,
          contactEmail: company.adminEmail,
          contactPhone: company.contactPhone || '',
          address: company.address || '',
          billingEnabled: company.billingEnabled || false,
        });
      }
    } catch (e) {
      // Fallback to basic data on error
      setCompanyForm({
        name: company.name,
        contactEmail: company.adminEmail,
        contactPhone: company.contactPhone || '',
        address: company.address || '',
        billingEnabled: company.billingEnabled || false,
      });
    }
    setShowCompanyModal(true);
  };

  const handleCreateCompany = async () => {
    if (!companyForm.name || !companyForm.contactEmail || !companyForm.contactPhone || !companyForm.address) {
      alert('Please fill in all required fields.');
      return;
    }
    setIsCreating(true);
    try {
      if (editingCompany) {
        // Update existing company
        await service.updateCompany(editingCompany.id, {
          name: companyForm.name,
          adminEmail: companyForm.contactEmail,
          contactPhone: companyForm.contactPhone,
          address: companyForm.address,
          billingEnabled: companyForm.billingEnabled,
        });
      } else {
        // Create new company
        await service.createCompany(companyForm);
      }
      setShowCompanyModal(false);
      setEditingCompany(null);
      setCompanyForm({ name: '', contactEmail: '', contactPhone: '', address: '', billingEnabled: false });
      loadData(); // Refresh list
    } catch (e: any) {
      alert(e.message || `Failed to ${editingCompany ? 'update' : 'create'} company.`);
    } finally {
      setIsCreating(false);
    }
  };

  const handleOpenEditProject = async (project: Project) => {
    setEditingProject(project);
    try {
      // Fetch full project data for editing (to get dailyLimit, monthlyLimit, etc.)
      const fullProject = await service.getProject(project.id);
      if (fullProject) {
        // Map frontend priority to form priority
        const formPriority = fullProject.priority === 'Low (Broadcast)' ? 'low' :
          fullProject.priority === 'High (Alert)' ? 'high' : 'normal';
        setProjectForm({
          companyId: fullProject.companyId,
          notifyKey: fullProject.name, // Project name is the notify_key
          description: fullProject.department, // Department is the description
          dailyLimit: fullProject.dailyLimit || 10000,
          monthlyLimit: fullProject.monthlyLimit || 300000,
          priority: formPriority as 'low' | 'normal' | 'high',
        });
      } else {
        // Fallback to basic data
        const formPriority = project.priority === 'Low (Broadcast)' ? 'low' :
          project.priority === 'High (Alert)' ? 'high' : 'normal';
        setProjectForm({
          companyId: project.companyId,
          notifyKey: project.name,
          description: project.department,
          dailyLimit: project.dailyLimit || 10000,
          monthlyLimit: project.monthlyLimit || 300000,
          priority: formPriority as 'low' | 'normal' | 'high',
        });
      }
    } catch (e) {
      // Fallback to basic data on error
      const formPriority = project.priority === 'Low (Broadcast)' ? 'low' :
        project.priority === 'High (Alert)' ? 'high' : 'normal';
      setProjectForm({
        companyId: project.companyId,
        notifyKey: project.name,
        description: project.department,
        dailyLimit: project.dailyLimit || 10000,
        monthlyLimit: project.monthlyLimit || 300000,
        priority: formPriority as 'low' | 'normal' | 'high',
      });
    }
    setShowProjectModal(true);
  };

  const handleCreateProject = async () => {
    if (!projectForm.companyId || !projectForm.notifyKey || !projectForm.description) {
      alert('Please fill in all required fields.');
      return;
    }
    setIsCreating(true);
    try {
      if (editingProject) {
        // Update existing project
        // Map form priority to frontend priority format
        const frontendPriority = projectForm.priority === 'low' ? 'Low (Broadcast)' :
          projectForm.priority === 'high' ? 'High (Alert)' : 'Normal';
        await service.updateProject(editingProject.id, {
          name: projectForm.notifyKey,
          department: projectForm.description,
          priority: frontendPriority,
        });
      } else {
        // Create new project
        // Generate a valid UUID for createdBy (in real app, this would come from auth context)
        // Using a simple UUID v4 format
        const createdBy = '550e8400-e29b-41d4-a716-446655440000'; // Valid UUID format
        await service.createProject({
          companyId: projectForm.companyId,
          notifyKey: projectForm.notifyKey,
          description: projectForm.description,
          dailyLimit: projectForm.dailyLimit,
          monthlyLimit: projectForm.monthlyLimit,
          priority: projectForm.priority,
          createdBy,
        });
      }
      setShowProjectModal(false);
      setEditingProject(null);
      setProjectForm({
        companyId: '',
        notifyKey: '',
        description: '',
        dailyLimit: 10000,
        monthlyLimit: 300000,
        priority: 'normal',
      });
      loadData(); // Refresh list
    } catch (e: any) {
      alert(e.message || `Failed to ${editingProject ? 'update' : 'create'} project.`);
    } finally {
      setIsCreating(false);
    }
  };

  const getCompanyName = (companyId: string) => {
    return companies.find(c => c.id === companyId)?.name || 'Unknown';
  };

  // Filter logic for Directory View
  const filteredUsers = directoryUsers.filter(u =>
    u.displayName.toLowerCase().includes(dirSearch.toLowerCase()) ||
    u.userPrincipalName.toLowerCase().includes(dirSearch.toLowerCase()) ||
    u.jobTitle.toLowerCase().includes(dirSearch.toLowerCase())
  );

  const filteredGroups = directoryGroups.filter(g =>
    g.displayName.toLowerCase().includes(dirSearch.toLowerCase()) ||
    g.mailNickname.toLowerCase().includes(dirSearch.toLowerCase())
  );

  // Filter projects for the billing dropdown based on selected company
  const filteredProjects = filterCompany
    ? projects.filter(p => p.companyId === filterCompany)
    : projects;

  const titles: Record<string, { title: string, subtitle: string }> = {
    health: { title: 'System Health', subtitle: 'Real-time infrastructure monitoring' },
    organization: { title: 'Organization Management', subtitle: 'Manage subsidiaries and project workspaces' },
    aad: { title: 'Directory Management', subtitle: 'Manage Azure AD Users and Groups' },
    bots: { title: 'Bot Registry', subtitle: 'Manage Teams Bot registrations' },
    billing: { title: 'Billing Center', subtitle: 'Cost analysis and transaction history' },
    ops: { title: 'System Operations', subtitle: 'Emergency controls and maintenance' },
    templates: { title: 'Global Templates', subtitle: 'Manage and distribute standardized templates' },
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between border-b border-gray-200 pb-5">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">{titles[view].title}</h2>
          <p className="text-gray-500 mt-1">{titles[view].subtitle}</p>
        </div>
        <div className="bg-indigo-100 text-indigo-800 px-3 py-1 rounded-full text-xs font-bold flex items-center gap-1">
          <ShieldAlert size={14} /> GLOBAL ADMIN
        </div>
      </div>

      <div className="min-h-[400px]">
        {/* === HEALTH === */}
        {view === 'health' && health && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
              <p className="text-sm font-medium text-gray-500">CPU Usage</p>
              <div className="mt-2 flex items-baseline gap-2">
                <span className="text-3xl font-bold text-gray-900">{health.cpuUsage}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2.5 mt-4">
                <div className="bg-blue-600 h-2.5 rounded-full" style={{ width: `${health.cpuUsage}%` }}></div>
              </div>
            </div>

            <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
              <p className="text-sm font-medium text-gray-500">Queue Length</p>
              <div className="mt-2 flex items-baseline gap-2">
                <span className="text-3xl font-bold text-gray-900">{health.queueLength}</span>
                <span className="text-sm text-gray-500">msgs</span>
              </div>
            </div>

            <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
              <p className="text-sm font-medium text-gray-500">API Latency</p>
              <div className="mt-2 flex items-baseline gap-2">
                <span className={`text-3xl font-bold ${health.apiLatencyMs > 500 ? 'text-red-600' : 'text-green-600'}`}>
                  {health.apiLatencyMs}
                </span>
                <span className="text-sm text-gray-500">ms</span>
              </div>
            </div>

            <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
              <p className="text-sm font-medium text-gray-500">AAD Sync Status</p>
              <div className="mt-2 flex items-center gap-2">
                <div className={`h-3 w-3 rounded-full ${health.aadSyncStatus === 'Healthy' ? 'bg-green-500' : 'bg-red-500'}`}></div>
                <span className="text-xl font-bold text-gray-900">{health.aadSyncStatus}</span>
              </div>
              <p className="text-xs text-gray-400 mt-2">Last: {new Date(health.lastSyncTime).toLocaleTimeString()}</p>
            </div>
          </div>
        )}

        {/* === GLOBAL TEMPLATES === */}
        {view === 'templates' && (
          <div className="space-y-6">
            <div className="flex justify-end">
              <Button onClick={handleOpenCreateGlobalTemplate}>
                <Plus size={16} className="mr-2" />
                New Global Template
              </Button>
            </div>

            <div className="bg-white shadow-sm border border-gray-200 rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Template Name</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Variables</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {globalTemplates.map((tpl) => (
                    <tr key={tpl.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{tpl.name}</td>
                      <td className="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">{tpl.description}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        <div className="flex gap-1 flex-wrap">
                          {tpl.variables.map(v => (
                            <span key={v} className="px-2 py-0.5 bg-gray-100 text-xs rounded border border-gray-200">{v}</span>
                          ))}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            onClick={() => handleOpenEditGlobalTemplate(tpl)}
                            className="text-gray-400 hover:text-blue-600 p-1"
                            title="Edit Template"
                          >
                            <Edit2 size={16} />
                          </button>
                          <button
                            onClick={() => handleDeleteGlobalTemplate(tpl.id)}
                            className="text-gray-400 hover:text-red-600 p-1"
                            title="Delete Template"
                          >
                            <Trash2 size={16} />
                          </button>
                          <button
                            onClick={() => openDistributeModal(tpl)}
                            className="text-indigo-600 hover:text-indigo-900 flex items-center gap-1"
                          >
                            <Share2 size={16} /> Distribute
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* === ORGANIZATION === */}
        {view === 'organization' && (
          <div className="space-y-6">
            {/* Org Sub-Navigation */}
            <div className="flex space-x-4 border-b border-gray-200">
              <button
                onClick={() => setOrgTab('companies')}
                className={`flex items-center gap-2 py-2 px-4 border-b-2 font-medium text-sm transition-colors ${orgTab === 'companies'
                  ? 'border-indigo-600 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
                  }`}
              >
                <Building2 size={16} />
                Companies
              </button>
              <button
                onClick={() => setOrgTab('projects')}
                className={`flex items-center gap-2 py-2 px-4 border-b-2 font-medium text-sm transition-colors ${orgTab === 'projects'
                  ? 'border-indigo-600 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
                  }`}
              >
                <Briefcase size={16} />
                Projects
              </button>
            </div>

            {orgTab === 'companies' && (
              <div className="bg-white shadow-sm border border-gray-200 rounded-lg overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center bg-gray-50">
                  <h3 className="text-sm font-semibold text-gray-700 uppercase">Registered Subsidiaries</h3>
                  <Button variant="secondary" onClick={() => setShowCompanyModal(true)}>
                    <Plus size={14} className="mr-1" /> Add Company
                  </Button>
                </div>
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Company ID</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Company Name</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Admin Email</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {companies.map((company) => (
                      <tr key={company.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-xs text-gray-500 font-mono">{company.id}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{company.name}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{company.adminEmail}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <button
                            onClick={() => handleOpenEditCompany(company)}
                            className="text-indigo-600 hover:text-indigo-900"
                          >
                            Edit
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            {orgTab === 'projects' && (
              <div className="bg-white shadow-sm border border-gray-200 rounded-lg overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center bg-gray-50">
                  <h3 className="text-sm font-semibold text-gray-700 uppercase">Active Workspaces</h3>
                  <Button variant="secondary" onClick={() => setShowProjectModal(true)}>
                    <Plus size={14} className="mr-1" /> Add Project
                  </Button>
                </div>
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Project Name</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Department</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Belongs To Company</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Project ID</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {projects.map((project) => (
                      <tr key={project.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{project.name}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{project.department}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                            <Building2 size={12} className="mr-1" />
                            {getCompanyName(project.companyId)}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-xs text-gray-400 font-mono">{project.id}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <button
                            onClick={() => handleOpenEditProject(project)}
                            className="text-indigo-600 hover:text-indigo-900"
                          >
                            Edit
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* === DIRECTORY MANAGEMENT (AAD) === */}
        {view === 'aad' && (
          <div className="space-y-6">
            {/* Tabs & Search */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-gray-200 pb-2">
              <div className="flex space-x-6">
                <button
                  onClick={() => setDirTab('users')}
                  className={`py-2 px-1 border-b-2 font-medium text-sm transition-colors ${dirTab === 'users'
                    ? 'border-indigo-600 text-indigo-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                    }`}
                >
                  Users
                </button>
                <button
                  onClick={() => setDirTab('groups')}
                  className={`py-2 px-1 border-b-2 font-medium text-sm transition-colors ${dirTab === 'groups'
                    ? 'border-indigo-600 text-indigo-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                    }`}
                >
                  Groups
                </button>
              </div>

              <div className="relative w-full sm:w-64">
                <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  type="text"
                  placeholder={`Search ${dirTab}...`}
                  className="w-full pl-9 pr-4 py-2 text-sm border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
                  value={dirSearch}
                  onChange={(e) => setDirSearch(e.target.value)}
                />
              </div>
            </div>

            {/* Content */}
            {dirTab === 'users' && (
              <div className="bg-white shadow-sm border border-gray-200 rounded-lg overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Display Name</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User Principal Name</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Job Title</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tags</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {filteredUsers.length === 0 ? (
                      <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-500">No users found.</td></tr>
                    ) : (
                      filteredUsers.map((user) => (
                        <tr key={user.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{user.displayName}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{user.userPrincipalName}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{user.jobTitle}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            <div className="flex flex-wrap gap-1">
                              {user.tags?.map(t => (
                                <span key={t} className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">{t}</span>
                              ))}
                            </div>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            )}

            {dirTab === 'groups' && (
              <div className="bg-white shadow-sm border border-gray-200 rounded-lg overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-3 w-8"></th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Group Name</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Mail Nickname</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Members</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {filteredGroups.length === 0 ? (
                      <tr><td colSpan={6} className="px-6 py-8 text-center text-gray-500">No groups found.</td></tr>
                    ) : (
                      filteredGroups.map((group) => (
                        <React.Fragment key={group.id}>
                          <tr className={`hover:bg-gray-50 cursor-pointer ${expandedGroupId === group.id ? 'bg-indigo-50 hover:bg-indigo-50' : ''}`} onClick={() => group.groupTypes.includes('Unified') && handleExpandGroup(group.id)}>
                            <td className="px-4 py-4 text-gray-400">
                              {group.groupTypes.includes('Unified') && (
                                expandedGroupId === group.id ? <ChevronDown size={16} /> : <ChevronRight size={16} />
                              )}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{group.displayName}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{group.mailNickname}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {group.groupTypes.includes('Unified') ? (
                                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">M365</span>
                              ) : (
                                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">Security</span>
                              )}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-500">{group.memberCount}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                              {group.groupTypes.includes('Unified') && (
                                <button onClick={(e) => handleSyncChannels(group.id, e)} className="text-indigo-600 hover:text-indigo-900 text-xs font-medium flex items-center justify-end gap-1 ml-auto">
                                  <RefreshCw size={12} /> Sync Channels
                                </button>
                              )}
                            </td>
                          </tr>
                          {/* Nested Channels Row */}
                          {expandedGroupId === group.id && (
                            <tr className="bg-gray-50">
                              <td colSpan={6} className="px-4 py-4 sm:px-10">
                                <div className="bg-white border border-gray-200 rounded-md p-4 shadow-inner">
                                  <h4 className="text-xs font-bold text-gray-500 uppercase tracking-wide mb-3 flex items-center gap-2">
                                    <Hash size={12} /> Channels
                                  </h4>

                                  {loadingChannels ? (
                                    <div className="text-sm text-gray-500 py-2 italic">Loading channels from Microsoft Graph...</div>
                                  ) : groupChannels.length === 0 ? (
                                    <div className="text-sm text-gray-500 py-2">No channels found or synced yet.</div>
                                  ) : (
                                    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
                                      {groupChannels.map(channel => (
                                        <div key={channel.id} className="flex items-center justify-between p-3 border border-gray-100 rounded bg-gray-50 hover:bg-white hover:border-gray-300 transition-all">
                                          <div className="flex items-center gap-2">
                                            {channel.membershipType === 'Private' ? <Lock size={14} className="text-amber-500" /> : <Globe size={14} className="text-gray-400" />}
                                            <span className="text-sm font-medium text-gray-700">{channel.displayName}</span>
                                          </div>
                                          <span className="text-xs text-gray-400">{channel.membershipType}</span>
                                        </div>
                                      ))}
                                    </div>
                                  )}
                                </div>
                              </td>
                            </tr>
                          )}
                        </React.Fragment>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* === BOTS === */}
        {view === 'bots' && (
          <div className="space-y-4">
            <div className="flex justify-end">
              <Button onClick={() => setShowBotModal(true)}>
                <Plus size={16} className="mr-2" />
                Register New Bot
              </Button>
            </div>

            <div className="bg-white shadow-sm border border-gray-200 rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Display Name</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">App ID</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tenant ID</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Usage</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {bots.map((bot) => (
                    <tr key={bot.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{bot.displayName}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-xs font-mono text-gray-500">{bot.appId}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-xs font-mono text-gray-500">{bot.tenantId}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm">
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${bot.status === 'Healthy' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                          }`}>
                          {bot.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-500">{bot.usageCount}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* === BILLING === */}
        {view === 'billing' && (
          <div className="space-y-8">
            {/* Chart Section */}
            <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
              <h3 className="text-lg font-medium text-gray-900 mb-6">Monthly Cost Overview</h3>
              <div className="h-80 w-full">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={billing}>
                    <CartesianGrid strokeDasharray="3 3" vertical={false} />
                    <XAxis dataKey="projectName" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="cost" fill="#3b82f6" name="Cost (USD)" radius={[4, 4, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>

            {/* Transaction History Section */}
            <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
              <div className="p-6 border-b border-gray-200 flex flex-col md:flex-row md:items-center justify-between gap-4">
                <h3 className="text-lg font-medium text-gray-900">Transaction History</h3>

                {/* Filters */}
                <div className="flex flex-col md:flex-row gap-2">
                  <select
                    className="form-select text-sm border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                    value={filterCompany}
                    onChange={(e) => { setFilterCompany(e.target.value); setFilterProject(''); }}
                  >
                    <option value="">All Companies</option>
                    {companies.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                  </select>

                  <select
                    className="form-select text-sm border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                    value={filterProject}
                    onChange={(e) => setFilterProject(e.target.value)}
                  >
                    <option value="">All Projects</option>
                    {filteredProjects.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
                  </select>

                  <input
                    type="month"
                    className="form-input text-sm border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                    value={filterMonth}
                    onChange={(e) => setFilterMonth(e.target.value)}
                  />

                  <Button variant="secondary" onClick={() => loadTransactions()}>
                    <Filter size={16} className="mr-2" /> Apply
                  </Button>
                </div>
              </div>

              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Company</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Project</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Qty</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                      <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {transactions.length === 0 ? (
                      <tr><td colSpan={7} className="px-6 py-8 text-center text-gray-500">No transactions found matching your filters.</td></tr>
                    ) : (
                      transactions.map((tx) => (
                        <tr key={tx.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {new Date(tx.timestamp).toLocaleDateString()} <span className="text-xs text-gray-400">{new Date(tx.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{tx.companyName}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{tx.projectName}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                              {tx.type.replace('_', ' ')}
                            </span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-500 font-mono">{tx.quantity}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium text-gray-900">
                            {tx.amount.toFixed(2)} <span className="text-gray-400 text-xs">{tx.currency}</span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-center">
                            <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${tx.status === 'Paid' ? 'bg-green-100 text-green-800' :
                              tx.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' : 'bg-red-100 text-red-800'
                              }`}>
                              {tx.status}
                            </span>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}

        {/* === OPS === */}
        {view === 'ops' && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Directory Synchronization</h3>
              <p className="text-sm text-gray-500 mb-6">
                Manually trigger an incremental sync with Azure AD. This is usually scheduled automatically every 4 hours.
              </p>
              <Button onClick={() => handleSystemAction('sync_aad')}>
                <RefreshCw size={16} className="mr-2" />
                Trigger Full Sync
              </Button>
            </div>

            <div className="bg-red-50 p-6 rounded-xl border border-red-200 shadow-sm">
              <h3 className="text-lg font-medium text-red-900 mb-4">Emergency Controls</h3>
              <p className="text-sm text-red-700 mb-6">
                Restarting services will cause a brief downtime of 10-30 seconds. All active WebSocket connections will be dropped.
              </p>
              <Button variant="danger" onClick={() => handleSystemAction('restart_service')}>
                <Server size={16} className="mr-2" />
                Restart Notification Service
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Bot Registration Modal */}
      {showBotModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowBotModal(false)}></div>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <div className="sm:flex sm:items-start">
                  <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                    <h3 className="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                      Register New Bot
                    </h3>
                    <div className="mt-4 space-y-4">
                      <Input
                        label="Display Name"
                        placeholder="e.g. Finance Support Bot"
                        value={botForm.displayName}
                        onChange={(e) => setBotForm({ ...botForm, displayName: e.target.value })}
                      />
                      <Input
                        label="Microsoft App ID (GUID)"
                        placeholder="00000000-0000-0000-0000-000000000000"
                        value={botForm.appId}
                        onChange={(e) => setBotForm({ ...botForm, appId: e.target.value })}
                      />
                      <Input
                        label="Azure Tenant ID (GUID)"
                        placeholder="00000000-0000-0000-0000-000000000000"
                        value={botForm.tenantId}
                        onChange={(e) => setBotForm({ ...botForm, tenantId: e.target.value })}
                      />
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Client Secret</label>
                        <input
                          type="password"
                          className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                          placeholder="••••••••••••••••"
                          value={botForm.clientSecret}
                          onChange={(e) => setBotForm({ ...botForm, clientSecret: e.target.value })}
                        />
                        <p className="text-xs text-gray-500 mt-1">This secret is never stored in plain text.</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <Button onClick={handleRegisterBot} isLoading={isRegistering} disabled={!botForm.displayName || !botForm.appId}>
                  Register Bot
                </Button>
                <Button variant="secondary" onClick={() => setShowBotModal(false)} className="mr-3">
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Template Distribution Modal */}
      {showDistributeModal && selectedTemplateForDist && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowDistributeModal(false)}></div>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-xl sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-2">Distribute Template</h3>
                <p className="text-sm text-gray-500 mb-4">
                  Select projects to copy <strong>"{selectedTemplateForDist.name}"</strong> to. This will create a local copy in each selected project.
                </p>

                <div className="max-h-60 overflow-y-auto border border-gray-200 rounded-md p-2 space-y-1">
                  {projects.map(p => (
                    <label key={p.id} className="flex items-center gap-3 p-2 hover:bg-gray-50 rounded cursor-pointer">
                      <input
                        type="checkbox"
                        className="rounded text-blue-600 focus:ring-blue-500"
                        checked={selectedProjectsForDist.includes(p.id)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setSelectedProjectsForDist([...selectedProjectsForDist, p.id]);
                          } else {
                            setSelectedProjectsForDist(selectedProjectsForDist.filter(id => id !== p.id));
                          }
                        }}
                      />
                      <div>
                        <div className="text-sm font-medium text-gray-900">{p.name}</div>
                        <div className="text-xs text-gray-500">{p.department}</div>
                      </div>
                    </label>
                  ))}
                </div>
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse gap-3">
                <Button onClick={handleDistribute} isLoading={isDistributing} disabled={selectedProjectsForDist.length === 0}>
                  Distribute to {selectedProjectsForDist.length} Projects
                </Button>
                <Button variant="secondary" onClick={() => setShowDistributeModal(false)}>
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Company Creation Modal */}
      {showCompanyModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowCompanyModal(false)}></div>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
                  {editingCompany ? 'Edit Company' : 'Create New Company'}
                </h3>
                <div className="space-y-4">
                  <Input
                    label="Company Name *"
                    placeholder="e.g. Contoso Corp"
                    value={companyForm.name}
                    onChange={(e) => setCompanyForm({ ...companyForm, name: e.target.value })}
                  />
                  <Input
                    label="Contact Email *"
                    type="email"
                    placeholder="admin@company.com"
                    value={companyForm.contactEmail}
                    onChange={(e) => setCompanyForm({ ...companyForm, contactEmail: e.target.value })}
                  />
                  <Input
                    label="Contact Phone *"
                    placeholder="+1-555-123-4567"
                    value={companyForm.contactPhone}
                    onChange={(e) => setCompanyForm({ ...companyForm, contactPhone: e.target.value })}
                  />
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Address *</label>
                    <textarea
                      className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      placeholder="Company address"
                      rows={3}
                      value={companyForm.address}
                      onChange={(e) => setCompanyForm({ ...companyForm, address: e.target.value })}
                    />
                  </div>
                  <div className="flex items-center">
                    <input
                      type="checkbox"
                      className="rounded text-blue-600 focus:ring-blue-500"
                      checked={companyForm.billingEnabled}
                      onChange={(e) => setCompanyForm({ ...companyForm, billingEnabled: e.target.checked })}
                    />
                    <label className="ml-2 text-sm text-gray-700">Enable Billing</label>
                  </div>
                </div>
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <Button onClick={handleCreateCompany} isLoading={isCreating} disabled={!companyForm.name || !companyForm.contactEmail || !companyForm.contactPhone || !companyForm.address}>
                  {editingCompany ? 'Update Company' : 'Create Company'}
                </Button>
                <Button variant="secondary" onClick={() => {
                  setShowCompanyModal(false);
                  setEditingCompany(null);
                  setCompanyForm({ name: '', contactEmail: '', contactPhone: '', address: '', billingEnabled: false });
                }} className="mr-3">
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Project Creation Modal */}
      {showProjectModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowProjectModal(false)}></div>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
                  {editingProject ? 'Edit Project' : 'Create New Project'}
                </h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Company *</label>
                    <select
                      className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      value={projectForm.companyId}
                      onChange={(e) => setProjectForm({ ...projectForm, companyId: e.target.value })}
                      disabled={!!editingProject}
                    >
                      <option value="">Select a company</option>
                      {companies.map(c => (
                        <option key={c.id} value={c.id}>{c.name}</option>
                      ))}
                    </select>
                    {editingProject && <p className="text-xs text-gray-500 mt-1">Company cannot be changed after creation</p>}
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Project Name (Notify Key) *
                    </label>
                    <input
                      type="text"
                      className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      placeholder="e.g. finance-core"
                      value={projectForm.notifyKey}
                      onChange={(e) => setProjectForm({ ...projectForm, notifyKey: e.target.value })}
                    />
                    <p className="mt-1 text-xs text-gray-500">Unique identifier for this project (alphanumeric, 3-50 characters)</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Description *</label>
                    <textarea
                      className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      placeholder="Project description"
                      rows={3}
                      value={projectForm.description}
                      onChange={(e) => setProjectForm({ ...projectForm, description: e.target.value })}
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <Input
                      label="Daily Limit"
                      type="number"
                      value={projectForm.dailyLimit.toString()}
                      onChange={(e) => setProjectForm({ ...projectForm, dailyLimit: parseInt(e.target.value) || 10000 })}
                    />
                    <Input
                      label="Monthly Limit"
                      type="number"
                      value={projectForm.monthlyLimit.toString()}
                      onChange={(e) => setProjectForm({ ...projectForm, monthlyLimit: parseInt(e.target.value) || 300000 })}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Priority</label>
                    <select
                      className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      value={projectForm.priority}
                      onChange={(e) => setProjectForm({ ...projectForm, priority: e.target.value as 'low' | 'normal' | 'high' })}
                    >
                      <option value="low">Low</option>
                      <option value="normal">Normal</option>
                      <option value="high">High</option>
                    </select>
                  </div>
                </div>
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <Button onClick={handleCreateProject} isLoading={isCreating} disabled={!projectForm.companyId || !projectForm.notifyKey || !projectForm.description}>
                  {editingProject ? 'Update Project' : 'Create Project'}
                </Button>
                <Button variant="secondary" onClick={() => {
                  setShowProjectModal(false);
                  setEditingProject(null);
                  setProjectForm({
                    companyId: '',
                    notifyKey: '',
                    description: '',
                    dailyLimit: 10000,
                    monthlyLimit: 300000,
                    priority: 'normal',
                  });
                }} className="mr-3">
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Global Template Modal */}
      {showGlobalTemplateModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowGlobalTemplateModal(false)}></div>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-6">{editingGlobalTemplate ? 'Edit Global Template' : 'Create New Global Template'}</h3>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <Input
                      label="Template Name"
                      value={globalTemplateForm.name || ''}
                      onChange={(e) => setGlobalTemplateForm({ ...globalTemplateForm, name: e.target.value })}
                    />
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                      <textarea
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
                        rows={3}
                        value={globalTemplateForm.description || ''}
                        onChange={(e) => setGlobalTemplateForm({ ...globalTemplateForm, description: e.target.value })}
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Variables</label>
                      <div className="flex gap-2 mb-2">
                        <input
                          type="text"
                          className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
                          placeholder="Add variable (e.g. title)"
                          value={globalTemplateVariableInput}
                          onChange={(e) => setGlobalTemplateVariableInput(e.target.value)}
                          onKeyDown={(e) => e.key === 'Enter' && addGlobalTemplateVariable()}
                        />
                        <Button onClick={addGlobalTemplateVariable} variant="secondary" className="px-3">Add</Button>
                      </div>
                      <div className="flex flex-wrap gap-2 p-2 bg-gray-50 rounded-md border border-gray-200 min-h-[40px]">
                        {globalTemplateForm.variables?.map(v => (
                          <span key={v} className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-white border border-gray-300 text-gray-700">
                            <Variable size={12} className="mr-1 text-blue-500" />
                            {v}
                            <button onClick={() => removeGlobalTemplateVariable(v)} className="ml-1 text-gray-400 hover:text-red-500"><Trash2 size={12} /></button>
                          </span>
                        ))}
                      </div>
                      <p className="text-xs text-gray-500 mt-1">Use <code>{`{variableName}`}</code> in your JSON to substitute values.</p>
                    </div>
                  </div>

                  <div className="flex flex-col h-full">
                    <label className="block text-sm font-medium text-gray-700 mb-1 flex items-center gap-2">
                      <Code size={16} /> JSON Structure
                    </label>
                    <textarea
                      className="flex-1 w-full p-4 font-mono text-xs bg-slate-900 text-slate-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
                      value={globalTemplateForm.defaultJsonStructure || ''}
                      onChange={(e) => setGlobalTemplateForm({ ...globalTemplateForm, defaultJsonStructure: e.target.value })}
                      spellCheck={false}
                    />
                  </div>
                </div>
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse gap-3">
                <Button onClick={handleSaveGlobalTemplate} isLoading={isSavingGlobalTemplate} disabled={!globalTemplateForm.name}>
                  Save Template
                </Button>
                <Button variant="secondary" onClick={() => setShowGlobalTemplateModal(false)}>
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
