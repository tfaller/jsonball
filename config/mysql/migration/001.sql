-- add encrypted flag
ALTER TABLE `document`
ADD COLUMN `encrypted` TINYINT(1) NOT NULL DEFAULT 0 AFTER `refreshedat`;

-- change data type of document itself.
-- this is because it is now probably encrypted
ALTER TABLE `document`
CHANGE COLUMN `document` `document` MEDIUMBLOB NOT NULL, LOCK = EXCLUSIVE;
