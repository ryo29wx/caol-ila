import { FC } from 'react';
import { Box } from '@mui/material';
import RecommendedUsers from '../components/RecommendedUsers';
import DashboardBottomNavigation from '../components/DashboardBottomNavigation';
import Appbar from '../components/Appbar';

const Dashboard: FC = () => {
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Appbar />
      <RecommendedUsers />
      <DashboardBottomNavigation />
    </Box>
  );
};

export default Dashboard;