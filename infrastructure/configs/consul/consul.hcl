# Lyss AI Platform Consul 配置

datacenter = "dc1"
data_dir = "/consul/data"
log_level = "INFO"
server = true
bootstrap_expect = 1

# UI 配置
ui_config {
  enabled = true
}

# 网络绑定
bind_addr = "0.0.0.0"
client_addr = "0.0.0.0"
retry_join = ["127.0.0.1"]

# Connect 服务网格
connect {
  enabled = true
}

# ACL 配置 (开发环境禁用)
acl = {
  enabled = false
  default_policy = "allow"
}

# 性能配置
performance {
  raft_multiplier = 1
}

# 端口配置
ports {
  grpc = 8502
}

# 服务定义目录
services {
  name = "consul"
  tags = ["consul", "service-discovery"]
  port = 8500
  check {
    http = "http://localhost:8500/v1/status/leader"
    interval = "10s"
  }
}

# Lyss AI Platform 服务预定义
services {
  name = "user-service"
  tags = ["lyss", "microservice", "user"]
  port = 8001
  check {
    http = "http://localhost:8001/health"
    interval = "30s"
  }
}

services {
  name = "auth-service"
  tags = ["lyss", "microservice", "auth"]
  port = 8002
  check {
    http = "http://localhost:8002/health"
    interval = "30s"
  }
}

services {
  name = "group-service"
  tags = ["lyss", "microservice", "group"]
  port = 8003
  check {
    http = "http://localhost:8003/health"
    interval = "30s"
  }
}

services {
  name = "credential-service"
  tags = ["lyss", "microservice", "credential"]
  port = 8004
  check {
    http = "http://localhost:8004/health"
    interval = "30s"
  }
}

services {
  name = "gateway-service"
  tags = ["lyss", "microservice", "gateway"]
  port = 8000
  check {
    http = "http://localhost:8000/health"
    interval = "30s"
  }
}

services {
  name = "billing-service"
  tags = ["lyss", "microservice", "billing"]
  port = 8005
  check {
    http = "http://localhost:8005/health"
    interval = "30s"
  }
}