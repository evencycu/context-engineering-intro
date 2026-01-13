import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Project } from '../types';
import { MOCK_PROJECTS } from '../constants';
import { service } from '../services';

interface ProjectContextType {
  currentProject: Project;
  setCurrentProject: (project: Project) => void;
  availableProjects: Project[];
  isLoading: boolean;
}

const ProjectContext = createContext<ProjectContextType | undefined>(undefined);

export const ProjectProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [availableProjects, setAvailableProjects] = useState<Project[]>(MOCK_PROJECTS);
  const [currentProject, setCurrentProject] = useState<Project>(MOCK_PROJECTS[0]);
  const [isLoading, setIsLoading] = useState(true);
  
  // Load projects from API on mount
  useEffect(() => {
    const loadProjects = async () => {
      try {
        setIsLoading(true);
        const projects = await service.getProjects();
        // Filter out global template project (should already be filtered in apiService, but double-check)
        const GLOBAL_PROJECT_ID = '00000000-0000-0000-0000-000000000000';
        const filteredProjects = projects.filter(p => p.id !== GLOBAL_PROJECT_ID);
        
        if (filteredProjects && filteredProjects.length > 0) {
          setAvailableProjects(filteredProjects);
          // If current project is global project, switch to first available project
          if (currentProject.id === GLOBAL_PROJECT_ID) {
            setCurrentProject(filteredProjects[0]);
          } else {
            // Try to keep current project if it still exists
            const currentStillExists = filteredProjects.find(p => p.id === currentProject.id);
            if (!currentStillExists) {
              setCurrentProject(filteredProjects[0]);
            }
          }
        } else {
          // Fallback to mock if API returns empty
          console.warn('No projects from API, using mock data');
          setAvailableProjects(MOCK_PROJECTS);
          setCurrentProject(MOCK_PROJECTS[0]);
        }
      } catch (error) {
        console.error('Failed to load projects from API, using mock data:', error);
        // Fallback to mock on error
        setAvailableProjects(MOCK_PROJECTS);
        setCurrentProject(MOCK_PROJECTS[0]);
      } finally {
        setIsLoading(false);
      }
    };
    
    loadProjects();
  }, []);
  
  return (
    <ProjectContext.Provider value={{ 
      currentProject, 
      setCurrentProject, 
      availableProjects,
      isLoading
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