import React, { useState, useEffect } from 'react';
import axios from 'axios';
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
import './styles.css';

const AnalyticsDashboard = () => {
    const [analytics, setAnalytics] = useState(null);
    const [relationships, setRelationships] = useState([]);
    const [period, setPeriod] = useState('24h');
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState('overview');

    // Traffic logs table state
    const [trafficLogs, setTrafficLogs] = useState([]);
    const [logsLoading, setLogsLoading] = useState(false);
    const [logsError, setLogsError] = useState(null);
    const [logsPage, setLogsPage] = useState(0);
    const [logsRowsPerPage, setLogsRowsPerPage] = useState(25);
    const [logsTotal, setLogsTotal] = useState(0);
    const [logsFilterIP, setLogsFilterIP] = useState('');
    const [logsFilterEmail, setLogsFilterEmail] = useState('');
    const [logsFilterUserAgent, setLogsFilterUserAgent] = useState('');
    const [logsFilterUsername, setLogsFilterUsername] = useState('');
    const [logsFilterCountry, setLogsFilterCountry] = useState('');
    const [logsFilterASN, setLogsFilterASN] = useState('');
    const [logsFilterResult, setLogsFilterResult] = useState('');
    const [logsOrderBy, setLogsOrderBy] = useState('timestamp');
    const [logsOrder, setLogsOrder] = useState('desc');
    const [logsStats, setLogsStats] = useState({ total: 0, allowed: 0, denied: 0, whitelisted: 0 });

    useEffect(() => {
        fetchAnalytics();
        fetchRelationships();
    }, [period]);

    useEffect(() => {
        if (activeTab === 'logs') {
            fetchTrafficLogs();
            fetchLogsStats();
        }
    }, [activeTab, logsPage, logsRowsPerPage, logsFilterIP, logsFilterEmail, logsFilterUserAgent, logsFilterUsername, logsFilterCountry, logsFilterASN, logsFilterResult, logsOrderBy, logsOrder]);

    const fetchAnalytics = async () => {
        try {
            const response = await axios.get(`/api/analytics/traffic?period=${period}`);
            setAnalytics(response.data);
        } catch (error) {
            console.error('Error fetching analytics:', error);
        } finally {
            setLoading(false);
        }
    };

    const fetchRelationships = async () => {
        try {
            const response = await axios.get('/api/analytics/relationships?limit=20');
            setRelationships(response.data.relationships);
        } catch (error) {
            console.error('Error fetching relationships:', error);
        }
    };

    const fetchTrafficLogs = async () => {
        setLogsLoading(true);
        setLogsError(null);
        try {
            const params = {
                page: logsPage + 1,
                limit: logsRowsPerPage,
                orderBy: logsOrderBy,
                order: logsOrder,
            };

            if (logsFilterIP) {
                params.ip_address = logsFilterIP;
            }
            if (logsFilterEmail) {
                params.email = logsFilterEmail;
            }
            if (logsFilterUserAgent) {
                params.user_agent = logsFilterUserAgent;
            }
            if (logsFilterUsername) {
                params.username = logsFilterUsername;
            }
            if (logsFilterCountry) {
                params.country = logsFilterCountry;
            }
            if (logsFilterASN) {
                params.asn = logsFilterASN;
            }
            if (logsFilterResult) {
                params.final_result = logsFilterResult;
            }

            const response = await axios.get('/api/analytics/logs', { params });
            
            if (response.data && response.data.logs) {
                setTrafficLogs(response.data.logs);
                setLogsTotal(response.data.total || 0);
            } else {
                setTrafficLogs([]);
                setLogsTotal(0);
            }
        } catch (error) {
            console.error('Error fetching traffic logs:', error);
            setLogsError('Failed to fetch traffic logs');
        } finally {
            setLogsLoading(false);
        }
    };

    const fetchLogsStats = async () => {
        try {
            const response = await axios.get('/api/analytics/logs/stats');
            setLogsStats(response.data);
        } catch (error) {
            console.error('Error fetching traffic log stats:', error);
        }
    };

    const formatNumber = (num) => {
        return new Intl.NumberFormat().format(num);
    };

    const formatPercentage = (num) => {
        return num.toFixed(1) + '%';
    };

    const formatTime = (ms) => {
        return ms.toFixed(2) + 'ms';
    };

    const handleLogsSort = (field) => {
        if (logsOrderBy === field) {
            setLogsOrder(logsOrder === 'asc' ? 'desc' : 'asc');
        } else {
            setLogsOrderBy(field);
            setLogsOrder('asc');
        }
    };

    const handleLogsChangePage = (event, newPage) => {
        setLogsPage(newPage);
    };

    const handleLogsChangeRowsPerPage = (event) => {
        setLogsRowsPerPage(parseInt(event.target.value, 10));
        setLogsPage(0);
    };

    const resetLogsFilters = () => {
        setLogsFilterIP('');
        setLogsFilterEmail('');
        setLogsFilterUserAgent('');
        setLogsFilterUsername('');
        setLogsFilterCountry('');
        setLogsFilterASN('');
        setLogsFilterResult('');
    };

    if (loading) {
        return (
            <div className="analytics-dashboard">
                <div className="loading">Loading analytics...</div>
            </div>
        );
    }

    return (
        <div className="analytics-dashboard">
            <div className="dashboard-header">
                <h1>Traffic Analytics Dashboard</h1>
                <div className="controls">
                    <select value={period} onChange={(e) => setPeriod(e.target.value)}>
                        <option value="1h">Last Hour</option>
                        <option value="24h">Last 24 Hours</option>
                        <option value="7d">Last 7 Days</option>
                        <option value="30d">Last 30 Days</option>
                    </select>
                </div>
            </div>

            <div className="tab-navigation">
                <button 
                    className={activeTab === 'overview' ? 'active' : ''} 
                    onClick={() => setActiveTab('overview')}
                >
                    Overview
                </button>
                <button 
                    className={activeTab === 'relationships' ? 'active' : ''} 
                    onClick={() => setActiveTab('relationships')}
                >
                    Data Relationships
                </button>
                <button 
                    className={activeTab === 'logs' ? 'active' : ''} 
                    onClick={() => setActiveTab('logs')}
                >
                    Traffic Logs
                </button>
            </div>

            {activeTab === 'overview' && analytics && (
                <div className="overview-tab">
                    <div className="metrics-grid">
                        <div className="metric-card">
                            <h3>Total Requests</h3>
                            <div className="metric-value">{formatNumber(analytics.total_requests)}</div>
                        </div>
                        <div className="metric-card">
                            <h3>Allowed</h3>
                            <div className="metric-value allowed">{formatNumber(analytics.allowed_requests)}</div>
                        </div>
                        <div className="metric-card">
                            <h3>Denied</h3>
                            <div className="metric-value denied">{formatNumber(analytics.denied_requests)}</div>
                        </div>
                        <div className="metric-card">
                            <h3>Whitelisted</h3>
                            <div className="metric-value whitelisted">{formatNumber(analytics.whitelisted_requests)}</div>
                        </div>
                        <div className="metric-card">
                            <h3>Avg Response Time</h3>
                            <div className="metric-value">{formatTime(analytics.avg_response_time_ms)}</div>
                        </div>
                        <div className="metric-card">
                            <h3>Cache Hit Rate</h3>
                            <div className="metric-value">{formatPercentage(analytics.cache_hit_rate)}</div>
                        </div>
                    </div>

                    <div className="charts-section">
                        <div className="chart-container">
                            <h3>Request Results Distribution</h3>
                            <div className="pie-chart">
                                <div className="pie-segment allowed" style={{
                                    '--percentage': analytics.total_requests > 0 ? 
                                        (analytics.allowed_requests / analytics.total_requests) * 100 : 0
                                }}>
                                    <span>Allowed: {formatNumber(analytics.allowed_requests)}</span>
                                </div>
                                <div className="pie-segment denied" style={{
                                    '--percentage': analytics.total_requests > 0 ? 
                                        (analytics.denied_requests / analytics.total_requests) * 100 : 0
                                }}>
                                    <span>Denied: {formatNumber(analytics.denied_requests)}</span>
                                </div>
                                <div className="pie-segment whitelisted" style={{
                                    '--percentage': analytics.total_requests > 0 ? 
                                        (analytics.whitelisted_requests / analytics.total_requests) * 100 : 0
                                }}>
                                    <span>Whitelisted: {formatNumber(analytics.whitelisted_requests)}</span>
                                </div>
                            </div>
                        </div>

                        <div className="chart-container">
                            <h3>Performance Metrics</h3>
                            <div className="performance-metrics">
                                <div className="metric-bar">
                                    <span>Response Time</span>
                                    <div className="bar">
                                        <div 
                                            className="bar-fill" 
                                            style={{width: `${Math.min(analytics.avg_response_time_ms / 100, 100)}%`}}
                                        ></div>
                                    </div>
                                    <span>{formatTime(analytics.avg_response_time_ms)}</span>
                                </div>
                                <div className="metric-bar">
                                    <span>Cache Hit Rate</span>
                                    <div className="bar">
                                        <div 
                                            className="bar-fill cache" 
                                            style={{width: `${analytics.cache_hit_rate}%`}}
                                        ></div>
                                    </div>
                                    <span>{formatPercentage(analytics.cache_hit_rate)}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {activeTab === 'relationships' && (
                <div className="relationships-tab">
                    <h3>Top Data Relationships</h3>
                    <div className="relationships-list">
                        {relationships.map((rel, index) => (
                            <div key={index} className="relationship-item">
                                <div className="relationship-header">
                                    <div className="relationship-type">{rel.relationship_type.replace('_', ' ')}</div>
                                    <div className="relationship-frequency">{rel.frequency} occurrences</div>
                                </div>
                                <div className="relationship-data">
                                    {rel.ip_address && <span className="data-item">IP: {rel.ip_address}</span>}
                                    {rel.email && <span className="data-item">Email: {rel.email}</span>}
                                    {rel.user_agent && <span className="data-item">User Agent: {rel.user_agent.substring(0, 50)}...</span>}
                                    {rel.username && <span className="data-item">Username: {rel.username}</span>}
                                    {rel.country && <span className="data-item">Country: {rel.country}</span>}
                                    {rel.charset && <span className="data-item">Charset: {rel.charset}</span>}
                                </div>
                                <div className="relationship-timeline">
                                    <span>First seen: {new Date(rel.first_seen).toLocaleString()}</span>
                                    <span>Last seen: {new Date(rel.last_seen).toLocaleString()}</span>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {activeTab === 'logs' && (
                <div className="logs-tab">
                    <h3>Traffic Logs</h3>
                    
                    <Box sx={{ display: 'flex', gap: 2, mb: 2, flexWrap: 'wrap', alignItems: 'center' }}>
                        <TextField
                            label="IP Address Filter"
                            value={logsFilterIP}
                            onChange={(e) => setLogsFilterIP(e.target.value)}
                            size="small"
                            sx={{ minWidth: 200 }}
                        />
                        <TextField
                            label="Email Filter"
                            value={logsFilterEmail}
                            onChange={(e) => setLogsFilterEmail(e.target.value)}
                            size="small"
                            sx={{ minWidth: 200 }}
                        />
                        <TextField
                            label="User Agent Filter"
                            value={logsFilterUserAgent}
                            onChange={(e) => setLogsFilterUserAgent(e.target.value)}
                            size="small"
                            sx={{ minWidth: 200 }}
                        />
                        <TextField
                            label="Username Filter"
                            value={logsFilterUsername}
                            onChange={(e) => setLogsFilterUsername(e.target.value)}
                            size="small"
                            sx={{ minWidth: 150 }}
                        />
                        <TextField
                            label="Country Filter"
                            value={logsFilterCountry}
                            onChange={(e) => setLogsFilterCountry(e.target.value)}
                            size="small"
                            sx={{ minWidth: 120 }}
                        />
                        <TextField
                            label="ASN Filter"
                            value={logsFilterASN}
                            onChange={(e) => setLogsFilterASN(e.target.value)}
                            size="small"
                            sx={{ minWidth: 120 }}
                        />
                        <FormControl size="small" sx={{ minWidth: 140 }}>
                            <InputLabel shrink>Result</InputLabel>
                            <Select
                                value={logsFilterResult}
                                label="Result"
                                onChange={(e) => setLogsFilterResult(e.target.value)}
                                displayEmpty
                                renderValue={(selected) => {
                                    if (!selected) return `All (${logsStats.total})`;
                                    if (selected === 'allowed') return `Allowed (${logsStats.allowed})`;
                                    if (selected === 'denied') return `Denied (${logsStats.denied})`;
                                    if (selected === 'whitelisted') return `Whitelisted (${logsStats.whitelisted})`;
                                    return selected;
                                }}
                            >
                                <MenuItem value="">All ({logsStats.total})</MenuItem>
                                <MenuItem value="allowed">Allowed ({logsStats.allowed})</MenuItem>
                                <MenuItem value="denied">Denied ({logsStats.denied})</MenuItem>
                                <MenuItem value="whitelisted">Whitelisted ({logsStats.whitelisted})</MenuItem>
                            </Select>
                        </FormControl>
                        <Button variant="outlined" size="small" onClick={resetLogsFilters}>
                            Reset
                        </Button>
                    </Box>

                    {logsLoading ? (
                        <div>Loading traffic logs...</div>
                    ) : logsError ? (
                        <div className="error">{logsError}</div>
                    ) : (
                        <TableContainer component={Paper}>
                            <TablePagination
                                component="div"
                                count={logsTotal}
                                page={logsPage}
                                onPageChange={handleLogsChangePage}
                                rowsPerPage={logsRowsPerPage}
                                onRowsPerPageChange={handleLogsChangeRowsPerPage}
                                rowsPerPageOptions={[10, 25, 50, 100]}
                                labelRowsPerPage="Entries per page:"
                            />
                            <Table className="list-table">
                                <TableHead>
                                    <TableRow>
                                        <TableCell sx={{ width: '140px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'timestamp'}
                                                direction={logsOrderBy === 'timestamp' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('timestamp')}
                                            >
                                                Timestamp
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '100px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'final_result'}
                                                direction={logsOrderBy === 'final_result' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('final_result')}
                                            >
                                                Result
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '140px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'ip_address'}
                                                direction={logsOrderBy === 'ip_address' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('ip_address')}
                                            >
                                                IP Address
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '200px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'email'}
                                                direction={logsOrderBy === 'email' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('email')}
                                            >
                                                Email
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '250px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'user_agent'}
                                                direction={logsOrderBy === 'user_agent' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('user_agent')}
                                            >
                                                User Agent
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '120px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'username'}
                                                direction={logsOrderBy === 'username' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('username')}
                                            >
                                                Username
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '80px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'country'}
                                                direction={logsOrderBy === 'country' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('country')}
                                            >
                                                Country
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '100px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'asn'}
                                                direction={logsOrderBy === 'asn' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('asn')}
                                            >
                                                ASN
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '120px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'response_time_ms'}
                                                direction={logsOrderBy === 'response_time_ms' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('response_time_ms')}
                                            >
                                                Response Time
                                            </TableSortLabel>
                                        </TableCell>
                                        <TableCell sx={{ width: '100px' }}>
                                            <TableSortLabel
                                                active={logsOrderBy === 'cache_hit'}
                                                direction={logsOrderBy === 'cache_hit' ? logsOrder : 'asc'}
                                                onClick={() => handleLogsSort('cache_hit')}
                                            >
                                                Cache Hit
                                            </TableSortLabel>
                                        </TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {trafficLogs.length === 0 ? (
                                        <TableRow>
                                            <TableCell colSpan={10} align="center">No traffic logs found</TableCell>
                                        </TableRow>
                                    ) : (
                                        trafficLogs.map((log, index) => (
                                            <TableRow key={log.id || index}>
                                                <TableCell sx={{ width: '140px' }}>{new Date(log.timestamp).toLocaleString()}</TableCell>
                                                <TableCell sx={{ width: '100px' }}>
                                                    <span className={`result ${log.final_result}`}>
                                                        {log.final_result}
                                                    </span>
                                                </TableCell>
                                                <TableCell sx={{ width: '140px' }}>{log.ip_address || '-'}</TableCell>
                                                <TableCell sx={{ width: '200px' }}>{log.email || '-'}</TableCell>
                                                <TableCell sx={{ width: '250px' }}>
                                                    {log.user_agent ? 
                                                        (log.user_agent.length > 50 ? 
                                                            log.user_agent.substring(0, 50) + '...' : 
                                                            log.user_agent) : 
                                                        '-'
                                                    }
                                                </TableCell>
                                                <TableCell sx={{ width: '120px' }}>{log.username || '-'}</TableCell>
                                                <TableCell sx={{ width: '80px' }}>{log.country || '-'}</TableCell>
                                                <TableCell sx={{ width: '100px' }}>{log.asn || '-'}</TableCell>
                                                <TableCell sx={{ width: '120px' }}>{log.response_time_ms}ms</TableCell>
                                                <TableCell sx={{ width: '100px' }}>
                                                    {log.cache_hit ? 
                                                        <span className="cache-hit">Yes</span> : 
                                                        <span className="cache-miss">No</span>
                                                    }
                                                </TableCell>
                                            </TableRow>
                                        ))
                                    )}
                                </TableBody>
                            </Table>
                            <TablePagination
                                component="div"
                                count={logsTotal}
                                page={logsPage}
                                onPageChange={handleLogsChangePage}
                                rowsPerPage={logsRowsPerPage}
                                onRowsPerPageChange={handleLogsChangeRowsPerPage}
                                rowsPerPageOptions={[10, 25, 50, 100]}
                                labelRowsPerPage="Entries per page:"
                            />
                        </TableContainer>
                    )}
                </div>
            )}
        </div>
    );
};

export default AnalyticsDashboard; 