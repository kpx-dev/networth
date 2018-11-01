import React from "react";
import {
  Card,
  CardBody,
  CardHeader,
  CardTitle,
  Table,
  Row,
  Col,
  DropdownToggle,
  DropdownMenu,
  DropdownItem,
  UncontrolledDropdown,
} from "reactstrap";
import { get } from "../../helpers/helpers.js";
import NotificationAlert from "react-notification-alert";

class Accounts extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      dropdownOpen: false,
      accounts: {}
    };
  }

  async componentDidMount() {
    try {
      const res = await get(`/accounts`);
      const body = await res.json();
      this.setState({ loading: false, accounts: body.data });
      console.log(body.data);
    } catch (e) {
      this.alert('Cannot get accounts. Problem connecting to REST API.');
      this.setState({ loading: false });
    }
  }

  dropdownToggle(e) {
    console.log(e);
    this.setState({
      dropdownOpen: !this.state.dropdownOpen
    });
  }

  alert(message, dismiss) {
    const options = {
      place: "tc",
      type: "danger",
      icon: "nc-icon nc-alert-circle-i",
      message: (
        <div>
          <div>
            {message}
          </div>
        </div>
      ),
    };
    if (dismiss) options.autoDismiss = dismiss;
    this.refs.alert.notificationAlert(options);
  }

  render() {
    return (
      <div className="content">
        <NotificationAlert ref="alert" />
        <h4>Accounts</h4>
        <Row>
          <Col xs={12}>
          {Object.keys(this.state.accounts).map((key, _) => {
            return (
              <Card>
              <CardHeader>
                <CardTitle tag="h6">{this.state.accounts[key].institution_name}</CardTitle>
              </CardHeader>
              <CardBody>
                <Table responsive>
                  <thead className="text-primary">
                    <tr>
                      <th width="40%">Name</th>
                      <th width="40%">Official Name</th>
                      <th width="10%">Balance</th>
                      <th className="text-right">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {this.state.accounts[key].accounts.map((account, idx) => {
                      return (
                        <tr>
                          <td>{account.name}</td>
                          <td>{account.official_name}</td>
                          <td>{
                            new Intl.NumberFormat('en-US', {
                              style: 'currency',
                              currency: 'USD',
                              minimumFractionDigits: 2,
                              maximumFractionDigits: 2,
                            }).format(account.balances.current)}</td>
                          <td key={`${key}_${idx}`} className="text-right">
                          <UncontrolledDropdown>
                            <DropdownToggle caret nav></DropdownToggle>
                            <DropdownMenu right>
                              <DropdownItem tag="a">Connect</DropdownItem>
                            </DropdownMenu>
                          </UncontrolledDropdown>

                          </td>
                        </tr>
                      )
                    })}
                    {/* {tbody.map((prop, key) => {
                      return (
                        <tr key={key}>
                          {prop.data.map((prop, key) => {
                            if (key === thead.length - 1)
                              return (
                                <td key={key} className="text-right">
                                  {prop}
                                </td>
                              );
                            return <td key={key}>{prop}</td>;
                          })}
                        </tr>
                      );
                    })} */}
                  </tbody>
                </Table>
              </CardBody>
            </Card>
            )
          })}

          </Col>
          </Row>



      </div>
    );
  }
}

export default Accounts;
