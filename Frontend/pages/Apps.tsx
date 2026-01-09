
import React from 'react';
import { ProjectSwitcher } from '../components/ProjectSwitcher';
import { Construction } from 'lucide-react';

export const Apps: React.FC = () => {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
         <div>
            <h2 className="text-2xl font-bold text-gray-900">Apps & Surveys</h2>
            <p className="text-gray-500 mt-1">Mini-apps for safety checks and polls.</p>
         </div>
         <ProjectSwitcher />
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-12 text-center">
        <div className="mx-auto h-16 w-16 bg-blue-50 rounded-full flex items-center justify-center mb-4">
          <Construction className="text-blue-500" size={32} />
        </div>
        <h3 className="text-lg font-medium text-gray-900">Under Construction</h3>
        <p className="text-gray-500 mt-2 max-w-md mx-auto">
          This module will allow you to create simple data collection apps and attach them to notifications.
        </p>
      </div>
    </div>
  );
};
