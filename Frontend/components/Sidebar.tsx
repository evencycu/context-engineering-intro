
import React from 'react';
import { NavLink, useLocation, useNavigate } from 'react-router-dom';
import { LayoutDashboard, Send, History, Settings, LogOut, ShieldCheck, Users, Bot, Server, Activity, Building2, CreditCard, ArrowRightLeft, Contact, LayoutTemplate } from 'lucide-react';
import { useProject } from '../context/ProjectContext';

export const Sidebar: React.FC = () => {
  const { currentProject } = useProject();
  const location = useLocation();
  const navigate = useNavigate();
  const isAdminMode = location.pathname.startsWith('/admin');

  const workspaceItems = [
    { icon: LayoutDashboard, label: 'Dashboard', path: '/' },
    { icon: Send, label: 'Compose & Schedule', path: '/compose' },
    { icon: LayoutTemplate, label: 'Templates', path: '/templates' },
    { icon: History, label: 'History & Logs', path: '/history' },
    { icon: Users, label: 'Audience & Tags', path: '/audience' },
    { icon: Settings, label: 'Settings', path: '/settings' },
  ];

  const adminItems = [
    { icon: Activity, label: 'System Health', path: '/admin/health' },
    { icon: Building2, label: 'Organization', path: '/admin/organization' },
    { icon: Contact, label: 'Directory Management', path: '/admin/aad' },
    { icon: LayoutTemplate, label: 'Global Templates', path: '/admin/templates' },
    { icon: Bot, label: 'Bot Registry', path: '/admin/bots' },
    { icon: CreditCard, label: 'Billing Center', path: '/admin/billing' },
    { icon: Server, label: 'System Ops', path: '/admin/ops' },
  ];

  const NavItem = ({ item }: { item: any }) => (
    <NavLink
      to={item.path}
      className={({ isActive }) =>
        `flex items-center gap-3 px-4 py-2.5 rounded-md transition-colors text-sm font-medium ${
          isActive
            ? isAdminMode 
              ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-900/50' 
              : 'bg-blue-600 text-white shadow-lg shadow-blue-900/50'
            : 'text-slate-300 hover:bg-slate-800 hover:text-white'
        }`
      }
    >
      <item.icon size={18} />
      {item.label}
    </NavLink>
  );

  return (
    <div className={`w-64 text-white flex flex-col h-screen fixed left-0 top-0 overflow-y-auto transition-colors duration-300 ${isAdminMode ? 'bg-slate-900' : 'bg-slate-900'}`}>
      
      {/* Header */}
      <div className="p-6 border-b border-slate-800 sticky top-0 bg-slate-900 z-10">
        <div className="flex items-center gap-3">
          <div className={`p-2 rounded-lg ${isAdminMode ? 'bg-indigo-600' : 'bg-blue-600'}`}>
             <ShieldCheck size={24} className="text-white" />
          </div>
          <div>
             <h1 className="text-lg font-bold tracking-tight">
               {isAdminMode ? 'Global Admin' : 'Teams Center'}
             </h1>
             <p className="text-xs text-slate-400">Enterprise Edition</p>
          </div>
        </div>
      </div>

      {/* Context / Mode Info */}
      <div className="px-6 py-4">
        {isAdminMode ? (
          <div className="bg-indigo-900/30 rounded-md p-3 border border-indigo-800/50">
             <p className="text-xs text-indigo-400 uppercase font-semibold mb-1">Role</p>
             <p className="text-sm font-medium text-white truncate">Super Administrator</p>
             <p className="text-xs text-slate-400 truncate">System Wide Access</p>
          </div>
        ) : (
          <div className="bg-slate-800 rounded-md p-3 border border-slate-700">
             <p className="text-xs text-slate-500 uppercase font-semibold mb-1">Current Workspace</p>
             <p className="text-sm font-medium text-blue-400 truncate">{currentProject.department}</p>
             <p className="text-xs text-slate-400 truncate">{currentProject.name}</p>
          </div>
        )}
      </div>

      {/* Navigation Links */}
      <nav className="flex-1 px-4 space-y-6 mt-2 pb-6">
        <div>
          <h3 className="px-4 text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">
            {isAdminMode ? 'Management Console' : 'Project Workspace'}
          </h3>
          <div className="space-y-1">
            {(isAdminMode ? adminItems : workspaceItems).map((item) => (
              <NavItem key={item.path} item={item} />
            ))}
          </div>
        </div>
      </nav>

      {/* Footer / Switcher */}
      <div className="p-4 border-t border-slate-800 sticky bottom-0 bg-slate-900 space-y-2">
        
        <button 
          onClick={() => navigate(isAdminMode ? '/' : '/admin/health')}
          className="flex items-center gap-3 px-4 py-2.5 w-full rounded-md text-sm font-medium text-slate-300 hover:bg-slate-800 hover:text-white transition-colors border border-slate-700 hover:border-slate-600"
        >
          <ArrowRightLeft size={18} />
          Switch to {isAdminMode ? 'Teams Center' : 'Global Admin'}
        </button>

        <button className="flex items-center gap-3 px-4 py-2 text-slate-400 hover:text-white text-sm w-full rounded-md hover:bg-slate-800/50 transition-colors">
          <LogOut size={18} />
          Sign Out
        </button>
      </div>
    </div>
  );
};
