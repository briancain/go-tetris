# ElastiCache Subnet Group
resource "aws_elasticache_subnet_group" "redis" {
  name       = "${var.project_name}-redis-subnet-group-${var.environment}"
  subnet_ids = aws_subnet.private[*].id

  tags = {
    Name = "${var.project_name}-redis-subnet-group-${var.environment}"
  }
}

# ElastiCache Redis Cluster
resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "${var.project_name}-redis-${var.environment}"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.redis.name
  security_group_ids   = [aws_security_group.redis.id]

  tags = {
    Name = "${var.project_name}-redis-${var.environment}"
  }
}

# DynamoDB Table - Players
resource "aws_dynamodb_table" "players" {
  name         = "${var.project_name}-players-${var.environment}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "playerId"

  attribute {
    name = "playerId"
    type = "S"
  }

  tags = {
    Name = "${var.project_name}-players-${var.environment}"
  }
}

# DynamoDB Table - Games
resource "aws_dynamodb_table" "games" {
  name         = "${var.project_name}-games-${var.environment}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "gameId"

  attribute {
    name = "gameId"
    type = "S"
  }

  tags = {
    Name = "${var.project_name}-games-${var.environment}"
  }
}

# DynamoDB Table - Leaderboard
resource "aws_dynamodb_table" "leaderboard" {
  name         = "${var.project_name}-leaderboard-${var.environment}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "playerId"

  attribute {
    name = "playerId"
    type = "S"
  }

  attribute {
    name = "highScore"
    type = "N"
  }

  global_secondary_index {
    name            = "highScore-index"
    hash_key        = "highScore"
    projection_type = "ALL"
  }

  tags = {
    Name = "${var.project_name}-leaderboard-${var.environment}"
  }
}
