-- CREATE DATABASE Netspace;

CREATE TABLE Locations (
	id SERIAL PRIMARY KEY,
	slug VARCHAR(55) UNIQUE NOT NULL,
	name VARCHAR(100) UNIQUE NOT NULL,
	isActive BOOL DEFAULT FALSE
);

CREATE TABLE Users (
	id TEXT UNIQUE NOT NULL,
	locationId INT REFERENCES Locations(id),
	name VARCHAR(30),
	slug VARCHAR(30),
	age INT,
	gender VARCHAR(6)
);

CREATE TABLE Interests (
	id SERIAL PRIMARY KEY,
	emoji TEXT NOT NULL,
	label TEXT NOT NULL
);

CREATE TABLE UserInterests (
	userId TEXT REFERENCES Users(id),
	interestId INT REFERENCES Interests(id),
	PRIMARY KEY(userId, interestId)
);

CREATE TABLE UserCustomInterests (
	userId TEXT REFERENCES Users(id),
	emoji TEXT NOT NULL,
	label TEXT NOT NULL,
	PRIMARY KEY(userId, label)
);

CREATE TABLE PrivateMessage (
	messageid TEXT PRIMARY KEY NOT NULL,
	locationId INT REFERENCES Locations(id),
	senderId TEXT REFERENCES Users(id),
	recipientId TEXT REFERENCES Users(id),
	"message" TEXT NOT NULL,
	"timestamp" TIMESTAMPTZ
);

CREATE TABLE PublicMessage (
	messageId  TEXT PRIMARY KEY NOT NULL,
	locationId INT REFERENCES Locations(id),
	senderId   TEXT REFERENCES Users(id),
	"message" TEXT NOT NULL,
	"timestamp" TIMESTAMPTZ
);

INSERT INTO Interests (emoji, label)
VALUES ('☕', 'Kopi'),
  ('🎮', 'Gaming'),
  ('📚', 'Buku'),
  ('🎵', 'Musik'),
  ('🍜', 'Kuliner'),
  ('✈️', 'Travel'),
  ('💻', 'Tech'),
  ('🎨', 'Seni'),
  ('🏋️', 'Olahraga'),
  ('🎬', 'Film'),
  ('📷', 'Fotografi'),
  ('🌱', 'Tanaman');

SELECT * FROM Users;