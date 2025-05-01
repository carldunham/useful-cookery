variable "aws_region" {
  description = "The AWS region to deploy resources"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, production)"
  type        = string
  default     = "dev"
}

# VPC Variables
variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones to use"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
}

# EKS Variables
variable "eks_cluster_version" {
  description = "Kubernetes version for EKS cluster"
  type        = string
  default     = "1.27"
}

variable "eks_node_group_desired_capacity" {
  description = "Desired number of nodes in the EKS node group"
  type        = number
  default     = 2
}

variable "eks_node_group_min_capacity" {
  description = "Minimum number of nodes in the EKS node group"
  type        = number
  default     = 1
}

variable "eks_node_group_max_capacity" {
  description = "Maximum number of nodes in the EKS node group"
  type        = number
  default     = 5
}

variable "eks_node_group_instance_types" {
  description = "Instance types for the EKS node group"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "eks_node_group_disk_size" {
  description = "Disk size in GB for EKS node group instances"
  type        = number
  default     = 50
}

# RDS Variables
variable "postgres_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "14.6"
}

variable "rds_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.medium"
}

variable "rds_allocated_storage" {
  description = "Allocated storage for RDS instance in GB"
  type        = number
  default     = 20
}

variable "rds_max_allocated_storage" {
  description = "Maximum allocated storage for RDS instance in GB"
  type        = number
  default     = 100
}

variable "rds_username" {
  description = "Username for RDS instance"
  type        = string
  default     = "postgres"
}

variable "rds_password" {
  description = "Password for RDS instance"
  type        = string
  sensitive   = true
}

variable "rds_db_name" {
  description = "Database name for RDS instance"
  type        = string
  default     = "usefulcookery"
}

# CloudFront Variables
variable "acm_certificate_arn" {
  description = "ARN of ACM certificate for CloudFront distribution"
  type        = string
}

# Application Variables
variable "jwt_secret" {
  description = "Secret key for JWT token generation"
  type        = string
  sensitive   = true
}

variable "openai_api_key" {
  description = "API key for OpenAI services"
  type        = string
  sensitive   = true
}
