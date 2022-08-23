if object_id('[source_message]', 'U') is not null
drop table [source_message];

if object_id('[message]', 'U') is not null
drop table [message];

CREATE TABLE [source_message]
(
    [id]              integer IDENTITY PRIMARY KEY,
    [category]        varchar(255),
    [message]         text,
    [localize_config] text
    );

CREATE TABLE [message]
(
    [id]              integer NOT NULL,
    [language]        varchar(16) NOT NULL,
    [translation]     text,
    [localize_config] text
    );

ALTER TABLE [message] ADD CONSTRAINT [pk_message_id_language] PRIMARY KEY ([id], [language]);
ALTER TABLE [message] ADD CONSTRAINT [fk_message_source_message] FOREIGN KEY ([id]) REFERENCES [source_message] ([id]) ON UPDATE NO ACTION ON DELETE CASCADE;

CREATE INDEX [idx_message_language] on [message] ([language]);
CREATE INDEX [idx_source_message_category] on [source_message] ([category]);