import axios from 'axios';

const imageApi = axios.create({
    baseURL: 'http://127.0.0.1:3000',
    headers: { 'Content-Type': 'multipart/form-data' }, // Explicitly set headers
});


export default imageApi;