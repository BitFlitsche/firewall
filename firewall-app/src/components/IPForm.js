import React, { useState, useEffect, useCallback } from 'react';
import axiosInstance from '../axiosConfig';
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
import { useLocation } from 'react-router-dom';

// Memoized Form Component
const IPFormComponent = React.memo(({ 
    ip, 
    status, 
    message, 
    error, 
    editId, 
    onIpChange, 
    onStatusChange, 
    onSubmit, 
    onCancelEdit 
}) => (
    <Box component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
        <TextField
            label="IP Address"
            value={ip}
            onChange={onIpChange}
            placeholder="Enter IP address"
            required
            fullWidth
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
    globalStatusCounts,
    onSearchChange, 
    onStatusChange, 
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
                        labelRowsPerPage="Einträge pro Seite:"
                    />
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'ID'}
                                        direction={orderBy === 'ID' ? order : 'asc'}
                                        onClick={() => onSort('ID')}
                                    >
                                        ID
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'Address'}
                                        direction={orderBy === 'Address' ? order : 'asc'}
                                        onClick={() => onSort('Address')}
                                    >
                                        IP Address
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'Status'}
                                        direction={orderBy === 'Status' ? order : 'asc'}
                                        onClick={() => onSort('Status')}
                                    >
                                        Status
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {ips.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={4} align="center">No IP addresses</TableCell>
                                </TableRow>
                            ) : (
                                ips.map((ipItem) => (
                                    <TableRow key={ipItem.ID}>
                                        <TableCell>{ipItem.ID}</TableCell>
                                        <TableCell>{ipItem.Address}</TableCell>
                                        <TableCell>{ipItem.Status}</TableCell>
                                        <TableCell>
                                            <IconButton onClick={() => onEdit(ipItem)} size="small">
                                                <EditIcon />
                                            </IconButton>
                                            <IconButton onClick={() => onDelete(ipItem.ID)} size="small" color="error">
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
                        labelRowsPerPage="Einträge pro Seite:"
                    />
                </TableContainer>
            </Paper>
        </Box>
    );
});

