import { FC } from 'react';
import SearchResultUsers from '../components/SearchResultUsers';
import DashboardBottomNavigation from '../components/DashboardBottomNavigation';
import Appbar from '../components/Appbar';
import BelowAppbarSpace from '../components/BelowAppbarSpace';

const Dashboard: FC = () => {
  return (
    <div>
      <Appbar />
      <BelowAppbarSpace />
      <SearchResultUsers />
      <DashboardBottomNavigation />
    </div>
  );
};

export default Dashboard;