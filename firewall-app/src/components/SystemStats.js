import React, { useEffect, useState } from 'react';
import { Box, Paper, Typography, Grid, LinearProgress, Chip, Tooltip } from '@mui/material';

const formatUptime = (seconds) => {
  const d = Math.floor(seconds / 86400);
  const h = Math.floor((seconds % 86400) / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);
  return `${d}d ${h}h ${m}m ${s}s`;
};

const SystemStats = () => {
  const [stats, setStats] = useState(null);
  const [error, setError] = useState('');

  useEffect(() => {
    let mounted = true;
    const fetchStats = async () => {
      try {
        const res = await fetch('/system-stats');
        const data = await res.json();
        if (mounted) setStats(data);
      } catch (err) {
        setError('Failed to load system stats');
      }
    };
    fetchStats();
    const interval = setInterval(fetchStats, 10000);
    return () => { mounted = false; clearInterval(interval); };
  }, []);

  if (error) return <Box sx={{ mb: 2 }}><Chip label={error} color="error" /></Box>;
  if (!stats) return <Box sx={{ mb: 2 }}><LinearProgress /></Box>;

  return (
    <Paper sx={{ p: 3, mb: 3 }} elevation={3}>
      <Typography variant="h5" gutterBottom>System Health</Typography>
      <Grid container spacing={2}>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="Uptime">
            <Chip label={`Uptime: ${formatUptime(stats.uptime)}`} color="primary" sx={{ width: '100%' }} />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="CPU Usage">
            <Box>
              <Typography variant="body2">CPU</Typography>
              <LinearProgress variant="determinate" value={stats.cpu_percent?.[0] || 0} sx={{ height: 10, borderRadius: 5 }} />
              <Typography variant="caption">{(stats.cpu_percent?.[0] || 0).toFixed(1)}%</Typography>
            </Box>
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="Memory Usage">
            <Box>
              <Typography variant="body2">Memory</Typography>
              <LinearProgress variant="determinate" value={stats.memory_percent || 0} sx={{ height: 10, borderRadius: 5 }} color="secondary" />
              <Typography variant="caption">{(stats.memory_used/1024/1024).toFixed(0)}MB / {(stats.memory_total/1024/1024).toFixed(0)}MB</Typography>
            </Box>
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="Disk Usage">
            <Box>
              <Typography variant="body2">Disk</Typography>
              <LinearProgress variant="determinate" value={stats.disk_percent || 0} sx={{ height: 10, borderRadius: 5 }} color="warning" />
              <Typography variant="caption">{(stats.disk_used/1024/1024/1024).toFixed(1)}GB / {(stats.disk_total/1024/1024/1024).toFixed(1)}GB</Typography>
            </Box>
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="Database Health">
            <Chip label={`DB: ${stats.db_health}`} color={stats.db_health === 'ok' ? 'success' : 'error'} sx={{ width: '100%' }} />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="DB Connections">
            <Chip label={`DB Conns: ${stats.db_connections}`} color="info" sx={{ width: '100%' }} />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="Elasticsearch Health">
            <Chip label={`ES: ${stats.es_health}`} color={stats.es_health === 'ok' ? 'success' : 'warning'} sx={{ width: '100%' }} />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="Go Routines">
            <Chip label={`GoRoutines: ${stats.go_routines}`} color="default" sx={{ width: '100%' }} />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="PID">
            <Chip label={`PID: ${stats.pid}`} color="default" sx={{ width: '100%' }} />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="Total Requests">
            <Chip label={`Requests: ${stats.request_count}`} color="primary" sx={{ width: '100%' }} />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Tooltip title="Total Errors">
            <Chip label={`Errors: ${stats.error_count}`} color={stats.error_count > 0 ? 'error' : 'success'} sx={{ width: '100%' }} />
          </Tooltip>
        </Grid>
      </Grid>
    </Paper>
  );
};

export default SystemStats; 