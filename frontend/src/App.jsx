import './App.css';
import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import Homepage from './pages/Homepage';
import Login from './pages/Login';
import Register from './pages/Register';
import Exam from './pages/Exam';
import Profile from './pages/Profile';
import Result from './pages/Result';
import Footer from './components/Footer';
import Navbar from './components/Navbar';
import { AuthProvider, useAuth } from './context/AuthContext';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Dashboard from './pages/admin/Dashboard';
import AdminNavbar from './components/AdminNavbar';
import { Toaster } from '@/components/ui/toaster';
import { Navigate } from 'react-router-dom';
import Student from './pages/admin/Student';
import StudentForm from './forms/StudentForm';
import { LoginForm } from './forms/LoginForm';
import { DataProvider } from './context/DataContext';

const ProtectedRoute = ({ element }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? element : <Navigate to="/login" />;
};

const AdminRoute = ({ element }) => {
  const { isAuthenticated, role } = useAuth();

  if (!isAuthenticated) return <Navigate to="/login" />;
  if (!role === "admin") return <Navigate to="/" />;

  return element;
};

const Layout = () => {
  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <div className="flex-grow">
        <Toaster />
        <Outlet />
      </div>
      <Footer />
    </div>
  );
};

const AdminLayout = () => {
  return (
    <div className="mt-16 p-0 w-full flex flex-col min-h-screen">
      <AdminNavbar />
      <div className="flex-grow">
        <Toaster />
        <Outlet />
      </div>
      <Footer />
    </div>
  );
};

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <DataProvider>
          <Router>
            <Routes>
              <Route element={<Layout />}>
                <Route path="/" element={<Homepage />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/exam" element={<ProtectedRoute element={<Exam />} />} />
                <Route path="/profile" element={<ProtectedRoute element={<Profile />} />} />
                <Route path="/result" element={<ProtectedRoute element={<Result />} />} />
              </Route>

              {/* Admin routes */}
              <Route element={<AdminLayout />}>
                <Route path="/admin" element={<AdminRoute element={<Dashboard />} />} />
                <Route path="/admin/exam" element={<AdminRoute element={<Dashboard />} />} />
                <Route path="/admin/result" element={<AdminRoute element={<Dashboard />} />} />
                <Route path="/admin/students" element={<AdminRoute element={<Student />} />} />
                <Route path="/admin/students/create" element={<AdminRoute element={<StudentForm />} />} />
              </Route>
            </Routes>
          </Router>
        </DataProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
