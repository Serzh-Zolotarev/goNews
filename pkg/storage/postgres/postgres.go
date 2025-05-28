package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func New(dbURL string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT p.id, p.title, p.content, a.id, a.name, p.created_at 
	FROM authors a, posts p
	WHERE p.author_id = a.id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []storage.Post

	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, rows.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), ` 
	INSERT INTO posts (title, content, author_id, created_at)
	VALUES ($1, $2, $3, $4)
	`,
		p.Title,
		p.Content,
		p.AuthorID,
		p.CreatedAt,
	)

	return err
}

func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), ` 
	UPDATE posts (title, content, author_id, created_at)
	SET  title = $1, content = $2, author_id = $3, created_at = $4
	WHERE id = $5
	`,
		p.Title,
		p.Content,
		p.AuthorID,
		p.CreatedAt,
		p.ID,
	)
	return err
}

func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
    DELETE FROM posts 
	WHERE id = $1
	`, p.ID)
	return err
}
