package repositories

import (
	"context"
	"errors"
	userrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/user"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

// UserRepositoryImpl реализует интерфейс для работы с пользователями в базе данных.
type UserRepositoryImpl struct {
	db *sqlx.DB // Подключение к базе данных
}

// NewUserRepositoryImpl создает новый экземпляр UserRepositoryImpl.
// db - подключение к базе данных.
// Возвращает реализацию интерфейса UserRepository.
func NewUserRepositoryImpl(db *sqlx.DB) userrepo.Repository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// FindByLogin ищет пользователя по его логину.
// ctx - контекст для управления таймаутами и отменой.
// login - логин пользователя, которого нужно найти.
// Возвращает сущность пользователя или ошибку, если пользователь не найден (ErrUserNotFound).
func (r *UserRepositoryImpl) FindByLogin(ctx context.Context, login string) (*entities.User, error) {
	user := &entities.User{}
	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		return nil, errors.Join(userrepo.ErrUserNotFound, err)
	}
	return user, nil
}

// Save сохраняет пользователя в базе данных.
// ctx - контекст для управления таймаутами и отменой.
// user - сущность пользователя, которую нужно сохранить.
// Возвращает ошибку, если:
//   - пользователь с таким логином уже существует (ErrUserAlreadyExists),
//   - произошла внутренняя ошибка сервера (ErrInternalServerError),
//
// В случае успешного выполнения возвращает nil.
func (r *UserRepositoryImpl) Save(ctx context.Context, user *entities.User) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO users (login, password) VALUES (:login, :password)", user)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // 23505 - это код ошибки для уникальных ограничений (duplicate key)
				return errors.Join(userrepo.ErrUserAlreadyExists, err)
			}

		}
		return userrepo.ErrInternalServerError
	}
	return nil
}
