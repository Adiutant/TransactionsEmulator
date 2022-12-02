CREATE TABLE IF NOT EXISTS Users (
                                     user_name text NOT NULL PRIMARY KEY,
                                     password text,
                                     balance text,
                                     activity int
);
CREATE TABLE IF NOT EXISTS Transactions(
                                           id serial NOT NULL PRIMARY KEY,
                                           sender_user_name text,
                                           sender_balance text,
                                           sender_result_balance text,
                                           recipient_user_name text,
                                           recipient_balance text,
                                           recipient_result_balance text,
                                           amount text

);