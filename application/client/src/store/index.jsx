// 1. Import necessary Redux functions
import { createStore, combineReducers } from 'redux';

// 2. Define Redux actions (optional)
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

// 3. Define Redux reducers
const counterReducer = (state = 0, action) => {
  switch (action.type) {
    case 'INCREMENT':
      return state + 1;
    case 'DECREMENT':
      return state - 1;
    default:
      return state;
  }
};

// 4. Combine reducers (if you have multiple reducers)
const rootReducer = combineReducers({
  counter: counterReducer,
  // Add more reducers here if needed
});

// 5. Create the Redux store
const store = createStore(rootReducer);

export default store;
