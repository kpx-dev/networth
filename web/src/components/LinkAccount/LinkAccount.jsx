import React from "react";
import PropTypes from "prop-types";
import PlaidLink from "react-plaid-link";
import AddIcon from "@material-ui/icons/Add";
import Button from "../../components/CustomButtons/Button.jsx";
import withStyles from "@material-ui/core/styles/withStyles";
import buttonStyle from "../../assets/jss/material-kit-react/components/buttonStyle.jsx";
import { post } from "../../helpers/helpers.js";
import { PLAID_CLIENT_NAME, PLAID_PUBLIC_KEY, PLAID_ENV, PLAID_PRODUCTS, PLAID_WEBHOOK } from "../../helpers/constants.js"

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

  async handleOnSuccess(access_token, metadata) {
    const body = {
      access_token,
      institution_id: metadata.institution.institution_id,
      institution_name: metadata.institution.name,
      account_id: metadata.account_id, // for stripe later (https://plaid.com/docs/link/stripe/)
      accounts: metadata.accounts || [],
    };
    await post('/tokens', body);
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
        webhook={PLAID_WEBHOOK}
        clientName={PLAID_CLIENT_NAME}
        env={PLAID_ENV}
        product={PLAID_PRODUCTS}
        publicKey={PLAID_PUBLIC_KEY}
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
