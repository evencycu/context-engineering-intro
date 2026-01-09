
import React from 'react';
import { HashRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ProjectProvider } from './context/ProjectContext';
import { Sidebar } from './components/Sidebar';
import { Dashboard } from './pages/Dashboard';
import { Compose } from './pages/Compose';
import { History } from './pages/History';
import { Settings } from './pages/Settings';
import { Audience } from './pages/Audience';
import { Templates } from './pages/Templates';
import { Admin } from './pages/Admin';

function App() {
  return (
    <ProjectProvider>
      <Router>
        <div className="flex min-h-screen bg-slate-50 font-sans">
          <Sidebar />
          <main className="flex-1 ml-64 p-8 overflow-y-auto h-screen">
            <div className="max-w-7xl mx-auto">
              <Routes>
                {/* Workspace Routes */}
                <Route path="/" element={<Dashboard />} />
                <Route path="/compose" element={<Compose />} />
                <Route path="/templates" element={<Templates />} />
                <Route path="/history" element={<History />} />
                <Route path="/audience" element={<Audience />} />
                <Route path="/settings" element={<Settings />} />
                
                {/* Admin Routes */}
                <Route path="/admin" element={<Navigate to="/admin/health" replace />} />
                <Route path="/admin/health" element={<Admin view="health" />} />
                <Route path="/admin/organization" element={<Admin view="organization" />} />
                <Route path="/admin/aad" element={<Admin view="aad" />} />
                <Route path="/admin/templates" element={<Admin view="templates" />} />
                <Route path="/admin/bots" element={<Admin view="bots" />} />
                <Route path="/admin/billing" element={<Admin view="billing" />} />
                <Route path="/admin/ops" element={<Admin view="ops" />} />

                {/* Catch all */}
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </div>
          </main>
        </div>
      </Router>
    </ProjectProvider>
  );
}

export default App;
