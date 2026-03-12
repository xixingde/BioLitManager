import React from 'react';
import { BrowserRouter, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

const LoginPage = React.lazy(() => import('@/pages/login/LoginPage'));
const HomePage = React.lazy(() => import('@/pages/home/HomePage'));
const PaperListPage = React.lazy(() => import('@/pages/paper/PaperListPage'));
const PaperCreatePage = React.lazy(() => import('@/pages/paper/PaperCreatePage'));
const PaperDetailPage = React.lazy(() => import('@/pages/paper/PaperDetailPage'));
const BatchImportPage = React.lazy(() => import('@/pages/paper/BatchImportPage'));
const BusinessReviewListPage = React.lazy(() => import('@/pages/review/BusinessReviewListPage'));
const PoliticalReviewListPage = React.lazy(() => import('@/pages/review/PoliticalReviewListPage'));
const ReviewPage = React.lazy(() => import('@/pages/review/ReviewPage'));

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  const location = useLocation();

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
};

const Router: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <HomePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/papers"
          element={
            <PrivateRoute>
              <PaperListPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/papers/create"
          element={
            <PrivateRoute>
              <PaperCreatePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/papers/:id"
          element={
            <PrivateRoute>
              <PaperDetailPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/papers/:id/edit"
          element={
            <PrivateRoute>
              <PaperCreatePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/papers/batch-import"
          element={
            <PrivateRoute>
              <BatchImportPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/reviews"
          element={
            <PrivateRoute>
              <HomePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/reviews/business"
          element={
            <PrivateRoute>
              <BusinessReviewListPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/reviews/business/:paperId"
          element={
            <PrivateRoute>
              <ReviewPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/reviews/political"
          element={
            <PrivateRoute>
              <PoliticalReviewListPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/reviews/political/:paperId"
          element={
            <PrivateRoute>
              <ReviewPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/statistics"
          element={
            <PrivateRoute>
              <HomePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/system/users"
          element={
            <PrivateRoute>
              <HomePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/system/projects"
          element={
            <PrivateRoute>
              <HomePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/system/journals"
          element={
            <PrivateRoute>
              <HomePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/system/config"
          element={
            <PrivateRoute>
              <HomePage />
            </PrivateRoute>
          }
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
