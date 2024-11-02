import React from 'react';
import keycloak from './../Keycloak';

const UserInfo: React.FC = () => {
  const handleLogout = () => {
    keycloak.logout();
  };

  return (
    <div>
      <h1>Welcome, {keycloak.tokenParsed?.preferred_username}</h1>
      <button onClick={handleLogout}>Logout</button>
    </div>
  );
};

export default UserInfo;