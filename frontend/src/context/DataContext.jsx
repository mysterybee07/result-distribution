// DataContext.jsx
import React, { createContext, useContext } from 'react';
import { useQuery } from '@tanstack/react-query';
import api from '../api'; // Adjust the import based on your api file location

const DataContext = createContext();

export const DataProvider = ({ children }) => {
    const { data: programs, isLoading: loadingPrograms, error: errorPrograms } = useQuery({
        queryKey: ['programs'],
        queryFn: async () => {
            const response = await api.get('/program');
            return response.data.programs;
        },
    });

    const { data: batches, isLoading: loadingBatches, error: errorBatches } = useQuery({
        queryKey: ['batches'],
        queryFn: async () => {
            const response = await api.get('/batch');
            return response.data.batches;
        },
    });

    // const { data:semesters, isLoading: loadingSemesters, error: errorSemesters } = useQuery({
    //     queryKey: ['semesters'],
    //     queryFn: async () => {
    //         const response = await api.get('/semester/by-program/');
    //         return response.data.semesters;
    //     },
    // });

    const { data: students, isLoading: loadingStudents, error: errorStudents } = useQuery({
        queryKey: ['students'],
        queryFn: async () => {
            const response = await api.get('/students');
            return response.data.students;
        },
    });

    return (
        <DataContext.Provider
            value={{
                programs,
                loadingPrograms,
                errorPrograms,
                batches,
                loadingBatches,
                errorBatches,
                // semesters,
                // loadingSemesters,
                // errorSemesters,
            }}>
            {children}
        </DataContext.Provider>
    );
};

export const useData = () => {
    const context = useContext(DataContext);
    if (!context) {
        throw new Error('useData must be used within a DataProvider');
    }
    return context;
};
