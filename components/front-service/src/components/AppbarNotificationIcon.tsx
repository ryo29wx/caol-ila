import IconButton from '@mui/material/IconButton';
import Badge from '@mui/material/Badge';
import NotificationIcon from '@mui/icons-material/Notifications';

function AppbarFavoriteIcon() {
  return (
    <div>
      <IconButton size="large" aria-label="show 12 new notifications" color="inherit">
        <Badge badgeContent={12} color="error">
          <NotificationIcon />
        </Badge>
      </IconButton>
    </div>
  );
}

export default AppbarFavoriteIcon;