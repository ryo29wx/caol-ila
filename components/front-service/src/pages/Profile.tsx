import React, { useState, useEffect} from 'react';
import { Card, CardContent, Typography, Box, Grid, Avatar } from '@mui/material';

interface UserData {
  userid: string;
  displayname: string;
  gender: string;
  age: number;
  title: string;
  company: string;
  companyemail: string;
  career: string;
  academic: string;
  description: string;
  mainimage: string;
  imagepath: string;
  registday: string;
  lastlogin: string;
}

const Profile: React.FC = () => {
  const [userData, setUserData] = useState<UserData | null>(null);

  // const apiUrl = process.env.REACT_APP_ADDITEM_API;
  const apiUrl = 'http://localhost:50051';
  const userid = '550e8400-e29b-41d4-a716-446655440000'; // for test 

  useEffect(() => {
    // データ取得
    const fetchData = async () => {
      try {
        const response = await fetch(`${apiUrl}/v1/profile?u=${userid}`);
        if (response.ok) {
          const data = await response.json();
          setUserData(data);
        } else {
          console.error("データの取得に失敗しました。");
        }
      } catch (error) {
        console.error("エラー:", error);
      }
    };

    fetchData();
  }, []);

  return (
    <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
      {userData && (
        <Card sx={{ maxWidth: 400, padding: 2 }}>
          <CardContent>
            <Grid container spacing={2} alignItems="center">
              <Grid item>
                <Avatar sx={{ bgcolor: 'primary.main', width: 56, height: 56 }}>
                  {userData.displayname.charAt(0)}
                </Avatar>
              </Grid>
              <Grid item xs>
                <Typography variant="h5" component="div">
                  {userData.displayname}
                </Typography>
                <Typography color="text.secondary">
                  {userData.companyemail}
                </Typography>
              </Grid>
            </Grid>
            <Box mt={2}>
              <Typography variant="body1" color="text.primary">
                年齢: {userData.age}歳
              </Typography>
              <Typography variant="body2" color="text.secondary" mt={1}>
                {userData.description}
              </Typography>
            </Box>
          </CardContent>
        </Card>
      )}
    </Box>
  );
};

export default Profile;