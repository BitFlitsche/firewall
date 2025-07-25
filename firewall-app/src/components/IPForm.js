import React, { useState, useEffect, useCallback } from 'react';
import axiosInstance from '../axiosConfig';
import ConflictTable from './ConflictTable';
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
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import TableSortLabel from '@mui/material/TableSortLabel';
import TablePagination from '@mui/material/TablePagination';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import Tooltip from '@mui/material/Tooltip';
import { useLocation } from 'react-router-dom';
import InfiniteScrollIPTable from './InfiniteScrollIPTable';


// Memoized Form Component
const IPFormComponent = React.memo(({ 
    ip, 
    status, 
    isCidr,
    source,
    message, 
    error, 
    editId, 
    onIpChange, 
    onStatusChange, 
    onCidrChange,
    onSourceChange,
    onSubmit, 
    onCancelEdit 
}) => (
    <Box component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
        <TextField
            label="IP Address or CIDR Block"
            value={ip}
            onChange={onIpChange}
            placeholder={isCidr ? "Enter CIDR (e.g., 192.168.1.0/24)" : "Enter IP address"}
            required
            fullWidth
            helperText={isCidr ? "CIDR format: IP/prefix (e.g., 10.0.0.0/8, 172.16.0.0/12)" : ""}
        />
        <TextField
            select
            label="Status"
            value={status}
            onChange={onStatusChange}
            fullWidth
        >
            <MenuItem value="denied">Denied</MenuItem>
            <MenuItem value="allowed">Allowed</MenuItem>
            <MenuItem value="whitelisted">Whitelisted</MenuItem>
        </TextField>
        <FormControlLabel
            control={
                <Checkbox
                    checked={isCidr}
                    onChange={onCidrChange}
                    color="primary"
                />
            }
            label={
                <Tooltip title="Check this to add a CIDR block (IP range) instead of a single IP address">
                    <span>CIDR Block</span>
                </Tooltip>
            }
        />
        <TextField
            select
            label="Source"
            value={source}
            onChange={onSourceChange}
            fullWidth
        >
            <MenuItem value="manual">Manual</MenuItem>
            <MenuItem value="stopforumspam_toxic_cidr">StopForumSpam Toxic CIDR</MenuItem>
        </TextField>
        <Button type="submit" variant="contained" color="primary">
            {editId ? 'Update IP' : 'Add IP'}
        </Button>
        {editId && (
            <Button variant="outlined" color="secondary" onClick={onCancelEdit}>
                Cancel Edit
            </Button>
        )}
    </Box>
));

// Memoized Filter Controls Component
const FilterControls = React.memo(({ 
    searchValue, 
    filterStatus, 
    filterType,
    globalStatusCounts,
    globalTypeCounts,
    onSearchChange, 
    onStatusChange, 
    onTypeChange,
    onReset,
    onSearchFocus,
    onSearchBlur,
    searchInputRef
}) => (
    <Box sx={{ display: 'flex', gap: 2, mb: 2, flexWrap: 'wrap', alignItems: 'center' }}>
        <TextField
            label="IP Address Filter"
            value={searchValue}
            onChange={onSearchChange}
            onFocus={onSearchFocus}
            onBlur={onSearchBlur}
            ref={searchInputRef}
            size="small"
        />
        <FormControl size="small" sx={{ minWidth: 140 }}>
            <InputLabel shrink>Status</InputLabel>
            <Select
                value={filterStatus}
                label="Status"
                onChange={onStatusChange}
                displayEmpty
                renderValue={(selected) => {
                    if (!selected) return `All (${globalStatusCounts.total})`;
                    if (selected === 'allowed') return `Allowed (${globalStatusCounts.allowed})`;
                    if (selected === 'denied') return `Denied (${globalStatusCounts.denied})`;
                    if (selected === 'whitelisted') return `Whitelisted (${globalStatusCounts.whitelisted})`;
                    return selected;
                }}
            >
                <MenuItem key="" value="">
                    All ({globalStatusCounts.total})
                </MenuItem>
                <MenuItem key="allowed" value="allowed">
                    Allowed ({globalStatusCounts.allowed})
                </MenuItem>
                <MenuItem key="denied" value="denied">
                    Denied ({globalStatusCounts.denied})
                </MenuItem>
                <MenuItem key="whitelisted" value="whitelisted">
                    Whitelisted ({globalStatusCounts.whitelisted})
                </MenuItem>
            </Select>
        </FormControl>
        <FormControl size="small" sx={{ minWidth: 140 }}>
            <InputLabel shrink>Type</InputLabel>
            <Select
                value={filterType}
                label="Type"
                onChange={onTypeChange}
                displayEmpty
                renderValue={(selected) => {
                    if (!selected) return `All (${globalTypeCounts.total})`;
                    if (selected === 'single') return `Single IP (${globalTypeCounts.single})`;
                    if (selected === 'cidr') return `CIDR Block (${globalTypeCounts.cidr})`;
                    return selected;
                }}
            >
                <MenuItem key="" value="">
                    All ({globalTypeCounts.total})
                </MenuItem>
                <MenuItem key="single" value="single">
                    Single IP ({globalTypeCounts.single})
                </MenuItem>
                <MenuItem key="cidr" value="cidr">
                    CIDR Block ({globalTypeCounts.cidr})
                </MenuItem>
            </Select>
        </FormControl>
        <Button variant="outlined" size="small" onClick={onReset}>
            Reset
        </Button>
    </Box>
));

