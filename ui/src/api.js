export async function replayPlexWebhook(payload) {
  const formData = new FormData();

  formData.append("payload", JSON.stringify(payload));

  const response = await fetch("/plex", {
    method: "POST",
    body: formData,
    headers: { "X-Request-Type": "replay" },
  });

  if (!response.ok) {
    throw new Error(response.statusText);
  }

  return `Successfully replayed Plex ${payload.event} webhook`;
}
