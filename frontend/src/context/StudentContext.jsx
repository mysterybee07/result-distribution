import React, { createContext, useContext, useEffect, useState } from 'react';
import api from '../api'; // Adjust your API import as necessary
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

const StudentContext = createContext();

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

export const StudentProvider = ({ children }) => {
    const authState = localStorage.getItem('isAuthenticated') === 'true';
    const queryClient = useQueryClient();
    const [isAuthenticated, setIsAuthenticated] = useState(authState);
    const [role, setRole] = useState(localStorage.getItem('user')?.role || '');
    console.log("🚀 ~ StudentProvider ~ role:", role)
    const [userData, setUserData] = useState(null);
    const [profileData, setProfileData] = useState(null);
    console.log("🚀 ~ StudentProvider ~ profileData:", profileData)

    // In StudentProvider:
    const { data: profile, isLoading: profileLoading } = useQuery({
        queryKey: ['profileData'],
        queryFn: fetchProfileData,
        enabled: isAuthenticated && role !== "admin", // Only fetch when authenticated
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
            setRole(profile.Users.role);
            localStorage.setItem('student_data', JSON.stringify(profile.Students));
        }
    }, [profile]);

    return (
        <StudentContext.Provider
            value={{
                profileData,
                profileLoading,
            }}
        >
            {children}
        </StudentContext.Provider>
    );
};

export const useStudent = () => useContext(StudentContext);
