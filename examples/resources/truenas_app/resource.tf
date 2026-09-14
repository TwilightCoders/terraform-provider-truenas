resource "truenas_app" "plex" {
  name    = "plex"
  include = ["/mnt/tank/apps/docker/plex.yml"]

  portals = {
    "Web UI" = "http://plex.lan/"
    "Direct" = "http://192.168.1.10:32400/web"
  }
  notes = "Library at /mnt/tank/media."

  # Redeploy when the included compose file changes.
  redeploy_trigger = filesha256("${path.module}/stacks/plex.yml")
}
