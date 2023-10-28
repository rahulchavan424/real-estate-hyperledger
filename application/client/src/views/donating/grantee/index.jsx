import React, { useState, useEffect } from 'react';

function DonatingGrantee() {
  const [loading, setLoading] = useState(true);
  const [donatingList, setDonatingList] = useState([]);
  const accountId = 'Your Account ID'; // Replace with your actual account ID
  const userName = 'Your Username'; // Replace with your actual username
  const balance = 'Your Balance'; // Replace with your actual balance

  const queryDonatingListByGrantee = () => {
    // Simulate the API call for querying donating list by grantee
    // Replace this with your actual API call
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([]);
      }, 1000);
    });
  };

  const updateDonating = (item, type) => {
    // Simulate the API call for updating donating
    // Replace this with your actual API call
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({});
      }, 1000);
    });
  };

  useEffect(() => {
    queryDonatingListByGrantee()
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

  const handleUpdateDonating = (item, type) => {
    let tip = '';
    if (type === 'done') {
      tip = 'Confirm Receipt';
    } else {
      tip = 'Cancel Donation';
    }

    if (window.confirm(`Do you want to ${tip}?`)) {
      setLoading(true);
      updateDonating({
        donor: item.donating.donor,
        grantee: item.donating.grantee,
        objectOfDonating: item.donating.objectOfDonating,
        status: type,
      })
        .then((response) => {
          setLoading(false);
          if (response !== null) {
            alert(`${tip} succeeded!`);
          } else {
            alert(`${tip} failed!`);
          }
          setTimeout(() => {
            window.location.reload();
          }, 1000);
        })
        .catch((_) => {
          setLoading(false);
        });
    } else {
      alert(`${tip} canceled`);
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
              <div className="el-card d-buy-card">
                <div className="clearfix">
                  <span>{val.donating.donatingStatus}</span>
                  {val.donating.donatingStatus === 'In Progress' && (
                    <>
                      <button
                        style={{
                          float: 'right',
                          padding: '3px 0',
                          background: 'none',
                          border: 'none',
                          textDecoration: 'underline',
                        }}
                        onClick={() => handleUpdateDonating(val, 'done')}
                      >
                        Confirm Receipt
                      </button>
                      <button
                        style={{
                          float: 'right',
                          padding: '3px 6px',
                          background: 'none',
                          border: 'none',
                          textDecoration: 'underline',
                        }}
                        onClick={() => handleUpdateDonating(val, 'cancelled')}
                      >
                        Cancel
                      </button>
                    </>
                  )}
                </div>
                <div className="item">
                  <span>Property ID: {val.donating.objectOfDonating}</span>
                </div>
                <div className="item">
                  <span>Donor ID: {val.donating.donor}</span>
                </div>
                <div className="item">
                  <span>Recipient ID: {val.donating.grantee}</span>
                </div>
                <div className="item">
                  <span>Created at: {val.donating.createTime}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default DonatingGrantee;
