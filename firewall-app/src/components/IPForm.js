import React, { useState } from 'react';
import axiosInstance from '../axiosConfig';
import ListView from './ListView';
// MUI imports
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import Button from '@mui/material/Button';
import InputLabel from '@mui/material/InputLabel';
import FormControl from '@mui/material/FormControl';
import Alert from '@mui/material/Alert';

const IPForm = () => {
    const [ip, setIp] = useState('');
    const [status, setStatus] = useState('denied');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            await axiosInstance.post('/ip', {
                Address: ip,
                Status: status
            });
            setMessage('IP address added successfully');
            setIp('');
            setRefresh(r => !r); // Toggle refresh to trigger ListView reload
        } catch (error) {
            setError(error.response?.data?.error || 'Error adding IP address');
        }
    };

    return (
        <Box sx={{ maxWidth: 600, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>IP Address Management</Typography>
                <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                    <TextField
                        label="IP Address"
                        value={ip}
                        onChange={(e) => setIp(e.target.value)}
                        required
                        placeholder="Enter IP address"
                        fullWidth
                    />
                    <FormControl fullWidth>
                        <InputLabel>Status</InputLabel>
                        <Select
                            value={status}
                            label="Status"
                            onChange={(e) => setStatus(e.target.value)}
                        >
                            <MenuItem value="denied">Denied</MenuItem>
                            <MenuItem value="allowed">Allowed</MenuItem>
                            <MenuItem value="whitelisted">Whitelisted</MenuItem>
                        </Select>
                    </FormControl>
                    <Button type="submit" variant="contained" color="primary">
                        Add IP
                    </Button>
                </Box>
                {message && <Alert severity="success" sx={{ mt: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
            </Paper>
            <Box className="list-section" sx={{ mt: 4 }}>
                <Paper elevation={2} sx={{ p: 2 }}>
                    <ListView endpoint="/ips" title="Stored IP Addresses" refresh={refresh} />
                </Paper>
            </Box>
        </Box>
    );
};

export default IPForm;
