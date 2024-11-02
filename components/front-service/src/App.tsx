import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import keycloak from './Keycloak';
import Admin from './pages/Admin';
import AdminLogin from './pages/AdminLogin';
import AdminCreateAccount from './pages/AdminCreateAccount';
import AdminEditAccount from './pages/AdminEditAccount';
import AdminDeleteAccount from './pages/AdminDeleteAccount';
import Dashboard from './pages/Dashboard';
import SearchResultDashboard from './pages/SearchResultDashboard';
import Chat from './pages/Chat';
import ChatList from './pages/ChatList';
import FavoriteGetList from './pages/FavoriteGetList';
import FavoritePostList from './pages/FavoritePostList';
import AccountCreate from './pages/AccountCreate';
import ProfileCreate from './pages/ProfileCreate';
import MyAccount from './pages/Account';
import Login from './pages/Login';
import UserDetails from './pages/UserDetails';
import './App.css';

const App: React.FC = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);


  useEffect(() => {
    keycloak.init({ onLoad: 'login-required' })
      .then(authenticated => {
        setIsAuthenticated(authenticated);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to initialize Keycloak:', error);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <div>loggin in...</div>;
  }

  return (
    <Router>
      <Routes>
        <Route path="/admin/login" element={<AdminLogin />} />
        <Route path="/admin" element={<Admin />} />
        <Route path="/admin/account/create" element={<AdminCreateAccount />} />
        <Route path="/admin/account/edit" element={<AdminEditAccount />} />
        <Route path="/admin/account/delete" element={<AdminDeleteAccount />} />
        <Route path="/" element={isAuthenticated ? <Dashboard /> : <Navigate to="/" />} />
        <Route path="/user/:id" element={<UserDetails />} />
        <Route path="/search" element={<SearchResultDashboard />} />
        <Route path="/like/post" element={<FavoritePostList />} />
        <Route path="/like/get" element={<FavoriteGetList />} />
        <Route path="/chat/" element={isAuthenticated ? <ChatList /> : <Navigate to="/login" />} />
        <Route path="/chat/:id" element={<Chat />} />
        <Route path="/account/create" element={<AccountCreate />} />
        <Route path="/account/create/p" element={<ProfileCreate />} />
        <Route path="/account" element={<MyAccount />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </Router>
  );
};

export default App;