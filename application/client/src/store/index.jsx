import { createStore, combineReducers } from 'redux';

// Define your initial account state
const initialAccountState = {
  accountId: null,
  roles: [],
  userName: '',
  balance: 0,
  // Add other properties as needed
};

// Define Redux actions (optional)
const incrementAction = () => {
  return {
    type: 'INCREMENT',
  };
};

const decrementAction = () => {
  return {
    type: 'DECREMENT',
  };
};

// Define the account reducer
const accountReducer = (state = initialAccountState, action) => {
  // Handle actions related to the account state here
  // For example, if you have actions that update the account state
  // you can handle them in this reducer.
  return state; // Replace with actual logic as needed
};

// Define other reducers as needed

// Combine reducers
const rootReducer = combineReducers({
  account: accountReducer,
  // Add more reducers here if needed
});

// Create the Redux store
const store = createStore(rootReducer);

export default store;