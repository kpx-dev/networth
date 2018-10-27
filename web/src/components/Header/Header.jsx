import React from "react";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import {
  Collapse,
  Navbar,
  NavbarToggler,
  NavbarBrand,
  Nav,
  NavItem,
  Container,
  Badge,
  Tooltip,
} from "reactstrap";
import dashboardRoutes from "routes/dashboard.jsx";
import { Auth } from 'aws-amplify';
import Connect from "components/Connect/Connect.jsx";
import { get } from "../../helpers/helpers.js";
import * as DB from '../../db.js';

class Header extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isOpen: false,
      dropdownOpen: false,
      color: "transparent",
      networth: 0
    };
    this.toggle = this.toggle.bind(this);
    this.dropdownToggle = this.dropdownToggle.bind(this);
  }

  static propTypes = {
    networth: PropTypes.number,
  };

  static defaultProps = {
    networth: 0,
  };

  async handleLogout() {
    await Auth.signOut();

    // TODO: make redirect work
    window.location.reload();
  };

  toggle() {
    if (this.state.isOpen) {
      this.setState({
        color: "transparent"
      });
    } else {
      this.setState({
        color: "dark"
      });
    }
    this.setState({
      isOpen: !this.state.isOpen
    });
  }
  dropdownToggle(e) {
    this.setState({
      dropdownOpen: !this.state.dropdownOpen
    });
  }

  openSidebar() {
    document.documentElement.classList.toggle("nav-open");
    this.refs.sidebarToggle.classList.toggle("toggled");
  }

  // function that adds color dark/transparent to the navbar on resize (this is for the collapse)
  updateColor() {
    if (window.innerWidth < 993 && this.state.isOpen) {
      this.setState({
        color: "dark"
      });
    } else {
      this.setState({
        color: "transparent"
      });
    }
  }

  async componentDidMount() {
    window.addEventListener("resize", this.updateColor.bind(this));
    const userInfo = await Auth.currentUserInfo();
    const username = userInfo.username;

    this.db = await DB.get();
    const query = await this.db.networth.findOne().where("username").eq(username);
    const queryRes = await query.exec();
    const nowTS = Date.now();

    // if cached for more than 1 min, get new one from api
    if (queryRes === null || queryRes.networth === 0 || (nowTS - new Date(queryRes.updated_at)) / 1000 > 60 ) {
      const networthRes = await get(`/networth`);
      const nwBody = await networthRes.json();
      const dbBody = {
        username,
        networth: nwBody.data.networth,
        updated_at: new Date().toISOString(),
      };

      if (queryRes === null) {
        const insertRes = await this.db.networth.insert(dbBody);
        console.log(insertRes);
      } else {
        const upsertRes = await query.update({$set: dbBody});
        console.log("upsertRes ", upsertRes);
      }
      this.setState({ networth: nwBody.data.networth });
    } else {
      this.setState({ networth: queryRes.networth });
    }
  }

  componentDidUpdate(e) {
    if (
      window.innerWidth < 993 &&
      e.history.location.pathname !== e.location.pathname &&
      document.documentElement.className.indexOf("nav-open") !== -1
    ) {
      document.documentElement.classList.toggle("nav-open");
      this.refs.sidebarToggle.classList.toggle("toggled");
    }
  }

  render() {
    return (
      // add or remove classes depending if we are on full-screen-maps page or not
      <Navbar
        color={
          this.props.location.pathname.indexOf("full-screen-maps") !== -1
            ? "dark"
            : this.state.color
        }
        expand="lg"
        className={
          this.props.location.pathname.indexOf("full-screen-maps") !== -1
            ? "navbar-absolute fixed-top"
            : "navbar-absolute fixed-top " +
              (this.state.color === "transparent" ? "navbar-transparent " : "")
        }
      >
        <Container fluid>
          <div className="navbar-wrapper">
            <div className="navbar-toggle">
              <button
                type="button"
                ref="sidebarToggle"
                className="navbar-toggler"
                onClick={() => this.openSidebar()}
              >
                <span className="navbar-toggler-bar bar1" />
                <span className="navbar-toggler-bar bar2" />
                <span className="navbar-toggler-bar bar3" />
              </button>
            </div>
            <NavbarBrand>
              <Badge color="dark" pill title="Net Worth">{
                new Intl.NumberFormat('en-US', {
                  style: 'currency',
                  currency: 'USD',
                  minimumFractionDigits: 0,
                  maximumFractionDigits: 0,
                }).format(this.state.networth)}</Badge>

            </NavbarBrand>
          </div>
          <NavbarToggler onClick={this.toggle}>
            <span className="navbar-toggler-bar navbar-kebab" />
            <span className="navbar-toggler-bar navbar-kebab" />
            <span className="navbar-toggler-bar navbar-kebab" />
          </NavbarToggler>
          <Collapse
            isOpen={this.state.isOpen}
            navbar
            className="justify-content-end"
          >
            <Nav navbar>
              <NavItem>
                <Connect></Connect>
              </NavItem>
              <NavItem>
                <Link to="/logout" onClick={this.handleLogout} className="nav-link btn-rotate">
                  <i className="nc-icon nc-lock-circle-open" title="Logout"/>
                  <p>
                    <span className="d-lg-none d-md-block">Logout</span>
                  </p>
                </Link>
              </NavItem>
            </Nav>
          </Collapse>
        </Container>
      </Navbar>
    );
  }
}

export default Header;
