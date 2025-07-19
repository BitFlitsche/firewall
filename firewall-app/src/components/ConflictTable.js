import React from 'react';
import {
    Box,
    Paper,
    Typography,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Chip,
    Alert,
    Collapse,
    Button,
    Stack
} from '@mui/material';
import ErrorIcon from '@mui/icons-material/Error';
import WarningIcon from '@mui/icons-material/Warning';
import DeleteSweepIcon from '@mui/icons-material/DeleteSweep';

const ConflictTable = ({ conflicts, open = true, onDeleteAllConflicts, isDeleting = false }) => {
    if (!conflicts || conflicts.length === 0) {
        return null;
    }

    const getSeverityColor = (severity) => {
        switch (severity) {
            case 'error':
                return '#d32f2f'; // Red
            case 'warning':
                return '#ed6c02'; // Orange
            default:
                return '#666';
        }
    };

    const getSeverityIcon = (severity) => {
        switch (severity) {
            case 'error':
                return <ErrorIcon sx={{ color: '#d32f2f' }} />;
            case 'warning':
                return <WarningIcon sx={{ color: '#ed6c02' }} />;
            default:
                return null;
        }
    };

    const getStatusColor = (status) => {
        switch (status) {
            case 'denied':
                return '#d32f2f'; // Red
            case 'allowed':
                return '#2e7d32'; // Green
            case 'whitelisted':
                return '#1976d2'; // Blue
            default:
                return '#666';
        }
    };

    const getConflictTypeLabel = (type) => {
        switch (type) {
            case 'cidr_covers_ip':
                return 'CIDR Covers IP';
            case 'ip_in_cidr':
                return 'IP in CIDR';
            case 'cidr_overlaps_cidr':
                return 'CIDR Overlaps CIDR';
            case 'exact_match':
                return 'Exact Match';
            default:
                return type;
        }
    };

    const errorCount = conflicts.filter(c => c.severity === 'error').length;
    const warningCount = conflicts.filter(c => c.severity === 'warning').length;

    return (
        <Collapse in={open}>
            <Box sx={{ mt: 2, mb: 2 }}>
                <Alert 
                    severity="warning" 
                    sx={{ mb: 2 }}
                    action={
                        <Stack direction="row" spacing={2} alignItems="center">
                            <Typography variant="body2" sx={{ fontWeight: 'bold' }}>
                                {errorCount} Errors, {warningCount} Warnings
                            </Typography>
                            {onDeleteAllConflicts && errorCount > 0 && (
                                <Button
                                    variant="outlined"
                                    color="error"
                                    size="small"
                                    startIcon={<DeleteSweepIcon />}
                                    onClick={onDeleteAllConflicts}
                                    disabled={isDeleting}
                                    sx={{ 
                                        borderColor: '#d32f2f',
                                        color: '#d32f2f',
                                        '&:hover': {
                                            borderColor: '#b71c1c',
                                            backgroundColor: '#ffebee'
                                        }
                                    }}
                                >
                                    {isDeleting ? 'Deleting...' : `Delete ${errorCount} Error Conflicts`}
                                </Button>
                            )}
                        </Stack>
                    }
                >
                    Conflicts detected with existing IP addresses/CIDR blocks
                </Alert>
                
                <Paper elevation={2}>
                    <TableContainer>
                        <Table size="small">
                            <TableHead>
                                <TableRow sx={{ backgroundColor: '#f5f5f5' }}>
                                    <TableCell sx={{ fontWeight: 'bold' }}>Type</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>Severity</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>Conflicting Address</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>Status</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>Message</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {conflicts.map((conflict, index) => (
                                    <TableRow 
                                        key={index}
                                        sx={{
                                            backgroundColor: conflict.severity === 'error' 
                                                ? '#ffebee' 
                                                : conflict.severity === 'warning' 
                                                    ? '#fff3e0' 
                                                    : 'inherit',
                                            '&:hover': {
                                                backgroundColor: conflict.severity === 'error' 
                                                    ? '#ffcdd2' 
                                                    : conflict.severity === 'warning' 
                                                        ? '#ffe0b2' 
                                                        : '#f5f5f5'
                                            }
                                        }}
                                    >
                                        <TableCell>
                                            <Chip 
                                                label={getConflictTypeLabel(conflict.type)}
                                                size="small"
                                                variant="outlined"
                                            />
                                        </TableCell>
                                        <TableCell>
                                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                {getSeverityIcon(conflict.severity)}
                                                <Chip 
                                                    label={conflict.severity.toUpperCase()}
                                                    size="small"
                                                    sx={{
                                                        backgroundColor: getSeverityColor(conflict.severity),
                                                        color: 'white',
                                                        fontWeight: 'bold'
                                                    }}
                                                />
                                            </Box>
                                        </TableCell>
                                        <TableCell>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                                                {conflict.conflicting[0]}
                                            </Typography>
                                        </TableCell>
                                        <TableCell>
                                            <Chip 
                                                label={conflict.status}
                                                size="small"
                                                sx={{
                                                    backgroundColor: getStatusColor(conflict.status),
                                                    color: 'white',
                                                    fontWeight: 'bold'
                                                }}
                                            />
                                        </TableCell>
                                        <TableCell>
                                            <Typography variant="body2" sx={{ fontSize: '0.875rem' }}>
                                                {conflict.message}
                                            </Typography>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Paper>
            </Box>
        </Collapse>
    );
};

export default ConflictTable; 