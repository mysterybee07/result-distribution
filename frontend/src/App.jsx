import './App.css'
import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom'
import Homepage from './pages/Homepage'
import Login from './pages/Login'
import Register from './pages/Register'
import Exam from './pages/Exam'
import Profile from './pages/Profile'
import Result from './pages/Result'
import Footer from './components/Footer'
import Navbar from './components/Navbar'
import { AuthProvider } from './context/AuthContext'
import { Query, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Dashboard from './pages/admin/Dashboard'
import AdminNavbar from './components/AdminNavbar'
import { Toaster } from '@/components/ui/toaster'

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
    <div className="m-0 p-0 w-full flex flex-col min-h-screen">
      <AdminNavbar />
      <div className="flex-grow">
        <Toaster />
        <Outlet />
      </div>
      <Footer />
    </div>
  )
};

// const Layout = ({ children }) => {
//   return (
//     <div>
//       <Navbar /> {/* Render the Navbar here */}
//       <main>{children}</main> {/* Render the rest of the app here */}
//       <Footer />
//     </div>
//   );
// };

const queryClient = new QueryClient();

function App() {

  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Router>
          <Routes>
            <Route element={<Layout />}>
              <Route path="/" element={<Homepage />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/exam" element={<Exam />} />
              <Route path="/profile" element={<Profile />} />
              <Route path="/result" element={<Result />} />
            </Route>

            {/* admin route */}
            <Route element={<AdminLayout />}>
              <Route path="/admin" element={<Dashboard />} />
              <Route path="/admin/exam" element={<Dashboard />} />
              <Route path="/admin/result" element={<Dashboard />} />

            </Route>
          </Routes>
        </Router>
      </AuthProvider>
    </QueryClientProvider>
  )
}

export default App
