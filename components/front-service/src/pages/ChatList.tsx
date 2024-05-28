import React from 'react';
import { useParams } from 'react-router-dom';
import UserDetail from '../components/userdetail/UserDetail';

const ChatList: React.FC = () => {
    const { userId } = useParams<{ userId: string }>();

    return (
        <div className="user-details-page">
            <h2>User Detail</h2>
            <UserDetail userId={Number(userId)} />
        </div>
    );
};

export default ChatList;