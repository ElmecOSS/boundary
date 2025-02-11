package vault

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/boundary/internal/credential/vault/store"
	"github.com/hashicorp/boundary/internal/db"
	dbassert "github.com/hashicorp/boundary/internal/db/assert"
	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/iam"
	"github.com/hashicorp/boundary/internal/kms"
	"github.com/hashicorp/boundary/internal/oplog"
	"github.com/hashicorp/boundary/internal/scheduler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestRepository_CreateCredentialLibrary(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)

	_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
	cs := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 1)[0]

	tests := []struct {
		name    string
		in      *CredentialLibrary
		opts    []Option
		want    *CredentialLibrary
		wantErr errors.Code
	}{
		{
			name:    "nil-CredentialLibrary",
			wantErr: errors.InvalidParameter,
		},
		{
			name:    "nil-embedded-CredentialLibrary",
			in:      &CredentialLibrary{},
			wantErr: errors.InvalidParameter,
		},
		{
			name: "invalid-no-store-id",
			in: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{},
			},
			wantErr: errors.InvalidParameter,
		},
		{
			name: "invalid-public-id-set",
			in: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:  cs.GetPublicId(),
					PublicId: "abcd_OOOOOOOOOO",
				},
			},
			wantErr: errors.InvalidParameter,
		},
		{
			name: "invalid-no-vault-path",
			in: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId: cs.GetPublicId(),
				},
			},
			wantErr: errors.InvalidParameter,
		},
		{
			name: "valid-no-options",
			in: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:    cs.GetPublicId(),
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:    cs.GetPublicId(),
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
		},
		{
			name: "valid-with-name",
			in: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:    cs.GetPublicId(),
					HttpMethod: "GET",
					Name:       "test-name-repo",
					VaultPath:  "/some/path",
				},
			},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:    cs.GetPublicId(),
					HttpMethod: "GET",
					Name:       "test-name-repo",
					VaultPath:  "/some/path",
				},
			},
		},
		{
			name: "valid-with-description",
			in: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:     cs.GetPublicId(),
					HttpMethod:  "GET",
					Description: "test-description-repo",
					VaultPath:   "/some/path",
				},
			},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:     cs.GetPublicId(),
					HttpMethod:  "GET",
					Description: "test-description-repo",
					VaultPath:   "/some/path",
				},
			},
		},
		{
			name: "valid-POST-method",
			in: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:     cs.GetPublicId(),
					HttpMethod:  "POST",
					Description: "test-description-repo",
					VaultPath:   "/some/path",
				},
			},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:     cs.GetPublicId(),
					HttpMethod:  "POST",
					Description: "test-description-repo",
					VaultPath:   "/some/path",
				},
			},
		},
		{
			name: "valid-POST-http-body",
			in: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:         cs.GetPublicId(),
					HttpMethod:      "POST",
					Description:     "test-description-repo",
					VaultPath:       "/some/path",
					HttpRequestBody: []byte(`{"common_name":"boundary.com"}`),
				},
			},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					StoreId:         cs.GetPublicId(),
					HttpMethod:      "POST",
					Description:     "test-description-repo",
					VaultPath:       "/some/path",
					HttpRequestBody: []byte(`{"common_name":"boundary.com"}`),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			ctx := context.Background()
			kms := kms.TestKms(t, conn, wrapper)
			sche := scheduler.TestScheduler(t, conn, wrapper)
			repo, err := NewRepository(rw, rw, kms, sche)
			require.NoError(err)
			require.NotNil(repo)
			got, err := repo.CreateCredentialLibrary(ctx, prj.GetPublicId(), tt.in, tt.opts...)
			if tt.wantErr != 0 {
				assert.Truef(errors.Match(errors.T(tt.wantErr), err), "want err: %q got: %q", tt.wantErr, err)
				assert.Nil(got)
				return
			}
			require.NoError(err)
			assert.Empty(tt.in.PublicId)
			require.NotNil(got)
			assertPublicId(t, CredentialLibraryPrefix, got.GetPublicId())
			assert.NotSame(tt.in, got)
			assert.Equal(tt.want.Name, got.Name)
			assert.Equal(tt.want.Description, got.Description)
			assert.Equal(got.CreateTime, got.UpdateTime)
			assert.NoError(db.TestVerifyOplog(t, rw, got.GetPublicId(), db.WithOperation(oplog.OpType_OP_TYPE_CREATE), db.WithCreateNotBefore(10*time.Second)))
		})
	}

	t.Run("invalid-duplicate-names", func(t *testing.T) {
		assert, require := assert.New(t), require.New(t)
		ctx := context.Background()
		kms := kms.TestKms(t, conn, wrapper)
		sche := scheduler.TestScheduler(t, conn, wrapper)
		repo, err := NewRepository(rw, rw, kms, sche)
		require.NoError(err)
		require.NotNil(repo)
		_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
		cs := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 1)[0]
		in := &CredentialLibrary{
			CredentialLibrary: &store.CredentialLibrary{
				StoreId:    cs.GetPublicId(),
				HttpMethod: "GET",
				VaultPath:  "/some/path",
				Name:       "test-name-repo",
			},
		}

		got, err := repo.CreateCredentialLibrary(ctx, prj.GetPublicId(), in)
		require.NoError(err)
		require.NotNil(got)
		assertPublicId(t, CredentialLibraryPrefix, got.GetPublicId())
		assert.NotSame(in, got)
		assert.Equal(in.Name, got.Name)
		assert.Equal(in.Description, got.Description)
		assert.Equal(got.CreateTime, got.UpdateTime)

		got2, err := repo.CreateCredentialLibrary(ctx, prj.GetPublicId(), in)
		assert.Truef(errors.Match(errors.T(errors.NotUnique), err), "want err: %q got: %q", errors.NotUnique, err)
		assert.Nil(got2)
	})

	t.Run("valid-duplicate-names-diff-stores", func(t *testing.T) {
		assert, require := assert.New(t), require.New(t)
		ctx := context.Background()
		kms := kms.TestKms(t, conn, wrapper)
		sche := scheduler.TestScheduler(t, conn, wrapper)
		repo, err := NewRepository(rw, rw, kms, sche)
		require.NoError(err)
		require.NotNil(repo)

		_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
		css := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 2)

		csA, csB := css[0], css[1]

		in := &CredentialLibrary{
			CredentialLibrary: &store.CredentialLibrary{
				HttpMethod: "GET",
				VaultPath:  "/some/path",
				Name:       "test-name-repo",
			},
		}
		in2 := in.clone()

		in.StoreId = csA.GetPublicId()
		got, err := repo.CreateCredentialLibrary(ctx, prj.GetPublicId(), in)
		require.NoError(err)
		require.NotNil(got)
		assertPublicId(t, CredentialLibraryPrefix, got.GetPublicId())
		assert.NotSame(in, got)
		assert.Equal(in.Name, got.Name)
		assert.Equal(in.Description, got.Description)
		assert.Equal(got.CreateTime, got.UpdateTime)

		in2.StoreId = csB.GetPublicId()
		got2, err := repo.CreateCredentialLibrary(ctx, prj.GetPublicId(), in2)
		require.NoError(err)
		require.NotNil(got2)
		assertPublicId(t, CredentialLibraryPrefix, got2.GetPublicId())
		assert.NotSame(in2, got2)
		assert.Equal(in2.Name, got2.Name)
		assert.Equal(in2.Description, got2.Description)
		assert.Equal(got2.CreateTime, got2.UpdateTime)
	})
}

