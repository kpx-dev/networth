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

  async componentDidMount() {
    const nw = await get('/networth');
    const body = await nw.json();
    this.setState({
      networth: body.data
    });
  }

  render() {
    const { classes } = this.props;

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
                currency: 'USD'
            }).format(this.state.networth)}
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
                currency: 'USD'
            }).format(this.state.networth)}
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
                currency: 'USD'
            }).format(this.state.networth)}
          </h2>
          <h3>Liabilities</h3>
          </CardContent>
          </Card>
          </GridItem>
        </GridContainer>

        </div>
      </div>
    );
  }
}

export default withStyles(pillsStyle)(DashboardTab);
