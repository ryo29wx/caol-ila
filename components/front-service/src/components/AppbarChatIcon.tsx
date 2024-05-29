import { useEffect, useState } from 'react';
import IconButton from '@mui/material/IconButton';
import Badge from '@mui/material/Badge';
import ChatIcon from '@mui/icons-material/Chat';

interface NewChat {
  amount: number;
}

function AppbarChatIcon() {
  const [newCahtAmount, setNewChatAmount] = useState<number>();

  // const apiUrl = process.env.REACT_APP_ACCOUNT_API;
  const apiUrl = 'http://localhost:8080';

  const handleIconClick = () => {
    window.open('/chat/list', '_blank');
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`${apiUrl}/v1/chat/amount?u=hoge`, {
          method: 'GET',
          headers: {
              'Content-Type': 'application/json'
          }
        });
        const chat: NewChat = await response.json();
        setNewChatAmount(chat.amount);
      } catch (error) {
        console.error('Error fetching the amount new chat:', error);
      }
    };

    fetchData();
  }, []);




  return (
    <div>
      <IconButton onClick={() => handleIconClick()} size="large" aria-label="show new chats" color="inherit">
        <Badge badgeContent={newCahtAmount} color="error">
          <ChatIcon />
        </Badge>
      </IconButton>
    </div>
  );
}

export default AppbarChatIcon;