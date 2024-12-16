# 🐎 **Horseq**

**Horseq** is a data pipeline designed to extract, normalize, and transform blockchain transaction data, store it in BigQuery, and visualize the results. It automates data ingestion from CSV files stored in Google Cloud Storage (GCS), processes the data using Go, and loads the results into BigQuery for analytics.

## **Table of Contents**
---

1. [Prerequisites](#prerequisites)
2. [Setup](#setup)
3. [Usage](#usage)
4. [License](#license)

---

## **Prerequisites**

Before setting up the project, ensure you have the following:

1. **Go Programming Language** (version 1.16 or higher)  
   Download: [https://golang.org/dl/](https://golang.org/dl/)

2. **GCP Project** with the following services enabled:
   - **BigQuery**
   - **Cloud Storage**


## **Setup**

### 1. **Clone the Repository**

```bash
git clone https://github.com/Sundbassen/horseq.git
cd horseq
```

### 2. **Set Up Environment Variables**

Create a `.env` file in the project root directory with the following variables:

```env
GOOGLE_APPLICATION_CREDENTIALS="/path/to/your-service-account.json"
```

### 3. **Create BigQuery Resources**

Create a BigQuery dataset and table using the provided SQL script.

```bash
bq mk --dataset --location=us your-gcp-project-id:your-dataset-name

bq query --use_legacy_sql=false < scripts/transactions.sql
```

### 4. **Upload Sample CSV to GCS**

Upload the sample CSV file to your GCS bucket:

```bash
gsutil cp /path/to/sample_data.csv gs://your-gcs-bucket-name/sample_data.csv
```

---

## **Usage**

### **Run the ETL Process**

To execute the ETL process and load data into BigQuery, run the following command:

```bash
1. cd cmd/horseq 
2. go build 
3. ./horseq -h and follow the help or simply: ./horseq datapipeline -p your-gcp-project-id -b your-bucket-name -c path-to-csv-in-bucket
```

This will:

1. Read rows of the csv file from GCS.
2. Fetch currency conversion rates from CoinGecko.
3. Insert the transformed data into BigQuery.

---

## **License**

This project is licensed under the [MIT License](LICEN1SE).

---
