/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//specifies that all fields in this package are required by default
//if fields are not requiered, they need to be specified as Optional
//+kubebuilder:validation:Optional
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type PortsContainer struct {
	//+kubebuilder:validation:Optional
	Dst int `json:"dst,omitempty"`

	Src   int      `json:"src"`
	Label []string `json:"label,omitempty"`
}

type ContainerDefinitions struct {
	Url         string   `json:"url"`
	Tag         string   `json:"tag"`
	Arguments   []string `json:"arguments"`
	Environment []string `json:"environment"`

	//+kubebuilder:validation:Optional
	Ports []PortsContainer `json:"ports,omitempty"`
}

type ModelDefinitions struct {
	//+kubebuilder:validation:Format:url
	Url string `json:"url"`

	Tag string `json:"tag"`
}

type ML_TriggerPipelines struct {
	//+kubebuilder:validation:Enum:[time,event]
	Type string `json:"type"`

	// TODO: update!!!!!
	Definition string `json:"definition,omitempty"`
}

type PipelinePipelines struct {
	Container     ContainerDefinitions `json:"container"`
	PersistOutput bool                 `json:"persistOutput"`
	From          *ModelDefinitions    `json:"from"`
	To            *ModelDefinitions    `json:"to"`
}

type PipelinesDefinitions struct {
	ML_Trigger ML_TriggerPipelines `json:"mltrigger"`
	Pipeline   []PipelinePipelines `json:"pipeline"`
	Sensors    []string            `json:"sensors"`
}

type ValidContract struct {
	//+kubebuilder:validation:Format:date-time
	Start string `json:"start"`
	End   string `json:"end"`
}

type PermissionsContract struct {
	Read  []string `json:"readPermission"`
	Write []string `json:"writePermission"`
}

type ContractBody struct {
	Valid ValidContract `json:"valid"`

	//+kubebuilder:validation:Format:date-time
	CreationTime string `json:"creationTime"`

	Partners    []string            `json:"partners"`
	Permissions PermissionsContract `json:"permissions"`
	Id          string              `json:"id"`

	//+kubebuilder:validation:Optional
	ParentContract string `json:"parentContract,omitempty"`
	Version        string `json:"version"`
}

type RequiredTechnicalContainersBody struct {
	System     string                 `json:"system"`
	Containers []ContainerDefinitions `json:"containers"`
}

type StorageDurationSensors struct {
	SystemName string `json:"systemName"`

	//+kubebuilder:validation:Format:duration
	Duration string `json:"duration"`
}

type MetaSensors struct {
}

type SensorsBody struct {
	Name            string                   `json:"name"`
	StorageDuration []StorageDurationSensors `json:"storageDuration"`
	Meta            MetaSensors              `json:"meta,omitempty"`
}

type ConnectionSystems struct {
	//+kubebuilder:validation:Format:url
	Url      string `json:"url"`
	UserMgmt string `json:"user-mgmt"`

	//+kubebuilder:validation:Format:duration
	Interval string `json:"interval"`

	Container ContainerDefinitions `json:"container"`
}

type SystemsAnalysis struct {
	System    string               `json:"system"`
	Enable    bool                 `json:"enable"`
	Pipelines PipelinesDefinitions `json:"pipelines"`

	//+kubebuilder:validation:Optional
	Connection ConnectionSystems `json:"connection,omitempty"`
}

type AnalysisBody struct {
	Enable  bool              `json:"enable"`
	Systems []SystemsAnalysis `json:"systems"`
}

type SensorsBlockchainConnectionBody struct {
}

type BlockchainConnectionBody struct {
	//+kubebuilder:validation:Format:url
	Uri string `json:"uri"`

	//+kubebuilder:validation:MinItems:=2
	//+kubebuilder:validation:MaxItems:=2
	ContainerList []ContainerDefinitions `json:"containerList"`

	//+kubebuilder:validation:MinItems:=0
	Sensors []SensorsBlockchainConnectionBody `json:"sensors"`
}

type CredentialsConnectionData struct {
}

type ConnectionData struct {
	//+kubebuilder:validation:Format:url
	Uri string `json:"uri"`

	Credentials CredentialsConnectionData `json:"credentials"`
	MachineID   string                    `json:"machineID"`
}

type SubscribeDataMachineConnectionsBody struct {
}

type MachineConnectionBody struct {
	Connector      ContainerDefinitions                `json:"connector"`
	ConnectionData ConnectionData                      `json:"connectionData"`
	SubscribeData  SubscribeDataMachineConnectionsBody `json:"subscribeData"`
}

type MetadataBody struct {
}

type BodyProperties struct {
	Machine string        `json:"machine"`
	Sensors []SensorsBody `json:"sensors"`

	//	//+kubebuilder:validation:Optional
	Contract                    ContractBody                      `json:"contract,omitempty"`
	RequiredTechnicalContainers []RequiredTechnicalContainersBody `json:"requiredTechnicalContainers,omitempty"`
	KosmosLocalSystems          []string                          `json:"kosmosLocalSystems,omitempty"`
	CheckSignature              bool                              `json:"checkSignature,omitempty"`
	Analysis                    AnalysisBody                      `json:"analysis,omitempty"`
	BlockchainConnection        BlockchainConnectionBody          `json:"blockchainConnection,omitempty"`
	MachineConnection           []MachineConnectionBody           `json:"machineConnection,omitempty"`
	Metadata                    MetadataBody                      `json:"metadata,omitempty"`
}

type MetaSignature struct {
	//+kubebuilder:validation:Format:date-time
	Date string `json:"date"`

	Algorithm string `json:"algorithm"`
}

type SignatureProperties struct {
	Meta      MetaSignature `json:"meta"`
	Signature string        `json:"signature"`
}

// KosmosKubeEdgeSpec defines the desired state of KosmosKubeEdge
type KosmosKubeEdgeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Body BodyProperties `json:"body"`

	//+kubebuilder:validation:Optional
	Signature SignatureProperties `json:"signature,omitempty"`
}

type KosmosKubeEdgeStatus struct {
	CreatedDeployments int `json:"createdDeployments"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KosmosKubeEdge is the Schema for the kosmoskubeedges API
type KosmosKubeEdge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KosmosKubeEdgeSpec   `json:"spec,omitempty"`
	Status KosmosKubeEdgeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KosmosKubeEdgeList contains a list of KosmosKubeEdge
type KosmosKubeEdgeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KosmosKubeEdge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KosmosKubeEdge{}, &KosmosKubeEdgeList{})
}