func TestRepository_UpdateCredentialLibrary(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)

	changeHttpRequestBody := func(b []byte) func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			l.HttpRequestBody = b
			return l
		}
	}

	changeHttpMethod := func(m Method) func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			l.HttpMethod = string(m)
			return l
		}
	}

	makeHttpMethodEmptyString := func() func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			l.HttpMethod = ""
			return l
		}
	}

	changeVaultPath := func(p string) func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			l.VaultPath = p
			return l
		}
	}

	changeName := func(n string) func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			l.Name = n
			return l
		}
	}

	changeDescription := func(d string) func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			l.Description = d
			return l
		}
	}

	makeNil := func() func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			return nil
		}
	}

	makeEmbeddedNil := func() func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			return &CredentialLibrary{}
		}
	}

	deletePublicId := func() func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			l.PublicId = ""
			return l
		}
	}

	nonExistentPublicId := func() func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			l.PublicId = "abcd_OOOOOOOOOO"
			return l
		}
	}

	combine := func(fns ...func(l *CredentialLibrary) *CredentialLibrary) func(*CredentialLibrary) *CredentialLibrary {
		return func(l *CredentialLibrary) *CredentialLibrary {
			for _, fn := range fns {
				l = fn(l)
			}
			return l
		}
	}

	tests := []struct {
		name      string
		orig      *CredentialLibrary
		chgFn     func(*CredentialLibrary) *CredentialLibrary
		masks     []string
		want      *CredentialLibrary
		wantCount int
		wantErr   errors.Code
	}{
		{
			name: "nil-credential-library",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
			chgFn:   makeNil(),
			masks:   []string{nameField, descriptionField},
			wantErr: errors.InvalidParameter,
		},
		{
			name: "nil-embedded-credential-library",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
			chgFn:   makeEmbeddedNil(),
			masks:   []string{nameField, descriptionField},
			wantErr: errors.InvalidParameter,
		},
		{
			name: "no-public-id",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
			chgFn:   deletePublicId(),
			masks:   []string{nameField, descriptionField},
			wantErr: errors.InvalidPublicId,
		},
		{
			name: "updating-non-existent-credential-library",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
					Name:       "test-name-repo",
				},
			},
			chgFn:   combine(nonExistentPublicId(), changeName("test-update-name-repo")),
			masks:   []string{nameField},
			wantErr: errors.RecordNotFound,
		},
		{
			name: "empty-field-mask",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
					Name:       "test-name-repo",
				},
			},
			chgFn:   changeName("test-update-name-repo"),
			wantErr: errors.EmptyFieldMask,
		},
		{
			name: "read-only-fields-in-field-mask",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
					Name:       "test-name-repo",
				},
			},
			chgFn:   changeName("test-update-name-repo"),
			masks:   []string{"PublicId", "CreateTime", "UpdateTime", "StoreId"},
			wantErr: errors.InvalidFieldMask,
		},
		{
			name: "unknown-field-in-field-mask",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
					Name:       "test-name-repo",
				},
			},
			chgFn:   changeName("test-update-name-repo"),
			masks:   []string{"Bilbo"},
			wantErr: errors.InvalidFieldMask,
		},
		{
			name: "change-name",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
					Name:       "test-name-repo",
				},
			},
			chgFn: changeName("test-update-name-repo"),
			masks: []string{nameField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
					Name:       "test-update-name-repo",
				},
			},
			wantCount: 1,
		},
		{
			name: "change-description",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Description: "test-description-repo",
				},
			},
			chgFn: changeDescription("test-update-description-repo"),
			masks: []string{descriptionField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Description: "test-update-description-repo",
				},
			},
			wantCount: 1,
		},
		{
			name: "change-name-and-description",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Name:        "test-name-repo",
					Description: "test-description-repo",
				},
			},
			chgFn: combine(changeDescription("test-update-description-repo"), changeName("test-update-name-repo")),
			masks: []string{nameField, descriptionField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Name:        "test-update-name-repo",
					Description: "test-update-description-repo",
				},
			},
			wantCount: 1,
		},
		{
			name: "delete-name",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Name:        "test-name-repo",
					Description: "test-description-repo",
				},
			},
			masks: []string{nameField},
			chgFn: combine(changeDescription("test-update-description-repo"), changeName("")),
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Description: "test-description-repo",
				},
			},
			wantCount: 1,
		},
		{
			name: "delete-description",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Name:        "test-name-repo",
					Description: "test-description-repo",
				},
			},
			masks: []string{descriptionField},
			chgFn: combine(changeDescription(""), changeName("test-update-name-repo")),
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
					Name:       "test-name-repo",
				},
			},
			wantCount: 1,
		},
		{
			name: "do-not-delete-name",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Name:        "test-name-repo",
					Description: "test-description-repo",
				},
			},
			masks: []string{descriptionField},
			chgFn: combine(changeDescription("test-update-description-repo"), changeName("")),
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Name:        "test-name-repo",
					Description: "test-update-description-repo",
				},
			},
			wantCount: 1,
		},
		{
			name: "do-not-delete-description",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Name:        "test-name-repo",
					Description: "test-description-repo",
				},
			},
			masks: []string{nameField},
			chgFn: combine(changeDescription(""), changeName("test-update-name-repo")),
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:  "GET",
					VaultPath:   "/some/path",
					Name:        "test-update-name-repo",
					Description: "test-description-repo",
				},
			},
			wantCount: 1,
		},
		{
			name: "change-vault-path",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/old/path",
				},
			},
			chgFn: changeVaultPath("/new/path"),
			masks: []string{vaultPathField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/new/path",
				},
			},
			wantCount: 1,
		},
		{
			name: "delete-vault-path",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
			chgFn:   changeVaultPath(""),
			masks:   []string{vaultPathField},
			wantErr: errors.NotNull,
		},
		{
			name: "change-http-method",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
			chgFn: changeHttpMethod(MethodPost),
			masks: []string{httpMethodField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "POST",
					VaultPath:  "/some/path",
				},
			},
			wantCount: 1,
		},
		{
			name: "delete-http-method",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "POST",
					VaultPath:  "/some/path",
				},
			},
			chgFn: makeHttpMethodEmptyString(),
			masks: []string{httpMethodField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
			wantCount: 1,
		},
		{
			name: "add-http-request-body",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "POST",
					VaultPath:  "/some/path",
				},
			},
			chgFn: changeHttpRequestBody([]byte("new request body")),
			masks: []string{httpRequestBodyField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:      "POST",
					VaultPath:       "/some/path",
					HttpRequestBody: []byte("new request body"),
				},
			},
			wantCount: 1,
		},
		{
			name: "delete-http-request-body",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:      "POST",
					VaultPath:       "/some/path",
					HttpRequestBody: []byte("request body"),
				},
			},
			chgFn: changeHttpRequestBody(nil),
			masks: []string{httpRequestBodyField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "POST",
					VaultPath:  "/some/path",
				},
			},
			wantCount: 1,
		},
		{
			name: "change-http-request-body",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:      "POST",
					VaultPath:       "/some/path",
					HttpRequestBody: []byte("old request body"),
				},
			},
			chgFn: changeHttpRequestBody([]byte("new request body")),
			masks: []string{httpRequestBodyField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:      "POST",
					VaultPath:       "/some/path",
					HttpRequestBody: []byte("new request body"),
				},
			},
			wantCount: 1,
		},
		{
			name: "change-method-to-GET-leave-request-body",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:      "POST",
					VaultPath:       "/some/path",
					HttpRequestBody: []byte("old request body"),
				},
			},
			chgFn:   changeHttpMethod(MethodGet),
			masks:   []string{httpMethodField},
			wantErr: errors.CheckConstraint,
		},
		{
			name: "change-method-to-POST-add-request-body",
			orig: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod: "GET",
					VaultPath:  "/some/path",
				},
			},
			chgFn: combine(changeHttpRequestBody([]byte("new request body")), changeHttpMethod(MethodPost)),
			masks: []string{httpRequestBodyField, httpMethodField},
			want: &CredentialLibrary{
				CredentialLibrary: &store.CredentialLibrary{
					HttpMethod:      "POST",
					VaultPath:       "/some/path",
					HttpRequestBody: []byte("new request body"),
				},
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			ctx := context.Background()
			kms := kms.TestKms(t, conn, wrapper)
			sche := scheduler.TestScheduler(t, conn, wrapper)
			repo, err := NewRepository(rw, rw, kms, sche)
			assert.NoError(err)
			require.NotNil(repo)

			_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
			cs := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 1)[0]

			tt.orig.StoreId = cs.GetPublicId()
			orig, err := repo.CreateCredentialLibrary(ctx, prj.GetPublicId(), tt.orig)
			assert.NoError(err)
			require.NotNil(orig)

			if tt.chgFn != nil {
				orig = tt.chgFn(orig)
			}
			got, gotCount, err := repo.UpdateCredentialLibrary(ctx, prj.GetPublicId(), orig, 1, tt.masks)
			if tt.wantErr != 0 {
				assert.Truef(errors.Match(errors.T(tt.wantErr), err), "want err: %q got: %q", tt.wantErr, err)
				assert.Equal(tt.wantCount, gotCount, "row count")
				assert.Nil(got)
				return
			}
			assert.NoError(err)
			assert.Empty(tt.orig.PublicId)
			require.NotNil(got)
			assertPublicId(t, CredentialLibraryPrefix, got.GetPublicId())
			assert.Equal(tt.wantCount, gotCount, "row count")
			assert.NotSame(tt.orig, got)
			assert.Equal(tt.orig.StoreId, got.StoreId)
			underlyingDB, err := conn.SqlDB(ctx)
			require.NoError(err)
			dbassert := dbassert.New(t, underlyingDB)
			if tt.want.Name == "" {
				dbassert.IsNull(got, "name")
				return
			}
			assert.Equal(tt.want.Name, got.Name)
			if tt.want.Description == "" {
				dbassert.IsNull(got, "description")
				return
			}
			assert.Equal(tt.want.Description, got.Description)
			if tt.wantCount > 0 {
				assert.NoError(db.TestVerifyOplog(t, rw, got.GetPublicId(), db.WithOperation(oplog.OpType_OP_TYPE_UPDATE), db.WithCreateNotBefore(10*time.Second)))
			}
		})
	}

	t.Run("invalid-duplicate-names", func(t *testing.T) {
		assert, require := assert.New(t), require.New(t)
		ctx := context.Background()
		kms := kms.TestKms(t, conn, wrapper)
		sche := scheduler.TestScheduler(t, conn, wrapper)
		repo, err := NewRepository(rw, rw, kms, sche)
		assert.NoError(err)
		require.NotNil(repo)

		name := "test-dup-name"
		_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
		cs := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 1)[0]
		libs := TestCredentialLibraries(t, conn, wrapper, cs.GetPublicId(), 2)

		lA, lB := libs[0], libs[1]

		lA.Name = name
		got1, gotCount1, err := repo.UpdateCredentialLibrary(ctx, prj.GetPublicId(), lA, 1, []string{"name"})
		assert.NoError(err)
		require.NotNil(got1)
		assert.Equal(name, got1.Name)
		assert.Equal(1, gotCount1, "row count")
		assert.NoError(db.TestVerifyOplog(t, rw, lA.GetPublicId(), db.WithOperation(oplog.OpType_OP_TYPE_UPDATE), db.WithCreateNotBefore(10*time.Second)))

		lB.Name = name
		got2, gotCount2, err := repo.UpdateCredentialLibrary(ctx, prj.GetPublicId(), lB, 1, []string{"name"})
		assert.Truef(errors.Match(errors.T(errors.NotUnique), err), "want err code: %v got err: %v", errors.NotUnique, err)
		assert.Nil(got2)
		assert.Equal(db.NoRowsAffected, gotCount2, "row count")
		err = db.TestVerifyOplog(t, rw, lB.GetPublicId(), db.WithOperation(oplog.OpType_OP_TYPE_UPDATE), db.WithCreateNotBefore(10*time.Second))
		assert.Error(err)
		assert.True(errors.IsNotFoundError(err))
	})

	t.Run("valid-duplicate-names-diff-CredentialStores", func(t *testing.T) {
		assert, require := assert.New(t), require.New(t)
		ctx := context.Background()
		kms := kms.TestKms(t, conn, wrapper)
		sche := scheduler.TestScheduler(t, conn, wrapper)
		repo, err := NewRepository(rw, rw, kms, sche)
		assert.NoError(err)
		require.NotNil(repo)

		_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
		css := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 2)

		csA, csB := css[0], css[1]

		in := &CredentialLibrary{
			CredentialLibrary: &store.CredentialLibrary{
				HttpMethod: "GET",
				VaultPath:  "/some/path",
				Name:       "test-name-repo",
			},
		}
		in2 := in.clone()

		in.StoreId = csA.GetPublicId()
		got, err := repo.CreateCredentialLibrary(ctx, prj.GetPublicId(), in)
		assert.NoError(err)
		require.NotNil(got)
		assertPublicId(t, CredentialLibraryPrefix, got.GetPublicId())
		assert.NotSame(in, got)
		assert.Equal(in.Name, got.Name)
		assert.Equal(in.Description, got.Description)

		in2.StoreId = csB.GetPublicId()
		in2.Name = "first-name"
		got2, err := repo.CreateCredentialLibrary(ctx, prj.GetPublicId(), in2)
		assert.NoError(err)
		require.NotNil(got2)
		got2.Name = got.Name
		got3, gotCount3, err := repo.UpdateCredentialLibrary(ctx, prj.GetPublicId(), got2, 1, []string{"name"})
		assert.NoError(err)
		require.NotNil(got3)
		assert.NotSame(got2, got3)
		assert.Equal(got.Name, got3.Name)
		assert.Equal(got2.Description, got3.Description)
		assert.Equal(1, gotCount3, "row count")
		assert.NoError(db.TestVerifyOplog(t, rw, got2.GetPublicId(), db.WithOperation(oplog.OpType_OP_TYPE_UPDATE), db.WithCreateNotBefore(10*time.Second)))
	})

	t.Run("change-scope-id", func(t *testing.T) {
		assert, require := assert.New(t), require.New(t)
		ctx := context.Background()
		kms := kms.TestKms(t, conn, wrapper)
		sche := scheduler.TestScheduler(t, conn, wrapper)
		repo, err := NewRepository(rw, rw, kms, sche)
		assert.NoError(err)
		require.NotNil(repo)

		_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
		css := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 2)

		csA, csB := css[0], css[1]

		lA := TestCredentialLibraries(t, conn, wrapper, csA.GetPublicId(), 1)[0]
		lB := TestCredentialLibraries(t, conn, wrapper, csB.GetPublicId(), 1)[0]

		assert.NotEqual(lA.StoreId, lB.StoreId)
		orig := lA.clone()

		lA.StoreId = lB.StoreId
		assert.Equal(lA.StoreId, lB.StoreId)

		got1, gotCount1, err := repo.UpdateCredentialLibrary(ctx, prj.GetPublicId(), lA, 1, []string{"name"})

		assert.NoError(err)
		require.NotNil(got1)
		assert.Equal(orig.StoreId, got1.StoreId)
		assert.Equal(1, gotCount1, "row count")
		assert.NoError(db.TestVerifyOplog(t, rw, lA.GetPublicId(), db.WithOperation(oplog.OpType_OP_TYPE_UPDATE), db.WithCreateNotBefore(10*time.Second)))
	})
}

