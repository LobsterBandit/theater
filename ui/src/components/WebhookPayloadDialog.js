import { Dialog, DialogContent, DialogTitle } from "@material-ui/core";

export function WebhookPayloadDialog({ handleClose, open, value }) {
  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogTitle>Plex Webhook Payload</DialogTitle>
      <DialogContent
        style={{
          fontSize: "12px",
          whiteSpace: "pre-wrap",
          wordWrap: "break-word",
        }}
      >
        {JSON.stringify(value, null, 2)}
      </DialogContent>
    </Dialog>
  );
}
