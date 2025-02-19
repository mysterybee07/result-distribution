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

const fetchProfileData = async () => {
    console.log("🔄 fetchProfileData is executing..."); // Debugging log

    try {
        const response = await api.get(`/profile`);
        console.log("✅ API Response:", response); // Log full response

        return response.data; // Axios responses store data in `response.data`
    } catch (error) {
        console.error("❌ Error fetching profile data:", error);
        throw new Error('Failed to fetch student data'); // Ensure error is thrown
    }
};



export const AuthProvider = ({ children }) => {
    const authState = localStorage.getItem('isAuthenticated') === 'true';
    const queryClient = useQueryClient();
    const [isAuthenticated, setIsAuthenticated] = useState(authState);
    const [role, setRole] = useState('');
    console.log("🚀 ~ AuthProvider ~ role:", role)
    const [userData, setUserData] = useState(null);
    const [profileData, setProfileData] = useState(null);
    console.log("🚀 ~ AuthProvider ~ profileData:", profileData)

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

    // In AuthProvider:
    const { data: profile, isLoading: profileLoading } = useQuery({
        queryKey: ['profileData'],
        queryFn: fetchProfileData,
        enabled: isAuthenticated, // Only fetch when authenticated
        onSuccess: (data) => {
            console.log("✅ Query onSuccess triggered");
            console.log("📦 Received data:", data);
            setProfileData(data);
            console.log("🔄 After setProfileData:", data);
            localStorage.setItem('student_data', JSON.stringify(data));
        },
        onError: (error) => {
            console.error("❌ Query failed:", error);
        }
    });

    // Effect to sync profile data with state
    useEffect(() => {
        if (profile) {
            setProfileData(profile);
            localStorage.setItem('student_data', JSON.stringify(profile.Students));
        }
    }, [profile]);

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
                profileData,
                profileLoading,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
