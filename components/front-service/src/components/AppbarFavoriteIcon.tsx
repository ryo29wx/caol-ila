import IconButton from '@mui/material/IconButton';
import Badge from '@mui/material/Badge';
import FavoriteIcon from '@mui/icons-material/Favorite';

function AppbarFavoriteIcon() {
  return (
    <div>
      <IconButton size="large" aria-label="show 7 new favorites" color="inherit">
        <Badge badgeContent={7} color="error">
          <FavoriteIcon />
        </Badge>
      </IconButton>
    </div>
  );
}

export default AppbarFavoriteIcon;