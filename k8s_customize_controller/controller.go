package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	bolingcavalryv1 "k8s_customize_controller/pkg/apis/bolingcavalry/v1"
	clientset "k8s_customize_controller/pkg/client/clientset/versioned"
	studentscheme "k8s_customize_controller/pkg/client/clientset/versioned/scheme"
	informers "k8s_customize_controller/pkg/client/informers/externalversions/bolingcavalry/v1"
	listers "k8s_customize_controller/pkg/client/listers/bolingcavalry/v1"
)

const controllerAgentName = "student-controller"

const (
	SuccessSynced = "Synced"

	MessageResourceSynced = "Student synced successfully"
)

// Controller is the controller implementation for Student resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// studentclientset is a clientset for our own API group
	studentclientset clientset.Interface

	studentsLister listers.StudentLister
	studentsSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface

	recorder record.EventRecorder
}

// NewController returns a new student controller
func NewController(
	kubeclientset kubernetes.Interface,
	studentclientset clientset.Interface,
	studentInformer informers.StudentInformer) *Controller {

	utilruntime.Must(studentscheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:    kubeclientset,
		studentclientset: studentclientset,
		studentsLister:   studentInformer.Lister(),
		studentsSynced:   studentInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Students"),
		recorder:         recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when Student resources change
	studentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueStudent,
		UpdateFunc: func(old, new interface{}) {
			oldStudent := old.(*bolingcavalryv1.Student)
			newStudent := new.(*bolingcavalryv1.Student)
			if oldStudent.ResourceVersion == newStudent.ResourceVersion {
                //??????????????????????????????????????????????????????????????????
				return
			}
			controller.enqueueStudent(new)
		},
		DeleteFunc: controller.enqueueStudentForDelete,
	})

	return controller
}

//???????????????controller?????????
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("??????controller???????????????????????????????????????")
	if ok := cache.WaitForCacheSync(stopCh, c.studentsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("worker??????")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("worker????????????")
	<-stopCh
	glog.Info("worker????????????")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// ???????????????
func (c *Controller) processNextWorkItem() bool {

	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {

			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// ???syncHandler???????????????
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// ??????
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// ?????????????????????
	student, err := c.studentsLister.Students(namespace).Get(name)
	if err != nil {
		// ??????Student???????????????????????????????????????????????????????????????????????????
		if errors.IsNotFound(err) {
			glog.Infof("Student?????????????????????????????????????????????????????????: %s/%s ...", namespace, name)

			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list student by: %s/%s", namespace, name))

		return err
	}

	glog.Infof("?????????student?????????????????????: %#v ...", student)
	glog.Infof("?????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????(??????????????????)")

	c.recorder.Event(student, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// ????????????????????????????????????
func (c *Controller) enqueueStudent(obj interface{}) {
	var key string
	var err error
	// ?????????????????????
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	// ???key????????????
	c.workqueue.AddRateLimited(key)
}

// ????????????
func (c *Controller) enqueueStudentForDelete(obj interface{}) {
	var key string
	var err error
	// ??????????????????????????????
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	//??????key????????????
	c.workqueue.AddRateLimited(key)
}