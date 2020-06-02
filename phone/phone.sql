CREATE TABLE `phone` (
  `phone_id` int NOT NULL AUTO_INCREMENT,
  `number` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`phone_id`),
  UNIQUE KEY `idx_unique_number` (`number`)
) ENGINE=InnoDB;

INSERT INTO phone (number) VALUES ("1234567890"), ("123 456 7891"), ("(123) 456 7892"), ("(123) 456-7893"), ("123-456-7894"), ("123-456-7890"), ("1234567892"), ("(123)456-7892");
