import React from "react";
import ReactDOM from "react-dom";
import Amplify from "aws-amplify";
import "./assets/scss/material-kit-react.css?v=1.2.0";
import App from "./App";
import registerServiceWorker from "./registerServiceWorker";
import { AWS_REGION, COGNITO_CLIENT_ID, COGNITO_POOL_ID } from "./helpers/constants.js";
Amplify.configure({
  Auth: {
    region: AWS_REGION,
    userPoolId: COGNITO_POOL_ID,
    userPoolWebClientId: COGNITO_CLIENT_ID,
    mandatorySignIn: false
  }
});

ReactDOM.render(<App />, document.getElementById("root"));

registerServiceWorker();
