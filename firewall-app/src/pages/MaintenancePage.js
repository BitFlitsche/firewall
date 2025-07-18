import React, { useState } from 'react';
import axios from '../axiosConfig';
import {
  Box,
  Typography,
  Button,
  Alert,
  CircularProgress,
  Card,
  CardContent,
  CardActions,
  Grid,
  Divider,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
  Tooltip
} from '@mui/material';
import {
  Sync as SyncIcon,
  Warning as WarningIcon,
  CheckCircle as CheckCircleIcon,
  Info as InfoIcon,
  Build as BuildIcon,
  Storage as StorageIcon,
  Timer as SpeedIcon,
  HealthAndSafety as HealthAndSafetyIcon
} from '@mui/icons-material';

const MaintenancePage = () => {
  const [loading, setLoading] = useState({
    fullSync: false,
    incrementalSync: false,
    flushCache: false,
    recreateIndex: false
  });
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('info');
  const [showConfirmDialog, setShowConfirmDialog] = useState(false);
  const [systemStats, setSystemStats] = useState(null);
  const [syncStatus, setSyncStatus] = useState({ full_sync_running: false });

  const fetchSystemStats = async () => {
    try {
      const response = await axios.get('/api/system-stats');
      setSystemStats(response.data);
    } catch (error) {
      console.error('Failed to fetch system stats:', error);
    }
  };

  const fetchSyncStatus = async () => {
    try {
      const response = await axios.get('/sync/status');
      setSyncStatus(response.data);
    } catch (error) {
      console.error('Failed to fetch sync status:', error);
    }
  };

  React.useEffect(() => {
    fetchSystemStats();
    fetchSyncStatus();
    // Refresh stats every 30 seconds
    const interval = setInterval(() => {
      fetchSystemStats();
      fetchSyncStatus();
    }, 30000);
    return () => clearInterval(interval);
  }, []);

  const handleFullSync = async () => {
    setLoading(prev => ({ ...prev, fullSync: true }));
    setMessage('');
    
    try {
      const response = await axios.post('/sync/full');
      const recordCount = response.data?.records_synced || 0;
      setMessage(`Full sync completed successfully! ${recordCount > 0 ? `${recordCount} records synced.` : ''}`);
      setMessageType('success');
      
      // Refresh system stats after sync
      setTimeout(fetchSystemStats, 2000);
    } catch (error) {
      setMessage(`Full sync failed: ${error.response?.data?.error || error.message}`);
      setMessageType('error');
    } finally {
      setLoading(prev => ({ ...prev, fullSync: false }));
    }
  };

  const handleIncrementalSync = async () => {
    setLoading(prev => ({ ...prev, incrementalSync: true }));
    setMessage('');
    
    try {
      const response = await axios.post('/sync/force');
      const recordCount = response.data?.records_synced || 0;
      setMessage(`Incremental sync completed successfully! ${recordCount > 0 ? `${recordCount} records synced.` : ''}`);
      setMessageType('success');
      
      // Refresh system stats after sync
      setTimeout(fetchSystemStats, 2000);
    } catch (error) {
      setMessage(`Incremental sync failed: ${error.response?.data?.error || error.message}`);
      setMessageType('error');
    } finally {
      setLoading(prev => ({ ...prev, incrementalSync: false }));
    }
  };

  const handleRecreateIndex = async (endpoint, indexName) => {
    setLoading(prev => ({ ...prev, recreateIndex: true }));
    setMessage('');
    
    try {
      const response = await axios.post(endpoint);
      const recordCount = response.data?.records_indexed || 0;
      setMessage(`${indexName} index recreated successfully! ${recordCount > 0 ? `${recordCount} records indexed.` : ''}`);
      setMessageType('success');
      
      // Refresh system stats after sync
      setTimeout(fetchSystemStats, 2000);
    } catch (error) {
      setMessage(`${indexName} index recreation failed: ${error.response?.data?.error || error.message}`);
      setMessageType('error');
    } finally {
      setLoading(prev => ({ ...prev, recreateIndex: false }));
    }
  };

  const handleFlushCache = async () => {
    setLoading(prev => ({ ...prev, flushCache: true }));
    setMessage('');
    
    try {
      const response = await axios.post('/cache/flush');
      const itemsCleared = response.data?.items_cleared || 0;
      setMessage(`Cache flushed successfully! ${itemsCleared > 0 ? `${itemsCleared} items cleared.` : ''}`);
      setMessageType('success');
      
      // Refresh system stats after flush
      setTimeout(fetchSystemStats, 1000);
    } catch (error) {
      setMessage(`Cache flush failed: ${error.response?.data?.error || error.message}`);
      setMessageType('error');
    } finally {
      setLoading(prev => ({ ...prev, flushCache: false }));
    }
  };

  const handleConfirmFullSync = () => {
    setShowConfirmDialog(false);
    handleFullSync();
  };

  const getHealthColor = (status) => {
    switch (status) {
      case 'green': return 'success';
      case 'yellow': return 'warning';
      case 'red': return 'error';
      default: return 'info';
    }
  };

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatUptime = (seconds) => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    
    if (days > 0) return `${days}d ${hours}h ${minutes}m`;
    if (hours > 0) return `${hours}h ${minutes}m`;
    return `${minutes}m`;
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
        <BuildIcon />
        System Maintenance
      </Typography>
      
      <Typography variant="body1" sx={{ mb: 4, color: 'text.secondary' }}>
        Manage system operations and perform maintenance tasks. Use these tools carefully as they can impact system performance.
      </Typography>

      {message && (
        <Alert severity={messageType} sx={{ mb: 3 }}>
          {message}
        </Alert>
      )}

      <Grid container spacing={3}>

        {/* System Status Card */}
        <Grid item xs={12} md={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <HealthAndSafetyIcon />
                System Status
              </Typography>
              
              {systemStats && (
                <Box>
                  <Grid container spacing={2}>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Uptime</Typography>
                      <Typography variant="h6">{formatUptime(systemStats.uptime)}</Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">CPU Usage</Typography>
                      <Typography variant="h6">{systemStats.cpu_percent?.[0]?.toFixed(1) || 'N/A'}%</Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Memory Usage</Typography>
                      <Typography variant="h6">{systemStats.memory_percent?.toFixed(1) || 'N/A'}%</Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="text.secondary">Disk Usage</Typography>
                      <Typography variant="h6">{systemStats.disk_percent?.toFixed(1) || 'N/A'}%</Typography>
                    </Grid>
                  </Grid>
                  
                  <Divider sx={{ my: 2 }} />
                  
                  <Typography variant="subtitle2" gutterBottom>Service Health</Typography>
                  <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                    <Chip 
                      label={`DB: ${systemStats.db_health || 'unknown'}`}
                      color={systemStats.db_health === 'ok' ? 'success' : 'error'}
                      size="small"
                    />
                    <Chip 
                      label={`ES: ${systemStats.es_health || 'unknown'}`}
                      color={getHealthColor(systemStats.es_health)}
                      size="small"
                    />
                    <Chip 
                      label={`Requests: ${systemStats.request_count || 0}`}
                      color="info"
                      size="small"
                    />
                  </Box>
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Cache Stats Card */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <StorageIcon />
                Cache Statistics
              </Typography>
              
              {systemStats?.cache_stats && (
                <Grid container spacing={2}>
                  <Grid item xs={6} md={3}>
                    <Typography variant="body2" color="text.secondary">Cache Items</Typography>
                    <Typography variant="h6">{systemStats.cache_stats?.items || 0}</Typography>
                  </Grid>
                  <Grid item xs={6} md={3}>
                    <Typography variant="body2" color="text.secondary">Memory Usage</Typography>
                    <Typography variant="h6">
                      {systemStats.cache_stats?.memory_usage ? formatBytes(systemStats.cache_stats.memory_usage) : '0 Bytes'}
                    </Typography>
                  </Grid>
                  <Grid item xs={6} md={3}>
                    <Typography variant="body2" color="text.secondary">Hit Rate</Typography>
                    <Typography variant="h6">
                      {systemStats.cache_stats?.hit_rate ? `${(systemStats.cache_stats.hit_rate * 100).toFixed(1)}%` : 'N/A'}
                    </Typography>
                  </Grid>
                  <Grid item xs={6} md={3}>
                    <Typography variant="body2" color="text.secondary">Evictions</Typography>
                    <Typography variant="h6">{systemStats.cache_stats?.evictions || 0}</Typography>
                  </Grid>
                </Grid>
              )}
            </CardContent>
            <CardActions>
              <Button
                variant="outlined"
                color="warning"
                startIcon={loading.flushCache ? <CircularProgress size={20} /> : <StorageIcon />}
                onClick={handleFlushCache}
                disabled={loading.flushCache}
                fullWidth
              >
                {loading.flushCache ? 'Flushing...' : 'Flush Cache'}
              </Button>
            </CardActions>
          </Card>
        </Grid>

        {/* Index Management Card */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <BuildIcon />
                Index Management
              </Typography>
              
              <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Manage Elasticsearch indices and data synchronization. Recreate individual indices or perform full/incremental syncs.
              </Typography>
              
              {/* Sync Operations */}
              <Typography variant="subtitle2" gutterBottom sx={{ mt: 2, mb: 1 }}>
                Data Synchronization
                {syncStatus.full_sync_running && (
                  <Chip 
                    label="Full Sync Running" 
                    color="warning" 
                    size="small" 
                    sx={{ ml: 1 }}
                  />
                )}
              </Typography>
              <Grid container spacing={2} sx={{ mb: 3 }}>
                <Grid item xs={12} sm={6}>
                  <Tooltip 
                    title={
                      <Box>
                        <Typography variant="body2" sx={{ fontWeight: 'bold', mb: 1 }}>
                          Full Sync - Complete data synchronization
                        </Typography>
                        <Typography variant="body2" sx={{ mb: 1 }}>
                          Performs a complete synchronization of all data from MySQL to Elasticsearch.
                        </Typography>
                        <Typography variant="body2" sx={{ mb: 1 }}>
                          <strong>Use Cases:</strong> Initial setup, data recovery, schema changes, troubleshooting
                        </Typography>
                        <Typography variant="body2" sx={{ mb: 1 }}>
                          <strong>Impact:</strong> High resource usage, may temporarily slow down the system
                        </Typography>
                        <Typography variant="body2">
                          <strong>Safety:</strong> No data loss, updates all sync timestamps
                        </Typography>
                      </Box>
                    }
                    arrow
                    placement="top"
                  >
                    <Button
                      variant="contained"
                      color="primary"
                      startIcon={loading.fullSync ? <CircularProgress size={20} /> : <SyncIcon />}
                      onClick={() => setShowConfirmDialog(true)}
                      disabled={loading.fullSync || loading.incrementalSync || loading.recreateIndex || syncStatus.full_sync_running}
                      fullWidth
                    >
                      {loading.fullSync ? 'Syncing...' : 'Full Sync'}
                    </Button>
                  </Tooltip>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Tooltip 
                    title={
                      <Box>
                        <Typography variant="body2" sx={{ fontWeight: 'bold', mb: 1 }}>
                          Incremental Sync - Safe data synchronization
                        </Typography>
                        <Typography variant="body2" sx={{ mb: 1 }}>
                          Forces an immediate incremental sync of changed data from MySQL to Elasticsearch.
                        </Typography>
                        <Typography variant="body2" sx={{ mb: 1 }}>
                          <strong>Use Cases:</strong> Force sync after manual data changes, troubleshooting sync issues
                        </Typography>
                        <Typography variant="body2" sx={{ mb: 1 }}>
                          <strong>Impact:</strong> Low resource usage, only syncs changed records
                        </Typography>
                        <Typography variant="body2">
                          <strong>Safety:</strong> Safe operation, updates sync timestamps
                        </Typography>
                      </Box>
                    }
                    arrow
                    placement="top"
                  >
                    <Button
                      variant="outlined"
                      color="primary"
                      startIcon={loading.incrementalSync ? <CircularProgress size={20} /> : <SyncIcon />}
                      onClick={handleIncrementalSync}
                      disabled={loading.fullSync || loading.incrementalSync || loading.recreateIndex || syncStatus.full_sync_running}
                      fullWidth
                    >
                      {loading.incrementalSync ? 'Syncing...' : 'Incremental Sync'}
                    </Button>
                  </Tooltip>
                </Grid>
              </Grid>
              
              {/* Index Recreation */}
              <Typography variant="subtitle2" gutterBottom sx={{ mt: 2, mb: 1 }}>
                Index Recreation
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6} md={4}>
                  <Button
                    variant="outlined"
                    color="secondary"
                    fullWidth
                    onClick={() => handleRecreateIndex('/ip/recreate-index', 'IP Address')}
                    disabled={loading.recreateIndex || loading.fullSync || loading.incrementalSync || syncStatus.full_sync_running}
                  >
                    Recreate IP Index
                  </Button>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Button
                    variant="outlined"
                    color="secondary"
                    fullWidth
                    onClick={() => handleRecreateIndex('/emails/recreate-index', 'Email')}
                    disabled={loading.recreateIndex || loading.fullSync || loading.incrementalSync || syncStatus.full_sync_running}
                  >
                    Recreate Email Index
                  </Button>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Button
                    variant="outlined"
                    color="secondary"
                    fullWidth
                    onClick={() => handleRecreateIndex('/user-agents/recreate-index', 'User Agent')}
                    disabled={loading.recreateIndex || loading.fullSync || loading.incrementalSync || syncStatus.full_sync_running}
                  >
                    Recreate User Agent Index
                  </Button>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Button
                    variant="outlined"
                    color="secondary"
                    fullWidth
                    onClick={() => handleRecreateIndex('/countries/recreate-index', 'Country')}
                    disabled={loading.recreateIndex || loading.fullSync || loading.incrementalSync || syncStatus.full_sync_running}
                  >
                    Recreate Country Index
                  </Button>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Button
                    variant="outlined"
                    color="secondary"
                    fullWidth
                    onClick={() => handleRecreateIndex('/charsets/recreate-index', 'Charset')}
                    disabled={loading.recreateIndex || loading.fullSync || loading.incrementalSync || syncStatus.full_sync_running}
                  >
                    Recreate Charset Index
                  </Button>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Button
                    variant="outlined"
                    color="secondary"
                    fullWidth
                    onClick={() => handleRecreateIndex('/usernames/recreate-index', 'Username')}
                    disabled={loading.recreateIndex || loading.fullSync || loading.incrementalSync || syncStatus.full_sync_running}
                  >
                    Recreate Username Index
                  </Button>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Confirmation Dialog */}
      <Dialog
        open={showConfirmDialog}
        onClose={() => setShowConfirmDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <WarningIcon color="warning" />
          Confirm Full Sync
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to perform a full sync? This operation will:
          </DialogContentText>
          <List dense>
            <ListItem>
              <ListItemIcon>
                <WarningIcon fontSize="small" color="warning" />
              </ListItemIcon>
              <ListItemText primary="Sync all data from MySQL to Elasticsearch" />
            </ListItem>
            <ListItem>
              <ListItemIcon>
                <SpeedIcon fontSize="small" color="warning" />
              </ListItemIcon>
              <ListItemText primary="Use significant system resources" />
            </ListItem>
            <ListItem>
              <ListItemIcon>
                <CheckCircleIcon fontSize="small" color="success" />
              </ListItemIcon>
              <ListItemText primary="Update all sync timestamps" />
            </ListItem>
          </List>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowConfirmDialog(false)}>
            Cancel
          </Button>
          <Button 
            onClick={handleConfirmFullSync} 
            variant="contained" 
            color="primary"
            startIcon={<SyncIcon />}
          >
            Start Full Sync
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default MaintenancePage; 