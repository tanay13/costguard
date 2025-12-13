'use client';

import { useState, useEffect } from 'react';
import { LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface ScanResult {
  total_current_cost_usd: number;
  total_optimal_cost_usd: number;
  total_potential_savings_usd: number;
  resources: Array<{
    resource: string;
    provider: string;
    costs: {
      current_cost_usd: number;
      optimal_cost_usd: number;
      potential_savings_usd: number;
    };
  }>;
}

interface DecisionLog {
  scan_id: string;
  timestamp: string;
  total_savings: number;
  actions_applied: number;
  pr_url: string;
  pr_number?: number;
  summary: string;
  repo_full_name?: string;
}

interface Repo {
  repo_full_name: string;
  last_scan: string;
  total_savings: number;
}

export default function Dashboard() {
  const [scanData, setScanData] = useState<ScanResult | null>(null);
  const [decisions, setDecisions] = useState<DecisionLog[]>([]);
  const [repos, setRepos] = useState<Repo[]>([]);
  const [selectedRepo, setSelectedRepo] = useState<string>('');
  const [loading, setLoading] = useState(true);

  const loadData = () => {
    const repoParam = selectedRepo ? `?repo=${encodeURIComponent(selectedRepo)}` : '';
    
    // Load repos list
    fetch('/api/repos')
      .then(res => res.json())
      .then(data => {
        setRepos(data);
        if (data.length > 0 && !selectedRepo) {
          // Auto-select first repo if none selected
          setSelectedRepo(data[0].repo_full_name);
        }
      })
      .catch(err => console.error('Failed to load repos:', err));

    // Load scan data
    fetch(`/api/scan${repoParam}`)
      .then(res => res.json())
      .then(data => {
        setScanData(data);
        setLoading(false);
      })
      .catch(err => {
        console.error('Failed to load scan data:', err);
        setLoading(false);
      });

    // Load decision logs
    fetch(`/api/decisions${repoParam}`)
      .then(res => res.json())
      .then(data => setDecisions(data))
      .catch(err => console.error('Failed to load decisions:', err));
  };

  useEffect(() => {
    loadData();
    
    // Auto-refresh every 30 seconds
    const interval = setInterval(loadData, 30000);
    return () => clearInterval(interval);
  }, [selectedRepo]);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-xl text-gray-600">Loading CostGuard Dashboard...</div>
      </div>
    );
  }

  const savingsData = scanData ? [
    { name: 'Current', value: scanData.total_current_cost_usd },
    { name: 'Optimized', value: scanData.total_optimal_cost_usd },
  ] : [];

  const resourceSavings = scanData?.resources.map(r => ({
    name: r.resource,
    savings: r.costs.potential_savings_usd,
    current: r.costs.current_cost_usd,
    optimal: r.costs.optimal_cost_usd,
  })) || [];

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">🚀 CostGuard</h1>
              <p className="mt-2 text-gray-600">Autonomous Cloud Cost Optimization via AI Orchestration</p>
            </div>
            <div className="flex items-center gap-4">
              <label className="text-sm font-medium text-gray-700">
                Filter by Repository:
              </label>
              <select
                value={selectedRepo}
                onChange={(e) => setSelectedRepo(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">All Repositories</option>
                {repos.map((repo) => (
                  <option key={repo.repo_full_name} value={repo.repo_full_name}>
                    {repo.repo_full_name} (${repo.total_savings}/mo savings)
                  </option>
                ))}
              </select>
              <button
                onClick={loadData}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                Refresh
              </button>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-sm font-medium text-gray-500">Current Monthly Cost</h3>
            <p className="mt-2 text-3xl font-bold text-gray-900">
              ${scanData?.total_current_cost_usd || '0.00'}
            </p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-sm font-medium text-gray-500">Optimized Cost</h3>
            <p className="mt-2 text-3xl font-bold text-green-600">
              ${scanData?.total_optimal_cost_usd || '0.00'}
            </p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-sm font-medium text-gray-500">Potential Savings</h3>
            <p className="mt-2 text-3xl font-bold text-blue-600">
              ${scanData?.total_potential_savings_usd || '0.00'}
            </p>
          </div>
        </div>

        {/* Cost Comparison Chart */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Cost Comparison</h2>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={savingsData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip formatter={(value: number) => `$${value}`} />
              <Legend />
              <Bar dataKey="value" fill="#3b82f6" name="Monthly Cost (USD)" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Resource Savings */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Savings by Resource</h2>
          <ResponsiveContainer width="100%" height={400}>
            <BarChart data={resourceSavings}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip formatter={(value: number) => `$${value}`} />
              <Legend />
              <Bar dataKey="current" fill="#ef4444" name="Current Cost" />
              <Bar dataKey="optimal" fill="#10b981" name="Optimized Cost" />
              <Bar dataKey="savings" fill="#3b82f6" name="Savings" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* AI Decision Log */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">🤖 AI Decision Log</h2>
          <div className="space-y-4">
            {decisions.length === 0 ? (
              <p className="text-gray-500">No decisions logged yet.</p>
            ) : (
              decisions.map((decision, idx) => (
                <div key={idx} className="border-l-4 border-blue-500 pl-4 py-2">
                  <div className="flex justify-between items-start">
                    <div>
                      {decision.repo_full_name && (
                        <p className="text-xs text-gray-400 mb-1">
                          📦 {decision.repo_full_name}
                        </p>
                      )}
                      <p className="font-medium">{decision.summary}</p>
                      <p className="text-sm text-gray-500 mt-1">
                        {new Date(decision.timestamp).toLocaleString()}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="text-lg font-bold text-green-600">
                        ${decision.total_savings}/mo
                      </p>
                      <p className="text-sm text-gray-500">
                        {decision.actions_applied} actions applied
                      </p>
                    </div>
                  </div>
                  {decision.pr_url && (
                    <a
                      href={decision.pr_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-sm text-blue-600 hover:underline mt-2 inline-block"
                    >
                      View PR #{decision.pr_number || 'N/A'} →
                    </a>
                  )}
                </div>
              ))
            )}
          </div>
        </div>
      </main>
    </div>
  );
}

