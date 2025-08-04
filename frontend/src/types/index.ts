// Application types
export interface Application {
  id: string;
  name: string;
  description: string;
  git_url: string;
  branch: string;
  port: number;
  environment: 'development' | 'staging' | 'production';
  status: 'pending' | 'building' | 'deploying' | 'running' | 'failed' | 'stopped';
  image_url?: string;
  deployment_url?: string;
  replicas: number;
  user_id: string;
  created_at: string;
  updated_at: string;
}

export interface ApplicationCreate {
  name: string;
  description: string;
  git_url: string;
  branch: string;
  port: number;
  environment: 'development' | 'staging' | 'production';
  replicas: number;
}

export interface ApplicationUpdate {
  name?: string;
  description?: string;
  git_url?: string;
  branch?: string;
  port?: number;
  environment?: 'development' | 'staging' | 'production';
  replicas?: number;
}

export interface ApplicationList {
  applications: Application[];
  total: number;
  page: number;
  limit: number;
}

// User types
export interface User {
  id: string;
  email: string;
  username: string;
  role: 'admin' | 'user';
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserCreate {
  email: string;
  username: string;
  password: string;
  role: 'admin' | 'user';
}

export interface UserLogin {
  email: string;
  password: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

// API Response types
export interface ApiResponse<T> {
  data?: T;
  message?: string;
  error?: string;
  details?: string;
}

// Pagination types
export interface PaginationParams {
  page?: number;
  limit?: number;
}

// Status badge colors
export const getStatusColor = (status: Application['status']) => {
  switch (status) {
    case 'running':
      return 'bg-green-100 text-green-800';
    case 'building':
    case 'deploying':
      return 'bg-yellow-100 text-yellow-800';
    case 'failed':
      return 'bg-red-100 text-red-800';
    case 'stopped':
      return 'bg-gray-100 text-gray-800';
    default:
      return 'bg-blue-100 text-blue-800';
  }
};

// Environment badge colors
export const getEnvironmentColor = (environment: Application['environment']) => {
  switch (environment) {
    case 'production':
      return 'bg-red-100 text-red-800';
    case 'staging':
      return 'bg-orange-100 text-orange-800';
    default:
      return 'bg-blue-100 text-blue-800';
  }
}; 