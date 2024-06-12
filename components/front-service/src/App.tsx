import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Admin from './pages/Admin';
import AdminLogin from './pages/AdminLogin';
import AdminCreateAccount from './pages/AdminCreateAccount';
import AdminEditAccount from './pages/AdminEditAccount';
import AdminDeleteAccount from './pages/AdminDeleteAccount';
import Dashboard from './pages/Dashboard';
import SearchResultDashboard from './pages/SearchResultDashboard';
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
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/login" element={<Login />} />
        <Route path="/user/:id" element={<UserDetails />} />
        <Route path="/search" element={<SearchResultDashboard />} />
        <Route path="/like/post" element={<FavoritePostList />} />
        <Route path="/like/get" element={<FavoriteGetList />} />
        <Route path="/chat/list" element={<ChatList />} />
        <Route path="/account/create" element={<AccountCreate />} />
        <Route path="/account/create/p" element={<ProfileCreate />} />
        <Route path="/account" element={<MyAccount />} />
        <Route path="/admin/login" element={<AdminLogin />} />
        <Route path="/admin" element={<Admin />} />
        <Route path="/admin/account/create" element={<AdminCreateAccount />} />
        <Route path="/admin/account/edit" element={<AdminEditAccount />} />
        <Route path="/admin/account/delete" element={<AdminDeleteAccount />} />
      </Routes>
    </Router>
  );
};

export default App;