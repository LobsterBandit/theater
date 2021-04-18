import { Box, Paper } from "@material-ui/core";
import { PlexWebhookTable } from "./PlexWebhookTable";
import { PlexWebhookToolbar } from "./PlexWebhookToolbar";
import { usePlexWebhooks } from "../hooks/usePlexWebhooks";

export function PlexWebhooks() {
  const {
    state: { loading, plexWebhooks, total },
    fetchPlexWebhooks,
    options,
  } = usePlexWebhooks();

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
        <PlexWebhookToolbar
          loading={loading}
          onRefreshClick={() => fetchPlexWebhooks(options)}
        />
        <PlexWebhookTable
          data={plexWebhooks}
          fetchData={fetchPlexWebhooks}
          loading={loading}
          totalCount={total}
        />
      </Paper>
    </Box>
  );
}
