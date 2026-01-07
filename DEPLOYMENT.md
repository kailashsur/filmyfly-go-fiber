# GCP Deployment Guide - Quick Start

## Prerequisites
- Google Cloud account
- `gcloud` CLI installed and authenticated
- GitHub repository with your code

## Step-by-Step Deployment

### 1. Create GCP VM
```bash
gcloud compute instances create filmyfly-server \
    --zone=us-central1-a \
    --machine-type=e2-medium \
    --boot-disk-size=30GB \
    --image-family=ubuntu-2204-lts \
    --image-project=ubuntu-os-cloud \
    --tags=http-server,https-server \
    --metadata-from-file=startup-script=vm-startup.sh
```

### 2. Configure Firewall
```bash
gcloud compute firewall-rules create allow-http --allow tcp:80 --target-tags http-server
gcloud compute firewall-rules create allow-https --allow tcp:443 --target-tags https-server
```

### 3. SSH into VM and Setup
```bash
gcloud compute ssh filmyfly-server --zone=us-central1-a

# Setup PostgreSQL
sudo -u postgres psql
CREATE DATABASE filmyfly;
CREATE USER filmyfly_user WITH ENCRYPTED PASSWORD 'YOUR_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE filmyfly TO filmyfly_user;
\q

# Restore database
gcloud compute scp ~/path/to/backup.dump filmyfly-server:~/backup.dump --zone=us-central1-a
sudo -u postgres pg_restore -d filmyfly -v ~/backup.dump
```

### 4. Deploy Application
```bash
# Clone repository
cd /opt/filmyfly
git clone https://github.com/YOUR_USERNAME/filmyfly-go-fiber.git .

# Create .env file
sudo nano /opt/filmyfly/.env
# Add your environment variables

# Build application
go build -o filmyfly cmd/server/main.go

# Setup systemd service
sudo cp filmyfly.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable filmyfly
sudo systemctl start filmyfly

# Setup Nginx
sudo cp nginx.conf /etc/nginx/sites-available/filmyfly
sudo ln -s /etc/nginx/sites-available/filmyfly /etc/nginx/sites-enabled/
sudo rm /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl restart nginx
```

### 5. Setup Auto-Deployment
```bash
# Connect GitHub to Cloud Build
gcloud builds triggers create github \
    --repo-name=filmyfly-go-fiber \
    --repo-owner=YOUR_GITHUB_USERNAME \
    --branch-pattern="^main$" \
    --build-config=cloudbuild.yaml
```

### 6. Get Your Application URL
```bash
gcloud compute instances describe filmyfly-server \
    --zone=us-central1-a \
    --format='get(networkInterfaces[0].accessConfigs[0].natIP)'
```

Visit `http://YOUR_VM_IP` to see your application!

## Useful Commands

### Check Application Status
```bash
sudo systemctl status filmyfly
```

### View Logs
```bash
sudo journalctl -u filmyfly -f
```

### Manual Deployment
```bash
./deploy.sh
```

For detailed instructions, see `gcp_deployment_plan.md`
