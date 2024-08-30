import { useEffect, useState } from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardMedia from '@mui/material/CardMedia';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Grid';
import { CardActionArea, Box } from '@mui/material';
import { Link } from 'react-router-dom';
import FavoriteIcon from './RecommendedUserFavoriteIcon';

type User = {
  user_id: string;
  nick_name: string;
  sex: number;
  title: string;
  company: string;
  like: boolean;
  image_url: string;
};

interface SearchApiResponse {
  users: User[];
  total: number;
}

function RecommendedUsers() {
    const [users, setUsers] = useState<User[]>([]);
    const [page] = useState(1); // ページネーションの現在のページ
    // const apiUrl = process.env.REACT_APP_SEARCH_API;
    const apiUrl = 'http://localhost:8080';
    
    useEffect(() => {
      const fetchData = async () => {
        try {
          const response = await fetch(`${apiUrl}/v1/search?q=sample&p=1`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
          });
          const data: SearchApiResponse = await response.json();
          setUsers(data.users);
        } catch (error) {
          console.error('Error fetching recommended users:', error);
        }
      };
  
      fetchData();
    }, [page]);
  
    return (
      <Box sx={{ flex: 1 }}>
        <Grid container spacing={3}>
            {users.map((user) => (
                <Grid item xs={12} sm={6} md={4} key={user.user_id}>
                  <Card>
                    <CardActionArea component={Link} to={`/user/${user.user_id}`} >
                      <CardMedia
                      component="img"
                      height="200"
                      image={user.image_url}
                      alt={user.title}
                      />
                      <CardContent>
                        <Typography variant="h6" component="div">
                            {user.nick_name} 
                            <br></br>
                            {user.title} at {user.company}
                        </Typography>
                      </CardContent>
                    </CardActionArea>
                    <FavoriteIcon />
                  </Card>
                </Grid>
            ))}
        </Grid>
      </Box>
    );
};
  

export default RecommendedUsers;

