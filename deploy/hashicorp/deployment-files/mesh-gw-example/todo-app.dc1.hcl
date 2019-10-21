job "todo-app"{

  datacenters = ["dc2-ncv"]
  region      = "dc2-region"
  type        = "service"

  group "todo-cluster" {
    count = 1

    task "todoapp1" {
      driver = "docker"
      config {
        image = "javier1/todo-app"
      }

      env {
        APP_PORT = "4000"
        APP_HOSTNAME = "localhost"
        DB_USER = "default"
        DB_PASSWORD = "pass"
        DB_HOST= "localhost"
        SECRET_KEY_BASE = "Y0uRvErYsecr3TANDL0ngStr1n"
      }

      resources {
        cpu    = 100
        memory = 300
      }

    }
    network {
      mode = "bridge"
      port "http" {
        static = 4000
        to     = 4000
      }
    }

    service  {
      name = "todoapp1"
      port = "4000"
      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgressdb1"
              local_bind_port = 5432
            }
          }
        }
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
          "-http-addr", "172.20.20.31:8500",
          "-grpc-addr", "172.20.20.31:8502",
          "-wan-address", "172.20.20.31:${NOMAD_PORT_proxy}",
          "-address", "172.20.20.31:${NOMAD_PORT_proxy}",
          "-bind-address", "default=172.20.20.31:${NOMAD_PORT_proxy}",
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