const IPForm = () => {
    const [ip, setIp] = useState('');
    const [status, setStatus] = useState('denied');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);
    const [ips, setIps] = useState([]);
    const [editId, setEditId] = useState(null);
    
    // Filtering and pagination state
    const [loading, setLoading] = useState(true);
    const [filterStatus, setFilterStatus] = useState('');
    const [orderBy, setOrderBy] = useState('ID');
    const [order, setOrder] = useState('desc');
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [total, setTotal] = useState(0);
    const [globalStatusCounts, setGlobalStatusCounts] = useState({ allowed: 0, denied: 0, whitelisted: 0, total: 0 });
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

    // Load global status counts
    useEffect(() => {
        axiosInstance.get('/ips/stats')
            .then(res => {
                setGlobalStatusCounts({
                    allowed: res.data.allowed || 0,
                    denied: res.data.denied || 0,
                    whitelisted: res.data.whitelisted || 0,
                    total: res.data.total || 0,
                });
            })
            .catch(() => setGlobalStatusCounts({ allowed: 0, denied: 0, whitelisted: 0, total: 0 }));
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

    // Fetch IPs with server-side filtering, sorting, and pagination
    useEffect(() => {
        const fetchIPs = async () => {
            setLoading(true);
            try {
                const response = await axiosInstance.get('/ips', {
                    params: {
                        page: page + 1,
                        limit: rowsPerPage,
                        status: filterStatus || undefined,
                        search: debouncedSearchValue || undefined,
                        orderBy,
                        order,
                    }
                });
                if (response.data && response.data.items) {
                    setIps(response.data.items);
                    setTotal(response.data.total || response.data.items.length);
                } else {
                    setIps([]);
                    setTotal(0);
                }
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch IP addresses');
                setLoading(false);
            }
        };
        fetchIPs();
    }, [refresh, page, rowsPerPage, filterStatus, debouncedSearchValue, orderBy, order]);

    const handleSubmit = useCallback(async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            if (editId) {
                await axiosInstance.put(`/ip/${editId}`, {
                    Address: ip,
                    Status: status
                });
                setMessage('IP address updated successfully');
            } else {
                await axiosInstance.post('/ip', {
                    Address: ip,
                    Status: status
                });
                setMessage('IP address added successfully');
            }
            setIp('');
            setStatus('denied');
            setEditId(null);
            setRefresh(r => !r);
        } catch (error) {
            setError(error.response?.data?.error || 'Error saving IP address');
        }
    }, [ip, status, editId]);

    const handleDelete = useCallback(async (id) => {
        if (!window.confirm('Delete this IP address?')) return;
        try {
            await axiosInstance.delete(`/ip/${id}`);
            setMessage('IP address deleted');
            setRefresh(r => !r);
        } catch {
            setError('Error deleting IP address');
        }
    }, []);

    const handleEdit = useCallback((ipItem) => {
        setIp(ipItem.Address);
        setStatus(ipItem.Status);
        setEditId(ipItem.ID);
    }, []);

    const handleSort = useCallback((field) => {
        if (orderBy === field) {
            setOrder(order === 'asc' ? 'desc' : 'asc');
        } else {
            setOrderBy(field);
            setOrder('asc');
        }
    }, [orderBy, order]);

    const handleChangePage = useCallback((event, newPage) => {
        setPage(newPage);
    }, []);

    const handleChangeRowsPerPage = useCallback((event) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    }, []);

    const handleSearchChange = useCallback((e) => {
        setSearchValue(e.target.value);
        setPage(0); // Reset to first page when searching
    }, []);

    const handleStatusChange = useCallback((e) => {
        setFilterStatus(e.target.value);
        setPage(0); // Reset to first page when filtering
    }, []);

    const handleReset = useCallback(() => {
        setSearchValue('');
        setFilterStatus('');
        setPage(0);
    }, []);

    const handleSearchFocus = useCallback(() => {
        setWasFocused(true);
    }, []);

    const handleSearchBlur = useCallback(() => {
        setWasFocused(false);
    }, []);

    const handleCancelEdit = useCallback(() => {
        setIp('');
        setStatus('denied');
        setEditId(null);
    }, []);

    const handleIpChange = useCallback((e) => {
        setIp(e.target.value);
    }, []);

    const handleStatusChangeForm = useCallback((e) => {
        setStatus(e.target.value);
    }, []);

    // Memoized values for components
    const formProps = React.useMemo(() => ({
        ip,
        status,
        message,
        error,
        editId,
        onIpChange: handleIpChange,
        onStatusChange: handleStatusChangeForm,
        onSubmit: handleSubmit,
        onCancelEdit: handleCancelEdit
    }), [ip, status, message, error, editId, handleIpChange, handleStatusChangeForm, handleSubmit, handleCancelEdit]);

    const filterProps = React.useMemo(() => ({
        searchValue,
        filterStatus,
        globalStatusCounts,
        onSearchChange: handleSearchChange,
        onStatusChange: handleStatusChange,
        onReset: handleReset,
        onSearchFocus: handleSearchFocus,
        onSearchBlur: handleSearchBlur,
        searchInputRef
    }), [searchValue, filterStatus, globalStatusCounts, handleSearchChange, handleStatusChange, handleReset, handleSearchFocus, handleSearchBlur]);

    const tableProps = React.useMemo(() => ({
        ips,
        loading,
        orderBy,
        order,
        page,
        rowsPerPage,
        total,
        onSort: handleSort,
        onEdit: handleEdit,
        onDelete: handleDelete,
        onChangePage: handleChangePage,
        onChangeRowsPerPage: handleChangeRowsPerPage
    }), [ips, loading, orderBy, order, page, rowsPerPage, total, handleSort, handleEdit, handleDelete, handleChangePage, handleChangeRowsPerPage]);

    return (
        <Box sx={{ maxWidth: 700, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>IP Address Management</Typography>
                <IPFormComponent {...formProps} />
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
                
                {/* Filter controls outside of table component */}
                <Box sx={{ mt: 4, mb: 2 }}>
                    <FilterControls {...filterProps} />
                </Box>
                
                <IPTable {...tableProps} />
            </Paper>
        </Box>
    );
};

export default IPForm;
