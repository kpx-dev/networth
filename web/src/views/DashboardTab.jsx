import React from "react";
import PropTypes from "prop-types";
import withStyles from "@material-ui/core/styles/withStyles";
import Dashboard from "@material-ui/icons/Dashboard";
import Schedule from "@material-ui/icons/Schedule";
import AddIcon from "@material-ui/icons/Add";
// import List from "@material-ui/icons/List";
import GridContainer from "../components/Grid/GridContainer.jsx";
import GridItem from "../components/Grid/GridItem.jsx";
import NavPills from "../components/NavPills/NavPills.jsx";
import pillsStyle from "../assets/jss/material-kit-react/views/componentsSections/pillsStyle.jsx";
import LinkAccount from "../components/LinkAccount/LinkAccount.jsx";
import { get } from "../helpers/helpers.js";

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
          <div id="navigation-pills">
            <GridContainer>
              <GridItem xs={12} sm={12} md={12} lg={12}>
                <NavPills
                  color="rose"
                  horizontal={{
                    tabsGrid: { xs: 12, sm: 4, md: 4, lg: 2 },
                    contentGrid: { xs: 12, sm: 8, md: 8 }
                  }}
                  tabs={[
                    {
                      tabButton: "Net Worth",
                      tabIcon: Dashboard,
                      tabContent: (
                        <span>
                          <h1>${this.state.networth}</h1>
                        </span>
                      )
                    },

                    {
                      tabButton: "Transactions",
                      tabIcon: Schedule,
                      tabContent: (
                        <span>
                          <p>
                            Efficiently unleash cross-media information without
                            cross-media value. Quickly maximize timely
                            deliverables for real-time schemas.
                          </p>
                        </span>
                      )
                    },

                    {
                      tabButton: "Link Account",
                      tabIcon: AddIcon,
                      tabContent: (
                        <span>
                          <LinkAccount
                            institution="ins_1"
                            text="Bank of America"
                          />
                          <LinkAccount />
                          {/* <p>Common Banks</p>

                          <LinkAccount institution="ins_3" text="Chase" />
                          <LinkAccount institution="ins_4" text="Wells Fargo" /> */}

                          {/* <LinkAccount institution="ins_5" text="Citi" />
                          <LinkAccount institution="ins_14" text="TD Bank" />
                          <LinkAccount institution="ins_13" text="PNC" />
                          <LinkAccount institution="ins_9" text="Capital One" />
                          <LinkAccount institution="ins_22" text="BB&T" />
                          <LinkAccount institution="ins_16" text="SunTrust" />
                          <LinkAccount institution="ins_25" text="Ally Bank" /> */}
                          {/* <LinkAccount
                            institution="ins_10"
                            text="American Express"
                          />
                          <br />
                          <br /> */}
                        </span>
                      )
                    },


                  ]}
                />
              </GridItem>
            </GridContainer>
          </div>
        </div>
      </div>
    );
  }
}

export default withStyles(pillsStyle)(DashboardTab);
