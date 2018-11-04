import React from "react";
import {
  Card,
  CardBody,
  CardHeader,
  CardTitle,
  Table,
  Row,
  Col
} from "reactstrap";
import { get } from "../../helpers/helpers.js";
import NotificationAlert from "react-notification-alert";
import last from "lodash/last";

class RegularTables extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      transactions: [],
    };
  }

  async componentDidMount() {
    const accountID = last(this.props.location.pathname.split('/'));

    try {
      this.setState({ loading: true });
      const res = await get(`/transactions/${accountID}`);
      const body = await res.json();
      console.log(body.data);
      this.setState({ loading: false, transactions: body.data });
    } catch (e) {
      // this.alert('Cannot get accounts. Problem connecting to REST API.');
      this.setState({ loading: false });
    }
  }

  render() {
    return (
      <div className="content">
        <Row>
          <Col xs={12}>
            <Card>
              <CardHeader>
                <CardTitle tag="h4">Transactions</CardTitle>
              </CardHeader>
              <CardBody>
                <Table responsive>
                  <thead className="text-primary">
                    <tr>
                      <th width="60%">Name</th>
                      <th width="20%">Date</th>
                      <th width="20%">Amount</th>
                    </tr>
                  </thead>
                  <tbody>
                  {this.state.transactions.map((transaction, idx) => {
                      return (
                        <tr>
                          <td>{transaction.name}</td>
                          <td>{transaction.date}</td>
                          <td>{transaction.amount * -1}</td>
                        </tr>
                      )
                    })}
                  </tbody>
                </Table>
              </CardBody>
            </Card>
          </Col>

          {/* <Col xs={12}>
            <Card className="card-plain">
              <CardHeader>
                <CardTitle tag="h4">Table on Plain Background</CardTitle>
                <p className="card-category"> Here is a subtitle for this table</p>
              </CardHeader>
              <CardBody>
                <Table responsive>
                  <thead className="text-primary">
                    <tr>
                      {thead.map((prop, key) => {
                        if (key === thead.length - 1)
                          return (
                            <th key={key} className="text-right">
                              {prop}
                            </th>
                          );
                        return <th key={key}>{prop}</th>;
                      })}
                    </tr>
                  </thead>
                  <tbody>
                    {tbody.map((prop, key) => {
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
                    })}
                  </tbody>
                </Table>
              </CardBody>
            </Card>
          </Col> */}
        </Row>
      </div>
    );
  }
}

export default RegularTables;
