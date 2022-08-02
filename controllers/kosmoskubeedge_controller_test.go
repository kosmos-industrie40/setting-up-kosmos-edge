package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	crdv1 "github.com/kosmos-industrie40/setting-up-kosmos-edge/api/v1"
)

var _ = Describe("snippet controller", func() {

	It("should do something", func() {
		snip := &crdv1.KosmosKubeEdge{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testsnip",
				Namespace: "default",
			},
			Spec: crdv1.KosmosKubeEdgeSpec{
				Body: crdv1.BodyProperties{
					Machine: "TestMaschine",
					Sensors: []crdv1.SensorsBody{
						{
							Name: "TestSensor",
							StorageDuration: []crdv1.StorageDurationSensors{
								{
									SystemName: "TestStorage",
								},
							},
						},
					},
					Contract: crdv1.ContractBody{
						Valid: crdv1.ValidContract{
							Start: "2021-04-15T00:00:00.000Z",
							End:   "2031-01-01T00:00:00.000Z",
						},
						CreationTime: "2021-04-15T11:32:22.1235Z",
						Partners:     []string{""},
						Permissions: crdv1.PermissionsContract{
							Read:  []string{},
							Write: []string{},
						},
						Id:             "0",
						ParentContract: "No",
						Version:        "0.0",
					},
					RequiredTechnicalContainers: []crdv1.RequiredTechnicalContainersBody{},
					KosmosLocalSystems:          []string{},
					CheckSignature:              false,
					Analysis: crdv1.AnalysisBody{
						Enable: false,
						Systems: []crdv1.SystemsAnalysis{
							{
								System:    "System",
								Enable:    false,
								Pipelines: []crdv1.PipelinesDefinitions{},
								Connection: crdv1.ConnectionSystems{
									Url:      "testUrl",
									UserMgmt: "TestUser",
									Interval: "TestInterval",
									Container: crdv1.ContainerDefinitions{
										Url:         "TestContainerURL",
										Tag:         "latest",
										Arguments:   []string{},
										Environment: []string{},
									},
								},
							},
						},
					},
					BlockchainConnection: crdv1.BlockchainConnectionBody{
						Uri: "http://blockchain-connection.de",
						ContainerList: []crdv1.ContainerDefinitions{
							{
								Url: "kosmos.idcp.io/ondics/blockchain-connector",
								Tag: "0.5.0",
								Arguments: []string{
									"--net kosmos-local",
									"--name blockchain-connector",
									"-p 1880:10005",
								},
								Environment: []string{
									"TZ=Europe/Berlin",
									"USE_TLS=false",
									"USE_STANDALONE_NO_MQTT=false",
									"BC_API_PRODDATA=http://kosmos_bc_local:3001/storage/prodData",
									"BC_API_MACHINES=http://0.0.0.0/machine",
									"BC_API_CONTRACTS=http://0.0.0.0/contract",
									"BCC_CONFIG=[{\"machineId\": \"1\",\"sensorId\": \"simulated_messages\",\"mqtt-topic\": \"kosmos/machine-data/<machineId>/sensor/<sensorId>/update\",\"blockchain-endpoint\": \"<BC_API_PRODDATA>\",\"mapping\": \"none\"}]",
									"MQTT_BROKER_FQDN=mqtt-broker.kosmos-local",
								},
							},
							{
								Url: "kosmos.idcp.io/datarella/local_app",
								Tag: "latest",
								Environment: []string{
									"CONNECTION_PROFILES_ENDPOINT=https://api.npoint.io/",
									"USE_CASE=asys",
								},
								Arguments: []string{
									"--name kosmos_bc_local",
									"--net kosmos-local",
								},
							},
						},
						Sensors: []crdv1.SensorsBlockchainConnectionBody{},
					},
					MachineConnection: []crdv1.MachineConnectionBody{},
					Metadata:          crdv1.MetadataBody{},
				},
			},
		}

		err := k8sClient.Create(context.Background(), snip)
		Expect(err).ToNot(HaveOccurred())

		entry := &crdv1.KosmosKubeEdge{}
		Eventually(func() bool {
			// check if deployment exists
			err = k8sClient.Get(context.Background(), types.NamespacedName{
				Name:      "testsnip",
				Namespace: "default",
			}, entry)
			return err == nil
		})

	})
})
