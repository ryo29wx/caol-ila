import Avatar from '@mui/material/Avatar';
import Box from '@mui/system/Box';
import Stack from '@mui/material/Stack';
import Grid from '@mui/material/Unstable_Grid2'; // Grid version 2

function BelowAppbarSpace() {
  const handleClickAvatorIcon = () => {
  };


  return (
    <div>
        <Box sx={{ flexGrow: 1 }}>
            <Grid container spacing={2}>
                <Stack
                direction={{ xs: 'column', sm: 'row' }}
                spacing={{ xs: 1, sm: 2, md: 4 }}
                >
                    <Avatar alt="Remy Sharp" src="/images/TestUser1.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser2.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser3.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser4.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser5.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser6.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser7.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser8.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser9.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                    <Avatar alt="Remy Sharp" src="/images/TestUser10.jpeg" sx={{ width: 112, height: 112 }} onClick={handleClickAvatorIcon}/>
                </Stack>
            </Grid>
        </Box>
    </div>
  );
}

export default BelowAppbarSpace;

