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
import axios from '../axiosConfig';
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
                        labelRowsPerPage="Einträge pro Seite:"
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
                        labelRowsPerPage="Einträge pro Seite:"
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

    return (
        <Accordion sx={{ 
            mt: 3,
            backgroundColor: '#f8f9fa',
            border: '1px solid #e3e6ea',
            borderRadius: 1,
            '&:before': {
                display: 'none',
            },
            '&.Mui-expanded': {
                margin: '24px 0',
            }
        }}>
            <AccordionSummary
                expandIcon={<ExpandMoreIcon />}
                aria-controls="charset-fields-content"
                id="charset-fields-header"
                sx={{
                    backgroundColor: '#f1f3f4',
                    borderBottom: '1px solid #e3e6ea',
                    '&:hover': {
                        backgroundColor: '#e8eaed',
                    },
                    '&.Mui-expanded': {
                        backgroundColor: '#e8eaed',
                        minHeight: '48px',
                    }
                }}
            >
                <Typography variant="h6" sx={{ color: '#1a73e8', fontWeight: 500 }}>
                    Charset Filter Fields
                </Typography>
            </AccordionSummary>
            <AccordionDetails sx={{ backgroundColor: '#ffffff', p: 3 }}>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    Select which fields should be checked for charset detection in the /filter endpoint.
                </Typography>

                {/* Standard Fields */}
                <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 'bold' }}>
                    Standard Fields
                </Typography>
                <Box sx={{ mb: 2 }}>
                    {standardFields.map(field => (
                        <FormControlLabel
                            key={field.name}
                            control={
                                <Checkbox
                                    checked={field.enabled}
                                    onChange={() => onStandardFieldToggle(field.name)}
                                    disabled={false}
                                />
                            }
                            label={
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    {field.name}
                                    <Chip label="Standard" size="small" color="primary" variant="outlined" />
                                </Box>
                            }
                        />
                    ))}
                </Box>

                {/* Custom Fields */}
                <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 'bold' }}>
                    Custom Fields
                </Typography>
                <Box sx={{ mb: 2 }}>
                    {customFields.length === 0 ? (
                        <Typography variant="body2" color="text.secondary">
                            No custom fields added yet.
                        </Typography>
                    ) : (
                        <List dense>
                            {customFields.map(field => (
                                <ListItem key={field.name} sx={{ py: 0.5 }}>
                                    <ListItemText 
                                        primary={field.name}
                                        secondary="Custom field"
                                    />
                                    <ListItemSecondaryAction>
                                        <IconButton 
                                            edge="end" 
                                            onClick={() => handleDeleteField(field.name)}
                                            size="small"
                                            color="error"
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
            </AccordionDetails>
        </Accordion>
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
    
    // Filtering and pagination state
    const [loading, setLoading] = useState(true);
    const [filterStatus, setFilterStatus] = useState('');
    const [orderBy, setOrderBy] = useState('charset');
    const [order, setOrder] = useState('asc');
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [total, setTotal] = useState(0);
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



    // Fetch charsets with server-side filtering, sorting, and pagination
    useEffect(() => {
        const fetchCharsets = async () => {
            setLoading(true);
            try {
                const response = await axios.get('/api/charsets', {
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
                    setCharsets(response.data.items);
                    setTotal(response.data.total || response.data.items.length);
                } else {
                    setCharsets([]);
                    setTotal(0);
                }
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch charsets');
                setLoading(false);
            }
        };
        fetchCharsets();
    }, [refresh, page, rowsPerPage, filterStatus, debouncedSearchValue, orderBy, order]);

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
        <Box sx={{ maxWidth: 700, mx: 'auto', mt: 4 }}>
            <Paper sx={{ p: 3 }} elevation={3}>
                <Typography variant="h5" gutterBottom>Charset Management</Typography>
                <CharsetFormComponent 
                    onSubmit={handleSubmit}
                    charset={charset}
                    setCharset={setCharset}
                    status={status}
                    setStatus={setStatus}
                    editId={editId}
                    setEditId={setEditId}
                />
                {message && <Alert severity="success" sx={{ mb: 2 }}>{message}</Alert>}
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
                
                {/* Custom Fields Management */}
                <CustomFieldsManager
                    standardFields={standardFields}
                    customFields={customFields}
                    onStandardFieldToggle={handleStandardFieldToggle}
                    onCustomFieldAdd={handleCustomFieldAdd}
                    onCustomFieldDelete={handleCustomFieldDelete}
                />
                
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
                
                <CharsetTable 
                    charsets={charsets}
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

export default CharsetForm; 