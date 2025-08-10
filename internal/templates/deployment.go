package templates

import (
	"fmt"
	"strings"

	"github.com/vibercode/cli/internal/models"
)

// GetDockerTemplate returns Docker template based on file type
func GetDockerTemplate(fileName string, config models.DeploymentConfig) string {
	switch fileName {
	case "Dockerfile", "Dockerfile.production":
		if config.MultiStage {
			return getMultiStageDockerfile(config)
		}
		return getStandardDockerfile(config)
	case "Dockerfile.multi-stage":
		return getMultiStageDockerfile(config)
	case "docker-compose.yml":
		return getDockerCompose(config)
	case "docker-compose.production.yml":
		return getDockerComposeProduction(config)
	case ".dockerignore":
		return getDockerIgnore()
	default:
		return ""
	}
}

// getMultiStageDockerfile generates multi-stage Dockerfile
func getMultiStageDockerfile(config models.DeploymentConfig) string {
	baseImage := config.GetBaseImage()
	finalImage := config.GetFinalImage()
	
	dockerfile := fmt.Sprintf(`# Build stage
FROM %s AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# Final stage
FROM %s

`, baseImage, finalImage)

	if config.Security {
		dockerfile += `# Security hardening
RUN apk --no-cache add ca-certificates && \
    apk update && apk upgrade && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN adduser -D -s /bin/sh -u 1001 appuser

WORKDIR /home/appuser

# Copy binary from builder stage
COPY --from=builder --chown=appuser:appuser /app/main .
COPY --from=builder --chown=appuser:appuser /app/config ./config

# Switch to non-root user
USER appuser

`
	} else {
		dockerfile += `WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

`
	}

	dockerfile += fmt.Sprintf(`# Expose port
EXPOSE %d

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:%d/health || exit 1

# Run the application
CMD ["./main"]
`, config.Port, config.Port)

	return dockerfile
}

// getStandardDockerfile generates standard Dockerfile
func getStandardDockerfile(config models.DeploymentConfig) string {
	return fmt.Sprintf(`FROM %s

WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy application
COPY . .

# Build application
RUN go build -o main .

# Expose port
EXPOSE %d

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:%d/health || exit 1

# Run application
CMD ["./main"]
`, config.GetBaseImage(), config.Port, config.Port)
}

// getDockerCompose generates docker-compose.yml
func getDockerCompose(config models.DeploymentConfig) string {
	return fmt.Sprintf(`version: '3.8'

services:
  %s:
    build:
      context: .
      dockerfile: %s
    ports:
      - "%d:%d"
    environment:
      - ENV=%s
      - PORT=%d
    volumes:
      - ./config:/app/config:ro
    depends_on:
      - db
      - redis
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:%d/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=%s
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  default:
    name: %s-network
`, config.AppName, config.GetDockerFileName(), config.Port, config.Port, 
   config.Environment, config.Port, config.Port, config.AppName, config.AppName)
}

// getDockerComposeProduction generates production docker-compose.yml
func getDockerComposeProduction(config models.DeploymentConfig) string {
	return fmt.Sprintf(`version: '3.8'

services:
  %s:
    image: ${DOCKER_REGISTRY}/%s:${VERSION}
    ports:
      - "%d:%d"
    environment:
      - ENV=production
      - PORT=%d
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
    volumes:
      - ./config:/app/config:ro
    restart: always
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:%d/health"]
      interval: 30s
      timeout: 10s
      retries: 5
      start_period: 60s
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
        window: 120s

networks:
  default:
    name: %s-production-network
    external: true
`, config.AppName, config.AppName, config.Port, config.Port, config.Port, 
   config.Port, config.AppName)
}

// getDockerIgnore generates .dockerignore
func getDockerIgnore() string {
	return `# Git
.git
.gitignore

# Documentation
README.md
docs/
*.md

# Tests
*_test.go
test/
coverage.out

# Build artifacts
*.exe
*.exe~
*.dll
*.so
*.dylib

# Development
.env
.env.local
.vscode/
.idea/

# Dependencies
vendor/

# Temporary files
*.tmp
*.log
.DS_Store

# Deployment
deployment/
k8s/
terraform/
.github/

# Node modules (if any)
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
`
}

