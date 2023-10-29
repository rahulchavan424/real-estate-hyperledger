const initialState = {
    accountId: '',
    userName: '',
    roles: [],
  };
  
  const accountReducer = (state = initialState, action) => {
    switch (action.type) {
      case 'SET_ACCOUNT':
        return {
          ...state,
          accountId: action.payload.accountId,
          userName: action.payload.userName,
          roles: action.payload.roles,
        };
      default:
        return state;
    }
  };
  
  export default accountReducer;