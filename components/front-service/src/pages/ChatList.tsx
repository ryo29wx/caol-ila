import React from 'react';
import DashboardBottomNavigation from '../components/DashboardBottomNavigation';
import Appbar from '../components/Appbar';

const ChatList: React.FC = () => {

    return (
        <div>
            <Appbar />
            <h1> Maintenance Now. </h1>
            <DashboardBottomNavigation />
        </div>
    );
};

export default ChatList;