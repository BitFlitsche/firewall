// src/axiosConfig.js
import axios from 'axios';

const axiosInstance = axios.create({
    baseURL: 'http://localhost:8081', // Setze hier die Basis-URL
});

export default axiosInstance;
