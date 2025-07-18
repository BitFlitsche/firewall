import React, { useState, useEffect, useCallback, memo } from 'react';
import axios from '../axiosConfig';
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
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Tooltip from '@mui/material/Tooltip';
import InfoIcon from '@mui/icons-material/Info';
import { useLocation } from 'react-router-dom';
import RecreateIndexButton from './RecreateIndexButton';

// Memoized Form Component
const UsernameFormComponent = memo(({ onSubmit, username, setUsername, status, setStatus, isRegex, setIsRegex, editId, setEditId }) => {
    return (
        <Box component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
            <TextField
                label="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Enter username or regex pattern"
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
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <FormControlLabel
                    control={
                        <Checkbox
                            checked={isRegex}
                            onChange={(e) => setIsRegex(e.target.checked)}
                            color="primary"
                        />
                    }
                    label="Use as Regular Expression"
                />
                <Tooltip 
                    title={
                        <Box>
                            <Typography variant="body2" sx={{ fontWeight: 'bold', mb: 1 }}>
                                Regex Examples:
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>admin.*</code> - Block admin-like usernames
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>.*bot.*</code> - Block bot usernames
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>test.*</code> - Block test usernames
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>^guest[0-9]+$</code> - Block guest123, guest456, etc.
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>.*admin.*</code> - Block any username containing "admin"
                            </Typography>
                        </Box>
                    }
                    arrow
                    placement="top"
                >
                    <InfoIcon color="action" sx={{ fontSize: 20 }} />
                </Tooltip>
            </Box>
            <Button type="submit" variant="contained" color="primary">
                {editId ? 'Update Username' : 'Add Username'}
            </Button>
            {editId && (
                <Button variant="outlined" color="secondary" onClick={() => { setEditId(null); setUsername(''); setStatus('denied'); setIsRegex(false); }}>
                    Cancel Edit
                </Button>
            )}
        </Box>
    );
}, (prevProps, nextProps) => {
    // Only re-render when form data changes
    return prevProps.username === nextProps.username && 
           prevProps.status === nextProps.status && 
           prevProps.isRegex === nextProps.isRegex &&
           prevProps.editId === nextProps.editId;
});

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
            label="Username Filter"
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
        <RecreateIndexButton indexType="username" />
    </Box>
));

// Memoized Table Component
const UsernameTable = React.memo(({ 
    usernames, 
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
                                        active={orderBy === 'Username'}
                                        direction={orderBy === 'Username' ? order : 'asc'}
                                        onClick={() => onSort('Username')}
                                    >
                                        Username
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
                                <TableCell>Regex</TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {usernames.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">No Usernames</TableCell>
                                </TableRow>
                            ) : (
                                usernames.map(usernameItem => (
                                    <TableRow key={usernameItem.ID}>
                                        <TableCell>{usernameItem.ID}</TableCell>
                                        <TableCell>{usernameItem.Username}</TableCell>
                                        <TableCell>{usernameItem.Status}</TableCell>
                                        <TableCell>{usernameItem.IsRegex ? 'Yes' : 'No'}</TableCell>
                                        <TableCell>
                                            <IconButton onClick={() => onEdit(usernameItem)} color="primary">
                                                <EditIcon />
                                            </IconButton>
                                            <IconButton onClick={() => onDelete(usernameItem.ID)} color="error">
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

const UsernameForm = () => {
    const [username, setUsername] = useState('');
    const [status, setStatus] = useState('denied');
    const [isRegex, setIsRegex] = useState(false);
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);
    const [usernames, setUsernames] = useState([]);
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
        axios.get('/usernames/stats')
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

    // Fetch usernames with server-side filtering, sorting, and pagination
    useEffect(() => {
        const fetchUsernames = async () => {
            setLoading(true);
            try {
                const response = await axios.get('/usernames', {
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
                    setUsernames(response.data.items);
                    setTotal(response.data.total || response.data.items.length);
                } else {
                    setUsernames([]);
                    setTotal(0);
                }
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch usernames');
                setLoading(false);
            }
        };
        fetchUsernames();
    }, [refresh, page, rowsPerPage, filterStatus, debouncedSearchValue, orderBy, order]);

    const handleSubmit = useCallback(async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            if (editId) {
                await axios.put(`/username/${editId}`, { username, status, IsRegex: isRegex });
                setMessage('Username updated successfully');
            } else {
                await axios.post('/username', { username, status, IsRegex: isRegex });
                setMessage('Username added successfully');
            }
            setUsername('');
            setStatus('denied');
            setIsRegex(false);
            setEditId(null);
            setRefresh(r => !r);
        } catch (err) {
            setError('Error saving username');
        }
    }, [username, status, editId, isRegex]);

    const handleDelete = useCallback(async (id) => {
        if (!window.confirm('Delete this username?')) return;
        try {
            await axios.delete(`/username/${id}`);
            setMessage('Username deleted');
            setRefresh(r => !r);
        } catch {
            setError('Error deleting username');
        }
    }, []);

    const handleEdit = useCallback((usernameItem) => {
        setUsername(usernameItem.Username);
        setStatus(usernameItem.Status);
        setIsRegex(usernameItem.IsRegex || false);
        setEditId(usernameItem.ID);
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
        setUsername('');
        setStatus('denied');
        setIsRegex(false);
        setEditId(null);
    }, []);

    const handleUsernameChange = useCallback((e) => {
        setUsername(e.target.value);
    }, []);

    const handleStatusChangeForm = useCallback((e) => {
        setStatus(e.target.value);
    }, []);

    // Memoized values for components
    const formProps = React.useMemo(() => ({
        username,
        status,
        message,
        error,
        editId,
        onUsernameChange: handleUsernameChange,
        onStatusChange: handleStatusChangeForm,
        onSubmit: handleSubmit,
        onCancelEdit: handleCancelEdit,
        isRegex,
        setIsRegex
    }), [username, status, message, error, editId, handleUsernameChange, handleStatusChangeForm, handleSubmit, handleCancelEdit, isRegex]);

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
        usernames,
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
    }), [usernames, loading, orderBy, order, page, rowsPerPage, total, handleSort, handleEdit, handleDelete, handleChangePage, handleChangeRowsPerPage]);

    return (
        <Box sx={{ maxWidth: 700, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>Username Management</Typography>
                <UsernameFormComponent 
                    onSubmit={handleSubmit}
                    username={username}
                    setUsername={setUsername}
                    status={status}
                    setStatus={setStatus}
                    isRegex={isRegex}
                    setIsRegex={setIsRegex}
                    editId={editId}
                    setEditId={setEditId}
                />
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
                
                {/* Filter controls outside of table component */}
                <Box sx={{ mt: 4, mb: 2 }}>
                    <FilterControls 
                        searchValue={searchValue}
                        onSearchChange={handleSearchChange}
                        filterStatus={filterStatus}
                        onStatusChange={handleStatusChange}
                        globalStatusCounts={globalStatusCounts}
                        onReset={handleReset}
                    />
                </Box>
                
                <UsernameTable 
                    usernames={usernames}
                    loading={loading}
                    total={total}
                    page={page}
                    rowsPerPage={rowsPerPage}
                    orderBy={orderBy}
                    order={order}
                    onSort={handleSort}
                    onPageChange={handleChangePage}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                />
            </Paper>
        </Box>
    );
};

export default UsernameForm; 