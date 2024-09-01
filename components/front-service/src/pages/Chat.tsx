import React, { useState, useEffect, useRef } from 'react';
import { Box, Button, TextField, List, ListItem, ListItemText } from '@mui/material';

interface Message {
  id: number;
  text: string;
}

const Chat: React.FC = () => {
  const [message, setMessage] = useState<string>('');
  const [messages, setMessages] = useState<Message[]>([]);
  const socket = useRef<WebSocket | null>(null);

  useEffect(() => {
    socket.current = new WebSocket('ws://localhost:50052/ws');

    // Recieve messages from websocket server(go)
    socket.current.onmessage = (event: MessageEvent) => {
      const newMessage: Message = {
        id: messages.length + 1,
        text: event.data,
      };
      setMessages((prevMessages) => [...prevMessages, newMessage]);
    };

    // close websocket when unmount the component
    return () => {
      socket.current?.close();
    };
  }, [messages.length]);

  const sendMessage = () => {
    if (message.trim() && socket.current) {
      // send messages to websocket server
      socket.current.send(message);
      setMessage('');
    }
  };

  return (
    <Box sx={{ width: '100%', maxWidth: 600, margin: 'auto', mt: 5 }}>
      <List sx={{ height: 300, overflowY: 'auto', border: '1px solid #ccc', mb: 2 }}>
        {messages.map((msg) => (
          <ListItem key={msg.id}>
            <ListItemText primary={msg.text} />
          </ListItem>
        ))}
      </List>
      <TextField
        fullWidth
        variant="outlined"
        label="Type your message..."
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        onKeyPress={(e) => {
          if (e.key === 'Enter') sendMessage();
        }}
      />
      <Button
        fullWidth
        variant="contained"
        color="primary"
        sx={{ mt: 2 }}
        onClick={sendMessage}
      >
        Send
      </Button>
    </Box>
  );
};

export default Chat;