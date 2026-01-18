# Helios Statement Extraction API - Setup Guide

## Current Status

✅ **Application is running successfully on port 8080**

The API has been built and deployed. All endpoints are operational, but Gemini LLM integration requires GCP credentials to be configured.

## Quick Test

### Test 1: PDF Text Extraction (Working Now)

```bash
curl -X POST http://localhost:8080/v1/statements/extract-text \
  -F "file=@src/main/resources/statement.pdf"
```

This endpoint extracts raw text from the PDF and returns it as plain text.

### Test 2: Statement Parsing with LLM (Requires GCP Setup)

```bash
curl -X POST http://localhost:8080/v1/statements/extract \
  -F "file=@src/main/resources/statement.pdf" \
  | jq .
```

This endpoint extracts and parses the statement into structured JSON using Gemini LLM.

## Setup Google Cloud & Gemini

To use the full statement parsing feature, you need to configure Google Cloud credentials:

### Step 1: Create a GCP Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Note your Project ID

### Step 2: Enable Vertex AI API

```bash
gcloud services enable aiplatform.googleapis.com --project=YOUR_PROJECT_ID
```

Or enable it in the [Google Cloud Console](https://console.cloud.google.com/apis/library/aiplatform.googleapis.com)

### Step 3: Create Service Account

```bash
# Create service account
gcloud iam service-accounts create helios-service-account \
  --display-name="Helios Statement Parser" \
  --project=YOUR_PROJECT_ID

# Grant Vertex AI User role
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:helios-service-account@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/aiplatform.user"

# Create and download key
gcloud iam service-accounts keys create ~/helios-gcp-key.json \
  --iam-account=helios-service-account@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

### Step 4: Set Environment Variables

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/Users/YOUR_USERNAME/helios-gcp-key.json"
export GEMINI_PROJECT_ID="YOUR_PROJECT_ID"
export GEMINI_LOCATION="us-central1"
export GEMINI_MODEL_NAME="gemini-1.5-flash"
```

### Step 5: Restart the Application

```bash
# Stop current application
pkill -f helios-0.0.1-SNAPSHOT.jar

# Start with environment variables
cd /Users/11411733/git/code/helios
java -jar target/helios-0.0.1-SNAPSHOT.jar
```

Or use Maven:

```bash
mvn spring-boot:run
```

## Alternative: Update application.yaml

Instead of environment variables, you can directly update `src/main/resources/application.yaml`:

```yaml
gemini:
  project-id: your-actual-project-id
  location: us-central1
  model-name: gemini-1.5-flash
```

Then rebuild and restart:

```bash
mvn clean package -DskipTests
java -jar target/helios-0.0.1-SNAPSHOT.jar
```

## Verify Setup

Once configured, test the full parsing endpoint:

```bash
curl -X POST http://localhost:8080/v1/statements/extract \
  -F "file=@src/main/resources/statement.pdf" \
  | jq .
```

Expected response (example):

```json
{
  "cardNumber": "5468 48XX XXXX 8032",
  "statementDate": "2025-08-05",
  "creditLine": 600000.00,
  "totalPaymentDue": 21367.39,
  "outstandingBalance": 38344.01,
  "minimumPaymentDue": 10942.35,
  "availableCredit": 560965.99,
  "availableCreditLimit": 600000.00,
  "paymentDueDate": "2025-08-25",
  "lastPaymentDate": "2025-07-28",
  "transactions": [
    {
      "transactionDate": "2025-07-05",
      "postingDate": "2025-07-06",
      "description": "OMISE.CO BANGKOK",
      "amount": 690.00,
      "isInstallment": false,
      "installmentCurrent": null,
      "installmentTotal": null,
      "installmentPlanAmount": null
    }
    // ... more transactions
  ]
}
```

## Troubleshooting

### Port 8080 Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Authentication Errors

If you see authentication errors:

1. Verify `GOOGLE_APPLICATION_CREDENTIALS` points to valid JSON key file
2. Ensure service account has `roles/aiplatform.user` role
3. Verify Vertex AI API is enabled in your project
4. Check the service account key file has correct permissions

### Maven Build Errors

```bash
# Clean and rebuild
mvn clean install

# Skip tests if needed
mvn clean package -DskipTests
```

## Application Management

### Check if Application is Running

```bash
curl http://localhost:8080/actuator/health 2>/dev/null || echo "Not running"
```

### Stop the Application

```bash
pkill -f helios-0.0.1-SNAPSHOT.jar
```

### View Logs

The application logs to console. To save logs:

```bash
java -jar target/helios-0.0.1-SNAPSHOT.jar > helios.log 2>&1 &
tail -f helios.log
```

## Next Steps

1. Configure GCP credentials (see above)
2. Test the `/extract` endpoint with your statement PDFs
3. Integrate the API with your frontend or other services
4. Review the [API_USAGE.md](API_USAGE.md) for detailed endpoint documentation

## Support

For issues or questions:
- Check application logs for detailed error messages
- Verify all environment variables are set correctly
- Ensure GCP billing is enabled for your project
- Review Gemini API quota limits in GCP Console
