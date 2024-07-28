SET FOREIGN_KEY_CHECKS = 0;
SET AUTOCOMMIT = 0;

DROP TABLE IF EXISTS Tweets;
CREATE TABLE Tweets (
    Id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    UserId BIGINT,
    Content TEXT,
    PostDate TIMESTAMP,
    FOREIGN KEY (UserId) REFERENCES Users(id)
);

DROP TABLE IF EXISTS Users;
CREATE TABLE Users (
    Id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    Name VARCHAR(255),
    Email VARCHAR(255),
    Password VARCHAR(255) NOT NULL,
    CreatedDate TIMESTAMP NOT NULL
);

DROP TABLE IF EXISTS Likes;
CREATE TABLE Likes (
    UserId BIGINT,
    TweetId BIGINT,
    PRIMARY KEY (UserId, TweetId),
    FOREIGN KEY (UserId) REFERENCES Users(id),
    FOREIGN KEY (TweetId) REFERENCES Tweets(id)
);

DROP TABLE IF EXISTS Interestings;
CREATE TABLE Interestings (
    UserId BIGINT,
    TweetId BIGINT,
    PRIMARY KEY (UserId, TweetId),
    FOREIGN KEY (UserId) REFERENCES Users(id),
    FOREIGN KEY (TweetId) REFERENCES Tweets(id)
);

DROP TABLE IF EXISTS Favorites;
CREATE TABLE Favorites (
    UserId BIGINT,
    TweetId BIGINT,
    PRIMARY KEY (UserId, TweetId),
    FOREIGN KEY (UserId) REFERENCES Users(id),
    FOREIGN KEY (TweetId) REFERENCES Tweets(id)
);

DROP TABLE IF EXISTS Follows;
CREATE TABLE Follows (
     FollowerId BIGINT,
     FolloweeId BIGINT,
     PRIMARY KEY (FollowerId, FolloweeId),
     FOREIGN KEY (FollowerId) REFERENCES Users(id),
     FOREIGN KEY (FolloweeId) REFERENCES Users(id)
);

INSERT INTO Users (Name, Email, Password, CreatedDate) VALUES
(
   'Jack Dorsey',
   'jack_dorsey@email.com',
   'password',
   '2021-01-01 00:00:00'
),
(
   'Britney Spears',
   'its_brittney@email.com',
   'password',
   '2021-01-01 00:00:00'
),
(
   'Kevin Durant',
   'kd@email.com',
   'password',
   '2021-01-01 00:00:00'
),
(
   'Joe Biden',
   'pres@email.com',
   'usa',
   '2021-01-01 00:00:00'
),
(
   'Elon Musk',
   'elon@x.com',
   'p@ssw0rd',
   '2022-10-28 00:00:00'
),
(
   'LeBron James',
   'goat@notmj.com',
   'king',
   '2020-01-01 00:00:00'
);

INSERT INTO Tweets (UserId, Content, PostDate) VALUES
(
   (SELECT Id FROM Users WHERE Name = 'Jack Dorsey'),
   'just setting up my twttr',
   '2011-01-01 00:00:00'
),
(
   (SELECT Id FROM Users WHERE Name = 'Britney Spears'),
   'Does anyone thing global warming is a good thing? I love Lady Gaga. I think she is a really interesting artist',
   '2011-02-01 00:00:00'
),
(
   (SELECT Id FROM Users WHERE Name = 'Kevin Durant'),
   'I\'m watching the History channel in the club and I\'m wondering how do these people kno what\'s goin on on the sun
    ... ain\'t nobody ever been',
   '2010-07-30 00:00:00'
);

INSERT INTO Likes (UserId, TweetId) VALUES
(
    (SELECT Id FROM Users WHERE Name = 'Elon Musk'),
    (SELECT Id FROM Tweets WHERE Content = 'just setting up my twttr')
),
(
    (SELECT Id FROM Users WHERE Name = 'Jack Dorsey'),
    (SELECT Id FROM Tweets WHERE Content = 'Does anyone thing global warming is a good thing? I love Lady Gaga. I think she is a really interesting artist')
),
(
    (SELECT Id FROM Users WHERE Name = 'LeBron James'),
    (SELECT Id FROM Tweets WHERE PostDate = '2010-07-30 00:00:00')
);

INSERT INTO Interestings (UserId, TweetId) VALUES
(
   (SELECT Id FROM Users WHERE Name = 'LeBron James'),
   (SELECT Id FROM Tweets WHERE Content = 'just setting up my twttr')
),
(
   (SELECT Id FROM Users WHERE Name = 'Joe Biden'),
   (SELECT Id FROM Tweets WHERE Content = 'Does anyone thing global warming is a good thing? I love Lady Gaga. I think she is a really interesting artist')
),
(
   (SELECT Id FROM Users WHERE Name = 'LeBron James'),
   (SELECT Id FROM Tweets WHERE PostDate = '2010-07-30 00:00:00')
);

INSERT INTO Favorites (UserId, TweetId) VALUES
(
    (SELECT Id FROM Users WHERE Name = 'Elon Musk'),
    (SELECT Id FROM Tweets WHERE Content = 'just setting up my twttr')
),
(
    (SELECT Id FROM Users WHERE Name = 'Jack Dorsey'),
    (SELECT Id FROM Tweets WHERE Content = 'Does anyone thing global warming is a good thing? I love Lady Gaga. I think she is a really interesting artist')
),
(
    (SELECT Id FROM Users WHERE Name = 'LeBron James'),
    (SELECT Id FROM Tweets WHERE PostDate = '2010-07-30 00:00:00')
);

INSERT INTO Follows (FollowerId, FolloweeId) VALUES
 (
     (SELECT Id FROM Users WHERE Name = 'Elon Musk'),
     (SELECT Id FROM Users WHERE Name = 'Jack Dorsey')
 ),
 (
     (SELECT Id FROM Users WHERE Name = 'Jack Dorsey'),
     (SELECT Id FROM Users WHERE Name = 'Britney Spears')
 ),
 (
     (SELECT Id FROM Users WHERE Name = 'Britney Spears'),
     (SELECT Id FROM Users WHERE Name = 'Kevin Durant')
 ),
 (
     (SELECT Id FROM Users WHERE Name = 'Kevin Durant'),
     (SELECT Id FROM Users WHERE Name = 'Joe Biden')
 ),
 (
     (SELECT Id FROM Users WHERE Name = 'Joe Biden'),
     (SELECT Id FROM Users WHERE Name = 'Elon Musk')
 ),
 (
     (SELECT Id FROM Users WHERE Name = 'Elon Musk'),
     (SELECT Id FROM Users WHERE Name = 'LeBron James')
 ),
 (
     (SELECT Id FROM Users WHERE Name = 'LeBron James'),
     (SELECT Id FROM Users WHERE Name = 'Jack Dorsey')
 );

SET FOREIGN_KEY_CHECKS = 1;
COMMIT;
