job "db-cluster" {

  datacenters = ["dc1-ncv"]
  region      = "dc1-region"
  type        = "service"

  group "db-cluster" {
    count = 1

    task "postgressdb1" {
      driver = "docker"
      config {
        image   = "postgres"
      }

      resources {
        cpu    = 100
        memory = 100
      }

      env {
        POSTGRES_DB = "demo_prod"
        POSTGRES_USER = "default"
        POSTGRES_PASSWORD = "pass"
      }
    }

    network {
      mode = "bridge"
      port "http" {
        static = 5432
        to     = 5432
      }
    }

    service {
      name = "postgressdb1"
      port = "5432"

      connect {
        sidecar_service {}
      }
    }
  }


  group "mesh-gateway" {
    count = 1

    task "mesh-gateway" {
      driver = "exec"

      config {
        command = "consul"
        args    = [
          "connect", "envoy",
          "-mesh-gateway",
          "-register",
          "-http-addr", "172.20.20.11:8500",
          "-grpc-addr", "172.20.20.11:8502",
          "-wan-address", "172.20.20.11:${NOMAD_PORT_proxy}",
          "-address", "172.20.20.11:${NOMAD_PORT_proxy}",
          "-bind-address", "default=172.20.20.11:${NOMAD_PORT_proxy}",
          "--",
          "-l", "debug"
        ]
      }

      resources {
        cpu    = 100
        memory = 100

        network {
          port "proxy" {
            static = 8433
          }
        }
      }
    }
  }
}
