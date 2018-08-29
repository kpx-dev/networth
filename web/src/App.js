import React, { Component } from "react";
import { Router, Route, Switch } from "react-router-dom";
import indexRoutes from "./routes/index.jsx";
import { createBrowserHistory } from "history";
import { Auth } from "aws-amplify";
import { withAuthenticator } from "aws-amplify-react";

const hist = createBrowserHistory();

class App extends Component {
  render() {
    return (
      <Router history={hist}>
        <Switch>
          {indexRoutes.map((prop, key) => {
            return (
              <Route path={prop.path} key={key} component={prop.component} />
            );
          })}
        </Switch>
      </Router>
    );
  }
}

export default withAuthenticator(App);
