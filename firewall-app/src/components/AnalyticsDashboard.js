import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Select from '@mui/material/Select';
import Button from '@mui/material/Button';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import InfiniteScrollCards from './InfiniteScrollCards';
import './styles.css';

const AnalyticsDashboard = () => {
    const [analytics, setAnalytics] = useState(null);
    const [period, setPeriod] = useState('24h');
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState('overview');

    // Traffic logs state
    const [trafficLogs, setTrafficLogs] = useState([]);
    const [logsLoading, setLogsLoading] = useState(false);
    const [logsError, setLogsError] = useState(null);
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
    const [logsInfiniteLoading, setLogsInfiniteLoading] = useState(false);
    const [logsHasMore, setLogsHasMore] = useState(true);

    useEffect(() => {
        console.log('useEffect triggered - period:', period);
        fetchAnalytics();
    }, [period]);

    useEffect(() => {
        if (activeTab === 'logs') {
            // Reset infinite scroll state when filters change
            setLogsHasMore(true);
            setTrafficLogs([]);
            fetchTrafficLogs();
            fetchLogsStats();
        }
    }, [activeTab, logsFilterIP, logsFilterEmail, logsFilterUserAgent, logsFilterUsername, logsFilterCountry, logsFilterASN, logsFilterResult, logsOrderBy, logsOrder]);

    const fetchAnalytics = async () => {
        try {
            console.log('Fetching analytics for period:', period);
            const response = await axios.get(`/api/analytics/traffic?period=${period}`);
            console.log('Analytics response:', response.data);
            setAnalytics(response.data);
            console.log('Analytics state set to:', response.data);
        } catch (error) {
            console.error('Error fetching analytics:', error);
        } finally {
            setLoading(false);
        }
    };



    const fetchTrafficLogs = async (isInfinite = false) => {
        if (isInfinite) {
            setLogsInfiniteLoading(true);
        } else {
            setLogsLoading(true);
        }
        setLogsError(null);
        
        try {
            const params = {
                page: isInfinite ? Math.floor(trafficLogs.length / 25) + 1 : 1,
                limit: 25,
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
                if (isInfinite) {
                    // Append to existing logs for infinite scroll
                    setTrafficLogs(prev => [...prev, ...response.data.logs]);
                    setLogsHasMore(response.data.logs.length === 25);
                } else {
                    // Replace logs for initial load
                    setTrafficLogs(response.data.logs);
                    setLogsHasMore(response.data.logs.length === 25);
                }
                setLogsTotal(response.data.total || 0);
            } else {
                if (!isInfinite) {
                    setTrafficLogs([]);
                }
                setLogsHasMore(false);
                setLogsTotal(0);
            }
        } catch (error) {
            console.error('Error fetching traffic logs:', error);
            setLogsError('Failed to fetch traffic logs');
        } finally {
            if (isInfinite) {
                setLogsInfiniteLoading(false);
            } else {
                setLogsLoading(false);
            }
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
        if (num === null || num === undefined) {
            console.log('formatNumber received null/undefined:', num);
            return '0';
        }
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



    const resetLogsFilters = () => {
        setLogsFilterIP('');
        setLogsFilterEmail('');
        setLogsFilterUserAgent('');
        setLogsFilterUsername('');
        setLogsFilterCountry('');
        setLogsFilterASN('');
        setLogsFilterResult('');
    };

    const handleLoadMoreLogs = () => {
        if (!logsInfiniteLoading && logsHasMore) {
            fetchTrafficLogs(true);
        }
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

            <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
                <Tabs 
                    value={activeTab} 
                    onChange={(event, newValue) => setActiveTab(newValue)}
                    aria-label="Analytics dashboard tabs"
                >
                    <Tab label="Overview" value="overview" />
                    <Tab label="Traffic Logs" value="logs" />
                </Tabs>
            </Box>

            {activeTab === 'overview' && analytics && (
                <div className="overview-tab">
                    {console.log('Rendering overview with analytics:', analytics)}
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
                        // Cards Layout with Infinite Scroll
                        <Box>
                            <InfiniteScrollCards
                                trafficLogs={trafficLogs}
                                loading={logsInfiniteLoading}
                                error={logsError}
                                total={logsTotal}
                                hasMore={logsHasMore}
                                onLoadMore={handleLoadMoreLogs}
                                onSort={handleLogsSort}
                                orderBy={logsOrderBy}
                                order={logsOrder}
                            />
                        </Box>
                    )}
                </div>
            )}
        </div>
    );
};

export default AnalyticsDashboard; 