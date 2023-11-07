import React, { FC } from 'react';
import RecommendedUsers from '../components/RecommendedUsers';
import DashboardBottomNavigation from '../components/DashboardBottomNavigation';
import Appbar from '../components/Appbar';
import BelowAppbarSpace from '../components/BelowAppbarSpace';

const Dashboard: FC = () => {
  return (
    <div>
      <Appbar />
      <BelowAppbarSpace />
      <RecommendedUsers />
      <DashboardBottomNavigation />
    </div>
  );
};

export default Dashboard;