import React from "react";
import PropTypes from "prop-types";
import {
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  // CardTitle,
  Row,
  Col
} from "reactstrap";
import { Line } from "react-chartjs-2";
import Stats from "components/Stats/Stats.jsx";
import Loader from 'react-loader-spinner';
import { each } from 'lodash/each';
import { get } from "../../helpers/helpers.js";

class Dashboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {networth: 0};
  }

  static propTypes = {
    networth: PropTypes.number,
  };

  static defaultProps = {
    networth: 0,
  };

  async componentDidMount() {
    const startDate = '2018-10-01';
    const endDate = '2018-11-01';
    this.setState({ loading: true });
    const nw = await get('/networth');
    const nwHistory = await get(`/networth?start_date=${startDate}&end_date=${endDate}`);
    const body = await nw.json();
    const nwHistoryBody = await nwHistory.json();
    this.chartData = this._generateChartData(nwHistoryBody.data);
    this.setState({
      networth: body.data,
      loading: false
    });
  }

  _generateChartData(data) {
    const networthSet = {
      label: 'Net Worth',
      fill: false,
      backgroundColor: 'rgba(75,192,192,0.4)',
      borderColor: 'rgba(75,192,192,1)',
      pointBorderColor: 'rgba(75,192,192,1)',
      pointBackgroundColor: '#fff',
      pointBorderWidth: 5,
      pointHoverRadius: 10,
      pointHoverBackgroundColor: 'rgba(75,192,192,1)',
      pointHoverBorderColor: 'rgba(220,220,220,1)',
      pointHoverBorderWidth: 2,
      pointRadius: 1,
      pointHitRadius: 10,
      data: [] };
    const assetsSet = {
      label: 'Assets',
      fill: false,
      backgroundColor: 'rgba(75,192,192,0.4)',
      borderColor: 'blue',
      pointBorderColor: 'blue',
      pointBackgroundColor: '#fff',
      pointBorderWidth: 5,
      pointHoverRadius: 10,
      pointHoverBackgroundColor: 'rgba(75,192,192,1)',
      pointHoverBorderColor: 'rgba(220,220,220,1)',
      pointHoverBorderWidth: 2,
      pointRadius: 1,
      pointHitRadius: 10,
      data: [] };
    const liabilitiesSet = {
      label: 'Liabilities',
      fill: false,
      backgroundColor: 'rgba(75,192,192,0.4)',
      borderColor: 'red',
      pointBorderColor: 'red',
      pointBackgroundColor: '#fff',
      pointBorderWidth: 5,
      pointHoverRadius: 10,
      pointHoverBackgroundColor: 'rgba(75,192,192,1)',
      pointHoverBorderColor: 'rgba(220,220,220,1)',
      pointHoverBorderWidth: 2,
      pointRadius: 1,
      pointHitRadius: 10,
      data: [] };
    const payload = { labels: [], datasets: [] };
    data.forEach(item => {
      payload.labels.push(item.sort);
      networthSet.data.push(item.networth);
      assetsSet.data.push(item.assets);
      liabilitiesSet.data.push(item.liabilities);
    });

    payload.datasets = [ networthSet, assetsSet, liabilitiesSet ];

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
          <Card className="card-chart">
              <CardHeader>
                {/* <CardTitle>Net Worth</CardTitle> */}
                {/* <p className="card-category">Net Worth over time.</p> */}
              </CardHeader>
              <CardBody>
                <Line
                  data={this.getChartData().data}
                  options={this.getChartData().options}
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
