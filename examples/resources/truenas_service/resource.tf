resource "truenas_service" "smb" {
  service = "cifs"
  enable  = true
}
