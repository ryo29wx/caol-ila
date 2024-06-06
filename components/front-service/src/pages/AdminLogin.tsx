import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Container, TextField, Button, IconButton, InputAdornment, Grid, Typography } from '@mui/material';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';

const AdminLoginPage: React.FC = () => {
    const navigate = useNavigate();
    const [showPassword, setShowPassword] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const [admin, setAdmin] = useState({
      name: '',
      password: ''
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      setAdmin({ ...admin, [e.target.name]: e.target.value });
    };

    const handleClickShowPassword = () => {
      setShowPassword(!showPassword);
    };

    const handleMouseDownPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
      event.preventDefault();
    };

    const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        // const apiUrl = process.env.REACT_APP_ADMIN_API;
        const apiUrl = 'http://localhost:8080';

        try {

          console.log(admin); // Success handling
          const response = await fetch(`${apiUrl}/v1/admin/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(admin)
          });
  
          if (!response.ok) {
            setErrorMessage('Passwords or Name do not match');
            throw new Error('Something went wrong');
          }
  
          const apiStatus = await response.json();
          console.log(apiStatus); // Success handling
          navigate('/admin');
          
        } catch (error) {
          console.error(error); // Error handling
        }
    };

    return (
        <Container maxWidth="sm">
        <Box component="form" onSubmit={handleLogin} noValidate sx={{ mt: 1 }}>
        <Grid container spacing={3} direction="column" alignItems="center" justifyContent="center" style={{ minHeight: '100vh' }}>
          <Grid item>
            <Typography variant="h4">Admin Login</Typography>
          </Grid>
          <Grid item>
            <TextField
              margin="normal"
              required
              fullWidth
              id="name"
              label="Name"
              name="name"
              autoComplete="name"
              autoFocus
              onChange={handleChange}
              />
          </Grid>
          <Grid item>
            <TextField
              margin="normal"
              required
              fullWidth
              id="password"
              label="Password"
              variant="outlined"
              autoComplete="password"
              name="password"
              type={showPassword ? 'text' : 'password'}
              onChange={handleChange}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      aria-label="toggle password visibility"
                      onClick={handleClickShowPassword}
                      onMouseDown={handleMouseDownPassword}
                      edge="end"
                    >
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item>
            <Button type="submit" variant="contained" color="primary" fullWidth>
              Login
            </Button>
          </Grid>
          {errorMessage && (
          <Grid item>
            <Typography variant="body2" color="error">
              {errorMessage}
            </Typography>
          </Grid>
        )}
        </Grid>
        </Box>
      </Container>
    );
};

export default AdminLoginPage;