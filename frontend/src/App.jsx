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
import BulkStudentForm from './forms/BulkStudentForm';
import CreateCourse from './pages/courses/CreateCourse';
import ListCourse from './pages/courses/ListCourse';
import CreateCollege from './pages/college/createCollege';
import ListCollege from './pages/college/listCollege';
import CreateNotice from './pages/notice/CreateNotice';
import EditNotice from './pages/notice/EditNotice';
import EditCourse from './pages/courses/EditCourse';
import { Demo } from './pages/Demo';
import AdminLayout from './layout/AuthLayout';

const ProtectedRoute = ({ element }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? element : <Navigate to="/login" />;
};

const AdminRoute = ({ element }) => {
  const { isAuthenticated, role, loading } = useAuth();
  if (loading) return <p>Loading...</p>;
  console.log("🚀 ~ AdminRoute ~ role:", role)

  if (!isAuthenticated) return <Navigate to="/login" />;
  if (role !== "admin") return <Navigate to="/" />;

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



const queryClient = new QueryClient();

function App() {
  // const { isAuthenticated } = useAuth();
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
              <Route path="/demo" element={<AdminRoute element={<Demo />} />} />
              <Route element={<AdminLayout />}>

                <Route path="/dashboard" element={<AdminRoute element={<Dashboard />} />} />
                <Route path="/admin/exam" element={<AdminRoute element={<Dashboard />} />} />
                <Route path="/admin/result" element={<AdminRoute element={<Dashboard />} />} />
                {/* student */}
                <Route path="/admin/students" element={<AdminRoute element={<Student />} />} />
                <Route path="/admin/students/create" element={<AdminRoute element={<StudentForm />} />} />
                <Route path="/admin/students/edit/:id" element={<AdminRoute element={<StudentForm />} />} />
                <Route path="/admin/students/create/bulk" element={<AdminRoute element={<BulkStudentForm />} />} />
                {/* course */}
                <Route path="/admin/courses" element={<AdminRoute element={<ListCourse />} />} />
                <Route path="/admin/courses/create" element={<AdminRoute element={<CreateCourse />} />} />
                <Route path="/admin/courses/edit/:id" element={<AdminRoute element={<EditCourse />} />} />
                {/* college */}
                <Route path='/admin/college' element={<AdminRoute element={<ListCollege />} />} />
                <Route path='/admin/college/create' element={<AdminRoute element={<CreateCollege />} />} />
                {/* Notice */}
                <Route path="/admin/notice" element={<AdminRoute element={<Dashboard />} />} />
                <Route path="/admin/notice/create" element={<AdminRoute element={<CreateNotice />} />} />
                <Route path="/admin/notice/edit/:id" element={<AdminRoute element={<EditNotice />} />} />
              </Route>
            </Routes>
          </Router>
        </DataProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
