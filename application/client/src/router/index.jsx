import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import RealEstate from '../views/realestate/list';
import AllSelling from '../views/selling/all';
import MeSelling from '../views/selling/me';
import BuySelling from '../views/selling/buy';
import AllDonating from '../views/donating/all';
import DonatingDonor from '../views/donating/donor';
import DonatingGrantee from '../views/donating/grantee';
import AddRealEstate from '../views/realestate/add';
import NotFound from '../views/404';

function AppRouter() {
  return (
    <Router>
        <Routes>
          <Route path="/realestate" component={RealEstate} />
          <Route path="/selling/all" component={AllSelling} />
          <Route path="/selling/me" component={MeSelling} />
          <Route path="/selling/buy" component={BuySelling} />
          <Route path="/donating/all" component={AllDonating} />
          <Route path="/donating/donor" component={DonatingDonor} />
          <Route path="/donating/grantee" component={DonatingGrantee} />
          <Route path="/addRealestate" component={AddRealEstate} />
          <Route path="/404" component={NotFound} />
          <Route exact path="/" component={RealEstate} />
          <Route path="*" component={NotFound} />
        </Routes>
    </Router>
  );
}

export default AppRouter;
