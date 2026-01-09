import React, { useState, useEffect } from 'react';
import { useProject } from '../context/ProjectContext';
import { MOCK_TEMPLATES } from '../constants';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { PeoplePicker } from '../components/PeoplePicker';
import { Code, LayoutTemplate, AlertCircle, CheckCircle, Calendar, Users, AlertTriangle, Radio } from 'lucide-react';
import { service } from '../services';
import { ProjectSwitcher } from '../components/ProjectSwitcher';
import { MessagePriority } from '../types';

export const Compose: React.FC = () => {
  const { currentProject } = useProject();
  const [mode, setMode] = useState<'template' | 'json'>('template');
  
  // Targeting
  const [targetType, setTargetType] = useState<'User' | 'Channel' | 'Tag' | 'List'>('User');
  const [recipient, setRecipient] = useState('');
  
  // Message Config
  const [selectedTemplateId, setSelectedTemplateId] = useState('');
  const [templateValues, setTemplateValues] = useState<Record<string, string>>({});
  const [jsonBody, setJsonBody] = useState('{\n  "type": "AdaptiveCard",\n  "version": "1.4",\n  "body": []\n}');
  const [priority, setPriority] = useState<MessagePriority>(MessagePriority.NORMAL);
  const [scheduleDate, setScheduleDate] = useState('');

  const [isSending, setIsSending] = useState(false);
  const [feedback, setFeedback] = useState<{ type: 'success' | 'error', message: string } | null>(null);

  const projectTemplates = MOCK_TEMPLATES.filter(t => t.projectId === currentProject.id);

  useEffect(() => {
    setSelectedTemplateId('');
    setTemplateValues({});
    setFeedback(null);
    setScheduleDate('');
  }, [currentProject]);

  const handleTemplateChange = (id: string) => {
    setSelectedTemplateId(id);
    const tpl = projectTemplates.find(t => t.id === id);
    if (tpl) {
      const initialValues: Record<string, string> = {};
      tpl.variables.forEach(v => initialValues[v] = '');
      setTemplateValues(initialValues);
    }
  };

  const handleSend = async () => {
    if (!recipient && targetType === 'User') {
      setFeedback({ type: 'error', message: 'Please select a recipient.' });
      return;
    }

    setIsSending(true);
    setFeedback(null);

    try {
      let finalContent = '';

      if (mode === 'template') {
        if (!selectedTemplateId) throw new Error('Please select a template');
        const tpl = projectTemplates.find(t => t.id === selectedTemplateId);
        if (!tpl) throw new Error('Invalid template');
        
        let jsonStr = tpl.defaultJsonStructure;
        Object.entries(templateValues).forEach(([key, val]) => {
          jsonStr = jsonStr.replace(`{${key}}`, val);
        });
        finalContent = jsonStr;
      } else {
        try {
          JSON.parse(jsonBody);
          finalContent = jsonBody;
        } catch (e) {
          throw new Error('Invalid JSON syntax');
        }
      }

      const res = await service.sendMessage({
        projectId: currentProject.id,
        recipient: targetType === 'User' ? recipient : `${targetType}: ${recipient || 'Selected Target'}`,
        recipientType: targetType,
        content: finalContent,
        isTemplate: mode === 'template',
        templateId: selectedTemplateId,
        priority,
        scheduledFor: scheduleDate || undefined
      });

      setFeedback({ type: 'success', message: res.message });
      setTemplateValues({});
    } catch (err: any) {
      setFeedback({ type: 'error', message: err.message || 'Failed to send' });
    } finally {
      setIsSending(false);
    }
  };

  return (
    <div className="space-y-6 max-w-5xl mx-auto pb-10">
      <div className="flex items-center justify-between">
         <div>
            <h2 className="text-2xl font-bold text-gray-900">Compose & Schedule</h2>
            <p className="text-gray-500 mt-1">Create notifications with priority and advanced targeting.</p>
         </div>
         <ProjectSwitcher />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        {/* Left Column: Configuration */}
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
            {/* Header / Mode Switcher */}
            <div className="bg-gray-50 border-b border-gray-200 px-6 py-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm w-fit">
                <button
                  onClick={() => setMode('template')}
                  className={`flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                    mode === 'template' ? 'bg-blue-50 text-blue-700' : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  <LayoutTemplate size={16} />
                  Template
                </button>
                <button
                  onClick={() => setMode('json')}
                  className={`flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                    mode === 'json' ? 'bg-blue-50 text-blue-700' : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  <Code size={16} />
                  JSON
                </button>
              </div>
            </div>

            {/* Content Body */}
            <div className="p-6">
              {mode === 'template' ? (
                <div className="space-y-6">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Select Template</label>
                    <select
                      className="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md border"
                      value={selectedTemplateId}
                      onChange={(e) => handleTemplateChange(e.target.value)}
                    >
                      <option value="">-- Choose a template --</option>
                      {projectTemplates.map((t) => (
                        <option key={t.id} value={t.id}>{t.name}</option>
                      ))}
                    </select>
                  </div>

                  {selectedTemplateId && (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 bg-gray-50 p-4 rounded-lg border border-gray-200">
                      {projectTemplates.find(t => t.id === selectedTemplateId)?.variables.map((v) => (
                        <Input
                          key={v}
                          label={v.charAt(0).toUpperCase() + v.slice(1)}
                          placeholder={`Enter ${v}...`}
                          value={templateValues[v] || ''}
                          onChange={(e) => setTemplateValues(prev => ({ ...prev, [v]: e.target.value }))}
                        />
                      ))}
                    </div>
                  )}
                </div>
              ) : (
                <div className="space-y-4">
                   <div className="flex items-center justify-between">
                     <label className="block text-sm font-medium text-gray-700">Adaptive Card JSON</label>
                   </div>
                   <textarea
                     value={jsonBody}
                     onChange={(e) => setJsonBody(e.target.value)}
                     className="w-full h-80 p-4 font-mono text-xs bg-slate-900 text-slate-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-y"
                     spellCheck={false}
                   />
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Right Column: Settings & Target */}
        <div className="space-y-6">
          
          {/* Targeting Card */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-4 flex items-center gap-2">
              <Users size={16} /> Target Audience
            </h3>
            
            <div className="space-y-4">
              <div className="flex gap-2">
                {(['User', 'List', 'Tag'] as const).map((t) => (
                   <button
                     key={t}
                     onClick={() => { setTargetType(t); setRecipient(''); }}
                     className={`flex-1 py-1.5 text-xs font-medium rounded border ${
                       targetType === t 
                         ? 'bg-blue-50 border-blue-200 text-blue-700' 
                         : 'border-gray-200 text-gray-600 hover:bg-gray-50'
                     }`}
                   >
                     {t}
                   </button>
                ))}
              </div>

              {targetType === 'User' && (
                <PeoplePicker value={recipient} onSelect={setRecipient} />
              )}
              
              {targetType === 'List' && (
                 <div>
                   <label className="block text-sm font-medium text-gray-700 mb-1">Select Audience List</label>
                   <select 
                     className="block w-full text-sm border-gray-300 rounded-md shadow-sm"
                     onChange={(e) => setRecipient(e.target.value)}
                   >
                     <option value="">Select list...</option>
                     <option value="list_01">All Employees (HQ)</option>
                     <option value="list_02">External Contractors</option>
                   </select>
                 </div>
              )}

              {targetType === 'Tag' && (
                 <Input label="Tag Name" placeholder="e.g. IT-Managers" value={recipient} onChange={(e) => setRecipient(e.target.value)} />
              )}
            </div>
          </div>

          {/* Delivery Settings */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-4 flex items-center gap-2">
              <Calendar size={16} /> Delivery Options
            </h3>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Priority</label>
                <div className="space-y-2">
                   {Object.values(MessagePriority).map((p) => (
                     <label key={p} className="flex items-center gap-2 cursor-pointer">
                        <input 
                          type="radio" 
                          name="priority" 
                          checked={priority === p} 
                          onChange={() => setPriority(p)}
                          className="text-blue-600 focus:ring-blue-500"
                        />
                        <span className="text-sm text-gray-700">{p}</span>
                     </label>
                   ))}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Schedule (Optional)</label>
                <Input 
                   type="datetime-local" 
                   value={scheduleDate}
                   onChange={(e) => setScheduleDate(e.target.value)}
                />
                <p className="text-xs text-gray-500 mt-1">Leave blank to send immediately.</p>
              </div>
            </div>
          </div>

          {/* Action Button */}
          <div className="bg-gray-50 rounded-xl border border-gray-200 p-4">
             {feedback && (
                <div className={`mb-4 p-3 rounded text-sm flex items-start gap-2 ${feedback.type === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                  {feedback.type === 'success' ? <CheckCircle size={16} className="mt-0.5" /> : <AlertTriangle size={16} className="mt-0.5" />}
                  {feedback.message}
                </div>
             )}
             <Button 
               onClick={handleSend} 
               isLoading={isSending} 
               className="w-full justify-center"
               disabled={isSending}
             >
               {scheduleDate ? 'Schedule Message' : 'Send Now'}
             </Button>
          </div>

        </div>
      </div>
    </div>
  );
};
