package tomcat

import (
	"context"

	"k8s.io/apimachinery/pkg/util/intstr"

	tomcatv1alpha1 "github.com/tomcat-operator/pkg/apis/tomcat/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_tomcat")

const (
	httpApplicationPort   int32  = 8080
	containerEnvNamespace string = "KUBERNETES_NAMESPACE"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Tomcat Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTomcat{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("tomcat-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Tomcat
	err = c.Watch(&source.Kind{Type: &tomcatv1alpha1.Tomcat{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Tomcat
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tomcatv1alpha1.Tomcat{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileTomcat implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileTomcat{}

// ReconcileTomcat reconciles a Tomcat object
type ReconcileTomcat struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Tomcat object and makes changes based on the state read
// and what is in the Tomcat.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTomcat) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Tomcat")

	// Fetch the Tomcat instance
	tomcat := &tomcatv1alpha1.Tomcat{}
	err := r.client.Get(context.TODO(), request.NamespacedName, tomcat)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Tomcat resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Tomcat")
		return reconcile.Result{}, err
	}

	// Check if the Service already exists, if not create a new one
	list := &corev1.ServiceList{}
	opts := &client.ListOptions{}
	err = r.client.List(context.TODO(), opts, list)
	if (err != nil && errors.IsNotFound(err)) || len(list.Items) == 1 {
		// Define a new Service
		ser := r.serviceForTomcat(tomcat)
		reqLogger.Info("Creating a new Service.", "Service.Namespace", ser.Namespace, "Service.Name", ser.Name)
		err = r.client.Create(context.TODO(), ser)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service.", "Service.Namespace", ser.Namespace, "Service.Name", ser.Name)
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service.")
		return reconcile.Result{}, err
	}

	// Check if the Deployment already exists, if not create a new one
	foundDeployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: tomcat.Spec.ApplicationName, Namespace: tomcat.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new DeploymentConfig
		dep := r.deploymentForTomcat(tomcat)
		reqLogger.Info("Creating a new DeploymentConfig.", "DeploymentConfig.Namespace", dep.Namespace, "DeploymentConfig.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new DeploymentConfig.", "DeploymentConfig.Namespace", dep.Namespace, "DeploymentConfig.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// DeploymentConfig created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service.")
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(tomcat)

	// Set Tomcat instance as the owner and controller
	if err := controllerutil.SetControllerReference(tomcat, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *tomcatv1alpha1.Tomcat) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func (r *ReconcileTomcat) serviceForTomcat(t *tomcatv1alpha1.Tomcat) *corev1.Service {

	applicationName := t.Spec.ApplicationName

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      applicationName,
			Namespace: t.Namespace,
			Labels: map[string]string{
				"application": applicationName,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Port:       8080,
				TargetPort: intstr.FromInt(8080),
			}},
			Selector: map[string]string{
				"application": applicationName,
			},
		},
	}

	return service
}

func (r *ReconcileTomcat) deploymentForTomcat(t *tomcatv1alpha1.Tomcat) *appsv1.Deployment {

	applicationName := t.Spec.ApplicationName
	applicationImage := t.Spec.ApplicationImage

	label := make(map[string]string)
	label["application"] = applicationName

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      applicationName,
			Namespace: t.Namespace,
			Labels:    label,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  applicationName,
						Image: applicationImage,
						Env: []corev1.EnvVar{
							{
								Name:  containerEnvNamespace,
								Value: applicationName,
							},
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: httpApplicationPort,
								Name:          "http",
							},
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: "/demo-1.0/health",
									Port: intstr.FromString("http"),
								},
							},
							InitialDelaySeconds: 3,
							PeriodSeconds:       3,
						},
					}},
				},
			},
		},
	}

	return deployment
}
