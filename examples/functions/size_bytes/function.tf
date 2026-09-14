output "quota" {
  value = provider::truenas::size_bytes("1T") # 1099511627776
}
