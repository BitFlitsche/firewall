import React, { useState, useEffect, useRef, memo } from 'react';
import {
    Box, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
    TableSortLabel, IconButton, CircularProgress, Typography
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import CountryFlag from './CountryFlag';

const InfiniteScrollCountryTable = memo(({ 
    countries, 
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
    if (loading && countries.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <CircularProgress />
            </Box>
        );
    }

    // Error state
    if (error && countries.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography color="error">{error}</Typography>
            </Box>
        );
    }

    // Empty state
    if (countries.length === 0 && !loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography>No countries found</Typography>
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
                                        active={orderBy === 'code'}
                                        direction={orderBy === 'code' ? order : 'asc'}
                                        onClick={() => onSort('code')}
                                    >
                                        Flag
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'code'}
                                        direction={orderBy === 'code' ? order : 'asc'}
                                        onClick={() => onSort('code')}
                                    >
                                        Country Code
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'name'}
                                        direction={orderBy === 'name' ? order : 'asc'}
                                        onClick={() => onSort('name')}
                                    >
                                        Country Name
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
                            {countries.map((countryItem) => (
                                <TableRow key={countryItem.id}>
                                    <TableCell>
                                        <CountryFlag countryCode={countryItem.code} size={24} />
                                    </TableCell>
                                    <TableCell>{countryItem.code}</TableCell>
                                    <TableCell>{countryItem.name || 'Unknown Country'}</TableCell>
                                    <TableCell>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                            <Box
                                                sx={{
                                                    width: 8,
                                                    height: 8,
                                                    borderRadius: '50%',
                                                    backgroundColor: 
                                                        countryItem.status === 'allowed' ? '#28a745' :
                                                        countryItem.status === 'denied' ? '#dc3545' :
                                                        countryItem.status === 'whitelisted' ? '#007bff' : '#6c757d'
                                                }}
                                            />
                                            {countryItem.status}
                                        </Box>
                                    </TableCell>
                                    <TableCell>
                                        <IconButton onClick={() => onEdit(countryItem)} size="small">
                                            <EditIcon />
                                        </IconButton>
                                        <IconButton onClick={() => onDelete(countryItem.id)} size="small" color="error">
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
                    {!hasMore && countries.length > 0 && (
                        <Box sx={{ 
                            display: 'flex', 
                            justifyContent: 'center', 
                            alignItems: 'center', 
                            py: 3 
                        }}>
                            <Typography variant="body2" color="text.secondary">
                                Showing {countries.length} of {total} entries
                            </Typography>
                        </Box>
                    )}
                </TableContainer>
            </Paper>
        </Box>
    );
});

export default InfiniteScrollCountryTable; 