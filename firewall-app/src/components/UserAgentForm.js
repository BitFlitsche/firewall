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
import InfiniteScrollUserAgentTable from './InfiniteScrollUserAgentTable';


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
                                        active={orderBy === 'user_agent'}
                                        direction={orderBy === 'user_agent' ? order : 'asc'}
                                        onClick={() => onSort('user_agent')}
                                    >
                                        User Agent
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
                                    <TableRow key={userAgentItem.id}>
                                        <TableCell>{userAgentItem.id}</TableCell>
                                        <TableCell>{userAgentItem.user_agent}</TableCell>
                                        <TableCell>{userAgentItem.status}</TableCell>
                                        <TableCell>{userAgentItem.is_regex ? 'Yes' : 'No'}</TableCell>
                                        <TableCell>
                                            <IconButton onClick={() => onEdit(userAgentItem)} size="small"><EditIcon /></IconButton>
                                            <IconButton onClick={() => onDelete(userAgentItem.id)} size="small" color="error"><DeleteIcon /></IconButton>
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
                        labelRowsPerPage="Entries per page:"
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
    
    // Filtering and infinite scroll state
    const [loading, setLoading] = useState(true);
    const [infiniteLoading, setInfiniteLoading] = useState(false);
    const [filterStatus, setFilterStatus] = useState('');
    const [orderBy, setOrderBy] = useState('id');
    const [order, setOrder] = useState('desc');
    const [total, setTotal] = useState(0);
    const [hasMore, setHasMore] = useState(true);
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





    // Fetch user agents with server-side filtering, sorting, and infinite scroll
    useEffect(() => {
        const fetchUserAgents = async () => {
            setLoading(true);
            
            try {
                const params = {
                    page: 1,
                    limit: 25,
                    status: filterStatus || undefined,
                    search: debouncedSearchValue || undefined,
                    orderBy,
                    order,
                };
                
                const response = await axios.get('/api/user-agents', { params });
                
                if (response.data && response.data.items) {
                    setUserAgents(response.data.items);
                    setHasMore(response.data.items.length === 25);
                    setTotal(response.data.total || 0);
                } else {
                    setUserAgents([]);
                    setHasMore(false);
                    setTotal(0);
                }
            } catch (err) {
                setError('Failed to fetch user agents');
            } finally {
                setLoading(false);
            }
        };
        
        // Reset infinite scroll state when filters change
        setHasMore(true);
        setUserAgents([]);
        
        fetchUserAgents();
    }, [refresh, filterStatus, debouncedSearchValue, orderBy, order]);

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
        setUserAgents([]); // Reset user agents when filtering
        setHasMore(true);
    };

    const handleReset = () => {
        setSearchValue('');
        setFilterStatus('');
        setUserAgents([]);
        setHasMore(true);
    };

    const handleSearchChange = useCallback((e) => {
        setSearchValue(e.target.value);
        setUserAgents([]); // Reset user agents when searching
        setHasMore(true);
    }, []);

    const handleSort = useCallback((field) => {
        if (orderBy === field) {
            setOrder(order === 'asc' ? 'desc' : 'asc');
        } else {
            setOrderBy(field);
            setOrder('asc');
        }
        // Reset user agents when sorting changes
        setUserAgents([]);
        setHasMore(true);
    }, [orderBy, order]);

    const handleLoadMore = useCallback(() => {
        if (!infiniteLoading && hasMore) {
            const fetchUserAgents = async () => {
                setInfiniteLoading(true);
                try {
                    const params = {
                        page: Math.floor(userAgents.length / 25) + 1,
                        limit: 25,
                        status: filterStatus || undefined,
                        search: debouncedSearchValue || undefined,
                        orderBy,
                        order,
                    };
                    
                    const response = await axios.get('/api/user-agents', { params });
                    
                    if (response.data && response.data.items) {
                        setUserAgents(prev => [...prev, ...response.data.items]);
                        setHasMore(response.data.items.length === 25);
                    } else {
                        setHasMore(false);
                    }
                } catch (err) {
                    setError('Failed to fetch more user agents');
                } finally {
                    setInfiniteLoading(false);
                }
            };
            fetchUserAgents();
        }
    }, [infiniteLoading, hasMore, userAgents.length, filterStatus, debouncedSearchValue, orderBy, order]);

    const handleEdit = useCallback((userAgentItem) => {
        setUserAgent(userAgentItem.user_agent);
        setStatus(userAgentItem.status);
        setEditId(userAgentItem.id);
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
        <Box sx={{ maxWidth: 1200, mx: 'auto', mt: 4, p: 2 }}>
            <Typography variant="h4" gutterBottom>
                User Agent Management
            </Typography>

            {/* Form Section */}
            <Paper sx={{ p: 3, mb: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>
                    Add or Modify User Agent
                </Typography>
                
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

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
            </Paper>

            {/* Filter Controls */}
            <FilterControls 
                searchValue={searchValue}
                onSearchChange={handleSearchChange}
                filterStatus={filterStatus}
                onStatusChange={handleStatusChange}
                globalStatusCounts={globalStatusCounts}
                onReset={handleReset}
            />

            {/* Table */}
            <InfiniteScrollUserAgentTable 
                userAgents={userAgents}
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

export default UserAgentForm;
