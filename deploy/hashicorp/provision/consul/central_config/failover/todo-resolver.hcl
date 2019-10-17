Kind = "service-resolver"
Name = "todoapp1"


failover = {
  "*" = {
    datacenters = ["dc2"]
  }
}
