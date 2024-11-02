import Keycloak from 'keycloak-js';

const keycloak = new Keycloak({
  url: 'http://localhost:8080/',
  realm: 'caolila-realm-dev',
  clientId: 'caolila-react-client',
});

export default keycloak;