// GetKubernetesTemplate returns Kubernetes template based on file type
func GetKubernetesTemplate(fileName string, config models.DeploymentConfig) string {
	switch fileName {
	case "namespace.yaml":
		return getKubernetesNamespace(config)
	case "deployment.yaml":
		return getKubernetesDeployment(config)
	case "service.yaml":
		return getKubernetesService(config)
	case "ingress.yaml":
		return getKubernetesIngress(config)
	case "configmap.yaml":
		return getKubernetesConfigMap(config)
	case "secret.yaml":
		return getKubernetesSecret(config)
	case "hpa.yaml":
		return getKubernetesHPA(config)
	default:
		return ""
	}
}

// getKubernetesNamespace generates namespace manifest
func getKubernetesNamespace(config models.DeploymentConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Namespace
metadata:
  name: %s
  labels:
    name: %s
    environment: %s
`, config.Namespace, config.Namespace, config.Environment)
}

// getKubernetesDeployment generates deployment manifest
func getKubernetesDeployment(config models.DeploymentConfig) string {
	labels := config.GetLabels()
	envVars := config.GetEnvironmentVariables()
	requests := config.GetResourceRequests()
	limits := config.GetResourceLimits()

	return fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
    version: %s
    environment: %s
spec:
  replicas: %d
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
        version: %s
        environment: %s
    spec:
      containers:
      - name: %s
        image: %s:%s
        ports:
        - containerPort: %d
          name: http
        env:
        - name: ENV
          value: %s
        - name: PORT
          value: "%d"
        resources:
          requests:
            memory: %s
            cpu: %s
          limits:
            memory: %s
            cpu: %s
        livenessProbe:
          httpGet:
            path: /health
            port: %d
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: %d
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          runAsUser: 1001
          capabilities:
            drop:
            - ALL
          readOnlyRootFilesystem: true
      securityContext:
        fsGroup: 1001
`,
		config.AppName, config.Namespace, labels["app"], labels["version"], labels["environment"],
		config.Replicas, labels["app"], labels["app"], labels["version"], labels["environment"],
		config.AppName, config.AppName, config.Version, config.Port,
		envVars["ENV"], config.Port, requests["memory"], requests["cpu"],
		limits["memory"], limits["cpu"], config.Port, config.Port)
}

// getKubernetesService generates service manifest
func getKubernetesService(config models.DeploymentConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  name: %s-service
  namespace: %s
  labels:
    app: %s
    environment: %s
spec:
  selector:
    app: %s
  ports:
  - name: http
    port: 80
    targetPort: %d
    protocol: TCP
  type: ClusterIP
`, config.AppName, config.Namespace, config.AppName, config.Environment,
   config.AppName, config.Port)
}

// getKubernetesIngress generates ingress manifest
func getKubernetesIngress(config models.DeploymentConfig) string {
	return fmt.Sprintf(`apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: %s-ingress
  namespace: %s
  labels:
    app: %s
    environment: %s
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - %s.example.com
    secretName: %s-tls
  rules:
  - host: %s.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: %s-service
            port:
              number: 80
`, config.AppName, config.Namespace, config.AppName, config.Environment,
   config.AppName, config.AppName, config.AppName, config.AppName)
}

// getKubernetesConfigMap generates configmap manifest
func getKubernetesConfigMap(config models.DeploymentConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: ConfigMap
metadata:
  name: %s-config
  namespace: %s
  labels:
    app: %s
    environment: %s
data:
  app.env: |
    ENV=%s
    PORT=%d
    LOG_LEVEL=info
    DATABASE_HOST=postgres-service
    REDIS_HOST=redis-service
`, config.AppName, config.Namespace, config.AppName, config.Environment,
   config.Environment, config.Port)
}

