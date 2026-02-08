import axios from 'axios';

const API_BASE = 'http://localhost:8080/api';

export const apiClient = axios.create({
  baseURL: API_BASE,
});
