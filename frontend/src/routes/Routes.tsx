import { Routes as RouterRoutes, Route } from 'react-router-dom';
import { lazy, Suspense } from 'react';
import { ProtectedRoute } from '../components/ProtectedRoute';

// Lazy load components
const HomePage = lazy(() => import('../pages/HomePage').then(module => ({ default: module.HomePage })));
const LoginPage = lazy(() => import('../pages/LoginPage').then(module => ({ default: module.LoginPage })));
const UserProfilePage = lazy(() => import('../pages/UserProfilePage').then(module => ({ default: module.UserProfilePage })));
const SettingsPage = lazy(() => import('../pages/SettingsPage').then(module => ({ default: module.SettingsPage })));
function Routes() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
    <RouterRoutes>
      <Route path="/login" element={<LoginPage />} />
      <Route element={<ProtectedRoute />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/user-profile" element={<UserProfilePage />} />
        <Route path="/settings" element={<SettingsPage />} />
      </Route>
    </RouterRoutes>
  </Suspense>
  );
}

export default Routes; 