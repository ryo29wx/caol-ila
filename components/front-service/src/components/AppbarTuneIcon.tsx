import React, { useState } from 'react';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import Button from '@mui/material/Button';
import FormControl from '@mui/material/FormControl';
import FormControlLabel from '@mui/material/FormControlLabel';
import Radio from '@mui/material/Radio';
import IconButton from '@mui/material/IconButton';
import RadioGroup from '@mui/material/RadioGroup';
import Checkbox from '@mui/material/Checkbox';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import InputLabel from '@mui/material/InputLabel';
import FormGroup from '@mui/material/FormGroup';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import TuneSharpIcon from '@mui/icons-material/TuneSharp';

// ダイアログ内のフォームアイテムの初期値を設定
const defaultFilterValues = {
  sex: '',
  jobTitles: [],
  liveAt: '',
  ageGroup: []
};

function AppbarTuneIcon() {
  const [open, setOpen] = useState(false);
  const [filters, setFilters] = useState(defaultFilterValues);
  const menuId = 'primary-search-account-menu';

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleSexChange = (event) => {
    setFilters({ ...filters, sex: event.target.value });
  };

  const handleJobTitleChange = (event) => {
    const jobTitles = [...filters.jobTitles];
    if (event.target.checked) {
      jobTitles.push(event.target.name);
    } else {
      const index = jobTitles.indexOf(event.target.name);
      if (index > -1) {
        jobTitles.splice(index, 1);
      }
    }
    setFilters({ ...filters, jobTitles: jobTitles });
  };

  const handleLiveAtChange = (event) => {
    setFilters({ ...filters, liveAt: event.target.value });
  };

  const handleAgeGroupChange = (event) => {
    const ageGroup = [...filters.ageGroup];
    if (event.target.checked) {
      ageGroup.push(event.target.name);
    } else {
      const index = ageGroup.indexOf(event.target.name);
      if (index > -1) {
        ageGroup.splice(index, 1);
      }
    }
    setFilters({ ...filters, ageGroup: ageGroup });
  };

  // 送信処理（実際にはAPIへのリクエストなど）
  const handleSubmit = () => {
    // フィルターを適用する処理
    console.log(filters);
    handleClose();
  };

  return (
    <div>
      <IconButton
              size="large"
              edge="end"
              aria-label="account of current user"
              aria-controls={menuId}
              aria-haspopup="true"
              onClick={handleClickOpen}
              color="inherit"
            >
        <TuneSharpIcon />
      </IconButton>
      <Dialog open={open} onClose={handleClose} fullWidth="true" maxWidth="md">
        <DialogTitle>フィルターする</DialogTitle>
        <DialogContent>
          <Box display="flex" justifyContent="space-around" flex="1 0 auto">
            <Typography variant="h6">性別</Typography>
            <FormControl component="fieldset">
              <RadioGroup name="sex" value={filters.sex} onChange={handleSexChange}>
                <FormControlLabel value="female" control={<Radio />} label="Female" />
                <FormControlLabel value="male" control={<Radio />} label="Male" />
              </RadioGroup>
            </FormControl>
          </Box>
          <Box display="flex" justifyContent="space-around" flex="1 0 auto">
            <Typography variant="h6">職種</Typography>
            <FormControl component="fieldset">
              <FormGroup>
                {['SWE', 'SRE', 'Manager', 'CxE', 'HR', 'Coordinator', 'DataScientist', 'Sales', 'InsideSales', 'CustomerSuccess'].map((jobTitle) => (
                  <FormControlLabel
                    key={jobTitle}
                    control={<Checkbox checked={filters.jobTitles.includes(jobTitle)} onChange={handleJobTitleChange} name={jobTitle} />}
                    label={jobTitle}
                  />
                ))}
              </FormGroup>
            </FormControl>
          </Box>
          <Box display="flex" justifyContent="space-around" flex="1 0 auto">
          <Typography variant="h6">住んでいるエリア</Typography>
            <FormControl fullWidth>
              <InputLabel>Live at</InputLabel>
              <Select value={filters.liveAt} onChange={handleLiveAtChange}>
                {/* 都道府県データをループで生成 */}
                {prefectures.map((prefecture) => (
                  <MenuItem key={prefecture} value={prefecture}>
                    {prefecture}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
          <Box display="flex" justifyContent="space-around" flex="1 0 auto">
            <Typography variant="h6">年齢</Typography>
            <FormControl component="fieldset">
              <FormGroup>
                {['20s', '30s', '40s', '50s+'].map((age) => (
                  <FormControlLabel
                    key={age}
                    control={<Checkbox checked={filters.ageGroup.includes(age)} onChange={handleAgeGroupChange} name={age} />}
                    label={age}
                  />
                ))}
              </FormGroup>
            </FormControl>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button onClick={handleSubmit}>Apply</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}

export default AppbarTuneIcon;

const prefectures = [
  '東京都', '大阪', '愛知', '福岡'
];