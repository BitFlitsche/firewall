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

const EmailForm = () => {
    const [email, setEmail] = useState('');
    const [status, setStatus] = useState('denied');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            await axiosInstance.post('/email', { address: email, status });
            setMessage('Email added successfully');
            setEmail('');
            setRefresh(r => !r); // Toggle refresh to trigger ListView reload
        } catch (error) {
            setError('Error adding email');
        }
    };

    return (
        <Box sx={{ maxWidth: 600, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>Email Management</Typography>
                <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                    <TextField
                        label="Email Address"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
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
                        Add Email
                    </Button>
                </Box>
                {message && <Alert severity="success" sx={{ mt: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
            </Paper>
            <Box className="list-section" sx={{ mt: 4 }}>
                <Paper elevation={2} sx={{ p: 2 }}>
                    <ListView endpoint="/emails" title="Stored Email Addresses" refresh={refresh} />
                </Paper>
            </Box>
        </Box>
    );
};

export default EmailForm;
