package repository

import "RIP/internal/app/ds"

//func GetUserInfo(r *Repository, id uint) (ds.User, error) {
//	var user ds.User
//
//	result := r.db.Where("id = ?", id).Find(&user)
//
//	if result.Error != nil {
//		return ds.User{}, result.Error
//	}
//
//	return user, nil
//}
//
//func (r *Repository) AddUser(newUser *ds.User) error {
//	result := r.db.Create(&newUser)
//	if result.Error != nil {
//		// Проверяем, является ли ошибка ошибкой уникального ключа
//		if isDuplicateKeyError(result.Error) {
//			return fmt.Errorf("login already exists")
//		}
//		// В противном случае, возвращаем оригинальную ошибку
//		return result.Error
//	}
//	return nil
//}
//
//// Функция для проверки, является ли ошибка ошибкой уникального ключа
//func isDuplicateKeyError(err error) bool {
//	pgError, isPGError := err.(*pgconn.PgError)
//	if isPGError && pgError.Code == "23505" {
//		// Код "23505" является кодом ошибки уникального ключа в PostgreSQL
//		return true
//	}
//	return false
//}
//
//func (r *Repository) SignUp(ctx context.Context, newUser ds.User) error {
//	return r.db.Create(&newUser).Error
//}
//
//func (r *Repository) GetByCredentials(ctx context.Context, user ds.User) (ds.User, error) {
//	err := r.db.First(&user, "login = ? AND password = ?", user.Login, user.Password).Error
//	return user, err
//}

func (r *Repository) UsersList() (*[]ds.User, error) {
	var users []ds.User
	result := r.db.Find(&users)
	return &users, result.Error
}

func (r *Repository) Register(user *ds.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	user := &ds.User{}
	res := r.db.Where("login = ?", login).First(user)
	return user, res.Error
}
