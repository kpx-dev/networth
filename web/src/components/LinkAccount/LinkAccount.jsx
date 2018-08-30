import React from "react";
import PropTypes from "prop-types";
import PlaidLink from "react-plaid-link";
import AddIcon from "@material-ui/icons/Add";
import Button from "../../components/CustomButtons/Button.jsx";
import withStyles from "@material-ui/core/styles/withStyles";
import buttonStyle from "../../assets/jss/material-kit-react/components/buttonStyle.jsx";
import { Auth } from "aws-amplify";

class LinkAccount extends React.Component {
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

  async handleOnSuccess(token, metadata) {
    const session = await Auth.currentSession();
    const exchangeUrl = "https://api.networth.app/tokens/exchange";
    const fetchOptions = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session.idToken.jwtToken}`
      },
      body: JSON.stringify({ token })
    };
    await fetch(exchangeUrl, fetchOptions);
  }

  // handleOnExit() {
  //   // handle the case when your user exits Link
  // }

  static propTypes = {
    classes: PropTypes.object.isRequired,
    institution: PropTypes.string,
    text: PropTypes.string
  };

  static defaultProps = {
    text: "Link Account"
  };

  render() {
    const { classes, institution, text } = this.props;
    const plaidStyle = { padding: 0, border: "none", borderRadius: 0 };

    return (
      <PlaidLink
        style={plaidStyle}
        className="plaid-link"
        institution={institution}
        webhook="https://api.networth.app/webhook"
        clientName="networth.app"
        env="sandbox"
        product={["transactions"]}
        publicKey="7e599ac974fb8343f50fac8535fcf1"
        // onExit={this.handleOnExit}
        onSuccess={this.handleOnSuccess}
      >
        <Button size="sm" round>
          <AddIcon className={classes.icons} /> {text}
        </Button>
      </PlaidLink>
    );
  }
}

export default withStyles(buttonStyle)(LinkAccount);
