import React, { createContext, useContext, useEffect, useState } from 'react';
import api from '../api'; // Adjust your API import as necessary

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [role, setRole] = useState('');
    const [userData, setUserData] = useState(null);
    const [loading, setLoading] = useState(true);

    const getUserData = async () => {
        try {
            const response = await fetch('http://127.0.0.1:3000/user/active', {
                method: 'GET',
                credentials: 'include', // Include cookies in the request
            });

            if (response.ok) {
                const data = await response.json();
                console.log("🚀 ~ getUserData ~ data:", data)
                setUserData(data.data);
                setIsAuthenticated(true);
                setRole(data.data.role);
            } else {
                throw new Error('Unauthorized');
            }
        } catch (error) {
            console.error('Error fetching user data:', error);
            setIsAuthenticated(false);
        } finally{
            setLoading(false);
        }
    };
    useEffect(() => {
        // Fetch user data on mount
        getUserData();
    }, []);

    const login = (user) => {
        setIsAuthenticated(true);
        setRole(user.role); // Store role on login
        setUserData(user); // Store user data on login

        // Persist user data to localStorage
        localStorage.setItem('user', JSON.stringify(user));
    };

    const logout = () => {
        setIsAuthenticated(false);
        setRole('');
        setUserData(null); // Clear user data on logout

        // Clear user data from localStorage
        // localStorage.removeItem('user');
    };

    return (
        <AuthContext.Provider value={{ isAuthenticated, login, logout, role, userData, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
