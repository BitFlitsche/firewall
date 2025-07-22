import React, { useState, useEffect, useRef, memo } from 'react';
import {
    Box, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
    TableSortLabel, IconButton, CircularProgress, Typography
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

const InfiniteScrollUserAgentTable = memo(({ 
    userAgents, 
    loading, 
    error, 
    total,
    hasMore,
    onLoadMore,
    onSort,
    orderBy,
    order,
    onEdit,
    onDelete
}) => {
    const observerRef = useRef();
    const loadingRef = useRef();

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

    // Loading state
    if (loading && userAgents.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <CircularProgress />
            </Box>
        );
    }

    // Error state
    if (error && userAgents.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography color="error">{error}</Typography>
            </Box>
        );
    }

    // Empty state
    if (userAgents.length === 0 && !loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography>No user agents found</Typography>
            </Box>
        );
    }

    return (
        <Box className="list-section" sx={{ mt: 4 }}>
            <Paper elevation={2} sx={{ p: 2 }}>
                <TableContainer>
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
                                        active={orderBy === 'user_agent'}
                                        direction={orderBy === 'user_agent' ? order : 'asc'}
                                        onClick={() => onSort('user_agent')}
                                    >
                                        User Agent
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
                                <TableCell>Is Regex</TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {userAgents.map(userAgentItem => (
                                <TableRow key={userAgentItem.id}>
                                    <TableCell>{userAgentItem.id}</TableCell>
                                    <TableCell>{userAgentItem.user_agent}</TableCell>
                                    <TableCell>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                            <Box
                                                sx={{
                                                    width: 8,
                                                    height: 8,
                                                    borderRadius: '50%',
                                                    backgroundColor: 
                                                        userAgentItem.status === 'allowed' ? '#28a745' :
                                                        userAgentItem.status === 'denied' ? '#dc3545' :
                                                        userAgentItem.status === 'whitelisted' ? '#007bff' : '#6c757d'
                                                }}
                                            />
                                            {userAgentItem.status}
                                        </Box>
                                    </TableCell>
                                    <TableCell>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                            <Box
                                                sx={{
                                                    width: 8,
                                                    height: 8,
                                                    borderRadius: '50%',
                                                    backgroundColor: userAgentItem.is_regex ? '#28a745' : '#6c757d'
                                                }}
                                            />
                                            {userAgentItem.is_regex ? 'Yes' : 'No'}
                                        </Box>
                                    </TableCell>
                                    <TableCell>
                                        <IconButton onClick={() => onEdit(userAgentItem)} size="small">
                                            <EditIcon />
                                        </IconButton>
                                        <IconButton onClick={() => onDelete(userAgentItem.id)} size="small" color="error">
                                            <DeleteIcon />
                                        </IconButton>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>

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
                    {!hasMore && userAgents.length > 0 && (
                        <Box sx={{ 
                            display: 'flex', 
                            justifyContent: 'center', 
                            alignItems: 'center', 
                            py: 3 
                        }}>
                            <Typography variant="body2" color="text.secondary">
                                Showing {userAgents.length} of {total} entries
                            </Typography>
                        </Box>
                    )}
                </TableContainer>
            </Paper>
        </Box>
    );
});

export default InfiniteScrollUserAgentTable; 