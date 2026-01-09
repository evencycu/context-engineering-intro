import React, { useEffect, useState } from 'react';
import { useProject } from '../context/ProjectContext';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Line } from 'recharts';
import { Activity, CheckCircle, AlertTriangle, Send } from 'lucide-react';
import { ProjectSwitcher } from '../components/ProjectSwitcher';

const StatCard = ({ title, value, subtext, icon: Icon, color }: any) => (
  <div className="bg-white p-6 rounded-xl border border-gray-100 shadow-sm flex items-start justify-between">
    <div>
      <p className="text-sm font-medium text-gray-500">{title}</p>
      <h3 className="text-2xl font-bold text-gray-900 mt-2">{value}</h3>
      <p className={`text-xs mt-1 ${subtext.includes('+') ? 'text-green-600' : 'text-gray-400'}`}>{subtext}</p>
    </div>
    <div className={`p-3 rounded-lg ${color}`}>
      <Icon size={24} className="text-white" />
    </div>
  </div>
);

export const Dashboard: React.FC = () => {
  const { currentProject } = useProject();
  
  // Mock data that changes when project changes
  const data = [
    { name: 'Mon', sent: 40, failed: 2 },
    { name: 'Tue', sent: 30, failed: 1 },
    { name: 'Wed', sent: 20, failed: 0 },
    { name: 'Thu', sent: 27, failed: 3 },
    { name: 'Fri', sent: 18, failed: 0 },
    { name: 'Sat', sent: 23, failed: 1 },
    { name: 'Sun', sent: 34, failed: 2 },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
          <p className="text-gray-500 mt-1">Overview for {currentProject.name}</p>
        </div>
        <ProjectSwitcher />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard title="Total Sent (7d)" value="1,284" subtext="+12% vs last week" icon={Send} color="bg-blue-500" />
        <StatCard title="Success Rate" value="98.2%" subtext="Consistent performance" icon={CheckCircle} color="bg-green-500" />
        <StatCard title="Failed Messages" value="12" subtext="Mostly connection timeouts" icon={AlertTriangle} color="bg-amber-500" />
        <StatCard title="API Latency" value="245ms" subtext="Average response time" icon={Activity} color="bg-indigo-500" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white p-6 rounded-xl border border-gray-100 shadow-sm">
          <h3 className="text-lg font-semibold text-gray-900 mb-6">Message Volume</h3>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={data}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                <XAxis dataKey="name" axisLine={false} tickLine={false} />
                <YAxis axisLine={false} tickLine={false} />
                <Tooltip 
                  contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
                />
                <Bar dataKey="sent" fill="#3b82f6" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white p-6 rounded-xl border border-gray-100 shadow-sm">
          <h3 className="text-lg font-semibold text-gray-900 mb-6">Error Trends</h3>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={data}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                <XAxis dataKey="name" axisLine={false} tickLine={false} />
                <YAxis axisLine={false} tickLine={false} />
                <Tooltip />
                <Line type="monotone" dataKey="failed" stroke="#ef4444" strokeWidth={3} dot={{ r: 4 }} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  );
};