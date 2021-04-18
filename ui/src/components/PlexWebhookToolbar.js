import { IconButton, Toolbar, Tooltip, Typography } from "@material-ui/core";
import RefreshIcon from "@material-ui/icons/Refresh";

export function PlexWebhookToolbar({ loading, onRefreshClick }) {
  return (
    <Toolbar disableGutters={true} variant="dense">
      <Typography flexGrow={1} variant="h6">
        Plex Webhooks
      </Typography>
      {loading && <Typography>Loading...</Typography>}
      <Tooltip title="Refresh">
        <IconButton
          onClick={onRefreshClick}
          color="inherit"
          aria-label="refresh"
          sx={{ ml: 2 }}
        >
          <RefreshIcon />
        </IconButton>
      </Tooltip>
    </Toolbar>
  );
}
