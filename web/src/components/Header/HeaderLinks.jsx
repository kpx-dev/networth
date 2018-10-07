import React from "react";
import withStyles from "@material-ui/core/styles/withStyles";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import { LockOpen, Add } from "@material-ui/icons";
import Button from "../CustomButtons/Button.jsx";
import headerLinksStyle from "../../assets/jss/material-kit-react/components/headerLinksStyle.jsx";
import { Auth } from 'aws-amplify';
import LinkAccount from "../LinkAccount/LinkAccount.jsx";

function HeaderLinks({ ...props }) {
  const { classes } = props;

  const handleLogout = async () => {
    await Auth.signOut();

    // TODO: make redirect work
    window.location.reload();
  };

  return (
    <List className={classes.list}>
      <ListItem className={classes.listItem}>
        <LinkAccount text="Connect" />
      </ListItem>
      <ListItem className={classes.listItem}>
        <Button
          color="transparent"
          className={classes.navLink}
          onClick={handleLogout}
        >
        <LockOpen className={classes.icons} /> Logout
        </Button>
      </ListItem>
    </List>
  );
}

export default withStyles(headerLinksStyle)(HeaderLinks);
