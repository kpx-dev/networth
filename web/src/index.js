import React from "react";
import ReactDOM from "react-dom";
import Amplify from "aws-amplify";
import "./assets/scss/material-kit-react.css?v=1.2.0";
import App from "./App";
import registerServiceWorker from "./registerServiceWorker";

Amplify.configure({
  Auth: {
    region: "us-east-1",
    userPoolId: "us-east-1_5cJz62UiG",
    userPoolWebClientId: "2tam11a22g38in2vqcd5kge3cu",
    mandatorySignIn: false
  }
});

ReactDOM.render(<App />, document.getElementById("root"));

registerServiceWorker();
