import React, { useEffect, useState } from 'react';
import DashboardBottomNavigation from '../components/DashboardBottomNavigation';
import Appbar from '../components/Appbar';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { Avatar, Box, Container, Typography, CircularProgress, Divider, List, ListItem, ListItemAvatar, ListItemButton, ListItemText } from '@mui/material';

interface ChatList {
    user_id: string;
    nick_name: string;
    sex: string;
    title: string;
    company: string;
    like: boolean;
    image_url: string;
}

const ChatList: React.FC = () => {

    const [chats, setChats] = useState<ChatList[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    // const apiUrl = process.env.REACT_APP_ADMIN_API;
    const apiUrl = 'http://localhost:8080';

    useEffect(() => {
      const fetchData = async () => {
        try {
          const response = await axios.get(`${apiUrl}/v1/chat/list`);
          setChats(response.data.users);
          console.log(chats)
        } catch (error) {
          setError('データの取得に失敗しました。');
        } finally {
          setLoading(false);
        }
      };
  
      fetchData();
    }, []);

    const handleClick = (path: string) => {
      navigate(path);
    };
  
    if (loading) {
      return <CircularProgress />;
    }
  
    if (error) {
      return <Typography color="error">{error}</Typography>;
    }
  
    return (
      <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Container>
      <Appbar />
      <List sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}>
        {chats.map((chat) => (
          <React.Fragment>
            <ListItem alignItems="flex-start">
              <ListItemButton onClick={() => handleClick(chat.user_id)}>
                <ListItemAvatar>
                  <Avatar alt={chat.nick_name} src={chat.image_url} />
                </ListItemAvatar>
                <ListItemText
                primary={
                  <React.Fragment>
                  <Typography
                    sx={{ display: 'inline' }}
                    component="span"
                    variant="body2"
                    color="text.primary"
                  >
                  {chat.nick_name}
                  </Typography>
                  <Divider/>
                  </React.Fragment>
                }
                secondary={
                  <React.Fragment>
                  <Typography
                      sx={{ display: 'inline' }}
                      component="span"
                      variant="body2"
                      color="text.secondary"
                  >
                      {chat.company} - {chat.title}
                  </Typography>
                  </React.Fragment>
              }
                />
              </ListItemButton>
            </ListItem>
            <Divider/>
          </React.Fragment>
        ))}
      </List>
      <DashboardBottomNavigation />
      </Container>
      </Box>
    );
};

export default ChatList;