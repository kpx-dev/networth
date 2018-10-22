import Dashboard from "views/Dashboard/Dashboard.jsx";
import TableList from "views/TableList/TableList.jsx";

var dashRoutes = [
  {
    path: "/networth",
    name: "Net Worth",
    icon: "nc-icon nc-money-coins",
    component: Dashboard
  },
  // {
  //   path: "/user-page",
  //   name: "User Profile",
  //   icon: "nc-icon nc-single-02",
  //   component: UserPage
  // },
  {
    path: "/accounts",
    name: "Accounts",
    icon: "nc-icon nc-bank",
    component: TableList
  },
  {
    path: "/transactions",
    name: "Transactions",
    icon: "nc-icon nc-tile-56",
    component: TableList
  },
  // {
  //   path: "/icons",
  //   name: "Icons",
  //   icon: "nc-icon nc-diamond",
  //   component: Icons
  // },
  // {
  //   path: "/typography",
  //   name: "Typography",
  //   icon: "nc-icon nc-caps-small",
  //   component: Typography
  // },
  { redirect: true, path: "/", pathTo: "/networth", name: "Net Worth" }
];
export default dashRoutes;
