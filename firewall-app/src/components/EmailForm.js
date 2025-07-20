import React, { useState, useEffect, useCallback, memo } from 'react';
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
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Tooltip from '@mui/material/Tooltip';
import InfoIcon from '@mui/icons-material/Info';
import { useLocation } from 'react-router-dom';


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
                label="Email Filter"
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

// Separate form component that only re-renders when form data changes
const EmailFormComponent = memo(({ onSubmit, email, setEmail, status, setStatus, isRegex, setIsRegex, editId, setEditId }) => {
    return (
        <Box component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
            <TextField
                label="Email Address"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Enter email address or regex pattern"
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
                                • <code>.*@spam\.com$</code> - Block all emails from spam.com
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>.*@.*\.ru$</code> - Block all .ru domain emails
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>^admin@.*</code> - Block admin emails from any domain
                            </Typography>
                            <Typography variant="body2" component="div">
                                • <code>.*@gmail\.com$</code> - Block all Gmail addresses
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
                {editId ? 'Update Email' : 'Add Email'}
            </Button>
            {editId && (
                <Button variant="outlined" color="secondary" onClick={() => { setEditId(null); setEmail(''); setStatus('denied'); setIsRegex(false); }}>
                    Cancel Edit
                </Button>
            )}
        </Box>
    );
}, (prevProps, nextProps) => {
    // Only re-render when form data changes
    return prevProps.email === nextProps.email && 
           prevProps.status === nextProps.status && 
           prevProps.isRegex === nextProps.isRegex &&
           prevProps.editId === nextProps.editId;
});

// Separate table component that only re-renders when data changes
const EmailTable = memo(({ 
    emails, 
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
                                        active={orderBy === 'address'}
                                        direction={orderBy === 'address' ? order : 'asc'}
                                        onClick={() => onSort('address')}
                                    >
                                        Email Address
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
                                <TableCell>Regex</TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {emails.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">No Emails</TableCell>
                                </TableRow>
                            ) : (
                                emails.map(emailItem => (
                                    <TableRow key={emailItem.id}>
                                        <TableCell>{emailItem.id}</TableCell>
                                        <TableCell>{emailItem.address}</TableCell>
                                        <TableCell>{emailItem.status}</TableCell>
                                        <TableCell>{emailItem.is_regex ? 'Yes' : 'No'}</TableCell>
                                        <TableCell>
                                            <IconButton onClick={() => onEdit(emailItem)} color="primary">
                                                <EditIcon />
                                            </IconButton>
                                            <IconButton onClick={() => onDelete(emailItem.id)} color="error">
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

const EmailForm = () => {
    const [email, setEmail] = useState('');
    const [status, setStatus] = useState('denied');
    const [isRegex, setIsRegex] = useState(false);
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);
    const [emails, setEmails] = useState([]);
    const [editId, setEditId] = useState(null);
    
    // Filtering and pagination state
    const [loading, setLoading] = useState(true);
    const [filterStatus, setFilterStatus] = useState('');
    const [orderBy, setOrderBy] = useState('id');
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
        axiosInstance.get('/api/emails/stats')
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



    // Fetch emails with server-side filtering, sorting, and pagination
    useEffect(() => {
        const fetchEmails = async () => {
            setLoading(true);
            try {
                const response = await axiosInstance.get('/api/emails', {
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
                    setEmails(response.data.items);
                    setTotal(response.data.total || response.data.items.length);
                } else {
                    setEmails([]);
                    setTotal(0);
                }
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch emails');
                setLoading(false);
            }
        };
        fetchEmails();
    }, [refresh, page, rowsPerPage, filterStatus, debouncedSearchValue, orderBy, order]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            if (editId) {
                await axiosInstance.put(`/api/email/${editId}`, { address: email, status, IsRegex: isRegex });
                setMessage('Email updated successfully');
            } else {
                await axiosInstance.post('/api/email', { address: email, status, IsRegex: isRegex });
                setMessage('Email added successfully');
            }
            setEmail('');
            setStatus('denied');
            setIsRegex(false);
            setEditId(null);
            setRefresh(r => !r);
        } catch (error) {
            setError('Error saving email');
        }
    };



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

    const handleEdit = useCallback((emailItem) => {
        setEmail(emailItem.address);
        setStatus(emailItem.status);
        setEditId(emailItem.id);
    }, []);

    const handleDelete = useCallback(async (id) => {
        if (!window.confirm('Delete this email?')) return;
        try {
            await axiosInstance.delete(`/api/email/${id}`);
            setMessage('Email deleted');
            setRefresh(r => !r);
        } catch {
            setError('Error deleting email');
        }
    }, []);



    return (
        <Box sx={{ maxWidth: 1200, mx: 'auto', mt: 4, p: 2 }}>
            <Typography variant="h4" gutterBottom>
                Email Management
            </Typography>

            {/* Form Section */}
            <Paper sx={{ p: 3, mb: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>
                    Add or Modify Email
                </Typography>
                
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

                <EmailFormComponent 
                    onSubmit={handleSubmit}
                    email={email}
                    setEmail={setEmail}
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
            <EmailTable 
                emails={emails}
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
        </Box>
    );
};

export default EmailForm;
