
import api from './api';

export const login = async (credentials: any) => {
  const response = await api.post('/auth/login', credentials);
  if (response.data.access_token) {
    localStorage.setItem('access_token', response.data.access_token);
  }
  return response.data;
};

export const register = async (userData: any) => {
  const response = await api.post('/auth/register', userData);
  if (response.data.access_token) {
    localStorage.setItem('access_token', response.data.access_token);
  }
  return response.data;
};

export const logout = () => {
  localStorage.removeItem('access_token');
};

export const getCurrentUser = async () => {
  const response = await api.get('/auth/me');
  return response.data;
};
