
import React, { useEffect, useState } from 'react';
import { useProject } from '../context/ProjectContext';
import { service } from '../services';
import { AudienceList, UserProfile, ChatGroup, SyncStatus } from '../types';
import { ProjectSwitcher } from '../components/ProjectSwitcher';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Upload, Users, Tag, FileSpreadsheet, Plus, Search, Filter, Trash2, X, MessageSquare, AlertTriangle, RefreshCw, CheckCircle, Clock } from 'lucide-react';

type Tab = 'lists' | 'tags' | 'chat-groups';

export const Audience: React.FC = () => {
  const { currentProject } = useProject();
  const [activeTab, setActiveTab] = useState<Tab>('lists');
  const [isLoading, setIsLoading] = useState(true);
  
  // Data States
  const [lists, setLists] = useState<AudienceList[]>([]);
  const [users, setUsers] = useState<UserProfile[]>([]);
  const [chatGroups, setChatGroups] = useState<ChatGroup[]>([]);
  
  // Tagging View States
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [tagInput, setTagInput] = useState<string>('');
  const [editingUser, setEditingUser] = useState<string | null>(null); // ID of user being edited

  // Upload Modal State
  const [showUploadModal, setShowUploadModal] = useState(false);
  const [uploadName, setUploadName] = useState('');
  const [isUploading, setIsUploading] = useState(false);

  // Chat Group Modal State
  const [showChatModal, setShowChatModal] = useState(false);
  const [chatForm, setChatForm] = useState({ name: '', chatId: '' });
  const [isRegisteringChat, setIsRegisteringChat] = useState(false);

  // Sync Status State
  const [syncStatus, setSyncStatus] = useState<SyncStatus | null>(null);
  const [isSyncing, setIsSyncing] = useState(false);
  const [syncPollInterval, setSyncPollInterval] = useState<NodeJS.Timeout | null>(null);

  useEffect(() => {
    // Wait for project to load, and ensure we have a valid UUID
    if (!currentProject || !currentProject.id) {
      return;
    }
    
    // Check if currentProject.id is a valid UUID (36 characters with dashes)
    const isValidUUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(currentProject.id);
    
    if (!isValidUUID) {
      console.warn(`[Audience] Invalid project ID format: ${currentProject.id}. Expected UUID. Skipping data load.`);
      setLists([]);
      setUsers([]);
      setChatGroups([]);
      setIsLoading(false);
      return;
    }
    
    loadData();
    loadSyncStatus();
    
    // Cleanup polling on unmount
    return () => {
      if (syncPollInterval) {
        clearInterval(syncPollInterval);
      }
    };
  }, [currentProject]);

  const loadData = async () => {
    if (!currentProject || !currentProject.id) {
      return;
    }
    
    setIsLoading(true);
    try {
      // Use project UUID to fetch data
      const [listsData, usersData, chatData] = await Promise.all([
        service.getAudienceLists(currentProject.id),
        service.getProjectUsers(currentProject.id),
        service.getChatGroups(currentProject.id)
      ]);
      setLists(listsData);
      setUsers(usersData);
      setChatGroups(chatData);
    } catch (e) {
      console.error('Failed to load audience data:', e);
      setLists([]);
      setUsers([]);
      setChatGroups([]);
    } finally {
      setIsLoading(false);
    }
  };

  const loadSyncStatus = async () => {
    try {
      const status = await service.getSyncStatus();
      setSyncStatus(status);
    } catch (e) {
      console.error('Failed to load sync status:', e);
    }
  };

  const handleSyncDirectory = async () => {
    setIsSyncing(true);
    try {
      await service.syncDirectory();
      // Start polling for sync status
      startSyncPolling();
    } catch (error) {
      console.error('Failed to start sync:', error);
      setIsSyncing(false);
    }
  };

  const startSyncPolling = () => {
    // Clear existing interval
    if (syncPollInterval) {
      clearInterval(syncPollInterval);
    }

    // Poll every 2 seconds
    const interval = setInterval(async () => {
      try {
        const status = await service.getSyncStatus();
        setSyncStatus(status);
        
        // Stop polling when sync is complete
        if (!status.isSyncing) {
          clearInterval(interval);
          setSyncPollInterval(null);
          setIsSyncing(false);
          // Reload data after sync completes
          loadData();
        }
      } catch (error) {
        console.error('Failed to poll sync status:', error);
        clearInterval(interval);
        setSyncPollInterval(null);
        setIsSyncing(false);
      }
    }, 2000);

    setSyncPollInterval(interval);
  };

  // --- Derived State for Tags ---
  const allTags = Array.from(new Set(users.flatMap(u => u.tags || []))).sort();
  const filteredUsers = users.filter(u => {
    const matchesSearch = u.displayName.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          u.userPrincipalName.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesTag = selectedTag ? u.tags?.includes(selectedTag) : true;
    return matchesSearch && matchesTag;
  });

  // --- Handlers ---

  const handleUploadList = async () => {
    if (!uploadName) return;
    setIsUploading(true);
    try {
      // Simulate reading a file
      const randomCount = Math.floor(Math.random() * 500) + 10;
      await service.uploadAudienceList(currentProject.id, uploadName, randomCount);
      const updatedLists = await service.getAudienceLists(currentProject.id);
      setLists(updatedLists);
      setShowUploadModal(false);
      setUploadName('');
    } finally {
      setIsUploading(false);
    }
  };

  const handleAddTag = async (userId: string, tag: string) => {
    if (!tag.trim()) return;
    // Optimistic Update
    const updatedUsers = users.map(u => {
       if (u.id === userId) {
         const newTags = u.tags ? [...u.tags, tag] : [tag];
         return { ...u, tags: newTags };
       }
       return u;
    });
    setUsers(updatedUsers);
    
    await service.addTagToUser(userId, tag);
    setTagInput('');
  };

  const handleRemoveTag = async (userId: string, tag: string) => {
    // Optimistic Update
    const updatedUsers = users.map(u => {
       if (u.id === userId && u.tags) {
         return { ...u, tags: u.tags.filter(t => t !== tag) };
       }
       return u;
    });
    setUsers(updatedUsers);

    await service.removeTagFromUser(userId, tag);
  };

  const handleRegisterChatGroup = async () => {
    if (!chatForm.name || !chatForm.chatId) return;
    setIsRegisteringChat(true);
    try {
      await service.registerChatGroup(currentProject.id, chatForm.name, chatForm.chatId);
      const updatedChats = await service.getChatGroups(currentProject.id);
      setChatGroups(updatedChats);
      setShowChatModal(false);
      setChatForm({ name: '', chatId: '' });
    } catch (e) {
      alert('Failed to register chat group');
    } finally {
      setIsRegisteringChat(false);
    }
  };

  return (
    <div className="space-y-6 h-full flex flex-col">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 flex-shrink-0">
         <div>
            <h2 className="text-2xl font-bold text-gray-900">Audience & Tags</h2>
            <p className="text-gray-500 mt-1">Manage target lists, dynamic user tags, and group chats.</p>
         </div>
         <ProjectSwitcher />
      </div>

      {/* Directory Sync Status */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div>
              <div className="flex items-center gap-2 mb-1">
                <h3 className="text-sm font-medium text-gray-900">Directory Status</h3>
                {syncStatus?.isSyncing && (
                  <span className="inline-flex items-center gap-1 text-xs text-blue-600">
                    <Clock size={12} className="animate-spin" />
                    Syncing...
                  </span>
                )}
                {syncStatus && !syncStatus.isSyncing && syncStatus.lastSyncStatus === 'success' && (
                  <span className="inline-flex items-center gap-1 text-xs text-green-600">
                    <CheckCircle size={12} />
                    Synced
                  </span>
                )}
                {syncStatus && syncStatus.lastSyncStatus === 'never' && (
                  <span className="inline-flex items-center gap-1 text-xs text-gray-500">
                    <AlertTriangle size={12} />
                    Not synced
                  </span>
                )}
              </div>
              <div className="text-xs text-gray-600 space-y-1">
                {syncStatus?.lastSyncAt ? (
                  <div>Last Sync: {new Date(syncStatus.lastSyncAt).toLocaleString()}</div>
                ) : (
                  <div>Last Sync: Never</div>
                )}
                <div>
                  Users: <span className="font-medium">{syncStatus?.userCount || 0}</span> | 
                  Groups: <span className="font-medium">{syncStatus?.groupCount || 0}</span>
                </div>
              </div>
            </div>
          </div>
          <Button
            onClick={handleSyncDirectory}
            disabled={isSyncing || (syncStatus?.isSyncing ?? false)}
            isLoading={isSyncing || (syncStatus?.isSyncing ?? false)}
            variant="primary"
            className="flex items-center gap-2"
          >
            <RefreshCw size={16} className={isSyncing || syncStatus?.isSyncing ? 'animate-spin' : ''} />
            Sync Directory
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200 flex-shrink-0">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('lists')}
            className={`${
              activeTab === 'lists'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
          >
            <FileSpreadsheet size={16} />
            Target Lists
          </button>
          <button
            onClick={() => setActiveTab('tags')}
            className={`${
              activeTab === 'tags'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
          >
            <Tag size={16} />
            User Tags
          </button>
          <button
            onClick={() => setActiveTab('chat-groups')}
            className={`${
              activeTab === 'chat-groups'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
          >
            <MessageSquare size={16} />
            Chat Groups
          </button>
        </nav>
      </div>

      {/* Content */}
      <div className="flex-1 min-h-0">
        
        {/* === TAB 1: TARGET LISTS === */}
        {activeTab === 'lists' && (
          <div className="space-y-6">
            <div className="flex justify-end">
              <Button onClick={() => setShowUploadModal(true)}>
                <Upload size={16} className="mr-2" />
                Upload New List
              </Button>
            </div>

            <div className="bg-white shadow-sm border border-gray-200 rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">List Name</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Recipients</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Last Updated</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {lists.length === 0 ? (
                    <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-500">No lists found. Upload one to get started.</td></tr>
                  ) : (
                    lists.map((list) => (
                      <tr key={list.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{list.name}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{list.count} users</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{list.lastUpdated}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                            {list.type}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <button className="text-red-600 hover:text-red-900"><Trash2 size={16} /></button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* === TAB 2: TAGS & USERS === */}
        {activeTab === 'tags' && (
          <div className="flex h-[calc(100vh-250px)] gap-6">
            
            {/* Sidebar: Tags */}
            <div className="w-64 bg-white border border-gray-200 rounded-lg flex flex-col flex-shrink-0">
              <div className="p-4 border-b border-gray-200">
                <h3 className="font-semibold text-gray-900 text-sm">Filter by Tag</h3>
              </div>
              <div className="flex-1 overflow-y-auto p-2 space-y-1">
                <button
                  onClick={() => setSelectedTag(null)}
                  className={`w-full text-left px-3 py-2 rounded-md text-sm font-medium flex items-center justify-between ${
                    selectedTag === null ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'
                  }`}
                >
                  <span>All Users</span>
                  <span className="bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full text-xs">{users.length}</span>
                </button>
                
                {allTags.map(tag => {
                  const count = users.filter(u => u.tags?.includes(tag)).length;
                  return (
                    <button
                      key={tag}
                      onClick={() => setSelectedTag(tag)}
                      className={`w-full text-left px-3 py-2 rounded-md text-sm font-medium flex items-center justify-between ${
                        selectedTag === tag ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'
                      }`}
                    >
                      <div className="flex items-center gap-2 truncate">
                         <Tag size={14} />
                         <span className="truncate">{tag}</span>
                      </div>
                      <span className="bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full text-xs">{count}</span>
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Main: User List */}
            <div className="flex-1 bg-white border border-gray-200 rounded-lg flex flex-col min-h-0">
              <div className="p-4 border-b border-gray-200 flex items-center justify-between gap-4">
                 <div className="relative flex-1">
                    <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input 
                      type="text"
                      placeholder="Search users by name, email, or tag..."
                      className="w-full pl-9 pr-4 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                    />
                 </div>
                 <div className="text-sm text-gray-500">
                    Showing {filteredUsers.length} users
                 </div>
              </div>

              <div className="flex-1 overflow-y-auto p-4 space-y-3">
                 {filteredUsers.length === 0 && (
                   <div className="text-center py-10 text-gray-500">No users match your filters.</div>
                 )}
                 {filteredUsers.map(user => (
                   <div key={user.id} className="flex items-start justify-between p-4 border border-gray-100 rounded-lg hover:bg-gray-50 transition-colors group">
                      <div className="flex items-center gap-3">
                         <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-700 font-bold">
                           {user.displayName.charAt(0)}
                         </div>
                         <div>
                            <h4 className="text-sm font-medium text-gray-900">{user.displayName}</h4>
                            <p className="text-xs text-gray-500">{user.userPrincipalName}</p>
                         </div>
                      </div>

                      <div className="flex items-center gap-2">
                         <div className="flex flex-wrap justify-end gap-1.5 max-w-[300px]">
                            {user.tags?.map(tag => (
                              <span key={tag} className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-700">
                                {tag}
                                <button 
                                  onClick={() => handleRemoveTag(user.id, tag)}
                                  className="ml-1.5 text-gray-400 hover:text-red-500"
                                >
                                  <X size={12} />
                                </button>
                              </span>
                            ))}
                            
                            {editingUser === user.id ? (
                               <div className="flex items-center gap-1">
                                  <input 
                                    type="text" 
                                    autoFocus
                                    className="w-24 text-xs border border-blue-300 rounded px-1 py-0.5 focus:outline-none"
                                    placeholder="New tag..."
                                    value={tagInput}
                                    onChange={(e) => setTagInput(e.target.value)}
                                    onKeyDown={(e) => {
                                      if(e.key === 'Enter') { handleAddTag(user.id, tagInput); setEditingUser(null); }
                                      if(e.key === 'Escape') setEditingUser(null);
                                    }}
                                    onBlur={() => setEditingUser(null)}
                                  />
                               </div>
                            ) : (
                               <button 
                                 onClick={() => { setEditingUser(user.id); setTagInput(''); }}
                                 className="opacity-0 group-hover:opacity-100 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium text-blue-600 hover:bg-blue-50 transition-opacity"
                               >
                                 <Plus size={12} className="mr-1" /> Add Tag
                               </button>
                            )}
                         </div>
                      </div>
                   </div>
                 ))}
              </div>
            </div>
          </div>
        )}

        {/* === TAB 3: CHAT GROUPS === */}
        {activeTab === 'chat-groups' && (
           <div className="space-y-6">
             <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 flex items-start gap-3">
                <AlertTriangle className="text-blue-600 mt-0.5 flex-shrink-0" size={18} />
                <div>
                   <h4 className="text-sm font-semibold text-blue-900">How to use Chat Groups</h4>
                   <p className="text-sm text-blue-800 mt-1">
                      To send notifications to a group chat, you must first invite the <strong>System Bot</strong> to the chat.
                      Then, copy the Chat Link or ID and register it here.
                   </p>
                </div>
             </div>

             <div className="flex justify-end">
               <Button onClick={() => setShowChatModal(true)}>
                 <Plus size={16} className="mr-2" />
                 Register Chat Group
               </Button>
             </div>

             <div className="bg-white shadow-sm border border-gray-200 rounded-lg overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Internal Alias</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Chat ID (Teams)</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Registered By</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {chatGroups.length === 0 ? (
                      <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-500">No chat groups registered.</td></tr>
                    ) : (
                      chatGroups.map((chat) => (
                        <tr key={chat.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{chat.name}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-xs text-gray-500 font-mono truncate max-w-[200px]" title={chat.chatId}>{chat.chatId}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{chat.registeredBy}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{chat.createdAt}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                            <button className="text-red-600 hover:text-red-900"><Trash2 size={16} /></button>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
             </div>
           </div>
        )}
      </div>

      {/* Upload Modal */}
      {showUploadModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75" onClick={() => setShowUploadModal(false)}></div>
            <div className="bg-white rounded-lg overflow-hidden shadow-xl transform transition-all max-w-md w-full p-6 relative z-10">
               <h3 className="text-lg font-medium text-gray-900 mb-4">Upload Target List</h3>
               <div className="space-y-4">
                  <Input 
                    label="List Name" 
                    placeholder="e.g. Q4 Vendors" 
                    value={uploadName}
                    onChange={(e) => setUploadName(e.target.value)}
                  />
                  
                  <div className="border-2 border-dashed border-gray-300 rounded-md p-6 flex flex-col items-center justify-center text-center">
                     <FileSpreadsheet className="h-10 w-10 text-gray-400 mb-2" />
                     <p className="text-sm text-gray-600">Drag and drop Excel/CSV file here</p>
                     <p className="text-xs text-gray-400 mt-1">or click to browse</p>
                  </div>

                  <div className="flex justify-end gap-3 pt-2">
                     <Button variant="secondary" onClick={() => setShowUploadModal(false)}>Cancel</Button>
                     <Button onClick={handleUploadList} isLoading={isUploading} disabled={!uploadName}>Upload</Button>
                  </div>
               </div>
            </div>
          </div>
        </div>
      )}

      {/* Chat Registration Modal */}
      {showChatModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75" onClick={() => setShowChatModal(false)}></div>
            <div className="bg-white rounded-lg overflow-hidden shadow-xl transform transition-all max-w-lg w-full p-6 relative z-10">
               <h3 className="text-lg font-medium text-gray-900 mb-4">Register Chat Group</h3>
               <div className="space-y-4">
                  <div className="bg-amber-50 border border-amber-200 rounded p-3 text-sm text-amber-800 flex gap-2">
                     <AlertTriangle size={18} className="flex-shrink-0" />
                     <p>You must add the <strong>System Bot</strong> to the group chat before registering it, otherwise messages will fail.</p>
                  </div>

                  <Input 
                    label="Internal Alias Name" 
                    placeholder="e.g. Project Alpha Standup" 
                    value={chatForm.name}
                    onChange={(e) => setChatForm({...chatForm, name: e.target.value})}
                  />

                  <Input 
                    label="Teams Chat ID (or Link)" 
                    placeholder="19:meeting_... or https://teams.microsoft.com/..." 
                    value={chatForm.chatId}
                    onChange={(e) => setChatForm({...chatForm, chatId: e.target.value})}
                  />
                  <p className="text-xs text-gray-500">You can find this in the link when you 'Get link to chat' in Teams.</p>

                  <div className="flex justify-end gap-3 pt-4">
                     <Button variant="secondary" onClick={() => setShowChatModal(false)}>Cancel</Button>
                     <Button onClick={handleRegisterChatGroup} isLoading={isRegisteringChat} disabled={!chatForm.name || !chatForm.chatId}>Register</Button>
                  </div>
               </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
