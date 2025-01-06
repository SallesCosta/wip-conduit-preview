package dto

// ##############  inputs ##############

type AuthenticationInput struct {
	User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

type UserDTO = struct {
	User struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
		Password string `json:"password"`
	} `json:"user"`
}

type CreateArticleInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList,omitempty"`
}

type UpdateArticleInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList,omitempty"`
}

type AddCommentInput struct {
	Comment struct {
		Body      string `json:"body"`
		ArticleID string `json:"article_id"`
	} `json:"comment"`
}

// ##############  outputs ##############

type AuthenticationOutput struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	UserName string `json:"user_name"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type Profile struct {
	UserName  string `json:"user_name"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type ProfileDTO struct {
	Profile Profile `json:"profile"`
}

type ArticleInput struct {
	Article struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		TagList     []string `json:"tag_list"`
	} `json:"article"`
}

type ArticleUpdateInput struct {
	Article struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Body        string `json:"body"`
	} `json:"article"`
}

type ArticleOutput struct {
	Slug           string       `json:"slug"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Body           string       `json:"body"`
	TagList        []string     `json:"tagList"`
	CreatedAt      string       `json:"createdAt"`
	UpdatedAt      string       `json:"updatedAt"`
	Favorited      bool         `json:"favorited"`
	FavoritesCount int          `json:"favoritesCount"`
	Author         AuthorOutput `json:"author"`
}

type AuthorOutput struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type Comment struct {
	ID        int          `json:"id"`
	CreatedAt string       `json:"createdAt"`
	UpdatedAt string       `json:"updatedAt"`
	Body      string       `json:"body"`
	Author    AuthorOutput `json:"author"`
}

type AllCommentsOutput struct {
	Comments []Comment `json:"comments"`
}

type AllTagsOutput struct {
	Tags []string `json:"tags"`
}

type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}
