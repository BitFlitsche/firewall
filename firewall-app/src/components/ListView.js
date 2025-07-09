import React, { useState, useEffect } from 'react';
import axiosInstance from '../axiosConfig';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import MenuItem from '@mui/material/MenuItem';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import TableSortLabel from '@mui/material/TableSortLabel';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Select from '@mui/material/Select';
import Button from '@mui/material/Button';
import TablePagination from '@mui/material/TablePagination';
import { useLocation } from 'react-router-dom';

const statusOptions = [
  { value: '', label: 'All' },
  { value: 'allowed', label: 'Allowed' },
  { value: 'denied', label: 'Denied' },
  { value: 'whitelisted', label: 'Whitelisted' },
];

const getValueField = (endpoint) => {
  if (endpoint === '/ips') return 'Address';
  if (endpoint === '/emails') return 'Address';
  if (endpoint === '/user-agents') return 'UserAgent';
  if (endpoint === '/countries') return 'Code';
  return '';
};

const getValueHeader = (endpoint) => {
  if (endpoint === '/ips') return 'IP Address';
  if (endpoint === '/emails') return 'Email Address';
  if (endpoint === '/user-agents') return 'User Agent';
  if (endpoint === '/countries') return 'Country Code';
  return 'Value';
};

const getStatsEndpoint = (endpoint) => {
  if (endpoint === '/ips') return '/ips/stats';
  if (endpoint === '/emails') return '/emails/stats';
  if (endpoint === '/user-agents') return '/user-agents/stats';
  if (endpoint === '/countries') return '/countries/stats';
  return null;
};

