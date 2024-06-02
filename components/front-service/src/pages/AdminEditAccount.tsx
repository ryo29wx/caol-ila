import React, { useState, ChangeEvent} from 'react';
import { Button, TextField, Container, Box } from '@mui/material';

type ImageUploadState = {
    file: File | null;
    previewUrl: string | null;
};

const AdminCreateUser = () => {
  const [product, setProduct] = useState({
    product_name: '',
    seller_name: '',
    category: 0,
    price: 0,
    stocks: 0,
    main_image: '',
    images: ['', ''],
    token: ''
  });
  

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setProduct({ ...product, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const apiUrl = process.env.REACT_APP_ADDITEM_API;

    try {
      const response = await fetch(`${apiUrl}/v1/admin/account/edit`, {
        method: 'POST',
        headers: {
          'Content-Type': 'multipart/form-data'
        },
        body: JSON.stringify(product)
      });

      if (!response.ok) {
        throw new Error('Something went wrong');
      }

      const data = await response.json();
      console.log(data); // Success handling
    } catch (error) {
      console.error(error); // Error handling
    }
  };

  const [imageUpload, setImageUpload] = useState<ImageUploadState>({ file: null, previewUrl: null });

  const handleImageChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files[0]) {
      const file = event.target.files[0];

      setImageUpload({
        file: file,
        previewUrl: URL.createObjectURL(file),
      });
    }
  };

  return (
    <Container component="main" maxWidth="xs">
    <h2>create user account! </h2>
      <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
        <TextField
          margin="normal"
          required
          fullWidth
          id="user_name"
          label="User Name"
          name="user_name"
          autoComplete="user_name"
          autoFocus
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          fullWidth
          id="mail_address1"
          label="Main Mail Address"
          name="mail_address1"
          autoComplete="mail_address1"
          autoFocus
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          fullWidth
          id="mail_address2"
          label="Mail Address 2"
          name="mail_address2"
          autoComplete="mail_address2"
          autoFocus
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          fullWidth
          id="mail_address3"
          label="Mail Address 3"
          name="mail_address3"
          autoComplete="mail_address3"
          autoFocus
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          fullWidth
          id="phone_number1"
          label="Main Phone Number"
          name="phone_number1"
          autoComplete="phone_number1"
          autoFocus
          type="number"
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          fullWidth
          id="phone_number2"
          label="Phone Number 2"
          name="phone_number2"
          autoComplete="phone_number2"
          autoFocus
          type="number"
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          fullWidth
          id="phone_number3"
          label="Phone Number 3"
          name="phone_number3"
          autoComplete="phone_number3"
          autoFocus
          type="number"
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id="address1"
          label="Address(Prefecture)"
          name="address1"
          autoComplete="address1"
          autoFocus
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id="address2"
          label="Address(City)"
          name="address2"
          autoComplete="address2"
          autoFocus
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          fullWidth
          id="address3"
          label="Address(Street)"
          name="address3"
          autoComplete="address3"
          autoFocus
          onChange={handleChange}
        />
        <TextField
          margin="normal"
          fullWidth
          id="postcode"
          label="Postcode"
          name="postcode"
          autoComplete="postcode"
          autoFocus
          type="number"
          onChange={handleChange}
        />
        <div>
          <input
            accept="image/*"
            style={{ display: 'none' }}
            id="raised-button-file"
            multiple
            type="file"
            onChange={handleImageChange}
          />
          <label htmlFor="raised-button-file">
            <Button variant="contained" component="span">
              Upload Image
            </Button>
          </label>
          {imageUpload.previewUrl && <img src={imageUpload.previewUrl} alt="Preview" style={{ width: '100%', marginTop: '10px' }} />}
        </div>
        <Button type="submit" fullWidth variant="contained" sx={{ mt: 3, mb: 2 }}>
          Create
        </Button>
      </Box>
    </Container>
  );
};

export default AdminCreateUser;