func TestRepository_LookupCredentialLibrary(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)

	_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
	cs := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 1)[0]
	l := TestCredentialLibraries(t, conn, wrapper, cs.GetPublicId(), 1)[0]

	badId, err := newCredentialLibraryId()
	require.NoError(t, err)
	require.NotNil(t, badId)

	tests := []struct {
		name    string
		in      string
		want    *CredentialLibrary
		wantErr errors.Code
	}{
		{
			name: "valid",
			in:   l.GetPublicId(),
			want: l,
		},
		{
			name:    "empty-public-id",
			in:      "",
			wantErr: errors.InvalidParameter,
		},
		{
			name: "not-found",
			in:   badId,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			ctx := context.Background()
			kms := kms.TestKms(t, conn, wrapper)
			sche := scheduler.TestScheduler(t, conn, wrapper)
			repo, err := NewRepository(rw, rw, kms, sche)
			assert.NoError(err)
			require.NotNil(repo)

			got, err := repo.LookupCredentialLibrary(ctx, tt.in)
			if tt.wantErr != 0 {
				assert.Truef(errors.Match(errors.T(tt.wantErr), err), "want err: %q got: %q", tt.wantErr, err)
				assert.Nil(got)
				return
			}
			require.NoError(err)

			switch {
			case tt.want == nil:
				assert.Nil(got)
			case tt.want != nil:
				assert.NotNil(got)
				assert.Equal(got, tt.want)
			}
		})
	}
}

