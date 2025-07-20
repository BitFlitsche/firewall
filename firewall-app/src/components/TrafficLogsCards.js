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
import './styles.css';

const TrafficLogsCards = ({ 
    trafficLogs, 
    loading, 
    error, 
    total, 
    page, 
    rowsPerPage, 
    onPageChange, 
    onRowsPerPageChange,
    onSort,
    orderBy,
    order
}) => {
    const [expandedCards, setExpandedCards] = useState(new Set());

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

    // Render card content based on detail level
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
                            <CountryFlag countryCode={log.country} size={20} />
                            <Typography variant="body2">
                                {log.country ? `${log.country} - ${getCountryName(log.country)}` : 'N/A'}
                            </Typography>
                        </Box>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <BusinessIcon fontSize="small" color="action" />
                            <Typography variant="body2">
                                {log.asn || 'N/A'}
                            </Typography>
                        </Box>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <SpeedIcon fontSize="small" color="action" />
                            <Typography variant="body2">
                                {log.response_time_ms}ms
                            </Typography>
                        </Box>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <CachedIcon fontSize="small" color="action" />
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
                                        <EmailIcon fontSize="small" color="action" />
                                        <Typography variant="body2" sx={{ wordBreak: 'break-all' }}>
                                            {log.email}
                                        </Typography>
                                    </Box>
                                </Grid>
                            )}
                            {log.username && (
                                <Grid item xs={12} sm={6}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                        <PersonIcon fontSize="small" color="action" />
                                        <Typography variant="body2">
                                            {log.username}
                                        </Typography>
                                    </Box>
                                </Grid>
                            )}
                            {log.user_agent && (
                                <Grid item xs={12}>
                                    <Box sx={{ display: 'flex', alignItems: 'flex-start', gap: 1 }}>
                                        <ComputerIcon fontSize="small" color="action" sx={{ mt: 0.5 }} />
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

    // Lazy loading with intersection observer
    const observerRef = useRef();
    const [visibleCards, setVisibleCards] = useState(50); // Start with 50 cards

    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        setVisibleCards(prev => Math.min(prev + 25, trafficLogs.length));
                    }
                });
            },
            { threshold: 0.1 }
        );
        observerRef.current = observer;
        return () => observer.disconnect();
    }, [trafficLogs.length]);

    // Render cards with lazy loading
    const renderCards = () => {
        return trafficLogs.slice(0, visibleCards).map((log, index) => {
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
                    ref={index === visibleCards - 1 ? (el) => {
                        if (el && observerRef.current) {
                            observerRef.current.observe(el);
                        }
                    } : null}
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
        });
    };

    // Loading state
    if (loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography>Loading traffic logs...</Typography>
            </Box>
        );
    }

    // Error state
    if (error) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography color="error">{error}</Typography>
            </Box>
        );
    }

    // Empty state
    if (trafficLogs.length === 0) {
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
                {renderCards()}
                {visibleCards < trafficLogs.length && (
                    <Box sx={{ display: 'flex', justifyContent: 'center', py: 2 }}>
                        <Typography variant="body2" color="text.secondary">
                            Loading more entries...
                        </Typography>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

export default TrafficLogsCards; 