package rhpamcentr

import (
	"github.com/bmozaffa/rhpam-operator/internal/constants"
	"github.com/bmozaffa/rhpam-operator/internal/pkg/defaults"
	"github.com/bmozaffa/rhpam-operator/internal/pkg/shared"
	"github.com/bmozaffa/rhpam-operator/pkg/apis/rhpam/v1alpha1"
	"github.com/openshift/api/apps"
	"github.com/openshift/api/apps/v1"
	"github.com/openshift/api/route"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GetRHMAPCentr(cr *v1alpha1.App) []runtime.Object {
	_, serviceName, labels := shared.GetCommonLabels(cr, constants.RhpamcentrServicePrefix)
	resources := shared.GetResources(cr.Spec.Console.Resources)
	image := shared.GetImage(cr.Spec.Console.Image, "rhpam70-businesscentral-openshift")

	dc := v1.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DeploymentConfig",
			APIVersion: apps.GroupName + "/v1", //TODO find out if there is a function that provides this
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "App",
				}),
			},
			Labels: labels,
		},
		Spec: v1.DeploymentConfigSpec{
			Strategy: v1.DeploymentStrategy{
				Type: v1.DeploymentStrategyTypeRecreate,
			},
			Triggers: v1.DeploymentTriggerPolicies{
				{
					Type: v1.DeploymentTriggerOnImageChange,
					ImageChangeParams: &v1.DeploymentTriggerImageChangeParams{
						Automatic:      true,
						ContainerNames: []string{serviceName},
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Namespace: constants.ImageStreamNamespace,
							Name:      constants.RhpamcentrImageStreamName + ":" + constants.ImageStreamTag,
						},
					},
				},
			},
			Replicas: 1,
			Selector: labels,
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &[]int64{60}[0],
					Containers: []corev1.Container{
						{
							Name:            serviceName,
							Image:           image,
							ImagePullPolicy: "Always",
							Resources:       resources,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"/bin/bash", "-c", "curl --fail --silent -u 'adminUser:RedHat' http://localhost:8080/kie-wb.jsp"},
									},
								},
								InitialDelaySeconds: 180,
								TimeoutSeconds:      2,
								PeriodSeconds:       15,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"/bin/bash", "-c", "curl --fail --silent -u 'adminUser:RedHat' http://localhost:8080/kie-wb.jsp"},
									},
								},
								InitialDelaySeconds: 60,
								TimeoutSeconds:      2,
								PeriodSeconds:       30,
								FailureThreshold:    6,
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "jolokia",
									ContainerPort: 8778,
									Protocol:      "TCP",
								},
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
								{
									Name:          "git-ssh",
									ContainerPort: 8001,
									Protocol:      "TCP",
								},
							},
						},
					},
				},
			},
		},
	}
	dc.Spec.Template.Spec.Containers[0].Env = shared.GetEnvVars(defaults.ConsoleEnvironmentDefaults(), cr.Spec.Console.Env)

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "App",
				}),
			},
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       8080,
					Protocol:   "TCP",
					TargetPort: intstr.FromInt(8080),
				},
				{
					Name:       "git-ssh",
					Port:       8001,
					Protocol:   "TCP",
					TargetPort: intstr.FromInt(8001),
				},
			},
			Selector: labels,
		},
	}

	openshiftRoute := routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: route.GroupName + "/v1", //TODO find out if there is a function that provides this
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "App",
				}),
			},
			Labels: labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Name: serviceName,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("http"),
			},
		},
	}
	return []runtime.Object{dc.DeepCopyObject(), service, openshiftRoute.DeepCopyObject()}
}
