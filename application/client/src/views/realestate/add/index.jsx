import React, { useState, useEffect } from 'react';

function AddRealEstate() {
  const [loading, setLoading] = useState(false);
  const [accountList, setAccountList] = useState([]);
  const [ruleForm, setRuleForm] = useState({
    proprietor: '',
    totalArea: 0,
    livingSpace: 0,
  });

  const checkArea = (value) => {
    return value > 0 ? '' : 'Must be greater than 0';
  };

  const submitForm = (e) => {
    e.preventDefault();
    setLoading(true);

    // Simulate the createRealEstate API call
    // Replace this with your actual API call
    setTimeout(() => {
      setLoading(false);
      // Simulate a successful response
      // Replace this with actual API response handling
      const response = true;

      if (response) {
        alert('Creation successful!'); // Replace with your preferred notification method
      } else {
        alert('Creation failed!'); // Replace with your preferred notification method
      }
    }, 2000);
  };

  const resetForm = () => {
    setRuleForm({
      proprietor: '',
      totalArea: 0,
      livingSpace: 0,
    });
  };

  const selectGet = (e) => {
    setRuleForm({
      ...ruleForm,
      proprietor: e.target.value,
    });
  };

  useEffect(() => {
    // Load account list on component mount
    // Simulate the queryAccountList API call
    // Replace this with your actual API call
    const queryAccountList = () => {
      return new Promise((resolve) => {
        setTimeout(() => {
          // Simulate a response with accountList data
          // Replace this with actual API response data
          const accountList = [
            { accountId: 1, userName: 'User1' },
            { accountId: 2, userName: 'User2' },
            // Add more account data as needed
          ];
          resolve(accountList);
        }, 1000);
      });
    };

    queryAccountList().then((response) => {
      if (response !== null) {
        // Filter out the admin user
        const filteredAccountList = response.filter(
          (item) => item.userName !== 'Admin'
        );
        setAccountList(filteredAccountList);
      }
    });
  }, []);

  return (
    <div className="app-container">
      <form className="el-form" onSubmit={submitForm}>
        <div className="el-form-item">
          <label className="el-form-item__label">Owner</label>
          <select
            value={ruleForm.proprietor}
            onChange={selectGet}
            className="el-select"
          >
            <option value="" disabled>
              Select Owner
            </option>
            {accountList.map((item) => (
              <option key={item.accountId} value={item.accountId}>
                {item.userName} - {item.accountId}
              </option>
            ))}
          </select>
        </div>
        <div className="el-form-item">
          <label className="el-form-item__label">Total Area (㎡)</label>
          <input
            type="number"
            value={ruleForm.totalArea}
            onChange={(e) => {
              setRuleForm({ ...ruleForm, totalArea: e.target.value });
            }}
            step="0.1"
            min="0"
          />
        </div>
        <div className="el-form-item">
          <label className="el-form-item__label">Living Space (㎡)</label>
          <input
            type="number"
            value={ruleForm.livingSpace}
            onChange={(e) => {
              setRuleForm({ ...ruleForm, livingSpace: e.target.value });
            }}
            step="0.1"
            min="0"
          />
        </div>
        <div className="el-form-item">
          <button type="submit" className="el-button" disabled={loading}>
            {loading ? 'Loading...' : 'Create Now'}
          </button>
          <button
            type="button"
            className="el-button"
            onClick={resetForm}
          >
            Reset
          </button>
        </div>
      </form>
    </div>
  );
}

export default AddRealEstate;
