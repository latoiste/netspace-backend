-- CREATE DATABASE Netspace;

CREATE TABLE Locations (
	id SERIAL PRIMARY KEY,
	slug VARCHAR(55) UNIQUE NOT NULL,
	name VARCHAR(100) UNIQUE NOT NULL,
	address TEXT UNIQUE NOT NULL,
	partnerId TEXT NOT NULL,
	joinDate TIMESTAMPTZ DEFAULT NOW(),
	capacity INT CHECK (capacity >= 1),
	timezone TEXT DEFAULT 'Asia/Jakarta',
	isActive BOOL DEFAULT TRUE,
	qrToken TEXT NOT NULL,
	qrLabel TEXT NOT NULL
);

CREATE TABLE Users (
	id TEXT UNIQUE NOT NULL,
	locationId INT REFERENCES Locations(id),
	name VARCHAR(30),
	slug VARCHAR(30),
	age INT,
	gender VARCHAR(6),
	createdAt TIMESTAMPTZ DEFAULT NOW(),
	isActive BOOL DEFAULT TRUE
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

CREATE TABLE PrivateMessages (
	messageid TEXT PRIMARY KEY NOT NULL,
	locationId INT REFERENCES Locations(id),
	senderId TEXT REFERENCES Users(id),
	recipientId TEXT REFERENCES Users(id),
	"message" TEXT NOT NULL,
	"timestamp" TIMESTAMPTZ DEFAULT NOW
);

CREATE TABLE PublicMessages (
	messageId  TEXT PRIMARY KEY NOT NULL,
	locationId INT REFERENCES Locations(id),
	senderId   TEXT REFERENCES Users(id),
	"message" TEXT NOT NULL,
	"timestamp" TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE Group (
	Id TEXT PRIMARY KEY NOT NULL,
	Name TEXT,
	IsActive BOOL DEFAULT TRUE	
);

CREATE TABLE GroupMessages (
	messageId  TEXT PRIMARY KEY NOT NULL,
	locationId INT REFERENCES Locations(id),
	senderId   TEXT REFERENCES Users(id),
	groupId TEXT REFERENCES Group(id),
	"message" TEXT NOT NULL,
	"timestamp" TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE Notifications (
	"id" TEXT PRIMARY KEY,
	userId TEXT REFERENCES Users(id),
	"type" TEXT CHECK ("type" IN ('message', 'group_invite', 'chat_request', 'system')),
	emoji TEXT NOT NULL,
	avatarGradient TEXT NOT NULL,
	title TEXT NOT NULL,
	description TEXT NOT NULL,    
	"timestamp" TIMESTAMPTZ DEFAULT NOW(),
	unread BOOL DEFAULT TRUE,
	primaryLabel TEXT,
	secondaryLabel TEXT
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

SELECT * FROM publicmessage;