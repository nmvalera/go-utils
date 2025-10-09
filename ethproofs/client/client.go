package ethproofs

import (
	"context"
)

// Package ethproofs provides a Go client for the EthProofs API.
//
// For more information about EthProofs, visit:
//   - API Documentation: https://staging--ethproofs.netlify.app/api.html
//   - App Preview: https://staging--ethproofs.netlify.app/
//   - Repository: https://github.com/ethproofs/ethproofs
//
//go:generate mockgen -source client.go -destination mock/client.go -package mock Client

// Client defines the interface for interacting with the EthProofs API
type Client interface {
	// Clusters
	CreateCluster(ctx context.Context, req *CreateClusterRequest) (*CreateClusterResponse, error)
	ListClusters(ctx context.Context) ([]Cluster, error)

	// Single Machine
	CreateMachine(ctx context.Context, req *CreateMachineRequest) (*CreateMachineResponse, error)

	// Proofs
	QueueProof(ctx context.Context, req *QueueProofRequest) (*ProofResponse, error)
	StartProving(ctx context.Context, req *StartProvingRequest) (*ProofResponse, error)
	SubmitProof(ctx context.Context, req *SubmitProofRequest) (*ProofResponse, error)

	// AWS Pricing
	ListAWSPricing(ctx context.Context) ([]AWSInstance, error)
}

// Request/Response types for Clusters
// CreateClusterRequest is the request to create a cluster
type CreateClusterRequest struct {
	Nickname      string          `json:"nickname"`
	Description   string          `json:"description,omitempty"`
	Hardware      string          `json:"hardware,omitempty"`
	CycleType     string          `json:"cycle_type,omitempty"`
	ProofType     string          `json:"proof_type,omitempty"`
	Configuration []ClusterConfig `json:"configuration"`
}

// ClusterConfig is the configuration for a cluster
type ClusterConfig struct {
	InstanceType  string `json:"instance_type"`
	InstanceCount int64  `json:"instance_count"`
}

// CreateClusterResponse is the response to create a cluster
type CreateClusterResponse struct {
	ID int64 `json:"id"`
}

// ListClustersResponse is the response to list clusters
type ListClustersResponse []Cluster

// Cluster is a cluster
type Cluster struct {
	ID                   int64           `json:"id"`
	Nickname             string          `json:"nickname"`
	Description          string          `json:"description"`
	Hardware             string          `json:"hardware"`
	CycleType            string          `json:"cycle_type"`
	ProofType            string          `json:"proof_type"`
	ClusterConfiguration []ClusterConfig `json:"cluster_configuration"`
}

// Request/Response types for Single Machine
type CreateMachineRequest struct {
	Nickname     string `json:"nickname"`
	Description  string `json:"description,omitempty"`
	Hardware     string `json:"hardware,omitempty"`
	CycleType    string `json:"cycle_type,omitempty"`
	ProofType    string `json:"proof_type,omitempty"`
	InstanceType string `json:"instance_type"`
}

// CreateMachineResponse is the response to create a machine

// CreateMachineResponse is the response to create a machine
type CreateMachineResponse struct {
	ID int64 `json:"id"`
}

// Request/Response types for Proofs
// QueueProofRequest is the request to queue a proof
type QueueProofRequest struct {
	BlockNumber int64 `json:"block_number"`
	ClusterID   int64 `json:"cluster_id"`
}

// StartProvingRequest is the request to start proving
type StartProvingRequest struct {
	BlockNumber int64 `json:"block_number"`
	ClusterID   int64 `json:"cluster_id"`
}

// SubmitProofRequest is the request to submit a proof
type SubmitProofRequest struct {
	BlockNumber   int64  `json:"block_number"`
	ClusterID     int64  `json:"cluster_id"`
	ProvingTime   int64  `json:"proving_time"`
	ProvingCycles *int64 `json:"proving_cycles,omitempty"`
	Proof         string `json:"proof"`
	VerifierID    string `json:"verifier_id,omitempty"`
}

// ProofResponse is the response to a proof
type ProofResponse struct {
	ProofID int64 `json:"proof_id"`
}

// Request/Response types for AWS Pricing
// ListAWSPricingResponse is the response to list AWS pricing
type ListAWSPricingResponse = []AWSInstance

// AWSInstance is an AWS instance
type AWSInstance struct {
	ID              int64   `json:"id"`
	InstanceType    string  `json:"instance_type"`
	Region          string  `json:"region"`
	HourlyPrice     float64 `json:"hourly_price"`
	InstanceMemory  float64 `json:"instance_memory"`
	VCPU            int64   `json:"vcpu"`
	InstanceStorage string  `json:"instance_storage"`
	CreatedAt       string  `json:"created_at"`
}
