import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, Navigate, useLocation } from 'react-router-dom';
import IPForm from './components/IPForm';
import EmailForm from './components/EmailForm';
import UserAgentForm from './components/UserAgentForm';
import CountryForm from './components/CountryForm';
import FilterForm from './components/FilterForm';
import DashboardStats from './components/DashboardStats';
import CharsetForm from './components/CharsetForm';
import UsernameForm from './components/UsernameForm';
import SystemHealthPage from './pages/SystemHealthPage';
import MaintenancePage from './pages/MaintenancePage';
import './components/styles.css';

// MUI imports
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import CssBaseline from '@mui/material/CssBaseline';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import DashboardIcon from '@mui/icons-material/Dashboard';
import PublicIcon from '@mui/icons-material/Public';
import EmailIcon from '@mui/icons-material/Email';
import DnsIcon from '@mui/icons-material/Dns';
import FilterListIcon from '@mui/icons-material/FilterList';
import PersonIcon from '@mui/icons-material/Person';
import HealthAndSafetyIcon from '@mui/icons-material/HealthAndSafety';
import BuildIcon from '@mui/icons-material/Build';

const drawerWidth = 220;

const navItems = [
  { text: 'Dashboard', icon: <DashboardIcon />, path: '/' },
  { text: 'System Health', icon: <HealthAndSafetyIcon />, path: '/system-health' },
  { text: 'Maintenance', icon: <BuildIcon />, path: '/maintenance' },
  { text: 'Charset List', icon: <DnsIcon />, path: '/charset-list' },
  { text: 'Country List', icon: <DnsIcon />, path: '/country-list' },
  { text: 'Email List', icon: <EmailIcon />, path: '/email-list' },
  { text: 'IP List', icon: <PublicIcon />, path: '/ip-list' },
  { text: 'User Agent List', icon: <PersonIcon />, path: '/useragent-list' },
  { text: 'Username List', icon: <DnsIcon />, path: '/username-list' },
  { text: 'Filter List', icon: <FilterListIcon />, path: '/filter-list' },
];

function Dashboard() {
  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>Welcome to the Admin Dashboard</Typography>
      <Typography variant="body1" sx={{ mb: 4 }}>
        Manage your firewall filters and monitor system statistics. Select a list from the menu to manage specific filters.
      </Typography>
      
      <DashboardStats />
    </Box>
  );
}

function Layout({ children }) {
  const location = useLocation();
  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
        <Toolbar>
          <Typography variant="h6" noWrap component="div">
            Firewall Management
          </Typography>
        </Toolbar>
      </AppBar>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box' },
        }}
      >
        <Toolbar />
        <Box sx={{ overflow: 'auto' }}>
          <List>
            {navItems.map((item) => (
              <ListItem key={item.text} disablePadding>
                <ListItemButton
                  component={Link}
                  to={item.path}
                  selected={location.pathname === item.path}
                >
                  <ListItemIcon>{item.icon}</ListItemIcon>
                  <ListItemText primary={item.text} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          <Divider />
        </Box>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, bgcolor: 'background.default', p: 3 }}>
        <Toolbar />
        {children}
      </Box>
    </Box>
  );
}

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/system-health" element={<SystemHealthPage />} />
          <Route path="/maintenance" element={<MaintenancePage />} />
          <Route path="/ip-list" element={<IPForm />} />
          <Route path="/email-list" element={<EmailForm />} />
          <Route path="/useragent-list" element={<UserAgentForm />} />
          <Route path="/country-list" element={<CountryForm />} />
          <Route path="/filter-list" element={<FilterForm />} />
          <Route path="/charset-list" element={<CharsetForm />} />
          <Route path="/username-list" element={<UsernameForm />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;
