import React, { useState, useEffect, useCallback, memo } from 'react';
import {
    Box, Paper, Typography, TextField, Button, Alert, Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
    TablePagination, TableSortLabel, IconButton, FormControl, InputLabel, Select, MenuItem, Checkbox, FormControlLabel,
    Chip, Dialog, DialogTitle, DialogContent, DialogActions, List, ListItem, ListItemText, ListItemSecondaryAction,
    Accordion, AccordionSummary, AccordionDetails
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import AddIcon from '@mui/icons-material/Add';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import axios from '../axiosConfig';
import { useLocation } from 'react-router-dom';
import InfiniteScrollCharsetTable from './InfiniteScrollCharsetTable';


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
                label="Charset Filter"
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
const CharsetFormComponent = memo(({ onSubmit, charset, setCharset, status, setStatus, editId, setEditId }) => {
    return (
        <Box component="form" onSubmit={onSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'stretch', mb: 2 }}>
            <TextField
                label="Charset"
                value={charset}
                onChange={(e) => setCharset(e.target.value)}
                placeholder="Enter charset (e.g. UTF-8)"
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
            <Button type="submit" variant="contained" color="primary">
                {editId ? 'Update Charset' : 'Add Charset'}
            </Button>
            {editId && (
                <Button variant="outlined" color="secondary" onClick={() => { setEditId(null); setCharset(''); setStatus('denied'); }}>
                    Cancel Edit
                </Button>
            )}
        </Box>
    );
}, (prevProps, nextProps) => {
    // Only re-render when form data changes
    return prevProps.charset === nextProps.charset && 
           prevProps.status === nextProps.status && 
           prevProps.editId === nextProps.editId;
});

