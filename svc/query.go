// Package svc contains the core logic for the service.
package svc

import (
	"encoding/json"
	"log"

	"example/go-service/db"
	"example/go-service/models"
)

// QueryVulnerabilities retrieves vulnerabilities from the database based on provided filters (only severity for now)
func QueryVulnerabilities(severity string) ([]models.QueryResult, error) {
	dbConn := db.GetDB()
	rows, err := dbConn.Query("SELECT payload FROM scans WHERE severity = ?", severity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.QueryResult
	for rows.Next() {
		var payloadJSON string
		err := rows.Scan(&payloadJSON)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		var payload map[string]interface{}
		err = json.Unmarshal([]byte(payloadJSON), &payload)
		if err != nil {
			log.Printf("Error unmarshaling payload: %v", err)
			continue
		}

		scanResults, ok := payload["scanResults"].(map[string]interface{})
		if !ok {
			log.Printf("Error: scanResults not found or not a map")
			continue
		}

		vulnerabilities, ok := scanResults["vulnerabilities"].([]interface{})
		if !ok {
			log.Printf("Error: vulnerabilities not found or not an array")
			continue
		}

		for _, vuln := range vulnerabilities {
			v, ok := vuln.(map[string]interface{})
			if !ok {
				log.Printf("Error: vulnerability is not a map")
				continue
			}

			if v["severity"].(string) != severity {
				continue
			}

			vulnerability := models.QueryResult{
				ID:             v["id"].(string),
				Severity:       v["severity"].(string),
				CVSS:           v["cvss"].(float64),
				Status:         v["status"].(string),
				PackageName:    v["package_name"].(string),
				CurrentVersion: v["current_version"].(string),
				FixedVersion:   v["fixed_version"].(string),
				Description:    v["description"].(string),
				PublishedDate:  v["published_date"].(string),
				RiskFactors:    make([]string, 0),
			}

			if link, ok := v["link"].(string); ok {
				vulnerability.Link = link
			}

			for _, factor := range v["risk_factors"].([]interface{}) {
				vulnerability.RiskFactors = append(vulnerability.RiskFactors, factor.(string))
			}

			results = append(results, vulnerability)
		}
	}

	return results, nil
}
