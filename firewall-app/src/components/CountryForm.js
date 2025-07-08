import React, { useState, useEffect } from 'react';
import axios from '../axiosConfig';
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Alert from '@mui/material/Alert';
import MenuItem from '@mui/material/MenuItem';
import ListView from './ListView';

const CountryForm = () => {
    const [country, setCountry] = useState('');
    const [status, setStatus] = useState('denied');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            await axios.post('/country', { code: country, status });
            setCountry('');
            setStatus('denied');
            setMessage('Country added successfully');
            setRefresh(r => !r); // Toggle refresh to trigger ListView reload
        } catch (err) {
            setError('Error adding country');
        }
    };

    return (
        <Box sx={{ maxWidth: 700, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>Country Management</Typography>
                <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
                    <TextField
                        label="Country Code"
                        value={country}
                        onChange={(e) => setCountry(e.target.value)}
                        placeholder="Enter Country Code (z.B. DE)"
                        required
                        fullWidth
                    />
                    <TextField
                        select
                        label="Status"
                        value={status}
                        onChange={(e) => setStatus(e.target.value)}
                        fullWidth
                    >
                        <MenuItem value="denied">Denied</MenuItem>
                        <MenuItem value="allowed">Allowed</MenuItem>
                        <MenuItem value="whitelisted">Whitelisted</MenuItem>
                    </TextField>
                    <Button type="submit" variant="contained" color="primary">
                        Add Country
                    </Button>
                </Box>
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
                <Box className="list-section" sx={{ mt: 4 }}>
                    <Paper elevation={2} sx={{ p: 2 }}>
                        <ListView endpoint="/countries" title="Stored Countries" refresh={refresh} />
                    </Paper>
                </Box>
            </Paper>
        </Box>
    );
};

export default CountryForm;
