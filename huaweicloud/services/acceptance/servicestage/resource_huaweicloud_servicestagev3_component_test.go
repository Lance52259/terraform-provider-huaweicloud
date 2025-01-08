package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
)

func getV3ComponentFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("servicestage", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}
	return servicestage.QueryV3Component(client, state.Primary.Attributes["application_id"], state.Primary.ID,
		state.Primary.Attributes["version"])
}

func TestAccV3Component_basic(t *testing.T) {
	var (
		component interface{}

		resourceName = "huaweicloud_servicestagev3_component.test"
		rc           = acceptance.InitResourceCheck(resourceName, &component, getV3ComponentFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
			// Make sure the networks of the CCE cluster and the CSE engine are same.
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
			acceptance.TestAccPreCheckImsImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Component_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "huaweicloud_servicestagev3_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "huaweicloud_servicestagev3_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.deploy_mode", "container"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.name", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.type", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "source",
						fmt.Sprintf("{\"auth\":\"iam\",\"kind\":\"image\",\"storage\":\"swr\",\"url\":\"%s\"}", acceptance.HW_IMS_IMAGE_URL)),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "replica", "2"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "limit_cpu", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "limit_memory", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "request_cpu", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "request_memory", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "envs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.name", "env_name"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.value", "env_value"),
					resource.TestCheckResourceAttr(resourceName, "storages.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.type", "HostPath"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "storages.0.parameters", "{\"default_mode\":0,\"path\":\"/tmp\"}"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.path", "/category"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.sub_path", "sub"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "command", "{\"args\":[\"-a\"],\"command\":[\"ls\"]}"),
					resource.TestCheckResourceAttr(resourceName, "post_start.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.command.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.type", "command"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.command.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.type", "command"),
					resource.TestCheckResourceAttr(resourceName, "mesher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "mesher.0.port", "60"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(resourceName, "logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.log_path", "/tmp"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.rotate", "Hourly"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.host_path", "/tmp"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.host_extend_path", "PodName"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.path", "/tmp"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.port", "600"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.dimensions", "cpu_usage,mem_usage"),
					resource.TestCheckResourceAttr(resourceName, "affinity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "affinity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.type", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.port", "800"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.type", "http"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.scheme", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.host", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.port", "8000"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.path", "/v1/test"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV3Component_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "huaweicloud_servicestagev3_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "huaweicloud_servicestagev3_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.deploy_mode", "container"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.name", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.type", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "source",
						fmt.Sprintf("{\"auth\":\"iam\",\"kind\":\"image\",\"storage\":\"swr\",\"url\":\"%s\"}", acceptance.HW_IMS_IMAGE_URL)),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "replica", "2"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "limit_cpu", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "limit_memory", "1"),
					resource.TestCheckResourceAttr(resourceName, "request_cpu", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "request_memory", "1"),
					resource.TestCheckResourceAttr(resourceName, "envs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.name", "new_env_name"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.value", "new_env_value"),
					resource.TestCheckResourceAttr(resourceName, "storages.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.type", "HostPath"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.name", fmt.Sprintf("%s-new", name)),
					resource.TestCheckResourceAttr(resourceName, "storages.0.parameters", "{\"default_mode\":0,\"path\":\"/tmp/new\"}"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.path", "/category/new"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.sub_path", "sub/new"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.read_only", "true"),
					resource.TestCheckResourceAttr(resourceName, "command", "{\"args\":[\"-l\"],\"command\":[\"ls\"]}"),
					resource.TestCheckResourceAttr(resourceName, "post_start.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.command.0", "newtest"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.type", "command"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.command.0", "newtest"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.type", "command"),
					resource.TestCheckResourceAttr(resourceName, "mesher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "mesher.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/HongKong"),
					resource.TestCheckResourceAttr(resourceName, "logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.log_path", "/tmp/new"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.rotate", "Daily"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.host_path", "/tmp/new"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.host_extend_path", "PodUID"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.path", "/tmp/new"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.port", "800"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.dimensions", "mem_usage"),
					resource.TestCheckResourceAttr(resourceName, "affinity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "affinity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.type", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.delay", "60"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.timeout", "60"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.port", "900"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.type", "http"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.delay", "60"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.timeout", "60"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.scheme", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.host", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.path", "/v1/test/new"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3ComponentImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func testAccV3ComponentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var applicationId, resourceId string
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of compnent is not found in the tfstate", resourceName)
		}
		applicationId = rs.Primary.Attributes["application_id"]
		resourceId = rs.Primary.ID
		if applicationId == "" || resourceId == "" {
			return "", fmt.Errorf("the component ID is not exist or application ID is missing")
		}
		return fmt.Sprintf("%s/%s", applicationId, resourceId), nil
	}
}

func testAccV3Component_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cce_clusters" "test" {
  cluster_id = "%[1]s"
}

data "huaweicloud_servicestagev3_runtime_stacks" "test" {}

locals {
  docker_runtime_stack = try([for v in data.huaweicloud_servicestagev3_runtime_stacks.test.runtime_stacks: v if
    v.type == "Docker" && v.status == "Supported"][0], null)
}

resource "huaweicloud_servicestagev3_application" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "0"
}

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  vpc_id                = try(data.huaweicloud_cce_clusters.test.clusters[0].vpc_id, "")
  enterprise_project_id = "0"
}

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  resources {
    id   = "%[1]s"
    type = "cce"
  }
  resources {
    id   = "%[3]s"
    type = "cse"
  }
}
`, acceptance.HW_CCE_CLUSTER_ID,
		name,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ID)
}

func testAccV3Component_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component" "test" {
  depends_on = [
    huaweicloud_servicestagev3_environment_associate.test
  ]

  application_id = huaweicloud_servicestagev3_application.test.id
  environment_id = huaweicloud_servicestagev3_environment.test.id
  name           = "%[2]s"

  runtime_stack {
    deploy_mode = try(local.docker_runtime_stack.deploy_mode, null)
    name        = try(local.docker_runtime_stack.name, null)
    type        = try(local.docker_runtime_stack.type, null)
    version     = try(local.docker_runtime_stack.version, null)
  }

  source = jsonencode({
    "auth": "iam",
    "kind": "image",
    "storage": "swr",
    "url": "%[3]s"
  })

  version = "1.0.1"
  replica = 2

  refer_resources {
    id         = "%[4]s"
    type       = "cce"
    parameters = jsonencode({
      "namespace": "default",
      "type": "VirtualMachine"
    })
  }
  refer_resources {
    id   = "%[5]s"
    type = "cse"
  }

  tags = {
    foo = "bar"
  }

  description    = "Created by terraform script"
  limit_cpu      = 0.25
  limit_memory   = 0.5
  request_cpu    = 0.25
  request_memory = 0.5

  envs {
    name  = "env_name"
    value = "env_value"
  }

  storages {
    type       = "HostPath"
    name       = "%[2]s"
    parameters = jsonencode({
      "default_mode": 0,
      "path": "/tmp"
    })
    mounts {
      path      = "/category"
      sub_path  = "sub"
      read_only = false
    }
  }

  command = jsonencode({
    "args": ["-a"],
    "command": ["ls"]
  })

  post_start {
    command = ["test"]
    type    = "command"
  }

  pre_stop {
    command = ["test"]
    type    = "command"
  }

  mesher {
    port = 60
  }

  timezone = "Asia/Shanghai"

  logs {
    log_path         = "/tmp"
    rotate           = "Hourly"
    host_path        = "/tmp"
    host_extend_path = "PodName"
  }

  custom_metric {
    path       = "/tmp"
    port       = 600
    dimensions = "cpu_usage,mem_usage"
  }

  affinity {
    condition = "required"
    kind      = "node"
    match_expressions {
      key       = "affinity1"
      value     = "foo"
      operation = "In"
    }
    weight = 100
  }
  affinity {
    condition = "preferred"
    kind      = "node"
    match_expressions {
      key       = "affinity2"
      value     = "bar"
      operation = "NotIn"
    }
    weight = 1
  }

  anti_affinity {
    condition = "required"
    kind      = "pod"
    match_expressions {
      key       = "anit-affinity1"
      operation = "Exists"
    }
    weight = 100
  }
  anti_affinity {
    condition = "preferred"
    kind      = "pod"
    match_expressions {
      key       = "anti-affinity2"
      operation = "DoesNotExist"
    }
    weight = 1
  }

  liveness_probe {
    type    = "tcp"
    delay   = 30
    timeout = 30
    port    = 800
  }

  readiness_probe {
    type    = "http"
    delay   = 30
    timeout = 30
    scheme  = "HTTPS"
    host    = "127.0.0.1"
    port    = 8000
    path    = "/v1/test"
  }
}
`, testAccV3Component_base(name),
		name,
		acceptance.HW_IMS_IMAGE_URL,
		acceptance.HW_CCE_CLUSTER_ID,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ID)
}

