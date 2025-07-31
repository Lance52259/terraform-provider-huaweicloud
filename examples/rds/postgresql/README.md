# Create a PostgreSQL RDS instance

This example provides best practice code for using Terraform to create a PostgreSQL RDS instance in HuaweiCloud RDS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the PostgreSQL RDS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The PostgreSQL RDS instance name
* `account_name` - Username with elevated privileges
* `database_name` - The name of the initial database
* `schema_name` - The name of the database schema
* `backup_name` - The name for instance backups
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format
* `instance_backup_keep_days` - The number of days to retain backups

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `availability_zone` - The availability zone to which the RDS instance belongs (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `instance_flavor_id` - The flavor ID of the RDS instance (default: "")
* `instance_db_type` - The database engine type (default: "PostgreSQL")
* `instance_db_version` - The database engine version (default: "16")
* `instance_db_port` - The database port (default: 5432)
* `instance_password` - The password for the RDS instance (default: "")
* `account_password` - The password for the database account (default: "")
* `instance_mode` - The instance mode for the RDS instance flavor (default: "single")
* `instance_flavor_group_type` - The group type for the RDS instance flavor (default: "general")
* `instance_flavor_vcpus` - The number of the RDS instance CPU cores for the RDS instance flavor (default: 2)
* `instance_volume_type` - The storage volume type (default: "CLOUDSSD")
* `instance_volume_size` - The storage volume size in GB (default: 40)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                    = "your_vpc_name"
  subnet_name                 = "your_subnet_name"
  security_group_name         = "your_security_group_name"
  instance_name               = "your_postgresql_instance_name"
  account_name                = "your_account_name"
  database_name               = "your_database_name"
  schema_name                 = "your_schema_name"
  backup_name                 = "your_backup_name"
  instance_backup_time_window = "08:00-09:00"
  instance_backup_keep_days   = 1
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Features

This example demonstrates the following features:

1. **PostgreSQL RDS Instance Creation**: Creates a complete PostgreSQL RDS instance with all necessary components
2. **Network Configuration**: Sets up VPC, subnet, and security group for the RDS instance
3. **Account Management**: Creates a PostgreSQL account with elevated privileges
4. **Database and Schema Management**: Creates a database and schema with proper ownership
5. **Backup Strategy**: Configures automated backup with customizable time window and retention period

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the PostgreSQL RDS instance takes about 5 minutes
* This example creates the PostgreSQL RDS instance, VPC, subnet, security group, account, database, schema, and backup
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.54.0 |
| random | >= 3.0.0 |
