-- +goose Up
-- +goose StatementBegin
create table if not exists mouse (
    id int primary key,
    mouseDistance int not null default 0,
    leftClick int not null default 0,
    rightClick int not null default 0
);

insert into mouse (id, mouseDistance, leftClick, rightClick) values (1, 0, 0, 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table mouse
-- +goose StatementEnd
