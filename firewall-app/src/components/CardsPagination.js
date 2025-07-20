import React from 'react';
import Box from '@mui/material/Box';
import Pagination from '@mui/material/Pagination';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';

const CardsPagination = ({
    total,
    page,
    rowsPerPage,
    onPageChange,
    onRowsPerPageChange,
    rowsPerPageOptions = [10, 25, 50, 100]
}) => {
    const totalPages = Math.ceil(total / rowsPerPage);
    const startItem = page * rowsPerPage + 1;
    const endItem = Math.min((page + 1) * rowsPerPage, total);

    const handlePageChange = (event, newPage) => {
        onPageChange(event, newPage - 1); // Convert to 0-based index
    };

    const handleRowsPerPageChange = (event) => {
        onRowsPerPageChange(event);
    };

    return (
        <Paper sx={{ p: 2, mb: 2 }}>
            <Box sx={{ 
                display: 'flex', 
                justifyContent: 'space-between', 
                alignItems: 'center',
                flexWrap: 'wrap',
                gap: 2
            }}>
                {/* Results Info */}
                <Typography variant="body2" color="text.secondary">
                    Showing {startItem}-{endItem} of {total} entries
                </Typography>

                {/* Pagination Controls */}
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    {/* Rows per page */}
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Typography variant="body2">Entries per page:</Typography>
                        <FormControl size="small" sx={{ minWidth: 80 }}>
                            <Select
                                value={rowsPerPage}
                                onChange={handleRowsPerPageChange}
                                displayEmpty
                            >
                                {rowsPerPageOptions.map((option) => (
                                    <MenuItem key={option} value={option}>
                                        {option}
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                    </Box>

                    {/* Page navigation */}
                    <Pagination
                        count={totalPages}
                        page={page + 1} // Convert to 1-based index
                        onChange={handlePageChange}
                        showFirstButton
                        showLastButton
                        size="small"
                        color="primary"
                    />
                </Box>
            </Box>
        </Paper>
    );
};

export default CardsPagination; 