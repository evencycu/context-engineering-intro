
import React, { useEffect, useState } from 'react';
import { useProject } from '../context/ProjectContext';
import { service } from '../services';
import { Template } from '../types';
import { ProjectSwitcher } from '../components/ProjectSwitcher';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Plus, Trash2, Edit2, Code, FileText, Variable } from 'lucide-react';

export const Templates: React.FC = () => {
  const { currentProject, isLoading: isProjectLoading } = useProject();
  const [templates, setTemplates] = useState<Template[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  
  // Modal State
  const [showModal, setShowModal] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState<Template | null>(null);
  const [formData, setFormData] = useState<Partial<Template>>({
    name: '',
    description: '',
    variables: [],
    defaultJsonStructure: '{}'
  });
  const [variableInput, setVariableInput] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    // Wait for project to load, and ensure we have a valid UUID
    if (isProjectLoading) {
      return;
    }
    
    // Check if currentProject.id is a valid UUID (36 characters with dashes)
    const isValidUUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(currentProject.id);
    
    if (!isValidUUID) {
      console.warn(`[Templates] Invalid project ID format: ${currentProject.id}. Expected UUID. Skipping template load.`);
      setIsLoading(false);
      return;
    }
    
    loadTemplates();
  }, [currentProject, isProjectLoading]);

  const loadTemplates = async () => {
    setIsLoading(true);
    try {
      const data = await service.getTemplates(currentProject.id);
      setTemplates(data);
    } catch (e: any) {
      console.error('Failed to load templates:', e);
      const errorMessage = e?.message || e?.toString() || 'Failed to load templates';
      alert(`Failed to load templates: ${errorMessage}`);
      setTemplates([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleOpenCreate = () => {
    setEditingTemplate(null);
    setFormData({
      name: '',
      description: '',
      variables: [],
      defaultJsonStructure: JSON.stringify({
        type: "AdaptiveCard",
        version: "1.4",
        body: [{ type: "TextBlock", text: "New Template" }]
      }, null, 2)
    });
    setVariableInput('');
    setShowModal(true);
  };

  const handleOpenEdit = (t: Template) => {
    setEditingTemplate(t);
    setFormData({
      name: t.name,
      description: t.description,
      variables: [...t.variables],
      defaultJsonStructure: t.defaultJsonStructure
    });
    setVariableInput('');
    setShowModal(true);
  };

  const handleSave = async () => {
    if (!formData.name || !formData.defaultJsonStructure) {
      alert('Please fill in template name and JSON structure');
      return;
    }
    
    // Validate JSON structure
    try {
      JSON.parse(formData.defaultJsonStructure!);
    } catch (e) {
      alert('Invalid JSON structure. Please check your JSON syntax.');
      return;
    }
    
    setIsSaving(true);
    try {
      if (editingTemplate) {
        await service.updateTemplate(editingTemplate.id, { ...formData, projectId: currentProject.id });
      } else {
        await service.createTemplate({
          projectId: currentProject.id,
          name: formData.name!,
          description: formData.description || '',
          variables: formData.variables || [],
          defaultJsonStructure: formData.defaultJsonStructure!
        });
      }
      setShowModal(false);
      loadTemplates();
    } catch (e: any) {
      console.error('Failed to save template:', e);
      const errorMessage = e?.message || e?.toString() || 'Failed to save template';
      alert(`Failed to save template: ${errorMessage}`);
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if(!window.confirm("Are you sure you want to delete this template?")) return;
    try {
      await service.deleteTemplate(id, currentProject.id);
      loadTemplates();
    } catch (e: any) {
      console.error('Failed to delete template:', e);
      const errorMessage = e?.message || e?.toString() || 'Failed to delete template';
      alert(`Failed to delete template: ${errorMessage}`);
    }
  };

  const addVariable = () => {
    if (variableInput && !formData.variables?.includes(variableInput)) {
      setFormData({ ...formData, variables: [...(formData.variables || []), variableInput] });
      setVariableInput('');
    }
  };

  const removeVariable = (v: string) => {
    setFormData({ ...formData, variables: formData.variables?.filter(item => item !== v) });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
         <div>
            <h2 className="text-2xl font-bold text-gray-900">Message Templates</h2>
            <p className="text-gray-500 mt-1">Manage reusable Adaptive Card templates for your project.</p>
         </div>
         <ProjectSwitcher />
      </div>

      <div className="flex justify-end">
        <Button onClick={handleOpenCreate}>
          <Plus size={16} className="mr-2" />
          Create Template
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {isLoading ? (
          <p className="col-span-3 text-center text-gray-500">Loading templates...</p>
        ) : templates.length === 0 ? (
          <p className="col-span-3 text-center text-gray-500 py-10 bg-white rounded-lg border border-gray-200 border-dashed">No templates found. Create one to get started.</p>
        ) : (
          templates.map((tpl) => (
            <div key={tpl.id} className="bg-white rounded-xl shadow-sm border border-gray-200 flex flex-col hover:border-blue-300 transition-colors">
              <div className="p-6 flex-1">
                <div className="flex justify-between items-start mb-2">
                  <div className="p-2 bg-blue-50 text-blue-600 rounded-lg">
                    <FileText size={20} />
                  </div>
                  <div className="flex gap-2">
                    <button onClick={() => handleOpenEdit(tpl)} className="text-gray-400 hover:text-blue-600 p-1">
                      <Edit2 size={16} />
                    </button>
                    <button onClick={() => handleDelete(tpl.id)} className="text-gray-400 hover:text-red-600 p-1">
                      <Trash2 size={16} />
                    </button>
                  </div>
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-1">{tpl.name}</h3>
                <p className="text-sm text-gray-500 line-clamp-2 mb-4">{tpl.description}</p>
                
                {tpl.variables.length > 0 && (
                  <div className="flex flex-wrap gap-1">
                    {tpl.variables.map(v => (
                      <span key={v} className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600">
                        {`{${v}}`}
                      </span>
                    ))}
                  </div>
                )}
              </div>
              <div className="px-6 py-3 bg-gray-50 border-t border-gray-200 text-xs text-gray-500 font-mono truncate">
                 ID: {tpl.id}
              </div>
            </div>
          ))
        )}
      </div>

      {/* Editor Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowModal(false)}></div>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-6">{editingTemplate ? 'Edit Template' : 'Create New Template'}</h3>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                   <div className="space-y-4">
                      <Input 
                        label="Template Name" 
                        value={formData.name}
                        onChange={(e) => setFormData({...formData, name: e.target.value})}
                      />
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea
                          className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
                          rows={3}
                          value={formData.description}
                          onChange={(e) => setFormData({...formData, description: e.target.value})}
                        />
                      </div>
                      
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Variables</label>
                        <div className="flex gap-2 mb-2">
                           <input 
                             type="text" 
                             className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
                             placeholder="Add variable (e.g. title)"
                             value={variableInput}
                             onChange={(e) => setVariableInput(e.target.value)}
                             onKeyDown={(e) => e.key === 'Enter' && addVariable()}
                           />
                           <Button onClick={addVariable} variant="secondary" className="px-3">Add</Button>
                        </div>
                        <div className="flex flex-wrap gap-2 p-2 bg-gray-50 rounded-md border border-gray-200 min-h-[40px]">
                           {formData.variables?.map(v => (
                             <span key={v} className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-white border border-gray-300 text-gray-700">
                               <Variable size={12} className="mr-1 text-blue-500" />
                               {v}
                               <button onClick={() => removeVariable(v)} className="ml-1 text-gray-400 hover:text-red-500"><Trash2 size={12}/></button>
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
                         value={formData.defaultJsonStructure}
                         onChange={(e) => setFormData({...formData, defaultJsonStructure: e.target.value})}
                         spellCheck={false}
                      />
                   </div>
                </div>
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse gap-3">
                <Button onClick={handleSave} isLoading={isSaving} disabled={!formData.name}>
                  Save Template
                </Button>
                <Button variant="secondary" onClick={() => setShowModal(false)}>
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
