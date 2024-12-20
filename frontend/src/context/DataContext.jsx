import React, { createContext, useContext } from 'react';
import { useQuery } from '@tanstack/react-query';
import api from '../api'; // Adjust the import based on your api file location

const DataContext = createContext();

export const DataProvider = ({ children, isAuthenticated }) => {
    // Only fetch data if the user is authenticated
    const { data: programs, isLoading: loadingPrograms, error: errorPrograms } = useQuery({
        queryKey: ['programs'],
        queryFn: async () => {
            const response = await api.get('/program');
            return response.data.programs;
        },
        enabled: isAuthenticated, // Query runs only if user is authenticated
    });

    const { data: batches, isLoading: loadingBatches, error: errorBatches } = useQuery({
        queryKey: ['batches'],
        queryFn: async () => {
            const response = await api.get('/batch');
            return response.data.batches;
        },
        enabled: isAuthenticated, // Query runs only if user is authenticated
    });

    // Uncomment this block if semesters data needs to be fetched
    // const { data: semesters, isLoading: loadingSemesters, error: errorSemesters } = useQuery({
    //     queryKey: ['semesters'],
    //     queryFn: async () => {
    //         const response = await api.get('/semester/by-program/');
    //         return response.data.semesters;
    //     },
    //     enabled: isAuthenticated, // Query runs only if user is authenticated
    // });

    const { data: students, isLoading: loadingStudents, error: errorStudents } = useQuery({
        queryKey: ['students'],
        queryFn: async () => {
            const response = await api.get('/students');
            return response.data.students;
        },
        enabled: isAuthenticated, // Query runs only if user is authenticated
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
                students,
                loadingStudents,
                errorStudents,
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
