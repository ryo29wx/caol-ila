import React from 'react';
import keycloak from './../Keycloak';


const Login: React.FC = () => {
  const handleLogin = () => {
    keycloak.login();
  };

  return (
    <div>
      <h2>Welcom to Caol-Ila!!</h2>
      <h3>Let's create your world!</h3>
      <button onClick={handleLogin}>Login / Create your account</button>
    </div>
  );
};

export default Login;