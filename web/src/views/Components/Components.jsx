import React from "react";
import classNames from "classnames";
// import { Link } from "react-router-dom";
import withStyles from "@material-ui/core/styles/withStyles";
import Header from "../../components/Header/Header.jsx";
import Footer from "../../components/Footer/Footer.jsx";
// import GridContainer from "components/Grid/GridContainer.jsx";
// import GridItem from "components/Grid/GridItem.jsx";
// import Button from "components/CustomButtons/Button.jsx";
// import Parallax from "components/Parallax/Parallax.jsx";
import HeaderLinks from "../../components/Header/HeaderLinks.jsx";
// import SectionBasics from "./Sections/SectionBasics.jsx";
// import SectionNavbars from "./Sections/SectionNavbars.jsx";
// import SectionTabs from "./Sections/SectionTabs.jsx";
import SectionPills from "./Sections/SectionPills.jsx";
// import SectionNotifications from "./Sections/SectionNotifications.jsx";
// import SectionTypography from "./Sections/SectionTypography.jsx";
// import SectionJavascript from "./Sections/SectionJavascript.jsx";
// import SectionCarousel from "./Sections/SectionCarousel.jsx";
// import SectionCompletedExamples from "./Sections/SectionCompletedExamples.jsx";
// import SectionLogin from "./Sections/SectionLogin.jsx";
// import SectionExamples from "./Sections/SectionExamples.jsx";
// import SectionDownload from "./Sections/SectionDownload.jsx";
import componentsStyle from "../../assets/jss/material-kit-react/views/components.jsx";

class Components extends React.Component {

  render() {
    const { classes, ...rest } = this.props;
    return (
      <div>
        <Header
          brand="networth.app"
          rightLinks={<HeaderLinks />}
          fixed
          color="transparent"
          changeColorOnScroll={{
            height: 400,
            color: "white"
          }}
          {...rest}
        />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />

        {/* <Parallax image={require("assets/img/bg4.jpg")}>
          <div className={classes.container}>
            <GridContainer>
              <GridItem>
                <div className={classes.brand}>
                  <h1 className={classes.title}>networth.app</h1>
                  <h3 className={classes.subtitle}>
                    A Badass Material-UI Kit based on Material Design.
                  </h3>
                </div>
              </GridItem>
            </GridContainer>
          </div>
        </Parallax> */}

        <div className={classNames(classes.main, classes.mainRaised)}>
          {/* <SectionBasics /> */}
          <SectionPills />

          {/* <SectionTabs /> */}
          {/* <SectionDownload /> */}

          {/* <SectionNavbars /> */}
          {/* <SectionNotifications /> */}
          {/* <SectionTypography />
          <SectionJavascript />
          <SectionCarousel />
          <SectionCompletedExamples />
          <SectionLogin /> */}
          {/* <GridItem md={12} className={classes.textCenter}>
            <Link to={"/login"} className={classes.link}>
              <Button color="primary" size="lg" simple>
                View Login Page
              </Button>
            </Link>
          </GridItem> */}
          {/* <SectionExamples /> */}
        </div>
        <Footer />
      </div>
    );
  }
}

export default withStyles(componentsStyle)(Components);
