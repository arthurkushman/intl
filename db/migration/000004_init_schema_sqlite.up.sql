drop table if exists `source_message`;
drop table if exists `message`;

CREATE TABLE `source_message`
(
    `id`          integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    `category`    varchar(255),
    `message`     text,
    `localize_config` text
);

CREATE TABLE `message`
(
    `id`              integer NOT NULL REFERENCES `source_message` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE,
    `language`        varchar(16) NOT NULL,
    `translation`     text,
    `localize_config` text,
    PRIMARY KEY (`id`, `language`)
);

CREATE INDEX idx_message_language ON message (language);
CREATE INDEX idx_source_message_category ON source_message (category);