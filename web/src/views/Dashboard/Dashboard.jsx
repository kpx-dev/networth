import React from "react";
import {
  Button, ButtonGroup,
  Card,
  CardBody,
  CardFooter,
  Row,
  Col,
} from "reactstrap";
import { Line } from "react-chartjs-2";
import Stats from "components/Stats/Stats.jsx";
import Loader from "react-loader-spinner";
import { get } from "../../helpers/helpers.js";
import { Auth } from 'aws-amplify';
import * as DB from '../../db.js';
import NotificationAlert from "react-notification-alert";
import Connect from "components/Connect/Connect.jsx";

class Dashboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = { resolution: "daily" };
    this.chartCache = {};
    this.chartData = {};
    this.chartOptions = {
      legend: {
        display: false,
      }
    };
  }

  async refreshNetworthChart(resolution) {
    if (this.chartCache[resolution] !== undefined) {
      this.setState({ resolution });
      this.chartData = this.chartCache[resolution];
      return this.chartData;
    }

    this.setState({ loading: true, resolution });

    try {
      const res = await get(`/networth_history?resolution=${resolution}`);
      const nwBody = await res.json();
      this.chartData = this._generateChartData(nwBody.data);
      this.setState({ loading: false });
      this.chartCache[resolution] = this.chartData;
    } catch (e) {
      this.alert('Cannot connect to REST API.');
      this.setState({ loading: false });
    }

    return this.chartData;
  }

  async componentDidMount() {
    await this.refreshNetworthChart(this.state.resolution);

    const userInfo = await Auth.currentUserInfo();
    const username = userInfo.username;
    this.db = await DB.get();
    const networth = await this.db.networth.findOne().where("username").eq(username).exec();
    if (networth && networth.updated_at) {
      const networthUpdatedAt = new Date(networth.updated_at).toDateString();
      this.setState({ networthUpdatedAt });
    }
    await this.getTokenError();
  }

  async getTokenError() {
    try {
      let res = await get(`/tokens/error`);
      const tokenErrors = await res.json();
      const badTokens = tokenErrors.data;
      if (badTokens && badTokens.length > 0) {
        const badToken = badTokens[0];
        res = await get(`/tokens/public?item_id=${badToken.item_id}`);
        const publicTokenData = await res.json();
        const publicToken = publicTokenData.data.public_token;

        this.alertRelinkAccount(`Accounts from ${badToken.institution_name} need re-linking again.`, publicToken);
      }
    } catch (e) {
      console.error(e);
    }
  }

  _generateChartData(data) {
    const networthSet = {
      label: 'Net Worth',
      fill: false,
      borderColor: "#51CACF",
      backgroundColor: "transparent",
      pointBorderColor: "#51CACF",
      pointRadius: 4,
      pointHoverRadius: 4,
      pointBorderWidth: 8,
      data: [],
    };

    const payload = {
      labels: [],
      datasets: [],
    };

    Object.keys(data).forEach(date => {
      const dateObj = new Date(date);
      let label = date;
      switch(this.state.resolution) {
        case "daily": label = `${dateObj.getMonth()}/${dateObj.getDate()} ${dateObj.toLocaleString('en-US', { hour: 'numeric', hour12: true })}`; break;
        case "monthly": label = `${dateObj.getMonth()}/${dateObj.getDate()}`; break;
        case "yearly": label = `${dateObj.getMonth()}/${dateObj.getFullYear()}`; break;
        default: label = date;
      }

      payload.labels.push(label);
      networthSet.data.push(data[date]);
    });

    payload.datasets = [ networthSet ];

    return payload;
  }

  alertRelinkAccount(message, publicToken) {
    const options = {
      place: "tc",
      type: "danger",
      icon: "nc-icon nc-alert-circle-i",
      message: (
        <div>
          <div>
            {message} <Connect text="Fix Now" token={publicToken}></Connect>
          </div>
        </div>
      ),
    };
    this.refs.alert.notificationAlert(options);
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
    if (this.state.loading) return (
      <div className="content">
        <NotificationAlert ref="alert" />

        <Loader type="ThreeDots" height={80} width={80} />
      </div>
    );

    return (
      <div className="content">
        <NotificationAlert ref="alert" />
        <br />
        <Row>
          <Col xs={12}>
            <ButtonGroup size="sm" className="pull-right">
              <Button onClick={() => this.refreshNetworthChart('daily')} active={this.state.resolution === 'daily'}>Daily</Button>
              <Button onClick={() => this.refreshNetworthChart('monthly')} active={this.state.resolution === 'monthly'}>Monthly</Button>
              <Button onClick={() => this.refreshNetworthChart('yearly')} active={this.state.resolution === 'yearly'}>Yearly</Button>
            </ButtonGroup>
          </Col>
        </Row>

        <Row>
          <Col xs={12}>
          <Card className="card-chart">
              <CardBody>
                <Line
                  data={this.chartData}
                  options={this.chartOptions}
                  width={400}
                  height={150}
                />
              </CardBody>
              <CardFooter>
                <hr />
                <Stats>
                  {[
                    {
                      i: "fas fa-check",
                      t: `Last updated: ${this.state.networthUpdatedAt}`
                    }
                  ]}
                </Stats>
              </CardFooter>
            </Card>
          </Col>
        </Row>
      </div>
    );
  }
}

export default Dashboard;
