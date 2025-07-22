import React, { useState, useEffect, useRef, memo } from 'react';
import {
    Box, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
    TableSortLabel, IconButton, CircularProgress, Typography
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

const InfiniteScrollASNTable = memo(({ 
    asns, 
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
    if (loading && asns.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <CircularProgress />
            </Box>
        );
    }

    // Error state
    if (error && asns.length === 0) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography color="error">{error}</Typography>
            </Box>
        );
    }

    // Empty state
    if (asns.length === 0 && !loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Typography>No ASNs found</Typography>
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
                                        active={orderBy === 'asn'}
                                        direction={orderBy === 'asn' ? order : 'asc'}
                                        onClick={() => onSort('asn')}
                                    >
                                        ASN
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'rir'}
                                        direction={orderBy === 'rir' ? order : 'asc'}
                                        onClick={() => onSort('rir')}
                                    >
                                        RIR
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'domain'}
                                        direction={orderBy === 'domain' ? order : 'asc'}
                                        onClick={() => onSort('domain')}
                                    >
                                        Domain
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'cc'}
                                        direction={orderBy === 'cc' ? order : 'asc'}
                                        onClick={() => onSort('cc')}
                                    >
                                        Country
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'asname'}
                                        direction={orderBy === 'asname' ? order : 'asc'}
                                        onClick={() => onSort('asname')}
                                    >
                                        Name
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
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'source'}
                                        direction={orderBy === 'source' ? order : 'asc'}
                                        onClick={() => onSort('source')}
                                    >
                                        Source
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {asns.map((asn) => (
                                <TableRow key={asn.id}>
                                    <TableCell>{asn.id}</TableCell>
                                    <TableCell>{asn.asn}</TableCell>
                                    <TableCell>{asn.rir}</TableCell>
                                    <TableCell>{asn.domain}</TableCell>
                                    <TableCell>{asn.cc}</TableCell>
                                    <TableCell>{asn.asname}</TableCell>
                                    <TableCell>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                            <Box
                                                sx={{
                                                    width: 8,
                                                    height: 8,
                                                    borderRadius: '50%',
                                                    backgroundColor: 
                                                        asn.status === 'allowed' ? '#28a745' :
                                                        asn.status === 'denied' ? '#dc3545' :
                                                        asn.status === 'whitelisted' ? '#007bff' : '#6c757d'
                                                }}
                                            />
                                            {asn.status}
                                        </Box>
                                    </TableCell>
                                    <TableCell>{asn.source}</TableCell>
                                    <TableCell>
                                        <IconButton
                                            size="small"
                                            onClick={() => onEdit(asn)}
                                        >
                                            <EditIcon />
                                        </IconButton>
                                        <IconButton
                                            size="small"
                                            onClick={() => onDelete(asn.id)}
                                            color="error"
                                        >
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
                    {!hasMore && asns.length > 0 && (
                        <Box sx={{ 
                            display: 'flex', 
                            justifyContent: 'center', 
                            alignItems: 'center', 
                            py: 3 
                        }}>
                            <Typography variant="body2" color="text.secondary">
                                Showing {asns.length} of {total} entries
                            </Typography>
                        </Box>
                    )}
                </TableContainer>
            </Paper>
        </Box>
    );
});

export default InfiniteScrollASNTable; 