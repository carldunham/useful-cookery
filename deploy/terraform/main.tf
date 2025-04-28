terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.10"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.5"
    }
  }

  backend "s3" {
    bucket         = "useful-cookery-terraform-state"
    key            = "terraform.tfstate"
    region         = "us-west-2"
    encrypt        = true
    dynamodb_table = "useful-cookery-terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "useful-cookery"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}

# EKS Cluster
module "eks" {
  source = "./modules/eks"

  cluster_name    = "useful-cookery-${var.environment}"
  cluster_version = var.eks_cluster_version
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnets

  node_groups = {
    main = {
      desired_capacity = var.eks_node_group_desired_capacity
      min_capacity     = var.eks_node_group_min_capacity
      max_capacity     = var.eks_node_group_max_capacity
      instance_types   = var.eks_node_group_instance_types
      disk_size        = var.eks_node_group_disk_size
    }
  }
}

# VPC for EKS
module "vpc" {
  source = "./modules/vpc"

  name               = "useful-cookery-${var.environment}"
  cidr               = var.vpc_cidr
  azs                = var.availability_zones
  private_subnets    = var.private_subnet_cidrs
  public_subnets     = var.public_subnet_cidrs
  enable_nat_gateway = true
  single_nat_gateway = var.environment != "production"
}

# RDS PostgreSQL
module "rds" {
  source = "./modules/rds"

  identifier           = "useful-cookery-${var.environment}"
  engine               = "postgres"
  engine_version       = var.postgres_version
  instance_class       = var.rds_instance_class
  allocated_storage    = var.rds_allocated_storage
  max_allocated_storage = var.rds_max_allocated_storage

  db_name              = "usefulcookery"
  username             = "postgres"
  password             = var.rds_password

  vpc_id               = module.vpc.vpc_id
  subnet_ids           = module.vpc.database_subnets

  multi_az             = var.environment == "production"
  backup_retention_period = var.environment == "production" ? 7 : 1
  deletion_protection  = var.environment == "production"

  allowed_security_group_ids = [module.eks.node_security_group_id]
}

# CloudFront CDN for UI
module "cloudfront" {
  source = "./modules/cloudfront"

  name                = "useful-cookery-ui-${var.environment}"
  domain_name         = "www.useful-cookery.com"
  alb_domain_name     = module.eks.alb_hostname
  acm_certificate_arn = var.acm_certificate_arn
}

# Kubernetes provider configuration
provider "kubernetes" {
  host                   = module.eks.cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)
  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
    command     = "aws"
  }
}

# Helm provider configuration
provider "helm" {
  kubernetes {
    host                   = module.eks.cluster_endpoint
    cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)
    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
      command     = "aws"
    }
  }
}

# Install AWS Load Balancer Controller
module "aws_load_balancer_controller" {
  source = "./modules/aws-load-balancer-controller"

  cluster_name = module.eks.cluster_name
  vpc_id       = module.vpc.vpc_id
}

# Install External DNS
module "external_dns" {
  source = "./modules/external-dns"

  cluster_name = module.eks.cluster_name
  domain_name  = "useful-cookery.com"
}

# Create Kubernetes namespace
resource "kubernetes_namespace" "useful_cookery" {
  metadata {
    name = "useful-cookery-${var.environment}"
  }
}

# Create Kubernetes secrets
resource "kubernetes_secret" "useful_cookery_secrets" {
  metadata {
    name      = "useful-cookery-secrets"
    namespace = kubernetes_namespace.useful_cookery.metadata[0].name
  }

  data = {
    "database.connection_string" = "postgres://${var.rds_username}:${var.rds_password}@${module.rds.endpoint}/${var.rds_db_name}?sslmode=require"
    "jwt.secret"                 = var.jwt_secret
    "openai.api_key"             = var.openai_api_key
  }
}