func TestRepository_DeleteCredentialLibrary(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)

	_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
	cs := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 1)[0]
	l := TestCredentialLibraries(t, conn, wrapper, cs.GetPublicId(), 1)[0]

	badId, err := newCredentialLibraryId()
	require.NoError(t, err)
	require.NotNil(t, badId)

	tests := []struct {
		name    string
		in      string
		want    int
		wantErr errors.Code
	}{
		{
			name: "found",
			in:   l.GetPublicId(),
			want: 1,
		},
		{
			name: "not-found",
			in:   badId,
		},
		{
			name:    "empty-public-id",
			in:      "",
			wantErr: errors.InvalidParameter,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			ctx := context.Background()
			kms := kms.TestKms(t, conn, wrapper)
			sche := scheduler.TestScheduler(t, conn, wrapper)
			repo, err := NewRepository(rw, rw, kms, sche)
			assert.NoError(err)
			require.NotNil(repo)

			got, err := repo.DeleteCredentialLibrary(ctx, prj.GetPublicId(), tt.in)
			if tt.wantErr != 0 {
				assert.Truef(errors.Match(errors.T(tt.wantErr), err), "want err: %q got: %q", tt.wantErr, err)
				return
			}
			assert.NoError(err)
			assert.Equal(tt.want, got, "row count")
		})
	}
}

