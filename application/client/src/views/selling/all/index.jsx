import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux'; // You need to set up Redux for state management

function AllSelling() {
  const dispatch = useDispatch();
  const accountId = useSelector((state) => state.account.accountId); // Assuming you have Redux state
  const roles = useSelector((state) => state.account.roles); // Replace with your Redux state structure
  const userName = useSelector((state) => state.account.userName); // Replace with your Redux state structure
  const balance = useSelector((state) => state.account.balance); // Replace with your Redux state structure

  const [sellingList, setSellingList] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchSellingList = () => {
      // Replace with your API call to fetch selling list
      // Simulating an API call here
      setTimeout(() => {
        const response = []; // Replace with actual response data
        setSellingList(response);
        setLoading(false);
      }, 1000);
    };

    // Fetch selling list
    fetchSellingList();
  }, []);

  const createSellingByBuy = (item) => {
    // Validation and API call logic goes here
  };

  const updateSelling = (item, type) => {
    // Validation and API call logic goes here
  };

  return (
    <div className="container">
      <div className="el-alert" type="success">
        <p>Account ID: {accountId}</p>
        <p>Username: {userName}</p>
        <p>Balance: ${balance}</p>
      </div>

      {sellingList.length === 0 && (
        <div style={{ textAlign: 'center' }}>
          <div className="el-alert" title="No data found" type="warning" />
        </div>
      )}

      <div className="el-row" style={{ marginBottom: '20px' }}>
        {sellingList.map((val, index) => (
          <div
            key={index}
            className="el-col"
            style={{ flex: '0 0 25%', marginRight: '8px' }}
          >
            <div className="all-card">
              <div className="clearfix" style={{ marginBottom: '10px' }}>
                <span>{val.sellingStatus}</span>
                {roles[0] !== 'admin' &&
                (val.seller === accountId || val.buyer === accountId) &&
                val.sellingStatus !== 'Completed' &&
                val.sellingStatus !== 'Expired' &&
                val.sellingStatus !== 'Cancelled' && (
                  <button
                    type="button"
                    onClick={() => updateSelling(val, 'cancelled')}
                  >
                    Cancel
                  </button>
                )}
                {roles[0] !== 'admin' &&
                val.seller === accountId &&
                val.sellingStatus === 'In Progress' && (
                  <button
                    type="button"
                    onClick={() => updateSelling(val, 'done')}
                  >
                    Confirm Payment
                  </button>
                )}
                {roles[0] !== 'admin' &&
                val.sellingStatus === 'On Sale' &&
                val.seller !== accountId && (
                  <button
                    type="button"
                    onClick={() => createSellingByBuy(val)}
                  >
                    Buy
                  </button>
                )}
              </div>

              <div className="item">
                <el-tag>Real Estate ID: </el-tag>
                <span>{val.objectOfSale}</span>
              </div>
              <div className="item">
                <el-tag type="success">Seller ID: </el-tag>
                <span>{val.seller}</span>
              </div>
              <div className="item">
                <el-tag type="danger">Price: </el-tag>
                <span>$ {val.price}</span>
              </div>
              <div className="item">
                <el-tag type="warning">Validity Period: </el-tag>
                <span>{val.salePeriod} days</span>
              </div>
              <div className="item">
                <el-tag type="info">Creation Time: </el-tag>
                <span>{val.createTime}</span>
              </div>
              <div className="item">
                <el-tag>Buyer ID: </el-tag>
                <span>
                  {val.buyer === '' ? 'Vacant' : val.buyer}
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default AllSelling;
