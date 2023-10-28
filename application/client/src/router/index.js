// AppRouter.js
import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';

import Layout from './Layout';
import NotFound from './NotFound';
import RealEstateList from './views/realestate/list/index';
import SellingAll from './views/selling/all/index';
import SellingMe from './views/selling/me/index';
import SellingBuy from './views/selling/buy/index';
import DonatingAll from './views/donating/all/index';
import DonatingDonor from './views/donating/donor/index';
import DonatingGrantee from './views/donating/grantee/index';
import AddRealEstate from './views/realestate/add/index';

function AppRouter() {
  return (
    <Router>
      <Layout>
        <Switch>
          <Route path="/realestate" component={RealEstateList} />
          <Route path="/selling/all" component={SellingAll} />
          <Route path="/selling/me" component={SellingMe} />
          <Route path="/selling/buy" component={SellingBuy} />
          <Route path="/donating/all" component={DonatingAll} />
          <Route path="/donating/donor" component={DonatingDonor} />
          <Route path="/donating/grantee" component={DonatingGrantee} />
          <Route path="/addRealestate" component={AddRealEstate} />
          <Route path="/404" component={NotFound} />
          <Route exact path="/" component={RealEstateList} />
          <Route path="*" component={NotFound} />
        </Switch>
      </Layout>
    </Router>
  );
}

export default AppRouter;
