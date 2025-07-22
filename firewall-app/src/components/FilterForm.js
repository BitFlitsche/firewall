import React, { useState, useEffect } from 'react';
import axios from '../axiosConfig'; // Importiere die konfigurierte Axios-Instanz
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Alert from '@mui/material/Alert';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import IconButton from '@mui/material/IconButton';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

const FilterForm = () => {
    const [filterData, setFilterData] = useState({
        ip: '',
        email: '',
        user_agent: '',
        country: '',
        asn: '',
        username: ''
    });
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');

    const handleInputChange = (field) => (event) => {
        setFilterData({
            ...filterData,
            [field]: event.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        
        // Validate that IP address is provided
        if (!filterData.ip.trim()) {
            setError('IP address is required');
            return;
        }
        
        try {
            const response = await axios.post('/api/filter', filterData);
            setMessage(`Filter applied! Result: ${JSON.stringify(response.data)}`);
            setFilterData({
                ip: '',
                email: '',
                user_agent: '',
                country: '',
                asn: '',
                username: ''
            });
        } catch (err) {
            // Extract error message from API response
            const errorMessage = err.response?.data?.error || err.message || 'Error applying filter';
            setError(errorMessage);
        }
    };

    return (
        <Box sx={{ maxWidth: 800, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>Apply Filter</Typography>
                <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                    <TextField
                        label="IP Address"
                        value={filterData.ip}
                        onChange={handleInputChange('ip')}
                        placeholder="Enter IP address to filter"
                        fullWidth
                        required
                        helperText="IP address is required for all filter requests"
                    />
                    <TextField
                        label="Email Address"
                        value={filterData.email}
                        onChange={handleInputChange('email')}
                        placeholder="Enter email address to filter"
                        fullWidth
                    />
                    <TextField
                        label="Username"
                        value={filterData.username}
                        onChange={handleInputChange('username')}
                        placeholder="Enter username to filter"
                        fullWidth
                    />
                    <TextField
                        label="User Agent"
                        value={filterData.user_agent}
                        onChange={handleInputChange('user_agent')}
                        placeholder="Enter user agent to filter"
                        fullWidth
                    />
                    <TextField
                        label="Country"
                        value={filterData.country}
                        onChange={handleInputChange('country')}
                        placeholder="Enter country code to filter"
                        fullWidth
                        helperText="Leave empty for automatic Country lookup from IP"
                    />
                    <TextField
                        label="ASN"
                        value={filterData.asn}
                        onChange={handleInputChange('asn')}
                        placeholder="Enter ASN to filter (e.g., AS7922)"
                        fullWidth
                        helperText="Leave empty for automatic ASN lookup from IP"
                    />
                    <Button type="submit" variant="contained" color="primary" size="large">
                        Apply Filter
                    </Button>
                </Box>
                {message && <Alert severity="success" sx={{ mt: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
            </Paper>
        </Box>
    );
};

export default FilterForm;