func TestRepository_ListCredentialLibraries(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)

	_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
	css := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 2)
	csA, csB := css[0], css[1]

	libs := TestCredentialLibraries(t, conn, wrapper, csA.GetPublicId(), 3)

	tests := []struct {
		name    string
		in      string
		opts    []Option
		want    []*CredentialLibrary
		wantErr errors.Code
	}{
		{
			name:    "with-no-credential-store-id",
			wantErr: errors.InvalidParameter,
		},
		{
			name: "CredentialStore-with-no-libraries",
			in:   csB.GetPublicId(),
			want: []*CredentialLibrary{},
		},
		{
			name: "CredentialStore-with-libraries",
			in:   csA.GetPublicId(),
			want: libs,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			ctx := context.Background()
			kms := kms.TestKms(t, conn, wrapper)
			sche := scheduler.TestScheduler(t, conn, wrapper)
			repo, err := NewRepository(rw, rw, kms, sche)
			assert.NoError(err)
			require.NotNil(repo)
			got, err := repo.ListCredentialLibraries(ctx, tt.in, tt.opts...)
			if tt.wantErr != 0 {
				assert.Truef(errors.Match(errors.T(tt.wantErr), err), "want err: %q got: %q", tt.wantErr, err)
				assert.Nil(got)
				return
			}
			require.NoError(err)
			opts := []cmp.Option{
				cmpopts.SortSlices(func(x, y *CredentialLibrary) bool { return x.PublicId < y.PublicId }),
				protocmp.Transform(),
			}
			assert.Empty(cmp.Diff(tt.want, got, opts...))
		})
	}
}

