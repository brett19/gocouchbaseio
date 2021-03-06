package gocouchbaseio

import (
	"encoding/json"
	"strings"
)

// RestPool represents a single pool returned from the pools REST API.
type cfgRestPool struct {
	Name         string `json:"name"`
	StreamingURI string `json:"streamingUri"`
	URI          string `json:"uri"`
}

// Pools represents the collection of pools as returned from the REST API.
type jcPools struct {
	ComponentsVersion     map[string]string `json:"componentsVersion,omitempty"`
	ImplementationVersion string            `json:"implementationVersion"`
	IsAdmin               bool              `json:"isAdminCreds"`
	UUID                  string            `json:"uuid"`
	Pools                 []cfgRestPool     `json:"pools"`
}

// A Node is a computer in a cluster running the couchbase software.
type cfgNode struct {
	ClusterCompatibility int                `json:"clusterCompatibility"`
	ClusterMembership    string             `json:"clusterMembership"`
	CouchAPIBase         string             `json:"couchApiBase"`
	Hostname             string             `json:"hostname"`
	InterestingStats     map[string]float64 `json:"interestingStats,omitempty"`
	MCDMemoryAllocated   float64            `json:"mcdMemoryAllocated"`
	MCDMemoryReserved    float64            `json:"mcdMemoryReserved"`
	MemoryFree           float64            `json:"memoryFree"`
	MemoryTotal          float64            `json:"memoryTotal"`
	OS                   string             `json:"os"`
	Ports                map[string]int     `json:"ports"`
	Status               string             `json:"status"`
	Uptime               int                `json:"uptime,string"`
	Version              string             `json:"version"`
	ThisNode             bool               `json:"thisNode,omitempty"`
}

type cfgNodeExt struct {
	Services struct {
		Kv      uint16 `json:"kv"`
		Capi    uint16 `json:"capi"`
		Mgmt    uint16 `json:"mgmt"`
		KvSsl   uint16 `json:"kvSSL"`
		CapiSsl uint16 `json:"capiSSL"`
		MgmtSsl uint16 `json:"mgmtSSL"`
	} `json:"services"`
	Hostname string `json:"hostname"`
}

// A Pool of nodes and buckets.
type cfgPool struct {
	BucketMap map[string]cfgBucket
	Nodes     []cfgNode

	BucketURL map[string]string `json:"buckets"`
}

// VBucketServerMap is the a mapping of vbuckets to nodes.
type cfgVBucketServerMap struct {
	HashAlgorithm string   `json:"hashAlgorithm"`
	NumReplicas   int      `json:"numReplicas"`
	ServerList    []string `json:"serverList"`
	VBucketMap    [][]int  `json:"vBucketMap"`
}

// Bucket is the primary entry point for most data operations.
type cfgBucket struct {
	SourceHostname      string
	AuthType            string             `json:"authType"`
	Capabilities        []string           `json:"bucketCapabilities"`
	CapabilitiesVersion string             `json:"bucketCapabilitiesVer"`
	Type                string             `json:"bucketType"`
	Name                string             `json:"name"`
	NodeLocator         string             `json:"nodeLocator"`
	Quota               map[string]float64 `json:"quota,omitempty"`
	Replicas            int                `json:"replicaNumber"`
	Password            string             `json:"saslPassword"`
	URI                 string             `json:"uri"`
	StreamingURI        string             `json:"streamingUri"`
	LocalRandomKeyURI   string             `json:"localRandomKeyUri,omitempty"`
	UUID                string             `json:"uuid"`
	DDocs               struct {
		URI string `json:"uri"`
	} `json:"ddocs,omitempty"`
	BasicStats  map[string]interface{} `json:"basicStats,omitempty"`
	Controllers map[string]interface{} `json:"controllers,omitempty"`

	// These are used for JSON IO, but isn't used for processing
	// since it needs to be swapped out safely.
	VBucketServerMap cfgVBucketServerMap `json:"vBucketServerMap"`
	Nodes            []cfgNode           `json:"nodes"`
	NodesExt         []cfgNodeExt        `json:"nodesExt,omitempty"`
}

func (cfg *cfgBucket) supports(needleCap string) bool {
	for _, cap := range cfg.Capabilities {
		if cap == needleCap {
			return true
		}
	}
	return false
}

func (cfg *cfgBucket) supportsCccp() bool {
	return cfg.supports("cccp")
}

func parseConfig(config []byte, srcHost string) (*cfgBucket, error) {
	configStr := strings.Replace(string(config), "$HOST", srcHost, -1)

	bk := new(cfgBucket)
	err := json.Unmarshal([]byte(configStr), bk)
	if err != nil {
		return nil, err
	}

	bk.SourceHostname = srcHost
	return bk, nil
}
