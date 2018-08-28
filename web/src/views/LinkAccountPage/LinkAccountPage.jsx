import React from "react";
import withStyles from "@material-ui/core/styles/withStyles";
import loginPageStyle from "assets/jss/material-kit-react/views/loginPage.jsx";
import PlaidLink from "react-plaid-link";

class LinkAccountPage extends React.Component {
  constructor(props) {
    super(props);
    // we use this to make the card to appear after the page has been rendered
    this.state = {
      cardAnimaton: "cardHidden"
    };
  }
  componentDidMount() {
    // we add a hidden class to the card and after 700 ms we delete it and the transition appears
    setTimeout(
      function() {
        this.setState({ cardAnimaton: "" });
      }.bind(this),
      700
    );
  }
  handleOnSuccess(token, metadata) {
    console.log(token, metadata);
  }
  handleOnExit() {
    // handle the case when your user exits Link
  }
  render() {
    // const { classes, ...rest } = this.props;
    return (
      <PlaidLink
        // institution={null}
        webhook="https://api.networth.app/webhook"
        clientName="networth.app"
        env="sandbox"
        product={["transactions"]}
        publicKey="7e599ac974fb8343f50fac8535fcf1"
        onExit={this.handleOnExit}
        onSuccess={this.handleOnSuccess}
      >
        Link Account
      </PlaidLink>
    );
  }
}

export default withStyles(loginPageStyle)(LinkAccountPage);
