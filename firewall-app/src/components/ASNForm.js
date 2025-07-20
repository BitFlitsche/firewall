import React, { useState, useEffect, useCallback } from 'react';
import axiosInstance from '../axiosConfig';
import { formatCountryDisplay } from '../utils/country_codes';
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
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import TableSortLabel from '@mui/material/TableSortLabel';
import TablePagination from '@mui/material/TablePagination';
import Link from '@mui/material/Link';
import { useLocation } from 'react-router-dom';

// Memoized Form Component
const ASNFormComponent = React.memo(({ 
    asn, 
    rir,
    domain,
    country,
    name,
    status, 
    source,
    message, 
    error, 
    editId, 
    onASNChange, 
    onRIRChange,
    onDomainChange,
    onCountryChange,
    onNameChange,
    onStatusChange, 
    onSourceChange,
    onSubmit, 
    onCancelEdit 
}) => (
    <Box component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
        <TextField
            label="ASN Number"
            value={asn}
            onChange={onASNChange}
            placeholder="Enter ASN (e.g., AS12345)"
            required
            fullWidth
            helperText="ASN format: AS followed by numbers (e.g., AS12345, AS7922)"
        />
        <TextField
            label="RIR (Regional Internet Registry) - Optional"
            value={rir}
            onChange={onRIRChange}
            placeholder="Enter RIR (e.g., arin, ripencc)"
            fullWidth
            helperText="Regional Internet Registry (e.g., arin, ripencc, apnic) - Optional"
        />
        <TextField
            label="Domain - Optional"
            value={domain}
            onChange={onDomainChange}
            placeholder="Enter domain name"
            fullWidth
            helperText="Domain name associated with the ASN - Optional"
        />
        <TextField
            label="Country Code - Optional"
            value={country}
            onChange={onCountryChange}
            placeholder="Enter country code (e.g., US, GB)"
            fullWidth
            helperText="ISO 3166-1 alpha-2 country code (e.g., US, GB, DE) - Optional"
        />
        <TextField
            label="ASN Name/Description"
            value={name}
            onChange={onNameChange}
            placeholder="Enter ASN name or description"
            required
            fullWidth
            helperText="Description of the ASN (e.g., Comcast Cable Communications)"
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
        <TextField
            label="Source"
            value={source}
            onChange={onSourceChange}
            placeholder="Enter source (e.g., manual, spamhaus)"
            fullWidth
            helperText="Source of the ASN data (e.g., manual, spamhaus)"
        />
        <Button type="submit" variant="contained" color="primary">
            {editId ? 'Update ASN' : 'Add ASN'}
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
    filterRIR,
    filterCountry,
    globalStatusCounts,
    globalRIRCounts,
    globalCountryCounts,
    onSearchChange, 
    onStatusChange,
    onRIRChange,
    onCountryChange,
    onReset,
    onSearchFocus,
    onSearchBlur,
    searchInputRef
}) => (
    <Box sx={{ display: 'flex', gap: 2, mb: 2, flexWrap: 'wrap', alignItems: 'center' }}>
        <TextField
            label="Search ASN, Domain, Name"
            value={searchValue}
            onChange={onSearchChange}
            onFocus={onSearchFocus}
            onBlur={onSearchBlur}
            ref={searchInputRef}
            size="small"
            placeholder="Search ASN, domain, or name..."
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
            <InputLabel shrink>RIR</InputLabel>
            <Select
                value={filterRIR}
                label="RIR"
                onChange={onRIRChange}
                displayEmpty
                renderValue={(selected) => {
                    if (!selected) return `All (${globalRIRCounts.total})`;
                    return `${selected} (${globalRIRCounts[selected] || 0})`;
                }}
            >
                <MenuItem key="" value="">
                    All ({globalRIRCounts.total})
                </MenuItem>
                {Object.entries(globalRIRCounts).filter(([key]) => key !== 'total').map(([rir, count]) => (
                    <MenuItem key={rir} value={rir}>
                        {rir} ({count})
                    </MenuItem>
                ))}
            </Select>
        </FormControl>
        <FormControl size="small" sx={{ minWidth: 200 }}>
            <InputLabel shrink>Country</InputLabel>
            <Select
                value={filterCountry}
                label="Country"
                onChange={onCountryChange}
                displayEmpty
                renderValue={(selected) => {
                    if (!selected) return `All (${globalCountryCounts.total})`;
                    const count = globalCountryCounts[selected] || 0;
                    return formatCountryDisplay(selected, count);
                }}
            >
                <MenuItem key="" value="">
                    All ({globalCountryCounts.total})
                </MenuItem>
                {Object.entries(globalCountryCounts).filter(([key]) => key !== 'total').map(([country, count]) => (
                    <MenuItem key={country} value={country}>
                        {formatCountryDisplay(country, count)}
                    </MenuItem>
                ))}
            </Select>
        </FormControl>
        <Button variant="outlined" size="small" onClick={onReset}>
            Reset
        </Button>
    </Box>
));

