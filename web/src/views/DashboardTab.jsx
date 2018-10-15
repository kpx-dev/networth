import React from "react";
import PropTypes from "prop-types";
import withStyles from "@material-ui/core/styles/withStyles";
import Dashboard from "@material-ui/icons/Dashboard";
import Schedule from "@material-ui/icons/Schedule";
import AddIcon from "@material-ui/icons/Add";
import GridContainer from "../components/Grid/GridContainer.jsx";
import GridItem from "../components/Grid/GridItem.jsx";
import NavPills from "../components/NavPills/NavPills.jsx";
import pillsStyle from "../assets/jss/material-kit-react/views/componentsSections/pillsStyle.jsx";
import LinkAccount from "../components/LinkAccount/LinkAccount.jsx";
import { get } from "../helpers/helpers.js";
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Loader from 'react-loader-spinner'
import {Line} from 'react-chartjs-2';
import { each } from 'lodash/each';

class DashboardTab extends React.Component {
  constructor(props) {
    super(props);
    this.state = {networth: 0};
  }

  static propTypes = {
    classes: PropTypes.object.isRequired,
    networth: PropTypes.number,
  };

  static defaultProps = {
    networth: 0,
  };

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

  render() {
    const { classes } = this.props;
    if (this.state.loading) return <Loader type="ThreeDots" color="#somecolor" height={80} width={80} />

    return (
      <div className={classes.section}>
        <div className={classes.container}>
        <GridContainer>
          <GridItem xs={4} sm={4} md={4}>
          <Card className={classes.card}>
          <CardContent>
          <h2>
            {new Intl.NumberFormat('en-us', {
                style: 'currency',
                currency: 'USD',
                minimumFractionDigits: 0,
            }).format(this.state.networth.networth)}
          </h2>
          <h3>Net Worth</h3>
          </CardContent>
          </Card>
          </GridItem>

          <GridItem xs={4} sm={4} md={4}>
          <Card className={classes.card}>
          <CardContent>
          <h2>
            {new Intl.NumberFormat('en-us', {
                style: 'currency',
                currency: 'USD',
                minimumFractionDigits: 0,
            }).format(this.state.networth.assets)}
          </h2>
          <h3>Assets</h3>
          </CardContent>
          </Card>
          </GridItem>

          <GridItem xs={4} sm={4} md={4}>
          <Card className={classes.card}>
          <CardContent>
          <h2>
            {new Intl.NumberFormat('en-us', {
                style: 'currency',
                currency: 'USD',
                minimumFractionDigits: 0,
            }).format(this.state.networth.liabilities)}
          </h2>
          <h3>Liabilities</h3>
          </CardContent>
          </Card>
          </GridItem>
        </GridContainer>

        <GridContainer>
          <GridItem xs={12} sm={12} md={12}>

          <Line data={this.chartData} width="600" height="250" />

          </GridItem>
        </GridContainer>
        </div>
      </div>
    );
  }
}

export default withStyles(pillsStyle)(DashboardTab);
