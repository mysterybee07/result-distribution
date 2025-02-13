import './App.css';
import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import Homepage from './pages/Homepage';
import Login from './pages/Login';
import Register from './pages/Register';
import Exam from './pages/Exam';
import Profile from './pages/Profile';
import Result from './pages/Result';
// import Footer from './components/Footer';
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
import CreateCenter from './pages/college/createCenter';
import ListCenter from './pages/college/listCenter';
import AssignCenter from './pages/exam/ExamSchedule';
import ExamSchedule from './pages/exam/ExamSchedule';
import ListExams from './pages/exam/ListExams';
import ListRoutine from './pages/exam/ListRoutine';
import ListMarks from './pages/marks/AddMarks';
import AddMarks from './pages/marks/AddMarks';

const ProtectedRoute = ({ element }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? element : <Navigate to="/login" />;
};

const AdminRoute = ({ element }) => {
  const { isAuthenticated, role, loading } = useAuth();
  if (loading) return <p>Loading...</p>;
  // console.log("🚀 ~ AdminRoute ~ role:", role)

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
      {/* <Footer /> */}
    </div>
  );
};


// User Routes
const userRoutes = [
  { path: "/", element: <Homepage /> },
  { path: "/login", element: <Login /> },
  { path: "/register", element: <Register /> },
  { path: "/exam", element: <ProtectedRoute element={<Exam />} /> },
  { path: "/profile", element: <ProtectedRoute element={<Profile />} /> },
  { path: "/result", element: <ProtectedRoute element={<Result />} /> },
];

// Admin Routes
const adminRoutes = [
  { path: "/dashboard", element: <Dashboard /> },
  { path: "/admin/result", element: <Dashboard /> },
  // Student
  { path: "/admin/students", element: <Student /> },
  { path: "/admin/students/create", element: <StudentForm /> },
  { path: "/admin/students/edit/:id", element: <StudentForm /> },
  { path: "/admin/students/create/bulk", element: <BulkStudentForm /> },
  // Course
  { path: "/admin/courses", element: <ListCourse /> },
  { path: "/admin/courses/create", element: <CreateCourse /> },
  { path: "/admin/courses/edit/:id", element: <EditCourse /> },
  // College & Center
  { path: "/admin/college", element: <ListCollege /> },
  { path: "/admin/college/create", element: <CreateCollege /> },
  { path: "/admin/center", element: <ListCenter /> },
  { path: "/admin/center/create", element: <CreateCenter /> },
  // Notice
  { path: "/admin/notice", element: <Dashboard /> },
  { path: "/admin/notice/create", element: <CreateNotice /> },
  { path: "/admin/notice/edit/:id", element: <EditNotice /> },
  // Exam
  { path: "/admin/exam", element: <ListExams /> },
  { path: "/admin/exam/routine", element: <ListRoutine /> },
  { path: "/admin/exam/create", element: <ExamSchedule /> },
  // Marks
  { path: "/admin/marks/create", element: <AddMarks /> },
];



const queryClient = new QueryClient();

function App() {
  // const { isAuthenticated } = useAuth();
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <DataProvider>
          <Router>
            <Routes>
              {/* User Routes */}
              <Route element={<Layout />}>
                {userRoutes.map(({ path, element }) => (
                  <Route key={path} path={path} element={element} />
                ))}
              </Route>

              {/* Demo Route */}
              <Route path="/demo" element={<AdminRoute element={<Demo />} />} />

              {/* Admin Routes */}
              <Route element={<AdminLayout />}>
                {adminRoutes.map(({ path, element }) => (
                  <Route key={path} path={path} element={<AdminRoute element={element} />} />
                ))}
              </Route>
            </Routes>
          </Router>
        </DataProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
