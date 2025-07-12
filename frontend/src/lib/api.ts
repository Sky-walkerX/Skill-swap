/**
 * API configuration and utilities for communicating with the backend
 */

// Get the API base URL from environment variables
const getApiBaseUrl = (): string => {
  // In browser environment, use NEXT_PUBLIC_API_URL
  if (typeof window !== 'undefined') {
    return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  }
  
  // In server environment, use GO_BACKEND_URL or fallback
  return process.env.GO_BACKEND_URL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
};

export const API_BASE_URL = getApiBaseUrl();

/**
 * Utility function to make API requests
 */
export const apiRequest = async (
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> => {
  const url = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint : `/${endpoint}`}`;
  
  const defaultOptions: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    credentials: 'include', // Include cookies for authentication
  };

  return fetch(url, {
    ...defaultOptions,
    ...options,
  });
};

/**
 * Utility function for GET requests
 */
export const apiGet = async (endpoint: string, options: RequestInit = {}): Promise<Response> => {
  return apiRequest(endpoint, {
    method: 'GET',
    ...options,
  });
};

/**
 * Utility function for POST requests
 */
export const apiPost = async (
  endpoint: string,
  data?: any,
  options: RequestInit = {}
): Promise<Response> => {
  return apiRequest(endpoint, {
    method: 'POST',
    body: data ? JSON.stringify(data) : undefined,
    ...options,
  });
};

/**
 * Utility function for PUT requests
 */
export const apiPut = async (
  endpoint: string,
  data?: any,
  options: RequestInit = {}
): Promise<Response> => {
  return apiRequest(endpoint, {
    method: 'PUT',
    body: data ? JSON.stringify(data) : undefined,
    ...options,
  });
};

/**
 * Utility function for DELETE requests
 */
export const apiDelete = async (endpoint: string, options: RequestInit = {}): Promise<Response> => {
  return apiRequest(endpoint, {
    method: 'DELETE',
    ...options,
  });
};
