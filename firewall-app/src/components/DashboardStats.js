import React, { useState, useEffect, useCallback } from 'react';
import axios from '../axiosConfig';
import {
  Box,
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  CircularProgress,
  Alert,
  Chip
} from '@mui/material';
import {
  Public as PublicIcon,
  Email as EmailIcon,
  Person as PersonIcon,
  Flag as FlagIcon
} from '@mui/icons-material';

const statusLabels = ['allowed', 'denied', 'whitelisted'];

const DashboardStats = () => {
  const [stats, setStats] = useState({
    ips: { total: 0, allowed: 0, denied: 0, whitelisted: 0 },
    emails: { total: 0, allowed: 0, denied: 0, whitelisted: 0 },
    userAgents: { total: 0, allowed: 0, denied: 0, whitelisted: 0 },
    countries: { total: 0, allowed: 0, denied: 0, whitelisted: 0 }
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const countStatuses = (arr, statusField = 'status') => {
    const counts = { allowed: 0, denied: 0, whitelisted: 0 };
    arr.forEach(item => {
      const status = (item[statusField] || '').toLowerCase();
      if (statusLabels.includes(status)) counts[status]++;
    });
    return counts;
  };

  const fetchStats = useCallback(async () => {
    try {
      setLoading(true);
      setError('');

      // Lade nur die Stats-Endpunkte
      const [ipsRes, emailsRes, userAgentsRes, countriesRes] = await Promise.all([
        axios.get('/ips/stats'),
        axios.get('/emails/stats'),
        axios.get('/user-agents/stats'),
        axios.get('/countries/stats')
      ]);

      setStats({
        ips: ipsRes.data,
        emails: emailsRes.data,
        userAgents: userAgentsRes.data,
        countries: countriesRes.data
      });
    } catch (err) {
      setError('Failed to load statistics');
      console.error('Error fetching stats:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  const getFilterIcon = (type) => {
    switch (type) {
      case 'IP Addresses':
        return <PublicIcon />;
      case 'Email Addresses':
        return <EmailIcon />;
      case 'User Agents':
        return <PersonIcon />;
      case 'Countries':
        return <FlagIcon />;
      default:
        return null;
    }
  };

  const getStatusColor = (count) => {
    if (count === 0) return 'default';
    if (count < 10) return 'success';
    if (count < 50) return 'warning';
    return 'error';
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 900, mx: 'auto', mt: 4 }}>
      <Paper sx={{ p: 3 }} elevation={3}>
        <Typography variant="h5" gutterBottom>
          Filter Statistics
        </Typography>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Filter Type</TableCell>
                <TableCell>Total</TableCell>
                <TableCell>Allowed</TableCell>
                <TableCell>Denied</TableCell>
                <TableCell>Whitelisted</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              <TableRow>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    {getFilterIcon('IP Addresses')}
                    IP Addresses
                  </Box>
                </TableCell>
                <TableCell>{stats.ips.total}</TableCell>
                <TableCell>
                  <Chip label={stats.ips.allowed} color={getStatusColor(stats.ips.allowed)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip label={stats.ips.denied} color={getStatusColor(stats.ips.denied)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip label={stats.ips.whitelisted} color={getStatusColor(stats.ips.whitelisted)} size="small" />
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    {getFilterIcon('Email Addresses')}
                    Email Addresses
                  </Box>
                </TableCell>
                <TableCell>{stats.emails.total}</TableCell>
                <TableCell>
                  <Chip label={stats.emails.allowed} color={getStatusColor(stats.emails.allowed)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip label={stats.emails.denied} color={getStatusColor(stats.emails.denied)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip label={stats.emails.whitelisted} color={getStatusColor(stats.emails.whitelisted)} size="small" />
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    {getFilterIcon('User Agents')}
                    User Agents
                  </Box>
                </TableCell>
                <TableCell>{stats.userAgents.total}</TableCell>
                <TableCell>
                  <Chip label={stats.userAgents.allowed} color={getStatusColor(stats.userAgents.allowed)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip label={stats.userAgents.denied} color={getStatusColor(stats.userAgents.denied)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip label={stats.userAgents.whitelisted} color={getStatusColor(stats.userAgents.whitelisted)} size="small" />
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    {getFilterIcon('Countries')}
                    Countries
                  </Box>
                </TableCell>
                <TableCell>{stats.countries.total}</TableCell>
                <TableCell>
                  <Chip label={stats.countries.allowed} color={getStatusColor(stats.countries.allowed)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip label={stats.countries.denied} color={getStatusColor(stats.countries.denied)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip label={stats.countries.whitelisted} color={getStatusColor(stats.countries.whitelisted)} size="small" />
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </TableContainer>
        <Box sx={{ mt: 3, p: 2, bgcolor: 'background.default', borderRadius: 1 }}>
          <Typography variant="body2" color="text.secondary">
            <strong>Total Filters:</strong> {stats.ips.total + stats.emails.total + stats.userAgents.total + stats.countries.total} entries
          </Typography>
        </Box>
      </Paper>
    </Box>
  );
};

export default DashboardStats; 