import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './styles.css';

const AnalyticsDashboard = () => {
    const [analytics, setAnalytics] = useState(null);
    const [relationships, setRelationships] = useState([]);
    const [period, setPeriod] = useState('24h');
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState('overview');

    useEffect(() => {
        fetchAnalytics();
        fetchRelationships();
    }, [period]);

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

    const formatNumber = (num) => {
        return new Intl.NumberFormat().format(num);
    };

    const formatPercentage = (num) => {
        return num.toFixed(1) + '%';
    };

    const formatTime = (ms) => {
        return ms.toFixed(2) + 'ms';
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
                    <h3>Recent Traffic Logs</h3>
                    <div className="logs-list">
                        {analytics?.logs?.map((log, index) => (
                            <div key={index} className="log-item">
                                <div className="log-header">
                                    <span className={`result ${log.final_result}`}>{log.final_result}</span>
                                    <span className="timestamp">{new Date(log.timestamp).toLocaleString()}</span>
                                    <span className="response-time">{log.response_time_ms}ms</span>
                                    {log.cache_hit && <span className="cache-hit">Cache Hit</span>}
                                </div>
                                <div className="log-data">
                                    {log.ip_address && <span>IP: {log.ip_address}</span>}
                                    {log.email && <span>Email: {log.email}</span>}
                                    {log.user_agent && <span>User Agent: {log.user_agent.substring(0, 50)}...</span>}
                                    {log.username && <span>Username: {log.username}</span>}
                                    {log.country && <span>Country: {log.country}</span>}
                                    {log.asn && <span>ASN: {log.asn}</span>}
                                    {log.charset && <span>Charset: {log.charset}</span>}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};

export default AnalyticsDashboard; 