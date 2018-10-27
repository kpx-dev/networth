import React from "react";
import PropTypes from "prop-types";
import PlaidLink from "react-plaid-link";
import { post } from "../../helpers/helpers.js";
import { PLAID_CLIENT_NAME, PLAID_PUBLIC_KEY, PLAID_ENV, PLAID_PRODUCTS, PLAID_WEBHOOK } from "../../helpers/constants.js"
import {
  Button,
} from "reactstrap";

class Connect extends React.Component {
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

  static propTypes = {
    institution: PropTypes.string,
    text: PropTypes.string
  };

  static defaultProps = {
    text: "+ Connect"
  };

  render() {
    const { institution, text } = this.props;
    const plaidStyle = { padding: 0, border: "none", borderRadius: 0, background: "none" };

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
        onSuccess={this.handleOnSuccess}
      >
        <Button color="danger" size="sm">{text}</Button>
      </PlaidLink>
    );
  }
}

export default Connect;
