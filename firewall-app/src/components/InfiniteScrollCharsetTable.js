import React, { useState, useEffect, useRef, memo } from 'react';
import {
    Box, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
    TableSortLabel, IconButton, CircularProgress, Typography
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

const InfiniteScrollCharsetTable = memo(({ 
    charsets, 
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
    if (loading && charsets.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <CircularProgress />
            </Box>
        );
    }

    // Error state
    if (error && charsets.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography color="error">{error}</Typography>
            </Box>
        );
    }

    // Empty state
    if (charsets.length === 0 && !loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography>No charsets found</Typography>
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
                            {charsets.map(charsetItem => (
                                <TableRow key={charsetItem.id}>
                                    <TableCell>{charsetItem.charset}</TableCell>
                                    <TableCell>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                            <Box
                                                sx={{
                                                    width: 8,
                                                    height: 8,
                                                    borderRadius: '50%',
                                                    backgroundColor: 
                                                        charsetItem.status === 'allowed' ? '#28a745' :
                                                        charsetItem.status === 'denied' ? '#dc3545' :
                                                        charsetItem.status === 'whitelisted' ? '#007bff' : '#6c757d'
                                                }}
                                            />
                                            {charsetItem.status}
                                        </Box>
                                    </TableCell>
                                    <TableCell>
                                        <IconButton onClick={() => onEdit(charsetItem)} size="small">
                                            <EditIcon />
                                        </IconButton>
                                        <IconButton onClick={() => onDelete(charsetItem.id)} size="small" color="error">
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
                    {!hasMore && charsets.length > 0 && (
                        <Box sx={{ 
                            display: 'flex', 
                            justifyContent: 'center', 
                            alignItems: 'center', 
                            py: 3 
                        }}>
                            <Typography variant="body2" color="text.secondary">
                                Showing {charsets.length} of {total} entries
                            </Typography>
                        </Box>
                    )}
                </TableContainer>
            </Paper>
        </Box>
    );
});

export default InfiniteScrollCharsetTable; 