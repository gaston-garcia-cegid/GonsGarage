package handler

// Tipos exportados solo para documentación OpenAPI (swag).

// SwaggerLoginOK respuesta de POST /auth/login.
type SwaggerLoginOK struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// SwaggerRegisterUser usuario creado (subconjunto estable para el contrato JSON).
type SwaggerRegisterUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// SwaggerRegisterOK respuesta de POST /auth/register.
type SwaggerRegisterOK struct {
	Message string              `json:"message"`
	User    SwaggerRegisterUser `json:"user"`
}

// SwaggerMeUser perfil en GET /auth/me.
type SwaggerMeUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// SwaggerMeOK envoltorio { "user": ... }.
type SwaggerMeOK struct {
	User SwaggerMeUser `json:"user"`
}

// SwaggerMessage error o mensaje genérico.
type SwaggerMessage struct {
	Error   string `json:"error,omitempty"`
	Details string `json:"details,omitempty"`
	Message string `json:"message,omitempty"`
}
