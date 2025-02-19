import { useQuery } from '@tanstack/react-query';
import api from '../api';

export const useExamSchedules = ({ batchId, programId, semesterId }) => {
  console.log("🚀 ~ useExamSchedules ~ semesterId:", semesterId)
  console.log("🚀 ~ useExamSchedules ~ programId:", programId)
  console.log("🚀 ~ useExamSchedules ~ batchId:", batchId)
  return useQuery({
    queryKey: ['examSchedules', { batchId, programId, semesterId }],
    queryFn: async ({ queryKey }) => {
      const [, { batchId, programId, semesterId }] = queryKey;
      const { data } = await api.get("/exam/schedules/by-batch-program", {
        params: { batch_id: batchId, program_id: programId, semester_id: semesterId },
      });
      return data.examSchedules;
    },
    enabled: !!batchId && !!programId && !!semesterId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 30 * 60 * 1000, // 30 minutes
  });
};