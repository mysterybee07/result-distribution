// hooks/useNotices.js
import { useQuery } from '@tanstack/react-query';
import api from '../api';

export const useNotices = ({ programId = null, batchId = null, enabled = true }) => {
  return useQuery({
    queryKey: ['notices', programId, batchId],
    queryFn: async () => {
      let apiUrl = '/notice';
      
      // If both programId and batchId are provided, use the filtered endpoint
      if (programId && batchId) {
        apiUrl = `/notice/by-program-and-batch?program_id=${programId}&batch_id=${batchId}`;
      }
      
      const response = await api.get(apiUrl);
      return Array.isArray(response.data.notices) ? response.data.notices : [];
    },
    enabled: enabled,
    staleTime: 2 * 60 * 1000, // Consider data fresh for 2 minutes
    cacheTime: 10 * 60 * 1000, // Keep data in cache for 10 minutes
  });
};