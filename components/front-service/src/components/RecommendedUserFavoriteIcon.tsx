import { useState } from 'react';
import IconButton from '@mui/material/IconButton';
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';


function AddShoppingCartIcon() {
  const [isFavorite, setIsFavorite] = useState(false);

  const handleIconClick = () => {
    setIsFavorite(!isFavorite);
  };

  return (
    <IconButton onClick={() => handleIconClick()} size="small" aria-label="show 4 new mails" color="inherit">
      {isFavorite ? <FavoriteIcon color="error" /> : <FavoriteBorderIcon />}
    </IconButton>
  );
}

export default AddShoppingCartIcon;