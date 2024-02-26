import React, { useEffect, useState } from 'react';

// ユーザーの詳細情報を保持する型
type UserDetails = {
    id: number;
    name: string;
    photoUrl: string;
    bio: string;
    // 他の必要なフィールド...
};

interface UserDetailProps {
    userId: number; // or string depending on your ID type
}

const UserDetail: React.FC<UserDetailProps> = ({ userId }) => {
    const [userDetails, setUserDetails] = useState<UserDetails | null>(null);

    useEffect(() => {
        const fetchUserData = async () => {
          try {
            const response = await fetch(`/v1/service/userdetail/${userId}`);
            if (response.ok) {
              const userDetail = await response.json();
              setUserDetails(userDetail);
            } else {
              console.error('Failed to fetch user data');
            }
          } catch (error) {
            console.error('Error fetching user data:', error);
            
          }
        };
    
        fetchUserData();
    
      }, [userId]);

    if (!userDetails) {
        return <div>Something wrong...</div>;
    }

    return (
        <div className="user-detail">
            <img src={userDetails.photoUrl} alt={`${userDetails.name}'s photo`} />
            <h1>{userDetails.name}</h1>
            <p>{userDetails.bio}</p>
            {/* 他の詳細情報の表示 */}
        </div>
    );
};

export default UserDetail;