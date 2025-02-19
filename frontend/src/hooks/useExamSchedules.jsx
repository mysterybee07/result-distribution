import { useQuery } from '@tanstack/react-query';
import api from '../api';

export const useExamSchedules = ({ batchID, programID, semesterID }) => {
  return useQuery({
    queryKey: ['examSchedules', { batchID, programID, semesterID }],
    queryFn: async ({ queryKey }) => {
      const [, { batchID, programID, semesterID }] = queryKey;
      const { data } = await api.get("/exam/schedules/by-batch-program", {
        params: { batch_id: batchID, program_id: programID, semester_id: semesterID },
      });
      return data.examSchedules;
    },
    enabled: !!batchID && !!programID && !!semesterID,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 30 * 60 * 1000, // 30 minutes
  });
};