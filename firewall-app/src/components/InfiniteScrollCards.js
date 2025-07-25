import React, { useState, useCallback, useRef, useEffect } from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardActions from '@mui/material/CardActions';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import Chip from '@mui/material/Chip';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Collapse from '@mui/material/Collapse';
import CircularProgress from '@mui/material/CircularProgress';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import EmailIcon from '@mui/icons-material/Email';
import ComputerIcon from '@mui/icons-material/Computer';
import PersonIcon from '@mui/icons-material/Person';
import PublicIcon from '@mui/icons-material/Public';
import BusinessIcon from '@mui/icons-material/Business';
import SpeedIcon from '@mui/icons-material/Speed';
import CachedIcon from '@mui/icons-material/Cached';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';
import WarningIcon from '@mui/icons-material/Warning';
import CountryFlag from './CountryFlag';
import { getCountryName } from '../utils/country_codes';
import Tooltip from '@mui/material/Tooltip';
import './styles.css';

const InfiniteScrollCards = ({ 
    trafficLogs, 
    loading, 
    error, 
    total,
    hasMore,
    onLoadMore,
    onSort,
    orderBy,
    order
}) => {
    const [expandedCards, setExpandedCards] = useState(new Set());
    const observerRef = useRef();
    const loadingRef = useRef();

    // Reset expanded cards when filters change
    useEffect(() => {
        setExpandedCards(new Set());
    }, [orderBy, order]);

    // Infinite scroll observer
    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting && hasMore && !loading) {
                        onLoadMore();
                    }
                });
            },
            { threshold: 0.1 }
        );
        observerRef.current = observer;

        if (loadingRef.current) {
            observer.observe(loadingRef.current);
        }

        return () => observer.disconnect();
    }, [hasMore, loading, onLoadMore]);

    // Toggle card expansion
    const toggleCard = useCallback((logId) => {
        setExpandedCards(prev => {
            const newSet = new Set(prev);
            if (newSet.has(logId)) {
                newSet.delete(logId);
            } else {
                newSet.add(logId);
            }
            return newSet;
        });
    }, []);

    // Get result icon and color
    const getResultIcon = (result) => {
        switch (result) {
            case 'allowed':
                return { icon: <CheckCircleIcon />, color: 'success' };
            case 'denied':
                return { icon: <CancelIcon />, color: 'error' };
            case 'whitelisted':
                return { icon: <CheckCircleIcon />, color: 'success' };
            default:
                return { icon: <WarningIcon />, color: 'warning' };
        }
    };

    // Format timestamp
    const formatTimestamp = (timestamp) => {
        const date = new Date(timestamp);
        return date.toLocaleString();
    };

    // Check if card has additional details to show
    const hasAdditionalDetails = (log) => {
        return !!(log.email || log.username || log.user_agent);
    };

    // Render card content
    const renderCardContent = (log, isExpanded) => {
        const resultInfo = getResultIcon(log.final_result);

        return (
            <CardContent sx={{ p: 2 }}>
                {/* Compact Header */}
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Typography variant="h6" component="div" sx={{ fontWeight: 'bold' }}>
                            {log.ip_address || 'N/A'}
                        </Typography>
                        <Chip 
                            icon={resultInfo.icon}
                            label={log.final_result}
                            color={resultInfo.color}
                            size="small"
                        />
                    </Box>
                    <Typography variant="caption" color="text.secondary">
                        {formatTimestamp(log.timestamp)}
                    </Typography>
                </Box>

                {/* Basic Info Row */}
                <Grid container spacing={1} sx={{ mb: 1 }}>
                    <Grid item xs={6} sm={3}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <Tooltip title="Country">
                                <CountryFlag countryCode={log.country} size={20} />
                            </Tooltip>
                            <Typography variant="body2">
                                {log.country ? `${log.country} - ${getCountryName(log.country)}` : 'N/A'}
                            </Typography>
                        </Box>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <Tooltip title="Autonomous System Number">
                                <BusinessIcon fontSize="small" color="action" />
                            </Tooltip>
                            <Typography variant="body2">
                                {log.asn || 'N/A'}
                            </Typography>
                        </Box>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <Tooltip title="Response Time">
                                <SpeedIcon fontSize="small" color="action" />
                            </Tooltip>
                            <Typography variant="body2">
                                {log.response_time_ms}ms
                            </Typography>
                        </Box>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <Tooltip title="Cache Hit">
                                <CachedIcon fontSize="small" color="action" />
                            </Tooltip>
                            <Typography variant="body2">
                                {log.cache_hit ? 'Yes' : 'No'}
                            </Typography>
                        </Box>
                    </Grid>
                </Grid>

                {/* Expanded Details */}
                <Collapse in={isExpanded} timeout="auto" unmountOnExit>
                    <Box sx={{ mt: 2, pt: 2, borderTop: 1, borderColor: 'divider' }}>
                        <Grid container spacing={2}>
                            {log.email && (
                                <Grid item xs={12} sm={6}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                        <Tooltip title="Email Address">
                                            <EmailIcon fontSize="small" color="action" />
                                        </Tooltip>
                                        <Typography variant="body2" sx={{ wordBreak: 'break-all' }}>
                                            {log.email}
                                        </Typography>
                                    </Box>
                                </Grid>
                            )}
                            {log.username && (
                                <Grid item xs={12} sm={6}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                        <Tooltip title="Username">
                                            <PersonIcon fontSize="small" color="action" />
                                        </Tooltip>
                                        <Typography variant="body2">
                                            {log.username}
                                        </Typography>
                                    </Box>
                                </Grid>
                            )}
                            {log.user_agent && (
                                <Grid item xs={12}>
                                    <Box sx={{ display: 'flex', alignItems: 'flex-start', gap: 1 }}>
                                        <Tooltip title="User Agent">
                                            <ComputerIcon fontSize="small" color="action" sx={{ mt: 0.5 }} />
                                        </Tooltip>
                                        <Typography variant="body2" sx={{ wordBreak: 'break-all' }}>
                                            {log.user_agent}
                                        </Typography>
                                    </Box>
                                </Grid>
                            )}
                        </Grid>
                    </Box>
                </Collapse>
            </CardContent>
        );
    };

    // Loading state
    if (loading && trafficLogs.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <CircularProgress />
            </Box>
        );
    }

    // Error state
    if (error && trafficLogs.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography color="error">{error}</Typography>
            </Box>
        );
    }

    // Empty state
    if (trafficLogs.length === 0 && !loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography>No traffic logs found</Typography>
            </Box>
        );
    }

    return (
        <Box sx={{ 
            height: 'calc(100vh - 400px)', 
            minHeight: 400,
            overflowY: 'auto',
            overflowX: 'hidden'
        }}>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                {/* Render all loaded cards */}
                {trafficLogs.map((log, index) => {
                    const isExpanded = expandedCards.has(log.id || index);
                    const hasDetails = hasAdditionalDetails(log);

                    return (
                        <Card 
                            key={log.id || index}
                            sx={{ 
                                m: 1, 
                                cursor: hasDetails ? 'pointer' : 'default',
                                '&:hover': hasDetails ? { 
                                    boxShadow: 3,
                                    transform: 'translateY(-1px)',
                                    transition: 'all 0.2s ease-in-out'
                                } : {}
                            }}
                            onClick={hasDetails ? () => toggleCard(log.id || index) : undefined}
                        >
                            {renderCardContent(log, isExpanded)}
                            {hasDetails && (
                                <CardActions sx={{ justifyContent: 'center', py: 1 }}>
                                    <IconButton size="small">
                                        {isExpanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                                    </IconButton>
                                </CardActions>
                            )}
                        </Card>
                    );
                })}

                {/* Loading indicator at bottom */}
                {(loading || hasMore) && (
                    <Box 
                        ref={loadingRef}
                        sx={{ 
                            display: 'flex', 
                            justifyContent: 'center', 
                            alignItems: 'center', 
                            py: 3 
                        }}
                    >
                        {loading ? (
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                <CircularProgress size={20} />
                                <Typography variant="body2" color="text.secondary">
                                    Loading more entries...
                                </Typography>
                            </Box>
                        ) : (
                            <Typography variant="body2" color="text.secondary">
                                Scroll to load more
                            </Typography>
                        )}
                    </Box>
                )}

                {/* End of results */}
                {!hasMore && trafficLogs.length > 0 && (
                    <Box sx={{ 
                        display: 'flex', 
                        justifyContent: 'center', 
                        alignItems: 'center', 
                        py: 3 
                    }}>
                        <Typography variant="body2" color="text.secondary">
                            Showing {trafficLogs.length} of {total} entries
                        </Typography>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

export default InfiniteScrollCards; 