// getKubernetesSecret generates secret manifest
func getKubernetesSecret(config models.DeploymentConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  name: %s-secret
  namespace: %s
  labels:
    app: %s
    environment: %s
type: Opaque
data:
  # Base64 encoded values
  # Use: echo -n 'your-secret' | base64
  database-password: cG9zdGdyZXM=  # postgres
  jwt-secret: eW91ci1qd3Qtc2VjcmV0  # your-jwt-secret
  api-key: eW91ci1hcGkta2V5  # your-api-key
`, config.AppName, config.Namespace, config.AppName, config.Environment)
}

// getKubernetesHPA generates horizontal pod autoscaler manifest
func getKubernetesHPA(config models.DeploymentConfig) string {
	return fmt.Sprintf(`apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: %s-hpa
  namespace: %s
  labels:
    app: %s
    environment: %s
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: %s
  minReplicas: %d
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
`, config.AppName, config.Namespace, config.AppName, config.Environment,
   config.AppName, config.Replicas)
}

// GetCloudTemplate returns cloud provider template based on file and provider
func GetCloudTemplate(fileName string, config models.DeploymentConfig) string {
	switch config.CloudProvider {
	case models.AWSProvider:
		return getAWSTemplate(fileName, config)
	case models.GCPProvider:
		return getGCPTemplate(fileName, config)
	case models.AzureProvider:
		return getAzureTemplate(fileName, config)
	default:
		return ""
	}
}

// getAWSTemplate returns AWS-specific templates
func getAWSTemplate(fileName string, config models.DeploymentConfig) string {
	switch fileName {
	case "ecs-task-definition.json":
		return getECSTaskDefinition(config)
	case "terraform/main.tf":
		return getAWSTerraform(config)
	case "cloudformation/template.yaml":
		return getCloudFormationTemplate(config)
	default:
		return ""
	}
}

// getECSTaskDefinition generates ECS task definition
func getECSTaskDefinition(config models.DeploymentConfig) string {
	return fmt.Sprintf(`{
  "family": "%s",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "%s",
      "image": "ACCOUNT.dkr.ecr.%s.amazonaws.com/%s:%s",
      "portMappings": [
        {
          "containerPort": %d,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "ENV",
          "value": "%s"
        },
        {
          "name": "PORT",
          "value": "%d"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/%s",
          "awslogs-region": "%s",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "healthCheck": {
        "command": ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:%d/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3,
        "startPeriod": 60
      },
      "essential": true
    }
  ]
}`, config.AppName, config.AppName, config.Region, config.AppName, config.Version,
    config.Port, config.Environment, config.Port, config.AppName, config.Region, config.Port)
}

// getAWSTerraform generates AWS Terraform configuration
func getAWSTerraform(config models.DeploymentConfig) string {
	return fmt.Sprintf(`terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "%s"
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "%s-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name        = "%s-cluster"
    Environment = "%s"
  }
}

# ECS Task Definition
resource "aws_ecs_task_definition" "main" {
  family                   = "%s"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "%s"
      image = "${aws_ecr_repository.main.repository_url}:%s"
      
      portMappings = [
        {
          containerPort = %d
          hostPort      = %d
        }
      ]

      environment = [
        {
          name  = "ENV"
          value = "%s"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.main.name
          "awslogs-region"        = "%s"
          "awslogs-stream-prefix" = "ecs"
        }
      }

      essential = true
    }
  ])

  tags = {
    Name        = "%s-task"
    Environment = "%s"
  }
}

# ECR Repository
resource "aws_ecr_repository" "main" {
  name                 = "%s"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name        = "%s-ecr"
    Environment = "%s"
  }
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "main" {
  name              = "/ecs/%s"
  retention_in_days = 14

  tags = {
    Name        = "%s-logs"
    Environment = "%s"
  }
}
`, config.Region, config.AppName, config.AppName, config.Environment,
   config.AppName, config.AppName, config.Version, config.Port, config.Port,
   config.Environment, config.Region, config.AppName, config.Environment,
   config.AppName, config.AppName, config.Environment, config.AppName,
   config.AppName, config.Environment)
}

// getCloudFormationTemplate generates CloudFormation template
func getCloudFormationTemplate(config models.DeploymentConfig) string {
	return fmt.Sprintf(`AWSTemplateFormatVersion: '2010-09-09'
Description: 'CloudFormation template for %s'

Parameters:
  ImageURI:
    Type: String
    Description: ECR Image URI
    Default: "ACCOUNT.dkr.ecr.%s.amazonaws.com/%s:latest"
  
  Environment:
    Type: String
    Default: %s
    AllowedValues: [development, staging, production]

Resources:
  # ECS Cluster
  ECSCluster:
    Type: AWS::ECS::Cluster
    Properties:
      ClusterName: !Sub "${AWS::StackName}-cluster"
      CapacityProviders:
        - FARGATE
      DefaultCapacityProviderStrategy:
        - CapacityProvider: FARGATE
          Weight: 1

  # Task Definition
  TaskDefinition:
    Type: AWS::ECS::TaskDefinition
    Properties:
      Family: !Sub "${AWS::StackName}-task"
      Cpu: 256
      Memory: 512
      NetworkMode: awsvpc
      RequiresCompatibilities:
        - FARGATE
      ExecutionRoleArn: !Ref TaskExecutionRole
      TaskRoleArn: !Ref TaskRole
      ContainerDefinitions:
        - Name: %s
          Image: !Ref ImageURI
          PortMappings:
            - ContainerPort: %d
          Environment:
            - Name: ENV
              Value: !Ref Environment
          LogConfiguration:
            LogDriver: awslogs
            Options:
              awslogs-group: !Ref LogGroup
              awslogs-region: !Ref AWS::Region
              awslogs-stream-prefix: ecs

  # Log Group
  LogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: !Sub "/ecs/${AWS::StackName}"
      RetentionInDays: 14

