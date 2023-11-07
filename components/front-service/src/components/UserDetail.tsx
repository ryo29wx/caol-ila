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

    // ユーザー詳細をフェッチする関数
    const fetchUserDetails = async (id: number) => {
        // ここでAPIからユーザー詳細をフェッチします
        // 例: const response = await fetch(`/api/users/${id}`);
        // setUserDetails(await response.json());
        
        // この例ではダミーデータを使用
        setUserDetails({
            id,
            name: 'Sample User',
            photoUrl: 'path_to_photo',
            bio: 'Sample bio text.',
            // 他のフィールド...
        });
    };

    useEffect(() => {
        fetchUserDetails(userId);
    }, [userId]);

    if (!userDetails) {
        return <div>Loading...</div>;
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