import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import SearchResultDashboard from './pages/SearchResultDashboard';
import ChatList from './pages/ChatList';
import FavoriteGetList from './pages/FavoriteGetList';
import FavoritePostList from './pages/FavoritePostList';
import AccountCreate from './pages/AccountCreate';
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
        <Route path="/fav/post" element={<FavoritePostList />} />
        <Route path="/fav/get" element={<FavoriteGetList />} />
        <Route path="/chat/list" element={<ChatList />} />
        <Route path="/account/create" element={<AccountCreate />} />
        <Route path="/account" element={<MyAccount />} />
      </Routes>
    </Router>
  );
};

export default App;