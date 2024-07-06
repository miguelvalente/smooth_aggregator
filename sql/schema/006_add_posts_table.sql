-- +goose Up
create table posts(
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title varchar(255),
    url varchar(255) not null unique,
    description varchar(255),
    published_at timestamp,
    feed_id uuid not null,
    constraint fk_feed_id
        foreign  key(feed_id)
        references feeds(id)
        on delete cascade
);


-- +goose Down
drop table posts;
