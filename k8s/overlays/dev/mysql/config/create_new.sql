CREATE DATABASE IF NOT EXISTS `caolila`;

USE `caolila`;

CREATE TABLE IF NOT EXISTS `user_info` (
    user_id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(36) NOT NULL,
    last_name VARCHAR(36) NOT NULL,
    mail_address1 VARCHAR(128) NOT NULL UNIQUE,
    mail_address2 VARCHAR(128),
    mail_address3 VARCHAR(128),
    phone_num1 VARCHAR(128) NOT NULL,
    phone_num2 VARCHAR(128),
    phone_num3 VARCHAR(128),
    address1 VARCHAR(128),
    address2 VARCHAR(128),
    address3 VARCHAR(128),
    post_code INT,
    rank INT,
    regist_day TIMESTAMP,
    sex INT,
    birth_day TIMESTAMP
);

LOCK TABLES `user_info` WRITE;

INSERT IGNORE INTO user_info VALUES 
('51096bba-a077-2c73-828a-3c8d2e8c08c7', 'Ryo', 'Kiuchi', 'ryo29wx@gmail.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c8', 'Sample1', 'hgo', 'nodje@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c9', 'Sample2', 'nfann', 'nodje2@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c0', 'Sample3', 'nfann', 'nodje3@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c1', 'Sample4', 'nfann', 'nodje4@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c2', 'Sample5', 'nfann', 'nodje5@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c3', 'Sample6', 'nfann', 'nodje6@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c4', 'Sample7', 'nfann', 'nodje7@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c5', 'Sample8', 'nfann', 'nodje8@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10'),
('51096bba-a077-2c73-828a-3c8d2e8c08c6', 'Sample9', 'nfann', 'nodje9@sample.com', '', '', '07012345678', '07012345678', '', 'Japan', 'chiyoda-ku', '', 1234567, 0, '2020-01-01 10:10:10', 0, '1994-01-01 10:10:10');

UNLOCK TABLES;