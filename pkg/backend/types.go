package backend

import "time"

// TODO types here were scrapped from an actual response, we need to:
// - verify this and cherry pick what we need
// - add omitempty where required
// - take note of the ResolutionRecommendation type at the bottom,
//   looks like it might either an array or an object, its data might be any

type VulnerabilityRef struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type IssueSemver struct {
	Vulnerable []string `json:"vulnerable"`
}

type IssueInsights struct {
	TriageAdvice any `json:"triageAdvice"`
}

type IssueCvssDetails struct {
	Assigner         string    `json:"assigner"`
	Severity         string    `json:"severity"`
	CvssV3Vector     string    `json:"cvssV3Vector"`
	CvssV3BaseScore  float64   `json:"cvssV3BaseScore"`
	ModificationTime time.Time `json:"modificationTime"`
}

type IssueEpssDetails struct {
	Percentile   string `json:"percentile"`
	Probability  string `json:"probability"`
	ModelVersion string `json:"modelVersion"`
}

type IssueIdentifiers struct {
	Cve  []string `json:"CVE"`
	Cwe  []string `json:"CWE"`
	Ghsa []string `json:"GHSA"`
}

type IssueMavenModuleName struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
}

type IssueRawData struct {
	ID              string               `json:"id"`
	Title           string               `json:"title"`
	CVSSv3          string               `json:"CVSSv3"`
	Credit          []string             `json:"credit"`
	Semver          IssueSemver          `json:"semver"`
	FixedIn         []any                `json:"fixedIn"`
	Insights        IssueInsights        `json:"insights"`
	Language        string               `json:"language"`
	Severity        string               `json:"severity"`
	CvssScore       float64              `json:"cvssScore"`
	IsDisputed      bool                 `json:"isDisputed"`
	ModuleName      string               `json:"moduleName"`
	CvssDetails     []IssueCvssDetails   `json:"cvssDetails"`
	Description     string               `json:"description"`
	EpssDetails     IssueEpssDetails     `json:"epssDetails"`
	Identifiers     IssueIdentifiers     `json:"identifiers"`
	PackageName     string               `json:"packageName"`
	Proprietary     bool                 `json:"proprietary"`
	DisclosureTime  time.Time            `json:"disclosureTime"`
	PackageManager  string               `json:"packageManager"`
	MavenModuleName IssueMavenModuleName `json:"mavenModuleName"`
}

type ReportedIssue struct {
	ID     string   `json:"id"`
	Source string   `json:"source"`
	Cves   []string `json:"cves"`
	//RawData IssueRawData `json:"rawData"`
}

type ResolutionMavenPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ResolutionRecommendation struct { // TODO ??
	CVE202326464 struct {
		IssueRef      string                 `json:"issueRef"`
		MavenPackage  ResolutionMavenPackage `json:"mavenPackage"`
		ProductStatus string                 `json:"productStatus"`
	} `json:"CVE-2023-26464"`
}
type DependencyAnalysisReport struct {
	Ref             VulnerabilityRef         `json:"ref"`
	Issues          []ReportedIssue          `json:"issues"`
	Transitive      []any                    `json:"transitive"`
	Recommendations ResolutionRecommendation `json:"recommendations"` // TODO ??
}
