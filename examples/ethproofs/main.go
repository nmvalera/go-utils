package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	ethproofs "github.com/kkrt-labs/go-utils/ethproofs/client"
	ethproofshttp "github.com/kkrt-labs/go-utils/ethproofs/client/http"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("ETHPROOFS_API_KEY")
	if apiKey == "" {
		log.Fatal("ETHPROOFS_API_KEY environment variable is required")
	}

	// Create client
	client, err := ethproofshttp.NewClient(&ethproofshttp.Config{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// // // Create a cluster
	cluster, err := client.CreateCluster(context.Background(), &ethproofs.CreateClusterRequest{
		Nickname:      "test-keth-cluster",
		Description:   "Test Keth Cluster",
		CycleType:     "cairo",
		ZKVMVersionID: 1,
		ProofType:     "stwo",
		Configuration: []*ethproofs.ClusterMachine{
			{
				Machine: &ethproofs.Machine{
					CPUModel:               "Intel Xeon 8375C (Ice Lake)",
					CPUCores:               32,
					GPUModels:              []string{},
					GPUCount:               []int64{},
					GPUMemoryGB:            []int64{},
					MemorySizeGB:           []int64{},
					MemoryCount:            []int64{},
					MemoryType:             []string{},
					StorageSizeGb:          16,
					TotalTeraFlops:         0,
					NetworkBetweenMachines: "100Gbps",
				},
				MachineCount: 1,
				CloudInstance: &ethproofs.CloudInstance{
					ID:                1,
					ProviderID:        1,
					InstanceName:      "r6i.8xlarge",
					Region:            "us-east-1",
					HourlyPrice:       0.5,
					CPUArch:           "x86_64",
					CPUCores:          32,
					CPUEffectiveCores: 32,
					CPUName:           "Intel Xeon 8375C (Ice Lake)",
					Memory:            128,
					GPUCount:          0,
				},
				CloudInstanceCount: 1,
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}
	fmt.Printf("Created cluster with ID: %d\n", cluster.ID)

	// List all clusters
	clusters, err := client.ListClusters(context.Background())
	if err != nil {
		log.Fatalf("Failed to list clusters: %v", err)
	}

	b, err := json.Marshal(clusters)
	if err != nil {
		log.Fatalf("Failed to marshal clusters: %v", err)
	}
	fmt.Println(string(b))

	// List all cloud instances
	cloudInstances, err := client.ListCloudInstances(context.Background())
	if err != nil {
		log.Fatalf("Failed to list cloud instances: %v", err)
	}

	b, err = json.Marshal(cloudInstances)
	if err != nil {
		log.Fatalf("Failed to marshal cloud instances: %v", err)
	}
	fmt.Println(string(b))

	// // Create a single machine
	// machine, err := client.CreateMachine(context.Background(), &ethproofs.CreateMachineRequest{
	// 	Nickname:     "test-machine",
	// 	Description:  "Test single machine for integration testing",
	// 	Hardware:     "RISC-V Prover",
	// 	CycleType:    "SP1",
	// 	ProofType:    "Groth16",
	// 	InstanceType: "t3.small",
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to create machine: %v", err)
	// }
	// fmt.Printf("\nCreated machine with ID: %d\n", machine.ID)

	// // List AWS pricing
	// instances, err := client.ListAWSPricing(context.Background())
	// if err != nil {
	// 	log.Fatalf("Failed to list AWS pricing: %v", err)
	// }
	// fmt.Println("\nAvailable AWS instances:")
	// for _, instance := range instances {
	// 	fmt.Printf("- %s: $%.3f/hour (%d vCPUs, %.1fGB RAM)\n",
	// 		instance.InstanceType,
	// 		instance.HourlyPrice,
	// 		instance.VCPU,
	// 		instance.InstanceMemory)
	// }

	// // Demonstrate proof lifecycle
	// // 1. Queue a proof
	// queuedProof, err := client.QueueProof(context.Background(), &ethproofs.QueueProofRequest{
	// 	BlockNumber: 12345,
	// 	ClusterID:   cluster.ID,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to queue proof: %v", err)
	// }
	// fmt.Printf("\nQueued proof with ID: %d\n", queuedProof.ProofID)

	// // 2. Start proving
	// startedProof, err := client.StartProving(context.Background(), &ethproofs.StartProvingRequest{
	// 	BlockNumber: 12345,
	// 	ClusterID:   cluster.ID,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to start proving: %v", err)
	// }
	// fmt.Printf("Started proving proof with ID: %d\n", startedProof.ProofID)

	// // 3. Submit completed proof
	// provingCycles := int64(1000000)
	// submittedProof, err := client.SubmitProof(context.Background(), &ethproofs.SubmitProofRequest{
	// 	BlockNumber:   12345,
	// 	ClusterID:     cluster.ID,
	// 	ProvingTime:   60000, // 60 seconds in milliseconds
	// 	ProvingCycles: &provingCycles,
	// 	Proof:         "base64_encoded_proof_data_here",
	// 	VerifierID:    "test-verifier",
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to submit proof: %v", err)
	// }
	// fmt.Printf("Submitted completed proof with ID: %d\n", submittedProof.ProofID)
}
