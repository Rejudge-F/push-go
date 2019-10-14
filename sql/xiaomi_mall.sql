DROP DATABASE IF EXISTS xiaomi_mall;
CREATE DATABASE xiaomi_mall;

USE xiaomi_mall;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` INT NOT NULL ,
    `name` VARCHAR(10),
    PRIMARY KEY (`id`)
);

drop PROCEDURE if exists test_insert;
delimiter //
CREATE PROCEDURE test_insert(n int)
    begin
        declare v int default 0;
        SET AUTOCOMMIT=0;
        while v < n
        do
            insert into user
            values (v, "99");
            set v = v + 1;
            if mod(v, 1000) = 0 then commit;
            end if;
        end while;
        SET AUTOCOMMIT=1;
    end //

CALL test_insert(20000000);