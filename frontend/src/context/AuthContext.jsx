// AuthContext.js
import React, { createContext, useContext, useState } from 'react';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [role, setRole] = useState('');
    console.log("🚀 ~ AuthProvider ~ role:", role)

    const login = () => setIsAuthenticated(true);
    const logout = () => setIsAuthenticated(false);
    const userRole = (value) => setRole(value);

    return (
        <AuthContext.Provider value={{ isAuthenticated, login, logout,userRole, role }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
