import React from 'react';
import AppRouter from './router';
import { Provider } from 'react-redux';
import store from './store';

function App() {
  return (
    <div id="app">
      <Provider store={store}>
        <AppRouter />
      </Provider>
    </div>
  );
}

export default App;
