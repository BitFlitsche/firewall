import React from 'react';
import Box from '@mui/material/Box';

// Flag sprite component using country codes
const CountryFlag = ({ countryCode, size = 24, style = {} }) => {
  if (!countryCode || countryCode.length !== 2) {
    return (
      <Box
        sx={{
          width: size,
          height: size,
          backgroundColor: '#f0f0f0',
          borderRadius: '2px',
          display: 'inline-flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: size * 0.4,
          color: '#999',
          ...style
        }}
      >
        ?
      </Box>
    );
  }

  // Convert country code to flag emoji
  const getFlagEmoji = (code) => {
    const first = code.charCodeAt(0) - 65 + 0x1F1E6;
    const second = code.charCodeAt(1) - 65 + 0x1F1E6;
    return String.fromCodePoint(first, second);
  };

  return (
    <Box
      sx={{
        width: size,
        height: size,
        fontSize: size * 0.8,
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
        borderRadius: '2px',
        overflow: 'hidden',
        ...style
      }}
      title={`${countryCode.toUpperCase()} flag`}
    >
      {getFlagEmoji(countryCode.toUpperCase())}
    </Box>
  );
};

export default CountryFlag; 