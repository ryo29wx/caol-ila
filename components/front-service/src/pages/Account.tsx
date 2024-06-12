import React, { useEffect, useState } from 'react';
import DashboardBottomNavigation from '../components/DashboardBottomNavigation';
import Appbar from '../components/Appbar';
import axios from 'axios';
import { Container, Typography, CircularProgress } from '@mui/material';

interface MyAccountData {
    id: number;
    name: string;
    email: string;
}

const Account: React.FC = () => {

    const [data, setData] = useState<MyAccountData | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
  
    useEffect(() => {
      // APIリクエストを送信
      axios.get<MyAccountData>('https://api.example.com/user/1')
        .then(response => {
          setData(response.data);
          setLoading(false);
        })
        .catch(error => {
          setError(error.message);
          setLoading(false);
        });
    }, []);
  
    if (loading) {
      return <CircularProgress />;
    }
  
    if (error) {
      return <Typography color="error">{error}</Typography>;
    }
  
    return (
      <Container>
      <Appbar />
        <Typography variant="h4">User Details</Typography>
        {data && (
          <>
            <Typography variant="body1">ID: {data.id}</Typography>
            <Typography variant="body1">Name: {data.name}</Typography>
            <Typography variant="body1">Email: {data.email}</Typography>
          </>
        )}
      <DashboardBottomNavigation />
      </Container>
    );
};

export default Account;