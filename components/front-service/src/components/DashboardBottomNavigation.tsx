import React from 'react';
import axios from 'axios';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import RestoreIcon from '@mui/icons-material/Restore';
import FavoriteIcon from '@mui/icons-material/Favorite';
import LocationOnIcon from '@mui/icons-material/LocationOn';

function DashboardBottomNavigation() {
  const [value, setValue] = React.useState(0);

  const handleIconClick = (apiEndpoint: string) => {
    axios.get(`http://your-backend-url/service/${apiEndpoint}`)
      .then(response => {
        console.log(response.data);
      })
      .catch(error => {
        console.error('There was an error fetching the data!', error);
      });
  };

  return (
    <BottomNavigation
      value={value}
      onChange={(event, newValue) => {
        setValue(newValue);
      }}
    >
      <BottomNavigationAction
        label="Recent"
        value="recent"
        icon={<RestoreIcon />}
        onClick={() => handleIconClick('recent')}
      />
      <BottomNavigationAction
        label="Favorites"
        value="favorites"
        icon={<FavoriteIcon />}
        onClick={() => handleIconClick('favo')}
      />
      <BottomNavigationAction
        label="Nearby"
        value="nearby"
        icon={<LocationOnIcon />}
        onClick={() => handleIconClick('nearby')}
      />
    </BottomNavigation>
  );
}

export default DashboardBottomNavigation;
