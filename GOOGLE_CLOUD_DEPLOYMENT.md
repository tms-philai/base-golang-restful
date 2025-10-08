# Google Cloud Platform Deployment Guide

Hướng dẫn deploy Base Golang RESTful API lên Google Cloud Platform.

---

## 📋 Prerequisites

### 1. Google Cloud Account
- Tạo tài khoản Google Cloud: https://cloud.google.com
- Enable billing cho project
- Install Google Cloud SDK: https://cloud.google.com/sdk/docs/install

### 2. Required APIs
Enable các APIs sau trong Google Cloud Console:
```bash
gcloud services enable \
    cloudbuild.googleapis.com \
    run.googleapis.com \
    sqladmin.googleapis.com \
    storage-api.googleapis.com \
    container.googleapis.com \
    cloudresourcemanager.googleapis.com
```

### 3. Tools Installation
```bash
# Install gcloud CLI
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# Install kubectl for GKE
gcloud components install kubectl

# Login to Google Cloud
gcloud auth login
gcloud config set project YOUR_PROJECT_ID
```

---

## 🗄️ Database Setup - Cloud SQL

### 1. Create Cloud SQL Instance
```bash
# Create PostgreSQL instance
gcloud sql instances create postgres-instance \
    --database-version=POSTGRES_15 \
    --tier=db-f1-micro \
    --region=asia-northeast1 \
    --root-password=YOUR_ROOT_PASSWORD \
    --storage-type=SSD \
    --storage-size=10GB \
    --backup \
    --backup-start-time=03:00

# Create database
gcloud sql databases create app_db \
    --instance=postgres-instance

# Create user
gcloud sql users create app_user \
    --instance=postgres-instance \
    --password=YOUR_USER_PASSWORD
```

### 2. Configure Cloud SQL Proxy (Local Development)
```bash
# Download Cloud SQL Proxy
curl -o cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.darwin.amd64
chmod +x cloud_sql_proxy

# Run proxy
./cloud_sql_proxy -instances=PROJECT_ID:asia-northeast1:postgres-instance=tcp:5432
```

### 3. Connection String
```bash
# For Cloud Run/GKE
DB_HOST=/cloudsql/PROJECT_ID:asia-northeast1:postgres-instance
DB_PORT=5432

# For local with proxy
DB_HOST=127.0.0.1
DB_PORT=5432
```

---

## 💾 Storage Setup - Google Cloud Storage

### 1. Create Storage Bucket
```bash
# Create bucket for file uploads
gsutil mb -l asia-northeast1 gs://YOUR_PROJECT_ID-uploads

# Set bucket permissions
gsutil iam ch allUsers:objectViewer gs://YOUR_PROJECT_ID-uploads

# Enable CORS
cat > cors.json <<EOF
[
  {
    "origin": ["*"],
    "method": ["GET", "POST", "PUT", "DELETE"],
    "responseHeader": ["Content-Type"],
    "maxAgeSeconds": 3600
  }
]
EOF

gsutil cors set cors.json gs://YOUR_PROJECT_ID-uploads
```

### 2. Create Service Account
```bash
# Create service account
gcloud iam service-accounts create go-api-sa \
    --display-name="Go API Service Account"

# Grant permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/storage.objectAdmin"

gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/cloudsql.client"

# Create key file
gcloud iam service-accounts keys create credentials.json \
    --iam-account=go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

---

## 🐳 Option 1: Deploy to Cloud Run (Recommended for Simple Apps)

### 1. Build and Push Image
```bash
# Build with Cloud Build
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/go-api

# Or build locally and push
docker build -t gcr.io/YOUR_PROJECT_ID/go-api .
docker push gcr.io/YOUR_PROJECT_ID/go-api
```

### 2. Deploy to Cloud Run
```bash
gcloud run deploy go-api \
    --image gcr.io/YOUR_PROJECT_ID/go-api \
    --platform managed \
    --region asia-northeast1 \
    --allow-unauthenticated \
    --add-cloudsql-instances YOUR_PROJECT_ID:asia-northeast1:postgres-instance \
    --set-env-vars "ENV=production,DB_HOST=/cloudsql/YOUR_PROJECT_ID:asia-northeast1:postgres-instance,DB_PORT=5432,DB_NAME=app_db,DB_USER=app_user,GCS_BUCKET=YOUR_PROJECT_ID-uploads" \
    --set-secrets "DB_PASSWORD=db-password:latest,JWT_SECRET=jwt-secret:latest" \
    --service-account go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com \
    --memory 512Mi \
    --cpu 1 \
    --max-instances 10 \
    --min-instances 0 \
    --timeout 300
```

### 3. Setup Custom Domain (Optional)
```bash
# Map custom domain
gcloud run domain-mappings create \
    --service go-api \
    --domain api.yourdomain.com \
    --region asia-northeast1
