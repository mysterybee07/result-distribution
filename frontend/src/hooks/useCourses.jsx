// In a shared hooks file (e.g., src/hooks/useCourses.js)
import { useQuery } from '@tanstack/react-query';
import api from '../api';

// Hook for all courses
export const useAllCourses = () => {
  return useQuery({
    queryKey: ['allCourses'],
    queryFn: async () => {
      const response = await api.get('/courses');
      return response.data.courses;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 60 * 60 * 1000, // 1 hour - courses likely change less frequently
  });
};

// Hook for filtered courses
export const useFilteredCourses = ({ programId, semesterId, enabled = true }) => {
  return useQuery({
    queryKey: ['filteredCourses', programId, semesterId],
    queryFn: async ({ queryKey }) => {
      const [, program_id, semester_id] = queryKey;
      if (!program_id || !semester_id) return [];
      
      const response = await api.get('/courses/filter', {
        params: { program_id, semester_id }
      });
      return response.data.courses;
    },
    enabled: enabled && !!programId && !!semesterId,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

// Combined hook that provides both datasets
export const useCourses = ({ programId, semesterId, enableFiltering = false }) => {
    // Hook for all courses
    const { 
      data: allCourses = [], 
      isLoading: loadingAllCourses,
      error: allCoursesError,
      ...allCoursesQuery
    } = useQuery({
      queryKey: ['allCourses'],
      queryFn: async () => {
        const response = await api.get('/courses');
        return response.data.courses;
      },
      staleTime: 5 * 60 * 1000, // 5 minutes
    });
    
    // Hook for filtered courses
    const {
      data: filteredCourses = [],
      isLoading: loadingFilteredCourses,
      error: filteredCoursesError,
      ...filteredCoursesQuery
    } = useQuery({
      queryKey: ['filteredCourses', programId, semesterId],
      queryFn: async ({ queryKey }) => {
        const [, program_id, semester_id] = queryKey;
        if (!program_id || !semester_id) return [];
        
        const response = await api.get('/courses/filter', {
          params: { program_id, semester_id }
        });
        return response.data.courses;
      },
      enabled: enableFiltering && !!programId && !!semesterId,
      staleTime: 2 * 60 * 1000, // 2 minutes
    });
    
    return {
      allCourses,
      filteredCourses,
      loadingAllCourses,
      loadingFilteredCourses,
      isLoading: loadingAllCourses || (enableFiltering && loadingFilteredCourses),
      error: allCoursesError || filteredCoursesError,
      allCoursesQuery,
      filteredCoursesQuery,
    };
  };