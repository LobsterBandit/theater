import { Box } from "@material-ui/core";
import { Route, Switch } from "react-router-dom";
import { Header, PlexWebhooks } from "./components";

export function App() {
  return (
    <Box display="flex" flexDirection="column" height="100vh">
      <Header />
      <Switch>
        <Route exact path="/">
          <PlexWebhooks />
        </Route>
      </Switch>
    </Box>
  );
}
