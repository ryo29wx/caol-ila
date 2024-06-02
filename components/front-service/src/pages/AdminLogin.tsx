import React, { useState, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, TextField, Button, Grid, Typography } from '@mui/material';

const AdminLoginPage: React.FC = () => {
    const [password, setPassword] = useState('');
    const [userID, setUserID] = useState('');
    const navigate = useNavigate();

    const handleLogin = (e: FormEvent) => {
        e.preventDefault();
        
        if (password == "1234" ) {
            navigate('/admin');
        } else {
            alert('Invalid username or password');
        }
    };

    return (
        <Container maxWidth="sm">
        <Grid container spacing={3} direction="column" alignItems="center" justifyContent="center" style={{ minHeight: '100vh' }}>
          <Grid item>
            <Typography variant="h4">Admin Login</Typography>
          </Grid>
          <Grid item>
            <TextField
              label="Admin ID"
              variant="outlined"
              fullWidth
              value={userID}
              onChange={(e) => setUserID(e.target.value)}
            />
          </Grid>
          <Grid item>
            <TextField
              label="Password"
              variant="outlined"
              type="password"
              fullWidth
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </Grid>
          <Grid item>
            <Button variant="contained" color="primary" onClick={handleLogin} fullWidth>
              Admin Login
            </Button>
          </Grid>
        </Grid>
      </Container>
    );
};

export default AdminLoginPage;