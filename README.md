# go-service

## Steps to build and run the service

1. Clone the github repository and change directory to the project's root folder: 
```go
git clone https://github.com/haardikdharma10/go-service.git

cd go-service
```

2. The project contains a `.devcontainer` directory which has the specifications for creating the container and a Dockerfile. If using VSCode, would suggest installing the [DevContainer](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension. When prompted, click "Reopen in Container" in VS Code.

3. VS Code will build the Dev Container and open the project inside it. This may take a few minutes the first time.

4. Once inside the container, before doing anything else, make sure to add your Github token in `svc/scan.go` which will be required to make the API calls to Github. 
```go
tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "<token>"},
	)
```

5. Build the binary using the following command:
```go
make build
```

6. Run the service:
```go
make run
```

7. The service will be running on port 8080. There are 2 endpoints to be tested:
- POST `/scan`: Scan a GitHub repository
- POST `/query`: Query vulnerabilities

8. Use Postman or curl to send API requests to the service.

```
curl -X POST -H "Content-Type: application/json" -d '{"owner":"velancio","repo":"vulnerability_scans"}' http://localhost:8080/scan
```

This returns the following response:
```json
{
  "total_files": 9,
  "processed_files": 9,
  "total_payloads": 26,
  "payloads": [
    {
      "source_file": "vulnscan19.json",
      "scan_time": "2025-02-12T21:23:13.609454-05:00",
      "payload_result": {
        "scanResults": {
          "resource_name": "ml-inference:2.0.0",
          "resource_type": "container",
          "scan_id": "VULN_scan_345mno",
          "scan_metadata": {
            "excluded_paths": [
              "/tmp",
              "/var/log"
            ],
            "policies_version": "2025.1.29",
            "scanner_version": "30.1.51",
            "scanning_rules": [
              "vulnerability",
              "compliance",
              "malware"
            ]
          },
          "scan_status": "completed",
          "summary": {
            "compliant": false,
            "fixable_count": 2,
            "severity_counts": {
              "CRITICAL": 0,
              "HIGH": 1,
              "LOW": 0,
              "MEDIUM": 1
            },
            "total_vulnerabilities": 2
          },
          "timestamp": "2025-01-29T13:00:00Z",
          "vulnerabilities": [
            {
              "current_version": "2.7.0",
              "cvss": 8.5,
              "description": "Remote code execution in TensorFlow model loading",
              "fixed_version": "2.7.1",
              "id": "CVE-2024-5555",
              "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-5555",
              "package_name": "tensorflow",
              "published_date": "2025-01-24T00:00:00Z",
              "risk_factors": [
                "Remote Code Execution",
                "High CVSS Score",
                "Public Exploit Available",
                "Exploit in Wild"
              ],
              "severity": "HIGH",
              "status": "active"
            },
            {
              "current_version": "1.0.0",
              "cvss": 6.7,
              "description": "Memory corruption in scikit-learn",
              "fixed_version": "1.0.1",
              "id": "CVE-2024-5556",
              "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-5556",
              "package_name": "scikit-learn",
              "published_date": "2025-01-25T00:00:00Z",
              "risk_factors": [
                "Memory Corruption",
                "Medium CVSS Score"
              ],
              "severity": "MEDIUM",
              "status": "active"
            }
          ]
        }
      }
    },
    .
    .
    .
    .
    }
  ]
}
```

9. A database will be created with 1 table `scans`. There are a couple ways to access the sqlite db:

- Use the sqlite3 command line tool:
```go
sqlite3 scans.db
```

- What I find simpler is to copy the db to the local machine:
```go
docker cp <container_id>:/workspace/scans.db ./scans.db
```

Then use something like [TablePlus](https://tableplus.com/) to view the database.

10. Now that we know our db has the data for vulnerabilities, we can query it using the POST `/query` endpoint. 

```
curl -X POST -H "Content-Type: application/json" -d '{"filters":{"severity":"LOW"}}' http://localhost:8080/query
```

This returns the following response:
```json
[
  {
    "id": "CVE-2024-8803",
    "severity": "LOW",
    "cvss": 3.2,
    "status": "active",
    "package_name": "nginx",
    "current_version": "1.20.1",
    "fixed_version": "1.20.2",
    "description": "Information disclosure in NGINX error logs",
    "published_date": "2025-01-26T00:00:00Z",
    "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-8803",
    "risk_factors": [
      "Information Disclosure",
      "Low CVSS Score"
    ]
  },
  {
    "id": "CVE-2024-3345",
    "severity": "LOW",
    "cvss": 3.5,
    "status": "fixed",
    "package_name": "busybox",
    "current_version": "1.33.0",
    "fixed_version": "1.33.1",
    "description": "Privilege escalation in BusyBox",
    "published_date": "2024-01-25T00:00:00Z",
    "link": "",
    "risk_factors": [
      "Privilege Escalation"
    ]
  },
  {
    "id": "CVE-2024-5557",
    "severity": "LOW",
    "cvss": 3.5,
    "status": "active",
    "package_name": "pandas",
    "current_version": "1.3.0",
    "fixed_version": "1.3.1",
    "description": "Information disclosure in pandas data frame handling",
    "published_date": "2025-01-26T00:00:00Z",
    "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-5557",
    "risk_factors": [
      "Information Disclosure",
      "Low CVSS Score"
    ]
  },
  {
    "id": "CVE-2024-7003",
    "severity": "LOW",
    "cvss": 3.2,
    "status": "active",
    "package_name": "moment",
    "current_version": "2.29.1",
    "fixed_version": "2.29.2",
    "description": "Regular expression denial of service in Moment.js",
    "published_date": "2025-01-26T00:00:00Z",
    "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-7003",
    "risk_factors": [
      "Regular Expression Denial of Service",
      "Low CVSS Score"
    ]
  },
  {
    "id": "CVE-2024-8903",
    "severity": "LOW",
    "cvss": 3.2,
    "status": "active",
    "package_name": "curl",
    "current_version": "7.74.0",
    "fixed_version": "7.74.1",
    "description": "Information disclosure in curl command line tool",
    "published_date": "2024-01-20T00:00:00Z",
    "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-8903",
    "risk_factors": [
      "Information Disclosure"
    ]
  },
  {
    "id": "CVE-2024-5568",
    "severity": "LOW",
    "cvss": 3.8,
    "status": "fixed",
    "package_name": "nodejs",
    "current_version": "16.13.0",
    "fixed_version": "16.13.1",
    "description": "Denial of service in Node.js HTTP parser",
    "published_date": "2024-01-29T00:00:00Z",
    "link": "",
    "risk_factors": [
      "Denial of Service"
    ]
  },
  {
    "id": "CVE-2024-3579",
    "severity": "LOW",
    "cvss": 3.2,
    "status": "fixed",
    "package_name": "curl",
    "current_version": "7.74.0",
    "fixed_version": "7.74.1",
    "description": "Information disclosure in curl command line tool",
    "published_date": "2024-01-05T00:00:00Z",
    "link": "",
    "risk_factors": [
      "Information Disclosure"
    ]
  },
  {
    "id": "CVE-2024-5557",
    "severity": "LOW",
    "cvss": 3.5,
    "status": "active",
    "package_name": "pandas",
    "current_version": "1.3.0",
    "fixed_version": "1.3.1",
    "description": "Information disclosure in pandas data frame handling",
    "published_date": "2025-01-26T00:00:00Z",
    "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-5557",
    "risk_factors": [
      "Information Disclosure",
      "Low CVSS Score"
    ]
  },
  {
    "id": "CVE-2024-8804",
    "severity": "LOW",
    "cvss": 3.2,
    "status": "active",
    "package_name": "nginx",
    "current_version": "1.20.1",
    "fixed_version": "1.20.2",
    "description": "Information disclosure in NGINX error logs",
    "published_date": "2025-01-25T00:00:00Z",
    "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-8804",
    "risk_factors": [
      "Information Disclosure",
      "Low CVSS Score"
    ]
  },
  {
    "id": "CVE-2024-9905",
    "severity": "LOW",
    "cvss": 3.7,
    "status": "active",
    "package_name": "logrotate",
    "current_version": "3.19.0",
    "fixed_version": "3.19.1",
    "description": "Information disclosure through improper log file permissions",
    "published_date": "2025-01-24T23:00:00Z",
    "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-9905",
    "risk_factors": [
      "Information Disclosure",
      "Low CVSS Score"
    ]
  }
]
```

11. To run tests:
```go
make test-coverage
```
Current coverage is at 72%, more exhaustive tests can be added to increase coverage.

12. Additional technical details:
- Concurrency: The service uses goroutines to handle multiple requests concurrently. The `sync.WaitGroup` is used to wait for all goroutines to finish before returning the response. Currently, the service can handle up to 3 concurrent requests. I tried to get the time difference betwwen sequential and concurrent requests as well and found out that sequential processing took 1.15 seconds to process 9 files, while concurrent processing took 0.43 seconds (~37% of the time taken by sequential processing) to process the same files.

- Error Handling: Failed calls to the GitHub API are retried 2 more times after the initial call.