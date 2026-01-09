
import React from 'react';
import { ProjectSwitcher } from '../components/ProjectSwitcher';
import { Bot } from 'lucide-react';

export const Automation: React.FC = () => {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
         <div>
            <h2 className="text-2xl font-bold text-gray-900">Automation Rules</h2>
            <p className="text-gray-500 mt-1">Configure keywords and webhooks.</p>
         </div>
         <ProjectSwitcher />
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-12 text-center">
        <div className="mx-auto h-16 w-16 bg-purple-50 rounded-full flex items-center justify-center mb-4">
          <Bot className="text-purple-500" size={32} />
        </div>
        <h3 className="text-lg font-medium text-gray-900">Coming Soon</h3>
        <p className="text-gray-500 mt-2 max-w-md mx-auto">
          Define keywords to trigger automatic webhooks or auto-replies for your project's bot.
        </p>
      </div>
    </div>
  );
};
