import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Grid, Button, Container } from '@mui/material';

const AdminLoginPage: React.FC = () => {
    const navigate = useNavigate();

    const handleCreateUserBottun = () => {
        navigate('/admin/account/create');
    };

    const handleEditUserBottun = () => {
        navigate('/admin/account/edit');
    };

    const handleDeleteUserBottun = () => {
        navigate('/admin/account/delete');
    };

    return (
        <Container>
          <h1> Admin Page</h1>
        <Grid container direction="column" spacing={2}>
          <Grid item >
            <Button onClick={handleCreateUserBottun} fullWidth variant="contained" color="primary">
              Create User Account
            </Button>
          </Grid>
          <Grid item>
            <Button onClick={handleEditUserBottun} fullWidth variant="contained" color="primary">
              Edit User Account
            </Button>
          </Grid>
          <Grid item>
            <Button onClick={handleDeleteUserBottun} fullWidth variant="contained" color="error">
              Delete User Account
            </Button>
          </Grid>
        </Grid>
      </Container>
    );
  };

export default AdminLoginPage;