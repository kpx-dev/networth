import { Auth } from "aws-amplify";
import { NW_API_BASE_URL } from "./constants.js";

async function post(path, body) {
  const session = await Auth.currentSession();
  const fetchOptions = {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      Authorization: session.idToken.jwtToken
    },
    body: JSON.stringify(body)
  };

  let absPath = path;
  console.log('REACT_APP_NW_API_HOST ', NW_API_BASE_URL, process.env.NODE_ENV);
  if (NW_API_BASE_URL) absPath = `${NW_API_BASE_URL}${path}`;

  return fetch(absPath, fetchOptions);
}

export { post };
