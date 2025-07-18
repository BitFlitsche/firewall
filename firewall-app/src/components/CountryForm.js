import React, { useState, useEffect, useCallback } from 'react';
import axios from '../axiosConfig';
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Alert from '@mui/material/Alert';
import MenuItem from '@mui/material/MenuItem';
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import TableSortLabel from '@mui/material/TableSortLabel';
import TablePagination from '@mui/material/TablePagination';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Select from '@mui/material/Select';
import { useLocation } from 'react-router-dom';
import RecreateIndexButton from './RecreateIndexButton';

// Memoized Form Component
const CountryFormComponent = React.memo(({ 
    country, 
    status, 
    message, 
    error, 
    editId, 
    onCountryChange, 
    onStatusChange, 
    onSubmit, 
    onCancelEdit 
}) => (
    <Box component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
        <TextField
            label="Country Code"
            value={country}
            onChange={onCountryChange}
            placeholder="Enter Country Code (e.g. DE)"
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
            {editId ? 'Update Country' : 'Add Country'}
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
            label="Country Code Filter"
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
        <RecreateIndexButton indexType="country" />
    </Box>
));

// Memoized Table Component
const CountryTable = React.memo(({ 
    countries, 
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
                                        active={orderBy === 'Code'}
                                        direction={orderBy === 'Code' ? order : 'asc'}
                                        onClick={() => onSort('Code')}
                                    >
                                        Country Code
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
                            {countries.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={4} align="center">No countries</TableCell>
                                </TableRow>
                            ) : (
                                countries.map((countryItem) => (
                                    <TableRow key={countryItem.ID}>
                                        <TableCell>{countryItem.ID}</TableCell>
                                        <TableCell>{countryItem.Code}</TableCell>
                                        <TableCell>{countryItem.Status}</TableCell>
                                        <TableCell>
                                            <IconButton onClick={() => onEdit(countryItem)} size="small">
                                                <EditIcon />
                                            </IconButton>
                                            <IconButton onClick={() => onDelete(countryItem.ID)} size="small" color="error">
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

const CountryForm = () => {
    const [country, setCountry] = useState('');
    const [status, setStatus] = useState('denied');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);
    const [countries, setCountries] = useState([]);
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
        axios.get('/countries/stats')
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

    // Fetch countries with server-side filtering, sorting, and pagination
    useEffect(() => {
        const fetchCountries = async () => {
            setLoading(true);
            try {
                const response = await axios.get('/countries', {
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
                    setCountries(response.data.items);
                    setTotal(response.data.total || response.data.items.length);
                } else {
                    setCountries([]);
                    setTotal(0);
                }
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch countries');
                setLoading(false);
            }
        };
        fetchCountries();
    }, [refresh, page, rowsPerPage, filterStatus, debouncedSearchValue, orderBy, order]);

    const handleSubmit = useCallback(async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            if (editId) {
                await axios.put(`/country/${editId}`, { code: country, status });
                setMessage('Country updated successfully');
            } else {
                await axios.post('/country', { code: country, status });
                setMessage('Country added successfully');
            }
            setCountry('');
            setStatus('denied');
            setEditId(null);
            setRefresh(r => !r);
        } catch (err) {
            setError('Error saving country');
        }
    }, [country, status, editId]);

    const handleDelete = useCallback(async (id) => {
        if (!window.confirm('Delete this country?')) return;
        try {
            await axios.delete(`/country/${id}`);
            setMessage('Country deleted');
            setRefresh(r => !r);
        } catch {
            setError('Error deleting country');
        }
    }, []);

    const handleEdit = useCallback((countryItem) => {
        setCountry(countryItem.Code);
        setStatus(countryItem.Status);
        setEditId(countryItem.ID);
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
        setCountry('');
        setStatus('denied');
        setEditId(null);
    }, []);

    const handleCountryChange = useCallback((e) => {
        setCountry(e.target.value);
    }, []);

    const handleStatusChangeForm = useCallback((e) => {
        setStatus(e.target.value);
    }, []);

    // Memoized values for components
    const formProps = React.useMemo(() => ({
        country,
        status,
        message,
        error,
        editId,
        onCountryChange: handleCountryChange,
        onStatusChange: handleStatusChangeForm,
        onSubmit: handleSubmit,
        onCancelEdit: handleCancelEdit
    }), [country, status, message, error, editId, handleCountryChange, handleStatusChangeForm, handleSubmit, handleCancelEdit]);

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
        countries,
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
    }), [countries, loading, orderBy, order, page, rowsPerPage, total, handleSort, handleEdit, handleDelete, handleChangePage, handleChangeRowsPerPage]);

    return (
        <Box sx={{ maxWidth: 700, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>Country Management</Typography>
                <CountryFormComponent {...formProps} />
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
                
                {/* Filter controls outside of table component */}
                <Box sx={{ mt: 4, mb: 2 }}>
                    <FilterControls {...filterProps} />
                </Box>
                
                <CountryTable {...tableProps} />
            </Paper>
        </Box>
    );
};

export default CountryForm;
