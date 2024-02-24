package postgres

import (
	"context"
	"fmt"
	"github.com/go-playground/assert/v2"
	"test/api/models"
	"test/config"
	"test/pkg/helper"
	"test/pkg/logger"
	"testing"
)

func TestUserRepo_Create(t *testing.T) {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Errorf("error while connection to db error: %v", err)
	}

	createUser := models.CreateUser{
		FullName: helper.GenerateFullName(),
		Phone:    helper.GeneratePhoneNumber(),
		Password: "password",
		Cash:     10,
		UserType: "customer",
		BranchID: "aa541fcc-bf74-11ee-ae0b-166244b65504",
	}
	userID, err := pgStore.User().Create(context.Background(), createUser)
	if err != nil {
		t.Errorf("error while creating user error: %v", err)
	}

	fmt.Println("phone", createUser.Phone)
	user, err := pgStore.User().GetByID(context.Background(), models.PrimaryKey{
		ID: userID,
	})
	if err != nil {
		t.Errorf("error while getting user error: %v", err)
	}

	assert.Equal(t, user.FullName, createUser.FullName)
	assert.Equal(t, user.Phone, createUser.Phone)
	assert.Equal(t, user.Cash, createUser.Cash)
}

func TestUserRepo_GetByID(t *testing.T) {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Errorf("error while connection to db error: %v", err)
	}

	createUser := models.CreateUser{
		FullName: helper.GenerateFullName(),
		Phone:    helper.GeneratePhoneNumber(),
		Password: "password",
		Cash:     10,
		UserType: "customer",
		BranchID: "aa541fcc-bf74-11ee-ae0b-166244b65504",
	}
	userID, err := pgStore.User().Create(context.Background(), createUser)
	if err != nil {
		t.Errorf("error while creating user error: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		user, err := pgStore.User().GetByID(context.Background(), models.PrimaryKey{
			ID: userID,
		})
		if err != nil {
			t.Errorf("error while getting user by id error: %v", err)
		}

		if user.ID != userID {
			t.Errorf("expected: %q, but got %q", userID, user.ID)
		}

		if user.FullName == "" {
			t.Error("expected some full name, but got nothing")
		}

		if user.Phone == "" {
			t.Error("expected some full name, but got nothing")
		} else if len(user.Phone) >= 14 || len(user.Phone) <= 12 {
			t.Errorf("expected phone length: 13, but got %d, user id is %s", len(user.Phone), user.ID)
		}

		if user.Cash < 0 {
			t.Errorf("expected > 0, but got %d", user.Cash)
		}

		if user.BranchID == "" {
			t.Error("expected some branch id, but got nothing")
		}
	})

	t.Run("failure", func(t *testing.T) {
		userID = ""
		user, err := pgStore.User().GetByID(context.Background(), models.PrimaryKey{
			ID: userID,
		})
		if err != nil {
			t.Errorf("error while getting user by id error: %v", err)
		}

		if user.ID != userID {
			t.Errorf("expected: %q, but got %q", userID, user.ID)
		}

		if user.FullName == "" {
			t.Error("expected some full name, but got nothing")
		}

		if user.Phone == "" {
			t.Error("expected some full name, but got nothing")
		} else if len(user.Phone) >= 14 || len(user.Phone) <= 12 {
			t.Errorf("expected phone length: 13, but got %d, user id is %s", len(user.Phone), user.ID)
		}

		if user.Cash < 0 {
			t.Errorf("expected > 0, but got %d", user.Cash)
		}

		if user.BranchID == "" {
			t.Error("expected some branch id, but got nothing")
		}
	})
}

func TestUserRepo_GetList(t *testing.T) {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Errorf("error while connection to db error: %v", err)
	}

	usersResp, err := pgStore.User().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 1000,
	})
	if err != nil {
		t.Errorf("error while getting usersResp error: %v", err)
	}

	assert.Equal(t, len(usersResp.Users), 46)
}

func TestUserRepo_Update(t *testing.T) {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Errorf("error while connection to db error: %v", err)
	}

	createUser := models.CreateUser{
		FullName: helper.GenerateFullName(),
		Phone:    helper.GeneratePhoneNumber(),
		Password: "password",
		Cash:     10,
		UserType: "customer",
		BranchID: "aa541fcc-bf74-11ee-ae0b-166244b65504",
	}

	userID, err := pgStore.User().Create(context.Background(), createUser)
	if err != nil {
		t.Errorf("error while creating user error: %v", err)
	}

	UpdateUser := models.UpdateUser{
		ID:       userID,
		FullName: helper.GenerateFullName(),
		Phone:    helper.GeneratePhoneNumber(),
		Cash:     10,
	}

	userUpdateID, err := pgStore.User().Update(context.Background(), UpdateUser)
	if err != nil {
		t.Errorf("error while creating user error: %v", err)
	}

	user, err := pgStore.User().GetByID(context.Background(), models.PrimaryKey{
		ID: userUpdateID,
	})
	if err != nil {
		t.Errorf("error while getting user error: %v", err)
	}

	assert.Equal(t, userID, user.ID)
	assert.Equal(t, user.FullName, UpdateUser.FullName)
	assert.Equal(t, user.Phone, UpdateUser.Phone)
}

func TestUserRepo_Delete(t *testing.T) {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Errorf("error while connection to db error: %v", err)
	}

	createUser := models.CreateUser{
		FullName: helper.GenerateFullName(),
		Phone:    helper.GeneratePhoneNumber(),
		Password: "password",
		Cash:     10,
		UserType: "customer",
		BranchID: "aa541fcc-bf74-11ee-ae0b-166244b65504",
	}

	userID, err := pgStore.User().Create(context.Background(), createUser)
	if err != nil {
		t.Errorf("error while creating user error: %v", err)
	}

	if err = pgStore.User().Delete(context.Background(), models.PrimaryKey{ID: userID}); err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}