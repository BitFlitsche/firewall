import React, { useState, useEffect, useRef, memo } from 'react';
import {
    Box, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
    TableSortLabel, IconButton, CircularProgress, Typography
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

const InfiniteScrollEmailTable = memo(({ 
    emails, 
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
    if (loading && emails.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <CircularProgress />
            </Box>
        );
    }

    // Error state
    if (error && emails.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography color="error">{error}</Typography>
            </Box>
        );
    }

    // Empty state
    if (emails.length === 0 && !loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography>No emails found</Typography>
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
                                        active={orderBy === 'address'}
                                        direction={orderBy === 'address' ? order : 'asc'}
                                        onClick={() => onSort('address')}
                                    >
                                        Email Address
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
                                <TableCell>Regex</TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {emails.map(emailItem => (
                                <TableRow key={emailItem.id}>
                                    <TableCell>{emailItem.id}</TableCell>
                                    <TableCell>{emailItem.address}</TableCell>
                                    <TableCell>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                            <Box
                                                sx={{
                                                    width: 8,
                                                    height: 8,
                                                    borderRadius: '50%',
                                                    backgroundColor: 
                                                        emailItem.status === 'allowed' ? '#28a745' :
                                                        emailItem.status === 'denied' ? '#dc3545' :
                                                        emailItem.status === 'whitelisted' ? '#007bff' : '#6c757d'
                                                }}
                                            />
                                            {emailItem.status}
                                        </Box>
                                    </TableCell>
                                    <TableCell>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                            <Box
                                                sx={{
                                                    width: 8,
                                                    height: 8,
                                                    borderRadius: '50%',
                                                    backgroundColor: emailItem.is_regex ? '#28a745' : '#6c757d'
                                                }}
                                            />
                                            {emailItem.is_regex ? 'Yes' : 'No'}
                                        </Box>
                                    </TableCell>
                                    <TableCell>
                                        <IconButton onClick={() => onEdit(emailItem)} size="small">
                                            <EditIcon />
                                        </IconButton>
                                        <IconButton onClick={() => onDelete(emailItem.id)} size="small" color="error">
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
                    {!hasMore && emails.length > 0 && (
                        <Box sx={{ 
                            display: 'flex', 
                            justifyContent: 'center', 
                            alignItems: 'center', 
                            py: 3 
                        }}>
                            <Typography variant="body2" color="text.secondary">
                                Showing {emails.length} of {total} entries
                            </Typography>
                        </Box>
                    )}
                </TableContainer>
            </Paper>
        </Box>
    );
});

export default InfiniteScrollEmailTable; 