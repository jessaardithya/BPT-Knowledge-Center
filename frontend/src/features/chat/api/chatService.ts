import { apiClient } from '@/utils/api';

// Handle Chat
export const sendChatMessage = async (message: string) => {
  const response = await apiClient.post('/chat', { message });
  return response.data;
};
