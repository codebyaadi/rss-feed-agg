package handlers

import (
	"time"

	"github.com/codebyaadi/rss-feed-agg/config"
	"github.com/codebyaadi/rss-feed-agg/internal/database"
	"github.com/google/uuid"
)

type Handler struct {
	*config.ApiConfig
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ApiKey    string    `json:"api_key"`
	Email     string    `json:"email"`
	AccessToken string	`json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// convertDatabaseUserToAPIUser converts a database user model to an API user model.
//
// Parameters:
//   - dbUser: The user model from the database, typically containing fields like ID, Name, CreatedAt, and UpdatedAt.
//
// Returns:
//   - A User struct containing the same ID, Name, CreatedAt, and UpdatedAt fields, formatted for API responses.
//
// This function maps the fields from the database-specific User struct to the API-specific
// User struct, ensuring the data is properly formatted for JSON serialization and API responses.
func convertDatabaseUserToAPIUser(dbUser database.User, accessToken, refreshToken string) User {
	user := User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		ApiKey:    dbUser.ApiKey,
		Email:     dbUser.Email,
	}

	if accessToken != "" {
		user.AccessToken = accessToken
	}
	if refreshToken != "" {
		user.RefreshToken = refreshToken
	}

	return user
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

// convertDatabaseFeedToAPIFeed converts a database feed model to an API feed model.
//
// Parameters:
//   - dbFeed: The feed model from the database, typically containing fields like ID, Name, CreatedAt, UpdatedAt, Url, and UserID.
//
// Returns:
//   - A Feed struct containing the same ID, Name, CreatedAt, UpdatedAt, Url, and UserID fields, formatted for API responses.
//
// This function maps the fields from the database-specific Feed struct to the API-specific
// Feed struct, ensuring the data is properly formatted for JSON serialization and API responses.
func convertDatabaseFeedToAPIFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		Name:      dbFeed.Name,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

// convertDatabaseFeedsToAPIFeeds converts a slice of database feed models to a slice of API feed models.
//
// Parameters:
//   - dbFeeds: A slice of database feed models, each containing fields like ID, Name, CreatedAt, UpdatedAt, Url, and UserID.
//
// Returns:
//   - A slice of Feed structs, each containing the same ID, Name, CreatedAt, UpdatedAt, Url, and UserID fields,
//     formatted for API responses.
//
// This function iterates over each database feed model, converts it to an API feed model using
// convertDatabaseFeedToAPIFeed, and returns a slice of these API feed models.
func convertDatabaseFeedsToAPIFeeds(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(dbFeeds))
	for i, dbFeed := range dbFeeds {
		feeds[i] = convertDatabaseFeedToAPIFeed(dbFeed)
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func convertDatabaseFeedFollowToAPIFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func convertDatabaseFeedFollowsToAPIFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollows := make([]FeedFollow, len(dbFeedFollows))
	for i, dbFeedFollow := range dbFeedFollows {
		feedFollows[i] = convertDatabaseFeedFollowToAPIFeedFollow(dbFeedFollow)
	}
	return feedFollows
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func convertDatabasePostToAPIPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
		PublishedAt: dbPost.PublishedAt,
		Description: description,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
	}
}

func convertDatabasePostsToAPIPosts(dbPosts []database.Post) []Post {
	posts := make([]Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		posts[i] = convertDatabasePostToAPIPost(dbPost)
	}
	return posts
}
