import React from "react";
import ReactDOM from "react-dom";
import { createBrowserHistory } from "history";
import { Router, Route, Switch } from "react-router-dom";
import Amplify from "aws-amplify";

import indexRoutes from "routes/index.jsx";

import "assets/scss/material-kit-react.css?v=1.2.0";

var hist = createBrowserHistory();

Amplify.configure({
  Auth: {
    region: "us-east-1",
    userPoolId: "us-east-1_5cJz62UiG",
    userPoolWebClientId: "2tam11a22g38in2vqcd5kge3cu",
    mandatorySignIn: false
  }
});

ReactDOM.render(
  <Router history={hist}>
    <Switch>
      {indexRoutes.map((prop, key) => {
        return <Route path={prop.path} key={key} component={prop.component} />;
      })}
    </Switch>
  </Router>,
  document.getElementById("root")
);
