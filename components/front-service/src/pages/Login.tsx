import React, { useState, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';

const LoginPage: React.FC = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navi = useNavigate();

    const handleSubmit = (e: FormEvent) => {
        e.preventDefault();
        
        // ここでログインロジックを実装します。成功した場合、ユーザーをダッシュボードにリダイレクトするなど。
        if (username && password) {
            // ログイン成功の場合（この部分は実際の認証ロジックに置き換える必要があります）
            navi.push('/dashboard');
        } else {
            // ログイン失敗の処理をここに書く
            alert('Invalid username or password');
        }
    };

    return (
        <div className="login-page">
            <h2>Login</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label htmlFor="username">Username:</label>
                    <input
                        type="text"
                        id="username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="password">Password:</label>
                    <input
                        type="password"
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>
                <button type="submit">Login</button>
            </form>
        </div>
    );
};

export default LoginPage;