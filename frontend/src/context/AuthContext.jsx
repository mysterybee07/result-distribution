// AuthContext.js
import { useQuery } from '@tanstack/react-query';
import React, { createContext, useContext, useEffect, useState } from 'react';
import api from '../api';
import { redirect } from 'react-router-dom';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [role, setRole] = useState('');

    // const { data, isLoading, error } = useQuery({
    //     queryKey: ['user'],
    //     queryFn: async () => {
    //         const token = localStorage.getItem('jwt_token');
    //         const response = await api.get('/user/active', {
    //             headers: {
    //                 Authorization: `Bearer ${token}`, // Include token in the request
    //             },
    //         });
    //         return response.data;
    //     },
    // });

    // console.log("active user:", data);
    const getUserData = async () => {
        try {
            const response = await fetch('http://127.0.0.1:3000/user/active', {
                method: 'GET',
                credentials: 'include', // Ensure cookies are included in the request
            });
    
            if (!response.ok) {
                throw new Error('Failed to fetch user data');
            }
    
            const data = await response.json();
            console.log('User data retrieved:', data);
        } catch (error) {
            console.error('Error fetching user data:', error);
        }
    };
    
    // Call the function to get user data
    getUserData();
    
    useEffect(() => {
        const savedToken = localStorage.getItem('jwt_token');

        if (savedToken) {
            // Validate token, fetch user data, and update user state

            setIsAuthenticated(true);
        } else {
            // Redirect to login if no valid token
            redirect('/login');
        }
    }, []);


    const login = () => setIsAuthenticated(true);
    const logout = () => setIsAuthenticated(false);
    const userRole = (value) => setRole(value);

    return (
        <AuthContext.Provider value={{ isAuthenticated, login, logout, userRole, role }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
