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
  Flag as FlagIcon,
  Dns as DnsIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
// SystemStats import entfernen
// import SystemStats from './SystemStats';

const statusLabels = ['allowed', 'denied', 'whitelisted'];

// Extrahiere die Filter-Statistik-Tabelle in eine eigene Komponente
const FilterStats = ({ stats, error, handleCountClick, getFilterIcon, getStatusColor }) => (
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
                {getFilterIcon('Charset Rules')}
                Charset Rules
              </Box>
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('charsets', null)}>{stats.charsets.total}</TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('charsets', 'allowed')}>
              <Chip label={stats.charsets.allowed} color={getStatusColor(stats.charsets.allowed)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('charsets', 'denied')}>
              <Chip label={stats.charsets.denied} color={getStatusColor(stats.charsets.denied)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('charsets', 'whitelisted')}>
              <Chip label={stats.charsets.whitelisted} color={getStatusColor(stats.charsets.whitelisted)} size="small" />
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                {getFilterIcon('Countries')}
                Countries
              </Box>
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('countries', null)}>{stats.countries.total}</TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('countries', 'allowed')}>
              <Chip label={stats.countries.allowed} color={getStatusColor(stats.countries.allowed)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('countries', 'denied')}>
              <Chip label={stats.countries.denied} color={getStatusColor(stats.countries.denied)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('countries', 'whitelisted')}>
              <Chip label={stats.countries.whitelisted} color={getStatusColor(stats.countries.whitelisted)} size="small" />
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                {getFilterIcon('Email Addresses')}
                Email Addresses
              </Box>
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('emails', null)}>{stats.emails.total}</TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('emails', 'allowed')}>
              <Chip label={stats.emails.allowed} color={getStatusColor(stats.emails.allowed)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('emails', 'denied')}>
              <Chip label={stats.emails.denied} color={getStatusColor(stats.emails.denied)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('emails', 'whitelisted')}>
              <Chip label={stats.emails.whitelisted} color={getStatusColor(stats.emails.whitelisted)} size="small" />
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                {getFilterIcon('IP Addresses')}
                IP Addresses
              </Box>
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('ips', null)}>{stats.ips.total}</TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('ips', 'allowed')}>
              <Chip label={stats.ips.allowed} color={getStatusColor(stats.ips.allowed)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('ips', 'denied')}>
              <Chip label={stats.ips.denied} color={getStatusColor(stats.ips.denied)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('ips', 'whitelisted')}>
              <Chip label={stats.ips.whitelisted} color={getStatusColor(stats.ips.whitelisted)} size="small" />
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                {getFilterIcon('User Agents')}
                User Agents
              </Box>
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('userAgents', null)}>{stats.userAgents.total}</TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('userAgents', 'allowed')}>
              <Chip label={stats.userAgents.allowed} color={getStatusColor(stats.userAgents.allowed)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('userAgents', 'denied')}>
              <Chip label={stats.userAgents.denied} color={getStatusColor(stats.userAgents.denied)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('userAgents', 'whitelisted')}>
              <Chip label={stats.userAgents.whitelisted} color={getStatusColor(stats.userAgents.whitelisted)} size="small" />
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                {getFilterIcon('Username Rules')}
                Username Rules
              </Box>
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('usernames', null)}>{stats.usernames.total}</TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('usernames', 'allowed')}>
              <Chip label={stats.usernames.allowed} color={getStatusColor(stats.usernames.allowed)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('usernames', 'denied')}>
              <Chip label={stats.usernames.denied} color={getStatusColor(stats.usernames.denied)} size="small" />
            </TableCell>
            <TableCell sx={{ cursor: 'pointer' }} onClick={() => handleCountClick('usernames', 'whitelisted')}>
              <Chip label={stats.usernames.whitelisted} color={getStatusColor(stats.usernames.whitelisted)} size="small" />
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </TableContainer>
  </Paper>
);

const DashboardStats = () => {
  const [stats, setStats] = useState({
    totalIPs: 0,
    totalEmails: 0,
    totalUserAgents: 0,
    totalCountries: 0,
    totalCharsets: 0,
    totalUsernames: 0,
    dbConnections: { current: 0, max: 0, idle: 0, inUse: 0 }
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleCountClick = (type, status) => {
    let path = '/';
    if (type === 'ips') path = '/ip-list';
    if (type === 'emails') path = '/email-list';
    if (type === 'userAgents') path = '/useragent-list';
    if (type === 'countries') path = '/country-list';
    if (type === 'charsets') path = '/charset-list';
    if (type === 'usernames') path = '/username-list';
    const query = status ? `?status=${status}` : '';
    navigate(`${path}${query}`);
  };

  const fetchStats = useCallback(async () => {
    try {
      setLoading(true);
      setError('');

      // Lade alle Stats-Endpunkte
      const [ipsRes, emailsRes, userAgentsRes, countriesRes, charsetsRes, usernamesRes] = await Promise.all([
        axios.get('/api/ips/stats'),
        axios.get('/api/emails/stats'),
        axios.get('/api/user-agents/stats'),
        axios.get('/api/countries/stats'),
        axios.get('/api/charsets/stats'),
        axios.get('/api/usernames/stats')
      ]);

      setStats({
        ips: ipsRes.data,
        emails: emailsRes.data,
        userAgents: userAgentsRes.data,
        countries: countriesRes.data,
        charsets: charsetsRes.data,
        usernames: usernamesRes.data
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
      case 'Charset Rules':
        return <DnsIcon />;
      case 'Username Rules':
        return <PersonIcon />;
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

  // Layout: FilterStats auf volle Breite
  return (
    <Box sx={{ maxWidth: 1200, mx: 'auto', mt: 4 }}>
      <Box sx={{ flexGrow: 1 }}>
        <Box sx={{ width: '100%' }}>
          <FilterStats
            stats={stats}
            error={error}
            handleCountClick={handleCountClick}
            getFilterIcon={getFilterIcon}
            getStatusColor={getStatusColor}
          />
        </Box>
      </Box>
    </Box>
  );
};

export default DashboardStats; 