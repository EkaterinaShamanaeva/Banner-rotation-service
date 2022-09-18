-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banners (
    id UUID NOT NULL PRIMARY KEY,
    description text
);

CREATE TABLE IF NOT EXISTS slots (
    id UUID NOT NULL PRIMARY KEY,
    description text
);

CREATE TABLE IF NOT EXISTS user_groups (
     id UUID NOT NULL PRIMARY KEY,
     description text
);

CREATE TABLE IF NOT EXISTS banner_slot (
     slot_id UUID NOT NULL,
     banner_id UUID NOT NULL,
     PRIMARY KEY (slot_id, banner_id)
);

CREATE TABLE IF NOT EXISTS banner_statistic (
     slot_id UUID NOT NULL,
     banner_id UUID NOT NULL,
     user_group_id UUID NOT NULL,
     click_count integer,
     show_count integer,
     PRIMARY KEY (banner_id, user_group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS slots;
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS banner_slot;
DROP TABLE IF EXISTS banner_statistic;
-- +goose StatementEnd
