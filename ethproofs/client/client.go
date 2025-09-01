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
	ClustersClient
	ProofsClient
	SingleMachineClient
	CloudInstanceClient
}

// ClustersClient defines the interface for interacting with the EthProofs Clusters API
type ClustersClient interface {
	// Clusters
	CreateCluster(ctx context.Context, req *CreateClusterRequest) (*CreateClusterResponse, error)
	ListClusters(ctx context.Context) ([]Cluster, error)
}

// CreateClusterRequest is the request for creating a cluster
type CreateClusterRequest struct {
	Nickname      string            `json:"nickname"`
	Description   string            `json:"description,omitempty"`
	ZKVMVersionID uint              `json:"zkvm_version_id"`
	Hardware      string            `json:"hardware,omitempty"`
	CycleType     string            `json:"cycle_type"`
	ProofType     string            `json:"proof_type"`
	Configuration []*ClusterMachine `json:"configuration"`
}

// CreateClusterResponse is the response for creating a cluster
type CreateClusterResponse struct {
	ID int64 `json:"id"`
}

// ListClustersResponse is the response for listing clusters
type ListClustersResponse []Cluster

// Cluster is a cluster
type Cluster struct {
	ID          int64             `json:"id"`
	Nickname    string            `json:"nickname"`
	Description string            `json:"description"`
	Hardware    string            `json:"hardware,omitempty"`
	CycleType   string            `json:"cycle_type,omitempty"`
	ProofType   string            `json:"proof_type,omitempty"`
	Machines    []*ClusterMachine `json:"machines"`
}

// ClusterMachine is the configuration for a cluster
type ClusterMachine struct {
	Machine            *Machine       `json:"machine"`
	MachineCount       int64          `json:"machine_count"`
	CloudInstance      *CloudInstance `json:"cloud_instance,omitempty"`
	CloudInstanceCount int64          `json:"cloud_instance_count,omitempty"`
}

// ProofsClient defines the interface for interacting with the EthProofs Proofs API
type ProofsClient interface {
	// QueueProof indicates the prover has started proving a block.
	QueueProof(ctx context.Context, req *QueueProofRequest) (*Proof, error)

	// StartProving indicates the prover has started to prove a block.
	StartProving(ctx context.Context, req *StartProvingRequest) (*Proof, error)

	// SubmitProof indicates the prover has completed proving a block and submits the proof.
	SubmitProof(ctx context.Context, req *SubmitProofRequest) (*Proof, error)
}

// QueueProofRequest is the request for queuing a proof
type QueueProofRequest struct {
	BlockNumber int64 `json:"block_number"`
	ClusterID   int64 `json:"cluster_id"`
}

// StartProvingRequest is the request for starting a proof
type StartProvingRequest struct {
	BlockNumber int64 `json:"block_number"`
	ClusterID   int64 `json:"cluster_id"`
}

// SubmitProofRequest is the request for submitting a proof
type SubmitProofRequest struct {
	BlockNumber   int64  `json:"block_number"`
	ClusterID     int64  `json:"cluster_id"`
	ProvingTime   int64  `json:"proving_time"`
	ProvingCycles *int64 `json:"proving_cycles,omitempty"`
	Proof         string `json:"proof"`
	VerifierID    string `json:"verifier_id,omitempty"`
}

// Proof is the response for queuing a proof
type Proof struct {
	ProofID int64 `json:"proof_id"`
}

// SingleMachineClient defines the interface for interacting with the EthProofs API
type SingleMachineClient interface {

	// Single Machine
	CreateMachine(ctx context.Context, req *Machine) (*CreateMachineResponse, error)
}

type CreateSingleMachineRequest struct {
	Nickname          string   `json:"nickname"`
	Description       string   `json:"description,omitempty"`
	ZKVMVersionID     int64    `json:"zkvm_version_id"`
	Hardware          string   `json:"hardware,omitempty"`
	CycleType         string   `json:"cycle_type"`
	ProofType         string   `json:"proof_type"`
	Machine           *Machine `json:"machine"`
	CloudInstanceName string   `json:"cloud_instance_name"`
}

type Machine struct {
	CPUModel               string   `json:"cpu_model"`
	CPUCores               int64    `json:"cpu_cores"`
	GPUModels              []string `json:"gpu_models"`
	GPUCount               []int64  `json:"gpu_count"`
	GPUMemoryGB            []int64  `json:"gpu_memory_gb"`
	MemorySizeGB           []int64  `json:"memory_size_gb"`
	MemoryCount            []int64  `json:"memory_count"`
	MemoryType             []string `json:"memory_type"`
	StorageSizeGb          int64    `json:"storage_size_gb,omitempty"`
	TotalTeraFlops         int64    `json:"total_tera_flops,omitempty"`
	NetworkBetweenMachines string   `json:"network_between_machines,omitempty"`
}

type CreateMachineResponse struct {
	ID int64 `json:"id"`
}

type CloudInstanceClient interface {
	ListCloudInstances(ctx context.Context) ([]CloudInstance, error)
}

type CloudInstance struct {
	ID                int64   `json:"id"`
	ProviderID        int64   `json:"provider_id"`
	InstanceName      string  `json:"instance_name"`
	Region            string  `json:"region"`
	HourlyPrice       float64 `json:"hourly_price"`
	CPUArch           string  `json:"cpu_arch"`
	CPUCores          int64   `json:"cpu_cores"`
	CPUEffectiveCores int64   `json:"cpu_effective_cores"`
	CPUName           string  `json:"cpu_name"`
	Memory            float64 `json:"memory"`
	GPUCount          int64   `json:"gpu_count"`
	GPUArch           string  `json:"gpu_arch"`
	GPUName           string  `json:"gpu_name"`
	GPUMemory         string  `json:"gpu_memory"`
	MoboName          string  `json:"mobo_name"`
	DiskName          string  `json:"disk_name"`
	DiskSpace         string  `json:"disk_space"`
	CreatedAt         string  `json:"created_at"`
	SnapshotDate      string  `json:"snapshot_date"`
	Provider          string  `json:"provider"`
}

type CloudProvider struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at"`
}
