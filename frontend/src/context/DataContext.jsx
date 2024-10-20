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

    return (
        <DataContext.Provider value={{ programs, loadingPrograms, errorPrograms, batches, loadingBatches, errorBatches }}>
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
