import React, { useEffect, useState } from 'react';
import { useProject } from '../context/ProjectContext';
import { service } from '../services';
import { NotificationLog, MessageStatus } from '../types';
import { ProjectSwitcher } from '../components/ProjectSwitcher';
import { Eye, CheckCircle, XCircle } from 'lucide-react';
import { Button } from '../components/ui/Button';

export const History: React.FC = () => {
  const { currentProject } = useProject();
  const [logs, setLogs] = useState<NotificationLog[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedLog, setSelectedLog] = useState<NotificationLog | null>(null);

  useEffect(() => {
    // Wait for project to load, and ensure we have a valid UUID
    if (!currentProject || !currentProject.id) {
      return;
    }
    
    // Check if currentProject.id is a valid UUID (36 characters with dashes)
    const isValidUUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(currentProject.id);
    
    if (!isValidUUID) {
      console.warn(`[History] Invalid project ID format: ${currentProject.id}. Expected UUID. Skipping logs load.`);
      setLogs([]);
      setIsLoading(false);
      return;
    }
    
    const fetchLogs = async () => {
      setIsLoading(true);
      try {
        // Use project UUID to fetch logs
        const data = await service.getLogs(currentProject.id);
        setLogs(data);
      } catch (e) {
        console.error("Failed to load logs", e);
        setLogs([]);
      } finally {
        setIsLoading(false);
      }
    };
    fetchLogs();
  }, [currentProject]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
         <div>
            <h2 className="text-2xl font-bold text-gray-900">Audit Logs</h2>
            <p className="text-gray-500 mt-1">History of sent notifications</p>
         </div>
         <ProjectSwitcher />
      </div>

      <div className="bg-white shadow-sm border border-gray-200 rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Timestamp</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Recipient</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Template</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {isLoading ? (
                <tr>
                  <td colSpan={5} className="px-6 py-10 text-center text-sm text-gray-500">Loading records...</td>
                </tr>
              ) : logs.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-10 text-center text-sm text-gray-500">No logs found for this project.</td>
                </tr>
              ) : (
                logs.map((log) => (
                  <tr key={log.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(log.timestamp).toLocaleString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {log.recipient}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                        {log.templateName}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        log.status === MessageStatus.SENT 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
                         {log.status === MessageStatus.SENT ? <CheckCircle size={12} /> : <XCircle size={12} />}
                         {log.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button 
                        onClick={() => setSelectedLog(log)}
                        className="text-blue-600 hover:text-blue-900"
                      >
                        <Eye size={18} />
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Detail Modal */}
      {selectedLog && (
        <div className="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setSelectedLog(null)}></div>
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <div className="sm:flex sm:items-start">
                  <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                    <h3 className="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                      Log Details: {selectedLog.id}
                    </h3>
                    <div className="mt-4 space-y-4">
                      <div className="grid grid-cols-2 gap-4 text-sm">
                         <div>
                           <span className="block text-gray-500">Sent By</span>
                           <span className="font-medium">{selectedLog.senderUpn}</span>
                         </div>
                         <div>
                           <span className="block text-gray-500">Recipient</span>
                           <span className="font-medium">{selectedLog.recipient}</span>
                         </div>
                      </div>
                      
                      {selectedLog.errorMessage && (
                         <div className="p-3 bg-red-50 text-red-700 text-sm rounded-md">
                            Error: {selectedLog.errorMessage}
                         </div>
                      )}

                      <div>
                        <span className="block text-sm font-medium text-gray-700 mb-2">Payload Content</span>
                        <div className="bg-slate-900 rounded-md p-4 overflow-x-auto">
                          <pre className="text-xs text-green-400 font-mono">
                            {JSON.stringify(JSON.parse(selectedLog.contentPreview), null, 2)}
                          </pre>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <Button onClick={() => setSelectedLog(null)} variant="secondary">Close</Button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};