import React from 'react';
import Box from '@mui/material/Box';
import ToggleButton from '@mui/material/ToggleButton';
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup';
import TableChartIcon from '@mui/icons-material/TableChart';
import ViewModuleIcon from '@mui/icons-material/ViewModule';
import Tooltip from '@mui/material/Tooltip';

const LayoutToggle = ({ layout, onLayoutChange }) => {
    const handleLayoutChange = (event, newLayout) => {
        if (newLayout !== null) {
            onLayoutChange(newLayout);
        }
    };

    return (
        <Box sx={{ mb: 2 }}>
            <ToggleButtonGroup
                value={layout}
                exclusive
                onChange={handleLayoutChange}
                size="small"
                color="primary"
            >
                <Tooltip title="Table Layout">
                    <ToggleButton value="table">
                        <TableChartIcon />
                    </ToggleButton>
                </Tooltip>
                <Tooltip title="Cards Layout">
                    <ToggleButton value="cards">
                        <ViewModuleIcon />
                    </ToggleButton>
                </Tooltip>
            </ToggleButtonGroup>
        </Box>
    );
};

export default LayoutToggle; 