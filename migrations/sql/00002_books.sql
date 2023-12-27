-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS books(
	"book_id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	"title" TEXT NOT NULL,
	"created" TIMESTAMPTZ DEFAULT NOW(),
	"updated" TIMESTAMPTZ,
	"deleted" TIMESTAMPTZ
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON books
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE books;