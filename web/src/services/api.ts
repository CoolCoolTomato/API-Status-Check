import type { APIConfig, CheckRecord, ApiResponse } from '@/types';

const API_BASE = '/api';

function authHeaders(): Record<string, string> {
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function handleResponse<T>(res: Response): Promise<ApiResponse<T>> {
  if (res.status === 401) {
    localStorage.removeItem('token');
    window.location.href = '/login';
  }
  return res.json();
}

export const apiService = {
  async getAPIs() {
    const res = await fetch(`${API_BASE}/admin/apis`, { headers: authHeaders() });
    return handleResponse<APIConfig[]>(res);
  },

  async createAPI(data: Omit<APIConfig, 'id' | 'created_at' | 'updated_at'>) {
    const res = await fetch(`${API_BASE}/admin/apis`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify(data),
    });
    return handleResponse<APIConfig>(res);
  },

  async updateAPI(id: string, data: Partial<APIConfig>) {
    const res = await fetch(`${API_BASE}/admin/apis/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify(data),
    });
    return handleResponse<null>(res);
  },

  async deleteAPI(id: string) {
    const res = await fetch(`${API_BASE}/admin/apis/${id}`, {
      method: 'DELETE',
      headers: authHeaders(),
    });
    return handleResponse<null>(res);
  },

  async getRecent() {
    const res = await fetch(`${API_BASE}/checks/recent`);
    return handleResponse<CheckRecord[]>(res);
  },

  async getHistory() {
    const res = await fetch(`${API_BASE}/checks/history`);
    return handleResponse<CheckRecord[]>(res);
  },

  async runCheck() {
    const res = await fetch(`${API_BASE}/checks/run`, {
      method: 'POST',
      headers: authHeaders(),
    });
    return handleResponse<string>(res);
  },
};
