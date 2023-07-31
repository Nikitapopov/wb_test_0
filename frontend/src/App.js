import React, {useEffect, useState} from 'react';
import './App.css';
import {OrderList} from './components/OrderList';
import {Order} from './components/Order';
import { Box } from '@mui/material';

const backendUrl = "http://127.0.0.1:8080"

function App() {
  const [id, setId] = useState('');
  const [order, setOrder] = useState('');
  useEffect(() => {
    const dataFetch = async () => {
      const data = await (
        await fetch(`${backendUrl}/order/${id}`)
      ).json();

      if (!data.error) {
        setOrder(data.data);
      }
      console.log("data:", data);
    };

    dataFetch();
  }, [id]);

  return (
    <div className="App">
      <Box sx={{ width: '100%', maxWidth: 500, padding: 3 }}>
        <OrderList id={id} setId={setId}/> 
        <Order orderBody={order}/> 
      </Box>
    </div>
  );
}

export default App;
