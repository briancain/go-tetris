variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "go-tetris"
}

variable "domain_name" {
  description = "Domain name for SSL certificate (e.g., api.tetris.example.com)"
  type        = string
  default     = ""
}

variable "enable_ssl" {
  description = "Enable SSL/HTTPS support with ACM certificate"
  type        = bool
  default     = false
}
