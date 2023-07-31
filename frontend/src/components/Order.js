import { Box, Typography } from '@mui/material';

const Order = ({orderBody}) => {
  console.log("orderBody", JSON.stringify(orderBody));
  return (
    <Box sx={{ width: '100%', maxWidth: 500 }}>
      <Typography variant="body1" align='left'>
        {orderBody && <div><pre>{JSON.stringify(orderBody, null, 2) }</pre></div>}
      </Typography>
    </Box>
  )
}

export {Order};