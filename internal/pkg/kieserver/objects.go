package kieserver

import (
	"github.com/bmozaffa/rhpam-operator/internal/constants"
	"github.com/bmozaffa/rhpam-operator/internal/pkg/defaults"
	"github.com/bmozaffa/rhpam-operator/internal/pkg/shared"
	"github.com/bmozaffa/rhpam-operator/pkg/apis/rhpam/v1alpha1"
	"github.com/imdario/mergo"
	"github.com/openshift/api/apps/v1"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func GetKieServer(cr *v1alpha1.App) []runtime.Object {
	_, serviceName, labels := shared.GetCommonLabels(cr, constants.KieServerServicePrefix)
	image := shared.GetImage(cr.Spec.Server.Image, "rhpam70-kieserver-openshift")
	resourceReqs := map[string]map[corev1.ResourceName]string{"Limits": {corev1.ResourceMemory: "220Mi"}, "Requests": {corev1.ResourceMemory: "220Mi"}}
	livenessProbeInts := map[string]int{"InitialDelaySeconds": 180, "TimeoutSeconds": 2, "PeriodSeconds": 15,}
	livenessProbeScript := map[string]string{"username": "adminUser", "password": "RedHat", "url": "http://localhost:8080/services/rest/server/healthcheck",}
	readinessProbeInts := map[string]int{"InitialDelaySeconds": 60, "TimeoutSeconds": 2, "PeriodSeconds": 30, "FailureThreshold": 6,}
	readinessProbeScript := map[string]string{"username": "adminUser", "password": "RedHat", "url": "http://localhost:8080/services/rest/server/readycheck",}

	dc := v1.DeploymentConfig{
		TypeMeta:   shared.GetDeploymentTypeMeta(),
		ObjectMeta: shared.GetObjectMeta(serviceName, cr, labels),
		Spec: v1.DeploymentConfigSpec{
			Strategy: v1.DeploymentStrategy{
				Type: v1.DeploymentStrategyTypeRecreate,
			},
			Triggers: shared.GetDeploymentTrigger(serviceName, constants.ImageStreamNamespace, constants.KieServerImageStreamName, constants.ImageStreamTag),
			Replicas: 1,
			Selector: labels,
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: shared.GetObjectMeta(serviceName, cr, labels),
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &[]int64{60}[0],
					Containers: []corev1.Container{
						{
							Name:            serviceName,
							Image:           image,
							ImagePullPolicy: "Always",
							Resources:       shared.GetResourceRequirements(resourceReqs),
							LivenessProbe:   shared.GetProbe(livenessProbeInts, livenessProbeScript),
							ReadinessProbe:  shared.GetProbe(readinessProbeInts, readinessProbeScript),
							Ports:           shared.GetContainerPorts(map[string]int{"http": 8080, "jolokia": 8778}),
						},
					},
				},
			},
		},
	}
	//defaultEnv := defaults.ServerEnvironmentDefaults()
	//rhpamcentrServiceName := cr.ObjectMeta.Name + "-" + constants.RhpamcentrServicePrefix
	//defaultEnv["KIE_SERVER_CONTROLLER_SERVICE"] = rhpamcentrServiceName
	//defaultEnv["RHPAMCENTR_MAVEN_REPO_SERVICE"] = rhpamcentrServiceName
	//defaultEnv["EXECUTION_SERVER_ROUTE_NAME"] = serviceName
	//shared.MergeContainerConfigs(dc.Spec.Template.Spec.Containers, cr.Spec.Server, defaultEnv)

	service := &corev1.Service{
		TypeMeta:   shared.GetServiceTypeMeta(),
		ObjectMeta: shared.GetObjectMeta(serviceName, cr, labels),
		Spec:       shared.GetServiceSpec(labels, map[string]int{"http": 8080}),
	}

	openshiftRoute := routev1.Route{
		TypeMeta:   shared.GetRouteTypeMeta(),
		ObjectMeta: shared.GetObjectMeta(serviceName, cr, labels),
		Spec:       shared.GetRouteSpec(serviceName),
	}
	return []runtime.Object{dc.DeepCopyObject(), service, openshiftRoute.DeepCopyObject()}
}

func ConstructObjects(object v1alpha1.OpenShiftObject, cr *v1alpha1.App) v1alpha1.OpenShiftObject {
	defaultObject := defaults.GetServerObject()
	mergo.Merge(&defaultObject, object, mergo.WithOverride)
	shared.SetReferences(&object, cr)
	return object
}
