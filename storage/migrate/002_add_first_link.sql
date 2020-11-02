-- +goose Up

insert into links (user_id, url) values (1, 'https://t.me');

-- +goose Down

delete from links where id=1;