Outputs:
  ClusterName:
    Description: Name of the ECS cluster
    Value: !Ref ECSCluster
    Export:
      Name: !Sub "${AWS::StackName}-cluster-name"
`, config.AppName, config.Region, config.AppName, config.Environment,
   config.AppName, config.Port)
}

// GetCICDTemplate returns CI/CD template based on provider
func GetCICDTemplate(fileName string, config models.DeploymentConfig) string {
	if strings.Contains(fileName, "github") {
		return getGitHubActionsTemplate(fileName, config)
	} else if strings.Contains(fileName, "gitlab") {
		return getGitLabCITemplate(config)
	} else if strings.Contains(fileName, "jenkins") {
		return getJenkinsTemplate(config)
	}
	return ""
}

// getGitHubActionsTemplate generates GitHub Actions workflow
func getGitHubActionsTemplate(fileName string, config models.DeploymentConfig) string {
	if strings.Contains(fileName, "ci.yml") {
		return getGitHubActionsCIWorkflow(config)
	} else if strings.Contains(fileName, "cd.yml") {
		return getGitHubActionsCDWorkflow(config)
	} else if strings.Contains(fileName, "security.yml") {
		return getGitHubActionsSecurityWorkflow(config)
	}
	return ""
}

// getGitHubActionsCIWorkflow generates CI workflow
func getGitHubActionsCIWorkflow(config models.DeploymentConfig) string {
	return fmt.Sprintf(`name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  GO_VERSION: '1.21'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
        
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
          
    - name: Download dependencies
      run: go mod download
      
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
      
    - name: Generate coverage report
      run: go tool cover -html=coverage.out -o coverage.html
      
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
        
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  build:
    needs: [test, lint]
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
        
    - name: Build
      run: go build -v ./...
`)
}

// getGitHubActionsCDWorkflow generates CD workflow
func getGitHubActionsCDWorkflow(config models.DeploymentConfig) string {
	return fmt.Sprintf(`name: CD

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    
    steps:
    - name: Checkout
      uses: actions/checkout@v4
      
    - name: Log in to Container Registry
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
        
    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=semver,pattern={{version}}
          type=semver,pattern={{major}}.{{minor}}
          
    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: .
        file: ./deployment/docker/%s
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    
    steps:
    - name: Deploy to production
      run: |
        echo "Deploying %s to production..."
        # Add your deployment commands here
        # kubectl apply -f deployment/kubernetes/
        # or
        # aws ecs update-service --cluster %s-cluster --service %s-service --force-new-deployment
`, config.GetDockerFileName(), config.AppName, config.AppName, config.AppName)
}

// getGitHubActionsSecurityWorkflow generates security workflow
func getGitHubActionsSecurityWorkflow(config models.DeploymentConfig) string {
	return `name: Security

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 2 * * 1'  # Weekly on Monday at 2 AM

jobs:
  gosec:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
        
    - name: Run Gosec Security Scanner
      uses: securecodewarrior/github-action-gosec@master
      with:
        args: '-fmt sarif -out gosec.sarif ./...'
        
    - name: Upload SARIF file
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: gosec.sarif

  trivy:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
      
    - name: Build image
      run: docker build -t test-image .
      
    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        image-ref: 'test-image'
        format: 'sarif'
        output: 'trivy-results.sarif'
        
    - name: Upload Trivy scan results to GitHub Security tab
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: 'trivy-results.sarif'
`
}

// getGitLabCITemplate generates GitLab CI configuration
func getGitLabCITemplate(config models.DeploymentConfig) string {
	return fmt.Sprintf(`stages:
  - test
  - build
  - security
  - deploy

variables:
  GO_VERSION: "1.21"
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"

before_script:
  - apt-get update -qq && apt-get install -y -qq git ca-certificates

test:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go mod download
    - go test -v -race -coverprofile=coverage.out ./...
    - go tool cover -func=coverage.out
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE/%s:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE/%s:$CI_COMMIT_SHA
  only:
    - main
    - develop

security:
  stage: security
  image: golang:${GO_VERSION}
  script:
    - go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    - gosec -fmt json -out gosec-report.json ./...
  artifacts:
    reports:
      sast: gosec-report.json
  allow_failure: true

deploy:production:
  stage: deploy
  image: alpine:latest
  script:
    - echo "Deploying %s to production"
    # Add deployment commands here
  environment:
    name: production
    url: https://%s.example.com
  only:
    - main
  when: manual
`, config.AppName, config.AppName, config.AppName, config.AppName)
}

// getJenkinsTemplate generates Jenkinsfile
func getJenkinsTemplate(config models.DeploymentConfig) string {
	return fmt.Sprintf(`pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.21'
        DOCKER_REGISTRY = 'your-registry.com'
        APP_NAME = '%s'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Test') {
            steps {
                sh '''
                    go version
                    go mod download
                    go test -v -race -coverprofile=coverage.out ./...
                    go tool cover -func=coverage.out
                '''
            }
            post {
                always {
                    publishCoverage adapters: [goAdapter('coverage.out')], sourceFileResolver: sourceFiles('STORE_LAST_BUILD')
                }
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                    go build -v -o ${APP_NAME} ./...
                '''
            }
        }
        
        stage('Security Scan') {
            steps {
                sh '''
                    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
                    gosec -fmt json -out gosec-report.json ./...
                '''
            }
            post {
                always {
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: '.',
                        reportFiles: 'gosec-report.json',
                        reportName: 'Security Report'
                    ])
                }
            }
        }
        
        stage('Docker Build') {
            when {
                branch 'main'
            }
            steps {
                script {
                    def image = docker.build("${DOCKER_REGISTRY}/${APP_NAME}:${BUILD_NUMBER}")
                    docker.withRegistry('https://${DOCKER_REGISTRY}', 'docker-registry-credentials') {
                        image.push()
                        image.push('latest')
                    }
                }
            }
        }
        
        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                script {
                    echo "Deploying ${APP_NAME} to production..."
                    // Add deployment commands here
                    // sh 'kubectl apply -f deployment/kubernetes/'
                }
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}
`, config.AppName)
}

// Helper functions for GCP and Azure templates (simplified for brevity)
func getGCPTemplate(fileName string, config models.DeploymentConfig) string {
	switch fileName {
	case "cloud-run.yaml":
		return fmt.Sprintf(`apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: %s
  annotations:
    run.googleapis.com/ingress: all
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "10"
        run.googleapis.com/cpu-throttling: "false"
    spec:
      containerConcurrency: 100
      containers:
      - image: gcr.io/PROJECT_ID/%s:%s
        ports:
        - containerPort: %d
        env:
        - name: ENV
          value: %s
        resources:
          limits:
            cpu: 1000m
            memory: 512Mi
`, config.AppName, config.AppName, config.Version, config.Port, config.Environment)
	default:
		return ""
	}
}

func getAzureTemplate(fileName string, config models.DeploymentConfig) string {
	switch fileName {
	case "container-instances.json":
		return fmt.Sprintf(`{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "containerGroupName": {
      "type": "string",
      "defaultValue": "%s"
    }
  },
  "resources": [
    {
      "type": "Microsoft.ContainerInstance/containerGroups",
      "apiVersion": "2021-03-01",
      "name": "[parameters('containerGroupName')]",
      "location": "[resourceGroup().location]",
      "properties": {
        "containers": [
          {
            "name": "%s",
            "properties": {
              "image": "your-registry.azurecr.io/%s:%s",
              "ports": [
                {
                  "port": %d,
                  "protocol": "TCP"
                }
              ],
              "environmentVariables": [
                {
                  "name": "ENV",
                  "value": "%s"
                }
              ],
              "resources": {
                "requests": {
                  "cpu": 0.5,
                  "memoryInGb": 1
                }
              }
            }
          }
        ],
        "osType": "Linux",
        "ipAddress": {
          "type": "Public",
          "ports": [
            {
              "port": %d,
              "protocol": "TCP"
            }
          ]
        }
      }
    }
  ]
}`, config.AppName, config.AppName, config.AppName, config.Version, 
    config.Port, config.Environment, config.Port)
	default:
		return ""
	}
}