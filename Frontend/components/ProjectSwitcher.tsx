import React from 'react';
import { useProject } from '../context/ProjectContext';
import { Briefcase } from 'lucide-react';

export const ProjectSwitcher: React.FC = () => {
  const { currentProject, setCurrentProject, availableProjects } = useProject();

  return (
    <div className="flex items-center gap-4 bg-white p-2 rounded-lg border border-gray-200 shadow-sm">
      <div className="flex items-center gap-2 text-gray-500 pl-2">
        <Briefcase size={18} />
        <span className="text-sm font-medium">Project:</span>
      </div>
      <select
        className="form-select block w-full pl-3 pr-10 py-1.5 text-base border-none focus:outline-none focus:ring-0 sm:text-sm font-medium text-gray-900 bg-transparent"
        value={currentProject.id}
        onChange={(e) => {
          const proj = availableProjects.find(p => p.id === e.target.value);
          if (proj) setCurrentProject(proj);
        }}
      >
        {availableProjects.map((proj) => (
          <option key={proj.id} value={proj.id}>
            {proj.name} ({proj.department})
          </option>
        ))}
      </select>
    </div>
  );
};