import Dashboard from "views/Dashboard/Dashboard.jsx";
import Accounts from "views/Accounts/Accounts.jsx";
import Transactions from "views/Transactions/Transactions.jsx";
import Connect from "views/Connect/Connect.jsx";

var dashRoutes = [
  {
    path: "/networth",
    name: "Net Worth",
    icon: "nc-icon nc-money-coins",
    component: Dashboard
  },
  {
    path: "/accounts",
    name: "Accounts",
    icon: "nc-icon nc-bank",
    component: Accounts
  },
  {
    path: "/transactions",
    name: "Transactions",
    icon: "nc-icon nc-tile-56",
    component: Transactions
  },
  {
    path: "/connect",
    name: "Connect New Bank",
    icon: "nc-icon nc-simple-add",
    component: Connect
  },
  { redirect: true, path: "/", pathTo: "/networth", name: "Net Worth" }
];
export default dashRoutes;
