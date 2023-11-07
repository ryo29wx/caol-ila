import React from 'react';
import { useParams } from 'react-router-dom';
import UserDetail from '../components/UserDetail';

const UserDetailsPage: React.FC = () => {
    // React Router の useParams フックを使って URL から userId を取得
    const { userId } = useParams<{ userId: string }>();

    return (
        <div className="user-details-page">
            <h2>User Details Page</h2>
            {/* 数値として userId を UserDetail コンポーネントに渡す */}
            <UserDetail userId={Number(userId)} />
        </div>
    );
};

export default UserDetailsPage;