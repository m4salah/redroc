alter table images drop CONSTRAINT user_id_fk;
alter table images drop column user_id;
drop table users;