func TestRepository_ListCredentialLibraries_Limits(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)
	sche := scheduler.TestScheduler(t, conn, wrapper)

	_, prj := iam.TestScopes(t, iam.TestRepo(t, conn, wrapper))
	cs := TestCredentialStores(t, conn, wrapper, prj.GetPublicId(), 1)[0]
	const count = 10
	libs := TestCredentialLibraries(t, conn, wrapper, cs.GetPublicId(), count)

	tests := []struct {
		name     string
		repoOpts []Option
		listOpts []Option
		wantLen  int
	}{
		{
			name:    "with-no-limits",
			wantLen: count,
		},
		{
			name:     "with-repo-limit",
			repoOpts: []Option{WithLimit(3)},
			wantLen:  3,
		},
		{
			name:     "with-negative-repo-limit",
			repoOpts: []Option{WithLimit(-1)},
			wantLen:  count,
		},
		{
			name:     "with-list-limit",
			listOpts: []Option{WithLimit(3)},
			wantLen:  3,
		},
		{
			name:     "with-negative-list-limit",
			listOpts: []Option{WithLimit(-1)},
			wantLen:  count,
		},
		{
			name:     "with-repo-smaller-than-list-limit",
			repoOpts: []Option{WithLimit(2)},
			listOpts: []Option{WithLimit(6)},
			wantLen:  6,
		},
		{
			name:     "with-repo-larger-than-list-limit",
			repoOpts: []Option{WithLimit(6)},
			listOpts: []Option{WithLimit(2)},
			wantLen:  2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			ctx := context.Background()
			kms := kms.TestKms(t, conn, wrapper)
			repo, err := NewRepository(rw, rw, kms, sche, tt.repoOpts...)
			assert.NoError(err)
			require.NotNil(repo)
			got, err := repo.ListCredentialLibraries(ctx, libs[0].StoreId, tt.listOpts...)
			require.NoError(err)
			assert.Len(got, tt.wantLen)
		})
	}
}
