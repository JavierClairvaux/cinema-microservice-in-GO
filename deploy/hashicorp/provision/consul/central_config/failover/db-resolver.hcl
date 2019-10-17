Kind = "service-resolver"
Name = "postgressdb1"

failover = {
  "*" = {
    datacenters = ["dc1"]
  }
}