// Memoized Table Component
const ASNTable = React.memo(({ 
    asns, 
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
                                        active={orderBy === 'asn'}
                                        direction={orderBy === 'asn' ? order : 'asc'}
                                        onClick={() => onSort('asn')}
                                    >
                                        ASN
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'rir'}
                                        direction={orderBy === 'rir' ? order : 'asc'}
                                        onClick={() => onSort('rir')}
                                    >
                                        RIR
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'domain'}
                                        direction={orderBy === 'domain' ? order : 'asc'}
                                        onClick={() => onSort('domain')}
                                    >
                                        Domain
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'cc'}
                                        direction={orderBy === 'cc' ? order : 'asc'}
                                        onClick={() => onSort('cc')}
                                    >
                                        Country
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'asname'}
                                        direction={orderBy === 'asname' ? order : 'asc'}
                                        onClick={() => onSort('asname')}
                                    >
                                        Name
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
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'source'}
                                        direction={orderBy === 'source' ? order : 'asc'}
                                        onClick={() => onSort('source')}
                                    >
                                        Source
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {asns.map((asn) => (
                                <TableRow key={asn.id}>
                                    <TableCell>{asn.id}</TableCell>
                                    <TableCell>{asn.asn}</TableCell>
                                    <TableCell>{asn.rir}</TableCell>
                                    <TableCell>{asn.domain}</TableCell>
                                    <TableCell>{asn.cc}</TableCell>
                                    <TableCell>{asn.asname}</TableCell>
                                    <TableCell>
                                        <span className={`status-${asn.status}`}>
                                            {asn.status}
                                        </span>
                                    </TableCell>
                                    <TableCell>{asn.source}</TableCell>
                                    <TableCell>
                                        <IconButton
                                            size="small"
                                            onClick={() => onEdit(asn)}
                                            color="primary"
                                        >
                                            <EditIcon />
                                        </IconButton>
                                        <IconButton
                                            size="small"
                                            onClick={() => onDelete(asn.id)}
                                            color="error"
                                        >
                                            <DeleteIcon />
                                        </IconButton>
                                    </TableCell>
                                </TableRow>
                            ))}
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

