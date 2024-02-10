import React from 'react';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import InfoIcon from '@mui/icons-material/Info';
import ContactMailIcon from '@mui/icons-material/ContactMail';

function DashboardBottomNavigation() {
  const [value, setValue] = React.useState(0);

  const handleIconClick = () => {
  };


  return (
    <BottomNavigation
      value={value}
      onChange={(event, newValue) => {
        setValue(newValue);
        console.log(event);
      }}
    >
      <BottomNavigationAction label="About Us" icon={<InfoIcon />} onClick={handleIconClick} />
      <BottomNavigationAction label="Contact Us" icon={<ContactMailIcon />} onClick={handleIconClick} />
      
    </BottomNavigation>
  );
}

export default DashboardBottomNavigation;
