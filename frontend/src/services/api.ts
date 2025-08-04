import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { 
  Application, 
  ApplicationCreate, 
  ApplicationUpdate, 
  ApplicationList,
  User,
  UserCreate,
  UserLogin,
  AuthResponse,
  ApiResponse 
} from '../types';

// Create axios instance
const api: AxiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Application API
export const applicationApi = {
  // Get all applications for current user
  getApplications: async (params?: { page?: number; limit?: number }): Promise<ApplicationList> => {
    const response: AxiosResponse<ApiResponse<ApplicationList>> = await api.get('/applications', { params });
    return response.data.data!;
  },

  // Get single application
  getApplication: async (id: string): Promise<Application> => {
    const response: AxiosResponse<ApiResponse<Application>> = await api.get(`/applications/${id}`);
    return response.data.data!;
  },

  // Create new application
  createApplication: async (data: ApplicationCreate): Promise<Application> => {
    const response: AxiosResponse<ApiResponse<Application>> = await api.post('/applications', data);
    return response.data.data!;
  },

  // Update application
  updateApplication: async (id: string, data: ApplicationUpdate): Promise<void> => {
    await api.put(`/applications/${id}`, data);
  },

  // Delete application
  deleteApplication: async (id: string): Promise<void> => {
    await api.delete(`/applications/${id}`);
  },

  // Deploy application
  deployApplication: async (id: string): Promise<void> => {
    await api.post(`/applications/${id}/deploy`);
  },

  // Get application logs
  getApplicationLogs: async (id: string, lines?: number): Promise<string[]> => {
    const response: AxiosResponse<ApiResponse<string[]>> = await api.get(`/applications/${id}/logs`, {
      params: { lines }
    });
    return response.data.data!;
  },

  // Get application metrics
  getApplicationMetrics: async (id: string): Promise<any> => {
    const response: AxiosResponse<ApiResponse<any>> = await api.get(`/applications/${id}/metrics`);
    return response.data.data!;
  },
};

// Auth API
export const authApi = {
  // Login
  login: async (credentials: UserLogin): Promise<AuthResponse> => {
    const response: AxiosResponse<ApiResponse<AuthResponse>> = await api.post('/auth/login', credentials);
    return response.data.data!;
  },

  // Register
  register: async (userData: UserCreate): Promise<AuthResponse> => {
    const response: AxiosResponse<ApiResponse<AuthResponse>> = await api.post('/auth/register', userData);
    return response.data.data!;
  },

  // Refresh token
  refreshToken: async (refreshToken: string): Promise<AuthResponse> => {
    const response: AxiosResponse<ApiResponse<AuthResponse>> = await api.post('/auth/refresh', {
      refresh_token: refreshToken
    });
    return response.data.data!;
  },

  // Logout
  logout: async (): Promise<void> => {
    await api.post('/auth/logout');
  },

  // Get current user
  getCurrentUser: async (): Promise<User> => {
    const response: AxiosResponse<ApiResponse<User>> = await api.get('/auth/me');
    return response.data.data!;
  },
};

// Health check
export const healthApi = {
  check: async (): Promise<{ status: string; message: string; time: string }> => {
    const response = await axios.get('/health');
    return response.data;
  },
};

export default api; 