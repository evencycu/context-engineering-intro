import React, { useState, useEffect } from 'react';
import { useProject } from '../context/ProjectContext';
import { ProjectSwitcher } from '../components/ProjectSwitcher';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Save, RefreshCw, Copy, Shield, Bell, Layout, CheckCircle, AlertTriangle, Eye, EyeOff } from 'lucide-react';
import { service } from '../services';

type Tab = 'general' | 'security' | 'notifications';

export const Settings: React.FC = () => {
  const { currentProject, setCurrentProject } = useProject();
  const [activeTab, setActiveTab] = useState<Tab>('general');
  const [isLoading, setIsLoading] = useState(false);
  const [feedback, setFeedback] = useState<{ type: 'success' | 'error', message: string } | null>(null);

  // Form States
  const [projectName, setProjectName] = useState(currentProject.name);
  const [department, setDepartment] = useState(currentProject.department);
  const [showApiKey, setShowApiKey] = useState(false);
  const [apiKey, setApiKey] = useState(currentProject.apiKey);

  // Reset form when project changes
  useEffect(() => {
    setProjectName(currentProject.name);
    setDepartment(currentProject.department);
    setApiKey(currentProject.apiKey);
    setShowApiKey(false);
    setFeedback(null);
  }, [currentProject]);

  const handleSaveGeneral = async () => {
    setIsLoading(true);
    setFeedback(null);
    try {
      await service.updateProject(currentProject.id, { name: projectName, department });
      
      // Update local context to reflect changes immediately
      setCurrentProject({
        ...currentProject,
        name: projectName,
        department: department
      });
      
      setFeedback({ type: 'success', message: 'Project settings updated successfully.' });
    } catch (error) {
      setFeedback({ type: 'error', message: 'Failed to update settings.' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleRegenerateKey = async () => {
    if (!window.confirm("Are you sure? The old API key will stop working immediately.")) return;
    
    setIsLoading(true);
    try {
      const newKey = await service.regenerateApiKey(currentProject.id);
      setApiKey(newKey);
      setCurrentProject({ ...currentProject, apiKey: newKey });
      setFeedback({ type: 'success', message: 'API Key regenerated. Please update your applications.' });
    } catch (error) {
      setFeedback({ type: 'error', message: 'Failed to regenerate API key.' });
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setFeedback({ type: 'success', message: 'Copied to clipboard!' });
    setTimeout(() => setFeedback(null), 2000);
  };

  const tabs = [
    { id: 'general', label: 'General', icon: Layout },
    { id: 'security', label: 'Security & API', icon: Shield },
    { id: 'notifications', label: 'Notification Defaults', icon: Bell },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Settings</h2>
          <p className="text-gray-500 mt-1">Manage configuration for {currentProject.name}</p>
        </div>
        <ProjectSwitcher />
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden min-h-[500px] flex flex-col md:flex-row">
        {/* Sidebar Tabs */}
        <div className="w-full md:w-64 bg-gray-50 border-r border-gray-200 p-4">
          <nav className="space-y-1">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => { setActiveTab(tab.id as Tab); setFeedback(null); }}
                className={`w-full flex items-center gap-3 px-3 py-3 text-sm font-medium rounded-md transition-colors ${
                  activeTab === tab.id
                    ? 'bg-white text-blue-600 shadow-sm ring-1 ring-gray-200'
                    : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'
                }`}
              >
                <tab.icon size={18} />
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        {/* Content Area */}
        <div className="flex-1 p-8">
          {feedback && (
            <div className={`mb-6 p-4 rounded-md flex items-center gap-2 text-sm ${
              feedback.type === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'
            }`}>
              {feedback.type === 'success' ? <CheckCircle size={16} /> : <AlertTriangle size={16} />}
              {feedback.message}
            </div>
          )}

          {activeTab === 'general' && (
            <div className="space-y-6 max-w-lg">
              <div>
                <h3 className="text-lg font-medium text-gray-900 border-b border-gray-200 pb-2 mb-4">Project Details</h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Project ID</label>
                    <div className="flex items-center gap-2">
                       <code className="bg-gray-100 px-2 py-1 rounded text-sm text-gray-600 font-mono">{currentProject.id}</code>
                       <span className="text-xs text-gray-400">(Immutable)</span>
                    </div>
                  </div>
                  <Input 
                    label="Project Name" 
                    value={projectName} 
                    onChange={(e) => setProjectName(e.target.value)} 
                  />
                  <Input 
                    label="Department / Cost Center" 
                    value={department} 
                    onChange={(e) => setDepartment(e.target.value)} 
                  />
                </div>
              </div>
              <div className="pt-4">
                <Button onClick={handleSaveGeneral} isLoading={isLoading}>
                  <Save size={16} className="mr-2" />
                  Save Changes
                </Button>
              </div>
            </div>
          )}

          {activeTab === 'security' && (
            <div className="space-y-8 max-w-2xl">
              <div>
                <h3 className="text-lg font-medium text-gray-900 border-b border-gray-200 pb-2 mb-4">API Credentials</h3>
                <p className="text-sm text-gray-500 mb-4">
                  This API key allows external systems to send notifications on behalf of this project. Keep it secure.
                </p>
                
                <div className="bg-slate-50 border border-slate-200 rounded-lg p-4">
                  <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Secret Key</label>
                  <div className="flex items-center gap-2">
                    <div className="flex-1 font-mono text-sm bg-white border border-gray-300 rounded px-3 py-2 text-gray-700 break-all">
                      {showApiKey ? apiKey : 'sk_live_•••••••••••••••••••••••••••'}
                    </div>
                    <Button variant="secondary" onClick={() => setShowApiKey(!showApiKey)} title={showApiKey ? "Hide" : "Show"}>
                      {showApiKey ? <EyeOff size={16} /> : <Eye size={16} />}
                    </Button>
                    <Button variant="secondary" onClick={() => copyToClipboard(apiKey)} title="Copy">
                      <Copy size={16} />
                    </Button>
                  </div>
                </div>

                <div className="mt-4">
                  <Button variant="danger" onClick={handleRegenerateKey} isLoading={isLoading}>
                    <RefreshCw size={16} className="mr-2" />
                    Regenerate Key
                  </Button>
                  <p className="mt-2 text-xs text-red-500">
                    Warning: Regenerating the key will immediately invalidate the current key.
                  </p>
                </div>
              </div>

              <div>
                <h3 className="text-lg font-medium text-gray-900 border-b border-gray-200 pb-2 mb-4">IP Whitelist</h3>
                <div className="bg-gray-50 p-4 rounded-md text-center text-sm text-gray-500">
                  IP Restrictions are managed via the Azure Portal for this environment.
                </div>
              </div>
            </div>
          )}

          {activeTab === 'notifications' && (
            <div className="space-y-6 max-w-lg">
              <div>
                <h3 className="text-lg font-medium text-gray-900 border-b border-gray-200 pb-2 mb-4">Message Defaults</h3>
                <div className="space-y-4">
                  <Input 
                    label="Default 'From' Name" 
                    placeholder="e.g. IT Operations"
                    defaultValue="IT Notifications"
                  />
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Global Footer Text</label>
                    <textarea 
                      className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      rows={3}
                      placeholder="e.g. Please do not reply to this message."
                      defaultValue="This is an automated message from the Notification Center."
                    ></textarea>
                  </div>
                </div>
              </div>

              <div className="pt-4">
                <Button variant="primary" onClick={() => setFeedback({ type: 'success', message: 'Notification defaults saved.' })}>
                  <Save size={16} className="mr-2" />
                  Save Defaults
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
