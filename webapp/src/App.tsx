import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import './App.css';

// Layout and Page Components
import MainLayout from './components/layout/MainLayout';
import PublicLayout from './components/layout/PublicLayout';
import { HomePage } from './pages/HomePage';
import { ProfilePage } from './pages/ProfilePage';
import { HealthPage } from './pages/HealthPage';
import { NotFoundPage } from './pages/NotFoundPage';

// Authentication components
import { useAuth } from './contexts/AuthContext';

// A component for public routes that don't require authentication
const PublicLayoutRoutes = () => {
    const { isAuthenticated } = useAuth();
    
    // If user is authenticated, redirect to protected layout for better UX
    if (isAuthenticated) {
        return (
            <MainLayout>
                <Outlet />
            </MainLayout>
        );
    }
    
    return (
        <PublicLayout>
            <Outlet />
        </PublicLayout>
    );
};

// A component to group routes that require authentication
const ProtectedLayoutRoutes = () => {
    const { isAuthenticated, isLoading } = useAuth();

    // Show loading state while checking authentication
    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-100">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
                    <p className="text-gray-600">Loading...</p>
                </div>
            </div>
        );
    }

    // Redirect to login if not authenticated
    if (!isAuthenticated) {
        window.location.href = '/_/login';
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-100">
                <div className="text-center">
                    <p className="text-gray-600 mb-4">Redirecting to login...</p>
                    <a href="/login" className="text-blue-500 hover:text-blue-700">
                        Click here if not redirected automatically
                    </a>
                </div>
            </div>
        );
    }

    // User is authenticated
    return (
        <MainLayout>
            <Outlet />
        </MainLayout>
    );
};


function App() {
    return (
        <BrowserRouter basename="/">
            <Routes>
                {/* Public routes that don't require authentication */}
                <Route element={<PublicLayoutRoutes />}>
                    <Route path="/home" element={<HomePage />} />
                    <Route path="/health" element={<HealthPage />} />
                </Route>

                {/* Protected routes that require authentication */}
                <Route element={<ProtectedLayoutRoutes />}>
                    <Route path="/profile" element={<ProfilePage />} />
                </Route>

                {/* Root path redirect to /home */}
                <Route path="/" element={<Navigate replace to="/home" />} />

                {/* Catch-all for unmatched routes */}
                <Route path="*" element={<NotFoundPage />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;
