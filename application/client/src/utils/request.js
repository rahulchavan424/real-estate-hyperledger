import axios from 'axios';
import { message } from 'antd'; // Assuming you're using Ant Design for notifications

const service = axios.create({
  baseURL: '/api/v1', // Update this environment variable as needed
  timeout: 5000,
});

service.interceptors.response.use(
  (response) => {
    const res = response.data;
    if (res.code !== 200) {
      message.warning('Server encountered an issue', 5); // Display a warning message
      return Promise.reject(new Error(res.msg || 'Error'));
    } else {
      return res.data;
    }
  },
  (error) => {
    if (!error.response) {
      message.error('Request failed: ' + error.message, 5); // Display an error message
      return Promise.reject(error);
    } else {
      message.error('Failure: ' + error.response.data.data, 5); // Display an error message
      return Promise.reject(error.response);
    }
  }
);

export default service;
