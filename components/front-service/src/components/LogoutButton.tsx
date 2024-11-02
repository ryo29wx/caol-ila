import React from 'react';
import keycloak from './../Keycloak';

const LogoutButton: React.FC = () => {
  const handleLogout = () => {
    keycloak.logout({
      redirectUri: window.location.origin,
    });
  };

  return (
    <button onClick={handleLogout}>Logout</button>
  );
};

export default LogoutButton;