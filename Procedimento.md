## Estruturando diretório

Sugestão de template para organização de pastas: www.github.com/golang-standards/project-layout
Eu criei:

> root/cmd - onde fica o main.go. Aqui também ficará o build gerado. Aqui eu gosto de por o `.env` e `tests.http`.
> root/configs - configurações que ajudam a rodar a app, onde pega as envs.
> root/internal - parte específica desta aplicação.
> root/pkg - partes que podem ser reaproveitadas por outras aplicações usando `go mod`.

rodar o `go mod init github.com/sallescosta/`

## Criando arquivo de configuração

> root/configs/config.go
```sh
package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

var cfg *conf

type conf struct {
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	WebServePort string `mapstructure:"WEB_SERVER_PORT"`
	JwtSecret    string `mapstructure:"JWT_SECRET"`
	JwtExpiresIn int    `mapstructure:"JWT_EXPIRESIN"`
	TokenAuth    *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)
	return cfg, err
}

```

## Criar entidade User

> root/internal/entity
```sh
package entity

import (
	"github.com/sallescosta/goexpert/api/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

func NewUser(name, email, password string) (*User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}
```

## Testar entidade User
## Criar entidade Product
## Testar entidade Product
## Criar UserDB
## Testar criação do usuário
## Criar principais métodos de ProductDB
## Criar FindAll com paginação
## Testar ProductDB
## Criar Handler para criar produto
## Testando endpoint  de criação de Product
## Roteadores - busca e alteração de products
## Listando e removendo Products
## Criando user
## Gerar JWT
## Protegendo Products
## Criando e trabalhando com middlewares
## Criar documentação para User
## Geração de access token
## Criando e listando products

