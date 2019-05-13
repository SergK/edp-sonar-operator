package service

import (
	"context"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sonar-operator/pkg/apis/edp/v1alpha1"
	"time"
)

type SonarService interface {
	// This is an entry point for service package. Invoked in err = r.service.Install(*instance) staticanalysistool_controller.go, Reconcile method.
	Install(instance v1alpha1.Sonar) error
}

func NewSonarService(platformService PlatformService, k8sClient client.Client) SonarService {
	return SonarServiceImpl{platformService: platformService, k8sClient: k8sClient}
}

type SonarServiceImpl struct {
	// Providing sonar service implementation through the interface (platform abstract)
	platformService PlatformService
	k8sClient       client.Client
}

// Invoking install method against SonarServiceImpl object should trigger list of methods, stored in client edp.PlatformService
func (s SonarServiceImpl) Install(instance v1alpha1.Sonar) error {
	log.Printf("Installing Sonar component has been started")
	s.updateStatus(instance, true)
	err := s.platformService.CreateSecret(instance)
	if err != nil {
		return err
	}

	err = s.platformService.CreateService(instance)
	if err != nil {
		return err
	}

	err = s.platformService.CreateDbDeployConf(instance)
	if err != nil {
		return err
	}

	sa, err := s.platformService.CreateServiceAccount(instance)
	if err != nil {
		return err
	}

	err = s.platformService.CreateSecurityContext(instance, sa)
	if err != nil {
		return err
	}

	//_ = s.platformService.CreateDeployConf(instance)
	//_ = s.platformService.CreateExternalEndpoint(instance)
	//_ = s.platformService.CreateVolume(instance)

	log.Printf("Installing Sonar component has been finished")
	return nil
}

func (s SonarServiceImpl) updateStatus(instance v1alpha1.Sonar, value bool) {
	if instance.Status.Available != value {
		instance.Status.Available = value
		instance.Status.LastTimeUpdated = time.Now()
		s.k8sClient.Update(context.TODO(), &instance)
	}
}