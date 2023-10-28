import React, { useState, useEffect } from 'react';

function DonatingDonor() {
  const [loading, setLoading] = useState(true);
  const [donatingList, setDonatingList] = useState([]);
  const accountId = 'Your Account ID'; // Replace with your actual account ID
  const userName = 'Your Username'; // Replace with your actual username
  const balance = 'Your Balance'; // Replace with your actual balance

  const queryDonatingList = () => {
    // Simulate the API call for querying donating list
    // Replace this with your actual API call
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([]);
      }, 1000);
    });
  };

  const updateDonating = (item) => {
    // Simulate the API call for updating donating
    // Replace this with your actual API call
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({});
      }, 1000);
    });
  };

  useEffect(() => {
    queryDonatingList()
      .then((response) => {
        if (response !== null) {
          setDonatingList(response);
        }
        setLoading(false);
      })
      .catch((_) => {
        setLoading(false);
      });
  }, []);

  const handleUpdateDonating = (item) => {
    if (window.confirm('Do you want to cancel the donation?')) {
      setLoading(true);
      updateDonating({
        donor: item.donor,
        grantee: item.grantee,
        objectOfDonating: item.objectOfDonating,
        status: 'cancelled',
      })
        .then((response) => {
          setLoading(false);
          if (response !== null) {
            alert('Operation Succeeded!');
          } else {
            alert('Operation Failed!');
          }
          setTimeout(() => {
            window.location.reload();
          }, 1000);
        })
        .catch((_) => {
          setLoading(false);
        });
    } else {
      alert('Operation Canceled');
    }
  };

  return (
    <div className="container">
      <div className="el-alert" style={{ backgroundColor: 'lightgreen' }}>
        <p>Account ID: {accountId}</p>
        <p>Username: {userName}</p>
        <p>Balance: ¥{balance} CNY</p>
      </div>
      {donatingList.length === 0 ? (
        <div style={{ textAlign: 'center' }}>
          <div className="el-alert" style={{ backgroundColor: 'lightcoral' }}>
            No data found
          </div>
        </div>
      ) : (
        <div className="el-row" style={{ gutter: 20 }}>
          {donatingList.map((val, index) => (
            <div key={index} style={{ span: 6, offset: 1 }}>
              <div className="el-card d-me-card">
                <div className="clearfix">
                  <span>{val.donatingStatus}</span>
                  {val.donatingStatus === 'In Progress' && (
                    <button
                      style={{
                        float: 'right',
                        padding: '3px 0',
                        background: 'none',
                        border: 'none',
                        textDecoration: 'underline',
                      }}
                      onClick={() => handleUpdateDonating(val)}
                    >
                      Cancel
                    </button>
                  )}
                </div>
                <div className="item">
                  <span>Property ID: {val.objectOfDonating}</span>
                </div>
                <div className="item">
                  <span>Donor ID: {val.donor}</span>
                </div>
                <div className="item">
                  <span>Recipient ID: {val.grantee}</span>
                </div>
                <div className="item">
                  <span>Created at: {val.createTime}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default DonatingDonor;
