// App.js
import React from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';

const App = () => {
  return (
    <div id="app">
      <Router>
        <Route path="/" component={MainView} />
      </Router>
    </div>
  );
};

export default App;