```

### 4. Setup Cloud Build Trigger (CI/CD)
```yaml
# cloudbuild.yaml
steps:
  # Run tests
  - name: 'golang:1.23'
    args: ['go', 'test', './...']
    dir: 'app'
  
  # Build image
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/go-api:$COMMIT_SHA', '.']
  
  # Push image
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/go-api:$COMMIT_SHA']
  
  # Deploy to Cloud Run
  - name: 'gcr.io/cloud-builders/gcloud'
    args:
      - 'run'
      - 'deploy'
      - 'go-api'
      - '--image=gcr.io/$PROJECT_ID/go-api:$COMMIT_SHA'
      - '--region=asia-northeast1'
      - '--platform=managed'

images:
  - 'gcr.io/$PROJECT_ID/go-api:$COMMIT_SHA'
```

---

## ☸️ Option 2: Deploy to Google Kubernetes Engine (GKE)

### 1. Create GKE Cluster
```bash
# Create cluster
gcloud container clusters create go-api-cluster \
    --region asia-northeast1 \
    --num-nodes 2 \
    --machine-type n1-standard-1 \
    --enable-autoscaling \
    --min-nodes 1 \
    --max-nodes 5 \
    --enable-autorepair \
    --enable-autoupgrade \
    --workload-pool=YOUR_PROJECT_ID.svc.id.goog

# Get credentials
gcloud container clusters get-credentials go-api-cluster \
    --region asia-northeast1
```

### 2. Create Kubernetes Secrets
```bash
# Create database credentials secret
kubectl create secret generic cloudsql-db-credentials \
    --from-literal=username=app_user \
    --from-literal=password=YOUR_USER_PASSWORD

# Create JWT secret
kubectl create secret generic jwt-secret \
    --from-literal=secret=YOUR_JWT_SECRET

# Create GCS credentials
kubectl create secret generic gcs-credentials \
    --from-file=credentials.json=./credentials.json
```

### 3. Deploy Application
```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml

# Check deployment status
kubectl get pods
kubectl get services
kubectl get ingress
```

### 4. Setup Ingress with SSL
```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-api-ingress
  annotations:
    kubernetes.io/ingress.class: "gce"
    kubernetes.io/ingress.global-static-ip-name: "go-api-ip"
    networking.gke.io/managed-certificates: "go-api-cert"
