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


// Separate memoized search field component that never re-renders
// Separate filter controls component that only re-renders when filter values change
const FilterControls = memo(({ 
    searchValue, 
    onSearchChange, 
    filterStatus, 
    onStatusChange, 
    globalStatusCounts, 
    onReset 
}) => {
    return (
        <Box sx={{ display: 'flex', gap: 2, mb: 2, flexWrap: 'wrap', alignItems: 'center' }}>
            <TextField
                label="User Agent Filter"
                value={searchValue}
                onChange={onSearchChange}
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
    );
}, (prevProps, nextProps) => {
    // Only re-render if filter values change
    return prevProps.searchValue === nextProps.searchValue && 
           prevProps.filterStatus === nextProps.filterStatus &&
           prevProps.globalStatusCounts.total === nextProps.globalStatusCounts.total &&
           prevProps.globalStatusCounts.allowed === nextProps.globalStatusCounts.allowed &&
           prevProps.globalStatusCounts.denied === nextProps.globalStatusCounts.denied &&
           prevProps.globalStatusCounts.whitelisted === nextProps.globalStatusCounts.whitelisted;
});

// Separate form component that never re-renders
// Separate form component that only re-renders when form data changes
const UserAgentFormComponent = memo(({ onSubmit, userAgent, setUserAgent, status, setStatus, isRegex, setIsRegex, editId, setEditId }) => {
    return (
        <Box component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
            <TextField
                label="User Agent"
                value={userAgent}
                onChange={(e) => setUserAgent(e.target.value)}
                placeholder="Enter user agent or regex pattern"
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
                                • <code>.*bot.*</code> - Block all bot user agents
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>.*crawler.*</code> - Block web crawlers
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>.*spider.*</code> - Block search engine spiders
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>Mozilla/5\.0.*Chrome</code> - Block Chrome browsers
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>.*curl.*</code> - Block curl requests
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
                {editId ? 'Update User Agent' : 'Add User Agent'}
            </Button>
            {editId && (
                <Button variant="outlined" color="secondary" onClick={() => { setEditId(null); setUserAgent(''); setStatus('denied'); setIsRegex(false); }}>
                    Cancel Edit
                </Button>
            )}
        </Box>
    );
}, (prevProps, nextProps) => {
    // Only re-render when form data changes
    return prevProps.userAgent === nextProps.userAgent && 
           prevProps.status === nextProps.status && 
           prevProps.isRegex === nextProps.isRegex &&
           prevProps.editId === nextProps.editId;
});

// Separate table component that only re-renders when data changes
const UserAgentTable = memo(({ 
    userAgents, 
    loading, 
    total, 
    page, 
    rowsPerPage, 
    orderBy, 
    order,
    onSort,
    onPageChange,
    onRowsPerPageChange,
    onEdit,
    onDelete
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
                        onPageChange={onPageChange}
                        rowsPerPage={rowsPerPage}
                        onRowsPerPageChange={onRowsPerPageChange}
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
                                        active={orderBy === 'UserAgent'}
                                        direction={orderBy === 'UserAgent' ? order : 'asc'}
                                        onClick={() => onSort('UserAgent')}
                                    >
                                        User Agent
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
                                <TableCell>Is Regex</TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {userAgents.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">No User Agents</TableCell>
                                </TableRow>
                            ) : (
                                userAgents.map(userAgentItem => (
                                    <TableRow key={userAgentItem.ID}>
                                        <TableCell>{userAgentItem.ID}</TableCell>
                                        <TableCell>{userAgentItem.UserAgent}</TableCell>
                                        <TableCell>{userAgentItem.Status}</TableCell>
                                        <TableCell>{userAgentItem.IsRegex ? 'Yes' : 'No'}</TableCell>
                                        <TableCell>
                                            <IconButton onClick={() => onEdit(userAgentItem)} size="small"><EditIcon /></IconButton>
                                            <IconButton onClick={() => onDelete(userAgentItem.ID)} size="small" color="error"><DeleteIcon /></IconButton>
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
                        onPageChange={onPageChange}
                        rowsPerPage={rowsPerPage}
                        onRowsPerPageChange={onRowsPerPageChange}
                        rowsPerPageOptions={[10, 25, 50, 100]}
                        labelRowsPerPage="Einträge pro Seite:"
                    />
                </TableContainer>
            </Paper>
        </Box>
    );
});

const UserAgentForm = () => {
    const [userAgent, setUserAgent] = useState('');
    const [status, setStatus] = useState('denied');
    const [isRegex, setIsRegex] = useState(false);
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);
    const [userAgents, setUserAgents] = useState([]);
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
        axios.get('/api/user-agents/stats')
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

    // Debounce search input with focus preservation
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearchValue(searchValue);
        }, 300);

        return () => clearTimeout(timer);
    }, [searchValue]);





    // Fetch user agents with server-side filtering, sorting, and pagination
    useEffect(() => {
        const fetchUserAgents = async () => {
            setLoading(true);
            try {
                const response = await axios.get('/api/user-agents', {
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
                    setUserAgents(response.data.items);
                    setTotal(response.data.total || response.data.items.length);
                } else {
                    setUserAgents([]);
                    setTotal(0);
                }
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch user agents');
                setLoading(false);
            }
        };
        fetchUserAgents();
    }, [refresh, page, rowsPerPage, filterStatus, debouncedSearchValue, orderBy, order]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            if (editId) {
                await axios.put(`/api/user-agent/${editId}`, { UserAgent: userAgent, Status: status, IsRegex: isRegex });
                setMessage('User Agent updated successfully');
            } else {
                await axios.post('/api/user-agent', { UserAgent: userAgent, Status: status, IsRegex: isRegex });
                setMessage('User Agent added successfully');
            }
            setUserAgent('');
            setStatus('denied');
            setIsRegex(false); // Reset regex checkbox
            setEditId(null);
            setRefresh(r => !r);
        } catch (error) {
            setError('Error saving user agent');
        }
    };





    const handleStatusChange = (e) => {
        setFilterStatus(e.target.value);
        setPage(0); // Reset to first page when filtering
    };

    const handleReset = () => {
        setSearchValue('');
        setFilterStatus('');
        setPage(0);
    };

    const handleSearchChange = useCallback((e) => {
        setSearchValue(e.target.value);
        setPage(0); // Reset to first page when searching
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

    const handleEdit = useCallback((userAgentItem) => {
        setUserAgent(userAgentItem.UserAgent);
        setStatus(userAgentItem.Status);
        setIsRegex(userAgentItem.IsRegex || false); // Set regex checkbox
        setEditId(userAgentItem.ID);
    }, []);

    const handleDelete = useCallback(async (id) => {
        if (!window.confirm('Delete this User Agent?')) return;
        try {
            await axios.delete(`/api/user-agent/${id}`);
            setMessage('User Agent deleted');
            setRefresh(r => !r);
        } catch {
            setError('Error deleting User Agent');
        }
    }, []);



    return (
        <Box sx={{ maxWidth: 700, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>User Agent Management</Typography>
                <UserAgentFormComponent 
                    onSubmit={handleSubmit}
                    userAgent={userAgent}
                    setUserAgent={setUserAgent}
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
                

                
                <UserAgentTable 
                    userAgents={userAgents}
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

export default UserAgentForm;
