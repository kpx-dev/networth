import React from "react";
import PropTypes from "prop-types";
import PlaidLink from "react-plaid-link";
import AddIcon from "@material-ui/icons/Add";
import Button from "../../components/CustomButtons/Button.jsx";
import withStyles from "@material-ui/core/styles/withStyles";
import buttonStyle from "../../assets/jss/material-kit-react/components/buttonStyle.jsx";
import { post } from "../../helpers/helpers.js";
import { NW_API_BASE_URL } from "../../helpers/constants.js"

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
    await post('/tokens/exchange', body);
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
    let webhook = 'http://networth.dev:3000/webhook';
    if (NW_API_BASE_URL) webhook = `${NW_API_BASE_URL}/webhook`;

    return (
      <PlaidLink
        style={plaidStyle}
        className="plaid-link"
        institution={institution}
        webhook={webhook}
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
