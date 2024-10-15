// api.js
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://127.0.0.1:3000',
  headers: { 
    'Content-Type': 'application/json',
    // 'Authorization': `Bearer ${localStorage.getItem('jwt_token')}`,
  },
});

export default api;
// This is a simple Axios instance that we can use to make requests to our backend API. 
//We set the baseURL to http://127.0.0.1:3000
