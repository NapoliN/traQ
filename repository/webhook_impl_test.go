package repository

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traQ/model"
	"github.com/traPtitech/traQ/rbac/role"
	"github.com/traPtitech/traQ/utils"
	"github.com/traPtitech/traQ/utils/optional"
	"strings"
	"testing"
)

func TestRepositoryImpl_CreateWebhook(t *testing.T) {
	t.Parallel()
	repo, _, _, user, channel := setupWithUserAndChannel(t, common)

	t.Run("Invalid name", func(t *testing.T) {
		t.Parallel()
		_, err := repo.CreateWebhook("", "", channel.ID, user.GetID(), "")
		assert.Error(t, err)
		_, err = repo.CreateWebhook(strings.Repeat("a", 40), "", channel.ID, user.GetID(), "")
		assert.Error(t, err)
	})

	t.Run("channel not found", func(t *testing.T) {
		t.Parallel()

		_, err := repo.CreateWebhook(utils.RandAlphabetAndNumberString(20), "aaa", uuid.Must(uuid.NewV4()), user.GetID(), "test")
		assert.Error(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		assert, _ := assertAndRequire(t)

		wb, err := repo.CreateWebhook("test", "aaa", channel.ID, user.GetID(), "test")
		if assert.NoError(err) {
			assert.Equal("test", wb.GetName())
			assert.Equal("aaa", wb.GetDescription())
			assert.Equal(channel.ID, wb.GetChannelID())
			assert.Equal(user.GetID(), wb.GetCreatorID())
			assert.Equal("test", wb.GetSecret())

			u, err := repo.GetUser(wb.GetBotUserID(), false)
			if assert.NoError(err) {
				assert.True(u.IsBot())
				assert.Equal(role.Bot, u.GetRole())
				assert.Equal(model.UserAccountStatusActive, u.GetState())
			}
		}
	})
}

func TestRepositoryImpl_UpdateWebhook(t *testing.T) {
	t.Parallel()
	repo, _, _, user, channel := setupWithUserAndChannel(t, common)

	t.Run("Nil id", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, repo.UpdateWebhook(uuid.Nil, UpdateWebhookArgs{}), ErrNilID.Error())
	})

	t.Run("Not found", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, repo.UpdateWebhook(uuid.Must(uuid.NewV4()), UpdateWebhookArgs{}), ErrNotFound.Error())
	})

	t.Run("Invalid name", func(t *testing.T) {
		t.Parallel()
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")
		err := repo.UpdateWebhook(wb.GetID(), UpdateWebhookArgs{
			Name: optional.StringFrom(strings.Repeat("a", 40)),
		})
		assert.Error(t, err)
	})

	t.Run("channel not found", func(t *testing.T) {
		t.Parallel()
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")
		err := repo.UpdateWebhook(wb.GetID(), UpdateWebhookArgs{
			ChannelID: optional.UUIDFrom(uuid.Must(uuid.NewV4())),
		})
		assert.Error(t, err)
	})

	t.Run("creator not found", func(t *testing.T) {
		t.Parallel()
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")
		err := repo.UpdateWebhook(wb.GetID(), UpdateWebhookArgs{
			CreatorID: optional.UUIDFrom(uuid.Must(uuid.NewV4())),
		})
		assert.Error(t, err)
	})

	t.Run("invalid creator", func(t *testing.T) {
		t.Parallel()
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")
		err := repo.UpdateWebhook(wb.GetID(), UpdateWebhookArgs{
			CreatorID: optional.UUIDFrom(wb.GetBotUserID()),
		})
		assert.Error(t, err)
	})

	t.Run("No changes", func(t *testing.T) {
		t.Parallel()
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")
		err := repo.UpdateWebhook(wb.GetID(), UpdateWebhookArgs{})
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")
		ch := mustMakeChannel(t, repo, random)
		assert, require := assertAndRequire(t)

		err := repo.UpdateWebhook(wb.GetID(), UpdateWebhookArgs{
			Description: optional.StringFrom("new description"),
			Name:        optional.StringFrom("new name"),
			Secret:      optional.StringFrom("new secret"),
			ChannelID:   optional.UUIDFrom(ch.ID),
			CreatorID:   optional.UUIDFrom(user.GetID()),
		})
		if assert.NoError(err) {
			wb, err := repo.GetWebhook(wb.GetID())
			require.NoError(err)
			assert.Equal("new name", wb.GetName())
			assert.Equal("new description", wb.GetDescription())
			assert.Equal("new secret", wb.GetSecret())
			assert.Equal(user.GetID(), wb.GetCreatorID())
			assert.Equal(ch.ID, wb.GetChannelID())
		}
	})
}

