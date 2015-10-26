package handler

import (
	"testing"

	log "github.com/dmuth/google-go-log4go"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestScheduler(t *testing.T) {
	Convey("eremeticScheduler", t, func() {
		s := eremeticScheduler{}
		id := "eremetic-task.9999"
		runningTasks[id] = eremeticTask{
			ID: id,
		}

		Convey("newTask", func() {
			task := eremeticTask{
				ID: "eremetic-task.1234",
			}
			offer := mesos.Offer{}
			newTask := s.newTask(&offer, &task)

			So(newTask.GetTaskId().GetValue(), ShouldEqual, task.ID)
		})

		Convey("createEremeticScheduler", func() {
			s := createEremeticScheduler()
			So(s.tasksCreated, ShouldEqual, 0)
		})

		Convey("API", func() {
			Convey("Registered", func() {
				fID := mesos.FrameworkID{Value: proto.String("1234")}
				mInfo := mesos.MasterInfo{}
				s.Registered(nil, &fID, &mInfo)
			})

			Convey("Reregistered", func() {
				s.Reregistered(nil, &mesos.MasterInfo{})
			})

			Convey("Disconnected", func() {
				s.Disconnected(nil)
			})

			Convey("ResourceOffers", func() {
				driver := NewMockScheduler()
				var offers []*mesos.Offer

				Convey("No offers", func() {
					s.ResourceOffers(driver, offers)
					So(driver.AssertNotCalled(t, "DeclineOffer"), ShouldBeTrue)
					So(driver.AssertNotCalled(t, "LaunchTasks"), ShouldBeTrue)
				})

				Convey("No tasks", func() {
					offers = append(offers, &mesos.Offer{Id: &mesos.OfferID{Value: proto.String("1234")}})
					driver.On("DeclineOffer").Return("declined").Once()
					s.ResourceOffers(driver, offers)
					So(driver.AssertCalled(t, "DeclineOffer"), ShouldBeTrue)
					So(driver.AssertNotCalled(t, "LaunchTasks"), ShouldBeTrue)
				})
			})

			Convey("StatusUpdate", func() {
				s.StatusUpdate(nil, &mesos.TaskStatus{
					TaskId: &mesos.TaskID{
						Value: proto.String(id),
					},
					State: mesos.TaskState_TASK_FAILED.Enum(),
				})

				So(runningTasks[id].Status, ShouldEqual, "TASK_FAILED")
			})

			Convey("FrameworkMessage", func() {
				driver := NewMockScheduler()
				message := `{"message": "this is a message"}`
				Convey("From Eremetic", func() {
					source := "eremetic-executor"
					executor := mesos.ExecutorID{
						Value: proto.String(source),
					}
					s.FrameworkMessage(driver, &executor, &mesos.SlaveID{}, message)
				})

				Convey("From an unknown source", func() {
					source := "other-source"
					executor := mesos.ExecutorID{
						Value: proto.String(source),
					}
					s.FrameworkMessage(driver, &executor, &mesos.SlaveID{}, message)
				})

				Convey("A bad json", func() {
					source := "eremetic-executor"
					executor := mesos.ExecutorID{
						Value: proto.String(source),
					}
					s.FrameworkMessage(driver, &executor, &mesos.SlaveID{}, "not a json")
				})
			})

			Convey("OfferRescinded", func() {
				s.OfferRescinded(nil, &mesos.OfferID{})
			})

			Convey("SlaveLost", func() {
				s.SlaveLost(nil, &mesos.SlaveID{})
			})

			Convey("ExecutorLost", func() {
				s.ExecutorLost(nil, &mesos.ExecutorID{}, &mesos.SlaveID{}, 2)
			})

			Convey("Error", func() {
				s.Error(nil, "Error")
			})
		})
	})
}

//------------------ Mock Scheduler ------------------------------------------//

type MockScheduler struct {
	mock.Mock
}

func NewMockScheduler() *MockScheduler {
	return &MockScheduler{}
}

func (sched *MockScheduler) Abort() (stat mesos.Status, err error) {
	sched.Called()
	return mesos.Status_DRIVER_ABORTED, nil
}

func (sched *MockScheduler) DeclineOffer(*mesos.OfferID, *mesos.Filters) (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_STOPPED, nil
}

func (sched *MockScheduler) Join() (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) KillTask(*mesos.TaskID) (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) ReconcileTasks([]*mesos.TaskStatus) (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) RequestResources([]*mesos.Request) (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) ReviveOffers() (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) Run() (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) Start() (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) Stop(bool) (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) SendFrameworkMessage(*mesos.ExecutorID, *mesos.SlaveID, string) (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) LaunchTasks([]*mesos.OfferID, []*mesos.TaskInfo, *mesos.Filters) (mesos.Status, error) {
	sched.Called()
	return mesos.Status_DRIVER_RUNNING, nil
}

func (sched *MockScheduler) Registered(sched.SchedulerDriver, *mesos.FrameworkID, *mesos.MasterInfo) {
	sched.Called()
}

func (sched *MockScheduler) Reregistered(sched.SchedulerDriver, *mesos.MasterInfo) {
	sched.Called()
}

func (sched *MockScheduler) Disconnected(sched.SchedulerDriver) {
	sched.Called()
}

func (sched *MockScheduler) ResourceOffers(sched.SchedulerDriver, []*mesos.Offer) {
	sched.Called()
}

func (sched *MockScheduler) OfferRescinded(sched.SchedulerDriver, *mesos.OfferID) {
	sched.Called()
}

func (sched *MockScheduler) StatusUpdate(sched.SchedulerDriver, *mesos.TaskStatus) {
	sched.Called()
}

func (sched *MockScheduler) FrameworkMessage(sched.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, string) {
	sched.Called()
}

func (sched *MockScheduler) SlaveLost(sched.SchedulerDriver, *mesos.SlaveID) {
	sched.Called()
}

func (sched *MockScheduler) ExecutorLost(sched.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, int) {
	sched.Called()
}

func (sched *MockScheduler) Error(d sched.SchedulerDriver, msg string) {
	log.Error(msg)
	sched.Called()
}