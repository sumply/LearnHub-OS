import { Router, Route } from '@solidjs/router';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Library from './pages/Library';
import Tasks from './pages/Tasks';
import Flashcards from './pages/Flashcards';
import Quiz from './pages/Quiz';
import ProtectedRoute from './components/ProtectedRoute';
import StudentDashboard from './pages/StudentDashboard';
import ParentDashboard from './pages/ParentDashboard';
import TeacherDashboard from './pages/TeacherDashboard';
import AdminDashboard from './pages/AdminDashboard';
import Profile from './pages/Profile';

const AppRoutes = () => (
  <Router>
    <Route path="/" component={Home} />
    <Route path="/login" component={Login} />
    <Route path="/register" component={Register} />
    <Route 
      path="/dashboard" 
      component={() => (
        <ProtectedRoute>
          <Dashboard />
        </ProtectedRoute>
      )} 
    />
    <Route 
      path="/library" 
      component={Library}
    />
    <Route 
      path="/materials" 
      component={() => <div style={{'text-align':'center','margin-top':'2rem'}}>Страница перемещена. Перейдите в <a href='/library'>Библиотеку</a>.</div>}
    />
    <Route 
      path="/tasks" 
      component={() => (
        <ProtectedRoute>
          <Tasks />
        </ProtectedRoute>
      )} 
    />
    <Route 
      path="/flashcards" 
      component={() => (
        <ProtectedRoute>
          <Flashcards />
        </ProtectedRoute>
      )} 
    />
    <Route 
      path="/quiz" 
      component={() => (
        <ProtectedRoute>
          <Quiz />
        </ProtectedRoute>
      )} 
    />
    <Route 
      path="/profile" 
      component={() => (
        <ProtectedRoute>
          <Profile />
        </ProtectedRoute>
      )} 
    />
    {/* Дашборды по ролям */}
    <Route 
      path="/student" 
      component={() => (
        <ProtectedRoute role="student">
          <StudentDashboard />
        </ProtectedRoute>
      )} 
    />
    <Route 
      path="/parent" 
      component={() => (
        <ProtectedRoute role="parent">
          <ParentDashboard />
        </ProtectedRoute>
      )} 
    />
    <Route 
      path="/teacher" 
      component={() => (
        <ProtectedRoute role="teacher">
          <TeacherDashboard />
        </ProtectedRoute>
      )} 
    />
    <Route 
      path="/admin" 
      component={() => (
        <ProtectedRoute role="admin">
          <AdminDashboard />
        </ProtectedRoute>
      )} 
    />
    <Route 
      path="/admin/profile/:id" 
      component={() => (
        <ProtectedRoute role="admin">
          <Profile />
        </ProtectedRoute>
      )} 
    />
  </Router>
);

export default AppRoutes; 