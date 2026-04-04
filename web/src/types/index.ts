export interface APIConfig {
  id: string;
  name: string;
  tag: string;
  api_url: string;
  token: string;
  model: string;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface CheckRecord {
  id: string;
  api_id: string;
  name: string;
  tag: string;
  api_url: string;
  model: string;
  available: boolean;
  latency_ms: number;
  checked_at: string;
  status_code: number;
  error_message: string;
  response_preview: string;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}
