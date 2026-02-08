import { apiClient } from '@/utils/api';
import { Document } from '../types';

// Fetch all documents
export const fetchDocuments = async (): Promise<Document[]> => {
  const response = await apiClient.get('/documents');
  return response.data;
};

// Upload document
export const uploadDocument = async (formData: FormData) => {
  const response = await apiClient.post('/documents/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};

// Update document metadata
export const updateDocument = async (
  id: string,
  updates: { display_name?: string; category?: string; description?: string }
) => {
  const response = await apiClient.put(`/documents/${id}`, updates);
  return response.data;
};

// Delete document
export const deleteDocument = async (id: string) => {
  const response = await apiClient.delete(`/documents/${id}`);
  return response.data;
};