spec:
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /*
        pathType: ImplementationSpecific
        backend:
          service:
            name: go-api-service
            port:
              number: 80
```

```bash
# Create static IP
gcloud compute addresses create go-api-ip --global

# Create managed certificate
cat > managed-cert.yaml <<EOF
apiVersion: networking.gke.io/v1
kind: ManagedCertificate
metadata:
  name: go-api-cert
spec:
  domains:
    - api.yourdomain.com
EOF

kubectl apply -f managed-cert.yaml
```

---

## 🔧 Environment Configuration

### 1. Create Secret Manager Secrets
```bash
# Create secrets in Secret Manager
echo -n "YOUR_DB_PASSWORD" | gcloud secrets create db-password --data-file=-
echo -n "YOUR_JWT_SECRET" | gcloud secrets create jwt-secret --data-file=-

# Grant access to service account
gcloud secrets add-iam-policy-binding db-password \
    --member="serviceAccount:go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding jwt-secret \
    --member="serviceAccount:go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"
```

### 2. Environment Variables
```bash
# .env.production
ENV=production
DB_HOST=/cloudsql/YOUR_PROJECT_ID:asia-northeast1:postgres-instance
DB_PORT=5432
DB_NAME=app_db
DB_USER=app_user
DB_SSL_MODE=require
GCS_PROJECT_ID=YOUR_PROJECT_ID
GCS_BUCKET=YOUR_PROJECT_ID-uploads
GCS_CREDENTIALS_FILE=/secrets/gcs/credentials.json
REDIS_HOST=10.0.0.3
REDIS_PORT=6379
LOG_LEVEL=info
```

---

## 📊 Monitoring & Logging

### 1. Setup Cloud Monitoring
```bash
# Enable monitoring
gcloud services enable monitoring.googleapis.com

# Create uptime check
gcloud monitoring uptime-checks create go-api-uptime \
    --display-name="Go API Uptime" \
    --resource-type=uptime-url \
    --monitored-resource=https://YOUR_CLOUD_RUN_URL/health
```

### 2. View Logs
```bash
# Cloud Run logs
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=go-api" \
    --limit 50 \
    --format json

# GKE logs
kubectl logs -f deployment/go-api

# Stream logs
gcloud logging tail "resource.type=cloud_run_revision"
```

### 3. Setup Alerts
```bash
# Create alert policy for high error rate
gcloud alpha monitoring policies create \
    --notification-channels=CHANNEL_ID \
    --display-name="High Error Rate" \
    --condition-display-name="Error rate > 5%" \
    --condition-threshold-value=0.05 \
    --condition-threshold-duration=300s
```

---

## 🔒 Security Best Practices

### 1. IAM Permissions
```bash
# Principle of least privilege
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/cloudsql.client"

# Remove default permissions
gcloud projects remove-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/editor"
```

### 2. Network Security
```bash
# Create VPC
gcloud compute networks create go-api-vpc \
    --subnet-mode=custom

# Create subnet
gcloud compute networks subnets create go-api-subnet \
    --network=go-api-vpc \
    --region=asia-northeast1 \
    --range=10.0.0.0/24

# Setup Cloud Armor (DDoS protection)
gcloud compute security-policies create go-api-policy \
    --description="Security policy for Go API"
```

### 3. Enable VPC Service Controls
```bash
# Create service perimeter
gcloud access-context-manager perimeters create go-api-perimeter \
    --title="Go API Perimeter" \
    --resources=projects/YOUR_PROJECT_ID \
    --restricted-services=storage.googleapis.com,sqladmin.googleapis.com
```

---

## 💰 Cost Optimization

### 1. Cloud Run
- Use `--min-instances=0` để scale to zero
- Set `--max-instances` để control costs
- Use `--cpu-throttling` để reduce costs
- Enable request timeout

### 2. Cloud SQL
- Use `db-f1-micro` or `db-g1-small` cho development
- Enable automatic storage increase
- Setup automated backups
- Use read replicas cho high traffic

### 3. Cloud Storage
- Use Standard storage class
- Setup lifecycle policies
- Enable CDN caching
- Compress files before upload

### 4. Monitoring Costs
```bash
# View current costs
gcloud billing accounts list
gcloud billing projects describe YOUR_PROJECT_ID

# Set budget alerts
gcloud billing budgets create \
    --billing-account=BILLING_ACCOUNT_ID \
    --display-name="Monthly Budget" \
    --budget-amount=100USD \
    --threshold-rule=percent=50 \
    --threshold-rule=percent=90
```

---

## 🚀 CI/CD Pipeline

### 1. Setup Cloud Build Trigger
```bash
# Connect GitHub repository
gcloud beta builds triggers create github \
    --repo-name=YOUR_REPO \
    --repo-owner=YOUR_GITHUB_USERNAME \
    --branch-pattern="^main$" \
    --build-config=cloudbuild.yaml
```

### 2. Multi-environment Deployment
```yaml
# cloudbuild-staging.yaml
steps:
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/go-api:staging', '.']
  
  - name: 'gcr.io/cloud-builders/gcloud'
    args:
      - 'run'
      - 'deploy'
      - 'go-api-staging'
      - '--image=gcr.io/$PROJECT_ID/go-api:staging'
      - '--region=asia-northeast1'

# cloudbuild-production.yaml  
steps:
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/go-api:production', '.']
  
  - name: 'gcr.io/cloud-builders/gcloud'
    args:
      - 'run'
      - 'deploy'
      - 'go-api'
      - '--image=gcr.io/$PROJECT_ID/go-api:production'
      - '--region=asia-northeast1'
```

---

## 🔧 Troubleshooting

### Common Issues

#### 1. Cloud SQL Connection Failed
```bash
# Check Cloud SQL status
gcloud sql instances describe postgres-instance

# Test connection with proxy
./cloud_sql_proxy -instances=PROJECT_ID:asia-northeast1:postgres-instance=tcp:5432

# Check service account permissions
gcloud projects get-iam-policy YOUR_PROJECT_ID \
    --flatten="bindings[].members" \
    --filter="bindings.members:serviceAccount:go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com"
```

#### 2. Cloud Run Deployment Failed
```bash
# Check build logs
gcloud builds list --limit=5
gcloud builds log BUILD_ID

# Check service logs
gcloud run services describe go-api --region=asia-northeast1
gcloud logging read "resource.type=cloud_run_revision" --limit=50
```

#### 3. GCS Permission Denied
```bash
# Check bucket permissions
gsutil iam get gs://YOUR_PROJECT_ID-uploads

# Grant permissions
gsutil iam ch serviceAccount:go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com:objectAdmin \
    gs://YOUR_PROJECT_ID-uploads
```

---

## 📚 Additional Resources

- [Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Cloud SQL Documentation](https://cloud.google.com/sql/docs)
- [GKE Documentation](https://cloud.google.com/kubernetes-engine/docs)
- [Cloud Storage Documentation](https://cloud.google.com/storage/docs)
- [Cloud Build Documentation](https://cloud.google.com/build/docs)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)

---

## 🎯 Quick Start Commands

```bash
# Setup project
export PROJECT_ID=your-project-id
export REGION=asia-northeast1
gcloud config set project $PROJECT_ID

# Enable APIs
gcloud services enable cloudbuild.googleapis.com run.googleapis.com sqladmin.googleapis.com storage-api.googleapis.com

# Build and deploy
gcloud builds submit --tag gcr.io/$PROJECT_ID/go-api
gcloud run deploy go-api --image gcr.io/$PROJECT_ID/go-api --region $REGION --allow-unauthenticated

# View logs
gcloud logging tail "resource.type=cloud_run_revision AND resource.labels.service_name=go-api"
```

---

**Happy Deploying! 🚀**
