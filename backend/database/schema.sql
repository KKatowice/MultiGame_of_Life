CREATE TABLE `users` (
	`userId`	INTEGER PRIMARY KEY AUTO_INCREMENT,
	`name`	VARCHAR(10) UNIQUE
);
CREATE TABLE `games` (
	`roomId`	INTEGER PRIMARY KEY AUTO_INCREMENT,
	`hei`	INTEGER,
	`wid`	INTEGER
);
CREATE TABLE `lobby` (
	`userId`	INTEGER,
	`roomId`	INTEGER,
	`id`	INTEGER PRIMARY KEY AUTO_INCREMENT,
	CONSTRAINT `fk_room` FOREIGN KEY(`roomId`) REFERENCES `games`(`roomId`) ON DELETE CASCADE,
	CONSTRAINT `fk_user` FOREIGN KEY(`userId`) REFERENCES `users`(`userId`) ON DELETE CASCADE
);
