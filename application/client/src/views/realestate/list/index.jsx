import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux'; // You need to set up Redux for state management

function RealEstate() {
  const dispatch = useDispatch();
  const accountId = useSelector((state) => state.account.accountId); // Assuming you have Redux state
  const roles = useSelector((state) => state.account.roles); // Replace with your Redux state structure
  const userName = useSelector((state) => state.account.userName); // Replace with your Redux state structure
  const balance = useSelector((state) => state.account.balance); // Replace with your Redux state structure

  const [realEstateList, setRealEstateList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingDialog, setLoadingDialog] = useState(false);
  const [dialogCreateSelling, setDialogCreateSelling] = useState(false);
  const [dialogCreateDonating, setDialogCreateDonating] = useState(false);

  const [realForm, setRealForm] = useState({
    price: 0,
    salePeriod: 0,
  });

  const [DonatingForm, setDonatingForm] = useState({
    proprietor: '',
  });

  const [accountList, setAccountList] = useState([]);

  const [valItem, setValItem] = useState({});

  useEffect(() => {
    const checkArea = (value) => {
      if (value <= 0) {
        return 'Must be greater than 0';
      } else {
        return '';
      }
    };

    const fetchRealEstateList = () => {
      // Replace with your API call to fetch real estate list
      // Simulating an API call here
      setTimeout(() => {
        const response = []; // Replace with actual response data
        setRealEstateList(response);
        setLoading(false);
      }, 1000);
    };

    const fetchAccountList = () => {
      // Replace with your API call to fetch account list
      // Simulating an API call here
      setTimeout(() => {
        const response = []; // Replace with actual response data
        const filteredAccountList = response.filter(
          (item) => item.userName !== 'Admin'
        );
        setAccountList(filteredAccountList);
      }, 1000);
    };

    // Fetch real estate list based on user role
    if (roles[0] === 'admin') {
      fetchRealEstateList();
    } else {
      // Fetch real estate list for a specific user
      // Replace 'accountType' with the actual user-specific identifier
      fetchRealEstateList(accountId);
    }

    // Fetch account list
    fetchAccountList();
  }, [roles, accountId]);

  const openDialog = (item) => {
    setDialogCreateSelling(true);
    setValItem(item);
  };

  const openDonatingDialog = (item) => {
    setDialogCreateDonating(true);
    setValItem(item);
  };

  const createSelling = (formName) => {
    // Validation and API call logic goes here
  };

  const createDonating = (formName) => {
    // Validation and API call logic goes here
  };

  const resetForm = (formName) => {
    // Reset form logic
  };

  const selectGet = (accountId) => {
    setDonatingForm({ proprietor: accountId });
  };

  return (
    <div className="container">
      <div className="el-alert" type="success">
        <p>Account ID: {accountId}</p>
        <p>Username: {userName}</p>
        <p>Balance: ${balance}</p>
        <p>When initiating a sale, donation, or pledge, the collateral status is true</p>
        <p>You can initiate a sale, donation, or pledge operation only when the collateral status is false</p>
      </div>

      {realEstateList.length === 0 && (
        <div style={{ textAlign: 'center' }}>
          <div className="el-alert" title="No data found" type="warning" />
        </div>
      )}

      <div className="el-row" style={{ marginBottom: '20px' }}>
        {realEstateList.map((val, index) => (
          <div
            key={index}
            className="el-col"
            style={{ flex: '0 0 25%', marginRight: '8px' }}
          >
            <div className="realEstate-card">
              <div className="clearfix" style={{ marginBottom: '10px' }}>
                Collateral Status:
                <span style={{ color: 'red' }}>{val.encumbrance}</span>
              </div>

              <div className="item">
                <el-tag>Real Estate ID: </el-tag>
                <span>{val.realEstateId}</span>
              </div>
              <div className="item">
                <el-tag type="success">Owner ID: </el-tag>
                <span>{val.proprietor}</span>
              </div>
              <div className="item">
                <el-tag type="warning">Total Area: </el-tag>
                <span>{val.totalArea} m²</span>
              </div>
              <div className="item">
                <el-tag type="danger">Living Space: </el-tag>
                <span>{val.livingSpace} m²</span>
              </div>

              {!val.encumbrance && roles[0] !== 'admin' && (
                <div>
                  <button
                    type="button"
                    onClick={() => openDialog(val)}
                  >
                    Sell
                  </button>
                  <el-divider direction="vertical" />
                  <button
                    type="button"
                    onClick={() => openDonatingDialog(val)}
                  >
                    Donate
                  </button>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      <div>
        <el-dialog
          visible={dialogCreateSelling}
          closeOnClickModal={false}
          onClose={() => resetForm('realForm')}
        >
          {/* Sell dialog content */}
        </el-dialog>

        <el-dialog
          visible={dialogCreateDonating}
          closeOnClickModal={false}
          onClose={() => resetForm('DonatingForm')}
        >
          {/* Donate dialog content */}
        </el-dialog>
      </div>
    </div>
  );
}

export default RealEstate;
