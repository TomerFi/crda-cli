package analyse

import "github.com/rhecosystemappeng/crda-cli/pkg/backend/api"

type VulnerabilitiesSummary struct {
	TotalScannedDependencies           int `json:"total_scanned_dependencies"`
	TotalScannedTransitiveDependencies int `json:"total_scanned_transitives"`
	TotalVulnerabilities               int `json:"total_vulnerabilities"`
	PubliclyAvailableVulnerabilities   int `json:"publicly_available_vulnerabilities"`
	VulnerabilitiesUniqueToSynk        int `json:"vulnerabilities_unique_to_synk"`
	DirectVulnerableDependencies       int `json:"direct_vulnerable_dependencies"`
	LowVulnerabilities                 int `json:"low_vulnerabilities"`
	MediumVulnerabilities              int `json:"medium_vulnerabilities"`
	HighVulnerabilities                int `json:"high_vulnerabilities"`
	CriticalVulnerabilities            int `json:"critical_vulnerabilities"`
	//ReportLink                            string `json:"report_link"`
	TotalDirectVulnerabilitiesIgnored     int  `json:"total_direct_vulns_ignored"`
	TotalTransitiveVulnerabilitiesIgnored int  `json:"total_transitive_vulns_ignored"`
	SnykTokenStatus                       bool `json:"snyk_token_status"`
}

func processVulnerabilities(reports []api.DependencyReport) (VulnerabilitiesSummary, error) {
	summary := VulnerabilitiesSummary{}
	// TODO waiting for this https://github.com/RHEcosystemAppEng/crda-backend/issues/28
	// we need to land on a response design that will good for all,
	// the VulnerabilitiesSummary fields were copied as-is from the old cli for reference
	return summary, nil
}