func testAccV3Component_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component" "test" {
  depends_on = [
    huaweicloud_servicestagev3_environment_associate.test
  ]

  application_id = huaweicloud_servicestagev3_application.test.id
  environment_id = huaweicloud_servicestagev3_environment.test.id
  name           = "%[2]s"

  runtime_stack {
    deploy_mode = try(local.docker_runtime_stack.deploy_mode, null)
    name        = try(local.docker_runtime_stack.name, null)
    type        = try(local.docker_runtime_stack.type, null)
    version     = try(local.docker_runtime_stack.version, null)
  }

  source = jsonencode({
    "auth": "iam",
    "kind": "image",
    "storage": "swr",
    "url": "%[3]s"
  })

  version = "1.0.2"
  replica = 2

  refer_resources {
    id         = "%[4]s"
    type       = "cce"
    parameters = jsonencode({
      "namespace": "default",
      "type": "VirtualMachine"
    })
  }
  refer_resources {
    id   = "%[5]s"
    type = "cse"
  }

  tags = {
    foo = "baar"
  }

  description    = "Updated by terraform script"
  limit_cpu      = 0.5
  limit_memory   = 1
  request_cpu    = 0.5
  request_memory = 1

  envs {
    name  = "new_env_name"
    value = "new_env_value"
  }

  storages {
    type       = "HostPath"
    name       = "%[2]s-new"
    parameters = jsonencode({
      "default_mode": 0,
      "path": "/tmp/new"
    })
    mounts {
      path      = "/category/new"
      sub_path  = "sub/new"
      read_only = true
    }
  }

  command = jsonencode({
    "args": ["-l"],
    "command": ["ls"]
  })

  post_start {
    command = ["newtest"]
    type    = "command"
  }

  pre_stop {
    command = ["newtest"]
    type    = "command"
  }

  mesher {
    port = 80
  }

  timezone = "Asia/HongKong"

  logs {
    log_path         = "/tmp/new"
    rotate           = "Daily"
    host_path        = "/tmp/new"
    host_extend_path = "PodUID"
  }

  custom_metric {
    path       = "/tmp/new"
    port       = 800
    dimensions = "mem_usage"
  }

  affinity {
    condition = "required"
    kind      = "node"
    match_expressions {
      key       = "new_affinity1"
      value     = "1"
      operation = "Gt"
    }
    weight = 100
  }
  affinity {
    condition = "preferred"
    kind      = "node"
    match_expressions {
      key       = "new_affinity2"
      value     = "100"
      operation = "Lt"
    }
    weight = 1
  }

  anti_affinity {
    condition = "required"
    kind      = "pod"
    match_expressions {
      key       = "new_anit-affinity1"
      operation = "Exists"
    }
    weight = 100
  }
  anti_affinity {
    condition = "preferred"
    kind      = "pod"
    match_expressions {
      key       = "new_anti-affinity2"
      operation = "DoesNotExist"
    }
    weight = 1
  }

  liveness_probe {
    type    = "tcp"
    delay   = 60
    timeout = 60
    port    = 900
  }

  readiness_probe {
    type    = "http"
    delay   = 60
    timeout = 60
    scheme  = "HTTP"
    host    = "192.168.0.1"
    port    = 8080
    path    = "/v1/test/new"
  }
}
`, testAccV3Component_base(name),
		name,
		acceptance.HW_IMS_IMAGE_URL,
		acceptance.HW_CCE_CLUSTER_ID,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ID)
}

func TestAccV3Component_deploy(t *testing.T) {
	var (
		component interface{}

		resourceName = "huaweicloud_servicestagev3_component.test"
		rc           = acceptance.InitResourceCheck(resourceName, &component, getV3ComponentFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckServiceStageEnabled(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckPublicRemoteIpAddress(t)
			acceptance.TestAccPreCheckJarObsStoragePaths(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Component_deploy_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "huaweicloud_servicestagev3_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "huaweicloud_servicestagev3_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttrSet(resourceName, "build"),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.deploy_mode", "container"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.name"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.type", "Java"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.version"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "limit_cpu", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "limit_memory", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "request_cpu", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "request_memory", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "replica", "2"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.0.forward_port", "8000"),
					resource.TestCheckResourceAttrPair(resourceName, "external_accesses.0.address",
						"huaweicloud_elb_certificate.test", "domain"),
					resource.TestCheckResourceAttr(resourceName, "envs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.name", "MOCK_ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.value", "false"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Shanghai"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV3Component_deploy_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "deploy_strategy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "deploy_strategy.0.type", "GrayRelease"),
					resource.TestCheckResourceAttrSet(resourceName, "deploy_strategy.0.gray_release"),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3ComponentImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func testAccV3Component_deploy_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 ="%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "ingress_ipv4_icmp_allow" {
  security_group_id = try(huaweicloud_networking_secgroup.test.id, "")
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "%[2]s/32" // Only specified PC can access it.
  priority          = 1
}

resource "huaweicloud_networking_secgroup_rule" "ingress_ipv4_tcp_allow" {
  security_group_id = try(huaweicloud_networking_secgroup.test.id, "")
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "22-23,443,30100-30130"
  remote_ip_prefix  = "%[2]s/32" // Only specified PC can access it.
  priority          = 1
}

resource "huaweicloud_networking_secgroup_rule" "ingress_ipv4_all_allow_inside_group" {
  security_group_id = try(huaweicloud_networking_secgroup.test.id, "")
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_group_id   = try(huaweicloud_networking_secgroup.test.id, "")
  priority          = 1
}

resource "huaweicloud_networking_secgroup_rule" "egress_ip6_all_allow" {
  security_group_id = try(huaweicloud_networking_secgroup.test.id, "")
  direction         = "egress"
  ethertype         = "IPv6"
  remote_ip_prefix  = "::/0"
  priority          = 1
}

resource "huaweicloud_networking_secgroup_rule" "egress_ipv4_all_allow" {
  security_group_id = try(huaweicloud_networking_secgroup.test.id, "")
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  priority          = 1
}

data "huaweicloud_cse_microservice_engine_flavors" "test" {}

locals {
  engine_flavor_id = try([for o in data.huaweicloud_cse_microservice_engine_flavors.test.flavors: o.id if
    strcontains(o.id, "s1") && strcontains(o.id, "small")][0], "")
}

resource "huaweicloud_vpc_eip" "cse_used" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cse_microservice_engine" "test" {
  auth_type             = "NONE"
  name                  = "%[1]s"
  description           = "Created by terraform test"
  flavor                = local.engine_flavor_id
  network_id            = huaweicloud_vpc_subnet.test.id
  eip_id                = huaweicloud_vpc_eip.cse_used.id
  enterprise_project_id = "%[3]s"
  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[1]s"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  flavor_id              = "cce.s2.small"
  container_network_type = "vpc-router"
  cluster_version        = "v1.19"
  cluster_type           = "VirtualMachine"

  kube_proxy_mode = "iptables"

  dynamic "masters" {
    for_each = slice(data.huaweicloud_availability_zones.test.names, 0, 3)

    content {
      availability_zone = masters.value
    }
  }
}

resource "huaweicloud_vpc_eip" "cce_used" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-cce-used"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }

  lifecycle {
    ignore_changes = [
      tags,
    ]
  }
}

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%[1]s"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_kps_keypair.test.name
  eip_id            = huaweicloud_vpc_eip.cce_used.id

  root_volume {
    volumetype = "SSD"
    size       = 100
  }

  data_volumes {
    volumetype = "SSD"
    size       = 100
  }
}

resource "huaweicloud_servicestagev3_environment" "test" {
  vpc_id                = huaweicloud_vpc.test.id
  name                  = "%[1]s"
  deploy_mode           = "container"
  enterprise_project_id = "%[3]s"
}

resource "huaweicloud_vpc_eip" "servicestage_used" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-cce-used"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }

  lifecycle {
    ignore_changes = [
      tags,
    ]
  }
}

resource "huaweicloud_vpc_eip" "loadbalancer_used" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-loadbalancer-used"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }

  lifecycle {
    ignore_changes = [
      tags,
    ]
  }
}

data "huaweicloud_elb_flavors" "l4flavors" {
  type            = "L4"
  max_connections = 1000000
  cps             = 20000
  bandwidth       = 100
}

data "huaweicloud_elb_flavors" "l7flavors" {
  type            = "L7"
  max_connections = 400000
  cps             = 4000
  bandwidth       = 100
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[1]s"
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  l4_flavor_id      = data.huaweicloud_elb_flavors.l4flavors.ids[0]
  l7_flavor_id      = data.huaweicloud_elb_flavors.l7flavors.ids[0]
  availability_zone = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  ipv4_eip_id       = huaweicloud_vpc_eip.loadbalancer_used.id
}

resource "huaweicloud_elb_certificate" "test" {
  name        = "%[1]s"
  domain      = "p2cserver.com"
  type        = "server"
  private_key = "%[4]s"
  certificate = "%[5]s"
}

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  resources {
    id   = huaweicloud_cce_cluster.test.id
    type = "cce"
  }
  resources {
    id   = huaweicloud_vpc_eip.servicestage_used.id
    type = "eip"
  }
  resources {
    id   = huaweicloud_cse_microservice_engine.test.id
    type = "cse"
  }
  resources {
    id   = huaweicloud_elb_loadbalancer.test.id
    type = "elb"
  }
}

resource "huaweicloud_servicestagev3_application" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[3]s"
}

data "huaweicloud_servicestagev3_runtime_stacks" "test" {}

locals {
  java_runtime_stack = try([for o in data.huaweicloud_servicestagev3_runtime_stacks.test.runtime_stacks: o if o.type == "Java" && o.deploy_mode == "container"][0], {})
}
`, name, acceptance.HW_PUBLIC_REMOTE_IP_ADDRESS,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY,
		acceptance.HW_CERTIFICATE_CONTENT)
}

