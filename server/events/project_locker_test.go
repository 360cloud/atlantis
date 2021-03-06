// Copyright 2017 HootSuite Media Inc.
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Modified hereafter by contributors to runatlantis/atlantis.

package events_test

import (
	"testing"

	. "github.com/petergtz/pegomock"
	"github.com/runatlantis/atlantis/server/events"
	"github.com/runatlantis/atlantis/server/events/locking"
	"github.com/runatlantis/atlantis/server/events/locking/mocks"
	"github.com/runatlantis/atlantis/server/events/models"
	"github.com/runatlantis/atlantis/server/logging"
	. "github.com/runatlantis/atlantis/testing"
)

func TestDefaultProjectLocker_TryLockWhenLocked(t *testing.T) {
	mockLocker := mocks.NewMockLocker()
	locker := events.DefaultProjectLocker{
		Locker: mockLocker,
	}
	expProject := models.Project{}
	expWorkspace := "default"
	expPull := models.PullRequest{}
	expUser := models.User{}

	lockingPull := models.PullRequest{
		Num: 2,
	}
	When(mockLocker.TryLock(expProject, expWorkspace, expPull, expUser)).ThenReturn(
		locking.TryLockResponse{
			LockAcquired: false,
			CurrLock: models.ProjectLock{
				Pull: lockingPull,
			},
			LockKey: "",
		},
		nil,
	)
	res, err := locker.TryLock(logging.NewNoopLogger(), expPull, expUser, expWorkspace, expProject)
	Ok(t, err)
	Equals(t, &events.TryLockResponse{
		LockAcquired:      false,
		LockFailureReason: "This project is currently locked by #2. The locking plan must be applied or discarded before future plans can execute.",
	}, res)
}

func TestDefaultProjectLocker_TryLockWhenLockedSamePull(t *testing.T) {
	RegisterMockTestingT(t)
	mockLocker := mocks.NewMockLocker()
	locker := events.DefaultProjectLocker{
		Locker: mockLocker,
	}
	expProject := models.Project{}
	expWorkspace := "default"
	expPull := models.PullRequest{Num: 2}
	expUser := models.User{}

	lockingPull := models.PullRequest{
		Num: 2,
	}
	lockKey := "key"
	When(mockLocker.TryLock(expProject, expWorkspace, expPull, expUser)).ThenReturn(
		locking.TryLockResponse{
			LockAcquired: false,
			CurrLock: models.ProjectLock{
				Pull: lockingPull,
			},
			LockKey: lockKey,
		},
		nil,
	)
	res, err := locker.TryLock(logging.NewNoopLogger(), expPull, expUser, expWorkspace, expProject)
	Ok(t, err)
	Equals(t, true, res.LockAcquired)

	// UnlockFn should work.
	mockLocker.VerifyWasCalled(Never()).Unlock(lockKey)
	err = res.UnlockFn()
	Ok(t, err)
	mockLocker.VerifyWasCalledOnce().Unlock(lockKey)
}

func TestDefaultProjectLocker_TryLockUnlocked(t *testing.T) {
	RegisterMockTestingT(t)
	mockLocker := mocks.NewMockLocker()
	locker := events.DefaultProjectLocker{
		Locker: mockLocker,
	}
	expProject := models.Project{}
	expWorkspace := "default"
	expPull := models.PullRequest{Num: 2}
	expUser := models.User{}

	lockingPull := models.PullRequest{
		Num: 2,
	}
	lockKey := "key"
	When(mockLocker.TryLock(expProject, expWorkspace, expPull, expUser)).ThenReturn(
		locking.TryLockResponse{
			LockAcquired: true,
			CurrLock: models.ProjectLock{
				Pull: lockingPull,
			},
			LockKey: lockKey,
		},
		nil,
	)
	res, err := locker.TryLock(logging.NewNoopLogger(), expPull, expUser, expWorkspace, expProject)
	Ok(t, err)
	Equals(t, true, res.LockAcquired)

	// UnlockFn should work.
	mockLocker.VerifyWasCalled(Never()).Unlock(lockKey)
	err = res.UnlockFn()
	Ok(t, err)
	mockLocker.VerifyWasCalledOnce().Unlock(lockKey)
}
