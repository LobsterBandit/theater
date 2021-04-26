import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Typography,
} from "@material-ui/core";
import Close from "@material-ui/icons/Close";

export function WebhookPayloadDialog({
  handleClose,
  handleReplay,
  open,
  value,
}) {
  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogTitle
        disableTypography
        sx={{
          alignItems: "center",
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-between",
        }}
      >
        <Typography variant="h6">Plex Webhook Payload</Typography>
        <IconButton onClick={handleClose}>
          <Close />
        </IconButton>
      </DialogTitle>
      <DialogContent
        style={{
          fontSize: "12px",
          whiteSpace: "pre-wrap",
          wordWrap: "break-word",
        }}
      >
        {JSON.stringify(value, null, 2)}
      </DialogContent>
      <DialogActions>
        <Button
          color="primary"
          onClick={(e) => {
            e.stopPropagation();
            handleReplay(e, value.payload);
          }}
          variant="contained"
        >
          Replay Event
        </Button>
      </DialogActions>
    </Dialog>
  );
}