func testAccV3Component_deploy_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component" "test" {
  depends_on = [
    huaweicloud_elb_certificate.test
    huaweicloud_servicestagev3_environment_associate.test
  ]

  application_id = huaweicloud_servicestagev3_application.test.id
  environment_id = huaweicloud_servicestagev3_environment.test.id
  name           = "%[2]s"
  version        = "1.0.1"

  build = jsonencode({
    "parameters": {
      "artifact_namespace": "%[3]s",
      "environment_id": huaweicloud_servicestagev3_environment.test.id,
      "cluster_namespace": "default",
      "use_public_cluster": false,
      "cluster_id": huaweicloud_cce_cluster.test.id,
      "dockerfile_path": "./",
      "build_env_selected": "current"
    }
  })

  source = jsonencode({
    "kind": "package",
    "url": try(element(split("%[4]s", ","), 0), ""),
    "storage": "obs"
  })

  runtime_stack {
    deploy_mode = try(local.java_runtime_stack.deploy_mode, "container")
    name        = try(local.java_runtime_stack.name, "OpenJDK17")
    type        = try(local.java_runtime_stack.type, "Java")
    version     = try(local.java_runtime_stack.version, null)
  }

  refer_resources {
    id         = huaweicloud_cce_cluster.test.id
    type       = "cce"
    parameters = jsonencode({
      "namespace": "default",
      "type": "VirtualMachine"
      "name": huaweicloud_cce_cluster.test.name
    })
  }
  refer_resources {
    id         = huaweicloud_elb_loadbalancer.test.id
    type       = "elb"
    parameters = jsonencode({
      "name": huaweicloud_elb_loadbalancer.test.name
    })
  }
  refer_resources {
    id   = huaweicloud_cse_microservice_engine.test.id
    type = "cse"
  }

  limit_cpu      = 0.25
  limit_memory   = 0.5
  request_cpu    = 0.25
  request_memory = 0.5
  replica        = 2

  external_accesses {
    protocol     = "HTTP"
    forward_port = "8000"
    address      = huaweicloud_elb_certificate.test.domain
  }

  envs {
    name  = "MOCK_ENABLED"
    value = "false"
  }

  timezone = "Asia/Shanghai"
}
`, testAccV3Component_deploy_base(name), name,
		acceptance.HW_SERVICESTAGE_ORGANIZATION_NAME,
		acceptance.HW_JAR_OBS_STORAGE_PATHS)
}

func testAccV3Component_deploy_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component" "test" {
  depends_on = [
    huaweicloud_elb_certificate.test
    huaweicloud_servicestagev3_environment_associate.test
  ]

  application_id = huaweicloud_servicestagev3_application.test.id
  environment_id = huaweicloud_servicestagev3_environment.test.id
  name           = "%[2]s"
  version        = "1.0.2"

  build = jsonencode({
    "parameters": {
      "artifact_namespace": "%[3]s",
      "environment_id": huaweicloud_servicestagev3_environment.test.id,
      "cluster_namespace": "default",
      "use_public_cluster": false,
      "cluster_id": huaweicloud_cce_cluster.test.id,
      "dockerfile_path": "./",
      "build_env_selected": "current"
    }
  })

  source = jsonencode({
    "kind": "package",
    "url": try(element(split("%[4]s", ","), 1), ""),
    "storage": "obs"
  })

  runtime_stack {
    deploy_mode = try(local.java_runtime_stack.deploy_mode, "container")
    name        = try(local.java_runtime_stack.name, "OpenJDK17")
    type        = try(local.java_runtime_stack.type, "Java")
    version     = try(local.java_runtime_stack.version, null)
  }

  refer_resources {
    id         = huaweicloud_cce_cluster.test.id
    type       = "cce"
    parameters = jsonencode({
      "namespace": "default",
      "type": "VirtualMachine"
      "name": huaweicloud_cce_cluster.test.name
    })
  }
  refer_resources {
    id         = huaweicloud_elb_loadbalancer.test.id
    type       = "elb"
    parameters = jsonencode({
      "name": huaweicloud_elb_loadbalancer.test.name
    })
  }
  refer_resources {
    id   = huaweicloud_cse_microservice_engine.test.id
    type = "cse"
  }

  limit_cpu      = 0.25
  limit_memory   = 0.5
  request_cpu    = 0.25
  request_memory = 0.5
  replica        = 2

  external_accesses {
    protocol     = "HTTP"
    forward_port = "8000"
    address      = huaweicloud_elb_certificate.test.domain
  }
  external_accesses {
    protocol = "tcp"
    address  = huaweicloud_vpc_eip.cce_used.address
  }

  envs {
    name  = "MOCK_ENABLED"
    value = "false"
  }

  timezone = "Asia/Shanghai"

  deploy_strategy {
    type         = "GrayRelease"
    gray_release = jsonencode({
      "type": "WEIGHT",
      "first_batch_replica": 1,
      "first_batch_weight": 10,
      "remaining_batch": 1,
      "deployment_mode": "4",
      "replica_surge_mode": "EXTRA"
    })
  }
}
`, testAccV3Component_deploy_base(name), name,
		acceptance.HW_SERVICESTAGE_ORGANIZATION_NAME,
		acceptance.HW_JAR_OBS_STORAGE_PATHS)
}
