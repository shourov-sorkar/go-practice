import { Routes as RouterRoutes, Route } from 'react-router-dom';
import { lazy, Suspense } from 'react';

// Lazy load components
const Home = lazy(() => import('../pages/Home'));
const Login = lazy(() => import('../pages/Login'));

function Routes() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <RouterRoutes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
      </RouterRoutes>
    </Suspense>
  );
}

export default Routes; 