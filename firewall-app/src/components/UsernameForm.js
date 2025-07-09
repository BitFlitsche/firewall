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
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';

const statusOptions = [
  { value: 'denied', label: 'Denied' },
  { value: 'allowed', label: 'Allowed' },
  { value: 'whitelisted', label: 'Whitelisted' },
];

const UsernameForm = () => {
  const [username, setUsername] = useState('');
  const [status, setStatus] = useState('denied');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [refresh, setRefresh] = useState(false);
  const [rules, setRules] = useState([]);
  const [editId, setEditId] = useState(null);

  useEffect(() => {
    axios.get('/usernames')
      .then(res => setRules(res.data))
      .catch(() => setRules([]));
  }, [refresh]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');
    setError('');
    try {
      if (editId) {
        await axios.put(`/username/${editId}`, { Username: username, Status: status });
        setMessage('Username rule updated');
      } else {
        await axios.post('/username', { Username: username, Status: status });
        setMessage('Username rule added');
      }
      setUsername('');
      setStatus('denied');
      setEditId(null);
      setRefresh(r => !r);
    } catch (err) {
      setError('Error saving username rule');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Delete this username rule?')) return;
    try {
      await axios.delete(`/username/${id}`);
      setMessage('Username rule deleted');
      setRefresh(r => !r);
    } catch {
      setError('Error deleting username rule');
    }
  };

  const handleEdit = (rule) => {
    setUsername(rule.Username);
    setStatus(rule.Status);
    setEditId(rule.ID);
  };

  return (
    <Box sx={{ maxWidth: 700, mx: 'auto', mt: 4 }}>
      <Paper sx={{ p: 3 }} elevation={3}>
        <Typography variant="h5" gutterBottom>Username Management</Typography>
        <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
          <TextField
            label="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
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
            {statusOptions.map(opt => (
              <MenuItem key={opt.value} value={opt.value}>{opt.label}</MenuItem>
            ))}
          </TextField>
          <Button type="submit" variant="contained" color="primary">
            {editId ? 'Update Rule' : 'Add Rule'}
          </Button>
          {editId && (
            <Button variant="outlined" color="secondary" onClick={() => { setEditId(null); setUsername(''); setStatus('denied'); }}>
              Cancel Edit
            </Button>
          )}
        </Box>
        {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        <Box className="list-section" sx={{ mt: 4 }}>
          <Paper elevation={2} sx={{ p: 2 }}>
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>ID</TableCell>
                    <TableCell>Username</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {rules.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={4} align="center">No username rules</TableCell>
                    </TableRow>
                  ) : (
                    rules.map(rule => (
                      <TableRow key={rule.ID}>
                        <TableCell>{rule.ID}</TableCell>
                        <TableCell>{rule.Username}</TableCell>
                        <TableCell>{rule.Status}</TableCell>
                        <TableCell>
                          <IconButton onClick={() => handleEdit(rule)} size="small"><EditIcon /></IconButton>
                          <IconButton onClick={() => handleDelete(rule.ID)} size="small" color="error"><DeleteIcon /></IconButton>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          </Paper>
        </Box>
      </Paper>
    </Box>
  );
};

export default UsernameForm; 