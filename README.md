# Setting Up KOSMoS Edge 
This repository provides the necessary Custom Resource Definitions (CRD) for the KOSMoS Edge. For easier implementation of the CRD's the [Kubebuilder framework](https://book.kubebuilder.io/) is used. Kubebuilder is a framework for building Kubernetes APIs using custom resource definitions (CRDs).

## Prerequisites
1. Kubernetes Cluster (Tested with version 1.24)
2. Golang Version 1.16+

## Setting Up KOSMoS Edge
1. Update your ```.kube/config``` File with information to the Cluster. Otherwise will assume running in cluster and use the cluster provided kubeconfig.
2. Build the package with 
```bash
go build -o bin/manager main.go
```
3. Run the programm with ```./bin/manager```

4. (You can run the the tests with ```make tests```)

## Under the Hood


### CRD Definitions

The CRD's can be found in [kosmoskubeedge_types.go](/api/v1/kosmoskubeedge_types.go). The important CRD is ```KosmosKubeEdge```:

```golang
type KosmosKubeEdge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KosmosKubeEdgeSpec   `json:"spec,omitempty"`
	Status KosmosKubeEdgeStatus `json:"status,omitempty"`
}
```
The important property is the ```KosmosKubeEdgeSpec```:

```golang
type KosmosKubeEdgeSpec struct {
	Body BodyProperties `json:"body"`
	Signature SignatureProperties `json:"signature,omitempty"`
}
```

with the BodyProperties containing:

```golang
type BodyProperties struct {
	Machine string        `json:"machine"`
	Sensors []SensorsBody `json:"sensors"`

	//	//+kubebuilder:validation:Optional
	Contract                    ContractBody                      `json:"contract,omitempty"`
	RequiredTechnicalContainers []RequiredTechnicalContainersBody `json:"requiredTechnicalContainers,omitempty"`
	KosmosLocalSystems          []string                          `json:"kosmosLocalSystems,omitempty"`
	CheckSignature              bool                              `json:"checkSignatures,omitempty"`
	Analysis                    AnalysisBody                      `json:"analysis,omitempty"`
	BlockchainConnection        BlockchainConnectionBody          `json:"blockchainConnection,omitempty"`
	MachineConnection           []MachineConnectionBody           `json:"machineConnection,omitempty"`
	Metadata                    MetadataBody                      `json:"metadata,omitempty"`
}
```

So a KosmosKubeEdge Object has information about the contract, the required containers that need to be started, as well aso Analysis,Blockchain or even Machinenen data. For a description of all possible fields investigate [based_on_this_contract.json](based_on_this_contract.json).


### Reconcile

The Kubernetes API heavely uses Controllers for Reconcile. For the ```KosmosKubeEdge``` CRD's [kosmoskubeedge_controller.go](/controllers/kosmoskubeedge_controller.go) is the associated controller. 

For every ```KosmosKubeEdge``` object it deploys the corresponding ```machineConnection,analysisSystem,blockchainConnection``` container and the ```RequiredTechnicalContainers```
