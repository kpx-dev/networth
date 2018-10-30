import Dashboard from "views/Dashboard/Dashboard.jsx";
import Accounts from "views/Accounts/Accounts.jsx";
import Transactions from "views/Transactions/Transactions.jsx";
import Connect from "views/Connect/Connect.jsx";

var dashRoutes = [
  {
    path: "/app/networth",
    name: "Net Worth",
    icon: "nc-icon nc-money-coins",
    component: Dashboard
  },
  {
    path: "/app/accounts",
    name: "Accounts",
    icon: "nc-icon nc-bank",
    component: Accounts
  },
  {
    path: "/app/transactions",
    name: "Transactions",
    icon: "nc-icon nc-tile-56",
    component: Transactions
  },
  {
    path: "/app/connect",
    name: "Connect New Bank",
    icon: "nc-icon nc-simple-add",
    component: Connect
  },
  { redirect: true,
    path: "/",
    pathTo: "/app/networth",
    name: "Net Worth"
  },
  { redirect: true,
    path: "/app",
    pathTo: "/app/networth",
    name: "Net Worth"
  }
];
export default dashRoutes;
