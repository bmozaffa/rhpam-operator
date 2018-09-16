package handler

import (
	"github.com/bmozaffa/rhpam-operator/pkg/apis/rhpam/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestTrialEnvironmentHandling(t *testing.T) {
	handler := NewHandler()
	event := sdk.Event{
		Object: &v1alpha1.App{
			Spec: v1alpha1.AppSpec{
				Environment: "trial-ephemeral",
			},
		},
		Deleted: false}
	logrus.Debugf("Testing with environment %v", event.Object.(*v1alpha1.App).Spec.Environment)

	defer func() {
		err := recover().(error)
		logrus.Debugf("Failed with error %v", err)
		assert.Equal(t, err.Error(), "lookup kubernetes.default.svc: no such host", "Did not get expected no such host error")
	}()

	handler.Handle(nil, event)
}

func TestAuthoringEnvironmentHandling(t *testing.T) {
	handler := NewHandler()
	event := sdk.Event{
		Object: &v1alpha1.App{
			Spec: v1alpha1.AppSpec{
				Environment: "authoring",
			},
		},
		Deleted: false}
	logrus.Debugf("Testing with environment %v", event.Object.(*v1alpha1.App).Spec.Environment)

	err := handler.Handle(nil, event)
	assert.Equal(t, err, nil, "Authoring environment not yet implemented so it should be a no-op and return nil")
}

func TestUnknownEnvironmentHandling(t *testing.T) {
	handler := NewHandler()
	event := sdk.Event{
		Object: &v1alpha1.App{
			Spec: v1alpha1.AppSpec{
				Environment: "unknown",
			},
		},
		Deleted: false}
	logrus.Debugf("Testing with environment %v", event.Object.(*v1alpha1.App).Spec.Environment)

	err := handler.Handle(nil, event)
	assert.Equal(t, err, nil, "Unknown environment should result in a no-op and return nil")
}

func TestUnknownResourceTypeHandling(t *testing.T) {
	handler := NewHandler()
	event := sdk.Event{
		Object:  nil,
		Deleted: false}
	logrus.Debugf("Testing with event object %v", reflect.TypeOf(event.Object))

	err := handler.Handle(nil, event)
	assert.Equal(t, err, nil, "Unknown event type should result in a no-op and return nil")
}
