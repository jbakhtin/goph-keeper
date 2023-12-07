package appservices

// PasswordAppServiceInterface - контракт описывающий взаимодействие с учетными данными пользователей, включая проверку пароля и хэширование пароля
type PasswordAppServiceInterface interface {
	HashPassword(rawPassword string) (string, error)
	CheckPassword(password, need string) (bool, error)
	// ...
}
