import IconButton from '@mui/material/IconButton';
import Badge from '@mui/material/Badge';
import ChatIcon from '@mui/icons-material/Chat';

function AppbarChatIcon() {
  return (
    <div>
      <IconButton size="large" aria-label="show 4 new mails" color="inherit">
        <Badge badgeContent={4} color="error">
          <ChatIcon />
        </Badge>
      </IconButton>
    </div>
  );
}

export default AppbarChatIcon;