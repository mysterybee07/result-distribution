import React, { createContext, useContext, useEffect, useState } from 'react';
import api from '../api'; // Adjust your API import as necessary
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

const AuthContext = createContext();

const fetchUserData = async () => {
    const response = await fetch('http://127.0.0.1:3000/user/active', {
        method: 'GET',
        credentials: 'include', // Include cookies in the request
    });

    if (!response.ok) {
        throw new Error('Unauthorized');
    }
    return response.json();
};

const fetchStudentData = async (userId) => {
    const response = await api.get(`/students/${userId}`);
    if (!response.ok) {
        throw new Error('Failed to fetch student data');
    }
    return response.json();
};

export const AuthProvider = ({ children }) => {
    const authState = localStorage.getItem('isAuthenticated') === 'true';
    const queryClient = useQueryClient();
    const [isAuthenticated, setIsAuthenticated] = useState(authState);
    const [role, setRole] = useState('');
    const [userData, setUserData] = useState(null);

    // Fetch user data with TanStack Query
    const { data: userResponse, isLoading } = useQuery({
        queryKey: ['user'],
        queryFn: fetchUserData,
        retry: false, // Don't retry if unauthorized
        onSuccess: (data) => {
            setUserData(data.data);
            setIsAuthenticated(true);
            setRole(data.data.role);
            localStorage.setItem('user', JSON.stringify(data.data));
            localStorage.setItem('isAuthenticated', 'true');
        },
        onError: () => {
            setIsAuthenticated(false);
            localStorage.setItem('isAuthenticated', 'false');
        },
    });

    // Fetch student data only if the role is 'user'
    // const { data: studentData, isLoading: studentLoading } = useQuery({
    //     queryKey: ['studentData', userData?.id],
    //     queryFn: () => fetchStudentData(userData.id),
    //     enabled: role === 'user' && !!userData?.id,
    //     onSuccess: (data) => {
    //         localStorage.setItem('student_data', JSON.stringify(data)); // Set student data to localStorage on success
    //     },
    // });

    // console.log("🚀 ~ AuthProvider ~ studentData:", studentData);


    // Login mutation
    const loginMutation = useMutation({
        mutationFn: async (user) => {
            setIsAuthenticated(true);
            setRole(user.role);
            setUserData(user);
            localStorage.setItem('user', JSON.stringify(user));
            localStorage.setItem('isAuthenticated', 'true');
        },
        onSuccess: () => {
            queryClient.invalidateQueries(['user']); // Refetch user data after login
        },
    });

    const logout = async () => {
        try {
            await fetch('http://127.0.0.1:3000/user/logout', {
                method: 'POST',
                credentials: 'include',
            });
        } catch (error) {
            console.error('Error logging out:', error);
        } finally {
            setIsAuthenticated(false);
            setRole('');
            setUserData(null);
            localStorage.removeItem('user');
            localStorage.setItem('isAuthenticated', 'false');
            queryClient.invalidateQueries(['user']); // Refetch user data after logout
        }
    };

    return (
        <AuthContext.Provider
            value={{
                isAuthenticated,
                login: loginMutation.mutate,
                logout,
                role,
                userData,
                loading: isLoading,
                // studentData,
                // studentLoading,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
