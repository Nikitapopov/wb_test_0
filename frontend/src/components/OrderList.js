import React, {useEffect, useState} from 'react';
import { FormControl, InputLabel, Select, MenuItem } from '@mui/material';

const backendUrl = "http://127.0.0.1:8080"

const OrderList = ({id, setId}) => {
  const [data, setData] = useState([]);

  useEffect(() => {
    const dataFetch = async () => {
      const data = await (
        await fetch(backendUrl + '/ordersIdsList')
      ).json();

      if (!data.error) {
        setData(data.data);
      }
      console.log("data:", data);
    };

    dataFetch();
  }, []);

  const handleChange = (event) => {
    setId(event.target.value)
  }

  return (
    <FormControl fullWidth>
      <InputLabel id="demo-simple-select-label">Orders</InputLabel>
      <Select
        labelId="demo-simple-select-label"
        id="demo-simple-select"
        value={id}
        label="Age"
        onChange={handleChange}
      >
        {data.map(item => {
          return <MenuItem value={item}>{item}</MenuItem>
        })}
      </Select>
    </FormControl>
  )
}

export {OrderList};