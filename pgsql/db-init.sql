CREATE DATABASE mydb;

\connect mydb;

CREATE SCHEMA manufacture;

CREATE TABLE manufacture.manufacturer
(
    id varchar,
    details jsonb
);

CREATE INDEX details_manufacturer_need_update_idx ON manufacture.manufacturer USING gin ((details -> 'needUpdate'));

INSERT INTO manufacture.manufacturer
VALUES ('asd1asd', '{"name": "Вася", "needUpdate": true}'),
       ('asd2as', '{"name": "Петя", "needUpdate": true}'),
       ('asd4123', '{"name": "Лева", "needUpdate": true}'),
       ('11235', '{"name": "Антон", "needUpdate": true}'),
       ('26asd', '{"name": "Петр", "needUpdate": true}'),
       ('37asd', '{"name": "Владимир", "needUpdate": true}'),
       ('48asda', '{"name": "Вовчик", "needUpdate": true}'),
       ('asdasd', '{"name": "Леня", "needUpdate": true}'),
       ('s2zxczc', '{"name": "Олег", "needUpdate": true}'),
       ('czxc2', '{"name": "Евгений", "needUpdate": true}'),
       ('adxcgn2', '{"name": "Курт Кобейн", "needUpdate": true}'),
       ('df2dfgd', '{"name": "Омар Одом", "needUpdate": true}'),
       ('dfg2gasdf', '{"name": "Леброн Джеймс", "needUpdate": true}'),
       ('dfg2dasx', '{"name": "Ники Минаж", "needUpdate": true}'),
       ('dfg2dfg', '{"name": "Рыбак", "needUpdate": true}'),
       ('dfg2rytjk', '{"name": "Филимон", "needUpdate": true}'),
       ('dfg2dg', '{"name": "Златовласка", "needUpdate": true}'),
       ('dfg2dfg', '{"name": "Андрей", "needUpdate": true}'),
       ('dfgdgc2', '{"name": "Анджей", "needUpdate": true}'),
       ('2cvbuy', '{"name": "Колян", "needUpdate": true}'),
       ('yui2', '{"name": "Димон", "needUpdate": true}'),
       ('2iuyi', '{"name": "Жека", "needUpdate": true}'),
       ('dfgertru2', '{"name": "Федор", "needUpdate": true}'),
       ('qwrtfhgm2', '{"name": "Гена", "needUpdate": true}'),
       ('cxchght2', '{"name": "Леха", "needUpdate": true}'),
       ('zx2qwe', '{"name": "Данил", "needUpdate": true}'),
       ('asdiu2', '{"name": "Илья", "needUpdate": true}'),
       ('oipuokh2', '{"name": "Бенедикт", "needUpdate": true}'),
       ('werthbv2', '{"name": "Кевин", "needUpdate": true}'),
       ('xcvcgu2', '{"name": "Амар", "needUpdate": true}'),
       ('zdrew2', '{"name": "Арман", "needUpdate": true}'),
       ('uyjnb2', '{"name": "Клинтон", "needUpdate": true}');
