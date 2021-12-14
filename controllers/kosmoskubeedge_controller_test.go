package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	crdv1 "gitlab.inovex.de/proj-kosmos/infra/setting-up-kosmos-edge/api/v1"
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
					RequiredTechnicalContainers: []crdv1.RequiredTechnicalContainersBody{
						{
							Containers: []crdv1.ContainerDefinitions{
								{
									Url: "harbor.kosmos.idcp.inovex.io/ondics/mqtt-broker",
									Tag: "rc3",
								},
							},
						},
					},
				},
			},
		}

		err := k8sClient.Create(context.Background(), snip)
		Expect(err).ToNot(HaveOccurred())

		entry := &v1.Deployment{}
		Eventually(func() bool {
			// check if deployment exists
			err = k8sClient.Get(context.Background(), types.NamespacedName{
				Name:      "kosmoskubeedge-requtechcont-1-1",
				Namespace: "default",
			}, entry)
			return err == nil
		}, time.Second*10, time.Second).Should(BeTrue())

	})
})