func TestRepositoryImpl_DeleteWebhook(t *testing.T) {
	t.Parallel()
	repo, _, _, user, channel := setupWithUserAndChannel(t, common)

	t.Run("Nil id", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, repo.DeleteWebhook(uuid.Nil), ErrNilID.Error())
	})

	t.Run("Not found", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, repo.DeleteWebhook(uuid.Must(uuid.NewV4())), ErrNotFound.Error())
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		assert, _ := assertAndRequire(t)
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")

		err := repo.DeleteWebhook(wb.GetID())
		if assert.NoError(err) {
			_, err := repo.GetWebhook(wb.GetID())
			assert.EqualError(err, ErrNotFound.Error())
		}
	})
}

func TestRepositoryImpl_GetWebhook(t *testing.T) {
	t.Parallel()
	repo, _, _, user, channel := setupWithUserAndChannel(t, common)

	t.Run("Nil id", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetWebhook(uuid.Nil)
		assert.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("Not found", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetWebhook(uuid.Must(uuid.NewV4()))
		assert.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		assert, _ := assertAndRequire(t)
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")

		w, err := repo.GetWebhook(wb.GetID())
		if assert.NoError(err) {
			assert.Equal(wb.GetID(), w.GetID())
			assert.Equal(wb.GetName(), w.GetName())
			assert.Equal(wb.GetChannelID(), w.GetChannelID())
			assert.Equal(wb.GetSecret(), w.GetSecret())
			assert.Equal(wb.GetDescription(), w.GetDescription())
			assert.Equal(wb.GetCreatorID(), w.GetCreatorID())
			assert.Equal(wb.GetBotUserID(), w.GetBotUserID())
		}
	})
}

func TestRepositoryImpl_GetWebhookByBotUserId(t *testing.T) {
	t.Parallel()
	repo, _, _, user, channel := setupWithUserAndChannel(t, common)

	t.Run("Nil id", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetWebhookByBotUserID(uuid.Nil)
		assert.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("Not found", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetWebhookByBotUserID(uuid.Must(uuid.NewV4()))
		assert.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		assert, _ := assertAndRequire(t)
		wb := mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")

		w, err := repo.GetWebhookByBotUserID(wb.GetBotUserID())
		if assert.NoError(err) {
			assert.Equal(wb.GetID(), w.GetID())
			assert.Equal(wb.GetName(), w.GetName())
			assert.Equal(wb.GetChannelID(), w.GetChannelID())
			assert.Equal(wb.GetSecret(), w.GetSecret())
			assert.Equal(wb.GetDescription(), w.GetDescription())
			assert.Equal(wb.GetCreatorID(), w.GetCreatorID())
			assert.Equal(wb.GetBotUserID(), w.GetBotUserID())
		}
	})
}

func TestRepositoryImpl_GetAllWebhooks(t *testing.T) {
	t.Parallel()
	repo, assert, _, user, channel := setupWithUserAndChannel(t, ex3)

	n := 10
	for i := 0; i < n; i++ {
		mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")
	}

	arr, err := repo.GetAllWebhooks()
	if assert.NoError(err) {
		assert.Len(arr, n)
	}
}

func TestRepositoryImpl_GetWebhooksByCreator(t *testing.T) {
	t.Parallel()
	repo, _, _, user, channel := setupWithUserAndChannel(t, common)

	n := 10
	for i := 0; i < n; i++ {
		mustMakeWebhook(t, repo, random, channel.ID, user.GetID(), "test")
	}
	user2 := mustMakeUser(t, repo, random)
	mustMakeWebhook(t, repo, random, channel.ID, user2.GetID(), "test")

	t.Run("Nil id", func(t *testing.T) {
		t.Parallel()
		arr, err := repo.GetWebhooksByCreator(uuid.Nil)
		if assert.NoError(t, err) {
			assert.Empty(t, arr)
		}
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		arr, err := repo.GetWebhooksByCreator(user.GetID())
		if assert.NoError(t, err) {
			assert.Len(t, arr, n)
		}
	})
}
