import Components from "../views/Components/Components.jsx";
import LandingPage from "../views/LandingPage/LandingPage.jsx";
import ProfilePage from "../views/ProfilePage/ProfilePage.jsx";
import LoginPage from "../views/LoginPage/LoginPage.jsx";
import LinkAccountPage from "../views/LinkAccountPage/LinkAccountPage.jsx";

var indexRoutes = [
  { path: "/landing-page", name: "LandingPage", component: LandingPage },
  { path: "/profile-page", name: "ProfilePage", component: ProfilePage },
  { path: "/login", name: "LoginPage", component: LoginPage },
  {
    path: "/link-account",
    name: "LinkAccountPage",
    component: LinkAccountPage
  },
  { path: "/", name: "Components", component: Components }
];

export default indexRoutes;
