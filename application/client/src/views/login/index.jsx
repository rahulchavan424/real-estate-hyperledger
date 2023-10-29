import React, { useState, useEffect } from 'react';
import { queryAccountList } from '../../api/account';
import { useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';

function Login() {
  const [loading, setLoading] = useState(false);
  const [redirect, setRedirect] = useState(undefined);
  const [accountList, setAccountList] = useState([]);
  const [value, setValue] = useState('');
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const handleLogin = () => {
    if (value) {
      setLoading(true);
      // Simulate the login action
      // Replace this with your actual login logic
      setTimeout(() => {
        setLoading(false);
        // Dispatch user data to Redux store
        dispatch({
          type: 'SET_ACCOUNT',
          payload: {
            accountId: value, // Assuming value is the selected user's accountId
            userName: 'User Name', // Replace with actual user name
            roles: ['user'], // Replace with actual user roles
          },
        });
        // Redirect after login
        // Replace this with your actual redirect logic
        navigate('/realestate'); // Redirect to the default path or your desired path
      }, 2000);
    } else {
      alert('Please select a user role'); // Replace this with your actual notification method
    }
  }

  const selectGet = (e) => {
    setValue(e.target.value);
  };

  useEffect(() => {
    // Load account list on component mount
    queryAccountList()
      .then((response) => {
        if (response !== null) {
          setAccountList(response);
        }
      });
  }, []);

  return (
    <div className="login-container">
      <div className="login-form" autoComplete="on">
        <div className="title-container">
          <h3 className="title">Real Estate Transaction System Based on Blockchain</h3>
        </div>
        <select value={value} onChange={selectGet} className="login-select">
          <option value="" disabled>Select user role</option>
          {accountList.map((item) => (
            <option key={item.accountId} value={item.accountId}>
              {item.userName} - {item.accountId}
            </option>
          ))}
        </select>
        <button
          className="el-button"
          style={{ width: '100%', marginBottom: '30px' }}
          onClick={handleLogin}
          disabled={loading}
        >
          {loading ? 'Loading...' : 'Enter'}
        </button>
        <div className="tips">
          <span>Tips: Choose different user roles to simulate transactions</span>
        </div>
      </div>
    </div>
  );
}

export default Login;