// Separate table component that only re-renders when data changes
const CharsetTable = memo(({ 
    charsets, 
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
                                        active={orderBy === 'charset'}
                                        direction={orderBy === 'charset' ? order : 'asc'}
                                        onClick={() => onSort('charset')}
                                    >
                                        Charset
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
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {charsets.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={3} align="center">No Charsets</TableCell>
                                </TableRow>
                            ) : (
                                charsets.map(charsetItem => (
                                    <TableRow key={charsetItem.id}>
                                        <TableCell>{charsetItem.charset}</TableCell>
                                        <TableCell>{charsetItem.status}</TableCell>
                                        <TableCell>
                                            <IconButton onClick={() => onEdit(charsetItem)} size="small"><EditIcon /></IconButton>
                                            <IconButton onClick={() => onDelete(charsetItem.id)} size="small" color="error"><DeleteIcon /></IconButton>
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

// Custom Fields Manager Component
const CustomFieldsManager = memo(({ 
    standardFields, 
    customFields, 
    onStandardFieldToggle, 
    onCustomFieldAdd, 
    onCustomFieldDelete 
}) => {
    const [newFieldDialog, setNewFieldDialog] = useState(false);
    const [newFieldName, setNewFieldName] = useState('');

    const handleAddField = () => {
        if (newFieldName.trim()) {
            onCustomFieldAdd(newFieldName.trim());
            setNewFieldName('');
            setNewFieldDialog(false);
        }
    };

    const handleDeleteField = (fieldName) => {
        if (window.confirm(`Delete custom field "${fieldName}"?`)) {
            onCustomFieldDelete(fieldName);
        }
    };

    const [showCharsetFields, setShowCharsetFields] = useState(false);

    return (
        <Box sx={{ mb: 3 }}>
            <Button
                variant="outlined"
                onClick={() => setShowCharsetFields(!showCharsetFields)}
                startIcon={showCharsetFields ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                sx={{ 
                    mb: 1,
                    borderColor: showCharsetFields ? '#1976d2' : '#ccc',
                    color: showCharsetFields ? '#1976d2' : '#666',
                    backgroundColor: showCharsetFields ? '#e3f2fd' : 'transparent',
                    '&:hover': {
                        backgroundColor: showCharsetFields ? '#bbdefb' : '#f5f5f5',
                        borderColor: '#1976d2',
                        color: '#1976d2'
                    },
                    fontWeight: showCharsetFields ? 600 : 400,
                    transition: 'all 0.2s ease'
                }}
            >
                🔧 Charset Filter Fields
            </Button>
            {showCharsetFields && (
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
                        Select which fields should be checked for charset detection in the /filter endpoint.
                    </Typography>

                {/* Standard Fields */}
                <Typography variant="subtitle1" sx={{ mb: 2, fontWeight: 600, color: '#1976d2' }}>
                    Standard Fields
                </Typography>
                <Box sx={{ mb: 3, display: 'flex', flexWrap: 'wrap', gap: 3 }}>
                    {standardFields.map(field => (
                        <FormControlLabel
                            key={field.name}
                            control={
                                <Checkbox
                                    checked={field.enabled}
                                    onChange={() => onStandardFieldToggle(field.name)}
                                    disabled={false}
                                    sx={{
                                        color: '#1976d2',
                                        '&.Mui-checked': {
                                            color: '#1976d2',
                                        },
                                    }}
                                />
                            }
                            label={
                                <Box sx={{ 
                                    display: 'flex', 
                                    alignItems: 'center', 
                                    gap: 1,
                                    p: 1,
                                    borderRadius: 1,
                                    backgroundColor: field.enabled ? '#e3f2fd' : 'transparent',
                                    border: field.enabled ? '1px solid #1976d2' : '1px solid transparent',
                                    transition: 'all 0.2s ease'
                                }}>
                                    {field.name}
                                    <Chip 
                                        label="Standard" 
                                        size="small" 
                                        color="primary" 
                                        variant="outlined"
                                        sx={{ 
                                            borderColor: '#1976d2',
                                            color: '#1976d2'
                                        }}
                                    />
                                </Box>
                            }
                            sx={{
                                
                            }}
                        />
                    ))}
                </Box>

                {/* Custom Fields */}
                <Typography variant="subtitle1" sx={{ mb: 2, fontWeight: 600, color: '#1976d2' }}>
                    Custom Fields
                </Typography>
                <Box sx={{ mb: 3 }}>
                    {customFields.length === 0 ? (
                        <Typography variant="body2" color="text.secondary" sx={{ 
                            p: 2, 
                            backgroundColor: '#f5f5f5', 
                            borderRadius: 1,
                            border: '1px dashed #ccc',
                            textAlign: 'center'
                        }}>
                            No custom fields added yet.
                        </Typography>
                    ) : (
                        <List dense>
                            {customFields.map(field => (
                                <ListItem key={field.name} sx={{ 
                                    py: 1,
                                    mb: 1,
                                    backgroundColor: 'white',
                                    borderRadius: 1,
                                    border: '1px solid #e0e0e0',
                                    '&:hover': {
                                        backgroundColor: '#f5f5f5',
                                        borderColor: '#1976d2',
                                        boxShadow: '0 2px 4px rgba(25, 118, 210, 0.2)',
                                        transform: 'translateY(-1px)'
                                    },
                                    transition: 'all 0.2s ease'
                                }}>
                                    <ListItemText 
                                        primary={field.name}
                                        secondary="Custom field"
                                        primaryTypographyProps={{ fontWeight: 500 }}
                                        secondaryTypographyProps={{ color: '#666' }}
                                    />
                                    <ListItemSecondaryAction>
                                        <IconButton 
                                            edge="end" 
                                            onClick={() => handleDeleteField(field.name)}
                                            size="small"
                                            color="error"
                                            sx={{
                                                '&:hover': {
                                                    backgroundColor: '#ffebee',
                                                    transform: 'scale(1.1)'
                                                },
                                                transition: 'all 0.2s ease'
                                            }}
                                        >
                                            <DeleteIcon />
                                        </IconButton>
                                    </ListItemSecondaryAction>
                                </ListItem>
                            ))}
                        </List>
                    )}
                </Box>

                {/* Add Custom Field Button */}
                <Button
                    variant="outlined"
                    startIcon={<AddIcon />}
                    onClick={() => setNewFieldDialog(true)}
                    size="small"
                    sx={{
                        borderColor: '#1976d2',
                        color: '#1976d2',
                        '&:hover': {
                            backgroundColor: '#e3f2fd',
                            borderColor: '#1565c0',
                            color: '#1565c0'
                        },
                        fontWeight: 500,
                        transition: 'all 0.2s ease'
                    }}
                >
                    Add Custom Field
                </Button>

                {/* Add Custom Field Dialog */}
                <Dialog open={newFieldDialog} onClose={() => setNewFieldDialog(false)} maxWidth="sm" fullWidth>
                    <DialogTitle>Add Custom Field</DialogTitle>
                    <DialogContent>
                        <TextField
                            autoFocus
                            margin="dense"
                            label="Field Name"
                            fullWidth
                            variant="outlined"
                            value={newFieldName}
                            onChange={(e) => setNewFieldName(e.target.value)}
                            placeholder="Enter field name (e.g., content, description, notes)"
                            helperText="This field will be checked for charset detection in the /filter endpoint"
                        />
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => setNewFieldDialog(false)}>Cancel</Button>
                        <Button onClick={handleAddField} variant="contained" disabled={!newFieldName.trim()}>
                            Add Field
                        </Button>
                    </DialogActions>
                </Dialog>
                </Paper>
                )}
            </Box>
        );
    }, (prevProps, nextProps) => {
        return JSON.stringify(prevProps) === JSON.stringify(nextProps);
    });

const CharsetForm = () => {
    const [charset, setCharset] = useState('');
    const [status, setStatus] = useState('denied');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [refresh, setRefresh] = useState(false);
    const [charsets, setCharsets] = useState([]);
    const [editId, setEditId] = useState(null);
    
    // Filtering and infinite scroll state
    const [loading, setLoading] = useState(true);
    const [infiniteLoading, setInfiniteLoading] = useState(false);
    const [filterStatus, setFilterStatus] = useState('');
    const [orderBy, setOrderBy] = useState('charset');
    const [order, setOrder] = useState('asc');
    const [total, setTotal] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [globalStatusCounts, setGlobalStatusCounts] = useState({ allowed: 0, denied: 0, whitelisted: 0, total: 0 });
    const location = useLocation();

    // Debounced search state
    const [searchValue, setSearchValue] = useState('');
    const [debouncedSearchValue, setDebouncedSearchValue] = useState('');

    // Custom fields management state
    const [standardFields, setStandardFields] = useState([
        { name: 'username', enabled: true },
        { name: 'email', enabled: true },
        { name: 'user_agent', enabled: true }
    ]);
    const [customFields, setCustomFields] = useState([]);

    // Load field configuration from backend
    useEffect(() => {
        const loadFieldConfig = async () => {
            try {
                const response = await axios.get('/api/charset-fields');
                if (response.data) {
                    setStandardFields(response.data.standard_fields || []);
                    setCustomFields(response.data.custom_fields || []);
                }
            } catch (err) {
                console.error('Failed to load field configuration:', err);
            }
        };
        loadFieldConfig();
    }, [refresh]);

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
        axios.get('/api/charsets/stats')
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



    // Fetch charsets with server-side filtering, sorting, and infinite scroll
    useEffect(() => {
        const fetchCharsets = async () => {
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
                
                const response = await axios.get('/api/charsets', { params });
                
                if (response.data && response.data.items) {
                    // Replace charsets for initial load
                    setCharsets(response.data.items);
                    setHasMore(response.data.items.length === 25);
                    setTotal(response.data.total || 0);
                } else {
                    setCharsets([]);
                    setHasMore(false);
                    setTotal(0);
                }
            } catch (err) {
                setError('Failed to fetch charsets');
            } finally {
                setLoading(false);
            }
        };
        
        // Reset infinite scroll state when filters change
        setHasMore(true);
        setCharsets([]);
        
        fetchCharsets();
    }, [refresh, filterStatus, debouncedSearchValue, orderBy, order]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            if (editId) {
                await axios.put(`/api/charset/${editId}`, { charset, status });
                setMessage('Charset updated successfully');
            } else {
                await axios.post('/api/charset', { charset, status });
                setMessage('Charset added successfully');
            }
            setCharset('');
            setStatus('denied');
            setEditId(null);
            setRefresh(r => !r);
        } catch (err) {
            setError('Error saving charset');
        }
    };



    const handleSearchChange = useCallback((e) => {
        setSearchValue(e.target.value);
        setCharsets([]); // Reset charsets when searching
        setHasMore(true);
    }, []);

    const handleStatusChange = useCallback((e) => {
        setFilterStatus(e.target.value);
        setCharsets([]); // Reset charsets when filtering
        setHasMore(true);
    }, []);

    const handleReset = useCallback(() => {
        setSearchValue('');
        setFilterStatus('');
        setCharsets([]);
        setHasMore(true);
    }, []);

    const handleLoadMore = useCallback(() => {
        if (!infiniteLoading && hasMore) {
            const fetchCharsets = async () => {
                setInfiniteLoading(true);
                try {
                    const params = {
                        page: Math.floor(charsets.length / 25) + 1,
                        limit: 25,
                        status: filterStatus || undefined,
                        search: debouncedSearchValue || undefined,
                        orderBy,
                        order,
                    };
                    
                    const response = await axios.get('/api/charsets', { params });
                    
                    if (response.data && response.data.items) {
                        setCharsets(prev => [...prev, ...response.data.items]);
                        setHasMore(response.data.items.length === 25);
                    } else {
                        setHasMore(false);
                    }
                } catch (err) {
                    setError('Failed to fetch more charsets');
                } finally {
                    setInfiniteLoading(false);
                }
            };
            fetchCharsets();
        }
    }, [infiniteLoading, hasMore, charsets.length, filterStatus, debouncedSearchValue, orderBy, order]);

    const handleSort = useCallback((field) => {
        if (orderBy === field) {
            setOrder(order === 'asc' ? 'desc' : 'asc');
        } else {
            setOrderBy(field);
            setOrder('asc');
        }
        // Reset charsets when sorting changes
        setCharsets([]);
        setHasMore(true);
    }, [orderBy, order]);

    const handleEdit = useCallback((charsetItem) => {
        setCharset(charsetItem.charset);
        setStatus(charsetItem.status);
        setEditId(charsetItem.id);
    }, []);

    const handleDelete = useCallback(async (id) => {
        if (!window.confirm('Delete this charset?')) return;
        try {
            await axios.delete(`/api/charset/${id}`);
            setMessage('Charset deleted');
            setRefresh(r => !r);
        } catch {
            setError('Error deleting charset');
        }
    }, []);

    // Custom fields management handlers
    const handleStandardFieldToggle = useCallback(async (fieldName) => {
        try {
            await axios.post('/api/charset-fields/toggle-standard', { field_name: fieldName });
            setStandardFields(prev => 
                prev.map(field => 
                    field.name === fieldName 
                        ? { ...field, enabled: !field.enabled }
                        : field
                )
            );
            setMessage(`Standard field "${fieldName}" toggled successfully`);
        } catch (err) {
            setError(`Failed to toggle standard field "${fieldName}"`);
        }
    }, []);

    const handleCustomFieldAdd = useCallback(async (fieldName) => {
        try {
            await axios.post('/api/charset-fields/add-custom', { field_name: fieldName });
            setCustomFields(prev => [...prev, { name: fieldName, enabled: true, type: 'custom' }]);
            setMessage(`Custom field "${fieldName}" added successfully`);
        } catch (err) {
            setError(err.response?.data?.error || `Failed to add custom field "${fieldName}"`);
        }
    }, []);

    const handleCustomFieldDelete = useCallback(async (fieldName) => {
        try {
            await axios.delete(`/api/charset-fields/custom/${fieldName}`);
            setCustomFields(prev => prev.filter(field => field.name !== fieldName));
            setMessage(`Custom field "${fieldName}" deleted successfully`);
        } catch (err) {
            setError(`Failed to delete custom field "${fieldName}"`);
        }
    }, []);


    return (
        <Box sx={{ maxWidth: 1200, mx: 'auto', mt: 4, p: 2 }}>
            <Typography variant="h4" gutterBottom>
                Charset Management
            </Typography>

            {/* Charset Filter Fields Section */}
            <CustomFieldsManager
                standardFields={standardFields}
                customFields={customFields}
                onStandardFieldToggle={handleStandardFieldToggle}
                onCustomFieldAdd={handleCustomFieldAdd}
                onCustomFieldDelete={handleCustomFieldDelete}
            />

            {/* Form Section */}
            <Paper sx={{ p: 3, mb: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>
                    Add or Modify Charset
                </Typography>
                
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

                <CharsetFormComponent 
                    onSubmit={handleSubmit}
                    charset={charset}
                    setCharset={setCharset}
                    status={status}
                    setStatus={setStatus}
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
            <InfiniteScrollCharsetTable 
                charsets={charsets}
                loading={infiniteLoading}
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

export default CharsetForm; 