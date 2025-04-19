package rbac

import (
	"context"
	"nostalgia/internal/app"
	"nostalgia/internal/app/query"
)

type Condition func(roles ...string) bool

type Context struct {
	app app.Application
}

func New(app app.Application) Context {
	return Context{
		app: app,
	}
}

func (c *Context) GetUserRoles(ctx context.Context, userId string) ([]string, error) {
	return c.app.Queries.GetRolesForUser.Handle(ctx, query.GetRolesForUser{
		UserId: userId,
	})
}

func All(expectedRoles ...string) func(actualRoles ...string) bool {
	return func(actualRoles ...string) bool {
		for _, expected := range expectedRoles {
			found := false
			for _, actual := range actualRoles {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func Any(expectedRoles ...string) func(actualRoles ...string) bool {
	return func(actualRoles ...string) bool {
		for _, expected := range expectedRoles {
			for _, actual := range actualRoles {
				if actual == expected {
					return true
				}
			}
		}
		return false
	}
}
