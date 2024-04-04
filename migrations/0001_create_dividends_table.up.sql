CREATE TABLE IF NOT EXISTS `dividends`
(
    `market`            varchar(10)    NOT NULL,
    `isin`              varchar(100)   NOT NULL,
    `name`              varchar(100)   NOT NULL,

    `value`             decimal(10, 4) NOT NULL,
    `yield`             decimal(10, 4),

    `ex_date`           date           NOT NULL,
    `pay_date`          date           NOT NULL,
    `registration_date` date           NOT NULL,

    PRIMARY KEY (`isin`, `ex_date`)
);
