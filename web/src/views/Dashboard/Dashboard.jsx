import React from "react";
import PropTypes from "prop-types";
import {
  Button, ButtonGroup,
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  Row,
  Col,
} from "reactstrap";
import { Line } from "react-chartjs-2";
import Stats from "components/Stats/Stats.jsx";
import Loader from "react-loader-spinner";
import { each } from "lodash/each";
import { get } from "../../helpers/helpers.js";

class Dashboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = { networth: 0, resolution: "monhtly" };
    this.chartCache = {};
  }

  static propTypes = {
    networth: PropTypes.number,
  };

  static defaultProps = {
    networth: 0,
  };

  async refreshNetworthChart(resolution) {
    if (this.chartCache[resolution] !== undefined) {
      this.setState({ resolution });
      this.chartData = this.chartCache[resolution];
      // console.log(this.chartData);
      return this.chartData;
    }

    this.setState({ loading: true, resolution });
    const res = await get(`/networth_history?resolution=${resolution}`);
    const nwBody = await res.json();
    this.chartData = this._generateChartData(nwBody.data);
    this.setState({ loading: false });
    this.chartCache[resolution] = this.chartData;

    return this.chartData;
  }

  async componentDidMount() {
    await this.refreshNetworthChart(this.state.resolution);
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
      data: []
    };

    const payload = {
      labels: [],
      datasets: [],
      options: {
        legend: {
          display: false,
          position: "top"
        }
      }
    };

    Object.keys(data).forEach(date => {
      const dateObj = new Date(date);
      let label = date;
      switch(this.state.resolution) {
        case "daily": label = `${dateObj.getHours()}`; break;
        case "monthly": label = `${dateObj.getMonth()}/${dateObj.getDate()}`; break;
        case "yearly": label = `${dateObj.getMonth()}`; break;
      }

      payload.labels.push(label);
      networthSet.data.push(data[date]);
    });

    payload.datasets = [ networthSet ];

    return payload;
  }

  getChartData() {
    const data = {
      data: {
        labels: [
          "Jan",
          "Feb",
          "Mar",
          "Apr",
          "May",
          "Jun",
          "Jul",
          "Aug",
          "Sep",
          "Oct",
          "Nov",
          "Dec"
        ],
        datasets: [
          // {
          //   data: [0, 19, 15, 20, 30, 40, 40, 50, 25, 30, 50, 70],
          //   fill: false,
          //   borderColor: "#fbc658",
          //   backgroundColor: "transparent",
          //   pointBorderColor: "#fbc658",
          //   pointRadius: 4,
          //   pointHoverRadius: 4,
          //   pointBorderWidth: 8
          // },
          {
            data: [0, 5, 10, 12, 20, 27, 30, 34, 42, 45, 55, 63],
            fill: false,
            borderColor: "#51CACF",
            backgroundColor: "transparent",
            pointBorderColor: "#51CACF",
            pointRadius: 4,
            pointHoverRadius: 4,
            pointBorderWidth: 8
          }
        ]
      },
      options: {
        legend: {
          display: false,
          position: "top"
        }
      }
    };

    return data;
  }

  render() {
    if (this.state.loading) return (
      <div className="content">
        <Loader type="ThreeDots" height={80} width={80} />
      </div>
    );

    return (
      <div className="content">
        <Row>
          <Col xs={12}>
          <ButtonGroup size="sm">
            <Button onClick={() => this.refreshNetworthChart('daily')} active={this.state.resolution === 'daily'}>Daily</Button>
            <Button onClick={() => this.refreshNetworthChart('monthly')} active={this.state.resolution === 'monthly'}>Monthly</Button>
            <Button onClick={() => this.refreshNetworthChart('yearly')} active={this.state.resolution === 'yearly'}>Yearly</Button>

          </ButtonGroup>

          <Card className="card-chart">
              <CardHeader>
                {/* <CardTitle>Net Worth</CardTitle> */}
                {/* <p className="card-category">Net Worth over time.</p> */}
              </CardHeader>
              <CardBody>
                <Line
                  data={this.chartData}
                  // data={this.getChartData().data}
                  // options={this.getChartData().options}
                  width={400}
                  height={150}
                />
              </CardBody>
              <CardFooter>
                {/* <div className="chart-legend">
                  <i className="fa fa-circle text-info" /> Net Worth
                  <i className="fa fa-circle text-warning" /> Liabilities
                </div> */}
                <hr />
                <Stats>
                  {[
                    {
                      i: "fas fa-check",
                      t: "Last updated: 10s ago"
                    }
                  ]}
                </Stats>
              </CardFooter>
            </Card>

            {/* <Card>
              <CardHeader>
                <CardTitle>Users Behavior</CardTitle>
                <p className="card-category">24 Hours performance</p>
              </CardHeader>
              <CardBody>
                <Line
                  data={dashboard24HoursPerformanceChart.data}
                  options={dashboard24HoursPerformanceChart.options}
                  width={400}
                  height={100}
                />
              </CardBody>
              <CardFooter>
                <hr />
                <Stats>
                  {[
                    {
                      i: "fas fa-history",
                      t: " Updated 3 minutes ago"
                    }
                  ]}
                </Stats>
              </CardFooter>
            </Card> */}
          </Col>
        </Row>

        {/* <Row>
          <Col xs={12} sm={6} md={6} lg={3}>
            <Card className="card-stats">
              <CardBody>
                <Row>
                  <Col xs={5} md={4}>
                    <div className="icon-big text-center">
                      <i className="nc-icon nc-globe text-warning" />
                    </div>
                  </Col>
                  <Col xs={7} md={8}>
                    <div className="numbers">
                      <p className="card-category">Capacity</p>
                      <CardTitle tag="p">150GB</CardTitle>
                    </div>
                  </Col>
                </Row>
              </CardBody>
              <CardFooter>
                <hr />
                <Stats>
                  {[
                    {
                      i: "fas fa-sync-alt",
                      t: "Update Now"
                    }
                  ]}
                </Stats>
              </CardFooter>
            </Card>
          </Col>
          <Col xs={12} sm={6} md={6} lg={3}>
            <Card className="card-stats">
              <CardBody>
                <Row>
                  <Col xs={5} md={4}>
                    <div className="icon-big text-center">
                      <i className="nc-icon nc-money-coins text-success" />
                    </div>
                  </Col>
                  <Col xs={7} md={8}>
                    <div className="numbers">
                      <p className="card-category">Revenue</p>
                      <CardTitle tag="p">$ 1,345</CardTitle>
                    </div>
                  </Col>
                </Row>
              </CardBody>
              <CardFooter>
                <hr />
                <Stats>
                  {[
                    {
                      i: "far fa-calendar",
                      t: "Last day"
                    }
                  ]}
                </Stats>
              </CardFooter>
            </Card>
          </Col>
          <Col xs={12} sm={6} md={6} lg={3}>
            <Card className="card-stats">
              <CardBody>
                <Row>
                  <Col xs={5} md={4}>
                    <div className="icon-big text-center">
                      <i className="nc-icon nc-vector text-danger" />
                    </div>
                  </Col>
                  <Col xs={7} md={8}>
                    <div className="numbers">
                      <p className="card-category">Errors</p>
                      <CardTitle tag="p">23</CardTitle>
                    </div>
                  </Col>
                </Row>
              </CardBody>
              <CardFooter>
                <hr />
                <Stats>
                  {[
                    {
                      i: "far fa-clock",
                      t: "In the last hour"
                    }
                  ]}
                </Stats>
              </CardFooter>
            </Card>
          </Col>
          <Col xs={12} sm={6} md={6} lg={3}>
            <Card className="card-stats">
              <CardBody>
                <Row>
                  <Col xs={5} md={4}>
                    <div className="icon-big text-center">
                      <i className="nc-icon nc-favourite-28 text-primary" />
                    </div>
                  </Col>
                  <Col xs={7} md={8}>
                    <div className="numbers">
                      <p className="card-category">Followers</p>
                      <CardTitle tag="p">+45K</CardTitle>
                    </div>
                  </Col>
                </Row>
              </CardBody>
              <CardFooter>
                <hr />
                <Stats>
                  {[
                    {
                      i: "fas fa-sync-alt",
                      t: "Update now"
                    }
                  ]}
                </Stats>
              </CardFooter>
            </Card>
          </Col>
        </Row> */}


        {/* <Row>
          <Col xs={12} sm={12} md={4}>
            <Card>
              <CardHeader>
                <CardTitle>Email Statistics</CardTitle>
                <p className="card-category">Last Campaign Performance</p>
              </CardHeader>
              <CardBody>
                <Pie
                  data={dashboardEmailStatisticsChart.data}
                  options={dashboardEmailStatisticsChart.options}
                />
              </CardBody>
              <CardFooter>
                <div className="legend">
                  <i className="fa fa-circle text-primary" /> Opened{" "}
                  <i className="fa fa-circle text-warning" /> Read{" "}
                  <i className="fa fa-circle text-danger" /> Deleted{" "}
                  <i className="fa fa-circle text-gray" /> Unopened
                </div>
                <hr />
                <Stats>
                  {[
                    {
                      i: "fas fa-calendar-alt",
                      t: " Number of emails sent"
                    }
                  ]}
                </Stats>
              </CardFooter>
            </Card>
          </Col>
          <Col xs={12} sm={12} md={8}>

          </Col>
        </Row> */}
      </div>
    );
  }
}

export default Dashboard;
