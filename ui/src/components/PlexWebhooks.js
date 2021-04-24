import { Box, Paper } from "@material-ui/core";
import { PlexWebhookTable } from "./PlexWebhookTable";

export function PlexWebhooks() {
  return (
    <Box
      backgroundColor="lightgray"
      component="main"
      display="flex"
      flexDirection="column"
      flexGrow={1}
      p={2}
    >
      <Paper elevation={4} sx={{ padding: "16px" }}>
        <PlexWebhookTable />
      </Paper>
    </Box>
  );
}
