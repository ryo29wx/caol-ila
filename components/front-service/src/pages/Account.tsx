import React, { useEffect, useState } from 'react';
import keycloak from './../Keycloak';
import {
  Card,
  CardContent,
  Typography,
  Avatar,
  Grid,
  Box,
  Button,
} from '@mui/material';
import AccountCircle from '@mui/icons-material/AccountCircle';

interface UserProfileProps {
  username?: string;
  firstname?: string;
  lastname?: string;
}

const Account: React.FC = () => {
  const [userProfile, setUserProfile] = useState<UserProfileProps>({});

  const handleLogout = () => {
    keycloak.logout();
  };

  useEffect(() => {
    if (keycloak.authenticated) {
      console.log(keycloak.tokenParsed);
      const { preferred_username, preferred_firstname, preferred_lastname } = keycloak.tokenParsed || {};
      setUserProfile({
        username: preferred_username,
        firstname: preferred_firstname,
        lastname: preferred_lastname,
      });
    }
  }, []);

  return (
    <Box display="flex" justifyContent="center" marginTop={4}>
      <Card sx={{ maxWidth: 400, padding: 2 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            <Grid item>
              <Avatar sx={{ bgcolor: 'primary.main', width: 56, height: 56 }}>
                <AccountCircle fontSize="large" />
              </Avatar>
            </Grid>
            <Grid item xs>
              <Typography variant="h5" component="div">
                {userProfile.firstname || ''} {userProfile.lastname || ''}
              </Typography>
              <Typography color="text.secondary">
                {userProfile.username || 'No username'}
              </Typography>
            </Grid>
          </Grid>
          <Box marginTop={2}>
            <Typography variant="body1" color="text.primary">
              Mail: {userProfile.username ? `${userProfile.username}@example.com` : 'メールアドレスなし'}
            </Typography>
          </Box>
          <Button variant="contained" color="secondary" onClick={handleLogout} sx={{ marginTop: 2 }} >
            Logout
          </Button>
        </CardContent>
      </Card>

    </Box>
  );
};

export default Account;