const ListView = ({ endpoint, title, refresh }) => {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filterValue, setFilterValue] = useState('');
  const [filterStatus, setFilterStatus] = useState('');
  const [orderBy, setOrderBy] = useState('ID');
  const [order, setOrder] = useState('desc');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const location = useLocation();
  const [globalStatusCounts, setGlobalStatusCounts] = useState({ allowed: 0, denied: 0, whitelisted: 0, total: 0 });

  // Set initial filterStatus from query param
  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const status = params.get('status');
    if (status && ['allowed','denied','whitelisted'].includes(status)) {
      setFilterStatus(status);
    }
  }, [location.search]);

  // Lade globale Status-Counts
  useEffect(() => {
    const statsEndpoint = getStatsEndpoint(endpoint);
    if (!statsEndpoint) return;
    axiosInstance.get(statsEndpoint)
      .then(res => {
        setGlobalStatusCounts({
          allowed: res.data.allowed || 0,
          denied: res.data.denied || 0,
          whitelisted: res.data.whitelisted || 0,
          total: res.data.total || 0,
        });
      })
      .catch(() => setGlobalStatusCounts({ allowed: 0, denied: 0, whitelisted: 0, total: 0 }));
  }, [endpoint, refresh]);

  const valueField = getValueField(endpoint);
  const valueHeader = getValueHeader(endpoint);

  // Serverseitige Daten für /ips, /emails, /user-agents, /countries
  useEffect(() => {
    const fetchItems = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await axiosInstance.get(endpoint, (endpoint === '/ips' || endpoint === '/emails' || endpoint === '/user-agents' || endpoint === '/countries') ? {
          params: {
            page: page + 1,
            limit: rowsPerPage,
            status: filterStatus || undefined,
            search: filterValue || undefined,
            orderBy,
            order,
          }
        } : undefined);
        if (response.data && Array.isArray(response.data)) {
          setItems(response.data);
          setTotal(response.data.length);
        } else if (response.data && response.data.items) {
          setItems(response.data.items);
          setTotal(response.data.total || response.data.items.length);
        } else {
          setItems([]);
          setTotal(0);
        }
        setLoading(false);
      } catch (err) {
        setError('Failed to fetch items');
        setLoading(false);
      }
    };
    fetchItems();
    // eslint-disable-next-line
  }, [endpoint, refresh, page, rowsPerPage, filterStatus, filterValue, orderBy, order]);

  // Für Endpunkte mit serverseitiger Pagination keinerlei clientseitige Filterung/Sortierung/Slicing
  let paginatedItems = items;
  let sortedItems = items; // immer definiert
  if (!(endpoint === '/ips' || endpoint === '/emails' || endpoint === '/user-agents' || endpoint === '/countries')) {
    // Legacy/sonstige Endpunkte: clientseitige Filterung/Sortierung/Slicing
    let filteredItems = items.filter((item) => {
      const value = (item[valueField] || '').toLowerCase();
      const status = (item.Status || '').toLowerCase();
      const valueMatch = filterValue === '' || value.includes(filterValue.toLowerCase());
      const statusMatch = filterStatus === '' || status === filterStatus;
      return valueMatch && statusMatch;
    });
    sortedItems = [...filteredItems].sort((a, b) => {
      let aValue = a[orderBy];
      let bValue = b[orderBy];
      if (orderBy === 'ID') {
        aValue = aValue || 0;
        bValue = bValue || 0;
        return order === 'asc' ? aValue - bValue : bValue - aValue;
      } else {
        aValue = (aValue || '').toString().toLowerCase();
        bValue = (bValue || '').toString().toLowerCase();
        if (aValue < bValue) return order === 'asc' ? -1 : 1;
        if (aValue > bValue) return order === 'asc' ? 1 : -1;
        return 0;
      }
    });
    paginatedItems = sortedItems.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);
  }

  // Für /ips, /emails, /user-agents, /countries: total immer aus API übernehmen
  const [total, setTotal] = useState(0);
  useEffect(() => {
    // Nur für Endpunkte ohne serverseitige Pagination (z.B. falls noch legacy)
    if (endpoint !== '/ips' && endpoint !== '/emails' && endpoint !== '/user-agents' && endpoint !== '/countries') {
      setTotal(sortedItems.length);
    }
    // eslint-disable-next-line
  }, [sortedItems, endpoint]);

  // Status-Counts berechnen
  const statusCounts = items.reduce((acc, item) => {
    const status = (item.Status || '').toLowerCase();
    if (!acc[status]) acc[status] = 0;
    acc[status]++;
    return acc;
  }, {});
  const totalCount = items.length;

  const handleSort = (field) => {
    if (orderBy === field) {
      setOrder(order === 'asc' ? 'desc' : 'asc');
    } else {
      setOrderBy(field);
      setOrder('asc');
    }
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <Box>
      <Box sx={{ display: 'flex', gap: 2, mb: 2, flexWrap: 'wrap', alignItems: 'center' }}>
        <TextField
          label={valueHeader + ' Filter'}
          value={filterValue}
          onChange={(e) => setFilterValue(e.target.value)}
          size="small"
        />
        <FormControl size="small" sx={{ minWidth: 140 }}>
          <InputLabel shrink>Status</InputLabel>
          <Select
            value={filterStatus}
            label="Status"
            onChange={(e) => setFilterStatus(e.target.value)}
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
        <Button variant="outlined" size="small" onClick={() => { setFilterValue(''); setFilterStatus(''); }}>
          Reset
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <TablePagination
          component="div"
          count={total}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          rowsPerPageOptions={[10, 25, 50, 100]}
          labelRowsPerPage="Einträge pro Seite:"
        />
        <Table className="list-table">
          <TableHead>
            <TableRow>
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'ID'}
                  direction={orderBy === 'ID' ? order : 'asc'}
                  onClick={() => handleSort('ID')}
                >
                  ID
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={orderBy === valueField}
                  direction={orderBy === valueField ? order : 'asc'}
                  onClick={() => handleSort(valueField)}
                >
                  {valueHeader}
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'Status'}
                  direction={orderBy === 'Status' ? order : 'asc'}
                  onClick={() => handleSort('Status')}
                >
                  Status
                </TableSortLabel>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedItems.length === 0 ? (
              <TableRow>
                <TableCell colSpan={3} align="center">No items found</TableCell>
              </TableRow>
            ) : (
              paginatedItems.map((item, idx) => (
                <TableRow key={item.ID || item.Address || item.UserAgent || item.Code || idx}>
                  <TableCell>{item.ID}</TableCell>
                  <TableCell>{item[valueField]}</TableCell>
                  <TableCell>{item.Status}</TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
        <TablePagination
          component="div"
          count={total}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          rowsPerPageOptions={[10, 25, 50, 100]}
          labelRowsPerPage="Einträge pro Seite:"
        />
      </TableContainer>
    </Box>
  );
};

export default ListView; 