// Memoized Table Component
const IPTable = React.memo(({ 
    ips, 
    loading, 
    orderBy, 
    order, 
    page, 
    rowsPerPage, 
    total,
    onSort, 
    onEdit, 
    onDelete, 
    onChangePage, 
    onChangeRowsPerPage 
}) => {
    if (loading) return <div>Loading...</div>;

    return (
        <Box className="list-section" sx={{ mt: 4 }}>
            <Paper elevation={2} sx={{ p: 2 }}>
                {/* Table with Pagination */}
                <TableContainer>
                    <TablePagination
                        component="div"
                        count={total}
                        page={page}
                        onPageChange={onChangePage}
                        rowsPerPage={rowsPerPage}
                        onRowsPerPageChange={onChangeRowsPerPage}
                        rowsPerPageOptions={[10, 25, 50, 100]}
                        labelRowsPerPage="Entries per page:"
                    />
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'id'}
                                        direction={orderBy === 'id' ? order : 'asc'}
                                        onClick={() => onSort('id')}
                                    >
                                        ID
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'address'}
                                        direction={orderBy === 'address' ? order : 'asc'}
                                        onClick={() => onSort('address')}
                                    >
                                        IP Address / CIDR
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'status'}
                                        direction={orderBy === 'status' ? order : 'asc'}
                                        onClick={() => onSort('status')}
                                    >
                                        Status
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>Type</TableCell>
                                <TableCell>Source</TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {ips.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={6} align="center">No IP addresses</TableCell>
                                </TableRow>
                            ) : (
                                ips.map((ipItem) => (
                                    <TableRow key={ipItem.id}>
                                        <TableCell>{ipItem.id}</TableCell>
                                        <TableCell>{ipItem.address}</TableCell>
                                        <TableCell>{ipItem.status}</TableCell>
                                        <TableCell>{ipItem.is_cidr ? 'CIDR Block' : 'Single IP'}</TableCell>
                                        <TableCell>{ipItem.source || 'manual'}</TableCell>
                                        <TableCell>
                                            <IconButton onClick={() => onEdit(ipItem)} size="small">
                                                <EditIcon />
                                            </IconButton>
                                            <IconButton onClick={() => onDelete(ipItem.id)} size="small" color="error">
                                                <DeleteIcon />
                                            </IconButton>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                    <TablePagination
                        component="div"
                        count={total}
                        page={page}
                        onPageChange={onChangePage}
                        rowsPerPage={rowsPerPage}
                        onRowsPerPageChange={onChangeRowsPerPage}
                        rowsPerPageOptions={[10, 25, 50, 100]}
                        labelRowsPerPage="Entries per page:"
                    />
                </TableContainer>
            </Paper>
        </Box>
    );
});

