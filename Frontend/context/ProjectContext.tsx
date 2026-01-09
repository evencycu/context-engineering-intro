import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Project } from '../types';
import { MOCK_PROJECTS } from '../constants';

interface ProjectContextType {
  currentProject: Project;
  setCurrentProject: (project: Project) => void;
  availableProjects: Project[];
}

const ProjectContext = createContext<ProjectContextType | undefined>(undefined);

export const ProjectProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  // Default to first project
  const [currentProject, setCurrentProject] = useState<Project>(MOCK_PROJECTS[0]);
  
  return (
    <ProjectContext.Provider value={{ 
      currentProject, 
      setCurrentProject, 
      availableProjects: MOCK_PROJECTS 
    }}>
      {children}
    </ProjectContext.Provider>
  );
};

export const useProject = () => {
  const context = useContext(ProjectContext);
  if (context === undefined) {
    throw new Error('useProject must be used within a ProjectProvider');
  }
  return context;
};