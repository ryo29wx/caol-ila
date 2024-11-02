import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Container, TextField, Button, IconButton, InputAdornment, Grid, Typography } from '@mui/material';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';

const LoginPage: React.FC = () => {
    const navigate = useNavigate();
    const [showPassword, setShowPassword] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const [user, setUser] = useState({
      mail: '',
      password: ''
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setUser({ ...user, [e.target.name]: e.target.value });
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

          console.log(user); // Success handling
          const response = await fetch(`${apiUrl}/v1/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(user)
          });
  
          if (!response.ok) {
            setErrorMessage('Passwords or Name do not match');
            throw new Error('Something went wrong');
          }
  
          const apiStatus = await response.json();
          console.log(apiStatus); // Success handling
          navigate('/');
          
        } catch (error) {
          console.error(error); // Error handling
        }
    };

    const handleCreateAccount = () => {
        navigate('/account/create');
    };

    return (
        <Container maxWidth="sm">
        <Box component="form" onSubmit={handleLogin} noValidate sx={{ mt: 1 }}>
        <Grid container spacing={3} direction="column" alignItems="center" justifyContent="center" style={{ minHeight: '100vh' }}>
          <Grid item>
            <Typography variant="h4">Caol-Ila Login</Typography>
          </Grid>
          <Grid item>
            <TextField
              margin="normal"
              required
              fullWidth
              id="userid"
              label="Mail Address"
              name="userid"
              autoComplete="userid"
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
        <Box component="form" onSubmit={handleCreateAccount} noValidate sx={{ mt: 1 }}>
            <h4>初めてアカウントを作る場合はこちら！</h4>
            <Grid item>
            <Button type="submit" variant="contained" color="secondary" fullWidth>
                Create Account
            </Button>
            </Grid>
        </Box>
      </Container>
    );
};

export default LoginPage;