const IPForm = () => {
    const [ip, setIp] = useState('');
    const [status, setStatus] = useState('denied');
    const [isCidr, setIsCidr] = useState(false);
    const [source, setSource] = useState('manual');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [conflicts, setConflicts] = useState([]);
    const [isDeletingConflicts, setIsDeletingConflicts] = useState(false);
    const [pendingOperation, setPendingOperation] = useState(null);
    const [refresh, setRefresh] = useState(false);
    const [ips, setIps] = useState([]);
    const [editId, setEditId] = useState(null);
    
    // Filtering and infinite scroll state
    const [loading, setLoading] = useState(true);
    const [infiniteLoading, setInfiniteLoading] = useState(false);
    const [filterStatus, setFilterStatus] = useState('');
    const [filterType, setFilterType] = useState('');
    const [orderBy, setOrderBy] = useState('id');
    const [order, setOrder] = useState('desc');
    const [total, setTotal] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [globalStatusCounts, setGlobalStatusCounts] = useState({ allowed: 0, denied: 0, whitelisted: 0, total: 0 });
    const [globalTypeCounts, setGlobalTypeCounts] = useState({ single: 0, cidr: 0, total: 0 });
    const location = useLocation();

    // Debounced search state
    const [searchValue, setSearchValue] = useState('');
    const [debouncedSearchValue, setDebouncedSearchValue] = useState('');
    const searchInputRef = React.useRef(null);
    const [wasFocused, setWasFocused] = useState(false);

    // Set initial filterStatus from query param
    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const status = params.get('status');
        if (status && ['allowed','denied','whitelisted'].includes(status)) {
            setFilterStatus(status);
        }
    }, [location.search]);

    // Load global status and type counts
    useEffect(() => {
        axiosInstance.get('/api/ips/stats')
            .then(res => {
                setGlobalStatusCounts({
                    allowed: res.data.allowed || 0,
                    denied: res.data.denied || 0,
                    whitelisted: res.data.whitelisted || 0,
                    total: res.data.total || 0,
                });
                setGlobalTypeCounts({
                    single: res.data.single || 0,
                    cidr: res.data.cidr || 0,
                    total: res.data.total || 0,
                });
            })
            .catch(() => {
                setGlobalStatusCounts({ allowed: 0, denied: 0, whitelisted: 0, total: 0 });
                setGlobalTypeCounts({ single: 0, cidr: 0, total: 0 });
            });
    }, [refresh]);

    // Debounce search input
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearchValue(searchValue);
        }, 300);

        return () => clearTimeout(timer);
    }, [searchValue]);

    // Restore focus after any re-render if the field was previously focused
    useEffect(() => {
        if (wasFocused && searchInputRef.current && searchInputRef.current.value) {
            searchInputRef.current.focus();
            // Restore cursor position to end of input
            const length = searchInputRef.current.value.length;
            searchInputRef.current.setSelectionRange(length, length);
        }
    });

    // Fetch IPs with server-side filtering, sorting, and infinite scroll
    useEffect(() => {
        const fetchIPs = async () => {
            setLoading(true);
            
            try {
                const params = {
                    page: 1,
                    limit: 25,
                    status: filterStatus || undefined,
                    type: filterType || undefined,
                    search: debouncedSearchValue || undefined,
                    orderBy,
                    order,
                };
                
                const response = await axiosInstance.get('/api/ips', { params });
                
                if (response.data && response.data.items) {
                    setIps(response.data.items);
                    setHasMore(response.data.items.length === 25);
                    setTotal(response.data.total || 0);
                } else {
                    setIps([]);
                    setHasMore(false);
                    setTotal(0);
                }
            } catch (err) {
                setError('Failed to fetch IP addresses');
            } finally {
                setLoading(false);
            }
        };
        
        // Reset infinite scroll state when filters change
        setHasMore(true);
        setIps([]);
        
        fetchIPs();
    }, [refresh, filterStatus, filterType, debouncedSearchValue, orderBy, order]);

    const handleSubmit = useCallback(async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        setConflicts([]);
        setPendingOperation(null);
        
        const operation = {
            address: ip,
            status: status,
            is_cidr: isCidr,
            editId: editId
        };
        
        try {
            if (editId) {
                await axiosInstance.put(`/api/ip/${editId}`, {
                    address: ip,
                    status: status,
                    is_cidr: isCidr,
                    source: 'manual'
                });
                setMessage('IP address updated successfully');
            } else {
                await axiosInstance.post('/api/ip', {
                    address: ip,
                    status: status,
                    is_cidr: isCidr,
                    source: 'manual'
                });
                setMessage('IP address added successfully');
            }
            setIp('');
            setStatus('denied');
            setIsCidr(false);
            setEditId(null);
            setRefresh(r => !r);
        } catch (error) {
            if (error.response?.status === 409 && error.response?.data?.conflicts) {
                // Handle conflict response
                setConflicts(error.response.data.conflicts);
                setPendingOperation(operation);
                const hasErrors = error.response.data.conflicts.some(c => c.severity === 'error');
                if (hasErrors) {
                    setError('IP address creation blocked due to conflicts');
                } else {
                    setMessage('IP address added with warnings');
                    setIp('');
                    setStatus('denied');
                    setIsCidr(false);
                    setEditId(null);
                    setRefresh(r => !r);
                }
            } else {
                setError(error.response?.data?.error || 'Error saving IP address');
            }
        }
    }, [ip, status, isCidr, editId]);

    const handleDeleteAllConflicts = useCallback(async () => {
        if (!pendingOperation || conflicts.length === 0) return;
        
        setIsDeletingConflicts(true);
        setError('');
        
        try {
            // Extract unique conflicting addresses from ERROR conflicts only
            const errorConflicts = conflicts.filter(c => c.severity === 'error');
            const conflictingAddresses = [...new Set(errorConflicts.map(c => c.conflicting[0]))];
            
            if (conflictingAddresses.length === 0) {
                setError('No error-level conflicts to delete');
                return;
            }
            
            // First, fetch all IPs to get their IDs
            const response = await axiosInstance.get('/api/ips', {
                params: {
                    limit: 1000, // Get all IPs to find the conflicting ones
                    page: 1
                }
            });
            
            const allIPs = response.data.items || [];
            
            // Find IDs of conflicting addresses
            const conflictingIds = [];
            for (const address of conflictingAddresses) {
                const foundIP = allIPs.find(ip => ip.address === address);
                if (foundIP) {
                    conflictingIds.push(foundIP.id);
                }
            }
            
            // Delete all conflicting entries by ID
            const deletePromises = conflictingIds.map(id => 
                axiosInstance.delete(`/api/ip/${id}`)
            );
            
            await Promise.all(deletePromises);
            
            // Retry the original operation
            const { address, status: opStatus, is_cidr, editId: opEditId } = pendingOperation;
            
            if (opEditId) {
                await axiosInstance.put(`/api/ip/${opEditId}`, {
                    address,
                    status: opStatus,
                    is_cidr
                });
                setMessage(`IP address updated successfully after removing ${conflictingIds.length} error conflicts`);
            } else {
                await axiosInstance.post('/api/ip', {
                    address,
                    status: opStatus,
                    is_cidr
                });
                setMessage(`IP address added successfully after removing ${conflictingIds.length} error conflicts`);
            }
            
            // Clear form and refresh
            setIp('');
            setStatus('denied');
            setIsCidr(false);
            setEditId(null);
            setConflicts([]);
            setPendingOperation(null);
            setRefresh(r => !r);
            
        } catch (error) {
            setError('Failed to delete conflicts or retry operation: ' + (error.response?.data?.error || error.message));
        } finally {
            setIsDeletingConflicts(false);
        }
    }, [conflicts, pendingOperation]);

    const handleDelete = useCallback(async (id) => {
        if (!window.confirm('Delete this IP address?')) return;
        try {
            await axiosInstance.delete(`/api/ip/${id}`);
            setMessage('IP address deleted');
            setRefresh(r => !r);
        } catch {
            setError('Error deleting IP address');
        }
    }, []);

    const handleEdit = useCallback((ipItem) => {
        setIp(ipItem.address);
        setStatus(ipItem.status);
        setIsCidr(ipItem.is_cidr || false);
        setSource('manual'); // Always set to manual when editing
        setEditId(ipItem.id);
    }, []);

    const handleSort = useCallback((field) => {
        if (orderBy === field) {
            setOrder(order === 'asc' ? 'desc' : 'asc');
        } else {
            setOrderBy(field);
            setOrder('asc');
        }
        // Reset IPs when sorting changes
        setIps([]);
        setHasMore(true);
    }, [orderBy, order]);

    const handleSearchChange = useCallback((e) => {
        setSearchValue(e.target.value);
        setIps([]); // Reset IPs when searching
        setHasMore(true);
    }, []);

    const handleCidrChange = useCallback((e) => {
        setIsCidr(e.target.checked);
    }, []);

    const handleStatusChange = useCallback((e) => {
        setFilterStatus(e.target.value);
        setIps([]); // Reset IPs when filtering
        setHasMore(true);
    }, []);

    const handleTypeChange = useCallback((e) => {
        setFilterType(e.target.value);
        setIps([]); // Reset IPs when filtering
        setHasMore(true);
    }, []);

    const handleReset = useCallback(() => {
        setSearchValue('');
        setFilterStatus('');
        setFilterType('');
        setIps([]);
        setHasMore(true);
    }, []);

    const handleLoadMore = useCallback(() => {
        if (!infiniteLoading && hasMore) {
            const fetchIPs = async () => {
                setInfiniteLoading(true);
                try {
                    const params = {
                        page: Math.floor(ips.length / 25) + 1,
                        limit: 25,
                        status: filterStatus || undefined,
                        type: filterType || undefined,
                        search: debouncedSearchValue || undefined,
                        orderBy,
                        order,
                    };
                    
                    const response = await axiosInstance.get('/api/ips', { params });
                    
                    if (response.data && response.data.items) {
                        setIps(prev => [...prev, ...response.data.items]);
                        setHasMore(response.data.items.length === 25);
                    } else {
                        setHasMore(false);
                    }
                } catch (err) {
                    setError('Failed to fetch more IP addresses');
                } finally {
                    setInfiniteLoading(false);
                }
            };
            fetchIPs();
        }
    }, [infiniteLoading, hasMore, ips.length, filterStatus, filterType, debouncedSearchValue, orderBy, order]);

    const handleSearchFocus = useCallback(() => {
        setWasFocused(true);
    }, []);

    const handleSearchBlur = useCallback(() => {
        setWasFocused(false);
    }, []);

    const handleCancelEdit = useCallback(() => {
        setIp('');
        setStatus('denied');
        setIsCidr(false);
        setSource('manual');
        setEditId(null);
    }, []);

    const handleStopForumSpamImport = useCallback(async () => {
        if (window.confirm('This will import toxic IP addresses in CIDR format from StopForumSpam. Existing StopForumSpam records will be replaced. Continue?')) {
            try {
                setLoading(true);
                await axiosInstance.post('/api/ips/import-stopforumspam');
                setMessage('StopForumSpam toxic CIDR data imported successfully!');
                // Wait a bit for the import to complete, then refresh
                setTimeout(() => {
                    setRefresh(prev => !prev);
                }, 2000);
            } catch (err) {
                const errorMessage = err.response?.data?.error || err.message || 'Failed to import StopForumSpam data';
                setError(errorMessage);
            } finally {
                setLoading(false);
            }
        }
    }, []);

    const handleStopForumSpamStats = useCallback(async () => {
        try {
            const response = await axiosInstance.get('/api/ips/stopforumspam-stats');
            const stats = response.data;
            const statsMessage = `StopForumSpam Toxic CIDRs: ${stats.total_stopforumspam_cidrs || 0}\nLast Import: ${stats.last_import ? new Date(stats.last_import).toLocaleString() : 'Never'}`;
            alert(statsMessage);
        } catch (err) {
            setError('Failed to get StopForumSpam stats');
        }
    }, []);

    const handleStopForumSpamStatus = useCallback(async () => {
        try {
            const response = await axiosInstance.get('/api/ips/stopforumspam-status');
            const status = response.data;
            const importStatus = status.import_enabled ? 'Enabled' : 'Disabled';
            const statusMessage = `Import Enabled: ${importStatus}\nImport Running: ${status.is_running ? 'Yes' : 'No'}\nTotal Imported: ${status.total_imported}\nLast Import: ${status.last_import}`;
            alert(statusMessage);
        } catch (err) {
            setError('Failed to get StopForumSpam import status');
        }
    }, []);

    const handleIpChange = useCallback((e) => {
        setIp(e.target.value);
    }, []);

    const handleStatusChangeForm = useCallback((e) => {
        setStatus(e.target.value);
    }, []);

    const handleSourceChange = useCallback((e) => {
        setSource(e.target.value);
    }, []);

    // Memoized values for components
    const formProps = React.useMemo(() => ({
        ip,
        status,
        isCidr,
        source,
        message,
        error,
        editId,
        onIpChange: handleIpChange,
        onStatusChange: handleStatusChangeForm,
        onCidrChange: handleCidrChange,
        onSourceChange: handleSourceChange,
        onSubmit: handleSubmit,
        onCancelEdit: handleCancelEdit
    }), [ip, status, isCidr, source, message, error, editId, handleIpChange, handleStatusChangeForm, handleCidrChange, handleSourceChange, handleSubmit, handleCancelEdit]);

    const filterProps = React.useMemo(() => ({
        searchValue,
        filterStatus,
        filterType,
        globalStatusCounts,
        globalTypeCounts,
        onSearchChange: handleSearchChange,
        onStatusChange: handleStatusChange,
        onTypeChange: handleTypeChange,
        onReset: handleReset,
        onSearchFocus: handleSearchFocus,
        onSearchBlur: handleSearchBlur,
        searchInputRef
    }), [searchValue, filterStatus, filterType, globalStatusCounts, globalTypeCounts, handleSearchChange, handleStatusChange, handleTypeChange, handleReset, handleSearchFocus, handleSearchBlur]);

    const tableProps = React.useMemo(() => ({
        ips,
        loading: loading || infiniteLoading,
        error,
        total,
        hasMore,
        onLoadMore: handleLoadMore,
        onSort: handleSort,
        onEdit: handleEdit,
        onDelete: handleDelete,
        orderBy,
        order
    }), [ips, loading, infiniteLoading, error, total, hasMore, handleLoadMore, handleSort, handleEdit, handleDelete, orderBy, order]);

    return (
        <Box sx={{ maxWidth: 1200, mx: 'auto', mt: 4, p: 2 }}>
            <Typography variant="h4" gutterBottom>
                IP Address Management
            </Typography>

            {/* StopForumSpam Import Buttons */}
            <Box sx={{ mb: 3 }}>
                <Button 
                    variant="contained" 
                    color="primary" 
                    onClick={handleStopForumSpamImport}
                    sx={{ mr: 2 }}
                >
                    Import stopforumspam Toxic IP Addresses (CIDR)
                </Button>
                <Button 
                    variant="outlined" 
                    onClick={handleStopForumSpamStats}
                    sx={{ mr: 2 }}
                >
                    View StopForumSpam Stats
                </Button>
                <Button 
                    variant="outlined" 
                    onClick={handleStopForumSpamStatus}
                >
                    Check Import Status
                </Button>
            </Box>

            {/* Form Section */}
            <Paper sx={{ p: 3, mb: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>
                    Add or Modify IP/CIDR
                </Typography>
                
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

                <IPFormComponent {...formProps} />
                <ConflictTable 
                    conflicts={conflicts} 
                    open={conflicts.length > 0} 
                    onDeleteAllConflicts={handleDeleteAllConflicts}
                    isDeleting={isDeletingConflicts}
                />
            </Paper>

            {/* Filter Controls */}
            <FilterControls
                searchValue={searchValue}
                filterStatus={filterStatus}
                filterType={filterType}
                globalStatusCounts={globalStatusCounts}
                globalTypeCounts={globalTypeCounts}
                onSearchChange={handleSearchChange}
                onStatusChange={handleStatusChange}
                onTypeChange={handleTypeChange}
                onReset={handleReset}
                onSearchFocus={handleSearchFocus}
                onSearchBlur={handleSearchBlur}
                searchInputRef={searchInputRef}
            />

            {/* Table */}
            <InfiniteScrollIPTable
                ips={ips}
                loading={loading || infiniteLoading}
                error={error}
                total={total}
                hasMore={hasMore}
                onLoadMore={handleLoadMore}
                onSort={handleSort}
                orderBy={orderBy}
                order={order}
                onEdit={handleEdit}
                onDelete={handleDelete}
            />
        </Box>
    );
};

export default IPForm;
