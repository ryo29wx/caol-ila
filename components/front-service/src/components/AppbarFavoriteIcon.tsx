import { useEffect, useState } from 'react';
import IconButton from '@mui/material/IconButton';
import Badge from '@mui/material/Badge';
import FavoriteIcon from '@mui/icons-material/Favorite';

interface NewLike {
  amount: number;
}

function AppbarFavoriteIcon() {
  const [newLikeAmount, setNewLikeAmount] = useState<number>();

  // const apiUrl = process.env.REACT_APP_ACCOUNT_API;
  const apiUrl = 'http://localhost:8080';

  const handleIconClick = () => {
    window.open('/like/get', '_blank');
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`${apiUrl}/v1/like/req?u=hoge`, {
          method: 'GET',
          headers: {
              'Content-Type': 'application/json'
          }
        });
        const like: NewLike = await response.json();
        setNewLikeAmount(like.amount);
      } catch (error) {
        console.error('Error fetching the amount new likes:', error);
      }
    };

    fetchData();
  }, []);


  return (
    <div>
      <IconButton onClick={() => handleIconClick()} size="large" aria-label="show new likes" color="inherit">
        <Badge badgeContent={newLikeAmount} color="error">
          <FavoriteIcon />
        </Badge>
      </IconButton>
    </div>
  );
}

export default AppbarFavoriteIcon;