const ASNForm = () => {
    const [asn, setASN] = useState('');
    const [rir, setRIR] = useState('');
    const [domain, setDomain] = useState('');
    const [country, setCountry] = useState('');
    const [name, setName] = useState('');
    const [status, setStatus] = useState('denied');
    const [source, setSource] = useState('manual');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);
    const [asns, setASNs] = useState([]);
    const [editId, setEditId] = useState(null);
    
    // Filtering and pagination state
    const [loading, setLoading] = useState(true);
    const [filterStatus, setFilterStatus] = useState('');
    const [filterRIR, setFilterRIR] = useState('');
    const [filterCountry, setFilterCountry] = useState('');
    const [orderBy, setOrderBy] = useState('id');
    const [order, setOrder] = useState('desc');
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [total, setTotal] = useState(0);
    const [globalStatusCounts, setGlobalStatusCounts] = useState({ allowed: 0, denied: 0, whitelisted: 0, total: 0 });
    const [globalRIRCounts, setGlobalRIRCounts] = useState({ total: 0 });
    const [globalCountryCounts, setGlobalCountryCounts] = useState({ total: 0 });
    const location = useLocation();

    // Debounced search state
    const [searchValue, setSearchValue] = useState('');
    const [debouncedSearchValue, setDebouncedSearchValue] = useState('');
    const searchInputRef = React.useRef(null);

    // UI state
    const [showResources, setShowResources] = useState(false);

    // Set initial filterStatus from query param
    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const status = params.get('status');
        if (status && ['allowed','denied','whitelisted'].includes(status)) {
            setFilterStatus(status);
        }
    }, [location.search]);

    // Debounced search effect
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearchValue(searchValue);
        }, 500);

        return () => clearTimeout(timer);
    }, [searchValue]);

    // Fetch global status counts
    useEffect(() => {
        axiosInstance.get('/api/asns/stats')
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

    // Fetch filter statistics (RIR and Country counts)
    useEffect(() => {
        axiosInstance.get('/api/asns/filter-stats')
            .then(res => {
                setGlobalRIRCounts(res.data.rir_counts || { total: 0 });
                setGlobalCountryCounts(res.data.country_counts || { total: 0 });
            })
            .catch(() => {
                setGlobalRIRCounts({ total: 0 });
                setGlobalCountryCounts({ total: 0 });
            });
    }, [refresh]);

    // Fetch ASNs with server-side filtering, sorting, and pagination
    useEffect(() => {
        const fetchASNs = async () => {
            setLoading(true);
            try {
                const response = await axiosInstance.get('/api/asns', {
                    params: {
                        page: page + 1,
                        limit: rowsPerPage,
                        status: filterStatus || undefined,
                        rir: filterRIR || undefined,
                        country: filterCountry || undefined,
                        search: debouncedSearchValue || undefined,
                        orderBy,
                        order,
                    }
                });
                if (response.data && response.data.items) {
                    setASNs(response.data.items);
                    setTotal(response.data.total || response.data.items.length);
                } else {
                    setASNs([]);
                    setTotal(0);
                }
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch ASNs');
                setLoading(false);
            }
        };
        fetchASNs();
    }, [refresh, page, rowsPerPage, filterStatus, filterRIR, filterCountry, debouncedSearchValue, orderBy, order]);

    const handleSubmit = useCallback(async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');

        try {
            if (editId) {
                await axiosInstance.put(`/api/asn/${editId}`, {
                    asn,
                    rir,
                    domain,
                    cc: country,
                    asname: name,
                    status,
                    source
                });
                setMessage('ASN updated successfully!');
            } else {
                await axiosInstance.post('/api/asn', {
                    asn,
                    rir,
                    domain,
                    cc: country,
                    asname: name,
                    status,
                    source
                });
                setMessage('ASN added successfully!');
            }

            // Reset form
            setASN('');
            setRIR('');
            setDomain('');
            setCountry('');
            setName('');
            setStatus('denied');
            setSource('manual');
            setEditId(null);
            setRefresh(prev => !prev);
        } catch (err) {
            const errorMessage = err.response?.data?.error || err.message || 'Error saving ASN';
            setError(errorMessage);
        }
    }, [asn, rir, domain, country, name, status, source, editId]);

    const handleEdit = useCallback((asnData) => {
        setASN(asnData.asn);
        setRIR(asnData.rir);
        setDomain(asnData.domain);
        setCountry(asnData.cc);
        setName(asnData.asname);
        setStatus(asnData.status);
        setSource('manual'); // Always set to manual when editing
        setEditId(asnData.id);
    }, []);

    const handleDelete = useCallback(async (id) => {
        if (window.confirm('Are you sure you want to delete this ASN?')) {
            try {
                await axiosInstance.delete(`/api/asn/${id}`);
                setMessage('ASN deleted successfully!');
                setRefresh(prev => !prev);
            } catch (err) {
                setError('Failed to delete ASN');
            }
        }
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

    const handleRIRFilterChange = useCallback((e) => {
        setFilterRIR(e.target.value);
        setPage(0); // Reset to first page when filtering
    }, []);

    const handleCountryFilterChange = useCallback((e) => {
        setFilterCountry(e.target.value);
        setPage(0); // Reset to first page when filtering
    }, []);

    const handleReset = useCallback(() => {
        setSearchValue('');
        setFilterStatus('');
        setFilterRIR('');
        setFilterCountry('');
        setPage(0);
    }, []);

    const handleSearchFocus = useCallback(() => {
        // Focus handler for search input
    }, []);

    const handleSearchBlur = useCallback(() => {
        // Blur handler for search input
    }, []);

    const handleCancelEdit = useCallback(() => {
        setASN('');
        setRIR('');
        setDomain('');
        setCountry('');
        setName('');
        setStatus('denied');
        setSource('manual');
        setEditId(null);
    }, []);

    const handleASNChange = useCallback((e) => {
        setASN(e.target.value);
    }, []);

    const handleNameChange = useCallback((e) => {
        setName(e.target.value);
    }, []);

    const handleStatusChangeForm = useCallback((e) => {
        setStatus(e.target.value);
    }, []);

    const handleRIRChange = useCallback((e) => {
        setRIR(e.target.value);
    }, []);

    const handleDomainChange = useCallback((e) => {
        setDomain(e.target.value);
    }, []);

    const handleCountryChange = useCallback((e) => {
        setCountry(e.target.value);
    }, []);

    const handleSourceChange = useCallback((e) => {
        setSource(e.target.value);
    }, []);

    const handleSpamhausImport = useCallback(async () => {
        if (window.confirm('This will import ASN data from Spamhaus ASN-DROP list. Existing Spamhaus records will be replaced. Continue?')) {
            try {
                setLoading(true);
                await axiosInstance.post('/api/asns/import-spamhaus');
                setMessage('Spamhaus ASN-DROP data imported successfully!');
                setRefresh(prev => !prev);
            } catch (err) {
                const errorMessage = err.response?.data?.error || err.message || 'Failed to import Spamhaus data';
                setError(errorMessage);
            } finally {
                setLoading(false);
            }
        }
    }, []);

    const handleSpamhausStats = useCallback(async () => {
        try {
            const response = await axiosInstance.get('/api/asns/spamhaus-stats');
            const stats = response.data;
            const statsMessage = `Spamhaus ASNs: ${stats.total_spamhaus_asns || 0}\nLast Sync: ${stats.last_sync ? new Date(stats.last_sync).toLocaleString() : 'Never'}`;
            alert(statsMessage);
        } catch (err) {
            setError('Failed to get Spamhaus stats');
        }
    }, []);



    const handleSpamhausStatus = useCallback(async () => {
        try {
            const response = await axiosInstance.get('/api/asns/spamhaus-status');
            const status = response.data;
            const autoImportStatus = status.auto_import_enabled ? 'Enabled' : 'Disabled';
            const statusMessage = `Auto Import: ${autoImportStatus}\nImport Running: ${status.is_running ? 'Yes' : 'No'}\nNext Scheduled: ${status.next_scheduled}\nTime Until Next: ${status.next_scheduled_relative}`;
            alert(statusMessage);
        } catch (err) {
            setError('Failed to get Spamhaus import status');
        }
    }, []);

    return (
        <Box sx={{ maxWidth: 1200, mx: 'auto', mt: 4, p: 2 }}>
            <Typography variant="h4" gutterBottom>
                ASN Management
            </Typography>

            {/* Spamhaus Import Buttons */}
            <Box sx={{ mb: 3 }}>
                <Button 
                    variant="contained" 
                    color="primary" 
                    onClick={handleSpamhausImport}
                    sx={{ mr: 2 }}
                >
                    Import Spamhaus ASN-DROP
                </Button>
                <Button 
                    variant="outlined" 
                    onClick={handleSpamhausStats}
                    sx={{ mr: 2 }}
                >
                    View Spamhaus Stats
                </Button>
                <Button 
                    variant="outlined" 
                    onClick={handleSpamhausStatus}
                >
                    Check Import Status
                </Button>
            </Box>

            {/* More Resources Section */}
            <Box sx={{ mb: 3 }}>
                <Button
                    variant="outlined"
                    onClick={() => setShowResources(!showResources)}
                    startIcon={showResources ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                    sx={{ 
                        mb: 1,
                        borderColor: showResources ? '#1976d2' : '#ccc',
                        color: showResources ? '#1976d2' : '#666',
                        backgroundColor: showResources ? '#e3f2fd' : 'transparent',
                        '&:hover': {
                            backgroundColor: showResources ? '#bbdefb' : '#f5f5f5',
                            borderColor: '#1976d2',
                            color: '#1976d2'
                        },
                        fontWeight: showResources ? 600 : 400,
                        transition: 'all 0.2s ease'
                    }}
                >
                    📚 More Resources
                </Button>
                {showResources && (
                    <Paper 
                        elevation={2} 
                        sx={{ 
                            p: 2, 
                            backgroundColor: '#f8f9fa',
                            border: '1px solid #e0e0e0',
                            borderRadius: 2,
                            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                            animation: 'slideDown 0.3s ease-out'
                        }}
                    >
                        <Typography variant="subtitle1" gutterBottom sx={{ color: '#1976d2', fontWeight: 600 }}>
                            External ASN Blacklist Resources:
                        </Typography>
                        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                            <Link 
                                href="https://cleantalk.org/blacklists/asn" 
                                target="_blank" 
                                rel="noopener noreferrer"
                                sx={{ 
                                    display: 'flex', 
                                    alignItems: 'center', 
                                    gap: 1,
                                    p: 1.5,
                                    borderRadius: 1,
                                    backgroundColor: 'white',
                                    textDecoration: 'none',
                                    color: '#1976d2',
                                    border: '1px solid #e0e0e0',
                                    transition: 'all 0.2s ease',
                                    '&:hover': {
                                        backgroundColor: '#f5f5f5',
                                        borderColor: '#1976d2',
                                        boxShadow: '0 2px 4px rgba(25, 118, 210, 0.2)',
                                        transform: 'translateY(-1px)'
                                    }
                                }}
                            >
                                <OpenInNewIcon fontSize="small" />
                                CleanTalk ASN Blacklist
                            </Link>
                            <Link 
                                href="https://ipapi.is/most-abusive-asn.html" 
                                target="_blank" 
                                rel="noopener noreferrer"
                                sx={{ 
                                    display: 'flex', 
                                    alignItems: 'center', 
                                    gap: 1,
                                    p: 1.5,
                                    borderRadius: 1,
                                    backgroundColor: 'white',
                                    textDecoration: 'none',
                                    color: '#1976d2',
                                    border: '1px solid #e0e0e0',
                                    transition: 'all 0.2s ease',
                                    '&:hover': {
                                        backgroundColor: '#f5f5f5',
                                        borderColor: '#1976d2',
                                        boxShadow: '0 2px 4px rgba(25, 118, 210, 0.2)',
                                        transform: 'translateY(-1px)'
                                    }
                                }}
                            >
                                <OpenInNewIcon fontSize="small" />
                                IPAPI Most Abusive ASN List
                            </Link>
                        </Box>
                    </Paper>
                )}
            </Box>

            {/* Form Section */}
            <Paper sx={{ p: 3, mb: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>
                    {editId ? 'Edit ASN' : 'Add New ASN'}
                </Typography>
                
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

                <ASNFormComponent
                    asn={asn}
                    rir={rir}
                    domain={domain}
                    country={country}
                    name={name}
                    status={status}
                    source={source}
                    message={message}
                    error={error}
                    editId={editId}
                    onASNChange={handleASNChange}
                    onRIRChange={handleRIRChange}
                    onDomainChange={handleDomainChange}
                    onCountryChange={handleCountryChange}
                    onNameChange={handleNameChange}
                    onStatusChange={handleStatusChangeForm}
                    onSourceChange={handleSourceChange}
                    onSubmit={handleSubmit}
                    onCancelEdit={handleCancelEdit}
                />
            </Paper>

            {/* Filter Controls */}
            <FilterControls
                searchValue={searchValue}
                filterStatus={filterStatus}
                filterRIR={filterRIR}
                filterCountry={filterCountry}
                globalStatusCounts={globalStatusCounts}
                globalRIRCounts={globalRIRCounts}
                globalCountryCounts={globalCountryCounts}
                onSearchChange={handleSearchChange}
                onStatusChange={handleStatusChange}
                onRIRChange={handleRIRFilterChange}
                onCountryChange={handleCountryFilterChange}
                onReset={handleReset}
                onSearchFocus={handleSearchFocus}
                onSearchBlur={handleSearchBlur}
                searchInputRef={searchInputRef}
            />

            {/* Table */}
            <ASNTable
                asns={asns}
                loading={loading}
                orderBy={orderBy}
                order={order}
                page={page}
                rowsPerPage={rowsPerPage}
                total={total}
                onSort={handleSort}
                onEdit={handleEdit}
                onDelete={handleDelete}
                onChangePage={handleChangePage}
                onChangeRowsPerPage={handleChangeRowsPerPage}
            />
        </Box>
    );
};

export default ASNForm; 