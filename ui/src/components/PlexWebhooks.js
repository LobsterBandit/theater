import { Box, Paper } from "@material-ui/core";
import { PlexWebhookTable } from "./PlexWebhookTable";
import { PlexWebhookToolbar } from "./PlexWebhookToolbar";
import { usePlexWebhooks } from "../hooks/usePlexWebhooks";

export function PlexWebhooks() {
  const [{ plexWebhooks }, refetch] = usePlexWebhooks();

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
        <PlexWebhookToolbar onRefreshClick={refetch} />
        <PlexWebhookTable data={plexWebhooks} />
      </Paper>
    </Box>
  );
}
